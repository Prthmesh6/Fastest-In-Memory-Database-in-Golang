package main

import (
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestDistributedCache_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Find available ports
	ports := []string{"16379", "16380", "16381"}
	addrs := []string{
		"127.0.0.1:16379",
		"127.0.0.1:16380",
		"127.0.0.1:16381",
	}

	// Check if ports are available
	for _, port := range ports {
		conn, err := net.Listen("tcp", ":"+port)
		if err != nil {
			t.Skipf("Port %s not available, skipping integration test", port)
			return
		}
		conn.Close()
	}

	peersCSV := strings.Join(addrs, ",")
	
	// Start 3 nodes
	cmds := make([]*exec.Cmd, 3)
	aofFiles := []string{"test-node-16379.aof", "test-node-16380.aof", "test-node-16381.aof"}
	
	// Clean up AOF files
	for _, f := range aofFiles {
		os.Remove(f)
	}
	defer func() {
		for _, f := range aofFiles {
			os.Remove(f)
		}
	}()

	for i := 0; i < 3; i++ {
		cmd := exec.Command("go", "run", ".", 
			"-addr", addrs[i],
			"-peers", peersCSV,
			"-aof", aofFiles[i])
		cmd.Stdout = nil // Suppress output for cleaner test runs
		cmd.Stderr = nil
		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start node %d: %v", i, err)
		}
		cmds[i] = cmd
		defer func(idx int) {
			if cmds[idx].Process != nil {
				cmds[idx].Process.Kill()
				cmds[idx].Wait() // Wait for process to exit
			}
		}(i)
	}

	// Wait for servers to start
	time.Sleep(2 * time.Second)

	// Test: Connect to node 1 and SET a key
	conn1, err := net.Dial("tcp", addrs[0])
	if err != nil {
		t.Fatalf("Failed to connect to node 1: %v", err)
	}
	defer conn1.Close()

	// SET key1=value1
	setCmd := Value{
		typ: "array",
		array: []Value{
			{typ: "bulk", bulk: "SET"},
			{typ: "bulk", bulk: "key1"},
			{typ: "bulk", bulk: "value1"},
		},
	}
	writer1 := NewWriter(conn1)
	if err := writer1.Write(setCmd); err != nil {
		t.Fatalf("Failed to write SET command: %v", err)
	}

	resp1 := NewResp(conn1)
	result1, err := resp1.Read()
	if err != nil {
		t.Fatalf("Failed to read SET response: %v", err)
	}
	if result1.typ != "string" || result1.str != "OK" {
		t.Errorf("SET failed, got %+v", result1)
	}

	// Test: GET from node 2 (should proxy to owner)
	conn2, err := net.Dial("tcp", addrs[1])
	if err != nil {
		t.Fatalf("Failed to connect to node 2: %v", err)
	}
	defer conn2.Close()

	getCmd := Value{
		typ: "array",
		array: []Value{
			{typ: "bulk", bulk: "GET"},
			{typ: "bulk", bulk: "key1"},
		},
	}
	writer2 := NewWriter(conn2)
	if err := writer2.Write(getCmd); err != nil {
		t.Fatalf("Failed to write GET command: %v", err)
	}

	resp2 := NewResp(conn2)
	result2, err := resp2.Read()
	if err != nil {
		t.Fatalf("Failed to read GET response: %v", err)
	}
	if result2.typ != "bulk" || result2.bulk != "value1" {
		t.Errorf("GET failed, expected 'value1', got %+v", result2)
	}

	// Test: GET from node 3 (should also proxy)
	conn3, err := net.Dial("tcp", addrs[2])
	if err != nil {
		t.Fatalf("Failed to connect to node 3: %v", err)
	}
	defer conn3.Close()

	writer3 := NewWriter(conn3)
	if err := writer3.Write(getCmd); err != nil {
		t.Fatalf("Failed to write GET command to node 3: %v", err)
	}

	resp3 := NewResp(conn3)
	result3, err := resp3.Read()
	if err != nil {
		t.Fatalf("Failed to read GET response from node 3: %v", err)
	}
	if result3.typ != "bulk" || result3.bulk != "value1" {
		t.Errorf("GET from node 3 failed, expected 'value1', got %+v", result3)
	}

	// Test: NODES command
	nodesCmd := Value{
		typ: "array",
		array: []Value{
			{typ: "bulk", bulk: "NODES"},
		},
	}
	if err := writer2.Write(nodesCmd); err != nil {
		t.Fatalf("Failed to write NODES command: %v", err)
	}

	resp4 := NewResp(conn2)
	result4, err := resp4.Read()
	if err != nil {
		t.Fatalf("Failed to read NODES response: %v", err)
	}
	if result4.typ != "array" || len(result4.array) != 3 {
		t.Errorf("NODES failed, expected 3 nodes, got %+v", result4)
	}

	// Test: Hash operations with proxying
	conn4, err := net.Dial("tcp", addrs[0])
	if err != nil {
		t.Fatalf("Failed to connect to node 1 for hash test: %v", err)
	}
	defer conn4.Close()

	// HSET on node 1
	hsetCmd := Value{
		typ: "array",
		array: []Value{
			{typ: "bulk", bulk: "HSET"},
			{typ: "bulk", bulk: "myhash"},
			{typ: "bulk", bulk: "field1"},
			{typ: "bulk", bulk: "value1"},
		},
	}
	writer4 := NewWriter(conn4)
	if err := writer4.Write(hsetCmd); err != nil {
		t.Fatalf("Failed to write HSET command: %v", err)
	}

	resp5 := NewResp(conn4)
	result5, err := resp5.Read()
	if err != nil {
		t.Fatalf("Failed to read HSET response: %v", err)
	}
	if result5.typ != "string" || result5.str != "OK" {
		t.Errorf("HSET failed, got %+v", result5)
	}

	// HGET from node 2 (should proxy)
	conn5, err := net.Dial("tcp", addrs[1])
	if err != nil {
		t.Fatalf("Failed to connect to node 2 for hash test: %v", err)
	}
	defer conn5.Close()

	hgetCmd := Value{
		typ: "array",
		array: []Value{
			{typ: "bulk", bulk: "HGET"},
			{typ: "bulk", bulk: "myhash"},
			{typ: "bulk", bulk: "field1"},
		},
	}
	writer5 := NewWriter(conn5)
	if err := writer5.Write(hgetCmd); err != nil {
		t.Fatalf("Failed to write HGET command: %v", err)
	}

	resp6 := NewResp(conn5)
	result6, err := resp6.Read()
	if err != nil {
		t.Fatalf("Failed to read HGET response: %v", err)
	}
	if result6.typ != "bulk" || result6.bulk != "value1" {
		t.Errorf("HGET failed, expected 'value1', got %+v", result6)
	}
}

package main

import (
	"testing"
)

func TestCluster_Enabled(t *testing.T) {
	tests := []struct {
		name   string
		self   string
		nodes  []string
		want   bool
	}{
		{"single node", "127.0.0.1:6379", []string{}, false},
		{"two nodes", "127.0.0.1:6379", []string{"127.0.0.1:6380"}, true},
		{"three nodes", "127.0.0.1:6379", []string{"127.0.0.1:6380", "127.0.0.1:6381"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCluster(tt.self, tt.nodes)
			if got := c.Enabled(); got != tt.want {
				t.Errorf("Cluster.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_Owner(t *testing.T) {
	c := NewCluster("127.0.0.1:6379", []string{"127.0.0.1:6380", "127.0.0.1:6381"})
	
	// Test that owner is deterministic
	key := "test-key"
	owner1 := c.Owner(key)
	owner2 := c.Owner(key)
	if owner1 != owner2 {
		t.Errorf("Owner() should be deterministic, got %v and %v", owner1, owner2)
	}
	
	// Test that owner is one of the nodes
	nodes := c.Nodes()
	found := false
	for _, n := range nodes {
		if n == owner1 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Owner() = %v, should be one of %v", owner1, nodes)
	}
	
	// Test that different keys can have different owners
	key2 := "test-key-2"
	owner2_key := c.Owner(key2)
	// At least one key should map to a different owner (with 3 nodes, this is likely)
	// But we can't guarantee it, so we just check it's valid
	found2 := false
	for _, n := range nodes {
		if n == owner2_key {
			found2 = true
			break
		}
	}
	if !found2 {
		t.Errorf("Owner() for key2 = %v, should be one of %v", owner2_key, nodes)
	}
}

func TestCluster_Nodes(t *testing.T) {
	c := NewCluster("127.0.0.1:6379", []string{"127.0.0.1:6380", "127.0.0.1:6381"})
	nodes := c.Nodes()
	
	if len(nodes) != 3 {
		t.Errorf("Nodes() length = %v, want 3", len(nodes))
	}
	
	// Test that nodes are sorted
	for i := 1; i < len(nodes); i++ {
		if nodes[i-1] > nodes[i] {
			t.Errorf("Nodes() should be sorted, got %v", nodes)
		}
	}
	
	// Test that self is included
	self := c.Self()
	found := false
	for _, n := range nodes {
		if n == self {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Self() = %v should be in Nodes() = %v", self, nodes)
	}
}

func TestCluster_Self(t *testing.T) {
	c := NewCluster("127.0.0.1:6379", []string{"127.0.0.1:6380"})
	self := c.Self()
	if self != "127.0.0.1:6379" {
		t.Errorf("Self() = %v, want 127.0.0.1:6379", self)
	}
}

func TestNormalizeDialAddr(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{"with host", "127.0.0.1:6379", "127.0.0.1:6379"},
		{"bare port", ":6379", "127.0.0.1:6379"},
		{"empty host", ":6380", "127.0.0.1:6380"},
		{"localhost explicit", "localhost:6379", "localhost:6379"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeDialAddr(tt.addr); got != tt.want {
				t.Errorf("normalizeDialAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestKey(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []Value
		wantKey string
		wantOk  bool
	}{
		{"GET with key", "GET", []Value{{typ: "bulk", bulk: "mykey"}}, "mykey", true},
		{"SET with key", "SET", []Value{{typ: "bulk", bulk: "mykey"}, {typ: "bulk", bulk: "value"}}, "mykey", true},
		{"HSET with hash", "HSET", []Value{{typ: "bulk", bulk: "myhash"}, {typ: "bulk", bulk: "field"}, {typ: "bulk", bulk: "value"}}, "myhash", true},
		{"HGET with hash", "HGET", []Value{{typ: "bulk", bulk: "myhash"}, {typ: "bulk", bulk: "field"}}, "myhash", true},
		{"HGETALL with hash", "HGETALL", []Value{{typ: "bulk", bulk: "myhash"}}, "myhash", true},
		{"PING no key", "PING", []Value{}, "", false},
		{"GET no args", "GET", []Value{}, "", false},
		{"SET no args", "SET", []Value{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotOk := requestKey(tt.command, tt.args)
			if gotKey != tt.wantKey {
				t.Errorf("requestKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotOk != tt.wantOk {
				t.Errorf("requestKey() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestCluster_Deduplication(t *testing.T) {
	// Test that duplicate nodes are removed
	c := NewCluster("127.0.0.1:6379", []string{"127.0.0.1:6380", "127.0.0.1:6379", "127.0.0.1:6380"})
	nodes := c.Nodes()
	if len(nodes) != 2 {
		t.Errorf("Nodes() should deduplicate, got %v nodes, want 2", len(nodes))
	}
}

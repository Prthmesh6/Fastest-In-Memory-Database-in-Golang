package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {
	var listenAddr string
	var peersCSV string
	var aofPath string
	flag.StringVar(&listenAddr, "addr", ":6379", "listen address")
	flag.StringVar(&peersCSV, "peers", "", "comma-separated peer addresses (include self) to enable clustering")
	flag.StringVar(&aofPath, "aof", "database.aof", "append-only file path")
	flag.Parse()

	selfAddr := normalizeDialAddr(listenAddr)
	var peers []string
	if peersCSV != "" {
		peers = strings.Split(peersCSV, ",")
	}
	cluster := NewCluster(selfAddr, peers)

	fmt.Printf("Listening on %s\n", listenAddr)
	if cluster.Enabled() {
		fmt.Printf("Cluster enabled. Self=%s Nodes=%v\n", cluster.Self(), cluster.Nodes())
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof(aofPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn, cluster, aof)
	}
}

func handleConn(conn net.Conn, cluster *Cluster, aof *Aof) {
	defer conn.Close()

	resp := NewResp(conn)
	writer := NewWriter(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			return
		}

		result := processRequest(cluster, aof, value)
		_ = writer.Write(result)
	}
}

func processRequest(cluster *Cluster, aof *Aof, value Value) Value {
	if value.typ != "array" {
		return Value{typ: "error", str: "ERR invalid request, expected array"}
	}
	if len(value.array) == 0 {
		return Value{typ: "error", str: "ERR invalid request, expected array length > 0"}
	}

	command := strings.ToUpper(value.array[0].bulk)
	args := value.array[1:]

	if command == "NODES" {
		nodes := cluster.Nodes()
		out := make([]Value, 0, len(nodes))
		for _, n := range nodes {
			out = append(out, Value{typ: "bulk", bulk: n})
		}
		return Value{typ: "array", array: out}
	}

	if key, ok := requestKey(command, args); ok && cluster.Enabled() {
		owner := cluster.Owner(key)
		if owner != cluster.Self() {
			v, err := proxyRequest(owner, value)
			if err != nil {
				return proxyError(owner, err)
			}
			return v
		}
	}

	handler, ok := Handlers[command]
	if !ok {
		return Value{typ: "error", str: "ERR unknown command"}
	}

	if command == "SET" || command == "HSET" {
		_ = aof.Write(value)
	}

	return handler(args)
}

func requestKey(command string, args []Value) (string, bool) {
	switch command {
	case "GET", "SET":
		if len(args) < 1 {
			return "", false
		}
		return args[0].bulk, true
	case "HSET", "HGET", "HGETALL":
		if len(args) < 1 {
			return "", false
		}
		return args[0].bulk, true // hash name
	default:
		return "", false
	}
}

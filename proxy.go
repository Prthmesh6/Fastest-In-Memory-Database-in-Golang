package main

import (
	"fmt"
	"net"
	"time"
)

func proxyRequest(ownerAddr string, req Value) (Value, error) {
	d := net.Dialer{Timeout: 2 * time.Second}
	conn, err := d.Dial("tcp", ownerAddr)
	if err != nil {
		return Value{}, err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return Value{}, err
	}

	if _, err := conn.Write(req.Marshal()); err != nil {
		return Value{}, err
	}

	resp := NewResp(conn)
	v, err := resp.Read()
	if err != nil {
		return Value{}, err
	}
	return v, nil
}

func proxyError(ownerAddr string, err error) Value {
	return Value{typ: "error", str: fmt.Sprintf("ERR proxy to %s failed: %v", ownerAddr, err)}
}

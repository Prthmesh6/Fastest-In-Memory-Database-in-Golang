package main

import (
	"testing"
)

func BenchmarkPing(b *testing.B) {
	// Test simple PING
	b.Run("Simple PING", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ping([]Value{})
		}
	})

	// Test PING with argument
	b.Run("PING with argument", func(b *testing.B) {
		args := []Value{{typ: "bulk", bulk: "hello"}}
		for i := 0; i < b.N; i++ {
			ping(args)
		}
	})
}

func BenchmarkSetGet(b *testing.B) {
	// Benchmark SET operation
	b.Run("SET", func(b *testing.B) {
		args := []Value{
			{typ: "bulk", bulk: "key"},
			{typ: "bulk", bulk: "value"},
		}
		for i := 0; i < b.N; i++ {
			set(args)
		}
	})

	// Benchmark GET operation
	b.Run("GET", func(b *testing.B) {
		args := []Value{{typ: "bulk", bulk: "key"}}
		for i := 0; i < b.N; i++ {
			get(args)
		}
	})
}

func BenchmarkHash(b *testing.B) {
	// Benchmark HSET operation
	b.Run("HSET", func(b *testing.B) {
		args := []Value{
			{typ: "bulk", bulk: "myhash"},
			{typ: "bulk", bulk: "field1"},
			{typ: "bulk", bulk: "value1"},
		}
		for i := 0; i < b.N; i++ {
			hset(args)
		}
	})

	// Benchmark HGET operation
	b.Run("HGET", func(b *testing.B) {
		args := []Value{
			{typ: "bulk", bulk: "myhash"},
			{typ: "bulk", bulk: "field1"},
		}
		for i := 0; i < b.N; i++ {
			hget(args)
		}
	})

	// Benchmark HGETALL operation
	b.Run("HGETALL", func(b *testing.B) {
		// First set up some data
		hset([]Value{
			{typ: "bulk", bulk: "myhash"},
			{typ: "bulk", bulk: "field1"},
			{typ: "bulk", bulk: "value1"},
		})
		hset([]Value{
			{typ: "bulk", bulk: "myhash"},
			{typ: "bulk", bulk: "field2"},
			{typ: "bulk", bulk: "value2"},
		})

		args := []Value{{typ: "bulk", bulk: "myhash"}}
		b.ResetTimer() // Reset timer after setup

		for i := 0; i < b.N; i++ {
			hgetall(args)
		}
	})
}

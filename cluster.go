package main

import (
	"hash/fnv"
	"net"
	"sort"
	"strings"
)

type Cluster struct {
	self  string
	nodes []string
}

func NewCluster(self string, nodes []string) *Cluster {
	self = normalizeDialAddr(self)

	set := map[string]struct{}{}
	var out []string
	add := func(a string) {
		a = strings.TrimSpace(a)
		if a == "" {
			return
		}
		a = normalizeDialAddr(a)
		if _, ok := set[a]; ok {
			return
		}
		set[a] = struct{}{}
		out = append(out, a)
	}

	add(self)
	for _, n := range nodes {
		add(n)
	}

	sort.Strings(out)

	return &Cluster{self: self, nodes: out}
}

func (c *Cluster) Enabled() bool {
	return len(c.nodes) > 1
}

func (c *Cluster) Self() string {
	return c.self
}

func (c *Cluster) Nodes() []string {
	cp := make([]string, 0, len(c.nodes))
	cp = append(cp, c.nodes...)
	return cp
}

func (c *Cluster) Owner(key string) string {
	if !c.Enabled() {
		return c.self
	}

	var best string
	var bestScore uint64
	for _, n := range c.nodes {
		s := rendezvousScore(key, n)
		if best == "" || s > bestScore {
			best = n
			bestScore = s
		}
	}
	return best
}

func rendezvousScore(key, node string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(node))
	return h.Sum64()
}

func normalizeDialAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return addr
	}

	host, port, err := net.SplitHostPort(addr)
	if err == nil {
		if host == "" {
			host = "127.0.0.1"
		}
		return net.JoinHostPort(host, port)
	}

	// If someone passes ":6379" without brackets/host parsing working in some edge cases,
	// fall back to prefixing localhost when it looks like a bare port.
	if strings.HasPrefix(addr, ":") && len(addr) > 1 {
		return "127.0.0.1" + addr
	}

	return addr
}

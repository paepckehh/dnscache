package main

import (
	"net"
	"testing"
	"time"

	"paepcke.de/dnscache"
)

const iterations = 250

var domains = []string{
	"bbc.co.uk",
	"codeberg.org",
	"www.ccc.de",
	"paepcke.de",
	"sslmate.com",
	"git.kernel.org",
	"www.github.com",
}

// Benchmark_stdlib ...
func Benchmark_stdlib(b *testing.B) {
	stdlib() // warm upstream cache & let it settle
	time.Sleep(2 * time.Second)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for ii := 0; ii < iterations; ii++ {
			stdlib()
		}
	}
}

// Benchmark_dnscache ...
func Benchmark_dnscache(b *testing.B) {
	cached() // warm upstream cache & let it settle
	time.Sleep(2 * time.Second)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for ii := 0; ii < iterations; ii++ {
			cached()
		}
	}
}

// simulate mixed dns payload via stdlib resolver
func stdlib() {
	for _, domain := range domains {
		_, _ = net.LookupIP(domain)
		_, _ = dnscache.LookupCNAME(domain)
		_, _ = net.LookupIP(domain)
	}
}

// simulate mixed dns payload via dnscache resolver
func cached() {
	for _, domain := range domains {
		_, _ = dnscache.LookupIP(domain)
		_, _ = dnscache.LookupCNAME(domain)
		_, _ = dnscache.LookupIP(domain)
	}
}

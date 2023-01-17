# OVERVIEW
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/dnscache.svg)](https://pkg.go.dev/paepcke.de/dnscache) [![Go Report Card](https://goreportcard.com/badge/paepcke.de/dnscache)](https://goreportcard.com/report/paepcke.de/dnscache) [![Go Build](https://github.com/paepckehh/dnscache/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/dnscache/actions/workflows/golang.yml)

[paepche.de/dnscache](https://paepcke.de/dnscache/)

- 100 % api compatible with stdlib dns net package (100% api coverage for simple api) just plug&play
- thread safe, memory efficient, low-latency, ignore (bungled) dns ttl and enforce caching for 24 hours
- less than 350 LOC, 100 % pure golang, stdlib only, external dependency free, easy to use
- see api.go for details

# SHOWTIME  

## default golang resolver vs. dnscache ( latency / alloc )

``` Shell
goos: freebsd
goarch: arm
pkg: paepcke.de/dnsache/cmd/dnsbench
Benchmark_stdlib-4     	       1	10403693854 ns/op	21246312 B/op	  245988 allocs/op
Benchmark_dnscache-4   	     531	    2247311 ns/op	       0 B/op	       0 allocs/op
PASS
```

# DOCS

[pkg.go.dev/paepcke.de/dnscache](https://pkg.go.dev/paepcke.de/dnscache)

# CONTRIBUTION

Yes, Please! PRs Welcome! 

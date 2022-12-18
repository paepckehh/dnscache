// package dnscache ...
package dnscache

import (
	"net"
	"runtime"
	"sync"
	"time"
)

const (
	_empty       = ""
	_sep         = "#"
	_mx          = "MX"
	_ns          = "NS"
	_srv         = "SRV"
	_txt         = "TXT"
	_ptr         = "PTR"
	_cname       = "CNAME"
	_minTTL      = 24 * 60 * 60     // force cache dns entries for one day, ignore (mostly bungled) ttl
	_expireEvery = 60 * time.Minute // run expire process ever 20 minutes
)

var (
	scale = runtime.NumCPU() * 10
	// IP Cache
	bgIP                       sync.WaitGroup
	dnsIPMap                   = make(map[string][]net.IP, scale)
	dnsIPExpire                = make(map[string]int64, scale)
	dnsIPChan                  = make(chan dnsIP, scale)
	dnsIPChanExpire            = make(chan dnsIP, scale*25)
	dnsIPLock, dnsIPLockExpire sync.Mutex
	// Host Cache
	bgHost                         sync.WaitGroup
	dnsHostMap                     = make(map[string]any, scale)
	dnsHostExpire                  = make(map[string]int64, scale)
	dnsHostChan                    = make(chan dnsHost, scale)
	dnsHostChanExpire              = make(chan dnsHost, scale*25)
	dnsHostLock, dnsHostLockExpire sync.Mutex
)

// dnsIP ...
type dnsIP struct {
	host string
	ip   []net.IP
}

// dnsHost ...
type dnsHost struct {
	host  string
	addrs any
}

// getIP ...
func getIP(host string) ([]net.IP, bool) {
	dnsIPLock.Lock()
	if ip, ok := dnsIPMap[host]; ok {
		dnsIPLock.Unlock()
		return ip, true
	}
	dnsIPLock.Unlock()
	return []net.IP{}, false
}

// getHost ...
func getHost(host string) (any, bool) {
	dnsHostLock.Lock()
	if addrs, ok := dnsHostMap[host]; ok {
		dnsHostLock.Unlock()
		return addrs, true
	}
	dnsHostLock.Unlock()
	return nil, false
}

// intit ...
func init() {
	spinUpCacheWriter()
	spinUpExpireWriter()
	spinUpExpireGC()
}

// spinUpCacheWriter ...
func spinUpCacheWriter() {
	go func() {
		for c := range dnsIPChan {
			dnsIPLock.Lock()
			dnsIPMap[c.host] = c.ip
			dnsIPLock.Unlock()
			dnsIPChanExpire <- c
		}
	}()
	go func() {
		for c := range dnsHostChan {
			dnsHostLock.Lock()
			dnsHostMap[c.host] = c.addrs
			dnsHostLock.Unlock()
			dnsHostChanExpire <- c
		}
	}()
}

// spinUpExpireWriter ...
func spinUpExpireWriter() {
	go func() {
		for c := range dnsIPChanExpire {
			dnsIPLockExpire.Lock()
			dnsIPExpire[c.host] = time.Now().Unix() + _minTTL
			dnsIPLockExpire.Unlock()
		}
	}()
	go func() {
		for c := range dnsHostChanExpire {
			dnsHostLockExpire.Lock()
			dnsHostExpire[c.host] = time.Now().Unix() + _minTTL
			dnsHostLockExpire.Unlock()
		}
	}()
}

// spinUpExpireGC ...
func spinUpExpireGC() {
	go func() {
		for {
			time.Sleep(_expireEvery)
			now := time.Now().Unix()
			{
				expired := []string{}
				dnsIPLockExpire.Lock()
				for k, v := range dnsIPExpire {
					if v < now {
						expired = append(expired, k)
					}
				}
				dnsIPLockExpire.Unlock()
				dnsIPLock.Lock()
				for _, e := range expired {
					delete(dnsIPMap, e)
				}
				dnsIPLock.Unlock()
			}
			{
				expired := []string{}
				dnsHostLockExpire.Lock()
				for k, v := range dnsHostExpire {
					if v < now {
						expired = append(expired, k)
					}
				}
				dnsHostLockExpire.Unlock()
				dnsHostLock.Lock()
				for _, e := range expired {
					delete(dnsHostMap, e)
				}
				dnsHostLock.Unlock()
			}
		}
	}()
}

// cleanCacheAll ...
func cleanCacheAll() {
	cleanIPCache()
	cleanHostCache()
}

// cleanIPCache ...
func cleanIPCache() {
	bgIP.Add(2)
	// lookup table
	go func() {
		dnsIPLock.Lock()
		for host := range dnsIPMap {
			delete(dnsIPMap, host)
		}
		bgIP.Done()
	}()
	go func() { // expire table
		dnsIPLockExpire.Lock()
		for host := range dnsIPExpire {
			delete(dnsIPExpire, host)
		}
		dnsIPLockExpire.Unlock()
		bgIP.Done()
	}()
	bgIP.Wait()
	dnsIPLock.Unlock()
}

// cleanHostCache ...
func cleanHostCache() {
	bgHost.Add(2)
	go func() { // lookup table
		dnsHostLock.Lock()
		for host := range dnsHostMap {
			delete(dnsHostMap, host)
		}
		bgHost.Done()
	}()
	go func() { // expire table
		dnsHostLockExpire.Lock()
		for host := range dnsHostExpire {
			delete(dnsHostExpire, host)
		}
		bgHost.Done()
	}()
	dnsHostLockExpire.Unlock()
	bgHost.Wait()
	dnsHostLock.Unlock()
}

// Package clientip selects client addresses behind an explicitly trusted proxy.
package clientip

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
)

var trustedNet = func() netip.Prefix {
	s := os.Getenv("HGMX_TRUSTED_PROXY")
	p, err := parseProxy(s)
	if err != nil {
		panic(fmt.Sprintf("clientip: invalid HGMX_TRUSTED_PROXY %q: %v", s, err))
	}
	return p
}()

func parseProxy(s string) (netip.Prefix, error) {
	if s == "" {
		return netip.Prefix{}, nil
	}
	if !strings.Contains(s, "/") {
		ip, err := netip.ParseAddr(s)
		if err != nil {
			return netip.Prefix{}, err
		}
		ip = ip.Unmap()
		return netip.PrefixFrom(ip, ip.BitLen()), nil
	}
	p, err := netip.ParsePrefix(s)
	if err != nil {
		return netip.Prefix{}, err
	}
	if p.Addr().Is4In6() {
		if p.Bits() < 96 {
			return netip.Prefix{}, fmt.Errorf("IPv4-mapped prefix shorter than /96")
		}
		p = netip.PrefixFrom(p.Addr().Unmap(), p.Bits()-96)
	}
	return p.Masked(), nil
}

// From returns the RemoteAddr host unless a trusted proxy supplies a valid IP.
// It accepts CF-Connecting-IP only from peers inside HGMX_TRUSTED_PROXY, an IP or CIDR.
// That setting is read once at package initialization; an invalid value panics.
func From(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	peer, err := netip.ParseAddr(host)
	if err != nil || !trustedNet.IsValid() || !trustedNet.Contains(peer.Unmap()) {
		return host
	}
	ip, err := netip.ParseAddr(r.Header.Get("CF-Connecting-IP"))
	if err != nil {
		return host
	}
	return ip.Unmap().String()
}

// Package clientip selects client addresses behind an explicitly trusted proxy.
package clientip

import (
	"net"
	"net/http"
	"net/netip"
	"os"
)

var trustedProxy = func() netip.Addr {
	ip, _ := netip.ParseAddr(os.Getenv("HGMX_TRUSTED_PROXY"))
	return ip.Unmap()
}()

// From returns the RemoteAddr host unless a trusted proxy supplies a valid IP.
// It accepts CF-Connecting-IP only from the IP in HGMX_TRUSTED_PROXY.
// That setting is read once at package initialization.
func From(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	peer, err := netip.ParseAddr(host)
	if err != nil || !trustedProxy.IsValid() || peer.Unmap() != trustedProxy {
		return host
	}
	ip, err := netip.ParseAddr(r.Header.Get("CF-Connecting-IP"))
	if err != nil {
		return host
	}
	return ip.Unmap().String()
}

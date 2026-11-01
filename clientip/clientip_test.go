package clientip

import (
	"net/http/httptest"
	"net/netip"
	"testing"
)

func TestFrom(t *testing.T) {
	defer func(proxy netip.Addr) { trustedProxy = proxy }(trustedProxy)
	for _, tt := range []struct {
		name, proxy, remote, header, want string
	}{
		{"spoof", "192.0.2.1", "198.51.100.9:1234", "203.0.113.77", "198.51.100.9"},
		{"trusted", "192.0.2.1", "192.0.2.1:1234", "203.0.113.77", "203.0.113.77"},
		{"no proxy", "", "192.0.2.1:1234", "203.0.113.77", "192.0.2.1"},
		{"bad proxy", "localhost", "192.0.2.1:1234", "203.0.113.77", "192.0.2.1"},
		{"no header", "192.0.2.1", "192.0.2.1:1234", "", "192.0.2.1"},
		{"bad header", "192.0.2.1", "192.0.2.1:1234", "not an IP", "192.0.2.1"},
		{"IPv6", "2001:db8::1", "[2001:db8::1]:1234", "2001:db8::2", "2001:db8::2"},
		{"mapped IPv4", "192.0.2.1", "[::ffff:192.0.2.1]:1234", "203.0.113.77", "203.0.113.77"},
		{"bare IP", "192.0.2.1", "198.51.100.9", "203.0.113.77", "198.51.100.9"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			trustedProxy, _ = netip.ParseAddr(tt.proxy)
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tt.remote
			r.Header.Set("CF-Connecting-IP", tt.header)
			if got := From(r); got != tt.want {
				t.Errorf("From = %q, want %q", got, tt.want)
			}
		})
	}
}

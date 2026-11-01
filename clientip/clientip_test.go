package clientip

import (
	"net/http/httptest"
	"net/netip"
	"testing"
)

func TestFrom(t *testing.T) {
	defer func(proxy netip.Prefix) { trustedNet = proxy }(trustedNet)
	for _, tt := range []struct {
		name, proxy, remote, header, want string
	}{
		{"spoof", "192.0.2.1", "198.51.100.9:1234", "203.0.113.77", "198.51.100.9"},
		{"trusted", "192.0.2.1", "192.0.2.1:1234", "203.0.113.77", "203.0.113.77"},
		{"no proxy", "", "192.0.2.1:1234", "203.0.113.77", "192.0.2.1"},
		{"no header", "192.0.2.1", "192.0.2.1:1234", "", "192.0.2.1"},
		{"bad header", "192.0.2.1", "192.0.2.1:1234", "not an IP", "192.0.2.1"},
		{"IPv6", "2001:db8::1", "[2001:db8::1]:1234", "2001:db8::2", "2001:db8::2"},
		{"mapped IPv4", "192.0.2.1", "[::ffff:192.0.2.1]:1234", "203.0.113.77", "203.0.113.77"},
		{"bare IP", "192.0.2.1", "198.51.100.9", "203.0.113.77", "198.51.100.9"},
		{"CIDR v4 match", "127.0.0.0/8", "127.1.2.3:1234", "203.0.113.77", "203.0.113.77"},
		{"CIDR v4 miss", "127.0.0.0/8", "128.0.0.1:1234", "203.0.113.77", "128.0.0.1"},
		{"CIDR v4 mapped peer", "127.0.0.0/8", "[::ffff:127.0.0.1]:1234", "203.0.113.77", "203.0.113.77"},
		{"CIDR v6 match", "2001:db8::/32", "[2001:db8:1::5]:1234", "2001:db8::2", "2001:db8::2"},
		{"CIDR v6 miss", "2001:db8::/32", "[2001:db9::5]:1234", "2001:db8::2", "2001:db9::5"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if trustedNet, err = parseProxy(tt.proxy); err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tt.remote
			r.Header.Set("CF-Connecting-IP", tt.header)
			if got := From(r); got != tt.want {
				t.Errorf("From = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseProxyInvalid(t *testing.T) {
	for _, s := range []string{"localhost", "127.0.0.0/33", "1.2.3.4/x", "::ffff:0.0.0.0/64"} {
		if _, err := parseProxy(s); err == nil {
			t.Errorf("parseProxy(%q) succeeded, want error", s)
		}
	}
}

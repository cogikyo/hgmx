package webhost

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

func Registrable(name string, alone bool) (string, error) {
	switch {
	case loopback(name) && alone:
		return name, nil
	case name == "localhost" || strings.HasSuffix(name, ".localhost"):
		return "localhost", nil
	}
	return publicsuffix.EffectiveTLDPlusOne(name)
}

func RelyingParty(raw string) (origin, rp string, err error) {
	origin, name, err := Parse(raw)
	if err != nil {
		return "", "", err
	}
	rp, err = Registrable(name, false)
	if err != nil {
		return "", "", fmt.Errorf("auth: host %s: %w", origin, err)
	}
	return origin, rp, nil
}

func Parse(raw string) (origin, name string, err error) {
	bad := fmt.Errorf("auth: host origin %q must be https://host, or http:// on localhost", raw)
	u, err := url.Parse(raw)
	if err != nil || u.User != nil || strings.TrimSuffix(u.Path, "/") != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", "", bad
	}
	key, name, ok := Authority(u.Scheme, u.Host)
	_, err = netip.ParseAddr(name)
	domain := err != nil
	if !ok || !(u.Scheme == "https" && (domain || loopback(name)) || u.Scheme == "http" && loopback(name)) {
		return "", "", bad
	}
	return u.Scheme + "://" + key, name, nil
}

func Authority(scheme, hostport string) (key, name string, ok bool) {
	u, err := url.Parse(scheme + "://" + hostport)
	if err != nil || hostport == "" || u.Host != hostport || strings.HasSuffix(hostport, ":") {
		return "", "", false
	}
	name, port := u.Hostname(), u.Port()
	if ip, err := netip.ParseAddr(name); err == nil {
		if ip.Zone() != "" {
			return "", "", false
		}
		name = ip.String()
	} else if name, err = idna.Lookup.ToASCII(name); err != nil || name == "" || strings.HasSuffix(name, ".") {
		return "", "", false
	}
	key = name
	if strings.Contains(name, ":") {
		key = "[" + name + "]"
	}
	if port != "" {
		n, err := strconv.ParseUint(port, 10, 16)
		if err != nil || n == 0 {
			return "", "", false
		}
		if p := strconv.FormatUint(n, 10); p != defaultPort[scheme] {
			key += ":" + p
		}
	}
	return key, name, true
}

var defaultPort = map[string]string{"https": "443", "http": "80"}

func loopback(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

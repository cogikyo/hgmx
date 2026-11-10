package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"time"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/internal/webhost"
)

// Host shares accounts and passkeys, but binds sessions to its origin.
type Host struct {
	// Origin is the public HTTPS origin; HTTP is allowed only on loopback hosts.
	Origin string
	Mode   Mode
	// Allow gates sign-in, CLI links, account changes, User and Require.
	// A nil Allow admits every user.
	Allow func(User) bool
	// Layout wraps auth's page root in a full document; required by New.
	Layout func(r *http.Request, title string, body templ.Component) templ.Component
	Views  Views
	Apps   []App // Native clients bound to this host's API origin.
}

type host struct {
	origin   string
	mode     Mode
	allow    func(User) bool
	layout   func(*http.Request, string, templ.Component) templ.Component
	views    Views
	sends    *limiter
	verifies *limiter
	begins   *limiter
}

func (h *host) admits(u User) bool { return h.allow == nil || h.allow(u) }

func parseHosts(hs []Host) (map[string]*host, string, error) {
	if len(hs) == 0 {
		return nil, "", errors.New("auth: Config.Hosts needs at least one host")
	}
	hosts := make(map[string]*host, len(hs))
	var rp string
	for _, h := range hs {
		origin, name, err := webhost.Parse(h.Origin)
		if err != nil {
			return nil, "", err
		}
		if h.Mode < Operator || h.Mode > Open {
			return nil, "", fmt.Errorf("auth: host %s: Mode must be Operator, Invite or Open", origin)
		}
		if hosts[origin] != nil {
			return nil, "", fmt.Errorf("auth: host %s is listed twice", origin)
		}
		if _, err := netip.ParseAddr(name); err == nil && len(hs) > 1 {
			return nil, "", fmt.Errorf("auth: IP host %s must be the only host", origin)
		}
		site, err := webhost.Registrable(name, len(hs) == 1)
		if err != nil {
			return nil, "", fmt.Errorf("auth: host %s: %w", origin, err)
		}
		if rp != "" && site != rp {
			return nil, "", fmt.Errorf("auth: host %s is not on %s", origin, rp)
		}
		rp = site
		hosts[origin] = &host{
			origin:   origin,
			mode:     h.Mode,
			allow:    h.Allow,
			layout:   h.Layout,
			views:    h.Views.withDefaults(),
			sends:    newLimiter(20, time.Hour),
			verifies: newLimiter(20, time.Hour),
			begins:   newLimiter(30, 10*time.Minute),
		}
	}
	return hosts, rp, nil
}

func (a *Auth) host(r *http.Request) (*host, bool) {
	for _, scheme := range [...]string{"https", "http"} {
		if key, _, ok := webhost.Authority(scheme, r.Host); ok {
			if h, ok := a.hosts[scheme+"://"+key]; ok {
				return h, true
			}
		}
	}
	return nil, false
}

func misdirected(w http.ResponseWriter) {
	http.Error(w, "unknown host", http.StatusMisdirectedRequest)
}

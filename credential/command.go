package credential

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
)

type Checks map[string]func(ctx context.Context, value string) error

// Failure is printed verbatim; Reason must not contain secrets or raw errors.
type Failure struct {
	Reason string
	Code   int
}

func (f *Failure) Error() string {
	if f.Code != 0 {
		return fmt.Sprintf("%s (%d)", f.Reason, f.Code)
	}
	return f.Reason
}

// Command handles the site's early credentials dispatch and returns an exit code.
// Pass the server's FlagSet before starting databases, queues or telemetry.
func Command(ctx context.Context, r *Registry, fs *flag.FlagSet, args []string, checks Checks, out io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(out, "usage: credentials manifest | credentials check [name...] [-- site flags...]")
		return 2
	}
	if err := r.wired(checks); err != nil {
		fmt.Fprintf(out, "credentials: %v\n", err)
		return 2
	}
	switch args[0] {
	case "manifest":
		if len(args) > 1 {
			fmt.Fprintln(out, "usage: credentials manifest")
			return 2
		}
		data, err := r.manifest(fs)
		if err != nil {
			fmt.Fprintf(out, "credentials: %v\n", err)
			return 1
		}
		out.Write(data)
		return 0
	case "check":
		return r.check(ctx, fs, args[1:], checks, out)
	}
	fmt.Fprintf(out, "credentials: unknown command %q; want manifest or check\n", args[0])
	return 2
}

func (r *Registry) wired(checks Checks) error {
	var errs []error
	for _, h := range r.handles {
		_, ok := checks[h.def.Name]
		if owned := checkIDs[h.def.Check] != ""; owned && !ok {
			errs = append(errs, fmt.Errorf("%s declares check %s, but Checks has no %s entry", h.def.Name, h.def.Check, h.def.Name))
		} else if !owned && ok {
			errs = append(errs, fmt.Errorf("Checks has %s, which declares no hgmx check", h.def.Name))
		}
	}
	for name := range checks {
		if !slices.ContainsFunc(r.handles, func(h Handle) bool { return h.def.Name == name }) {
			errs = append(errs, fmt.Errorf("Checks has %s, which is not declared", name))
		}
	}
	return errors.Join(errs...)
}

func (r *Registry) kinds(fs *flag.FlagSet) (map[string]string, error) {
	kinds := map[string]string{}
	var errs []error
	fs.VisitAll(func(f *flag.Flag) {
		kind := flagKind(f.Value)
		if kind == "" {
			errs = append(errs, fmt.Errorf("flag -%s has kind %T; credentials.json takes only string and bool flags", f.Name, f.Value))
		}
		kinds["-"+f.Name] = kind
	})
	for _, name := range slices.Sorted(maps.Keys(r.features)) {
		err := r.features[name].walk(func(c Condition) error { return leafKind(c, kinds) })
		if err != nil {
			errs = append(errs, fmt.Errorf("feature %s: %w", name, err))
		}
	}
	return kinds, errors.Join(errs...)
}

func (r *Registry) on(fs *flag.FlagSet) (map[string]bool, error) {
	if _, err := r.kinds(fs); err != nil {
		return nil, err
	}
	flags := map[string]string{}
	fs.Visit(func(f *flag.Flag) { flags["-"+f.Name] = f.Value.String() })
	on := map[string]bool{Always.name: true}
	for _, name := range slices.Sorted(maps.Keys(r.features)) {
		ok, err := r.features[name].holds(flags)
		if err != nil {
			return nil, fmt.Errorf("feature %s: %w", name, err)
		}
		on[name] = ok
	}
	return on, nil
}

func (r *Registry) manifest(fs *flag.FlagSet) ([]byte, error) {
	kinds, err := r.kinds(fs)
	if err != nil {
		return nil, err
	}

	type entry struct {
		Name     string `json:"name"`
		Env      string `json:"env"`
		Purpose  string `json:"purpose"`
		Required string `json:"required"`
		Source   Source `json:"source"`
		Shared   string `json:"shared"`
		Check    string `json:"check,omitempty"`
	}
	m := struct {
		Version     int                  `json:"version"`
		Flags       map[string]string    `json:"flags"`
		Features    map[string]Condition `json:"features"`
		Credentials []entry              `json:"credentials"`
	}{1, kinds, map[string]Condition{}, []entry{}}
	maps.Copy(m.Features, r.features)
	for _, h := range r.handles {
		d := h.def
		m.Credentials = append(m.Credentials, entry{d.Name, d.Env, d.Purpose, h.required, d.Source, h.shared, d.Check})
	}
	slices.SortFunc(m.Credentials, func(a, b entry) int { return strings.Compare(a.Name, b.Name) })

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func flagKind(v flag.Value) string {
	if b, ok := v.(interface{ IsBoolFlag() bool }); ok && b.IsBoolFlag() {
		return "bool"
	}
	if g, ok := v.(flag.Getter); ok {
		if _, ok := g.Get().(string); ok {
			return "string"
		}
	}
	return ""
}

func leafKind(c Condition, kinds map[string]string) error {
	if c.sel.kind != "flag" {
		return nil
	}
	kind, ok := kinds[c.sel.name]
	switch {
	case !ok:
		return fmt.Errorf("flag %s is not registered", c.sel.name)
	case kind != "bool":
		return nil
	case c.op == "set" || c.op == "gte":
		return fmt.Errorf("bool flag %s takes eq or present, not %s", c.sel.name, c.op)
	case c.op == "eq" && c.arg != "true" && c.arg != "false":
		return fmt.Errorf("bool flag %s eq %q: want true or false", c.sel.name, c.arg)
	}
	return nil
}

func (r *Registry) check(ctx context.Context, fs *flag.FlagSet, args []string, checks Checks, out io.Writer) int {
	names, flags := args, []string(nil)
	if i := slices.Index(args, "--"); i >= 0 {
		names, flags = args[:i], args[i+1:]
	}
	fs.SetOutput(io.Discard)
	if err := fs.Parse(flags); err != nil {
		fmt.Fprintf(out, "credentials: site flags: %v\n", err)
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(out, "credentials: site flags: unexpected argument %q\n", fs.Arg(0))
		return 2
	}

	var on map[string]bool
	if len(names) == 0 {
		var err error
		if on, err = r.on(fs); err != nil {
			fmt.Fprintf(out, "credentials: %v\n", err)
			return 2
		}
	}
	var selected, skipped []Handle
	for _, h := range r.handles {
		switch {
		case len(names) > 0:
			if slices.Contains(names, h.def.Name) {
				selected = append(selected, h)
			}
		case checkIDs[h.def.Check] == "":
		case on[h.required]:
			selected = append(selected, h)
		default:
			skipped = append(skipped, h)
		}
	}
	code := 0
	for _, name := range names {
		i := slices.IndexFunc(selected, func(h Handle) bool { return h.def.Name == name })
		switch {
		case i < 0:
			fmt.Fprintf(out, "%s: not declared\n", name)
			code = 2
		case selected[i].def.Check == "":
			fmt.Fprintf(out, "%s: declares no check\n", name)
			code = 2
		case checkIDs[selected[i].def.Check] == "":
			fmt.Fprintf(out, "%s: %s check owned by cogictl\n", name, selected[i].def.Check)
			code = 2
		}
	}
	if len(selected)+len(skipped) == 0 {
		fmt.Fprintln(out, "credentials: no credential declares an hgmx check")
		code = 2
	}
	if code != 0 {
		return code
	}

	for _, h := range skipped {
		fmt.Fprintf(out, "%s: skipped; feature %s is off\n", h.def.Name, h.required)
	}
	for _, h := range selected {
		if !h.probe(ctx, checks[h.def.Name], out) {
			code = 1
		}
	}
	return code
}

func (h Handle) probe(ctx context.Context, fn func(context.Context, string) error, out io.Writer) bool {
	name, check := h.def.Name, h.def.Check
	value, err := h.Read()
	switch {
	case errors.Is(err, ErrMissing):
		fmt.Fprintf(out, "%s: missing from $CREDENTIALS_DIRECTORY and env %s\n", name, h.def.Env)
		return false
	case err != nil:
		fmt.Fprintf(out, "%s: credential file unreadable\n", name)
		return false
	case value == "":
		fmt.Fprintf(out, "%s: empty\n", name)
		return false
	}
	err = fn(ctx, value)
	var f *Failure
	switch {
	case err == nil:
		fmt.Fprintf(out, "%s: %s\n", name, checkIDs[check])
		return true
	case errors.As(err, &f):
		fmt.Fprintf(out, "%s: %s\n", name, f.Error())
	default:
		fmt.Fprintf(out, "%s: %s check failed\n", name, check)
	}
	return false
}

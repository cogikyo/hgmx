package credential

import (
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
)

var label = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

var checkIDs = map[string]string{
	"smtp":          "SMTP authentication passed; no mail sent",
	"google":        "token endpoint accepted client authentication; browser sign-in not tested",
	"stripe":        "",
	"porter-key-id": "",
	"porter-secret": "",
}

type Definition struct {
	Name    string
	Env     string
	Purpose string
	Source  Source
	Check   string
}

type Source interface {
	source() error
}

// Generate requests entropy bytes encoded as unpadded base64url text.
type Generate struct {
	Bytes int
}

type Issued struct {
	Prefix []string
	Format string
}

func (g Generate) source() error {
	if g.Bytes < 1 || g.Bytes > 4096 {
		return fmt.Errorf("generate %d bytes: want 1 to 4096", g.Bytes)
	}
	return nil
}

func (i Issued) source() error {
	for _, p := range i.Prefix {
		if p == "" || strings.ContainsAny(p, "\r\n") {
			return fmt.Errorf("issued prefix %q: want a nonempty single line", p)
		}
	}
	if strings.ContainsAny(i.Format, "\r\n") {
		return fmt.Errorf("issued format: want a single line")
	}
	return nil
}

func (g Generate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind     string `json:"kind"`
		Bytes    int    `json:"bytes"`
		Encoding string `json:"encoding"`
	}{"generate", g.Bytes, "base64url"})
}

func (i Issued) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind   string   `json:"kind"`
		Prefix []string `json:"prefix,omitempty"`
		Format string   `json:"format,omitempty"`
	}{"issued", i.Prefix, i.Format})
}

type Requirement struct {
	name string
}

var Always = Requirement{"always"}

// Registry holds one site's declarations; its zero value is ready to use.
type Registry struct {
	features map[string]Condition
	handles  []Handle
}

type Handle struct {
	def      Definition
	required string
	shared   string
}

// Feature names a deployment condition; it does not switch runtime behavior.
func (r *Registry) Feature(name string, c Condition) Requirement {
	switch {
	case !label.MatchString(name) || name == Always.name:
		panic(fmt.Sprintf("credential: feature %q: want %s, other than always", name, label))
	case c.op == "":
		panic(fmt.Sprintf("credential: feature %s: empty condition", name))
	}
	if _, ok := r.features[name]; ok {
		panic(fmt.Sprintf("credential: feature %s declared twice", name))
	}
	if r.features == nil {
		r.features = map[string]Condition{}
	}
	r.features[name] = c
	return Requirement{name}
}

// Declare copies d and panics on invalid or duplicate metadata.
// required is [Always] or a feature of r; shared is "site" or a shared label.
func (r *Registry) Declare(d Definition, required Requirement, shared string) Handle {
	if err := r.valid(d, required, shared); err != nil {
		panic(fmt.Sprintf("credential: %s: %v", d.Name, err))
	}
	if i, ok := d.Source.(Issued); ok {
		i.Prefix = slices.Clone(i.Prefix)
		d.Source = i
	}
	h := Handle{def: d, required: required.name, shared: shared}
	r.handles = append(r.handles, h)
	return h
}

func (r *Registry) valid(d Definition, required Requirement, shared string) error {
	switch {
	case !label.MatchString(d.Name):
		return fmt.Errorf("name: want %s", label)
	case !envName.MatchString(d.Env):
		return fmt.Errorf("env %q: want %s", d.Env, envName)
	case d.Purpose == "" || strings.ContainsAny(d.Purpose, "\r\n"):
		return fmt.Errorf("purpose: want a nonempty single line")
	case d.Source == nil:
		return fmt.Errorf("needs a Generate or Issued source")
	case shared != "site" && !label.MatchString(shared):
		return fmt.Errorf("shared %q: want site or a label matching %s", shared, label)
	}
	if err := d.Source.source(); err != nil {
		return err
	}
	if _, ok := checkIDs[d.Check]; d.Check != "" && !ok {
		return fmt.Errorf("check %q is not one of %s", d.Check, strings.Join(slices.Sorted(maps.Keys(checkIDs)), ", "))
	}
	if _, ok := r.features[required.name]; !ok && required != Always {
		return fmt.Errorf("requirement %q is not Always or a feature of this registry", required.name)
	}
	for _, h := range r.handles {
		if h.def.Name == d.Name {
			return fmt.Errorf("declared twice")
		}
		if h.def.Env == d.Env {
			return fmt.Errorf("env %s is also %s's", d.Env, h.def.Name)
		}
	}
	return nil
}

func (h Handle) Name() string { return h.def.Name }

// Read tries $CREDENTIALS_DIRECTORY/<name>, then the declared env fallback.
// Only missing files permit fallback; even an empty file takes priority.
// Values lose trailing CR and LF; missing files and empty env return [ErrMissing].
func (h Handle) Read() (string, error) { return read(h.def.Name, h.def.Env) }

package credential

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var envName = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

type Selector struct {
	kind     string
	name     string
	fallback string
}

// Flag requires one leading dash, as in "-accounts".
func Flag(name string) Selector {
	if len(name) < 2 || name[0] != '-' || name[1] == '-' || strings.Contains(name, "=") {
		panic(fmt.Sprintf("credential: flag %q: want one leading dash, such as -gate", name))
	}
	return Selector{kind: "flag", name: name}
}

func Env(name string) Selector {
	if !envName.MatchString(name) {
		panic(fmt.Sprintf("credential: env %q: want %s", name, envName))
	}
	return Selector{kind: "env", name: name}
}

// Fallback uses env only if the flag is absent; empty flags override env.
func (s Selector) Fallback(env Selector) Selector {
	if s.kind != "flag" || s.fallback != "" || env.kind != "env" {
		panic("credential: Fallback takes a Flag and then an Env")
	}
	s.fallback = env.name
	return s
}

// Set requires a supplied, nonempty value.
func (s Selector) Set() Condition { return s.leaf("set", nil) }

// Present includes explicitly supplied empty flags; an empty env value is unset.
func (s Selector) Present() Condition { return s.leaf("present", nil) }

func (s Selector) Eq(v string) Condition { return s.leaf("eq", v) }

// GTE compares nonnegative base-10 int64 values; invalid values are errors.
func (s Selector) GTE(n int64) Condition {
	if n < 0 {
		panic(fmt.Sprintf("credential: %s gte %d: want a nonnegative integer", s.name, n))
	}
	return s.leaf("gte", n)
}

func (s Selector) leaf(test string, arg any) Condition {
	if s.kind == "" {
		panic("credential: condition on an empty selector; use Flag or Env")
	}
	c := Condition{op: test, sel: Selector{kind: s.kind, name: s.name}, arg: arg}
	if s.fallback == "" {
		return c
	}
	env := Condition{op: test, sel: Selector{kind: "env", name: s.fallback}, arg: arg}
	return Any(c, All(Not(Flag(s.name).Present()), env))
}

// Condition describes deployment policy, not site runtime behavior.
// Absent inputs fail leaf tests; flag defaults are not used.
type Condition struct {
	op   string
	sel  Selector
	arg  any
	subs []Condition
}

func All(c ...Condition) Condition { return group("all", c) }
func Any(c ...Condition) Condition { return group("any", c) }

func Not(c Condition) Condition {
	return group("not", []Condition{c})
}

func group(op string, c []Condition) Condition {
	if len(c) == 0 {
		panic("credential: " + op + " needs at least one condition")
	}
	for _, sub := range c {
		if sub.op == "" {
			panic("credential: " + op + " holds an empty condition")
		}
	}
	return Condition{op: op, subs: append([]Condition(nil), c...)}
}

func (c Condition) walk(fn func(Condition) error) error {
	if len(c.subs) == 0 {
		return fn(c)
	}
	for _, sub := range c.subs {
		if err := sub.walk(fn); err != nil {
			return err
		}
	}
	return nil
}

func (c Condition) holds(flags map[string]string) (bool, error) {
	switch c.op {
	case "all", "any":
		decisive := c.op == "any"
		var first error
		for _, sub := range c.subs {
			ok, err := sub.holds(flags)
			if err == nil && ok == decisive {
				return decisive, nil
			}
			if first == nil {
				first = err
			}
		}
		return !decisive, first
	case "not":
		ok, err := c.subs[0].holds(flags)
		return !ok, err
	}
	value, ok := flags[c.sel.name]
	if c.sel.kind == "env" {
		value = os.Getenv(c.sel.name)
		ok = value != ""
	}
	switch {
	case !ok || c.op == "gte" && value == "":
		return false, nil
	case c.op == "present":
		return true, nil
	case c.op == "set":
		return value != "", nil
	case c.op == "eq":
		return value == c.arg, nil
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil || strings.Trim(value, "0123456789") != "" {
		return false, fmt.Errorf("%s is not a nonnegative int64 for gte", c.sel.name)
	}
	return n >= c.arg.(int64), nil
}

func (c Condition) MarshalJSON() ([]byte, error) {
	switch c.op {
	case "":
		return nil, fmt.Errorf("credential: empty condition")
	case "all", "any":
		return json.Marshal(map[string][]Condition{c.op: c.subs})
	case "not":
		return json.Marshal(map[string]Condition{"not": c.subs[0]})
	}
	name, err := json.Marshal(c.sel.name)
	if err != nil {
		return nil, err
	}
	arg := c.arg
	if arg == nil {
		arg = true
	}
	test, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil, `{%q:%s,%q:%s}`, c.sel.kind, name, c.op, test), nil
}

package telemetry

import (
	"log/slog"
	"strings"
	"unicode/utf8"
)

const maxString = 1 << 10

func clip(s string) string {
	if len(s) > maxString {
		i := maxString
		for i > 0 && !utf8.RuneStart(s[i]) {
			i--
		}
		s = s[:i]
	}
	return strings.Clone(s)
}

func snapshot(dst []slog.Attr, a slog.Attr) []slog.Attr {
	v := a.Value.Resolve()
	if a.Key == "" && v.Kind() != slog.KindGroup {
		return dst
	}
	switch v.Kind() {
	case slog.KindGroup:
		var members []slog.Attr
		for _, m := range v.Group() {
			members = snapshot(members, m)
		}
		if len(members) == 0 {
			return dst
		}
		if a.Key == "" {
			return append(dst, members...)
		}
		return append(dst, slog.Attr{Key: clip(a.Key), Value: slog.GroupValue(members...)})
	case slog.KindString:
		v = slog.StringValue(clip(v.String()))
	case slog.KindAny:
		v = slog.StringValue(clip(v.String()))
	}
	return append(dst, slog.Attr{Key: clip(a.Key), Value: v})
}

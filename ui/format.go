package ui

import (
	"strconv"
	"strings"
)

// Expand expands a format string using `git log` message syntax.
func Expand(format string, values map[string]string) string {
	f := &expander{values: values}
	return f.Expand(format)
}

// An expander is a stateful helper to expand a format string.
type expander struct {
	// formatted holds the
	formatted []string

	// values is the map of values that should be expanded.
	values map[string]string

	// skipNext is true if the next placeholder is not a place holder and can be
	// output directly as such.
	skipNext bool
}

func (f *expander) Expand(format string) string {
	parts := strings.Split(format, "%")
	f.formatted = make([]string, 0, len(parts))
	f.append(parts[0])
	for _, p := range parts[1:] {
		v, t := f.expandOneVar(p)
		f.append(v, t)
	}
	return f.crush()
}

func (f *expander) append(formattedText ...string) {
	f.formatted = append(f.formatted, formattedText...)
}

func (f *expander) crush() string {
	s := strings.Join(f.formatted, "")
	f.formatted = nil
	return s
}

var colorMap = map[string]string{
	"red":   "31",
	"green": "32",
	"blue":  "34",
	"reset": "",
}

func (f *expander) expandOneVar(format string) (expand string, untouched string) {
	if f.skipNext {
		f.skipNext = false
		return "", format
	}
	if format == "" {
		f.skipNext = true
		return "", "%"
	}

	if e, u, ok := f.expandSpecialChar(format[0], format[1:]); ok {
		return e, u
	}

	if f.values != nil {
		for i := 1; i <= len(format); i++ {
			if v, exists := f.values[format[0:i]]; exists {
				return v, format[i:]
			}
		}
	}

	return "", "%" + format
}

func (f *expander) expandSpecialChar(firstChar byte, format string) (expand string, untouched string, wasExpanded bool) {
	switch firstChar {
	case 'n':
		return "\n", format, true
	case 'C':
		for k, v := range colorMap {
			if strings.HasPrefix(format, k) {
				return "\033[" + v + "m", format[len(k):], true
			}
		}
		// TODO: Add custom color as specified in color.branch.* options.
    // TODO: Handle auto-coloring.
	case 'x':
		if len(format) >= 2 {
			if v, err := strconv.ParseInt(format[:2], 16, 32); err == nil {
				return string(v), format[2:], true
			}
		}
	case '+':
		if e, u := f.expandOneVar(format); e != "" {
			return "\n" + e, u, true
		} else {
			return "", u, true
		}
	case ' ':
		if e, u := f.expandOneVar(format); e != "" {
			return " " + e, u, true
		} else {
			return "", u, true
		}
	case '-':
		if e, u := f.expandOneVar(format); e != "" {
			return e, u, true
		} else {
			f.append(strings.TrimRight(f.crush(), "\n"))
			return "", u, true
		}
	}
	return "", "", false
}

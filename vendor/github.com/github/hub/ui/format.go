package ui

import (
	"regexp"
	"strconv"
	"strings"
)

// Expand expands a format string using `git log` message syntax.
func Expand(format string, values map[string]string, colorize bool) string {
	f := &expander{values: values, colorize: colorize}
	return f.Expand(format)
}

// An expander is a stateful helper to expand a format string.
type expander struct {
	// formatted holds the parts of the string that have already been formatted.
	formatted []string

	// values is the map of values that should be expanded.
	values map[string]string

	// colorize is a flag to indiciate whether to use colors.
	colorize bool

	// skipNext is true if the next placeholder is not a placeholder and can be
	// output directly as such.
	skipNext bool

	// padNext is an object that should be used to pad the next placeholder.
	padNext *padder
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
	"black":   "30",
	"red":     "31",
	"green":   "32",
	"yellow":  "33",
	"blue":    "34",
	"magenta": "35",
	"cyan":    "36",
	"white":   "37",
	"reset":   "",
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

	if f.padNext != nil {
		p := f.padNext
		f.padNext = nil
		e, u := f.expandOneVar(format)
		return f.pad(e, p), u
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
				if f.colorize {
					return "\033[" + v + "m", format[len(k):], true
				}
				return "", format[len(k):], true
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
	case '<', '>':
		if m := paddingPattern.FindStringSubmatch(string(firstChar) + format); len(m) == 7 {
			if p := padderFromConfig(m[1], m[2], m[3], m[4], m[5]); p != nil {
				f.padNext = p
				return "", m[6], true
			}
		}
	}
	return "", "", false
}

func (f *expander) pad(s string, p *padder) string {
	size := int(p.size)
	if p.sizeAsColumn {
		previous := f.crush()
		f.append(previous)
		size -= len(previous) - strings.LastIndex(previous, "\n") - 1
	}

	numPadding := size - len(s)
	if numPadding == 0 {
		return s
	}

	if numPadding < 0 {
		if p.usePreviousSpace {
			previous := f.crush()
			noBlanks := strings.TrimRight(previous, " ")
			f.append(noBlanks)
			numPadding += len(previous) - len(noBlanks)
		}

		if numPadding <= 0 {
			return p.truncate(s, -numPadding)
		}
	}

	switch p.orientation {
	case padLeft:
		return strings.Repeat(" ", numPadding) + s
	case padMiddle:
		return strings.Repeat(" ", numPadding/2) + s + strings.Repeat(" ", (numPadding+1)/2)
	}

	// Pad right by default.
	return s + strings.Repeat(" ", numPadding)
}

type paddingOrientation int

const (
	padRight paddingOrientation = iota
	padLeft
	padMiddle
)

type truncingMethod int

const (
	noTrunc truncingMethod = iota
	truncLeft
	truncRight
	truncMiddle
)

type padder struct {
	orientation      paddingOrientation
	size             int64
	sizeAsColumn     bool
	usePreviousSpace bool
	truncing         truncingMethod
}

var paddingPattern = regexp.MustCompile(`^(>)?([><])(\|)?\((\d+)(,[rm]?trunc)?\)(.*)$`)

func padderFromConfig(alsoLeft, orientation, asColumn, size, trunc string) *padder {
	p := &padder{}

	if orientation == ">" {
		p.orientation = padLeft
	} else if alsoLeft == "" {
		p.orientation = padRight
	} else {
		p.orientation = padMiddle
	}

	p.sizeAsColumn = asColumn != ""

	var err error
	if p.size, err = strconv.ParseInt(size, 10, 64); err != nil {
		return nil
	}

	p.usePreviousSpace = alsoLeft != "" && p.orientation == padLeft

	switch trunc {
	case ",trunc":
		p.truncing = truncLeft
	case ",rtrunc":
		p.truncing = truncRight
	case ",mtrunc":
		p.truncing = truncMiddle
	}

	return p
}

func (p *padder) truncate(s string, numReduce int) string {
	if numReduce == 0 {
		return s
	}
	numLeft := len(s) - numReduce - 2
	if numLeft < 0 {
		numLeft = 0
	}

	switch p.truncing {
	case truncRight:
		return ".." + s[len(s)-numLeft:len(s)]
	case truncMiddle:
		return s[:numLeft/2] + ".." + s[len(s)-(numLeft+1)/2:len(s)]
	}

	// Trunc left by default.
	return s[:numLeft] + ".."
}

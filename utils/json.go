package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type state struct {
	isObject    bool
	isArray     bool
	arrayIndex  int
	objectKey   string
	parentState *state
}

func stateKey(s *state) string {
	k := ""
	if s.parentState != nil {
		k = stateKey(s.parentState)
	}
	if s.isObject {
		return fmt.Sprintf("%s.%s", k, s.objectKey)
	} else if s.isArray {
		return fmt.Sprintf("%s.[%d]", k, s.arrayIndex)
	} else {
		return k
	}
}

func JSONPath(out io.Writer, src io.Reader, colorize bool) (hasNextPage bool, endCursor string) {
	dec := json.NewDecoder(src)
	dec.UseNumber()

	s := &state{}
	postEmit := func() {
		if s.isObject {
			s.objectKey = ""
		} else if s.isArray {
			s.arrayIndex++
		}
	}

	color := func(c string, t interface{}) string {
		if colorize {
			return fmt.Sprintf("\033[%sm%s\033[m", c, t)
		} else if tt, ok := t.(string); ok {
			return tt
		} else {
			return fmt.Sprintf("%s", t)
		}
	}

	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				s = &state{isObject: true, parentState: s}
			case '[':
				s = &state{isArray: true, parentState: s}
			case '}', ']':
				s = s.parentState
				postEmit()
			default:
				panic("unknown delim")
			}
		} else {
			if s.isObject && s.objectKey == "" {
				s.objectKey = token.(string)
			} else {
				k := stateKey(s)
				fmt.Fprintf(out, "%s\t", color("0;36", k))

				switch tt := token.(type) {
				case string:
					fmt.Fprintf(out, "%s\n", strings.Replace(tt, "\n", "\\n", -1))
					if strings.HasSuffix(k, ".pageInfo.endCursor") {
						endCursor = tt
					}
				case json.Number:
					fmt.Fprintf(out, "%s\n", color("0;35", tt))
				case nil:
					fmt.Fprintf(out, "\n")
				case bool:
					fmt.Fprintf(out, "%s\n", color("1;33", fmt.Sprintf("%v", tt)))
					if strings.HasSuffix(k, ".pageInfo.hasNextPage") {
						hasNextPage = tt
					}
				default:
					panic("unknown type")
				}
				postEmit()
			}
		}
	}
	return
}

package commands

import (
	"fmt"
	"strings"
)

type stringSliceValue []string

func (s *stringSliceValue) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func (s *stringSliceValue) String() string {
	return fmt.Sprintf("%s", *s)
}

type mapValue map[string]string

func (m mapValue) Set(val string) error {
	v := strings.SplitN(val, "=", 2)
	if len(v) != 2 {
		return fmt.Errorf("Flag should be in the format of <name>=<value>")
	}

	m[v[0]] = v[1]

	return nil
}

func (m mapValue) String() string {
	s := make([]string, 0)
	for k, v := range m {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(s, ",")
}

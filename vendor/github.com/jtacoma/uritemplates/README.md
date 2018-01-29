# uritemplates
--
    import "github.com/jtacoma/uritemplates"

[![Build Status](https://travis-ci.org/jtacoma/uritemplates.png)](https://travis-ci.org/jtacoma/uritemplates) [![Coverage Status](https://coveralls.io/repos/jtacoma/uritemplates/badge.png)](https://coveralls.io/r/jtacoma/uritemplates)

Package uritemplates is a level 4 implementation of RFC 6570 (URI
Template, http://tools.ietf.org/html/rfc6570).

To use uritemplates, parse a template string and expand it with a value
map:

	template, _ := uritemplates.Parse("https://api.github.com/repos{/user,repo}")
	values := make(map[string]interface{})
	values["user"] = "jtacoma"
	values["repo"] = "uritemplates"
	expanded, _ := template.Expand(values)
	fmt.Printf(expanded)

## License

Use of this source code is governed by a BSD-style license that can be found in
the LICENSE file.

// Package hypermedia provides helpers for parsing hypermedia links in resources
// and expanding the links to make further requests.
package hypermedia

import (
	"fmt"
	"github.com/jtacoma/uritemplates"
	"net/url"
	"reflect"
)

// Hyperlink is a string url.  If it is a uri template, it can be converted to
// a full URL with Expand().
type Hyperlink string

// Expand converts a uri template into a url.URL using the given M map.
func (l Hyperlink) Expand(m M) (*url.URL, error) {
	template, err := uritemplates.Parse(string(l))
	if err != nil {
		return nil, err
	}

	// clone M to map[string]interface{}
	// if we don't do this type assertion will
	// fail on jtacoma/uritemplates
	// see https://github.com/jtacoma/uritemplates/blob/master/uritemplates.go#L189
	mm := make(map[string]interface{}, len(m))
	for k, v := range m {
		mm[k] = v
	}

	expanded, err := template.Expand(mm)
	if err != nil {
		return nil, err
	}

	return url.Parse(expanded)
}

// M represents a map of values to expand a Hyperlink.
type M map[string]interface{}

// Relations is a map of keys that point to Hyperlink objects.
type Relations map[string]Hyperlink

// Rel fetches and expands the Hyperlink by its given key in the Relations map.
func (h Relations) Rel(name string, m M) (*url.URL, error) {
	if rel, ok := h[name]; ok {
		return rel.Expand(m)
	}
	return nil, fmt.Errorf("No %s relation found", name)
}

// A HypermediaResource has link relations for next actions of a resource.
type HypermediaResource interface {
	Rels() Relations
}

// The HypermediaDecoder gets the link relations from any HypermediaResource.
func HypermediaDecoder(res HypermediaResource) Relations {
	return res.Rels()
}

// HALResource is a resource with hypermedia specified as JSON HAL.
//
// http://stateless.co/hal_specification.html
type HALResource struct {
	Links Links `json:"_links"`
	rels  Relations
}

// Rels gets the link relations from the HALResource's Links field.
func (r *HALResource) Rels() Relations {
	if r.rels == nil {
		r.rels = make(map[string]Hyperlink)
		for name, link := range r.Links {
			r.rels[name] = link.Href
		}
	}
	return r.rels
}

// Links is a collection of Link objects in a HALResource.  Note that the HAL
// spec allows single link objects or an array of link objects.  Sawyer
// currently only supports single link objects.
type Links map[string]Link

// Link represents a single link in a HALResource.
type Link struct {
	Href Hyperlink `json:"href"`
}

// Expand converts a uri template into a url.URL using the given M map.
func (l *Link) Expand(m M) (*url.URL, error) {
	return l.Href.Expand(m)
}

// The HyperFieldDecoder gets link relations from a resource by reflecting on
// its Hyperlink properties.  The relation name is taken either from the name
// of the field, or a "rel" struct tag.
//
//   type Foo struct {
//     Url         Hyperlink `rel:"self" json:"url"`
//     CommentsUrl Hyperlink `rel:"comments" json:"comments_url"`
//   }
//
func HyperFieldDecoder(res interface{}) Relations {
	rels := make(Relations)
	t := reflect.TypeOf(res).Elem()
	v := reflect.ValueOf(res).Elem()
	fieldlen := t.NumField()
	for i := 0; i < fieldlen; i++ {
		fillRelation(rels, t, v, i)
	}
	return rels
}

func fillRelation(rels map[string]Hyperlink, t reflect.Type, v reflect.Value, index int) {
	f := t.Field(index)

	if hyperlinkType != f.Type {
		return
	}

	hl := v.Field(index).Interface().(Hyperlink)
	name := f.Name
	if rel := f.Tag.Get("rel"); len(rel) > 0 {
		name = rel
	}
	rels[name] = hl
}

var hyperlinkType = reflect.TypeOf(Hyperlink("foo"))

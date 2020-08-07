package ssssg

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
)

type Component struct {
	Path    string
	Content string
}

// readComponents recursively from dir, and return the path and content of each.
func readComponents(dir string) ([]Component, error) {
	var components []Component
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subComponents, err := readComponents(path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			components = append(components, subComponents...)
			continue
		}

		p := path.Join(dir, entry.Name())
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		components = append(components, Component{
			Path:    p,
			Content: string(content),
		})
	}

	return components, nil
}

// parseComponents into the given template, returning the new template.
func parseComponents(t *template.Template, dir string) (*template.Template, error) {
	components, err := readComponents(dir)
	if err != nil {
		return nil, err
	}

	for _, c := range components {
		p := strings.TrimPrefix(c.Path, dir)
		p = strings.TrimPrefix(p, "/")
		fmt.Println("Parsing component", p)
		t, err = t.New(p).Parse(c.Content)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

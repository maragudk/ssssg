package ssssg

import (
	"io/ioutil"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Page struct {
	Path string
	Meta struct {
		Title       string
		Description string
	}
	Body   string
	Layout *string
}

// readPages recursively from dir, saving the config content, the content, and the path.
func readPages(dir string) ([]Page, error) {
	var pages []Page
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subPages, err := readPages(path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			pages = append(pages, subPages...)
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		configContent, err := ioutil.ReadFile(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var page Page
		if err := yaml.Unmarshal(configContent, &page); err != nil {
			return nil, err
		}

		body, err := ioutil.ReadFile(path.Join(dir, strings.TrimSuffix(entry.Name(), ".yaml")+".html"))
		if err != nil {
			return nil, err
		}
		page.Path = path.Join(dir, entry.Name())
		page.Body = string(body)
		pages = append(pages, page)
	}
	return pages, nil
}

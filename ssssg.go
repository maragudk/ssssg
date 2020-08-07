package ssssg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

type Page struct {
	Path   string
	Config Config
	Body   string
	Layout string
}

type Config struct {
	Meta struct {
		Title       string
		Description string
	}
}

type BuildOptions struct {
	LayoutsDir string
	PagesDir   string
	BuildDir   string
}

func Build(options BuildOptions) error {
	pages, err := exploreDir(options.PagesDir)
	if err != nil {
		return err
	}
	t, err := template.New("layouts").ParseFiles(path.Join(options.LayoutsDir, "default.html"))
	if err != nil {
		return err
	}

	for _, page := range pages {
		fmt.Println("Building", page.Path)
		if page.Layout == "" {
			page.Layout = "default.html"
		}

		outputPath := path.Join(options.BuildDir, strings.TrimSuffix(strings.TrimPrefix(page.Path, options.PagesDir+"/"), ".yaml")) + ".html"

		if err := os.MkdirAll(path.Dir(outputPath), 0766); err != nil {
			return err
		}

		output, err := os.Create(outputPath)
		defer func(path string) {
			if err := output.Close(); err != nil {
				log.Println("Could not close", path, ":", err)
			}
		}(outputPath)
		if err != nil {
			return err
		}
		if err := t.ExecuteTemplate(output, page.Layout, page); err != nil {
			return err
		}
	}
	return nil
}

func exploreDir(dir string) ([]Page, error) {
	var pages []Page
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subPages, err := exploreDir(path.Join(dir, entry.Name()))
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
		var config Config
		if err := yaml.Unmarshal(configContent, &config); err != nil {
			return nil, err
		}

		body, err := ioutil.ReadFile(path.Join(dir, strings.TrimSuffix(entry.Name(), ".yaml")+".html"))
		page := Page{
			Path:   path.Join(dir, entry.Name()),
			Config: config,
			Body:   string(body),
		}
		pages = append(pages, page)
	}
	return pages, nil
}

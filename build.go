package ssssg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type BuildOptions struct {
	BuildDir      string
	ComponentsDir string
	LayoutsDir    string
	PagesDir      string
	StaticsDir    string
}

func Build(options BuildOptions) error {
	copyStatics := exec.Command("cp", "-va", options.StaticsDir+"/", options.BuildDir)
	copyOutput, err := copyStatics.CombinedOutput()
	fmt.Println(string(copyOutput))
	if err != nil {
		return fmt.Errorf("could not copy static files")
	}

	pages, err := readPages(options.PagesDir)
	if err != nil {
		return err
	}

	t, err := template.New("layouts").ParseGlob(path.Join(options.LayoutsDir, "*.html"))
	if err != nil {
		return err
	}

	t = template.Must(t.New("nolayout").Parse("{{.Body}}"))

	t, err = parseComponents(t, options.ComponentsDir)
	if err != nil {
		return err
	}

	for _, page := range pages {
		fmt.Println("Building", page.Path)
		var content strings.Builder
		t, err = t.New(page.Path).Parse(page.Body)
		if err != nil {
			return err
		}

		err = t.ExecuteTemplate(&content, page.Path, nil)
		if err != nil {
			return err
		}
		page.Body = content.String()

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

		layoutName := "default.html"
		if page.Layout != nil {
			layoutName = *page.Layout
		}
		if layoutName == "" {
			layoutName = "nolayout"
		}
		if err := t.ExecuteTemplate(output, layoutName, page); err != nil {
			return err
		}
	}

	return nil
}

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

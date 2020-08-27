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

	"github.com/fsnotify/fsnotify"

	"ssssg/errors2"
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

	t, err = parseComponents(t, options.ComponentsDir)
	if err != nil {
		return err
	}

	for _, page := range pages {
		fmt.Println("Building", page.Path)
		t, err := t.New(page.Path).Parse(page.Content)
		if err != nil {
			return err
		}

		outputPath := path.Join(options.BuildDir, strings.TrimPrefix(page.Path, options.PagesDir+"/"))

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

		if err := t.ExecuteTemplate(output, page.Path, nil); err != nil {
			return err
		}
	}

	return nil
}

// Watch for changes and build.
func Watch(options BuildOptions) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors2.Wrap(err, "could not create file watcher")
	}
	defer func() {
		_ = watcher.Close()
	}()

	done := make(chan error)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Println(event.Name, "changed, building", event.String())
				if err := Build(options); err != nil {
					done <- err
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				done <- err
			}
		}
	}()

	for _, dir := range []string{options.ComponentsDir, options.LayoutsDir, options.PagesDir, options.StaticsDir} {
		if err := watcher.Add(dir); err != nil {
			return errors2.Wrap(err, "could not add %v to watcher", dir)
		}
	}
	err = <-done

	return errors2.Wrap(err, "error watching")
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
	Path    string
	Content string
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
		if !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		var page Page
		content, err := ioutil.ReadFile(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		page.Path = path.Join(dir, entry.Name())
		page.Content = string(content)
		pages = append(pages, page)
	}
	return pages, nil
}

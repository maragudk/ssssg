package ssssg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
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

	t, err := template.New("layouts/default.html").ParseFiles(path.Join(options.LayoutsDir, "default.html"))
	if err != nil {
		return err
	}

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

		if err := t.ExecuteTemplate(output, "default.html", page); err != nil {
			return err
		}
	}

	return nil
}

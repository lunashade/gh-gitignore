package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cli/go-gh"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Template struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type Result struct {
	apiPath string
	tmpl    *Template
}

func (r *Result) Print() {
	fmt.Println("# ===", r.tmpl.Name, "===")
	fmt.Println("# source: gh api", r.apiPath)
	fmt.Println(r.tmpl.Source)
}

func main() {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	list := []string{}
	err = client.Get("gitignore/templates", &list)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	indices, err := fuzzyfinder.FindMulti(
		list,
		func(i int) string {
			return list[i]
		},
		fuzzyfinder.WithHeader("choose template name(s). [ESC]:abort/[TAB]:toggle select"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	results := make([]*Result, 0, len(indices))
	for _, idx := range indices {
		item := list[idx]
		r := &Result{
			apiPath: path.Join("gitignore/templates", item),
			tmpl:    new(Template),
		}
		err = client.Get(r.apiPath, r.tmpl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		results = append(results, r)
	}

	for _, r := range results {
		r.Print()
	}
}

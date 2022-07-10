package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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
	// parse flags
	flag.Parse()
	args := flag.Args()
	useFuzzy := len(args) == 0

	// connect
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	// get avaliable kinds
	kinds := []string{}
	err = client.Get("gitignore/templates", &kinds)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// find template kind to print
	var indices []int
	if useFuzzy {
		indices = getIndicesFuzzy(kinds)
	} else {
		indices = getIndicesFromArgs(kinds, args)
	}

	// print results
	results := make([]*Result, 0, len(indices))
	for _, idx := range indices {
		item := kinds[idx]
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

func getIndicesFuzzy(kinds []string) []int {
	indices, err := fuzzyfinder.FindMulti(
		kinds,
		func(i int) string {
			return kinds[i]
		},
		fuzzyfinder.WithHeader("choose template name(s). [ESC]:abort/[TAB]:toggle select"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return indices
}

func getIndicesFromArgs(kinds, args []string) []int {
	// split comma-separated argument
	cands := make([]string, 0, len(args))
	for _, arg := range args {
		cands = append(cands, strings.Split(arg, ",")...)
	}

	// find indices
	kind_to_index := make(map[string]int)
	for i, kind := range kinds {
		kind = strings.ToLower(kind)
		kind_to_index[kind] = i
	}
	indices := make([]int, 0, len(cands))
	for _, cand := range cands {
		cand = strings.ToLower(cand)
		i, ok := kind_to_index[cand]
		if ok {
			indices = append(indices, i)
		} else {
			fmt.Fprintln(os.Stderr, "no available gitignore for:", cand)
		}
	}
	return indices
}

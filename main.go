package main

import (
	"fmt"
	"log"
	"path"

	"github.com/cli/go-gh"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Template struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

func main() {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	response := []string{}
	err = client.Get("gitignore/templates", &response)
	if err != nil {
		log.Fatal(err)
		return
	}

	indices, err := fuzzyfinder.FindMulti(
		response,
		func(i int) string {
			return response[i]
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, idx := range indices {
		item := response[idx]
		pat := path.Join("gitignore/templates", item)
		resp := new(Template)
		err = client.Get(pat, resp)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("# ===", resp.Name, "===")
		fmt.Println("# gh api", pat)
		fmt.Println(resp.Source)
	}
}

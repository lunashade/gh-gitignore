package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cli/go-gh"
	"github.com/ktr0731/go-fuzzyfinder"
)

type RepoContent struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

func main() {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	response := []RepoContent{}
	err = client.Get("repos/github/gitignore/contents", &response)
	if err != nil {
		log.Fatal(err)
		return
	}

	ignores := make([]RepoContent, 0, len(response))
	for _, c := range response {
		if strings.HasSuffix(c.Name, ".gitignore") {
			ignores = append(ignores, c)
		}
	}

	indices, err := fuzzyfinder.FindMulti(
		ignores,
		func(i int) string {
			return ignores[i].Name
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, idx := range indices {
		item := ignores[idx]
		resp, err := http.Get(item.DownloadURL)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatal(resp.Status)
			return
		}
		fmt.Println("#", item.Name)
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

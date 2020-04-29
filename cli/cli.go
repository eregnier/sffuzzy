package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/eregnier/sffuzzy"
)

// Args : Entry point parameters for fuzzy cli
var Args struct {
	Search    string `arg:"positional" arg:"required" help:"Search terms to find in given data"`
	Limit     int    `arg:"-l" help:"Results limit, use -1 for no limit" default:"10"`
	Sort      bool   `arg:"-s" help:"Whether or not results are sorted" default:"true"`
	Normalize bool   `arg:"-n" help:"normalize search string and data string for searching. It fuzzy search with no accents/special characters" default:"true"`
}

func main() {
	arg.MustParse(&Args)
	if Args.Search == "" {
		log.Println("No search term provided, aborting")
		os.Exit(1)
	}
	options := sffuzzy.Options{
		Sort:      Args.Sort,
		Limit:     Args.Limit,
		Normalize: Args.Normalize,
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("Error while reading data content from stdin. Aborting")
		os.Exit(1)
	}

	targets := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	results := sffuzzy.SearchOnce(Args.Search, &targets, options)
	j, err := json.MarshalIndent(results.Results, "", "  ")
	if err != nil {
		log.Println("Error while converting search results to json. Aborting")
		os.Exit(1)
	}
	fmt.Println(string(j))
}

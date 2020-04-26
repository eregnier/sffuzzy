package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	d, _ := ioutil.ReadFile("sample.csv")
	names := strings.Split(string(d), "\n")

	search := "osakajapan"
	searchOnce(&names, search)
	cacheSearch(&names, search)
}

func cacheSearch(names *[]string, search string) {
	fmt.Println("\n + Perform cache search, first search is slower.")

	options := Options{Sort: true, AllowedTypos: 5, Normalize: false}

	//First search with manual caching, this is slower
	s := time.Now().UnixNano()
	cacheTargets := Prepare(names, options)
	Search(search, cacheTargets, options)
	fmt.Println("duration ms>", (time.Now().UnixNano()-s)/int64(time.Millisecond))

	justSearch := func(search string) *SearchResult {
		s = time.Now().UnixNano()
		result := Search(search, cacheTargets, options)
		fmt.Println("duration ms>", (time.Now().UnixNano()-s)/int64(time.Millisecond))
		return result
	}

	//Fast subsequents searches
	fmt.Println(" + Perform cached searches")
	fmt.Println(justSearch("san fransisco").Results[0:5])
	fmt.Println(justSearch("mumbai").Results[0:5])
}

func searchOnce(names *[]string, search string) {
	options := Options{Sort: true, AllowedTypos: 5, Normalize: true}
	s := time.Now().UnixNano()
	results := SearchOnce(search, names, options)
	fmt.Println("duration ms>", (time.Now().UnixNano()-s)/int64(time.Millisecond))
	j, _ := json.MarshalIndent(results.Results[:10], "", "  ")
	fmt.Println(string(j))
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if os.Getenv("DEBUG") != "1" {
		log.SetOutput(ioutil.Discard)
	}
}

func TestMinimalSearch(t *testing.T) {
	names := []string{"super man", "super noel", "super du"}
	results := SearchOnce("perdu", &names, Options{Sort: true, AllowedTypos: 5, Normalize: true})
	log.Println("TestMinimalSearch", results)
}

func TestMinimalSearchCache(t *testing.T) {
	names := []string{"super man", "super noel", "super du"}
	options := Options{Sort: true, AllowedTypos: 5, Normalize: true}
	cacheTargets := Prepare(&names, options)
	results := Search("perdu", cacheTargets, options)
	log.Println("TestMinimalSearchCache", results)
}

func TestCacheSearch(t *testing.T) {
	d, _ := ioutil.ReadFile("sample.csv")
	names := strings.Split(string(d), "\n")
	search := "osakajapan"
	options := Options{Sort: true, AllowedTypos: 5, Normalize: false}

	log.Println("\n + Perform cache search, first search is slower.")

	//First search with manual caching, this is slower
	s := time.Now().UnixNano()
	cacheTargets := Prepare(&names, options)
	Search(search, cacheTargets, options)
	deltaFirstSearch := (time.Now().UnixNano() - s) / int64(time.Millisecond)
	log.Println("duration ms>", deltaFirstSearch)

	justSearch := func(search string, deltaFirstSearch int64, t *testing.T) *SearchResult {
		s = time.Now().UnixNano()
		result := Search(search, cacheTargets, options)
		deltaCacheSearch := (time.Now().UnixNano() - s) / int64(time.Millisecond)
		log.Println("duration ms>", deltaCacheSearch)
		if float64(deltaCacheSearch) > float64(deltaFirstSearch/2) {
			t.Errorf("Expected cache search is at least 2 times faster than first search")
		}
		return result
	}

	//Fast subsequents searches
	log.Println(" + Perform cached searches")
	log.Println(justSearch("san fransisco", deltaFirstSearch, t).Results[0:5])
	log.Println(justSearch("mumbai", deltaFirstSearch, t).Results[0:5])
}

func TestSearchOnce(t *testing.T) {

	d, _ := ioutil.ReadFile("sample.csv")
	names := strings.Split(string(d), "\n")
	search := "osakajapan"
	options := Options{Sort: true, AllowedTypos: 5, Normalize: true}

	s := time.Now().UnixNano()
	results := SearchOnce(search, &names, options)
	log.Println("duration ms>", (time.Now().UnixNano()-s)/int64(time.Millisecond))

	tables := []struct {
		Target     string
		Score      int
		MatchCount int
		Typos      int
		Complete   bool
	}{
		{"Ōsaka;Japan", 10, 10, 1, true},
		{"Yuzhno-Sakhalinsk;Russia", 5, 5, 5, false},
		{"Oshakati;Namibia", 5, 5, 5, false},
		{"Makedonska Kamenica;Macedonia", 5, 5, 5, false},
		{"Zhosaly;Kazakhstan", 5, 5, 5, false},
	}
	for x, result := range tables {
		if result.Target != results.Results[x].Target {
			t.Errorf("Expecting Target to be %s got %s", result.Target, results.Results[x].Target)
		}
		if result.Score != results.Results[x].Score {
			t.Errorf("Expecting Score to be %d got %d", result.Score, results.Results[x].Score)
		}
		if result.MatchCount != results.Results[x].MatchCount {
			t.Errorf("Expecting MatchCount to be %d got %d", result.MatchCount, results.Results[x].MatchCount)
		}
		if result.Typos != results.Results[x].Typos {
			t.Errorf("Expecting Typos to be %d got %d", result.Typos, results.Results[x].Typos)
		}
		if result.Complete != results.Results[x].Complete {
			t.Errorf("Expecting Complete to be %t got %t", result.Complete, results.Results[x].Complete)
		}

	}

	j, _ := json.MarshalIndent(results.Results[:10], "", "  ")
	log.Println(string(j))
}

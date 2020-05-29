package sffuzzy

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

func computeDuration(s int64) float64 {
	duration := float64((time.Now().UnixNano()-s)/int64(time.Nanosecond)) / 1000000.0
	log.Printf(" 🕑 Duration: %fms", duration)
	return duration
}

func getTargets() []Target {
	d, _ := ioutil.ReadFile("sample.csv")
	names := strings.Split(string(d), "\n")
	targets := make([]Target, len(names))
	for i, name := range names {
		targets[i] = Target{token: name}
	}
	return targets
}

func TestMinimalSearch(t *testing.T) {
	names := []Target{
		Target{token: "super man"},
		Target{token: "super noel"},
		Target{token: "super du"},
	}
	results := SearchOnce("perdu", &names, Options{Sort: true, Normalize: true})
	log.Println("TestMinimalSearch", results)
}

func TestRelatedDocuments(t *testing.T) {
	names := []Target{
		Target{token: "super man", relatedDocument: "A"},
		Target{token: "super noel", relatedDocument: "B"},
		Target{token: "super du", relatedDocument: "C"},
	}
	results := SearchOnce("perdu", &names, Options{Sort: true, Normalize: true})
	log.Println("TestRelatedDocuments", results)
	for i, document := range []string{"C", "A", "B"} {
		if results.Results[i].RelatedDocument != document {
			t.Errorf("Unable to find related document in result for search term %s", results.Results[i].Target)
		}
	}
}

func TestMinimalSearchCache(t *testing.T) {
	names := []Target{
		Target{token: "super man"},
		Target{token: "super noel"},
		Target{token: "super du"},
	}
	options := Options{Sort: true, Normalize: true}
	cacheTargets := Prepare(&names, options)
	results := Search("perdu", cacheTargets, options)
	log.Println("TestMinimalSearchCache", results)
}

func TestCacheSearch(t *testing.T) {
	targets := getTargets()
	search := "osakajapan"
	options := Options{Sort: true, Normalize: false, Limit: 3}

	log.Println(" + Cache search, first search is slower.")

	//First search with manual caching, this is slower
	s := time.Now().UnixNano()
	cacheTargets := Prepare(&targets, options)
	Search(search, cacheTargets, options)
	deltaFirstSearch := computeDuration(s)

	justSearch := func(search string, deltaFirstSearch float64, t *testing.T) *SearchResult {
		s = time.Now().UnixNano()
		result := Search(search, cacheTargets, options)
		deltaCacheSearch := computeDuration(s)
		if float64(deltaCacheSearch) > float64(deltaFirstSearch/2) {
			t.Errorf("Expected cache search is at least 2 times faster than first search")
		}
		return result
	}

	//Fast subsequents searches
	log.Println(" + Cached searches")
	log.Println(justSearch("san fransisco", deltaFirstSearch, t).Results)
	log.Println(justSearch("mumbai", deltaFirstSearch, t).Results)
	log.Println(justSearch("hong kong", deltaFirstSearch, t).Results)
	log.Println(justSearch("agadez", deltaFirstSearch, t).Results)
	log.Println(justSearch("Palma", deltaFirstSearch, t).Results)
	log.Println(justSearch("sucre bolivia", deltaFirstSearch, t).Results)
	log.Println(justSearch("ibb yemen", deltaFirstSearch, t).Results)
	log.Println(justSearch("west view", deltaFirstSearch, t).Results)
}

func TestSearchOnce(t *testing.T) {
	log.Println(" + Search all at once")

	targets := getTargets()
	search := "osakajapan"
	options := Options{Sort: true, Normalize: true, Limit: 5}

	s := time.Now().UnixNano()
	results := SearchOnce(search, &targets, options)
	computeDuration(s)

	tables := []struct {
		Target     string
		Score      int
		MatchCount int
		Typos      int
	}{
		{"Ōsaka;Japan", 13, 10, 1},
		{"Sri Jayewardenepura Kotte;Sri Lanka", 8, 5, 9},
		{"South Salt Lake;United States", 7, 5, 23},
		{"Vientiane;Laos", 7, 2, 0},
		{"Chimboy Shahri;Uzbekistan", 7, 5, 15},
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

	}

	j, _ := json.MarshalIndent(results.Results, "", "  ")
	log.Println("Print plain unmarshaled json results")
	log.Println(string(j))
}

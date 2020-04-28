package sffuzzy

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var noResults []string

// Options : Search options
type Options struct {
	Sort         bool
	AllowedTypos int
	Normalize    bool
}

// SearchResult : Current search results wrapper
type SearchResult struct {
	TotalComplete int               `json:"totalComplete"`
	Results       []algorithmResult `json:"results"`
	BestScore     int               `json:"bestScore"`
}

// algorithmResult : Data generated on string evaluation
type algorithmResult struct {
	Target     string `json:"target"`
	Score      int    `json:"score"`
	MatchCount int    `json:"matchCount"`
	Typos      int    `json:"typos"`
	Complete   bool   `json:"complete"`
}

// batchItem : handles item range to process per batch
type batchItem struct {
	start int
	stop  int
}

// CacheTarget : data structure handling search payload cache
type CacheTarget struct {
	target string
	cache  []string
}

type searchTerm struct {
	term []string
	len  int
}

//Unicode clean callback
func cleanUnicode(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func normalize(text string, options Options) string {
	if options.Normalize {
		t := transform.Chain(norm.NFD, transform.RemoveFunc(cleanUnicode), norm.NFC)
		text, _, _ = transform.String(t, text)
	}
	return strings.ToLower(text)
}

//Prepare data set for multi searches
func Prepare(targets *[]string, options Options) *[]CacheTarget {

	cacheTargets := make([]CacheTarget, len(*targets))

	for i, target := range *targets {
		preparedTerm := strings.Split(normalize(target, options), "")
		cacheTargets[i] = CacheTarget{target: target, cache: preparedTerm}
	}
	return &cacheTargets
}

func prepareSearch(search string, options Options) [][]searchTerm {
	searchTerms := strings.Split(normalize(search, options), " ")
	result := make([][]searchTerm, len(searchTerms))
	for x, term := range searchTerms {
		result[x].term = strings.Split(term, "")
	}
	return result
}

// SearchOnce : shorthand function to trigger search and caching at once
func SearchOnce(search string, targets *[]string, options Options) *SearchResult {
	cacheTargets := Prepare(targets, options)
	return Search(search, cacheTargets, options)
}

// Search : function to perform the fuzzy search
func Search(search string, cacheTargets *[]CacheTarget, options Options) *SearchResult {
	if search == "" {
		return nil
	}

	preparedSearch := prepareSearch(search, options)
	searchLen := len(preparedSearch)

	targetLen := len(*cacheTargets)
	results := make([]algorithmResult, targetLen)
	resultWrapper := SearchResult{TotalComplete: 0, Results: results, BestScore: 0}

	for i, cacheTarget := range *cacheTargets {
		result := algorithmResult{Target: cacheTarget.target, Score: 0, Typos: 0, MatchCount: 0, Complete: false}
		algorithm(preparedSearch, &cacheTarget.cache, &result, options)
		if result.Complete {
			resultWrapper.TotalComplete++
		}
		if result.Typos == 0 {
			result.Score++
		}
		if result.Score > resultWrapper.BestScore {
			resultWrapper.BestScore = result.Score
		}
		results[i] = result
	}
	if options.Sort {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Score > results[j].Score
		})
	}
	return &resultWrapper
}

func algorithm(search [][]string, target *[]string, result *algorithmResult, options Options) {

	targetLen := len(*target)
	for _, term := range search {
		searchI := 0
		targetI := 0
		for {
			if term[searchI] == (*target)[targetI] {
				result.MatchCount++
				result.Score++
				searchI++
				if searchI == searchLen || result.MatchCount == searchLen {
					result.Complete = true
					return
				}
			} else {
				if searchI != 0 {
					result.Typos++
				}
				if result.MatchCount > 0 && (options.AllowedTypos != -1 && result.Typos >= options.AllowedTypos) {
					return
				}
			}
			targetI++
			if searchI == searchLen || targetI == targetLen {
				return
			}
		}
	}

}

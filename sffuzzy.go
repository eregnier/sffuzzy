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
	Sort      bool
	Normalize bool
	Limit     int
}

// SearchResult : Current search results wrapper
type SearchResult struct {
	Results   []algorithmResult `json:"results"`
	BestScore int               `json:"bestScore"`
}

// algorithmResult : Data generated on string evaluation
type algorithmResult struct {
	Target     string `json:"target"`
	Score      int    `json:"score"`
	MatchCount int    `json:"matchCount"`
	Typos      int    `json:"typos"`
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
	len    int
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
		cacheTargets[i] = CacheTarget{target: target, cache: preparedTerm, len: len(preparedTerm)}
	}
	return &cacheTargets
}

func prepareSearch(search string, options Options) []searchTerm {
	searchTerms := strings.Split(normalize(search, options), " ")
	result := make([]searchTerm, len(searchTerms))
	for x, term := range searchTerms {
		result[x].term = strings.Split(term, "")
		result[x].len = len(term)
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

	targetLen := len(*cacheTargets)
	results := make([]algorithmResult, targetLen)
	resultWrapper := SearchResult{Results: results, BestScore: 0}

	for i, cacheTarget := range *cacheTargets {
		result := algorithmResult{Target: cacheTarget.target, Score: 0, Typos: 0, MatchCount: 0}
		algorithm(preparedSearch, cacheTarget, &result, options)
		accurateTokens := cacheTarget.len - result.MatchCount
		if result.Typos < accurateTokens {
			result.Score++
		}
		if result.Typos < accurateTokens/2 {
			result.Score++
		}
		if result.Typos < accurateTokens*2 {
			result.Score++
		}
		if result.Typos == 0 && result.MatchCount > 0 {
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
	if options.Limit > 0 && options.Limit < targetLen {
		resultWrapper.Results = make([]algorithmResult, options.Limit)
		for i, result := range results[:options.Limit] {
			resultWrapper.Results[i] = result
		}
	}
	return &resultWrapper
}

func algorithm(searchTerms []searchTerm, target CacheTarget, result *algorithmResult, options Options) {

	for _, term := range searchTerms {
		searchI := 0
		hasFirstMatch := false
		currentTermMatchCount := 0
		currentTermTypos := 0
		for targetI := 0; targetI < target.len; targetI++ {
			if term.term[searchI] == target.cache[targetI] {
				hasFirstMatch = true
				currentTermMatchCount++
				result.MatchCount++
				result.Score++
				searchI++
			} else {
				if hasFirstMatch {
					currentTermTypos++
					result.Typos++
				}
			}
			if searchI == term.len {
				break
			}
		}
		if currentTermMatchCount == term.len {
			result.Score += 2
		}
		if currentTermTypos == 0 && currentTermMatchCount > 0 {
			result.Score++
		}
	}

}

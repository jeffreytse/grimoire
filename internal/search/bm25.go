// Package search provides BM25-based skill retrieval.
// Skills are scored against a free-text query using the BM25 ranking function,
// which handles lexical matching better than simple substring search
// (e.g. "write better commit messages" matches propose-conventional-commit).
package search

import (
	"math"
	"sort"
	"strings"

	"github.com/jeffreytse/grimoire/internal/skills"
)

const (
	bm25K1 = 1.2
	bm25B  = 0.75
)

// Result is a skill with its BM25 relevance score.
type Result struct {
	Skill skills.Skill
	Score float64
}

// BM25Search scores corpus skills against query and returns top-N results sorted
// by score descending. If topN <= 0, all non-zero-score documents are returned.
func BM25Search(query string, corpus []skills.Skill, topN int) []Result {
	if len(corpus) == 0 || strings.TrimSpace(query) == "" {
		return nil
	}

	tokens := tokenize(query)
	if len(tokens) == 0 {
		return nil
	}

	// Document lengths and average
	docTexts := make([][]string, len(corpus))
	docLen := make([]float64, len(corpus))
	for i := range corpus {
		docTexts[i] = tokenize(docText(&corpus[i]))
		docLen[i] = float64(len(docTexts[i]))
	}
	avgDL := avgFloat(docLen)
	if avgDL == 0 {
		avgDL = 1
	}

	// Document frequency per token
	df := make(map[string]int)
	for _, toks := range docTexts {
		seen := make(map[string]bool, len(toks))
		for _, t := range toks {
			if !seen[t] {
				df[t]++
				seen[t] = true
			}
		}
	}

	N := float64(len(corpus))
	results := make([]Result, 0, len(corpus))

	for i := range corpus {
		sk := &corpus[i]
		tf := termFreq(docTexts[i])
		score := 0.0
		for _, t := range tokens {
			if tf[t] == 0 {
				continue
			}
			idf := math.Log((N-float64(df[t])+0.5)/(float64(df[t])+0.5) + 1)
			num := tf[t] * (bm25K1 + 1)
			den := tf[t] + bm25K1*(1-bm25B+bm25B*docLen[i]/avgDL)
			score += idf * num / den
		}
		if score > 0 {
			results = append(results, Result{Skill: *sk, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Skill.Name < results[j].Skill.Name
	})

	if topN > 0 && len(results) > topN {
		results = results[:topN]
	}
	return results
}

// docText concatenates all searchable text fields for a skill.
func docText(sk *skills.Skill) string {
	parts := make([]string, 0, 5+len(sk.Tags))
	parts = append(parts, sk.Name, sk.Description, sk.Domain, sk.Subdomain)
	parts = append(parts, sk.Tags...)
	return strings.Join(parts, " ")
}

// tokenize lowercases, splits on whitespace and non-alpha boundaries, and
// expands hyphen/underscore-separated terms so "propose-conventional-commit"
// matches both "propose" and "conventional" and "commit".
func tokenize(s string) []string {
	s = strings.ToLower(s)
	s = strings.NewReplacer("-", " ", "_", " ", "/", " ", ".", " ").Replace(s)
	fields := strings.Fields(s)
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if len(f) > 1 { // skip single chars
			out = append(out, f)
		}
	}
	return out
}

func termFreq(tokens []string) map[string]float64 {
	m := make(map[string]float64, len(tokens))
	for _, t := range tokens {
		m[t]++
	}
	return m
}

func avgFloat(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	var s float64
	for _, x := range xs {
		s += x
	}
	return s / float64(len(xs))
}

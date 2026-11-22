// Package search provides BM25-based skill retrieval.
// Skills are scored against a free-text query using the BM25 ranking function,
// which handles lexical matching better than simple substring search
// (e.g. "write better commit messages" matches propose-conventional-commit).
package search

import (
	"math"
	"sort"
	"strings"
	"unicode"

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
	// Matched is how many distinct query terms the skill contains; Terms is
	// how many distinct terms the query has after stopword removal.
	Matched, Terms int
}

// BM25Search scores corpus skills against query and returns top-N results sorted
// by score descending. If topN <= 0, all non-zero-score documents are returned.
func BM25Search(query string, corpus []skills.Skill, topN int) []Result {
	if len(corpus) == 0 || strings.TrimSpace(query) == "" {
		return nil
	}

	tokens := uniqueTokens(tokenize(query))
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
		matched := 0
		for _, t := range tokens {
			if tf[t] == 0 {
				continue
			}
			matched++
			idf := math.Log((N-float64(df[t])+0.5)/(float64(df[t])+0.5) + 1)
			num := tf[t] * (bm25K1 + 1)
			den := tf[t] + bm25K1*(1-bm25B+bm25B*docLen[i]/avgDL)
			score += idf * num / den
		}
		if score > 0 {
			results = append(results, Result{Skill: *sk, Score: score, Matched: matched, Terms: len(tokens)})
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

// tokenize lowercases and splits on anything that is not a letter or digit, so
// "propose-conventional-commit" yields "propose", "conventional", "commit" and
// "workflows?" yields "workflow". Single characters and stopwords are dropped,
// and each term is lightly stemmed so "tests"/"testing" match "test".
func tokenize(s string) []string {
	fields := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if len(f) <= 1 || stopwords[f] {
			continue
		}
		out = append(out, stem(f))
	}
	return out
}

func uniqueTokens(tokens []string) []string {
	seen := make(map[string]bool, len(tokens))
	out := tokens[:0:0]
	for _, t := range tokens {
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

// stemSuffixes are stripped longest-first; a stem must keep minStem letters.
// Crude next to Porter, but applied to queries and skills alike, which is all
// matching needs: "renovating" and "renovation" both become "renov".
var stemSuffixes = []struct{ suffix, repl string }{
	{"ations", ""}, {"ation", ""}, {"ating", ""}, {"ated", ""}, {"izing", ""},
	{"izes", ""}, {"ized", ""}, {"ize", ""}, {"ies", "y"}, {"ied", "y"},
	{"ing", ""}, {"ed", ""}, {"es", ""}, {"ate", ""}, {"s", ""},
}

const minStem = 3

func stem(w string) string {
	if strings.HasSuffix(w, "ss") || strings.HasSuffix(w, "us") || strings.HasSuffix(w, "is") {
		return w // class, status, analysis
	}
	for _, s := range stemSuffixes {
		if strings.HasSuffix(w, s.suffix) && len(w)-len(s.suffix)+len(s.repl) >= minStem {
			return w[:len(w)-len(s.suffix)] + s.repl
		}
	}
	return w
}

// stopwords carry no topic: function words plus conversational filler common
// in agent prompts ("please help me ...", "can you ...").
var stopwords = func() map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(`
		a about above after again against all also am an and any are as at be
		because been before being below between both but by can could did do
		does doing done down during each few for from further had has have
		having he her here hers him his how however if in into is it its itself
		just me more most my myself no nor not now of off on once only or other
		our ours out over own same she should so some such than that the their
		them then there these they this those through to too under until up
		upon us very was we were what when where which while who whom why will
		with would you your yours yourself
		im ive id ill youre dont doesnt didnt cant wont isnt arent
		please thanks thank ok okay yes yeah sure hi hello hey lol nice cool
		help want wanted need needs let lets get got make made go goes going
		ahead like well good great really actually maybe also whether way thing
		things something anything stuff kind sort bit lot lots one ones see
		look looks tell give show try know think said say says now latest new
		use using used able`) {
		m[w] = true
	}
	return m
}()

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

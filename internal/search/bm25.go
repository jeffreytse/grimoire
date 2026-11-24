// Package search provides BM25-based skill retrieval.
// Skills are scored against a free-text query using the BM25 ranking function,
// which handles lexical matching better than simple substring search
// (e.g. "write better commit messages" matches propose-conventional-commit).
package search

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/jeffreytse/grimoire/internal/skills"
)

const (
	bm25K1 = 1.2
	bm25B  = 0.75
)

// field is one searchable part of a skill. BM25F (Robertson & Zaragoza 2009)
// weights each field's term frequency, normalized by that field's own average
// length, before saturation. A skill's name and tags state what it is for;
// its description also carries incidental words ("in the 24-48 hours before a
// competition"), which plain BM25 over the concatenated text scored as highly
// as the topic itself.
type field struct {
	weight float64
	text   func(sk *skills.Skill) string
}

var fields = []field{
	{3, func(sk *skills.Skill) string { return sk.Name }},
	{2, func(sk *skills.Skill) string { return strings.Join(sk.Tags, " ") }},
	{1, func(sk *skills.Skill) string { return sk.Description }},
	{1, func(sk *skills.Skill) string { return sk.Domain + " " + sk.Subdomain }},
}

// BodyWeight is the BM25F weight of a skill's body (its full SKILL.md text)
// when Options.Body is set. The body holds the words practitioners use for
// a problem ("rollback", "go-live") that the name and tags often lack, but
// also many incidental words, so it counts for less than the description.
var BodyWeight = 0.3

// bodyField is the optional body field.
var bodyField = field{0, func(sk *skills.Skill) string { return sk.Body }}

// Options narrow and widen a search.
type Options struct {
	// Body also matches against each skill's full text (BodyWeight).
	Body bool
	// Domain and Subdomain, when set, keep only skills there.
	Domain, Subdomain string
	// Tags, when set, keep only skills carrying every one of them.
	Tags []string
}

// Result is a skill with its BM25 relevance score.
type Result struct {
	Skill skills.Skill
	Score float64
	// Matched is how many distinct query terms the skill contains; Terms is
	// how many distinct terms the query has after stopword removal.
	Matched, Terms int
}

// BM25Search scores corpus skills against query with BM25F and returns top-N
// results sorted by score descending. If topN <= 0, all non-zero-score
// documents are returned.
func BM25Search(query string, corpus []skills.Skill, topN int) []Result {
	return Search(query, corpus, topN, Options{})
}

// Search is BM25Search with options: filters applied before scoring, and the
// skill body as an extra field.
func Search(query string, corpus []skills.Skill, topN int, opts Options) []Result {
	corpus = filterCorpus(corpus, opts)
	fields := fields
	if opts.Body {
		body := bodyField
		body.weight = BodyWeight
		fields = append(append([]field{}, fields...), body)
	}
	if len(corpus) == 0 || strings.TrimSpace(query) == "" {
		return nil
	}
	tokens := uniqueTokens(tokenize(query))
	if len(tokens) == 0 {
		return nil
	}

	// Per-field term frequencies and lengths, and per-field average lengths.
	n := len(corpus)
	tfs := make([][]map[string]float64, n)
	lens := make([][]float64, n)
	avg := make([]float64, len(fields))
	df := make(map[string]int)
	for i := range corpus {
		tfs[i] = make([]map[string]float64, len(fields))
		lens[i] = make([]float64, len(fields))
		seen := map[string]bool{}
		for f := range fields {
			toks := tokenize(fields[f].text(&corpus[i]))
			tfs[i][f] = termFreq(toks)
			lens[i][f] = float64(len(toks))
			avg[f] += float64(len(toks))
			for _, t := range toks {
				if !seen[t] {
					seen[t] = true
					df[t]++
				}
			}
		}
	}
	for f := range avg {
		avg[f] /= float64(n)
		if avg[f] == 0 {
			avg[f] = 1
		}
	}

	N := float64(n)
	idf := make(map[string]float64, len(tokens))
	for _, t := range tokens {
		idf[t] = math.Log((N-float64(df[t])+0.5)/(float64(df[t])+0.5) + 1)
	}

	results := make([]Result, 0, n)
	for i := range corpus {
		score := 0.0
		matched := 0
		for _, t := range tokens {
			wtf := 0.0
			for f := range fields {
				if tf := tfs[i][f][t]; tf > 0 {
					wtf += fields[f].weight * tf / (1 - bm25B + bm25B*lens[i][f]/avg[f])
				}
			}
			if wtf == 0 {
				continue
			}
			matched++
			score += idf[t] * wtf * (bm25K1 + 1) / (wtf + bm25K1)
		}
		if score > 0 {
			results = append(results, Result{Skill: corpus[i], Score: score, Matched: matched, Terms: len(tokens)})
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

// filterCorpus keeps the skills matching opts' domain, subdomain and tags.
func filterCorpus(corpus []skills.Skill, opts Options) []skills.Skill {
	if opts.Domain == "" && opts.Subdomain == "" && len(opts.Tags) == 0 {
		return corpus
	}
	out := make([]skills.Skill, 0, len(corpus))
next:
	for i := range corpus {
		sk := &corpus[i]
		if opts.Domain != "" && !strings.EqualFold(sk.Domain, opts.Domain) {
			continue
		}
		if opts.Subdomain != "" && !strings.EqualFold(sk.Subdomain, opts.Subdomain) {
			continue
		}
		for _, want := range opts.Tags {
			found := false
			for _, tag := range sk.Tags {
				if strings.EqualFold(tag, want) {
					found = true
					break
				}
			}
			if !found {
				continue next
			}
		}
		out = append(out, *sk)
	}
	return out
}

// tokenize lowercases and splits on anything that is not a letter or digit, so
// "propose-conventional-commit" yields "propose", "conventional", "commit" and
// "workflows?" yields "workflow". Single characters and stopwords are dropped,
// and each term is lightly stemmed so "tests"/"testing" match "test".
// slashAcronym matches short slash-joined acronyms — A/B, CI/CD, I/O — which
// splitting would reduce to letters too short to keep.
var slashAcronym = regexp.MustCompile(`\b([a-z0-9]{1,3})/([a-z0-9]{1,3})\b`)

func tokenize(s string) []string {
	s = slashAcronym.ReplaceAllString(strings.ToLower(s), "$1$2")
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		// Bare numbers ("48", "10") are quantities, not topics; they matched
		// whatever skill happened to quote a number. Mixed tokens ("oauth2",
		// "s3", "h264") are kept.
		if len(f) <= 1 || stopwords[f] || isNumber(f) {
			continue
		}
		out = append(out, stem(f))
	}
	return out
}

func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
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
		them then there these they this those through to too under until
		upon us very was we were what when where which while who whom why will
		with would you your yours yourself
		im ive id ill youre dont doesnt didnt cant wont isnt arent
		please thanks thank ok okay yes yeah sure hi hello hey lol nice cool
		help want wanted need needs let lets get got make made go goes going
		ahead like well good great really actually maybe also whether way thing
		things something anything stuff kind sort bit lot lots see
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

// rrfK is the reciprocal rank fusion constant — the standard 60 (Cormack et
// al., 2009); larger values flatten the advantage of top ranks.
const rrfK = 60

// rrfDepth is how deep each ranking goes into the fusion.
const rrfDepth = 50

// HybridSearch ranks corpus for the queries by fusing, with reciprocal rank
// fusion, the field-weighted ranking (name, tags, description) and the same
// ranking with skill bodies added. The body finds skills whose practitioner
// vocabulary lives in their text; fusion keeps the precise field ranking on
// top. Several queries (phrasings of one need) are fused the same way.
// Scores are fused ranks, comparable only within one call.
//
// Measured on testdata/hook_eval.json (54 labeled prompts) against
// field-only BM25: top-1 39→40, hit@3 45→48, hit@10 49→52, MRR .790→.812.
func HybridSearch(queries []string, corpus []skills.Skill, topN int, opts Options) []Result {
	type acc struct {
		res   Result
		score float64
	}
	fused := map[string]*acc{}
	var order []string
	add := func(rank []Result) {
		for i := range rank {
			name := rank[i].Skill.Name
			a, ok := fused[name]
			if !ok {
				a = &acc{res: rank[i]}
				fused[name] = a
				order = append(order, name)
			}
			a.score += 1 / float64(rrfK+i+1)
			if rank[i].Matched > a.res.Matched {
				a.res.Matched, a.res.Terms = rank[i].Matched, rank[i].Terms
			}
		}
	}
	for _, q := range queries {
		if strings.TrimSpace(q) == "" {
			continue
		}
		plain := opts
		plain.Body = false
		add(Search(q, corpus, rrfDepth, plain))
		body := opts
		body.Body = true
		add(Search(q, corpus, rrfDepth, body))
	}
	out := make([]Result, 0, len(order))
	for _, name := range order {
		a := fused[name]
		a.res.Score = a.score
		out = append(out, a.res)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Skill.Name < out[j].Skill.Name
	})
	if topN > 0 && len(out) > topN {
		out = out[:topN]
	}
	return out
}

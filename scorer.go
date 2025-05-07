package assessment

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type Scorer interface {
	Calculate([]Answer) int
}

type ScorerFunc func([]Answer) int

func (sf ScorerFunc) Calculate(a []Answer) int { return sf(a) }

func Rule(fn MatchFunc, is int) Scorer { return count(fn, is) }

type MatchFunc func(Answer) bool

func count(match MatchFunc, matchValue int) ScorerFunc {
	return func(a []Answer) int {
		score := 0
		for _, answer := range a {
			if match(answer) {
				score += matchValue
			}
		}
		return score
	}
}

func matchAll(all ...MatchFunc) MatchFunc {
	return func(a Answer) bool {
		for _, fn := range all {
			if !fn(a) {
				return false
			}
		}

		return true
	}
}

func matchAny(any ...MatchFunc) MatchFunc {
	return func(a Answer) bool {
		for _, fn := range any {
			if fn(a) {
				return true
			}
		}

		return false
	}
}

func matchQuestionNumbers(n ...int) MatchFunc {
	return func(a Answer) bool { return slices.Contains(n, a.Number) }
}

func matchYes(a Answer) bool { return a.Response }
func matchNo(a Answer) bool  { return !a.Response }

func ParseScorer(name string) (Scorer, error) {
	r := GetReader(name)
	if r == nil {
		return nil, fmt.Errorf("unable to open location %q", name)
	}

	defer func() { _ = r.Close() }()

	return parseScorer(r)
}

func parseScorer(from io.Reader) (Scorer, error) {
	scorer := &ruleScorer{}
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		if err := scorer.Parse(scanner.Text()); err != nil {
			return scorer, err
		}
	}
	return scorer, scanner.Err()
}

//
// ruleScorer
//

type ruleScorer struct {
	rules []Scorer
}

func (s *ruleScorer) Calculate(a []Answer) int {
	score := 0

	for _, rule := range s.rules {
		score += rule.Calculate(a)
	}

	return score
}

/*
sum
all:yes:1
matching 7,8,12,13,17,21,22,26,27,35,40,43,44,47,49,60:yes:5
matching 2,18,25,32,50,54:yes:2
matching 4,9:no:2
*/

func (s *ruleScorer) Parse(line string) error {
	tok := strings.Split(line, ":")
	switch len(tok) {
	case 3:
		return s.parseTriple(tok)
	default:
		return nil
	}
}

func (s *ruleScorer) parseTriple(tokens []string) error {
	filter, response, weightStr := tokens[0], tokens[1], tokens[2]
	weight, err := strconv.Atoi(weightStr)
	if err != nil {
		return fmt.Errorf("parsing triple: weight %w", err)
	}

	var rFn MatchFunc
	switch response {
	case "yes":
		rFn = matchYes
	case "no":
		rFn = matchNo
	default:
		return fmt.Errorf("unconfigured response %q", response)
	}

	filterFn, err := s.parseFilter(filter)
	if err != nil {
		return fmt.Errorf("parsing triple: filter %w", err)
	}

	s.rules = append(s.rules, Rule(matchAll(filterFn, rFn), weight))

	return nil
}

func (s *ruleScorer) parseFilter(filter string) (MatchFunc, error) {
	segments := strings.Split(filter, " ")
	if len(segments) == 0 {
		return nil, fmt.Errorf("parse filter: empty")
	}

	switch segments[0] {
	case "all":
		return matchAll(), nil
	case "matching":
		if len(segments) != 2 {
			return nil, fmt.Errorf("parse filter: unknown filter args: %+v", segments)
		}
		items := strings.Split(segments[1], ",")
		numbers := make([]int, len(items))
		for ix, s := range items {
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("parse filter: expecting a number, got %q: %w", s, err)
			}
			numbers[ix] = v
		}

		return matchQuestionNumbers(numbers...), nil

	default:
		return nil, fmt.Errorf("unknown filter %q", segments[0])
	}
}

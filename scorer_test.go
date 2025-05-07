package assessment

import (
	"testing"
)

func TestParseRules(t *testing.T) {
	for ix, str := range []string{"testdata/rules.txt", "testdata/alt-rules.txt"} {
		s, err := ParseScorer(str)
		if err != nil {
			t.Errorf("%d: %v", ix, err)
		}
		if s == nil {
			t.Errorf("%d: expected a scorer", ix)
		}
		t.Logf("results: %+v", s)
	}
}

func TestScorer(t *testing.T) {
	type testcase struct {
		questions string
		responses string
		rules     string
		expect    int
	}

	for name, tc := range map[string]testcase{
		"response1":    {questions: "testdata/questions.txt", responses: "testdata/answers1.txt", rules: "testdata/rules.txt", expect: 41},
		"response2":    {questions: "testdata/questions.txt", responses: "testdata/answers2.txt", rules: "testdata/rules.txt", expect: 21},
		"response2alt": {questions: "testdata/questions.txt", responses: "testdata/answers2.txt", rules: "testdata/alt-rules.txt", expect: 11},
	} {
		t.Run(name, func(t *testing.T) {
			scorer, err := ParseScorer(tc.rules)
			if err != nil {
				t.Fatal(err)
			}

			questions, err := ParseQuestions(tc.questions)
			if err != nil {
				t.Fatal(err)
			}

			responses, err := ParseAnswers(tc.responses)
			if err != nil {
				t.Fatal(err)
			}

			if len(questions) != len(responses) {
				t.Fatalf("expected length of questions (%d) and responses (%d) to match", len(questions), len(responses))
			}

			if score, expect := scorer.Calculate(responses), tc.expect; score != expect {
				t.Errorf("expected score: %d, got %d", expect, score)
			}
		})
	}
}

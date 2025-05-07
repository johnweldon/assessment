package assessment

import (
	"testing"
)

func TestParseQuestions(t *testing.T) {
	for ix, str := range []string{"1 question one\n 2. two questions\n 03: tres", "testdata/questions.txt"} {
		q, err := ParseQuestions(str)
		if err != nil {
			t.Errorf("%d: %v", ix, err)
		}
		if len(q) == 0 {
			t.Errorf("%d: expected some questions", ix)
		}
		t.Logf("results: %+v", q)
	}
}

func TestParseAnswers(t *testing.T) {
	for ix, str := range []string{"1 yes\n 2. n\n 03: t", "testdata/answers1.txt", "testdata/answers2.txt"} {
		a, err := ParseAnswers(str)
		if err != nil {
			t.Errorf("%d: %v", ix, err)
		}
		if len(a) == 0 {
			t.Errorf("%d: expected some answers", ix)
		}
		t.Logf("results: %+v", a)
	}
}

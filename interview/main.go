package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/johnweldon/assessment"
)

var (
	ErrHelp       = errors.New("help")
	ErrValidation = errors.New("validation error")
	rules         = flag.String("rules", "", "file or url for rules (ASSESSMENT_RULES)")
	qsrc          = flag.String("questions", "", "file or url of questions (ASSESSMENT_QUESTIONS)")
)

func initialize() (questions []assessment.Question, scorer assessment.Scorer, asker assessment.Asker, err error) {
	flag.Set("rules", os.Getenv("ASSESSMENT_RULES"))
	flag.Set("questions", os.Getenv("ASSESSMENT_QUESTIONS"))
	flag.Parse()

	if rules == nil || *rules == "" {
		err = fmt.Errorf("must specify location of rules: %w", ErrValidation)
		return
	}

	if scorer, err = assessment.ParseScorer(*rules); err != nil {
		err = fmt.Errorf("unable to create scorer: %w", errors.Join(err, ErrValidation))
		return
	}

	if qsrc == nil || *qsrc == "" {
		err = fmt.Errorf("must specify location of questions: %w", ErrValidation)
		return
	}

	if questions, err = assessment.ParseQuestions(*qsrc); err != nil {
		err = fmt.Errorf("unable to load questions: %w", errors.Join(err, ErrValidation))
		return
	}

	if assessment.IsInteractive(os.Stdin) {
		asker = assessment.NewBubbleTeaAsker()
	} else {
		asker = assessment.NewConsoleAsker(os.Stdin, os.Stdout)
	}

	return
}

func main() {
	questions, scorer, asker, err := initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "initialization failure: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	answers := make([]assessment.Answer, len(questions))
	for ix, q := range questions {
		a, err := asker.Ask(q)
		if err != nil {
			log.Fatalf("quitting: %v", err)
		}
		answers[ix] = a
		fmt.Printf("  > score: %d\n", scorer.Calculate(answers))
	}

	fmt.Printf("\nFinal score: %d\n", scorer.Calculate(answers))
}

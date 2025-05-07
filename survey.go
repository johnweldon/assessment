package assessment

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

type Question struct {
	Number int
	Text   string
}

type Answer struct {
	Number   int
	Response bool
}

type Asker interface {
	Ask(Question) (Answer, error)
}

type AskerFunc func(Question) (Answer, error)

func (a AskerFunc) Ask(q Question) (Answer, error) { return a(q) }

func ParseQuestions(name string) ([]Question, error) {
	r := GetReader(name)
	if r == nil {
		return nil, fmt.Errorf("unable to open location %q", name)
	}

	defer func() { _ = r.Close() }()

	return parseQuestions(r)
}

func ParseAnswers(name string) ([]Answer, error) {
	r := GetReader(name)
	if r == nil {
		return nil, fmt.Errorf("unable to open location %q", name)
	}

	defer func() { _ = r.Close() }()

	return parseAnswers(r)
}

var questionRE = regexp.MustCompile(`^[[:space:]]*([[:digit:]]+)[.: ]*(.*?)[[:space:]]*$`)

func parseQuestions(r io.Reader) ([]Question, error) {
	questions := []Question{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		parts := questionRE.FindStringSubmatch(sc.Text())
		if len(parts) == 3 {
			if num, err := strconv.Atoi(parts[1]); err != nil {
				log.Printf("WARN: unable to parse %q: %v", parts[1], err)
			} else {
				questions = append(questions, Question{Number: num, Text: parts[2]})
			}
		}
	}
	return questions, sc.Err()
}

var answerRE = regexp.MustCompile(`^[[:space:]]*([[:digit:]]+)[.: ]*(.*?)[[:space:]]*$`)

func parseAnswers(r io.Reader) ([]Answer, error) {
	answers := []Answer{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		parts := answerRE.FindStringSubmatch(sc.Text())
		if len(parts) == 3 {
			if num, err := strconv.Atoi(parts[1]); err != nil {
				log.Printf("WARN: unable to parse number %q: %v", parts[1], err)
			} else {
				switch parts[2] {
				case "y", "Y", "yes", "Yes", "YES":
					answers = append(answers, Answer{Number: num, Response: true})
				case "n", "N", "no", "No", "NO":
					answers = append(answers, Answer{Number: num, Response: false})
				default:
					if resp, err := strconv.ParseBool(parts[2]); err != nil {
						log.Printf("WARN: unable to parse bool %q: %v", parts[2], err)
					} else {
						answers = append(answers, Answer{Number: num, Response: resp})
					}
				}
			}
		}
	}
	return answers, sc.Err()
}

package assessment

import (
	"bufio"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func IsInteractive(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func NewConsoleAsker(r io.Reader, w io.Writer) Asker {
	return &consoleAsker{sc: bufio.NewScanner(r), write: w}
}

type consoleAsker struct {
	write io.Writer
	sc    *bufio.Scanner
}

func (a *consoleAsker) Ask(q Question) (Answer, error) {
	_, _ = fmt.Fprintf(a.write, "%2d. %s. (y|n|q)\n", q.Number, q.Text)
	for a.sc.Scan() {
		r := a.sc.Bytes()
		if len(r) == 0 {
			continue
		}

		switch r[0] {
		case 'y', 'Y':
			return Answer{Number: q.Number, Response: true}, nil
		case 'n', 'N':
			return Answer{Number: q.Number, Response: false}, nil
		case 'q', 'Q', 'x', 'X':
			return Answer{}, fmt.Errorf("quit")
		default:
			continue
		}
	}

	return Answer{}, a.sc.Err()
}

type bbtModel struct {
	q        Question
	response bool
	yes      bool
}

var _ tea.Model = (*bbtModel)(nil)

func (m *bbtModel) Init() tea.Cmd { return nil }

func (m *bbtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			m.response = true
			m.yes = true
			return m, tea.Quit
		case "n":
			m.response = true
			m.yes = false
			return m, tea.Quit
		case "q":
			m.response = false
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *bbtModel) View() string { return fmt.Sprintf("%2d. %s (y|n|q)", m.q.Number, m.q.Text) }

func NewBubbleTeaAsker() Asker { return &bbtAsker{} }

type bbtAsker struct{}

func (bbtAsker) Ask(q Question) (Answer, error) {
	model := &bbtModel{q: q}
	if _, err := tea.NewProgram(model).Run(); err != nil {
		return Answer{}, fmt.Errorf("error running asker: (%v) %v", model, err)
	}

	if model.response {
		return Answer{Number: q.Number, Response: model.yes}, nil
	}

	return Answer{}, fmt.Errorf("quit, no response")
}

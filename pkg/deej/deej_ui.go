package deej

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 40
)

type DeejUi struct {
	deej       *Deej
	tea        *tea.Program
	progresses map[int]float32
	progress   progress.Model

	input chan SliderMoveEvent
}

func newDeejUi(deej *Deej) *DeejUi {
	ui := &DeejUi{
		deej:       deej,
		input:      make(chan SliderMoveEvent),
		progresses: map[int]float32{},
		progress:   progress.New(progress.WithDefaultGradient()),
	}

	return ui
}

func (du *DeejUi) Run() {
	for sliderIdxString := range du.deej.config.SliderMapping.m {
		du.progresses[sliderIdxString] = 0
	}
	go func() {
		for event := range du.input {
			if du.tea != nil {
				du.tea.Send(event)
			}
		}
	}()
	du.tea = tea.NewProgram(du)
	if _, err := du.tea.Run(); err != nil {
		fmt.Println("Oh no!", err)
		close(du.deej.stopChannel)
	}
}

func (du *DeejUi) Init() tea.Cmd {
	return nil
}

func (du *DeejUi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		//for _, progress := range du.progresses {
		//	progress.Width = msg.Width - padding*2 - 4
		//	if progress.Width > maxWidth {
		//		progress.Width = maxWidth
		//	}
		//}
		du.progress.Width = msg.Width - padding*2 - 40
		if du.progress.Width > maxWidth {
			du.progress.Width = maxWidth
		}
		return du, nil

	case SliderMoveEvent:
		du.progresses[msg.SliderID] = msg.PercentValue
		return du, nil

	// FrameMsg is sent when the progress bar wants to animate itself
	//case progress.FrameMsg:
	//	var cmds []tea.Cmd
	//	for _, model := range du.progresses {
	//		progressModel, c := model.Update(msg)
	//		model = progressModel.(progress.Model)
	//		cmds = append(cmds, c)
	//	}
	//	progressModel, cmd := du.progress.Update(msg)
	//	du.progress = progressModel.(progress.Model)
	//	cmds = append(cmds, cmd)
	//	return du, tea.Batch(cmds...)
	default:
		return du, nil
	}
}

func (du *DeejUi) View() string {
	pad := strings.Repeat(" ", padding)
	ss := "\n"
	sessions, _ := du.deej.sessions.sessionFinder.GetAllSessions()
	for _, session := range sessions {
		name := session.Key()
		if len(name) > 10 {
			name = name[:10] + "..."
		}
		ss += pad + name + "\t" + pad + RenderMute(session.GetMute()) + pad + du.progress.ViewAs(float64(session.GetVolume())) + "\n\n"
	}
	return ss
}

func RenderMute(t bool) string {
	if t {
		return "✖"
	}
	return "✔"
}

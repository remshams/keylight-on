package hue_group_scenes

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"
	hue_groups "keylight-charm/pages/hue/groups"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var scenesStyle = lipgloss.NewStyle().Margin(1, 2)

type sceneItem struct {
	title string
}

func (sceneItem sceneItem) Title() string {
	return sceneItem.title
}

func (sceneItem sceneItem) FilterValue() string {
	return sceneItem.title
}

type Model struct {
	adapter *hue.HueAdapter
	group   hue_control.Group
	scenes  list.Model
}

func InitModel(adapter *hue.HueAdapter, group hue_control.Group) Model {
	return Model{
		adapter: adapter,
		group:   group,
		scenes:  createScenes(group.GetScenes()),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := scenesStyle.GetFrameSize()
		m.scenes.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = hue_groups.CreateBackToGroupDetailsAction()
		default:
			m.scenes, cmd = m.scenes.Update(msg)
		}
	default:
		m.scenes, cmd = m.scenes.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return scenesStyle.Render(m.scenes.View())
}

func createScenes(scenes []hue_control.Scene) list.Model {
	var items []list.Item
	for _, scene := range scenes {
		items = append(items, sceneItem{title: scene.Name()})
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Scenes"
	return list
}

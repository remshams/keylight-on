package home

import (
	"keylight-charm/lights/hue"
	"keylight-charm/lights/keylight"
	"keylight-charm/pages"
	hue_home "keylight-charm/pages/hue/home"
	keylight_home "keylight-charm/pages/keylight/home"
	"keylight-charm/stores"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var menuStyle = lipgloss.NewStyle().Margin(1, 2)

type menuItem struct {
	title, desc string
}

func (menuItem menuItem) Title() string {
	return menuItem.title
}

func (menuItem menuItem) Description() string {
	return menuItem.desc
}

func (menuItem menuItem) FilterValue() string { return menuItem.title }

type viewState string

const (
	menu      viewState = "menu"
	keylights viewState = "keylights"
	hueLights viewState = "hueLights"
)

type Model struct {
	keylight keylight_home.Model
	hue      hue_home.Model
	menu     list.Model
	state    viewState
}

func InitModel(keylightAdapter *keylight.KeylightAdapter, hueAdapter *hue.HueAdapter) Model {
	return Model{
		keylight: keylight_home.InitModel(keylightAdapter),
		hue:      hue_home.InitModel(hueAdapter),
		menu:     createMenu(),
		state:    menu,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if pages.IsSystemMsg(msg) {
		cmd = m.processSystemUpdate(msg)
	} else {
		switch m.state {
		case menu:
			cmd = m.processMenuUpdate(msg)
		case keylights:
			cmd = m.processKeylightsUpdate(msg)
		case hueLights:
			cmd = m.processHueUpate(msg)
		}
	}
	return m, cmd
}

func (m *Model) processSystemUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		stores.LayoutStore.Width = msg.Width
		stores.LayoutStore.Height = msg.Height
		cmd = pages.CreateWindowResizeAction(msg.Width, msg.Height)
	}
	return cmd
}

func (m *Model) processMenuUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.WindowResizeAction:
		h, v := menuStyle.GetFrameSize()
		m.menu.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch m.state {
		case menu:
			switch msg.String() {
			case "ctrl+c", "q":
				cmd = tea.Quit
			case "enter":
				if m.menu.Index() == 0 {
					m.state = keylights
					cmd = m.keylight.Init()
				} else {
					m.state = hueLights
					cmd = m.hue.Init()
				}
			default:
				m.menu, cmd = m.menu.Update(msg)
			}
		}
	default:
		m.menu, cmd = m.menu.Update(msg)
	}
	return cmd
}

func (m *Model) processKeylightsUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.BackToMenuAction:
		m.state = menu
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			m.keylight, cmd = m.keylight.Update(msg)
		}
	default:
		m.keylight, cmd = m.keylight.Update(msg)
	}
	return cmd
}

func (m *Model) processHueUpate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.BackToMenuAction:
		m.state = menu
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			m.hue, cmd = m.hue.Update(msg)
		}
	default:
		m.hue, cmd = m.hue.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case menu:
		return menuStyle.Render(m.menu.View())
	case keylights:
		return m.keylight.View()
	case hueLights:
		return m.hue.View()
	default:
		return ""
	}
}

func createMenu() list.Model {
	items := []list.Item{
		menuItem{title: "Keylights", desc: "Control keylights"},
		menuItem{title: "HueLights", desc: "Control huelights"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Lights"
	return list
}

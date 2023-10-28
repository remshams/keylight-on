package keylight_list

import (
	"keylight-charm/components/toast"
	"keylight-charm/keylight"
	"keylight-charm/utils"
	"keylight-control/control"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SelectedKeylight struct {
	Keylight *control.Keylight
}

func createSelectedKeylightAction(keylight *control.Keylight) tea.Cmd {
	return func() tea.Msg {
		return SelectedKeylight{keylight}
	}
}

type AddKeylight struct{}

type EditKeylight struct {
	Keylight *control.Keylight
}

type RemoveKeylight struct {
	Keylight *control.Keylight
}

type ReloadKeylights struct{}

type viewState string

type Model struct {
	keylights       []control.Keylight
	keylightAdapter *keylight.KeylightAdapter
	table           table.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter, keylights []control.Keylight) Model {
	model := Model{keylights: keylights, keylightAdapter: keylightAdapter, table: createTable(keylights)}
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd = m.selectedKeylight(m.table.SelectedRow()[0])
		case "a":
			cmd = m.addNewKeylight()
		case "e":
			cmd = m.editKeylight(m.table.SelectedRow()[0])
		case "d":
			cmd = m.removeKeylight(m.table.SelectedRow()[0])
		case "r":
			cmd = m.reloadKeylights()
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}

	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func createTable(keylights []control.Keylight) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Ip", Width: 15},
		{Title: "Port", Width: 5},
		{Title: "Connected", Width: 10},
	}
	rows := []table.Row{}

	for _, keylight := range keylights {
		rows = append(
			rows,
			table.Row{
				strconv.Itoa(keylight.Metadata.Id),
				keylight.Metadata.Name,
				keylight.Metadata.Ip.String(),
				strconv.Itoa(keylight.Metadata.Port),
				strconv.FormatBool(keylight.Metadata.Connected)},
		)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func (m *Model) selectedKeylight(keylightId string) tea.Cmd {
	keylight := utils.FindKeylightWithId(m.keylights, keylightId)
	if keylight != nil && keylight.Metadata.Connected {
		return createSelectedKeylightAction(keylight)
	} else {
		return toast.CreateErrorToastAction("Keylight could not be found or not connected")
	}
}

func (m *Model) addNewKeylight() tea.Cmd {
	return func() tea.Msg {
		return AddKeylight{}
	}
}

func (m *Model) editKeylight(keylightId string) tea.Cmd {
	return func() tea.Msg {
		index, _ := strconv.Atoi(keylightId)
		keylight := &m.keylights[index]
		return EditKeylight{
			Keylight: keylight,
		}
	}
}

func (m *Model) removeKeylight(keylightId string) tea.Cmd {
	return func() tea.Msg {
		keylight := utils.FindKeylightWithId(m.keylights, keylightId)
		if keylight != nil {
			return RemoveKeylight{
				Keylight: keylight,
			}
		} else {
			return toast.CreateInfoToastAction("Keylight does not exist")
		}
	}
}

func (m *Model) reloadKeylights() tea.Cmd {
	return func() tea.Msg {
		return ReloadKeylights{}
	}
}

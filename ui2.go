package main

import (
	"fmt"
	"gollector/managers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/color"
)

type model struct {
	installs   []item
	list, item int
}

type item struct {
	index   int
	checked bool
}

var selectedInstalls []*managers.Installation

func (m *model) Init() tea.Cmd {

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeyMsg(typed)
	}
	return m, nil
}

func (m *model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc", "ctrl+c":
		return tea.Quit
	case " ":
		switch m.list {
		case 0:
			m.installs[m.item].checked = !m.installs[m.item].checked
		}
	case "enter":
		for i := 0; i < len(m.installs); i++ {
			item := m.installs[i]
			if item.checked {
				selectedInstalls = append(selectedInstalls, &installations[i])
			}
		}
		return tea.Quit
	case "up":
		if m.item > 0 {
			m.item--
		} else if m.list > 0 {
			m.list--
			m.item = len(m.installs) - 1
		}
	case "down":
		switch m.list {
		case 0:
			if m.item+1 < len(m.installs) {
				m.item++
			} else {
				m.list++
				m.item = 0
			}
		}
	}
	return nil
}

func (m *model) View() string {
	curInstall := -1
	switch m.list {
	case 0:
		curInstall = m.item
	}
	return m.renderList("Que voulez-vous installer ?", m.installs, curInstall)
}

func (m *model) renderList(header string, items []item, selected int) string {
	out := color.HEX("be5bff").Sprintf("\n-> " + header + "\n\n")
	for i, item := range items {
		sel := "  "
		if i == selected {
			sel = "->"
		}
		sel = color.HEX("62f6ff").Sprintf(sel)
		check := " "
		installName := color.HEX("af0000").Sprintf(installations[item.index].Name)
		if items[i].checked {
			check = color.New(color.FgGreen, color.BgBlack).Sprintf("✓")
			installName = color.HEX("00af18").Sprintf(installations[item.index].Name)
		}

		installDesc := color.HEX("3f3f3f").Sprintf(installations[item.index].Description)
		out += fmt.Sprintf("\t%s [%s] %s\n\t\t%s\n\n", sel, check, installName, installDesc)
	}

	out += color.HEX("3f3f3f").Sprintf("\n\n  Space: ") +
		color.HEX("8700c1").Sprintf("Select install") + "\n  " +
		color.HEX("3f3f3f").Sprintf("Enter: ") +
		color.HEX("8700c1").Sprintf("Start install on all provided machines\n\n")
	return out
}

func GetItems() []item {
	var items []item
	for a := 0; a < 5; a++ {
		for i := 0; i < len(installations); i++ {
			item := item{index: i, checked: false}
			items = append(items, item)
		}
	}

	return items
}

func ui2() {
	m := &model{
		installs: GetItems(),
	}
	if err := tea.NewProgram(m).Start(); err != nil {
		panic(fmt.Sprintf("failed to run program: %v", err))
	}
	if len(selectedInstalls) > 0 {
		color.HEX("62f6ff").Printf("\n\n[----- Installation en cours -----]\n\n")
	}
	InstallAll(selectedInstalls)
}

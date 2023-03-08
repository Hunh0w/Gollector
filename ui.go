package main

import (
	"fmt"
	"gollector/managers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/color"
)

type model struct {
	installs      []item
	item          int
	currentOffset int
}

type item struct {
	index   int
	checked bool
}

var selectedInstalls []*managers.Installation

func GetItemsWithOffset(items []item, offset int, amount int) []item {
	var itemsWithOffset []item
	items_len := len(items)
	for i := 0; i < amount; i++ {
		if items_len-1 < i+offset {
			break
		}
		itemsWithOffset = append(itemsWithOffset, items[i+offset])
	}
	return itemsWithOffset
}

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
		m.installs[m.item+m.currentOffset].checked = !m.installs[m.item+m.currentOffset].checked
	case "enter":
		for i := 0; i < len(m.installs); i++ {
			item := m.installs[i]
			if item.checked {
				selectedInstalls = append(selectedInstalls, &installations[item.index])
			}
		}
		return tea.Quit
	case "up":
		if m.item > 0 {
			m.item--
		} else {
			if m.currentOffset > 0 {
				m.currentOffset--
			}
		}

	case "down":
		currentItems := GetItemsWithOffset(m.installs, m.currentOffset, 4)
		if m.item+1 < 4 && m.item+1 < len(currentItems) {
			m.item++
		} else {

			itemsTest := GetItemsWithOffset(m.installs, m.currentOffset+1, 4)
			if len(itemsTest) >= 4 {
				m.currentOffset++
			}
		}
	}
	return nil
}

func (m *model) View() string {
	curInstall := m.item
	title := "Que voulez-vous installer ? " + color.HEX("40fcff").Sprintf("(%d modules)", len(m.installs))
	return m.renderList(title, m.installs, curInstall)
}

func (m *model) renderList(header string, items []item, selected int) string {
	out := color.HEX("be5bff").Sprintf("\n-> " + header + "\n\n")
	for i, item := range GetItemsWithOffset(items, m.currentOffset, 4) {
		sel := "  "
		if i == selected {
			sel = "->"
		}
		sel = color.HEX("62f6ff").Sprintf(sel)
		check := " "
		installName := color.HEX("af0000").Sprintf(installations[item.index].Name)
		if items[i+m.currentOffset].checked {
			check = color.New(color.FgGreen, color.BgBlack).Sprintf("✓")
			installName = color.HEX("00af18").Sprintf(installations[item.index].Name)
		}

		installDesc := color.HEX("3f3f3f").Sprintf(installations[item.index].Description)
		out += fmt.Sprintf("\t%s%d) %s [%s] %s\n\t\t%s\n\n", "\033[0m", i+m.currentOffset+1, sel, check, installName, installDesc)
	}

	out += color.HEX("3f3f3f").Sprintf("\n\n  Space: ") +
		color.HEX("8700c1").Sprintf("Select install") + "\n  " +
		color.HEX("3f3f3f").Sprintf("Enter: ") +
		color.HEX("8700c1").Sprintf("Start install on all provided machines\n\n")
	return out
}

func GetItems() []item {
	var items []item
	for i := 0; i < len(installations); i++ {
		item := item{index: i, checked: false}
		items = append(items, item)
	}

	return items
}

func ui() {
	m := &model{
		installs:      GetItems(),
		currentOffset: 0,
	}
	if err := tea.NewProgram(m).Start(); err != nil {
		panic(fmt.Sprintf("failed to run program: %v", err))
	}
	if len(selectedInstalls) > 0 {
		color.HEX("62f6ff").Printf("\n\n[----- Installation en cours -----]\n\n")
	}
	InstallAll(selectedInstalls)
}

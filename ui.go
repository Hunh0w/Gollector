package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gollector/managers"
)

func ui() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Que voulez-vous installer ?").
			SetExpansion(1).
			SetTextColor(tcell.ColorPurple).
			SetAlign(tview.AlignCenter))
	
	for i := 0; i < len(installations); i++ {
		table.SetCell(i+1, 0, tview.NewTableCell(installations[i].Name).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter))
	}

	table.SetCell(len(installations)+1, 0, tview.NewTableCell("Valider").
			SetTextColor(tcell.ColorGreen).
			SetAlign(tview.AlignCenter))

	table.SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		if row >= 1 && row <= len(installations) {
			cell := table.GetCell(row, column)
			color := cell.Color
			if color == tcell.ColorBlue {
				cell.SetTextColor(tcell.ColorWhite)
				return
			}
			cell.SetTextColor(tcell.ColorBlue)
			return
		}
		if row == len(installations)+1 {
			app.Stop()
			var installs []*managers.Installation
			for i := 0; i < len(installations); i++ {
				color := table.GetCell(i+1, 0).Color
				if color == tcell.ColorBlue {
					installs = append(installs, &installations[i])
				}
			}
			InstallAll(installs)
		}
	}).SetSelectable(true, true)

	flex := tview.NewFlex().
		AddItem(table, 0, 1, false).
		SetDirection(tview.FlexRow)

	if err := app.SetRoot(flex, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"fmt"
	"os/exec"
	"os"
	"encoding/json"
	"bytes"
	"strings"
	"io/ioutil"
)

type InstallationWay struct {
	Name	string `json:name`
	TestCmd	string `json:testCmd`
	Commands	[]string `json:commands`
}

type Installation struct {
	Name	string `json:name`
	Ways	[]InstallationWay `json:way`
}

var installations []Installation

func initModules() {
	files, err := ioutil.ReadDir("modules/")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		jsonFile, err := os.Open("modules/"+file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result Installation
		json.Unmarshal(byteValue, &result)

		installations = append(installations, result)
	}
}

func executeCommand(command string, output bool) bool {
	fmt.Println("Executing '"+command+"'...")
	array := strings.Fields(command)
	cmd := exec.Command(array[0], array[1:]...)

	var out bytes.Buffer
	cmd.Stdout = &out

	
    if err := cmd.Run(); err != nil {
		if !output {
			return false
		}
        fmt.Println(err)
		return false
    }


	if !output{
		return true
	}
	fmt.Println(out.String())
	return true
}

func install(installation Installation) {
	fmt.Println("Installing", installation.Name, "...")
	way, isSuccess := getWay(installation)
	if !isSuccess {
		fmt.Println("Installation impossible car aucun package manager détecté")
		return
	}
	installWay(way)
}

func installWay(way InstallationWay) {
	for _, command := range way.Commands {
		executeCommand(command, true)
	}
}

func getWay(installation Installation) (InstallationWay, bool) {
	var return_way InstallationWay
	success := false
	for _, way := range installation.Ways {
		isSuccess := executeCommand(way.TestCmd, false)
		if isSuccess {
			return_way = way
			success = true
			break
		}
	}
	return return_way, success
}

func main() {

	initModules()

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
			for i := 0; i < len(installations); i++ {
				color := table.GetCell(i+1, 0).Color
				if color == tcell.ColorBlue {
					install(installations[i])
				}
			}
		}
	}).SetSelectable(true, true)

	flex := tview.NewFlex().
		AddItem(table, 0, 1, false).
		SetDirection(tview.FlexRow)

	if err := app.SetRoot(flex, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
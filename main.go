package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"fmt"
	"os/exec"
	"os"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

type installation_way struct {
	name	string `json:name`
	commands	[]string `json:commands`
}

type installation struct {
	name	string `json:name`
	installation	[]installation_way `json:installation`
}

var installations []installation

func initModules() {
	files, err := ioutil.ReadDir("modules/")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		fmt.Println(file.Name(), "...")
		jsonFile, err := os.Open("modules/"+file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)

		array := result["installation"].([]interface{})
		for _, way := range array {
			fmt.Println(way)
		}
		
	}
}

func executeCommand(command string, output bool) {

	cmd := exec.Command("/bin/bash", "-c", command)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if !output {
		return
	}
	if err != nil {
		fmt.Println(command, ":", err)
		return
	}

	fmt.Println(out.String())
}

func install(installation installation) {
	fmt.Println("Installing", installation.name, "...")
	for i := 0; i < len(installation.installation[0].commands); i++ {
		command := installation.installation[0].commands[i]
		go executeCommand(command, true)
	}
}

func main() {

	initModules()

	if(true){
		return
	}

	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Que voulez-vous installer ?").
			SetExpansion(1).
			SetTextColor(tcell.ColorPurple).
			SetAlign(tview.AlignCenter))
	
	for i := 0; i < len(installations); i++ {
		table.SetCell(i+1, 0, tview.NewTableCell(installations[i].name).
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
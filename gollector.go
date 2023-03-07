package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"fmt"
	"os/exec"
	"bytes"
	"strings"
	"sync"
	"gollector/managers"
)


//type SSH_Host managers.SSH_Host

// Global Vars
var installations []managers.Installation
var hosts []managers.SSH_Host

/*
// Entry Point
*/
func main() {

	// Init required content
	installations = managers.InitModules();
	hosts = managers.InitHosts();	

	/*
	// UI Initialization
	*/
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

/*
// Required Initializations
*/





/*
// Goroutines Runner
*/
func InstallAll(installs []*managers.Installation){
	var wg sync.WaitGroup
	for i1 := 0; i1 < len(hosts); i1++ {
		host := hosts[i1]
		wg.Add(1)
		go InstallToHost(installs, &host, &wg)
	}
	wg.Wait()
}





/*
// Local Command Execution Utils
*/
func ExecuteCommand(command string, output bool) bool {
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

func Install(installation managers.Installation) {
	fmt.Println("Installing", installation.Name, "...")
	way, isSuccess := GetWay(installation)
	if !isSuccess {
		fmt.Println("Installation impossible car aucun package manager détecté")
		return
	}
	installWay(way)
}

func installWay(way managers.InstallationWay) {
	for _, command := range way.Commands {
		ExecuteCommand(command, true)
	}
}

func GetWay(installation managers.Installation) (managers.InstallationWay, bool) {
	var return_way managers.InstallationWay
	success := false
	for _, way := range installation.Ways {
		isSuccess := ExecuteCommand(way.TestCmd, false)
		if isSuccess {
			return_way = way
			success = true
			break
		}
	}
	return return_way, success
}
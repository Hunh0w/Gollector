package main

import (
	"bytes"
	"fmt"
	"gollector/managers"
	"os/exec"
	"strings"
	"sync"
)

// Global Vars
var installations []managers.Installation
var hosts []managers.SSH_Host

/*
// Entry Point
*/
func main() {

	// Init required content
	installations = managers.InitModules()
	hosts = managers.InitHosts()

	ui()
}

/*
// Goroutines Runner
*/
func InstallAll(installs []*managers.Installation) {
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
	fmt.Println("Executing '" + command + "'...")
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

	if !output {
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

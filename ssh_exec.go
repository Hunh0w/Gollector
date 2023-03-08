package main

import (
	"fmt"
	"gollector/managers"
	"sync"

	"github.com/sfreiberg/simplessh"
)

func InstallToHost(installs []*managers.Installation, ssh_host *managers.SSH_Host, wg *sync.WaitGroup) {

	// Utils Variables
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorCyan := "\033[36m"
	colorPurple := "\033[35m"
	colorGray := "\033[30m"

	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(ssh_host.HostAddress, ssh_host.User, ssh_host.Password); err != nil {
		fmt.Println(err)
		return
	}

	defer client.Close()
	defer wg.Done()

	for i1 := 0; i1 < len(installs); i1++ {
		install := installs[i1]
		pkg_detected := false
		for i2 := 0; i2 < len(install.Ways); i2++ {
			iway := install.Ways[i2]
			_, err := client.Exec(iway.TestCmd)
			if err != nil {
				continue
			}
			pkg_detected = true
			fmt.Printf("%s[%s]%s -> Installing %s%s%s...%s\n", colorPurple, ssh_host.HostAddress, colorRed, colorCyan, install.Name, colorRed, colorReset)

			result_messages := fmt.Sprintf("%s[%s]%s -> %s%s%s installation result: \n", colorPurple, ssh_host.HostAddress, colorGreen, colorCyan, install.Name, colorGreen)
			for i2 := 0; i2 < len(iway.Commands); i2++ {
				Command := iway.Commands[i2]
				result, err := client.Exec(Command)
				if err != nil {
					fmt.Printf("%s[%s]%s ERROR: \n%s\n", colorPurple, ssh_host.HostAddress, colorRed, err)
				}
				result_messages += fmt.Sprintf("%s%s%s\n", colorGray, string(result), colorReset)
			}
			fmt.Print(result_messages)

			break
		}
		if !pkg_detected {
			fmt.Printf("%s[%s] %s%s%s n'a pas pu être installé.\n", colorPurple, ssh_host.HostAddress, colorPurple, install.Name, colorRed)
		}
	}

}

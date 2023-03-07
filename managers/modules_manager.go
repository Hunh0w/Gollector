package managers

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
)

type InstallationWay struct {
	Name	string `json:name`
	TestCmd	string `json:testCmd`
	Commands	[]string `json:commands`
}

type Installation struct {
	Name	string `json:name`
	Description	string `json:description`
	Ways	[]InstallationWay `json:ways`
}


func InitModules() []Installation {
	files, err := ioutil.ReadDir("modules/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}

	var installations []Installation

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
		
		fmt.Println(result)

		installations = append(installations, result)
	}
	if len(installations) <= 0 {
		fmt.Println("No installations provided")
		os.Exit(4)
	}
	return installations
}
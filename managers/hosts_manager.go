package managers

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type SSH_Host struct {
	HostAddress		string `json:HostAddress`
	User			string `json:User`
	Password		string `json:Password`
}


func InitHosts() []SSH_Host {
	jsonFile, err := os.Open("hosts.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []SSH_Host
	json.Unmarshal(byteValue, &result)
	if len(result) <= 0 {
		fmt.Println("No hosts provided")
		os.Exit(3)
	}

	return result
}
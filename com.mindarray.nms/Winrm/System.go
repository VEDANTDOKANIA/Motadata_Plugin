package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func SystemData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	result := make(map[string]interface{})
	var errors []string
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		errors = append(errors, err.Error())
	}
	clients, er := client.CreateShell()
	defer clients.Close()
	if er != nil {
		errors = append(errors, er.Error())
		result["status"] = "fail"
		result["error"] = errors
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		a := "aa"
		output := ""
		ac := "(Get-WmiObject win32_operatingsystem).name;(Get-WMIObject win32_operatingsystem).version;whoami;(Get-WMIObject win32_operatingsystem).LastBootUpTime;" // Command jo humko run karna hain
		output, _, _, err = client.RunPSWithString(ac, a)
		res1 := strings.Split(output, "\n")
		result["system.os.name"] = strings.Split(res1[0], "\r")[0]
		result["system.os.version"] = strings.Split(res1[1], "\r")[0]
		result["system.user.name"] = strings.Split(res1[2], "\r")[0]
		result["system.up.time"] = strings.Split(res1[3], "\r")[0]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

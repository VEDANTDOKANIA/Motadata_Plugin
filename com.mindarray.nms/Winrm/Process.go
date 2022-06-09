package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
)

func ProcessData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	var errors []string
	result := make(map[string]interface{})
	if err != nil {
		errors = append(errors, err.Error())
	}
	clients, er := client.CreateShell()
	defer func(clients *winrm.Shell) {
		err := clients.Close()
		if err != nil {
			errors = append(errors, err.Error())
		}
	}(clients)
	if er != nil {
		errors = append(errors, er.Error())
		result["status"] = "fail"
		result["error"] = errors
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		a := "aa"
		output := ""
		ac := "(Get-Counter '\\Process(*)\\ID Process','\\Process(*)\\% Processor Time','\\Process(*)\\Thread Count' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
		output, _, _, err = client.RunPSWithString(ac, a)
		re := regexp.MustCompile("Path\\s*\\:\\s*\\\\+[\\w\\-#.]+\\\\\\w*\\(([\\w\\-#.]+)\\)\\\\%?\\s*(\\w*\\s*\\w*)\\s*\\w*\\s*:\\s*([\\d\\.]+)")
		value := re.FindAllStringSubmatch(output, -1)
		var processes []map[string]interface{}
		processes = append(processes, result)
		var count int
		for index := 0; index < len(value); index++ {
			temp := make(map[string]interface{})
			temp1 := make(map[string]interface{})
			processName := value[index][1]
			for j := 0; j < len(processes); j++ {
				temp = processes[j]
				if temp[processName] != nil {
					count = 1
					break
				} else {
					count = 0
				}
			}
			if count == 0 {
				temp1["process.name"] = processName
				if (value[index][2]) == "id process" {
					temp1["process.id"] = value[index][3]
				} else if value[index][2] == "processor time" {
					temp1["process.processor.time.percent"] = value[index][3]
				} else if value[index][2] == "thread count" {
					temp1["process.thread.count"] = value[index][3]
				}
				processes = append(processes, temp1)

			} else {
				if (value[index][2]) == "id process" {
					temp["process.id"] = value[index][3]
				} else if value[index][2] == "processor time" {
					temp["process.processor.time.percent"] = value[index][3]
				} else if value[index][2] == "thread count" {
					temp["process.thread.count"] = value[index][3]
				}
				count = 1
				processes = append(processes, temp)
			}
		}
		processes = processes[1:len(processes)]
		size := (len(processes)) / 3
		var values []map[string]interface{}
		for index := 0; index < len(processes)/3; index = index + 1 {
			count := index
			temp2 := make(map[string]interface{})
			temp2 = processes[index]
			temp2["process.processor.time.percent"] = processes[count+size]["process.processor.time.percent"]
			temp2["process.thread.count"] = processes[count+size+size]["process.thread.count"]
			values = append(values, temp2)
		}
		result["process"] = values
		result["status"] = "success"
		data, err2 := json.Marshal(result)
		if err2 != nil {
			out := make(map[string]interface{})
			out["status"] = "fail"
			out["error"] = err2.Error()
			output, _ := json.Marshal(out)
			fmt.Print(string(output))
		} else {
			fmt.Print(string(data))
		}
	}
}

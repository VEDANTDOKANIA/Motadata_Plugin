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
	result := make(map[string]interface{})
	if err != nil {
		result["error"] = "yes"
		result["Cause"] = err
	} else {
		result["error"] = "no"
	}
	a := "aa"
	output := ""
	ac := "(Get-Counter '\\Process(*)\\ID Process','\\Process(*)\\% Processor Time','\\Process(*)\\Thread Count' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
	output, _, _, err = client.RunPSWithString(ac, a)
	re := regexp.MustCompile("Path\\s*:\\s\\\\\\\\(\\w*-\\w*)\\\\\\w*\\((\\S*)\\)\\\\([\\w\\d\\s%]+)\\n\\w*\\s\\:\\s(\\d*)")
	value := re.FindAllStringSubmatch(output, -1)

	var processes []map[string]interface{}
	processes = append(processes, result)
	var count int
	for i := 0; i < len(value); i++ {
		temp := make(map[string]interface{})
		temp1 := make(map[string]interface{})
		processName := value[i][2]
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
			if (value[i][3]) == "id process\r" {
				temp1["process.id"] = value[i][4]
			} else if value[i][3] == "% processor time\r" {
				temp1["process.processor.time.percent"] = value[i][4]
			} else if value[i][3] == "thread count\r" {
				temp1["process.thread.count"] = value[i][4]
			}
			processes = append(processes, temp1)

		} else {
			if (value[i][3]) == "id process\r" {
				temp["process.id"] = value[i][4]
			} else if value[i][3] == "% processor time\r" {
				temp["process.processor.time.percent"] = value[i][4]
			} else if value[i][3] == "thread count\r" {
				temp["process.thread.count"] = value[i][4]
			}
			count = 1
			processes = append(processes, temp)
		}

	}
	processes = processes[1:len(processes)]
	size := (len(processes)) / 3
	var Values []map[string]interface{}
	for k := 0; k < len(processes)/3; k = k + 1 {
		count := k
		temp2 := make(map[string]interface{})
		temp2 = processes[k]
		temp2["process.processor.time.percent"] = processes[count+size]["process.processor.time.percent"]
		temp2["process.thread.count"] = processes[count+size+size]["process.thread.count"]
		Values = append(Values, temp2)
	}
	result["process"] = Values
	result["ip"] = credentials["ip"]
	result["metric.group"] = credentials["metric.group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

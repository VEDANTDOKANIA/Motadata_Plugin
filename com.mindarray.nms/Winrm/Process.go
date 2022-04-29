package Winrm

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
)

func ProcessData(credentials map[string]interface{}) {
	host := (credentials["IP_Address"]).(string)
	port := int(credentials["Port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	result := make(map[string]interface{})
	if err != nil {
		result["Error"] = "yes"
		result["Cause"] = err
	} else {
		result["Error"] = "no"
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
			temp1["Process.Name"] = processName
			if (value[i][3]) == "id process\r" {
				temp1["Process.ID"] = value[i][4]
			} else if value[i][3] == "% processor time\r" {
				temp1["Process.Processor.Time.Percent"] = value[i][4]
			} else if value[i][3] == "thread count\r" {
				temp1["Process.thread.Count"] = value[i][4]
			}
			processes = append(processes, temp1)

		} else {
			if (value[i][3]) == "id process\r" {
				temp["Process.ID"] = value[i][4]
			} else if value[i][3] == "% processor time\r" {
				temp["Process.Processor.Time.Percent"] = value[i][4]
			} else if value[i][3] == "thread count\r" {
				temp["Process.thread.Count"] = value[i][4]
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
		temp2["Process.Processor.Time.Percent"] = processes[count+size]["Process.Processor.Time.Percent"]
		temp2["Process.Thread.Count"] = processes[count+size+size]["Process.thread.Count"]
		Values = append(Values, temp2)
	}
	result["Process"] = Values
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
	"strings"
)

func CpuData(credentials map[string]interface{}) {
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
		ac := "(Get-Counter '\\Processor(*)\\% Idle Time','\\Processor(*)\\% Processor Time','\\Processor(*)\\% user time' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
		output, _, _, err = client.RunPSWithString(ac, a)
		re := regexp.MustCompile("Path\\s*:\\s\\\\\\\\(\\w*-\\w*)\\\\\\w*\\((\\S*)\\)\\\\([\\w\\d\\s%]+)\\n\\w*\\s\\:\\s(\\d*)")
		value := re.FindAllStringSubmatch(output, -1)
		var counters = 3
		var cores []map[string]interface{}
		size := len(value) / counters
		for i := 0; i < len(value)/counters; i++ {
			count := 0
			core := make(map[string]interface{})
			if value[i][1] == "_total" {
				result["system.cpu.idle.percent"] = value[i][3]
				result["system.cpu.process.percent"] = value[count+size][3]
				result["system.cpu.user.percent"] = strings.Split(value[count+size+size][3], "\r")[0]
			} else {
				core["core.name"] = value[i][2]
				core["core.idle.percent"] = value[i][4]
				core["core.process.percent"] = value[count+size][4]
				core["core.user.percent"] = strings.Split(value[count+size+size][4], "\r")[0]
				cores = append(cores, core)
			}
		}
		result["cores"] = cores
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

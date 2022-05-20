package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
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
	ac := "(Get-Counter '\\Processor(*)\\% Idle Time','\\Processor(*)\\% Processor Time','\\Processor(*)\\% user time' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
	output, _, _, err = client.RunPSWithString(ac, a)
	re := regexp.MustCompile("Path\\s*:\\s\\\\\\\\\\w*-\\w*\\\\\\w*\\((\\S*)\\)\\\\([\\w\\d\\s%]+)\\n\\w*\\s\\:\\s(\\d*.\\d*)")
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
			core["core.name"] = value[i][1]
			core["core.idle.percent"] = value[i][3]
			core["core.process.percent"] = value[count+size][3]
			core["core.user.percent"] = strings.Split(value[count+size+size][3], "\r")[0]
			cores = append(cores, core)
		}
	}
	result["cores"] = cores
	result["ip"] = credentials["ip"]
	result["metric.group"] = credentials["metric.group"]
	fmt.Println(value)
}

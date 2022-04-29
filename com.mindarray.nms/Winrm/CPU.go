package Winrm

import (
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
	"strings"
)

func CpuData(credentials map[string]interface{}) {
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
			result["System.Cpu.Idle.Percent"] = value[i][3]
			result["System.Cpu.Process.Percent"] = value[count+size][3]
			result["System.Cpu.User.Percent"] = strings.Split(value[count+size+size][3], "\r")[0]
		} else {
			core["Core.Name"] = value[i][1]
			core["Core.Idle.Percent"] = value[i][3]
			core["Core.Process.Percent"] = value[count+size][3]
			core["Core.User.Percent"] = strings.Split(value[count+size+size][3], "\r")[0]
			cores = append(cores, core)
		}
	}
	result["Cores"] = cores
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	fmt.Println(value)
}

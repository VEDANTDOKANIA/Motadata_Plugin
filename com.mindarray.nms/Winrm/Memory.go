package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strconv"
	"strings"
)

func MemoryData(credentials map[string]interface{}) {
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
		ac := "Get-WmiObject win32_OperatingSystem |%{\"{0} \n{1} \n{2} \n" +
			"{3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}" // Command jo humko run karna hain
		output, _, _, err = client.RunPSWithString(ac, a)
		res1 := strings.Split(output, "\n")

		total_space_memory, _ := strconv.ParseInt(strings.TrimSpace(res1[0]), 10, 64)
		total_space_virtual, _ := strconv.ParseInt(strings.TrimSpace(res1[2]), 10, 64)
		free_space_memory, _ := strconv.ParseInt(strings.TrimSpace(res1[1]), 10, 64)
		free_space_virtual, _ := strconv.ParseInt(strings.TrimSpace(res1[3]), 10, 64)
		total_space := float64(total_space_memory + total_space_virtual)
		free_space := float64(free_space_virtual + free_space_memory)
		percent := float64(free_space/total_space) * 100
		result["memory.total.bytes"] = total_space_memory * 1000
		result["memory.free.bytes"] = free_space_memory * 1000
		result["memory.used.bytes"] = (total_space_memory - free_space_memory) * 1000
		result["memory.virtual.total.bytes"] = total_space_virtual * 1000
		result["memory.virtual.free.bytes"] = free_space_virtual * 1000
		result["memory.virtual.used.bytes"] = (total_space_virtual - free_space_virtual) * 1000
		result["memory.used.percent"] = percent
		result["memory.available.percent"] = 100.0 - percent
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

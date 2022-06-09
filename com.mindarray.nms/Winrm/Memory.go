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
		ac := "Get-WmiObject win32_OperatingSystem |%{\"{0} \n{1} \n{2} \n" +
			"{3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}"
		output, _, _, err = client.RunPSWithString(ac, a)
		res1 := strings.Split(output, "\n")

		totalSpaceMemory, _ := strconv.ParseInt(strings.TrimSpace(res1[0]), 10, 64)
		totalSpaceVirtual, _ := strconv.ParseInt(strings.TrimSpace(res1[2]), 10, 64)
		freeSpaceMemory, _ := strconv.ParseInt(strings.TrimSpace(res1[1]), 10, 64)
		freeSpaceVirtual, _ := strconv.ParseInt(strings.TrimSpace(res1[3]), 10, 64)
		totalSpace := float64(totalSpaceMemory + totalSpaceVirtual)
		freeSpace := float64(freeSpaceVirtual + freeSpaceMemory)
		percent := float64(freeSpace/totalSpace) * 100
		result["memory.total.bytes"] = totalSpaceMemory * 1000
		result["memory.free.bytes"] = freeSpaceMemory * 1000
		result["memory.used.bytes"] = (totalSpaceMemory - freeSpaceMemory) * 1000
		result["memory.virtual.total.bytes"] = totalSpaceVirtual * 1000
		result["memory.virtual.free.bytes"] = freeSpaceVirtual * 1000
		result["memory.virtual.used.bytes"] = (totalSpaceVirtual - freeSpaceVirtual) * 1000
		result["memory.used.percent"] = percent
		result["memory.available.percent"] = 100.0 - percent
		result["metric.group"] = credentials["metric.group"]
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

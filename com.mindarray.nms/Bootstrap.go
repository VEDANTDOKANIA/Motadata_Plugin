package main

import (
	"MotadataPlugin/com.mindarray.nms/SNMP"
	"MotadataPlugin/com.mindarray.nms/SSH"
	"MotadataPlugin/com.mindarray.nms/Winrm"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	argument, _ := base64.StdEncoding.DecodeString(os.Args[1])
	credentials := make(map[string]interface{})
	var err = json.Unmarshal([]byte(string(argument)), &credentials)
	var errors []string
	result := make(map[string]interface{})
	if err != nil {
		errors = append(errors, err.Error())

	}

	if credentials["category"] == "discovery" {
		if credentials["type"] == "linux" {
			SSH.Discovery(credentials)
		} else if credentials["type"] == "windows" {
			Winrm.Discovery(credentials)
		} else if credentials["type"] == "snmp" {
			SNMP.Discovery(credentials)
		} else {
			errors = append(errors, "wrong tpe provided")
		}

	} else if credentials["category"] == "polling" {
		if credentials["type"] == "linux" {
			switch credentials["metric.group"] {
			case "system":
				SSH.SystemData(credentials)
				break
			case "disk":
				SSH.DiskData(credentials)
				break
			case "memory":
				SSH.MemoryData(credentials)
				break
			case "process":
				SSH.ProcessData(credentials)
				break
			case "cpu":
				SSH.CpuData(credentials)
			default:
				errors = append(errors, "Wrong metric group selected for metric type linux")

			}
		} else if credentials["type"] == "windows" {
			switch credentials["metric.group"] {
			case "system":
				Winrm.SystemData(credentials)
				break
			case "disk":
				Winrm.DiskData(credentials)
				break
			case "memory":
				Winrm.MemoryData(credentials)
				break
			case "process":
				Winrm.ProcessData(credentials)
				break
			case "cpu":
				Winrm.CpuData(credentials)
			default:
				errors = append(errors, "Wrong metric group selected for metric type Windows")

			}
		} else if credentials["type"] == "snmp" {
			switch credentials["metric.group"] {
			case "system":
				SNMP.SystemData(credentials)
				break
			case "interface":
				SNMP.InterfaceData(credentials)
				break
			default:
				errors = append(errors, "Wrong metric group selected for metric type Network Devices")
			}
		} else {
			errors = append(errors, "wrong type provided")
		}
	} else {
		errors = append(errors, "wrong category given")
	}
	if errors != nil {
		result["status"] = "fail"
		result["error"] = errors
	}

	data, _ := json.Marshal(result)
	if string(data) != "{}" {
		fmt.Print(string(data))
	}

}

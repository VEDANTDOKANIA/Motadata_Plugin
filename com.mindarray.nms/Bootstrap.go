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

//TODO CHange error msg in all metric group with status

func main() {
	argument, _ := base64.StdEncoding.DecodeString(os.Args[1])
	credentials := make(map[string]interface{})
	var error = json.Unmarshal([]byte(string(argument)), &credentials)
	result := make(map[string]interface{})
	if error != nil {
		result["error"] = "yes"
		result["Cause"] = error

	}

	if credentials["category"] == "discovery" {
		if credentials["type"] == "linux" {
			SSH.Discovery(credentials)
		} else if credentials["type"] == "windows" {
			Winrm.Discovery(credentials)
		} else if credentials["type"] == "snmp" {
			SNMP.Discovery(credentials)
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
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type linux"

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
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Windows"

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
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Network Devices"
			}
		}
	} else {
		result["error"] = "yes"
		result["Cause"] = "Wrong Category Given"

	}
	data, _ := json.Marshal(result)
	if string(data) != "{}" {
		fmt.Print(string(data))
	}

}

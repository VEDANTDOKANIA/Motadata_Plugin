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
		if credentials["metric.type"] == "linux" {
			SSH.Discovery(credentials)
		} else if credentials["metric.type"] == "windows" {
			Winrm.Discovery(credentials)
		} else if credentials["metric.type"] == "network" {
			SNMP.Discovery(credentials)
		}

	} else if credentials["category"] == "polling" {
		if credentials["metric.type"] == "linux" {
			switch credentials["metric.group"] {
			case "System":
				SSH.SystemData(credentials)
				break
			case "Disk":
				SSH.DiskData(credentials)
				break
			case "Memory":
				SSH.MemoryData(credentials)
				break
			case "Process":
				SSH.ProcessData(credentials)
				break
			case "CPU":
				SSH.CpuData(credentials)
			default:
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type linux"

			}
		} else if credentials["metric.type"] == "windows" {
			switch credentials["metric.group"] {
			case "System":
				Winrm.SystemData(credentials)
				break
			case "Disk":
				Winrm.DiskData(credentials)
				break
			case "Memory":
				Winrm.MemoryData(credentials)
				break
			case "Process":
				Winrm.ProcessData(credentials)
				break
			case "CPU":
				Winrm.CpuData(credentials)
			default:
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Windows"

			}
		} else if credentials["metric.type"] == "network" {
			switch credentials["metric.group"] {
			case "System":
				SNMP.SystemData(credentials)
				break
			case "Interface":
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

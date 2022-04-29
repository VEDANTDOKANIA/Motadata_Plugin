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
	var error = json.Unmarshal([]byte(string(argument)), &credentials)
	result := make(map[string]interface{})
	if error != nil {
		result["Error"] = "yes"
		result["Cause"] = error

	}
	fmt.Println(credentials)
	if credentials["category"] == "Discovery" {
		if credentials["Metric_Type"] == "linux" {
			SSH.Discovery(credentials)
		} else if credentials["Metric_Type"] == "windows" {
			Winrm.Discovery(credentials)
		} else if credentials["Metric_Type"] == "network" {
			SNMP.Discovery(credentials)
		}

	} else if credentials["category"] == "Polling" {
		if credentials["Metric_Type"] == "linux" {
			switch credentials["Metric_Group"] {
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
				result["Error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type linux"

			}
		} else if credentials["Metric_Type"] == "windows" {
			switch credentials["Metric_Group"] {
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
				result["Error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Windows"

			}
		} else if credentials["Metric_Type"] == "network" {
			switch credentials["Metric_Group"] {
			case "System":
				SNMP.SystemData(credentials)
				break
			case "Interface":
				SNMP.InterfaceData(credentials)
				break
			default:
				result["Error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Network Devices"
			}
		}
	} else {
		result["Error"] = "yes"
		result["Cause"] = "Wrong Category Given"

	}
	data, _ := json.Marshal(result)
	if string(data) != "{}" {
		fmt.Print(string(data))
	}

}

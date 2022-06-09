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
	var err = json.Unmarshal(argument, &credentials)
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
			errors = append(errors, "wrong type provided")
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
				errors = append(errors, "wrong metric group selected for metric type linux")

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
				errors = append(errors, "wrong metric group selected for metric type windows")

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
				errors = append(errors, "wrong metric group selected for metric type network devices")
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

	data, err2 := json.Marshal(result)
	if err2 != nil {
		out := make(map[string]interface{})
		out["status"] = "fail"
		out["error"] = err2.Error()
		output, _ := json.Marshal(out)
		fmt.Print(string(output))
	} else {
		if string(data) != "{}" {
			fmt.Print(string(data))
		}
	}
}

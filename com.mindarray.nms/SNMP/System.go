package SNMP

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

func SystemData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	result := make(map[string]interface{})
	var version = g.Version1
	switch credentials["version"] {
	case "version1":
		version = g.Version1
		break
	case "version2":
		version = g.Version2c
		break
	case "version3":
		version = g.Version3
		break
	}

	params := &g.GoSNMP{
		Target:    credentials["ip"].(string),
		Port:      uint16(int(credentials["port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(1) * time.Second,
	}
	err := params.Connect()
	var errors []string
	if err != nil {
		result["error"] = err
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		oid := []string{"1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.6.0", "1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0"}
		value, _ := params.Get(oid)
		for _, variable := range value.Variables {
			switch variable.Name {
			case ".1.3.6.1.2.1.1.5.0":
				result["system.name"] = string(variable.Value.([]byte))
				break
			case ".1.3.6.1.2.1.1.1.0":
				result["system.description"] = strings.Replace(string(variable.Value.([]byte)), "\r\n", " ", 3)
				break
			case ".1.3.6.1.2.1.1.6.0":
				if len(variable.Value.([]uint8)) == 0 {
					result["system.location"] = "empty"
				} else {
					result["system.location"] = string(variable.Value.([]byte))
				}
				break
			case ".1.3.6.1.2.1.1.2.0":
				result["system.oid"] = variable.Value
				break
			case ".1.3.6.1.2.1.1.3.0":
				result["system.up.time"] = variable.Value
				break
			default:
				errors = append(errors, "unknown interface")
			}

		}
		if len(errors) == 0 {
			result["status"] = "success"
		} else {
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
			fmt.Print(string(data))
		}

	}
}

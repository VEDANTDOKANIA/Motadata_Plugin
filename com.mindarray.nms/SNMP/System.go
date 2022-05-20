package SNMP

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
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
	if err != nil {
		result["error"] = "yes"
		result["Cause"] = err
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
		//return
	} else {
		result["error"] = "no"
	}

	oid := []string{"1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.6.0", "1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0"}
	value, _ := params.Get(oid)
	for _, variable := range value.Variables {

		switch variable.Name {
		case ".1.3.6.1.2.1.1.5.0":
			result["system_name"] = string(variable.Value.([]byte))
			break
		case ".1.3.6.1.2.1.1.1.0":
			result["system_description"] = string(variable.Value.([]byte))
			break
		case ".1.3.6.1.2.1.1.6.0":
			if len(variable.Value.([]uint8)) == 0 {
				result["system_loaction"] = "No location Specified"
			} else {
				result["system_loaction"] = string(variable.Value.([]byte))
			}

			break
		case ".1.3.6.1.2.1.1.2.0":
			result["system_oid"] = variable.Value
			break
		case ".1.3.6.1.2.1.1.3.0":
			result["system_upTime"] = variable.Value
			break
		default:
			result["error"] = "Unknown Interface"
		}

	}
	result["ip"] = credentials["ip"]
	result["metric.group"] = credentials["metric.group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

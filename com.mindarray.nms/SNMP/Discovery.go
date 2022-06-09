package SNMP

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

func Discovery(credentials map[string]interface{}) {
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
	var errors []string
	params := &g.GoSNMP{
		Target:    credentials["ip"].(string),
		Port:      uint16(int(credentials["port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(2) * time.Second,
	}
	err := params.Connect()
	if err != nil {
		result["status"] = "fail"
		result["error"] = err.Error()
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		walkOid := "1.3.6.1.2.1.2.2.1.1"
		snmpIndex := "1.3.6.1.2.1.2.2.1.1."
		description := "1.3.6.1.2.1.2.2.1.2."
		name := "1.3.6.1.2.1.31.1.1.1.1."
		operationalStatus := "1.3.6.1.2.1.2.2.1.8."
		alias := "1.3.6.1.2.1.31.1.1.1.18."

		var walkOidArray []string
		walk := params.Walk(walkOid, func(pdu g.SnmpPDU) error {
			switch pdu.Type {
			case g.IPAddress:
				result := pdu.Value
				walkOidArray = append(walkOidArray, result.(string))
				break
			case g.Integer:
				result := g.ToBigInt(pdu.Value)
				walkOidArray = append(walkOidArray, result.String())
				break
			case g.OctetString:
				result := pdu.Value.([]byte)
				walkOidArray = append(walkOidArray, string(result))
				break
			case g.Gauge32:
				result := pdu.Value
				walkOidArray = append(walkOidArray, result.(string))
			default:
				result := pdu.Value
				walkOidArray = append(walkOidArray, result.(string))
			}
			return nil
		},
		)
		if walk != nil {
			errors = append(errors, walk.Error())
		} else {
			var oids []string
			for i := 0; i < len(walkOidArray); i++ {
				oids = append(oids, snmpIndex+walkOidArray[i])
				oids = append(oids, description+walkOidArray[i])
				oids = append(oids, name+walkOidArray[i])
				oids = append(oids, operationalStatus+walkOidArray[i])
				oids = append(oids, alias+walkOidArray[i])
			}
			var startIndex = 0
			var endIndex = 50

			var resultArray []interface{}

			for {
				if len(resultArray) == len(oids) {
					break
				}
				output, error := params.Get(oids[startIndex:endIndex])
				if error != nil {
					errors = append(errors, walk.Error())
					return
				}
				for _, variable := range output.Variables {
					resultArray = append(resultArray, SnmpData(variable))
				}
				startIndex = endIndex
				endIndex = endIndex + 40
				if endIndex > len(oids) {
					endIndex = len(oids)
				}
			}
			var interfaces []map[string]interface{}
			for index := 0; index < len(resultArray); index = index + 5 {
				interfaceValue := make(map[string]interface{})
				interfaceValue["interface.index"] = resultArray[index].(int)
				interfaceValue["interface.description"] = resultArray[index+1]
				interfaceValue["interface.name"] = resultArray[index+2]
				if resultArray[index+3] == 1 {
					interfaceValue["interface.operational.status"] = "up"
				} else {
					interfaceValue["interface.operational.status"] = "down"
				}
				if resultArray[index+4] == "" {
					interfaceValue["interface.alias.name"] = ""
				} else {
					interfaceValue["interface.alias.name"] = strings.Trim(fmt.Sprintf("%v", resultArray[index+4]), "\"")
				}
				interfaces = append(interfaces, interfaceValue)
			}
			result["interfaces"] = interfaces
			value, _ := params.Get([]string{"1.3.6.1.2.1.1.5.0"})
			variable := value.Variables[0]
			result["host"] = string(variable.Value.([]byte))
		}
		if len(errors) == 0 {
			result["status"] = "success"
		} else {
			result["status"] = "fail"
			result["error"] = errors
		}
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

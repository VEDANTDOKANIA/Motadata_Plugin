package SNMP

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

func InterfaceData(credentials map[string]interface{}) {
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
		Timeout:   time.Duration(1) * time.Second,
	}
	err := params.Connect()
	if err != nil {
		result["error"] = err
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		walkOid := "1.3.6.1.2.1.2.2.1.1"
		snmpIndex := "1.3.6.1.2.1.2.2.1.1."
		description := "1.3.6.1.2.1.2.2.1.2."
		name := "1.3.6.1.2.1.31.1.1.1.1."
		operationalStatus := "1.3.6.1.2.1.2.2.1.8."
		adminStatus := "1.3.6.1.2.1.2.2.1.7."
		alias := "1.3.6.1.2.1.31.1.1.1.18."
		sentError := "1.3.6.1.2.1.2.2.1.20."
		receiveError := "1.3.6.1.2.1.2.2.1.14."
		sentOctets := "1.3.6.1.2.1.2.2.1.16."
		receiveOctets := "1.3.6.1.2.1.2.2.1.10."
		ifSpeed := "1.3.6.1.2.1.2.2.1.5."

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
		}

		var oids []string
		for index := 0; index < len(walkOidArray); index++ {
			oids = append(oids, snmpIndex+walkOidArray[index])
			oids = append(oids, description+walkOidArray[index])
			oids = append(oids, name+walkOidArray[index])
			oids = append(oids, operationalStatus+walkOidArray[index])
			oids = append(oids, adminStatus+walkOidArray[index])
			oids = append(oids, alias+walkOidArray[index])
			oids = append(oids, sentError+walkOidArray[index])
			oids = append(oids, receiveError+walkOidArray[index])
			oids = append(oids, sentOctets+walkOidArray[index])
			oids = append(oids, receiveOctets+walkOidArray[index])
			oids = append(oids, ifSpeed+walkOidArray[index])
		}
		var startIndex = 0
		var endIndex = 40
		var resultArray []interface{}
		for {
			if len(resultArray) == len(oids) {
				break
			}
			output, err2 := params.Get(oids[startIndex:endIndex])
			if err2 != nil {
				errors = append(errors, err2.Error())
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
		for index := 0; index < len(resultArray); index = index + 11 {
			interfaceValue := make(map[string]interface{})
			interfaceValue["interface.index"] = resultArray[index].(int)
			interfaceValue["interface.description"] = resultArray[index+1]
			interfaceValue["interface.name"] = resultArray[index+2]
			if resultArray[index+3] == 1 {
				interfaceValue["interface.operational.status"] = "up"
			} else {
				interfaceValue["interface.operational.status"] = "down"
			}
			if resultArray[index+4] == 1 {
				interfaceValue["interface.admin.status"] = "up"
			} else {
				interfaceValue["interface.admin.status"] = "down"
			}
			interfaceValue["interface.alias.name"] = strings.Trim(fmt.Sprintf("%v", resultArray[index+5]), "\"")
			interfaceValue["interface.sent.errors"] = resultArray[index+6]
			interfaceValue["interface.receive.errors"] = resultArray[index+7]
			interfaceValue["interface.sent.octets"] = resultArray[index+8]
			interfaceValue["interface.receive.octets"] = resultArray[index+9]
			interfaceValue["interface.speed"] = resultArray[index+10]
			interfaces = append(interfaces, interfaceValue)
		}
		result["interfaces"] = interfaces
		if len(errors) == 0 {
			result["status"] = "success"
		} else {
			result["status"] = "fail"
			result["error"] = errors
		}
		data, error := json.Marshal(result)
		if error != nil {
			out := make(map[string]interface{})
			out["status"] = "fail"
			out["error"] = error.Error()
			output, _ := json.Marshal(out)
			fmt.Print(string(output))
		} else {
			fmt.Print(string(data))
		}
	}
}
func SnmpData(pdu g.SnmpPDU) interface{} {

	if pdu.Value == " " {
		return pdu.Value
	}
	switch pdu.Type {
	case g.IPAddress:
		return pdu.Value
	case g.Integer:
		return pdu.Value
	case g.OctetString:
		return string(pdu.Value.([]byte))
	default:
		return pdu.Value
	}

}

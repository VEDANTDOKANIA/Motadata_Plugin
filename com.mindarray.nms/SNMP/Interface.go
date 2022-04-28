package SNMP

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"time"
)

func InterfaceData(credentials map[string]interface{}) {
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
		Target:    credentials["IP_Address"].(string),
		Port:      uint16(int(credentials["Port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(1) * time.Second,
	}
	err := params.Connect()
	if err != nil {
		result["Error"] = "yes"
		result["Cause"] = err
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
		//return
	} else {
		result["Error"] = "no"
	}
	walkOid := "1.3.6.1.2.1.2.2.1.1"
	index := "1.3.6.1.2.1.2.2.1.1."
	description := "1.3.6.1.2.1.2.2.1.2."
	//address := "1.3.6.1.2.1.2.2.1.6."
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
		result["Error"] = walk
	}

	var oids []string
	for i := 0; i < len(walkOidArray); i++ {
		oids = append(oids, index+walkOidArray[i])
		oids = append(oids, description+walkOidArray[i])
		//oids = append(oids, address+walkOidArray[i])
		oids = append(oids, name+walkOidArray[i])
		oids = append(oids, operationalStatus+walkOidArray[i])
		oids = append(oids, adminStatus+walkOidArray[i])
		oids = append(oids, alias+walkOidArray[i])
		oids = append(oids, sentError+walkOidArray[i])
		oids = append(oids, receiveError+walkOidArray[i])
		oids = append(oids, sentOctets+walkOidArray[i])
		oids = append(oids, receiveOctets+walkOidArray[i])
		oids = append(oids, ifSpeed+walkOidArray[i])
	}
	var startIndex = 0
	var endIndex = 60

	var resultArray []interface{}

	for {
		if len(resultArray) == len(oids) {
			break
		}
		output, error := params.Get(oids[startIndex:endIndex])
		if error != nil {
			result["Error"] = error
			return
		}
		for _, variable := range output.Variables {
			resultArray = append(resultArray, SnmpData(variable))
		}
		startIndex = endIndex
		endIndex = endIndex + 60

		if endIndex > len(oids) {
			endIndex = len(oids)
		}

	}
	var interfaces []map[string]interface{}
	for i := 0; i < len(resultArray); i = i + 11 {
		interfaceValue := make(map[string]interface{})
		interfaceValue["Interface.Index"] = resultArray[i].(int)
		interfaceValue["Interface.Description"] = resultArray[i+1]
		interfaceValue["Interface.Name"] = resultArray[i+2]
		if resultArray[i+3] == 1 {
			interfaceValue["Interface.Operational.Status"] = "Up"
		} else {
			interfaceValue["Interface.Operational.Status"] = "Down"
		}
		if resultArray[i+4] == 1 {
			interfaceValue["Interface.Admin.Status"] = "Up"
		} else {
			interfaceValue["Interface.Admin.Status"] = "Down"
		}
		if resultArray[i+5] == "" {
			interfaceValue["Interface.Alias.Name"] = "Empty"
		} else {
			interfaceValue["Interface.Alias.Name"] = resultArray[i+5]
		}

		interfaceValue["Interface.Sent.Errors"] = resultArray[i+6]
		interfaceValue["Interface.Receive.Errors"] = resultArray[i+7]
		interfaceValue["Interface.Sent.Octets"] = resultArray[i+8]
		interfaceValue["Interface.Receive.Octets"] = resultArray[i+8]
		interfaceValue["Interface.Speed"] = resultArray[i+9]

		interfaces = append(interfaces, interfaceValue)
	}
	result["Interface"] = interfaces
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
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
		//data, _ := hex.DecodeString(string(pdu.Value.([]byte)))
		//return data
	default:
		return pdu.Value
	}

}

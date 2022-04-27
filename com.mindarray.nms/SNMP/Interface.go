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

}

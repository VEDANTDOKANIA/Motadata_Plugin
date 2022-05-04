package SNMP

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
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

	params := &g.GoSNMP{
		Target:    credentials["ip.address"].(string),
		Port:      uint16(int(credentials["port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(2) * time.Second,
	}
	err := params.Connect()
	//oid := []string{"1.3.6.1.2.1.1.5.0"}
	if err != nil {
		result["status"] = "Unsuccessful"
		result["error"] = err
	} else {
		result["status"] = "successful"
	}
	_, error := params.Get([]string{"1.3.6.1.2.1.1.5.0"})
	if error != nil {
		result["status"] = "Unsuccessful"
		result["error"] = error.Error()
	} else {
		result["status"] = "successful"
	}

	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

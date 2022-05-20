package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	var errorOccurred []string
	defer exception.ErrorHandle(credentials)
	result := make(map[string]interface{})
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		errorOccurred = append(errorOccurred, err.Error())
	}
	_, err2 := client.CreateShell()
	if err2 != nil {
		errorOccurred = append(errorOccurred, err2.Error())
	}
	if len(errorOccurred) == 0 {
		result["status"] = "success"
	} else {
		result["status"] = "fail"
		result["error"] = errorOccurred
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

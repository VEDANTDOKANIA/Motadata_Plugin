package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	result := make(map[string]interface{})
	var errors []string
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		errors = append(errors, err.Error())
	}
	_, er := client.CreateShell()

	if er != nil {
		if strings.Contains(er.Error(), "http response error: 401 - invalid content type") {
			errors = append(errors, "wrong username or password")
		} else if strings.Contains(er.Error(), " connect: connection refused") {
			errors = append(errors, "wrong ip or port")
		}

		result["status"] = "fail"
		result["error"] = errors
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		a := "aa"
		output := ""
		ac := "hostname"
		output, _, _, err = client.RunPSWithString(ac, a)
		result["host"] = strings.Split(output, "\r\n")[0]
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

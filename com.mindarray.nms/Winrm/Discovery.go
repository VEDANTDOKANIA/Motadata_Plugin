package Winrm

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
)

func Discovery(credentials map[string]interface{}) {
	host := (credentials["IP_Address"]).(string)
	port := int(credentials["Port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	_, err := winrm.NewClient(endpoint, username, password)
	result := make(map[string]interface{})
	if err != nil {
		result["Error"] = "yes"
		result["Cause"] = err
	} else {
		result["Error"] = "no"
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

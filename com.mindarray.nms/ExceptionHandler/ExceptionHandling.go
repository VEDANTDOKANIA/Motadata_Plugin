package ExceptionHandler

import (
	"encoding/json"
	"fmt"
)

func ErrorHandle(credentials map[string]interface{}) {

	var data = make(map[string]interface{})
	data["ip"] = credentials["ip"]
	data["metric.group"] = credentials["metric.group"]
	error := recover()
	if error != nil {
		data["error"] = fmt.Sprintf("%v", error)
		data["status"] = "fail"
		result, _ := json.Marshal(data)
		fmt.Println(string(result))
	}

}

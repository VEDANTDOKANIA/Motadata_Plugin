package ExceptionHandler

import "fmt"

func ErrorHandle(credentials map[string]interface{}) {

	var data = make(map[string]interface{})
	data["ip"] = credentials["ip"]
	data["metric.group"] = credentials["metric.group"]
	error := recover()
	if error != nil {
		data["Panic"] = "Yes"
		data["error"] = error
		fmt.Println(data)

	}

}

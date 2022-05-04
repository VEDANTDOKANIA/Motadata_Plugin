package Winrm

import exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"

func SystemData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
}

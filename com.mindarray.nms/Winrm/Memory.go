package Winrm

import exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"

func MemoryData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
}

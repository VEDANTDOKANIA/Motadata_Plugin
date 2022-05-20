package Winrm

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"math"
	"strconv"
	"strings"
)

func DiskData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	result := make(map[string]interface{})
	if err != nil {
		result["error"] = "yes"
		result["Cause"] = err
	} else {
		result["error"] = "no"
	}
	a := "aa"
	output := ""
	ac := "Get-WmiObject win32_logicaldisk |Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size}"
	output, _, _, err = client.RunPSWithString(ac, a)
	res := strings.Split(output, "\r\n")
	var disks []map[string]interface{}
	var usedBytes int64
	var totalBytes int64
	if math.Mod(float64(len(res)), 3) != 0 {

	}
	for i := 0; i < len(res); i = i + 3 {
		disk := make(map[string]interface{})
		disk["Disk.Name"] = strings.Split(res[i], ":")[0]
		if (i+1) > len(res) || res[i+1] == "" {
			disk["disk.free.bytes"] = 0
			disk["disk.total.bytes"] = 0
			disk["disk.available.bytes"] = 0
			disk["disk.used.percent"] = 0
			disk["disk.free.percent"] = 0
			disks = append(disks, disk)
			break
		}
		bytes, _ := strconv.ParseInt(res[i+1], 10, 64)
		usedBytes = usedBytes + bytes
		disk["disk.available.bytes"], _ = strconv.ParseInt(res[i+1], 10, 64)
		bytes, _ = strconv.ParseInt(res[i+2], 10, 64)
		totalBytes = totalBytes + bytes
		disk["disk.total.bytes"] = bytes
		disk["disk.used.bytes"] = (disk["disk.total.bytes"]).(int64) - (disk["disk.available.bytes"]).(int64)
		disk["disk.used.percent"] = (float64((float64((disk["disk.total.bytes"]).(int64)) - float64((disk["disk.used.bytes"]).(int64))) / float64((disk["disk.total.bytes"].(int64))))) * 100
		disk["disk.free.percent"] = 100 - disk["disk.used.percent"].(float64)
		disks = append(disks, disk)
	}
	result["disk.total.bytes"] = totalBytes
	result["disk.used.byes"] = usedBytes
	result["disk.available.bytes"] = totalBytes - usedBytes
	result["disk.used.percent"] = ((float64(totalBytes) - float64(usedBytes)) / float64(totalBytes)) * 100
	result["disk.available.percent"] = 100.00 - (result["Disk.Used.Percent"]).(float64)
	result["disks"] = disks
	result["ip"] = credentials["ip"]
	result["metric.group"] = credentials["metric.group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

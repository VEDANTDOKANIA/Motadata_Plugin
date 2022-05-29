package SSH

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

func DiskData(credentials map[string]interface{}) {
	const cmd = "df | awk  '{if ($1 != \"Filesystem\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$5}'"
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip"].(string)
	sshPort := int(credentials["port"].(float64))
	sshUser := credentials["username"].(string)
	sshPassword := credentials["password"].(string)

	config := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, er := ssh.Dial("tcp", addr, config)
	var errors []string
	result := make(map[string]interface{})
	if er != nil {
		errors = append(errors, er.Error())
	}
	session, err := sshClient.NewSession()
	if err != nil {
		errors = append(errors, err.Error())
	}
	combo, er := session.CombinedOutput(cmd)
	if er != nil {
		errors = append(errors, er.Error())
		result["status"] = "fail"
		result["error"] = errors
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		output := string(combo)
		res := strings.Split(output, "\n")
		utilization := 0.0
		totalBytes := 0
		usedBytes := 0
		availableBytes := 0
		var disks []map[string]interface{}
		for i := 0; i < len(res)-1; i++ {
			disk := make(map[string]interface{})
			value := strings.Split(res[i], " ")
			disk["disk.name"] = value[0]
			total, _ := (strconv.ParseInt(value[1], 10, 64))
			totalBytes = int(int64(totalBytes) + total*1024)
			disk["disk.bytes.total"] = total * 1024
			used, _ := (strconv.ParseInt(value[2], 10, 64))
			usedBytes = int(int64(usedBytes) + used*1024)
			disk["disk.bytes.used"] = used * 1024
			available, _ := (strconv.ParseInt(value[3], 10, 64))
			availableBytes = int(int64(availableBytes) + available*1024)
			disk["disk.bytes.available"] = available * 1024
			usedPercent, _ := (strconv.ParseInt(strings.Split(value[4], "%")[0], 10, 64))
			disk["disk.use.percent"] = usedPercent
			disk["disk.free.percent"] = 100 - usedPercent
			disks = append(disks, disk)
		}
		result["disks"] = disks
		result["disk.total.bytes"] = totalBytes
		result["disk.used.bytes"] = usedBytes
		result["disk.available.bytes"] = availableBytes
		utilization = ((float64(totalBytes-availableBytes) / float64(totalBytes)) * 100)
		result["disk.utilization.percent"] = utilization
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}

}

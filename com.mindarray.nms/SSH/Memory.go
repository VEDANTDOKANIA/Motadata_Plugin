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

func MemoryData(credentials map[string]interface{}) {
	const cmd = "free -b | awk  '{if ($1 != \"total\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$7}'"
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
		memoryValue := strings.Split(res[0], " ")
		totalBytes, _ := strconv.ParseInt(memoryValue[1], 10, 64)
		result["memory.total.bytes"] = totalBytes
		usedBytes, _ := strconv.ParseInt(memoryValue[2], 10, 64)
		result["memory.used.bytes"] = usedBytes
		result["memory.free.bytes"], _ = strconv.ParseInt(memoryValue[3], 10, 64)
		result["memory.available.bytes"], _ = strconv.ParseInt(memoryValue[4], 10, 64)
		swapValue := strings.Split(res[1], " ")
		result["memory.swap.total.bytes"], _ = strconv.ParseInt(swapValue[1], 10, 64)
		result["memory.swap.used.bytes"], _ = strconv.ParseInt(swapValue[2], 10, 64)
		result["memory.swap.free.bytes"], _ = strconv.ParseInt(swapValue[3], 10, 64)
		usedPercent := float64(float64(totalBytes-usedBytes) / float64(totalBytes))
		result["memory.used.percent"] = usedPercent
		result["memory.available.percent"] = 100 - usedPercent

		result["status"] = "success"
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

package SSH

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func ProcessData(credentials map[string]interface{}) {
	const cmd = "ps -aux | awk  '{if ($1 != \"USER\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$11}'"
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
		var processes []map[string]interface{}
		for i := 0; i < len(res)-1; i++ {
			processValue := make(map[string]interface{})
			value := strings.Split(res[i], " ")
			processValue["process.user"] = value[0]
			processValue["process.id"] = value[1]
			processValue["process.cpu.percent"] = value[2]
			processValue["process.memory.percent"] = value[3]
			processValue["process.command"] = value[4]
			processes = append(processes, processValue)
		}
		result["processes"] = processes
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

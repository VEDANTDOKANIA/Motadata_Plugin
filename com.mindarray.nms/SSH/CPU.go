package SSH

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func CpuData(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	const cmd = "mpstat -P ALL |awk  '{if ($4 != \"CPU\") print $4 \" \" $5 \" \" $7 \" \" $14}'"
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
		system := strings.Split(res[2], " ")
		result["system.cpu.user.percent"] = system[1]
		result["system.cpu.system.percent"] = system[2]
		result["system.cpu.idle.percent"] = system[3]
		var cores []map[string]interface{}
		for i := 3; i < len(res)-1; i++ {
			core := make(map[string]interface{})
			value := strings.Split(res[i], " ")
			core["core.name"] = value[0]
			core["core.user.percent"] = value[1]
			core["core.system.percent"] = value[2]
			core["core.idle.percent"] = value[3]
			cores = append(cores, core)
		}
		result["cores"] = cores
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}

}

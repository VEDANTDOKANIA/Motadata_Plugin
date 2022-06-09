package SSH

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func SystemData(credentials map[string]interface{}) {
	const cmd = "uname -a | awk  '{ print $1 \" \" $2  \" \" $4 \" \"$6 \" \" $7 \" \" $8 \" \"$9 }' && vmstat | awk  '{if ($1 != \"procs\" && $1 !=\"r\") print $1 \" \" $2 \" \"  $12}'"
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
		systemValue := strings.Split(res[0], " ")
		result["system.os.name"] = systemValue[0]
		result["system.user.name"] = systemValue[1]
		result["system.os.version"] = systemValue[2]
		result["system.up.time"] = systemValue[3] + " " + systemValue[4] + " " + systemValue[5] + " " + systemValue[6]
		processValue := strings.Split(res[1], " ")
		result["system.running.processes"] = processValue[0]
		result["system.blocking.processes"] = processValue[1]
		result["system.context.switching"] = processValue[2]

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

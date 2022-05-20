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
	defer exception.ErrorHandle(credentials)
	const cmd = "uname -a | awk  '{ print $1 \" \" $2  \" \" $4 \" \"$6 \" \" $7 \" \" $8 \" \"$9 }' && vmstat | awk  '{if ($1 != \"procs\" && $1 !=\"r\") print $1 \" \" $2 \" \"  $12}'"
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

	result := make(map[string]interface{})
	if er != nil {
		result["error"] = "yes"
		result["Cause"] = er
	} else {
		result["error"] = "no"
	}
	session, err := sshClient.NewSession()

	if err != nil {
		result["error"] = "yes"
		result["Cause"] = er
	} else {
		result["error"] = "no"
	}

	combo, er := session.CombinedOutput(cmd)
	output := string(combo)
	res := strings.Split(output, "\n")
	systemValue := strings.Split(res[0], " ")
	result["system.os.name"] = systemValue[0]
	result["system.user.name"] = systemValue[1]
	result["system.os.version"] = systemValue[2]
	result["system.up.time"] = systemValue[3] + " " + systemValue[4] + " " + systemValue[5] + " " + systemValue[6]

	processValue := strings.Split(res[1], " ")
	result["system.running.process"] = processValue[0]
	result["system.blocking.process"] = processValue[1]
	result["system.context.switching"] = processValue[2]
	result["ip"] = credentials["ip"]
	result["metric.group"] = credentials["metric.group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

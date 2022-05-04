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
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip.address"].(string)
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
	combo, er := session.CombinedOutput("ps -aux | awk  '{if ($1 != \"USER\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$11}'")
	output := string(combo)
	res := strings.Split(output, "\n")
	var processes []map[string]interface{}
	for i := 0; i < len(res)-1; i++ {
		//cpu := make(map[string]interface{})
		processValue := make(map[string]interface{})
		value := strings.Split(res[i], " ")
		processValue["Process.User"] = value[0]
		processValue["Process.ID"] = value[1]
		processValue["Process.CPU.Percent"] = value[2]
		processValue["Process.Memory.Percent"] = value[3]
		processValue["Process.Command"] = value[4]
		processes = append(processes, processValue)
	}
	result["Process"] = processes
	result["ip.address"] = credentials["ip.address"]
	result["metric.group"] = credentials["metric.group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

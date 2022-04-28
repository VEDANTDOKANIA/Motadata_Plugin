package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func SystemData(credentials map[string]interface{}) {
	const cmd = "uname -a | awk  '{ print $1 \" \" $2  \" \" $4 \" \"$6 \" \" $7 \" \" $8 \" \"$9 }' && vmstat | awk  '{if ($1 != \"procs\" && $1 !=\"r\") print $1 \" \" $2 \" \"  $12}'"
	sshHost := credentials["IP_Address"].(string)
	sshPort := int(credentials["Port"].(float64))
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
		result["Error"] = "yes"
		result["Cause"] = er
	} else {
		result["Error"] = "no"
	}
	session, err := sshClient.NewSession()

	if err != nil {
		result["Error"] = "yes"
		result["Cause"] = er
	} else {
		result["Error"] = "no"
	}

	combo, er := session.CombinedOutput(cmd)
	output := string(combo)
	res := strings.Split(output, "\n")
	systemValue := strings.Split(res[0], " ")
	result["System.OS.Name"] = systemValue[0]
	result["System.User.Name"] = systemValue[1]
	result["System.OS.Version"] = systemValue[2]
	result["System.Up.Time"] = systemValue[3] + " " + systemValue[4] + " " + systemValue[5] + " " + systemValue[6]

	processValue := strings.Split(res[1], " ")
	result["System.Running.Process"] = processValue[0]
	result["System.Blocking.Process"] = processValue[1]
	result["System.Context.Switching"] = processValue[2]
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

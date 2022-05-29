package SSH

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Discovery(credentials map[string]interface{}) {
	const cmd = "uname -n"
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip"].(string)
	sshPort := int(credentials["port"].(float64))
	sshUser := credentials["username"].(string)
	sshPassword := credentials["password"].(string)
	config := &ssh.ClientConfig{
		Timeout:         6 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	var errors []string
	sshClient, er := ssh.Dial("tcp", addr, config)
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
	}
	output := string(combo)
	if len(output) == 0 {
		errors = append(errors, "unable to gather hostname wrong credentials")
	}
	result := make(map[string]interface{})
	if len(errors) > 0 {
		result["status"] = "fail"
		result["error"] = errors
	} else {
		result["status"] = "success"
		result["host"] = strings.Split(output, "\n")[0]
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

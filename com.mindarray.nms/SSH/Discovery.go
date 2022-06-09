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
		Timeout:         4 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	var errors []string
	result := make(map[string]interface{})
	sshClient, er := ssh.Dial("tcp", addr, config)
	if er != nil {
		if strings.Contains(er.Error(), "unable to authenticate, attempted methods [none password], no supported methods remain") {
			errors = append(errors, "wrong username or password")
		} else if strings.Contains(er.Error(), "network is unreachable") {
			errors = append(errors, "wrong ip address")
		} else if strings.Contains(er.Error(), "connection refused") {
			errors = append(errors, "wrong ip address or port")
		} else {
			errors = append(errors, "wrong credentials provided")
		}
	} else {
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
		} else {
			result["host"] = strings.Split(output, "\n")[0]
		}
	}

	if len(errors) > 0 {
		result["status"] = "fail"
		result["error"] = errors
	} else {
		result["status"] = "success"

	}
	data, error := json.Marshal(result)
	if error != nil {
		out := make(map[string]interface{})
		out["status"] = "fail"
		out["error"] = error.Error()
		output, _ := json.Marshal(out)
		fmt.Print(string(output))
	} else {
		fmt.Print(string(data))
	}

}

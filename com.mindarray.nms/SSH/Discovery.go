package SSH

import (
	exception "MotadataPlugin/com.mindarray.nms/ExceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip.address"].(string)
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
	_, er := ssh.Dial("tcp", addr, config)

	result := make(map[string]interface{})
	if er != nil {
		result["status"] = "Unsuccessful"
		result["error"] = er.Error()
	} else {
		result["status"] = "successful"
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

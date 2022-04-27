package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

func Discovery(credentials map[string]interface{}) {
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
	_, er := ssh.Dial("tcp", addr, config)

	result := make(map[string]interface{})
	if er != nil {
		result["Error"] = "yes"
		result["Cause"] = er
	} else {
		result["Error"] = "no"
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))

}

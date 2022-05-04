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

	combo, er := session.CombinedOutput(cmd)
	output := string(combo)
	res := strings.Split(output, "\n")
	system := strings.Split(res[2], " ")
	result["System.Cpu.User.Percent"] = system[1]
	result["System.Cpu.System.Percent"] = system[2]
	result["System.Cpu.Idle.Percent"] = system[3]
	var cores []map[string]interface{}
	for i := 3; i < len(res)-1; i++ {
		//cpu := make(map[string]interface{})
		core := make(map[string]interface{})
		value := strings.Split(res[i], " ")
		core["Core.Name"] = value[0]
		core["Core.User.Percent"] = value[1]
		core["Core.System.Percent"] = value[2]
		core["Core.Idle.Percent"] = value[3]
		cores = append(cores, core)
	}
	result["Cores"] = cores
	result["ip.address"] = credentials["ip.address"]
	result["metric.group"] = credentials["metric.group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func CpuData(credentials map[string]interface{}) {
	const cmd = "mpstat -P ALL |awk  '{if ($4 != \"CPU\") print $4 \" \" $5 \" \" $7 \" \" $14}'"
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
	system := strings.Split(res[2], " ")
	result["System.Cpu.User.Percent"] = system[1]
	result["System.Cpu.System.Percent"] = system[2]
	result["System.Cpu.Idle.Percent"] = system[3]
	var cores []map[string]interface{}
	for i := 3; i < len(res)-1; i++ {
		//cpu := make(map[string]interface{})
		core := make(map[string]interface{})
		value := strings.Split(res[i], " ")
		core["Cpu.Name"] = value[0]
		core["Cpu.User.Percent"] = value[1]
		core["Cpu.System.Percent"] = value[2]
		core["Cpu.Idle.Percent"] = value[3]
		cores = append(cores, core)
	}
	result["Cores"] = cores
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

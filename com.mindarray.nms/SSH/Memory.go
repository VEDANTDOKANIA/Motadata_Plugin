package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

func MemoryData(credentials map[string]interface{}) {
	const cmd = "free -b | awk  '{if ($1 != \"total\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$7}'"
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

	memoryValue := strings.Split(res[0], " ")
	totalBytes, _ := strconv.ParseInt(memoryValue[1], 10, 64)
	result["Memory.Total.Bytes"] = totalBytes
	usedBytes, _ := strconv.ParseInt(memoryValue[2], 10, 64)
	result["Memory.Used.Bytes"] = usedBytes
	result["Memory.Free.Bytes"], _ = strconv.ParseInt(memoryValue[3], 10, 64)
	result["Memory.Available.Bytes"], _ = strconv.ParseInt(memoryValue[4], 10, 64)

	swapValue := strings.Split(res[1], " ")
	result["Memory.Swap.Total.Bytes"], _ = strconv.ParseInt(swapValue[1], 10, 64)
	result["Memory.Swap.Used.Bytes"], _ = strconv.ParseInt(swapValue[2], 10, 64)
	result["Memory.Swap.Free.Bytes"], _ = strconv.ParseInt(swapValue[3], 10, 64)
	usedPercent := float64(float64(totalBytes-usedBytes) / float64(totalBytes))
	result["Memory.Used.Percent"] = usedPercent
	result["Memory.Available.Percent"] = 100 - usedPercent
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

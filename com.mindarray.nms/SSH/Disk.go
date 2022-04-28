package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

func DiskData(credentials map[string]interface{}) {
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
	combo, er := session.CombinedOutput("df | awk  '{if ($1 != \"Filesystem\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$5}'")
	output := string(combo)
	res := strings.Split(output, "\n")
	utilization := 0.0
	totalBytes := 0
	usedBytes := 0
	availableBytes := 0
	var disks []map[string]interface{}
	for i := 0; i < len(res)-1; i++ {
		disk := make(map[string]interface{})
		value := strings.Split(res[i], " ")
		disk["Disk.Name"] = value[0]
		total, _ := (strconv.ParseInt(value[1], 10, 64))
		totalBytes = int(int64(totalBytes) + total*1024)
		disk["Disk.Bytes.Total"] = total * 1024
		used, _ := (strconv.ParseInt(value[2], 10, 64))
		usedBytes = int(int64(usedBytes) + used*1024)
		disk["Disk.Bytes.Used"] = used * 1024
		available, _ := (strconv.ParseInt(value[3], 10, 64))
		availableBytes = int(int64(availableBytes) + available*1024)
		disk["Disk.Bytes.Available"] = available * 1024
		usedPercent, _ := (strconv.ParseInt(strings.Split(value[4], "%")[0], 10, 64))
		disk["Disk.Use.Percent"] = usedPercent
		disk["Disk.Free.Percent"] = 100 - usedPercent
		disks = append(disks, disk)
	}
	result["Disk"] = disks
	result["Disk.Total.Bytes"] = totalBytes
	result["Disk.Used.Bytes"] = usedBytes
	result["Disk.Available.Bytes"] = availableBytes
	utilization = ((float64(totalBytes-availableBytes) / float64(totalBytes)) * 100)
	result["Disk.Utilization.Percent"] = utilization
	result["IP_Address"] = credentials["IP_Address"]
	result["Metric_Group"] = credentials["Metric_Group"]
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

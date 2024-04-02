package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

func IsHostFound(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("%s: %s", target, err.Error())
		return false
	}

	return true
}

// 192.168.0.10:5000
// 192.168.0.11:5000
// 192.168.0.12:5000
// 192.168.0.10:5001
// 192.168.0.10:5002
// 192.168.0.10:5003

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbours(myHost string, myPort uint16, startIp uint8, endIp uint8,
	startPort uint16, endPort uint16) []string {

	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbours := make([]string, 0)

	for port := startPort; port <= endPort; port++ {
		for ip := startIp; ip <= endIp; ip++ {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsHostFound(guessHost, port) {
				neighbours = append(neighbours, guessTarget)
			}
		}
	}
	return neighbours
}

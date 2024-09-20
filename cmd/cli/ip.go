package main

import (
	"strconv"
	"strings"
)

type IP [4]byte
type IP3 [3]byte

func (ip IP) ToIP3() IP3 {
	return IP3{ip[1], ip[2], ip[3]}
}

func parseIP(strIp string) IP {
	parts := strings.Split(strIp, ".")
	var ip IP

	for i := range 4 {
		octet, err := strconv.Atoi(parts[i])
		if err != nil {
			panic(err)
		}
		ip[i] = byte(octet)
	}

	return ip
}

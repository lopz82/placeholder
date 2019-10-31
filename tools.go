package main

import (
	"log"
	"net"
	"os"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func GetHostName() (string, error) {
	h, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return h, nil
}

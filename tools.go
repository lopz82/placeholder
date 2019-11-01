package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sort"
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

func sortHeaders(h http.Header) []string {
	l := make([]string, 0)
	for k := range h {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

func getLongestHeader(h http.Header) int {
	max := 0
	for k := range h {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}

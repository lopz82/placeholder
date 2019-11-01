package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"golang.org/x/sys/unix"
)

const (
	_         = iota
	KB uint64 = 1 << (10 * iota)
	MB
	GB
	TB
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	h, err := GetHostName()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(w, h)
}

func ReturnHeadersHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	if len(param) == 0 {
		http.Error(w,
			`Please, use parameters to set the header and the values to be returned.
Values of existing headers are appended.
Example: http://host.com//headers?header1=value1&header2=value2`,
			400)
		return
	}
	fmt.Fprintln(w, "Requested headers")
	for header, values := range param {
		for _, value := range values {
			w.Header().Add(header, value)
			fmt.Fprintf(w, "%-s: %v\n", header, value)
		}
	}
}

func ClientRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Remote Address: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "Requested URI: %s\n", r.RequestURI)
	fmt.Fprintln(w, "Headers:")
	sortH := sortHeaders(r.Header)
	max := getLongestHeader(r.Header)
	for _, h := range sortH {
		fmt.Fprintf(w, "  %*s: %v\n", max, h, r.Header.Get(h))
	}
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,
		`Please, use the following urls to access to the available tools:
/request -> Shows your request headers
/memory -> Shows host information
/return -> Returns the specified HTTP response status code
/headers -> Returns the specified headers
/health -> Health check
/interfaces -> Host information`)
}

func SystemHandler(w http.ResponseWriter, r *http.Request) {
	info := &unix.Sysinfo_t{}
	err := unix.Sysinfo(info)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, "Loads: %v\n", info.Loads)
	fmt.Fprintf(w, "Uptime: %d\n", info.Uptime)
	fmt.Fprintf(w, "Total RAM: %d MB\n", info.Totalram/MB)
	fmt.Fprintf(w, "Free RAM: %d MB\n", info.Freeram/MB)
	fmt.Fprintf(w, "Shared RAM: %d MB\n", info.Sharedram/MB)
	fmt.Fprintf(w, "Buffer RAM: %d MB\n", info.Bufferram/MB)
	fmt.Fprintf(w, "Total Swap: %d MB\n", info.Totalswap/MB)
	fmt.Fprintf(w, "Free Swap: %d MB\n", info.Freeswap/MB)
	fmt.Fprintf(w, "Processes: %d\n", info.Procs)
}

func InterfacesHandler(w http.ResponseWriter, r *http.Request) {
	h, err := GetHostName()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, "Hostname: %s\n", h)
	fmt.Fprintf(w, "IP Addr: %s\n\n", GetOutboundIP())

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintln(w, "Interfaces:\n")
	for _, i := range interfaces {
		fmt.Fprintf(w, "Name: %s\n", i.Name)
		fmt.Fprintf(w, "  Addr: %s\n", i.HardwareAddr)
		fmt.Fprintf(w, "  Index: %d\n", i.Index)
		fmt.Fprintf(w, "  MTU: %v\n", i.MTU)
		fmt.Fprintf(w, "  Flags: %v\n\n", i.HardwareAddr)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Healthcheck")
	w.WriteHeader(200)
}

func ReturnCodeHandler(w http.ResponseWriter, r *http.Request) {
	strcode := r.URL.Query().Get("code")
	if strcode == "" {
		http.Error(w,
			`Please, use a 'code' parameter to set the response status code to be returned.
The response status code should be a integer between 100 and 599.
Example: http://host.com//return?code=404`,
			400)
		return
	}
	code, err := strconv.Atoi(strcode)
	if err != nil {
		http.Error(w,
			`Please, the response status code should be a integer between 100 and 599.
Check https://developer.mozilla.org/en-US/docs/Web/HTTP/Status for more info.`,
			400)
		return
	}
	if (code < 100) || (code > 599) {
		http.Error(w,
			fmt.Sprintf(`%d is not a valid HTTP response status code.
Check https://developer.mozilla.org/en-US/docs/Web/HTTP/Status for more info.`, code),
			400)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, "A %d HTTP response status code was sent.", code)
}

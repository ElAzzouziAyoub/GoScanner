package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"goscan/Ascii"
)

var wg sync.WaitGroup

type PortService struct {
    Port        int    `json:"port"`
    IsTCP       bool   `json:"is_tcp"`
    IsUDP       bool   `json:"is_udp"`
    Description string `json:"description"`
    IsOfficial  bool   `json:"is_official"`
}


func main() {

	ascii_art.PrintAsciiArt()

	data, err := os.ReadFile("ports.json")
	if err != nil {
    log.Fatal(err)
	}

	var services []PortService

	err = json.Unmarshal(data, &services)
	if err != nil {
    log.Fatal(err)
	}

	serviceMap := make(map[int]PortService)

	for _, svc := range services {
    serviceMap[svc.Port] = svc
	}




	target := flag.String("target", "", "Target host")
	StartPort := flag.Int("sp", 80, "Starting Port to scan")
	EndPort := flag.Int("ep", 80, "Ending Port to scan")
	Type := flag.String("type","tcp","TCP/UDP")

	flag.Parse()

	if *target == "" {
		fmt.Println("Error: target is required")
		flag.Usage()
		return
	}
	if *StartPort < 1 || *EndPort > 65535 {
		fmt.Println("Error : invalid port range")
		flag.Usage()
		return
	} 

	if *StartPort > *EndPort {
		fmt.Println("Error : starting port must be less than ending port")
		flag.Usage()
		return
	}

	fmt.Println("--------------------------------")

	for i := *StartPort; i <= *EndPort; i++ {
		wg.Add(1)
		go func(a int) {

			start := time.Now()
			defer wg.Done()
			conn, err := net.DialTimeout(*Type,fmt.Sprintf("%s:%d",*target,a),2*time.Second)
			port := serviceMap[a]

			state := "unknown"

			if err == nil {
    		state = "open"
    		conn.Close()
			} else {
    	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
        	state = "filtered"
    		}else {
        	state = "closed"
    		}
			}



			if err == nil {
				fmt.Println("Port : ",a)
				fmt.Println("State : ",state)
				fmt.Println("Description :",port.Description)
				if port.IsTCP {
					fmt.Println("Type : TCP")
				}else{
					fmt.Println("Type : UDP")
				}
				duration := time.Since(start)
				fmt.Println("Duration :",duration.Milliseconds(),"ms")
				fmt.Println("IsOfficial :",port.IsOfficial)
				fmt.Println("--------------------------------")
			}

		}(i)
	}
	wg.Wait()
}

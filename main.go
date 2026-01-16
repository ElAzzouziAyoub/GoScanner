package main

import (
	"fmt"
	"net"
	"sync"
	"flag"
)

var wg sync.WaitGroup


func main() {

	target := flag.String("target", "", "Target host")
	StartPort := flag.Int("sp", 80, "Starting Port to scan")
	EndPort := flag.Int("ep", 80, "Ending Port to scan")

	flag.Parse()

	if *target == "" {
		fmt.Println("Error: target is required")
		flag.Usage()
		return
	}

	if *StartPort > *EndPort {
		fmt.Println("Error : starting port must be less than ending port")
		flag.Usage()
		return
	}

	for i := *StartPort; i <= *EndPort; i++ {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			_, err := net.Dial("tcp",fmt.Sprintf("%s:%d",*target,a))
			if err == nil {
				fmt.Println("Port Opened : ",a)
			}

		}(i)
	}
	wg.Wait()
}

package main

import (
	"fmt"
	"net"
	"sync"
)

var wg sync.WaitGroup


func main() {

	for i := 0; i <= 1000; i++ {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			_, err := net.Dial("tcp",fmt.Sprintf("scanme.nmap.org:%d",a))
			if err == nil {
				fmt.Println("Port Opened : ",a)
			}

		}(i)
	}
	wg.Wait()
}

package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, result chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			result <- 0
			continue
		}
		conn.Close()
		result <- p
	}
}

func main() {
	ports := make(chan int, 100)
	result := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, result)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-result
		if port != 0 {
			fmt.Printf("Obtained result is %d\n", port)
			openports = append(openports, port)
		}
		fmt.Printf("Current iteration is %d\n", i)
	}

	close(ports)
	close(result)
	sort.Ints(openports)

	for _, port := range openports {
		fmt.Printf("%d is open\n", port)
	}
}

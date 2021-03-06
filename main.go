package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var tcp200Ports = []int{
	1,3,7,9,13,17,19,21,22,23,25,26,37,
	53,79,80,81,82,88,100,106,110,111,
	113,119,135,139,143,144,179,199,
	254,255,280,311,389,427,443,444,
	445,464,465,497,513,514,515,543,
	544,548,554,587,593,625,631,636,
	646,787,808,873,902,990,993,995,
	1000,1022,1024,1025,1026,1027,1028,
	1029,1030,1031,1032,1033,1035,1036,
	1037,1038,1039,1040,1041,1044,1048,
	1049,1050,1053,1054,1056,1058,1059,
	1064,1065,1066,1069,1071,1074,1080,
	1110,1234,1433,1494,1521,1720,1723,
	1755,1761,1801,1900,1935,1998,2000,
	2001,2002,2003,2005,2049,2103,2105,
	2107,2121,2161,2301,2383,2401,2601,
	2717,2869,2967,3000,3001,3128,3268,
	3306,3389,3689,3690,3703,3986,4000,
	4001,4045,4899,5000,5001,5003,5009,
	5050,5051,5060,5101,5120,5190,5357,
	5432,5555,5631,5666,5800,5900,5901,
	6000,6001,6002,6004,6112,6646,6666,
	7000,7070,7937,7938,8000,8002,8008,
	8009,8010,8031,8080,8081,8443,8888,
	9000,9001,9090,9100,9102,9999,10000,
	10001,10010,32768,32771,49152,49153,
	49154,49155,49156,49157,50000,
}

var udp200Ports = []int{
	7,9,13,17,19,21,22,23,37,42,49,53,67,68,
	69,80,88,111,120,123,135-139,158,161,162,
	177,192,199,389,407,427,443,445,464,497,
	500,514,515,517,518,520,593,623,626,631,
	664,683,800,989-990,996,997,998,999,1001,
	1008,1019,1021,1022,1023,1024,1025,1026,
	1027,1028,1029,1030,1031,1032,1033,1034,
	1036,1038,1039,1041,1043,1044,1045,1049,
	1068,1419,1433,1434,1645,1646,1701,1718,
	1719,1782,1812,1813,1885,1900,2000,2002,
	2048,2049,2148,2222,2223,2967,3052,3130,
	3283,3389,3456,3659,3703,4000,4045,4444,
	4500,4672,5000,5001,5060,5093,5351,5353,
	5355,5500,5632,6000,6001,6346,7938,9200,
	9876,10000,10080,11487,16680,17185,19283,
	19682,20031,22986,27892,30718,31337,32768,
	32769,32769,32770,32771,32772,32773,32815,
	33281,33354,34555,34861,34862,37444,39213,
	41524,44968,49152,49153,49154,49156,49158,
	49159,49162,49163,49165,49166,49168,49171,
	49172,49179,49180,49181,49182,49184,49185,
	49186,19187,19188,19189,19190,19191,19192,
	19193,19194,19195,19196,49199,49200,49201,
	49202,49205,49208,49209,49210,49211,58002,
	65024,
}


func main() {
	var wg sync.WaitGroup

	isTCP := true
	isUDP := true

	if len(os.Args) < 2 {
		fmt.Println("Not enough args")
		return
	}
	start := time.Now()
	log.Printf("I have started at %s\n", start)
	wg.Add(2)

	go func() {
		scanning(os.Args[1], tcp200Ports, isTCP)
		wg.Done()
	}()

	go func() {
		scanning(os.Args[1], udp200Ports, isUDP)
		wg.Done()
	}()

	wg.Wait()
	elapsed := time.Since(start) / time.Second
	log.Printf("%s this much time passed\n", elapsed)
}

func scanning(url string, toports []int, isTCP bool) {
	ports := make(chan int, 100)
	result := make(chan int)
	var openports []int
	var method string

	for i := 0; i < cap(ports); i++ {
		go worker(ports, result, url, isTCP)
	}

	go func() {
		for _, i := range toports {
			ports <- i
		}
	}()

	for i := 0; i < len(toports); i++ {
		port := <-result
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(result)
	sort.Ints(openports)

	if isTCP {
		method = "TCP"
	} else {
		method = "UDP"
	}

	for _, port := range openports {
		fmt.Printf("%v: %d is open\n", method, port)
	}
}

func worker(ports, result chan int, url string, isTCP bool) {
	var method string

	if isTCP {
		method = "tcp"
	} else {
		method = "udp"
	}

	for p := range ports {
		address := fmt.Sprintf("%v:%d", url, p)
		conn, err := net.DialTimeout(method, address, 5 * time.Second)
		if err != nil {
			result <- 0
			continue
		}
		defer conn.Close()
		result <- p
	}
}

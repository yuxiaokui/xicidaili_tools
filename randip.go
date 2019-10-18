package main

import (
	"fmt"
	"math/rand"
	"time"
)


func getIpaddr() string {
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func getPort() string {
	port := fmt.Sprintf("%d", rand.Intn(65535))
	return port
}



func main() {
	rand.Seed(time.Now().Unix())
	for i := 0;i< 10000;i++ {
		fmt.Println(getIpaddr() + ":" + getPort())
	}
}

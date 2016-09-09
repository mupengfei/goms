package main

import (
	. "fmt"
	"net"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	channel := make(chan int, 1024)
	for i := 2049; i < 22049; i++ {
		go func() {
			tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:2046")
			localAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:"+string(i))
			conn, err := net.DialTCP("tcp", localAddr, tcpAddr)
			//			conn, err := net.Dial("tcp", "127.0.0.1:2046")
			defer func() {
				if conn != nil {
					conn.Close()
				}
			}()
			if err != nil {
				//		panic("error")
			} else {
				channel <- 0
				conn.Write([]byte("Hello World"))
				buf := make([]byte, 1024)
				for {
					lenght, err := conn.Read(buf)
					if err != nil {
						conn.Close()
						break
					}
					if lenght > 0 {
						buf[lenght] = 0
					}
					//fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
					//					reciveStr := string(buf[0:lenght])
					//					Println(reciveStr)
				}
			}
		}()
	}
	go colConNum(channel)
	Println("END")
	for {
		i := 2
		i++
	}
}

func colConNum(channel chan int) {
	num := 0
	for {
		select {
		case <-channel:
			{
				num++
				Println(num)
			}
		}
	}
}

package main

import (
	. "fmt"
	"net"
	"runtime"

	log "github.com/thinkboy/log4go"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.LoadConfiguration("log.xml")
	defer log.Close()
	channel := make(chan int, 1024)
	for i := 2049; i < 2249; i++ {
		go func(i int) {
			Println(i)
			tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:2046")
			localAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:"+string(i))
			conn, err := net.DialTCP("tcp", localAddr, tcpAddr)
			//			conn, err := net.Dial("tcp", "127.0.0.1:2046")
			defer func() {
				if conn != nil {
					log.Error("wolegequ")
					conn.Close()
				}
			}()
			if err != nil {
				//		panic("error")
			} else {
				//				channel <- 0
				conn.Write([]byte("Hello World"))
				buf := make([]byte, 1024)
				for {
					log.Info(i)
					lenght, err := conn.Read(buf)
					if err != nil {
						conn.Close()
						log.Error(err)
						break
					}
					if lenght > 0 {
						buf[lenght] = 0
					}
					//fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
					log.Info("Rec[", conn.RemoteAddr().String(), "  ", i, "] Say :", string(buf[0:lenght]))
					//					reciveStr := string(buf[0:lenght])
					//					Println(reciveStr)
				}
			}
		}(i)
	}
	go colConNum(channel)
	select {}
	Println("END")
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

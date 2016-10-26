package main

import (
	"net"
	"runtime"
	"sync"
	"time"

	log "github.com/thinkboy/log4go"
)

var locker sync.Mutex

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.LoadConfiguration("log.xml")
	defer log.Close()
	StartServer("2046")
}

func StartServer(port string) {
	service := ":" + port //strconv.Itoa(port);
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	conns := make(map[string]net.Conn)
	//	messages := make(chan string, 10)
	//启动服务器广播线程
	//	go echoHandler(&conns, messages)
	var i int = 0
	for {
		log.Info("Listening ...")
		conn, err := l.Accept()
		checkError(err, "Accept")
		//		log.Info("Accepting ...")
		conns[conn.RemoteAddr().String()] = conn
		i++
		//启动一个新线程
		//		go Handler(conn, messages)
		go func(conn *net.Conn) {
			for {
				_, err := (*conn).Write([]byte(time.Now().String()))
				log.Info(time.Now().String(), i)
				if err != nil {
					log.Info(string(i) + "*")
					//					locker.Unlock()
					break
				}
				time.Sleep(time.Second * 2)
			}
		}(&conn)
	}
}

func echoHandler(conns *map[string]net.Conn, messages chan string) {
	for {
		msg := <-messages
		//		log.Info(msg)
		for key, value := range *conns {
			//			log.Info("connection is connected from ...", key)
			locker.Lock()
			_, err := value.Write([]byte(msg))
			if err != nil {
				log.Error(err.Error())
				delete(*conns, key)
			}
			locker.Unlock()
		}
	}
}

func Handler(conn net.Conn, messages chan string) {
	//	log.Info("connection is connected from ...", conn.RemoteAddr().String())
	buf := make([]byte, 1024)
	for {
		locker.Lock()
		lenght, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}
		//fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
		reciveStr := string(buf[0:lenght])
		messages <- reciveStr
		locker.Unlock()
	}
}

func checkError(error error, info string) {
	if error != nil {
		//		log.Error("ERROR: " + info + " " + error.Error()) // terminate program
	}
}

package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Port int

	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//broadcast channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	Server := &Server{
		Ip : ip,
		Port : port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return Server
}

func (s *Server) ListenMessage(){
	for{
		msg := <- s.Message
		s.mapLock.Lock()
		for _,cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) BroadCast(user *User, msg string)  {
	sendMsg := "["+user.Addr+"]"+user.Name+":"+msg
	s.Message<-sendMsg
}

func (s *Server) Handler(conn net.Conn)  {

	user := NewUser(conn, s)
	isLive := make(chan bool)

	user.Online()

	//读取用户msg
	go func() {
		buf := make([]byte,4096)
		for{
			n, err := conn.Read(buf)
			if n==0{
				user.Offline()
				return
			}
			if err != nil&&err != io.EOF {
				fmt.Println("Conn Read err:",err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)

			//active a user
			isLive <- true
		}
	}()

	for{
		select {
			case <- isLive:
			case <- time.After(time.Second*100):
				user.SendMsg("kick off")
				close(user.C)
				conn.Close()
				return
		}
	}

}

func (s *Server) Start(){
	fmt.Println("start!")
	//socket listen
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d",s.Ip,s.Port),
	)
	if err != nil {
		fmt.Println("net.Listen err:",err)
		return 
	}
	defer listener.Close()

	//启动监听
	go s.ListenMessage()

	for  {

		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:",err)
			continue
		}
		//do handler
		go s.Handler(conn)

	}
}
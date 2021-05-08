package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

func NewClient(serverIp string, serverPort int) *Client{
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
	}

	conn, err := net.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", serverIp, serverPort),
	)
	if err != nil {
		fmt.Println("net.Dial error:",err)
		return nil
	}

	client.conn = conn

	return client
}

func (c *Client) DealResponse() {
	io.Copy(os.Stdout,c.conn)
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
	fmt.Println("3.rename")
	fmt.Println("0.quit")

	fmt.Scanln(&flag)
	if flag>=0&&flag<=3{
		c.flag = flag
		return true
	}else{
		fmt.Println("error number")
		return false
	}
}
func (c *Client) SelectUsers(){
	sendMsg := "who"
	_,err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}
}
func (c *Client) PrivateChat(){
	var remoteName string
	var chatMsg string

	c.SelectUsers()
	fmt.Println("please input who you want to talk, input exit to leave")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("please input content, input exit to leave")
		fmt.Scanln(&chatMsg)

		for chatMsg!="exit"{
			if len(chatMsg) != 0{
				sendMsg := "to|"+remoteName+"|"+chatMsg
				_,err:=c.conn.Write([]byte(sendMsg))
				if err!=nil {
					fmt.Println("conn.Write",err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("please input content, input exit to leave")
			fmt.Scanln(&chatMsg)
		}

		c.SelectUsers()
		fmt.Println("please input content, input exit to leave")
		fmt.Scanln(&chatMsg)
	}

}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println("please input chat content, input 'exit' to leave chat room")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit"{

		if len(chatMsg)!=0{
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:",err)
				break
			}
		}
		fmt.Println("please input chat content, input 'exit' to leave chat room")
		fmt.Scanln(&chatMsg)
	}
}

func (c *Client) UpdateName() bool{
	fmt.Println("please input username:")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|"+c.Name+"\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:",err)
		return false
	}

	return true
}

func (c *Client) run() {
	for c.flag != 0{
		for c.menu() != true {
		}

		switch c.flag {
		case 1:
			c.PublicChat()
			fmt.Println("public chat")
		case 2:
			c.PrivateChat()
			fmt.Println("private chat")
		case 3:
			c.UpdateName()
			fmt.Println("rename")
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(
		&serverIp,
		"ip",
		"127.0.0.1",
		"指定IP地址，默认为127.0.0.1",
	)
	flag.IntVar(
		&serverPort,
		"port",
		8080,
		"指定端口号，默认为8888",
	)
}



func main() {

	flag.Parse()

	client := NewClient(serverIp,serverPort)
	if client==nil {
		fmt.Println("建立失败")
		return
	}
	fmt.Println("建立成功")

	go client.DealResponse()
	client.run()
}

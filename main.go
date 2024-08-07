package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func newClient(ws *websocket.Conn) *Client {
	return &Client{
		conn: ws,
	}
}

func (c *Client) readLoop(){
	buf := make([]byte , 1024)
	for {

		n, err := c.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				c.conn.Close()
				break
			}
			break
		}

		msg := buf[:n]

		fmt.Printf("message: %s", string(msg))
	}
}

func (c *Client) listener() error{
	reader := bufio.NewReader(os.Stdin)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)

		if err != nil {
			fmt.Println("Read error: ",err)
		}
		msg := buf[:n]
		str := string(msg)
		str = strings.TrimSpace(str)

		if str == "/close" {
			c.conn.Close()
			fmt.Println("CLOSED")
			
			return errors.New("websocket closed")
		}

		_ , err = c.conn.Write(msg)

		if err != nil {
			fmt.Println("write error: ",err)
		}
	}
}

func (c *Client) Close(){
	c.conn.Close()
}

func main(){

	ws, err := websocket.Dial("ws://178.42.237.220:42069/ws", "", "http://localhost")

	if err != nil {
		fmt.Println(err)
	}

	client := newClient(ws)

	defer client.Close()

	go func() {
			if err := client.listener(); err != nil {
				fmt.Println(err)
			}
	}()
	client.readLoop()



}
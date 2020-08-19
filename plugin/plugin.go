package plugin

import (
	"fmt"
	"net"

	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/utils"
)

//ReceiveHandler the receive handler
//clientID is the random id of the tcp connection
//data  is the origin payload
type ReceiveHandler func(clientID string, perfix string, data []byte)

//NetReceiver the interface for the tcp receiver
//Read is the only one method should be implemented! it's just run the tcp server on a port
type NetReceiver interface {
	Write(clientID string, buf []byte) error
	Listen(port int, max int)
}

type netReceiver struct {
	handler ReceiveHandler
	clients map[string]net.Conn
}

//Read
//start to open a socket server on the port
//use go routine to read, call the handler after received the message data.
func (r *netReceiver) Listen(port int, max int) {
	go func() {
		// listen on all interfaces
		ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))

		// run loop forever (or until ctrl-c)
		for {
			// accept connection on port
			conn, _ := ln.Accept()
			clientID := utils.GenShortID()
			r.clients[clientID] = conn
			go func() {
				for {
					buf := make([]byte, max)
					reqLen, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Error to read message because of ", err)
						continue
					}
					if reqLen < 10 {
						fmt.Println("too short")
						continue
					}
					data := buf[0:reqLen]
					perfix := fmt.Sprintf("%x", buf[0:2])
					fmt.Printf("perfix: %s\n", perfix)
					// output message received
					go r.handler(clientID, perfix, data)
				}
			}()

		}
	}()

}

func (r *netReceiver) Write(clientID string, buf []byte) error {
	return nil
}

//NewNetReceiver create a new receiver
func NewNetReceiver(f ReceiveHandler) NetReceiver {
	return &netReceiver{
		handler: f,
		clients: make(map[string]net.Conn),
	}
}

func init() {
	fpm.Register(func(app *fpm.Fpm) {
		// 配置 socket
		if !app.HasConfig("socket") {
			panic("socket config node required")
		}
		socketConfig := app.GetConfig("socket").(map[string]interface{})

		port := socketConfig["port"].(float64)
		app.Logger.Debugf("Socket Config port: %d", port)

		server := NewNetReceiver(func(clientID string, perfix string, data []byte) {
			// publish here

			app.Publish("#tcp/receive/"+perfix, map[string]interface{}{
				"clientID": clientID,
				"data":     data,
			})
			app.Logger.Infof("tcp receive: %s, %v", clientID, data)
		})
		server.Listen(int(port), 1024)

	})
}

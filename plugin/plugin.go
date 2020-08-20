package plugin

import (
	"errors"
	"fmt"
	"net"
	"strings"

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
	Listen()
	Clients() map[string]string
	SetID(id, connID string) (bool, error)
}

type Options struct {
	Port   int
	Max    int
	Perfix []string
}
type netReceiver struct {
	options *Options
	handler ReceiveHandler
	clients map[string]net.Conn
	ids     map[string]string
}

//Read
//start to open a socket server on the port
//use go routine to read, call the handler after received the message data.
func (r *netReceiver) Listen() {
	port := r.options.Port
	max := r.options.Max
	perfixArr := r.options.Perfix
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
						if err.Error() == "EOF" {
							// client closed
							conn.Close()
							break
						}
						fmt.Println("Error to read message because of ", err)
						continue
					}
					if reqLen < 10 {
						// fmt.Println("too short")
						continue
					}
					data := buf[0:reqLen]
					perfix := fmt.Sprintf("%x", buf[0:2])
					matched := false
					for _, p := range perfixArr {
						matched = p == perfix
						if matched {
							break
						}
					}
					if !matched {
						// fmt.Printf("perfix: %s \t not matched\n", perfix)
						continue
					}
					// output message received
					go r.handler(clientID, perfix, data)
				}
			}()

		}
	}()

}

func (r *netReceiver) Clients() map[string]string {
	return r.ids
}

func (r *netReceiver) Write(id string, buf []byte) error {
	clientID, ok := r.ids[id]
	if !ok {
		clientID = id
	}
	conn, ok := r.clients[clientID]
	if !ok {
		return errors.New("clientID/id: " + clientID + " not exists")
	}
	_, err := conn.Write(buf)
	return err
}

func (r *netReceiver) SetID(id, connID string) (bool, error) {
	if _, ok := r.ids[id]; ok {
		return false, errors.New("id exists")
	}
	if _, ok := r.clients[connID]; !ok {
		return false, errors.New("conn id not exists")
	}
	r.ids[id] = connID
	return true, nil
}

//NewNetReceiver create a new receiver
func NewNetReceiver(options *Options, f ReceiveHandler) NetReceiver {
	return &netReceiver{
		options: options,
		handler: f,
		clients: make(map[string]net.Conn),
		ids:     make(map[string]string),
	}
}

func init() {
	fpm.Register(func(app *fpm.Fpm) {
		// 配置 socket
		if !app.HasConfig("socket") {
			panic("socket config node required")
		}
		socketConfig := app.GetConfig("socket").(map[string]interface{})
		options := &Options{
			Port:   int(socketConfig["port"].(float64)),
			Max:    int(socketConfig["max"].(float64)),
			Perfix: strings.Split(socketConfig["perfix"].(string), ","),
		}

		app.Logger.Debugf("Socket Config port: %+v", *options)

		server := NewNetReceiver(options, func(clientID string, perfix string, data []byte) {
			// publish here

			app.Publish("#tcp/receive/"+perfix, map[string]interface{}{
				"clientID": clientID,
				"data":     data,
			})
			app.Logger.Infof("tcp receive: %s, %v", clientID, data)
		})

		bizModule := make(fpm.BizModule, 0)

		bizModule["send"] = func(param *fpm.BizParam) (data interface{}, err error) {
			clientID := (*param)["clientID"].(string)
			// TODO: get type of the data, it can be string / []byte
			buf := (*param)["data"].([]byte)
			err = server.Write(clientID, buf)
			data = 1
			return
		}
		bizModule["clients"] = func(param *fpm.BizParam) (data interface{}, err error) {
			data = server.Clients()
			return
		}
		bizModule["setID"] = func(param *fpm.BizParam) (data interface{}, err error) {
			clientID := (*param)["clientID"].(string)
			id := (*param)["id"].(string)
			data, err = server.SetID(id, clientID)
			return
		}

		app.AddBizModule("socket", &bizModule)

		server.Listen()

	})
}

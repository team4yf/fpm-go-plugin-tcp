package main

import (
	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/fpm-go-plugin-tcp/plugin"
)

func main() {

	app := fpm.New()
	app.Init()
	// app.Execute("mqttclient.subscribe", &fpm.BizParam{
	// 	"topics": "$s2d/+/ipc/demo/execute",
	// })

	app.Subscribe("#tcp/receive/6162", func(topic string, data interface{}) {
		//data 通常是 byte[] 类型，可以转成 string 或者 map

		app.Logger.Debugf("data: %+v", data)
		payload := data.(map[string]interface{})
		clientID := payload["clientID"].(string)
		app.Execute("socket.setID", &fpm.BizParam{
			"clientID": clientID,
			"id":       "abc",
		})
		app.Execute("socket.send", &fpm.BizParam{
			"clientID": "abc",
			"data":     []byte{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97},
		})
	})
	app.Subscribe("#tcp/disconnect", func(_ string, data interface{}) {
		// data: { "id": "abc", "clientID": "bcd" }
		app.Logger.Debugf("data: %+v", data)
	})
	// app.Execute("mqttclient.publish", &fpm.BizParam{
	// 	"topic":   "$s2d/111/ipc/demo/feedback",
	// 	"payload": ([]byte)(`{"test":1}`),
	// })

	app.Run()

}

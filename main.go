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
	})
	// app.Execute("mqttclient.publish", &fpm.BizParam{
	// 	"topic":   "$s2d/111/ipc/demo/feedback",
	// 	"payload": ([]byte)(`{"test":1}`),
	// })

	app.Run(":9999")

}

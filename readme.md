## fpm-go-plugin-tcp

#### Config

```json
{
    "socket": {
        "port": 5002,
        "prefix": ["6162","6163"]
    }
}
```

```sh
echo -n -e "\xfe\xdc\x31\x38\x03\x34\x43\x0d\x0a\x0a\x0a\x0a\x0a\x0a\x0a\x0a\x0a" | nc localhost 5002
```
```sh
# subscribe the tcp data
mosquitto_sub -h open.yunplus.io -t '$d2s/tcp' -u "fpmuser" -P
```


ref:
https://shimo.im/docs/YhvpcDkCkPxgqrY9/ 《精讯云上传协议2.0+(2)》，可复制链接后用石墨文档 App 或小程序打开


```golang
// subscribe the event
app.Subscribe("#tcp/receive", func(topic string, data interface{}) {
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

app.Subscribe("#tcp/disconnect", func(_ string, data interface{} ) {
    // data: { "id": "abc", "clientID": "bcd" }
    app.Logger.Debugf("data: %+v", data)
})
```
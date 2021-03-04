package main

import (
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/Sin46/MqttServer"
)

const (
	GetDidKey string = "get_did_key"
)

func main() {
	rand.Seed(time.Now().Unix())
	client, err := mqtt.NewClient(mqtt.ClientInfo{
		ClientId:  fmt.Sprintf("go%04d", rand.Intn(9999)),
		BrokerUrl: "mqtt.fjdzzh.com:1883",
		Username:  "FjdzMacUser",
		Password:  "geomacuser",
		KeepAlive: 60,
	})
	if err != nil {
		fmt.Println("登录失败")
		return
	}
	client.Subscribe("+")
	for msg := range client.Out() {
		msg := NewMsg(msg)
		from, err := msg.From()
		if err != nil {
			continue
		}
		params := msg.Params()
		if params["cmd"] == GetDidKey {
			fmt.Printf("设备: %s 请求消息: %s\n", from, msg.base.Payload)
			sn := params["device_sn"]
			did := genDid(sn)
			reply := fmt.Sprintf("$cmd=set_did_key&device_sn=%s&did=%s&key=%036s", sn, did, did)
			client.In() <- mqtt.Msg{
				Topic:   "Subscription_" + from,
				Payload: reply,
			}
			fmt.Printf("设备: %s 授权消息: %s\n\n", from, reply)
		}

	}
}

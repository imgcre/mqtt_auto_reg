package main

import (
	"errors"
	"regexp"
	"strings"

	mqtt "github.com/Sin46/MqttServer"
)

type Msg struct {
	base   mqtt.Msg
	from   string
	params map[string]string
}

func NewMsg(base mqtt.Msg) *Msg {
	msg := &Msg{base, "", nil}
	return msg
}

func (msg *Msg) From() (string, error) {
	if msg.from != "" {
		return msg.from, nil
	}
	topicReg := regexp.MustCompile(`^Publish_(.*?)$`)
	result := topicReg.FindStringSubmatch(msg.base.Topic)
	if len(result) < 2 {
		return "", errors.New("unknown gateway")
	}
	msg.from = result[1]
	return result[1], nil
}

func (msg *Msg) Params() map[string]string {
	if msg.params != nil {
		return msg.params
	}
	params := make(map[string]string)
	payload := msg.base.Payload
	if strings.HasPrefix(payload, "$") {
		payload = payload[1:]
		payload = strings.ReplaceAll(payload, "\r", "")
		payload = strings.Split(payload, "\n")[0]
		for _, item := range strings.Split(payload, "&") {
			kvp := strings.Split(item, "=")
			key := kvp[0]
			value := kvp[1]
			params[key] = value
		}
	}
	msg.params = params
	return params
}

func genDid(sn string) string {
	m := map[string]string{
		"WG": "10",
		"YL": "11",
		"LF": "12",
		"QJ": "13",
	}
	return m[sn[0:2]] + sn[len(sn)-4:]
}

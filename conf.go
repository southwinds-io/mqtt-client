/*
  MQTT Client - © 2018-Present - SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package mqtt

import (
	"fmt"
	"os"
	"strconv"
)

type ConfKey string

const (
	ConfMQTTUri                 ConfKey = "MQTT_URI"
	ConfMQTTUser                ConfKey = "MQTT_USER"
	ConfMQTTPwd                 ConfKey = "MQTT_PWD"
	ConfMQTTInsecureSkipVerify  ConfKey = "MQTT_INSECURE_SKIP_VERIFY"
	ConfMQTTClientId            ConfKey = "MQTT_CLIENT_ID"
	ConfMQTTQoS                 ConfKey = "MQTT_QoS"
	ConfMQTTTopic               ConfKey = "MQTT_TOPIC"
	ConfMQTTShutdownGracePeriod ConfKey = "MQTT_SHUTDOWN_GRACE_PERIOD"
	ConfMQTTDebug                       = "MQTT_DEBUG"
)

type Conf struct {
}

func NewConf() *Conf {
	return &Conf{}
}

func (c *Conf) get(key ConfKey) string {
	return os.Getenv(string(key))
}

func (c *Conf) getConfOxMsgBrokerUri() string {
	return c.getValue(ConfMQTTUri)
}

func (c *Conf) getConfOxMsgBrokerClientId() string {
	return c.getValue(ConfMQTTClientId)
}

func (c *Conf) getConfOxMsgBrokerQoS() int {
	if len(c.getValue(ConfMQTTQoS)) > 0 {
		i, err := strconv.Atoi(c.getValue(ConfMQTTQoS))
		if err != nil {
			fmt.Printf("ERROR: invalid value for variable %s, value [%s] \n", ConfMQTTQoS, c.getValue(ConfMQTTQoS))
			os.Exit(1)
		}
		return i
	} else {
		return 0
	}
}

func (c *Conf) getConfOxMsgBrokerTopic() string {
	return c.getValue(ConfMQTTTopic)
}

func (c *Conf) getConfOxMsgBrokerUser() string {
	return c.getValue(ConfMQTTUser)
}

func (c *Conf) getConfOxMsgBrokerPwd() string {
	return c.getValue(ConfMQTTPwd)
}

func (c *Conf) getConfOxMsgBrokerInsecureSkipVerify() bool {
	b, err := strconv.ParseBool(c.getValue(ConfMQTTInsecureSkipVerify))
	if err != nil {
		fmt.Printf("ERROR: invalid value for variable %s \n", ConfMQTTInsecureSkipVerify)
		os.Exit(1)
	}
	return b
}

func (c *Conf) getConfOxMsgBrokerShutdownGracePeriod() uint {
	if len(c.getValue(ConfMQTTShutdownGracePeriod)) > 0 {
		i, err := strconv.ParseUint(c.getValue(ConfMQTTShutdownGracePeriod), 0, 32)
		if err != nil {
			fmt.Printf("ERROR: invalid value for variable %s, value [%s] \n", ConfMQTTShutdownGracePeriod, c.getValue(ConfMQTTShutdownGracePeriod))
			return 250
		}
		return uint(i)
	} else {
		return 250
	}
}

func (c *Conf) getValue(key ConfKey) string {
	value := os.Getenv(string(key))
	return value
}

func IsMqttConfigAvailable() bool {
	if len(os.Getenv(string(ConfMQTTUri))) > 0 {
		return true
	}
	return false
}

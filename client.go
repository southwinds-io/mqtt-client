/*
  MQTT Client - © 2018-Present - SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package mqtt

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	conf *Conf
	mq   mqtt.Client
}

var (
	mqc *Client
)

type connStatus struct {
	bool
	err error
}

// NewClient construct a new MsgClient with defaut configuration
func NewClient() *Client {
	var (
		err    error
		newmqc *Client
	)
	if newmqc == nil {
		newmqc, err = newMsgClient(new(Conf))
		if err != nil {
			log.Fatalf("ERROR: fail to create MsgClient : %s \n", err)
		}
		mqc = newmqc
	}
	return mqc
}

func buildClientOptions(c *Conf) *mqtt.ClientOptions {

	var cid string
	if len(c.getConfOxMsgBrokerClientId()) > 0 {
		cid = c.getConfOxMsgBrokerClientId()
	} else {
		cid = "runner-client"
	}
	clientId := flag.String("clientid", cid, "A client Id for the connection")
	username := flag.String("username", c.getConfOxMsgBrokerUser(), "A username to authenticate to the MQTT server")
	password := flag.String("password", c.getConfOxMsgBrokerPwd(), "Password to match username")
	server := flag.String("server", c.getConfOxMsgBrokerUri(), "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	flag.Parse()
	connOpts := mqtt.NewClientOptions().AddBroker(*server).SetClientID(*clientId).SetCleanSession(true)
	if len(c.getConfOxMsgBrokerUser()) > 0 {
		connOpts.SetUsername(*username)
		if len(c.getConfOxMsgBrokerUser()) > 0 {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: c.getConfOxMsgBrokerInsecureSkipVerify(), ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	return connOpts
}

// func newMsgClient(c *Conf, f OnMessageReceived) (*MsgClient, error) {
func newMsgClient(c *Conf) (*Client, error) {
	connOpts := buildClientOptions(c)
	debug("build client option for mqtt client ")
	client := mqtt.NewClient(connOpts)
	debug("new mqtt client created")
	return &Client{
		conf: c,
		mq:   client,
	}, nil
}

func (mqc *Client) Publish(topic, msg string) error {
	// TODO boolean retain or not to be fixed
	qos := mqc.conf.getConfOxMsgBrokerQoS()
	q := &qos
	if token := mqc.mq.Publish(topic, byte(*q), false, msg); token.Wait() && token.Error() != nil {
		debug("MQ client failed to public message to topic %s : %s\n", topic, token.Error())
		return token.Error()
	} else {
		debug("Published message to topic %s\n", topic)
		return nil
	}
}

func (mqc *Client) Subscribe(handler mqtt.MessageHandler) error {
	qos := mqc.conf.getConfOxMsgBrokerQoS()
	q := &qos
	topic := mqc.conf.getConfOxMsgBrokerTopic()
	if token := mqc.mq.Subscribe(topic, byte(*q), handler); token.Wait() && token.Error() != nil {
		debug("failed to subscribe topic %s : \n %s \n", topic, token.Error())
		return token.Error()
	}
	debug(" subscribed to topic %s \n", topic)
	return nil
}

func (mqc *Client) Start(waitInSeconds int) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// this channel receives a connection
	conn_status := make(chan connStatus, 1)
	// this channel receives a timeout flag
	timeout := make(chan connStatus, 1)

	go func() {
		if token := mqc.mq.Connect(); token.Wait() && token.Error() != nil {
			debug("MQ client failed to connect : %s", token.Error())
			conn_status <- connStatus{bool: false, err: token.Error()}
		} else {
			op := mqc.mq.OptionsReader()
			fmt.Printf("Connected to mqtt broker at [ %s ] with client id [ %s ]\n", mqc.conf.getConfOxMsgBrokerUri(), op.ClientID())
			conn_status <- connStatus{bool: true, err: nil}
		}
	}()

	go func() {
		// timeout period is In Seconds
		time.Sleep(time.Duration(waitInSeconds) * time.Second)
		er := errors.New("mqtt client failed to connect mqtt broker, the timed out period has elapsed\n")
		conn_status <- connStatus{bool: false, err: er}
	}()

	select {
	// the connection has been established before the timeout
	case c := <-conn_status:
		{
			return c.err
		}
	// the connection has not yet returned when the timeout happens
	case t := <-timeout:
		{
			return t.err
		}
	}

	<-stop
	debug("disconnecting mqtt client..")
	mqc.mq.Disconnect(mqc.conf.getConfOxMsgBrokerShutdownGracePeriod())
	debug("mqtt client disconnected..")
	return nil
}

func (mqc *Client) IsConnected() bool {
	return mqc.mq.IsConnected()
}

func debug(format string, args ...any) {
	if debugEnabled() {
		log.Printf(fmt.Sprintf("DEBUG %s", format), args)
	}
}

func debugEnabled() bool {
	return len(os.Getenv(ConfMQTTDebug)) > 0
}

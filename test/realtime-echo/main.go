package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xtaci/kcp-go"
)

const (
	// client num
	clientNum = 1000
	// io internal
	ioInterval = time.Second / 20

	serverAddr string = "127.0.0.1" // write your own server address here
	serverPort int    = 9999        // write your own server listening port here

	bufferLength = 800
)

type myClient struct {
	*kcp.UDPSession

	// clientID       uint16
	// userID         []byte
	receiveBuffer  []byte
	sendingMessage []byte
	// status         uint8

	interval time.Duration
}

var (
	payloadSize   = 64
	eventCodeEcho = 0x00f0

	messagePool = &sync.Pool{
		New: func() interface{} {
			var message = make([]byte, payloadSize+3)
			message[0] = 0x30 | 0x0f                 // header, hard coding
			message[1] = byte(payloadSize>>8) | 0xff // length high bits
			message[2] = byte(payloadSize | 0xff)    // length low bits
			message[3] = byte(eventCodeEcho >> 8)    // event code high bits
			message[4] = byte(eventCodeEcho | 0xff)  // event code low bits
			message[5] = 'h'
			message[6] = 'a'
			message[7] = 'y'
			message[8] = 'a'
			message[9] = 'b'
			message[10] = 'u'
			message[11] = 's'
			message[12] = 'a'
			// message[13] ...
			// ...
			// total length 64 will be sent
			return message
		},
	}
)

var (
	qpsCount int64 = 0
	qpsAt    int64
)

func simulateClient(client *myClient) {
	var kcpConn, err = kcp.DialWithOptions(fmt.Sprintf("%s:%d", serverAddr, serverPort), nil, 0, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	kcpConn.SetNoDelay(1, 10, 2, 1) // hard coding
	var (
		headerBytes = []byte{0x34}
		tokenBytes  = []byte("ZWNoby10ZXN0LWFwcCN0ZXN0LXVzZXIjdGVzdC1yYW5kb20tc3RyaW5n") // test user token
		lengthBytes = []byte{byte(len(tokenBytes) & 0xff)}

		heartbeatBytes = []byte{0x30, 0, 0, 0, 0, 0, 0, 0, 0}
	)
	var authMessage = append(headerBytes, lengthBytes...)
	authMessage = append(authMessage, tokenBytes...)
	kcpConn.Write(authMessage)

	client.UDPSession = kcpConn
	client.receiveBuffer = make([]byte, bufferLength)
	client.sendingMessage = messagePool.Get().([]byte)
	defer messagePool.Put(client.sendingMessage)

	// echo loop
	for {
		time.Sleep(client.interval)
		// when connection stablished, server will proactively send message
		n, err := client.UDPSession.Read(client.receiveBuffer)
		if n < 1 || err != nil {
			panic(nil)
		}
		// check received message
		if client.receiveBuffer[0] == 0x3f && !bytes.Equal(client.receiveBuffer[5:13], []byte("hayabusa")) {
			fmt.Printf("%s\n", client.receiveBuffer[5:13])
			panic("incorrected response")
		}
		// count
		atomic.AddInt64(&qpsCount, 1)
		var now = time.Now().Unix()
		var oldQpsAt = qpsAt
		qpsAt = now
		if now > oldQpsAt {
			fmt.Println("qps:", qpsCount)
			atomic.StoreInt64(&qpsCount, 0)
			client.UDPSession.Write(heartbeatBytes)
		}
		// sending message
		n, err = client.UDPSession.Write(client.sendingMessage)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		if n < 1 {
			panic(n)
		}
	}
}

func main() {
	fmt.Println("start test")
	qpsAt = time.Now().Unix()
	for i := 0; i < clientNum; i++ {
		var client = &myClient{
			interval: ioInterval,
		}
		go simulateClient(client)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}

package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/xtaci/kcp-go"
)

const (
	// room num
	roomNum = 10
	// room players num
	roomSize = 10
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
			return make([]byte, payloadSize+3)
		},
	}
)

var (
	receivedCount   int64 = 0
	receivedCountAt int64
)

func simulateClient(client *myClient) {
	var kcpConn, err = kcp.DialWithOptions(fmt.Sprintf("%s:%d", serverAddr, serverPort), nil, 0, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	kcpConn.SetNoDelay(1, 10, 2, 1) // hard coding
	var (
		authHeaderBytes = []byte{0x34}
		authTokenBytes  = []byte("ZWNoby10ZXN0LWFwcCN0ZXN0LXVzZXIjdGVzdC1yYW5kb20tc3RyaW5n") // test user token
		authLengthBytes = []byte{byte(len(authTokenBytes) & 0xff)}

		heartbeatBytes       = []byte{0x30, 0, 0, 0, 0, 0, 0, 0, 0}
		lastHeartAt    int64 = 0
	)
	var authMessage = append(authHeaderBytes, authLengthBytes...)
	authMessage = append(authMessage, authTokenBytes...)
	client.UDPSession = kcpConn
	client.receiveBuffer = make([]byte, bufferLength)
	client.sendingMessage = messagePool.Get().([]byte)
	defer messagePool.Put(client.sendingMessage)
	// basic authentication
	client.UDPSession.Write(authMessage)
	client.receiveResponseWithNonProcess()
	// auto enter room(builtin basic matching function)
	var pkt = newPacket()
	// header and event code
	pkt.SetHeader(0x38).SetEventCode(0x0110)
	// mux string
	pkt.WriteBytes([]byte("__mux_test__"))
	pkt.WriteFloat32(0).WriteFloat32(math.MaxFloat32)
	// room players num
	pkt.WriteUint16(roomSize).WriteUint16(roomSize)
	client.UDPSession.Write(pkt.ToBytes())
	client.receiveResponseWithNonProcess()
	// read loop
	go func() {
		var checkMessage = []byte("hello hayabusa")
		for {
			n, _ := client.UDPSession.Read(client.receiveBuffer)
			var message = client.receiveBuffer[14:n]
			// check message
			if client.receiveBuffer[0] == 0x28 && // response header
				client.receiveBuffer[3] == 0x80 && // event code high bits
				client.receiveBuffer[4] == 0xff && // event code low bits
				bytes.Equal(message, checkMessage) {
				receivedCount++
				var now = time.Now().Unix()
				if now > receivedCountAt {
					receivedCountAt = now
					fmt.Println("receiving:", receivedCount, " message(s)/s")
					receivedCount = 0
				}
			}
		}
	}()
	// write loop
	for {
		time.Sleep(ioInterval)
		pkt.Reset()
		pkt.SetHeader(0x38).SetEventCode(0x00ff)
		pkt.WriteUint8(2)  // 0=self 1=single user 2=room
		pkt.WriteUint16(0) // user id, omit when room broadcast
		pkt.WriteUint16(0) // client side callback id
		pkt.WriteString("hello hayabusa")
		client.UDPSession.Write(pkt.ToBytes())

		var now = time.Now().Unix()
		if now >= lastHeartAt+2 {
			lastHeartAt = now
			client.UDPSession.Write(heartbeatBytes)
		}
	}
}

func main() {
	fmt.Println("start test")
	receivedCountAt = time.Now().Unix()
	for i := 0; i < roomNum*roomSize; i++ {
		var client = &myClient{
			interval: ioInterval,
		}
		go simulateClient(client)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}

func (client *myClient) receiveResponseWithNonProcess() {
	client.UDPSession.Read(client.receiveBuffer)
}

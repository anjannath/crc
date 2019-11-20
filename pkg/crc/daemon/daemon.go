package daemon

import (
	"net"
	"os"
	"fmt"
	"bytes"
//	"strings"

	"github.com/code-ready/crc/pkg/crc/logging"
)

const (
	Start = "start"
	Stop = "stop"
	Status = "status"
	Delete = "delete"
	Version = "version"
	
	socketPath = "/tmp/crc.sock"
)

var (
	writeBuffer bytes.Buffer
)

func RunDaemon() {
	err := os.Remove(socketPath)
	if err != nil {
		logging.Error("Failed removing the old socket at", socketPath)
	}
	lnr, err := net.Listen("unix", socketPath)
	if err != nil {
		logging.Error("Failed to create daemon socket", socketPath)
	}
	conn, err := lnr.Accept()
	if err != nil {
		logging.Error("Error establishing connection with client")
	}
	// Check if conn is dropped and exit the loop
	for {
		logging.Info("Loop starts")
		readBuffer := make([]byte, 20)
		numBytes, err := conn.Read(readBuffer)
		if err != nil && numBytes == 0 {
			logging.Error("Error reading from socket")
		}
		command := string(readBuffer[0:numBytes])
		readBuffer = nil
		switch command {
		case Start:
			logging.Info("received start command")
			r, _ := startHandler()
			fmt.Fprint(&writeBuffer, r)
			conn.Write(writeBuffer.Bytes())
			writeBuffer.Reset()
		case Stop:
			logging.Info("handle Stop")
		case Delete:
			logging.Info("Handle delete")
		case Status:
			logging.Info("handle status")
		case Version:
			logging.Info("handle version")
		default:
			logging.Info("Unknown command received", command)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.bug.st/serial"
)

func main() {
	baud := flag.Int("baud", 19200, "Baud rate (9600 or 19200)")
	file := flag.String("file", "", "File to send (required)")
	port := flag.String("port", "", "Serial port, e.g. /dev/ttyUSB0 or COM3 (required)")
	txdelay := flag.Int("txdelay", 10, "Delay between bytes in milliseconds (0–1000)")
	echo := flag.Bool("echo", false, "Print each sent byte in hex to stdout")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `CL9 Core Universal Remote Firmware Loader

Sends a binary firmware file to a CL9 device over a serial port.
The data is XOR-encoded in 256-byte blocks before transmission.

Usage:
  cl9core --port <port> --file <file> [options]

Required:
  --port string    Serial port to use (e.g. /dev/ttyUSB0 or COM3)
  --file string    Firmware file to send

Options:
  --baud int       Baud rate, 9600 or 19200 (default 19200)
  --txdelay int    Delay between bytes in milliseconds, 0–1000 (default 10)
  --echo           Print each sent (XOR-encoded) byte in hex to stdout

Example:
  cl9core --port /dev/ttyUSB0 --file firmware.bin --baud 19200 --echo

`)
	}
	flag.Parse()

	if *file == "" || *port == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *baud != 9600 && *baud != 19200 {
		log.Fatal("--baud must be 9600 or 19200")
	}
	if *txdelay < 0 || *txdelay > 1000 {
		log.Fatal("--txdelay must be between 0 and 1000")
	}

	data, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
	}

	mode := &serial.Mode{
		BaudRate: *baud,
	}

	s, err := serial.Open(*port, mode)
	if err != nil {
		log.Fatalf("cannot open serial port %s: %v", *port, err)
	}
	defer s.Close()

	var prev byte = 0
	for i, b := range data {
		if i%256 == 0 {
			prev = 0
		}
		encoded := b ^ prev
		if *echo {
			if i%256 == 0 && i > 0 {
				fmt.Println()
			} else if i%256 != 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%02X", encoded)
		}
		if _, err := s.Write([]byte{encoded}); err != nil {
			log.Fatalf("write error at byte %d: %v", i, err)
		}
		prev = encoded
		if *txdelay > 0 {
			time.Sleep(time.Duration(*txdelay) * time.Millisecond)
		}
	}

	if *echo {
		fmt.Println()
	}
	fmt.Printf("Done. Sent %d bytes.\n", len(data))
}

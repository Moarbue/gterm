package main

import (
	"log"

	"github.com/albenik/go-serial/v2"
)

type serial_port struct {
	serialport *serial.Port
	serialopts serial_options
	isOpen     bool
}

type serial_options struct {
	baudrate     int
	databits     int
	HUPCL        bool
	parity       serial.Parity
	readtimeout  int
	stopbits     serial.StopBits
	writetimeout int
}

var serial_opts serial_options

func (sp *serial_port) get_ports() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return []string{""}
	} else {
		return ports
	}
}

func (sp *serial_port) start(portname string, baudrate int) {
	sp.serialopts.baudrate = baudrate
	sp.serialopts.databits = 8
	sp.serialopts.HUPCL = false
	sp.serialopts.parity = serial.NoParity
	sp.serialopts.readtimeout = 1000
	sp.serialopts.writetimeout = 1000
	sp.serialopts.stopbits = serial.OneStopBit
	sp.isOpen = false
	sp.change_port(portname)
	sp.apply_options()
	sp.isOpen = true
}

func (sp *serial_port) configure(baudrate int, databits int, HUPCL bool, parity serial.Parity, readtimeout int, writetimeout int, stopbits serial.StopBits) {
	sp.serialopts.baudrate = baudrate
	sp.serialopts.databits = databits
	sp.serialopts.HUPCL = HUPCL
	sp.serialopts.parity = parity
	sp.serialopts.readtimeout = readtimeout
	sp.serialopts.writetimeout = writetimeout
	sp.serialopts.stopbits = stopbits

	sp.apply_options()
}

func (sp *serial_port) change_port(portname string) {
	sp.isOpen = false
	sp.serialport.Close()

	port, err := serial.Open(portname)
	if err != nil {
		log.Fatal(err)
	}

	sp.serialport = port
	sp.apply_options()
	sp.isOpen = true
}

func (sp *serial_port) change_baudrate(baudrate int) {
	sp.serialopts.baudrate = baudrate
	sp.apply_options()
}

func (sp *serial_port) change_parity(parity serial.Parity) {
	sp.serialopts.parity = parity
	sp.apply_options()
}

func (sp *serial_port) apply_options() {
	sp.serialport.Reconfigure(
		serial.WithBaudrate(sp.serialopts.baudrate), serial.WithDataBits(sp.serialopts.databits), serial.WithHUPCL(sp.serialopts.HUPCL),
		serial.WithParity(sp.serialopts.parity), serial.WithReadTimeout(sp.serialopts.readtimeout), serial.WithStopBits(sp.serialopts.stopbits),
		serial.WithWriteTimeout(sp.serialopts.writetimeout))
}

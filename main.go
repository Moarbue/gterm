package main

import (
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/albenik/go-serial/v2"
)

var arduino serial_port

func get_baud() []string {
	return []string{"300", "1200", "2400", "4800", "9600", "19200", "38400", "57600", "74880", "115200", "230400", "250000", "500000"}
}
func get_parity() []string {
	return []string{"None", "Odd", "Even", "Mark", "Space"}
}

func main() {
	a := app.New()
	w := a.NewWindow("Serial Plotter")
	c := w.Canvas()

	content := make_interface()

	c.SetContent(content)
	w.ShowAndRun()
}

func change_baudrate(baud string) {
	baudrate, err := strconv.Atoi(baud)
	if err != nil {
		log.Fatal(err)
	}
	arduino.change_baudrate(baudrate)
}

func change_parity(parity string) {
	pars := get_parity()
	for i := 0; i < len(pars); i++ {
		if pars[i] == parity {
			arduino.change_parity(serial.Parity(i))
			break
		}
	}
}

func make_interface() *fyne.Container {
	portSel := &widget.Select{Options: arduino.get_ports(), PlaceHolder: "\t\t\t\t\t\t", OnChanged: arduino.change_port}
	baudSel := &widget.Select{Options: get_baud(), PlaceHolder: "9600  ", OnChanged: change_baudrate}
	paritySel := &widget.Select{Options: get_parity(), PlaceHolder: "None  ", OnChanged: change_parity}

	go func() {
		for range time.Tick(time.Second) {
			portSel.Options = arduino.get_ports()
		}
	}()

	return container.NewBorder(nil,
		container.NewPadded(container.NewHBox(layout.NewSpacer(), container.NewHBox(
			widget.NewLabel("Serial Port:"), portSel, layout.NewSpacer(),
			widget.NewLabel("Baud Rate:"), baudSel, layout.NewSpacer(),
			widget.NewLabel("Parity Mode:"), paritySel))), nil, nil)
}

package main

import (
	"image/color"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/albenik/go-serial/v2"
)

const WINDOW_WIDTH = 1080
const WINDOW_HEIGHT = 720

var arduino serial_port

func get_baud() []string {
	return []string{"300", "1200", "2400", "4800", "9600", "19200", "38400", "57600", "74880", "115200", "230400", "250000", "500000"}
}
func get_parity() []string {
	return []string{"None", "Odd", "Even", "Mark", "Space"}
}

func get_write_ends() []string {
	return []string{"None", "NL", "CR", "CR + NL"}
}

func main() {
	a := app.New()
	w := a.NewWindow("Serial Plotter")
	w.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT))
	w.SetFixedSize(true)
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

func serial_write(data string) {

}

func change_write_end(data string) {

}

func make_interface() *fyne.Container {
	portSel := &widget.Select{Options: arduino.get_ports(), PlaceHolder: "\t\t\t\t\t\t", OnChanged: arduino.change_port}
	baudSel := &widget.Select{Options: get_baud(), PlaceHolder: "9600  ", OnChanged: change_baudrate}
	paritySel := &widget.Select{Options: get_parity(), PlaceHolder: "None  ", OnChanged: change_parity}
	writeappendSel := &widget.Select{Options: get_write_ends(), Selected: get_write_ends()[0], OnChanged: change_write_end}
	writeBox := &widget.Entry{OnChanged: serial_write}
	readList := &widget.Entry{MultiLine: true}
	readList.Disable()

	go func() {
		for range time.Tick(time.Second) {
			portSel.Options = arduino.get_ports()
		}
	}()

	plot_bg := color.NRGBA{R: 30, G: 30, B: 30, A: 150}
	plot_fg := color.NRGBA{R: 180, G: 0, B: 0, A: 255}
	plot_area := canvas.NewRectangle(plot_bg)
	plot_area.SetMinSize(fyne.NewSize(WINDOW_WIDTH/2, WINDOW_HEIGHT/2))

	lines := []canvas.Line{
		{Position1: fyne.NewPos(0.1, 0.2), Position2: fyne.NewPos(100, 50), StrokeColor: plot_fg, StrokeWidth: 1, Hidden: false},
		{Position1: fyne.NewPos(0.1, 0.2), Position2: fyne.NewPos(0.1, 0.2), StrokeColor: plot_fg, StrokeWidth: 1, Hidden: false},
		{Position1: fyne.NewPos(0.1, 0.2), Position2: fyne.NewPos(0.1, 0.2), StrokeColor: plot_fg, StrokeWidth: 1, Hidden: false}}

	plot_container := container.NewPadded(plot_area, &lines[0], &lines[1], &lines[2])

	return container.NewBorder(container.NewPadded(container.NewVBox(
		container.NewHBox(widget.NewLabel("Send message:"), layout.NewSpacer(), widget.NewLabel("Send on enter"), writeappendSel), writeBox)),
		container.NewPadded(container.NewHBox(layout.NewSpacer(), container.NewHBox(
			widget.NewLabel("Serial Port:"), portSel, layout.NewSpacer(),
			widget.NewLabel("Baud Rate:"), baudSel, layout.NewSpacer(),
			widget.NewLabel("Parity Mode:"), paritySel))), nil, nil,
		container.NewPadded(container.NewHSplit(readList, plot_container)))
}

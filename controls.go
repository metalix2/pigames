
package main

import (

	
	"fmt"
	"log"

	"periph.io/x/conn/gpio"
	"periph.io/x/conn/gpio/gpioreg"
   	"periph.io/x/host"
    
)


func controls() {

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Lookup a pin by its number:
	p := gpioreg.ByName("GPIO27")
	if p == nil {
		log.Fatal("Failed to find GPIO27")
	}

	fmt.Printf("%s: %s\n", p, p.Function())

	// Set it as input, with an internal pull up resistor:
	if err := p.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
}

// button_L = DigitalInOut(board.D27)
// button_L.direction = Direction.INPUT
// button_L.pull = Pull.UP

// button_R = DigitalInOut(board.D23)
// button_R.direction = Direction.INPUT
// button_R.pull = Pull.UP

// button_U = DigitalInOut(board.D17)
// button_U.direction = Direction.INPUT
// button_U.pull = Pull.UP

// button_D = DigitalInOut(board.D22)
// button_D.direction = Direction.INPUT
// button_D.pull = Pull.UP

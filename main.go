package main

import (
	"embed"
    "image"
    "image/gif"
    "log"
    "time"

    "github.com/metalix2/pigames/scenarios"

    "periph.io/x/conn/v3/i2c/i2creg"
    "periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
    "periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)


//go:embed images
var image_path embed.FS


type Difficulty struct {
    Width   int
    Height  int
    Level   int
    Time    time.Duration
}

var difficulty = []Difficulty{{ Width:128, Height: 64, Level: 1, Time: 0.0}, { Width:128, Height: 128, Level: 2, Time: 0.0}, { Width: 384, Height: 64, Level: 3, Time: 0.0}, { Width:256, Height: 128, Level: 4, Time: 0.0}, { Width:256, Height: 192, Level: 5, Time: 0.0}, { Width: 128, Height: 64, Level: 0, Time: 0.0}}

var level int
var titleShown = false
var introShown = false
var showLevel = false
var introFrames = 0

var a_event = 0

func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
 
    // Open a handle to the first available I²C bus:
    bus, err := i2creg.Open("")

    // Open a handle to a ssd1306 connected on the I²C bus:
    dev, err := ssd1306.NewI2C(bus, &ssd1306.DefaultOpts)
    
    titlePath := "images/title.gif"
    f, err := image_path.Open(titlePath)
	if err != nil {
		log.Fatal(err)
	}
    titleGif, err := gif.DecodeAll(f)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    avatarpath := "images/avatar.gif"
    avatarReader, err := image_path.Open(avatarpath)
	if err != nil {
		log.Fatal(err)
	}
    avatarGif, err := gif.DecodeAll(avatarReader)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    matetopath := "images/mateto.gif"
    matetoReader, err := image_path.Open(matetopath)
	if err != nil {
		log.Fatal(err)
	}
    matetoGif, err := gif.DecodeAll(matetoReader)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }
    heartpath := "images/heart.gif"
    heartReader, err := image_path.Open(heartpath)
	if err != nil {
		log.Fatal(err)
	}
    heartGif, err := gif.DecodeAll(heartReader)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    // Lookup a pin by its number:
	p_r := gpioreg.ByName("GPIO23")
    p_l := gpioreg.ByName("GPIO27")
    p_u := gpioreg.ByName("GPIO17")
    p_d := gpioreg.ByName("GPIO22")
    p_c := gpioreg.ByName("GPIO4")
    p_a := gpioreg.ByName("GPIO6")
    // p_a := gpioreg.ByName("GPIO6")

	// Set it as input, with an internal pull up resistor:
	if err := p_r.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_l.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_u.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_d.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_c.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_a.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

    // Intial co-ordinates we need to capture previous co-ord and next co-ords:
    prev_coords := map[string]int{"x": 2, "y": 2}
    next_coords := map[string]int{"x": 2, "y": 2}

    // Global params
    var dir = 1
    var fps = 10
    level = 0
    screenX := 0
    screenY := 0
    ts := time.Now();
    p_counter := 0 // reset's when passes 5
    a_counter := 0
    index := 0
    // Display the frames in a loop:
    for i := 1; ;  {
        // fps currently 300ms per avatar frame so ~3FPS for avatar animation; screen refreshes every 100ms so 10FPS for each cycle. 
        c := time.After(time.Duration(10*fps) * time.Millisecond)
        if titleShown && introShown && difficulty[level].Level == 0 {
            // end game
            index = i % len(avatarGif.Image)
            // movement
            if p_r.Read() == gpio.Low {
                dir = 1
                next_coords["x"] += 2
            } else if p_u.Read() == gpio.Low {
                next_coords["y"] -= 2
            }
            if p_l.Read() == gpio.Low {
                dir = 0
                next_coords["x"] -= 2
            } else if p_d.Read() == gpio.Low {
                next_coords["y"] += 2
            }
            var img *image.Paletted
            var coords map[string]int
            img, coords = scenarios.DrawEnding(difficulty[level].Width, difficulty[level].Height, avatarGif.Image[index], matetoGif.Image[index],  heartGif.Image[index], prev_coords, next_coords, dir)
            prev_coords["x"] = coords["x"]
            prev_coords["y"] = coords["y"]
            next_coords["x"] = coords["x"]
            next_coords["y"] = coords["y"]
            dev.Draw(img.Bounds(), img, image.Point{0, 0})
            <-c
        } else if showLevel {
            introFrames += 1
            img := scenarios.DrawLevelText(128, 64, difficulty[level].Level, introFrames)
            dev.Draw(img.Bounds(), img, image.Point{0, 0})

            if p_a.Read() == gpio.Low {
                a_counter += 1
                if a_counter < 2{ 
                    showLevel = false
                }
            }

            if p_a.Read() == gpio.High {
                a_counter = 0            
            }

            <-c 
        } else if titleShown && introShown && difficulty[level].Level != 0 {
            index = i % len(avatarGif.Image)
            
            if ts.Add(time.Duration(10 * avatarGif.Delay[index]) * time.Millisecond).Sub(time.Now()) < time.Duration(10 * 1) * time.Millisecond {
                i++
                ts = time.Now();
                index = i % len(avatarGif.Image)
            }
            //reset mechanic
            if p_c.Read() == gpio.Low {
                p_counter += 1
            }
            if p_c.Read() != gpio.Low {
                p_counter = 0
            }
            if p_counter > 5  && p_counter < 10 {
                screenX = 0
                screenY = 0
                prev_coords = map[string]int{"x": 2, "y": 2}
                next_coords = map[string]int{"x": 2, "y": 2}
            } else if p_counter > 10 {
                level = 0
                // currentMaze = createMaze(int(math.Round(float64(difficulty[level].Width/14))), int(math.Round(float64(difficulty[level].Height/16))))
                screenX = 0
                screenY = 0
                prev_coords = map[string]int{"x": 2, "y": 2}
                next_coords = map[string]int{"x": 2, "y": 2}
            }
            // movement 
            if p_r.Read() == gpio.Low {
                dir = 1
                next_coords["x"] += 2
            } else if p_u.Read() == gpio.Low {
                next_coords["y"] -= 2
            }
            if p_l.Read() == gpio.Low {
                dir = 0
                next_coords["x"] -= 2
            } else if p_d.Read() == gpio.Low {
                next_coords["y"] += 2
            } 
            var img *image.Paletted
            var coords map[string]int
            var x, y int
    
            img, coords, x, y, level, sLevel, iFrames := scenarios.DrawCanvas(difficulty[level].Width, difficulty[level].Height, avatarGif.Image[index], prev_coords, next_coords, dir, screenX, screenY, difficulty[level].Width, difficulty[level].Height, difficulty[level].Level, showLevel, introFrames)
            showLevel = sLevel
            introFrames = iFrames
            // img, next_coords, screenX, screenY, 
            screenX = x
            screenY = y
            prev_coords["x"] = coords["x"]
            prev_coords["y"] = coords["y"]
            next_coords["x"] = coords["x"]
            next_coords["y"] = coords["y"]
            dev.Draw(img.Bounds(), img, image.Point{screenX, screenY})
            <-c
        } else {
            // We wait on title until button pressed
            if p_a.Read() == gpio.Low {
                a_counter += 1
                if a_counter < 2 { 
   
                    a_event += 1
                    introFrames = 0
                }
            }
            if p_a.Read() == gpio.High {
                a_counter = 0            
            }
            if i > len(titleGif.Image) - 2 {
                i = len(titleGif.Image) - 2
            }
            
            if ts.Add(time.Duration(10 * titleGif.Delay[i]) * time.Millisecond).Sub(time.Now()) < time.Duration(10 * 1) * time.Millisecond {
                i++
                ts = time.Now();
            }
            img, tShown, iShown, iFrames := scenarios.DrawIntro(128, 64, titleGif.Image[i], a_event, titleShown, introShown, introFrames)
            titleShown = tShown
            introShown = iShown
            introFrames = iFrames
            dev.Draw(img.Bounds(), img, image.Point{0, 0})
            <-c
        }
    }
}
package main

import (
	"embed"
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/gif"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
    "image/color/palette"
    "log"
    "math/rand"
    "math"
    "time"

    "pigames/imageflip"
    "github.com/itchyny/maze"

    "periph.io/x/conn/v3/i2c/i2creg"
    "periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
    "periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)


//go:embed images
var image_path embed.FS

var currentMaze *maze.Maze

type Config struct {
	Width       int
	Height      int
	Start       *maze.Point
	Goal        *maze.Point
	Image       bool
	Scale		int
	Seed        int64
}

func createMaze(w, h int) *maze.Maze {

    seed := rand.Int63n(10)
    config := &Config{
        Width: w,
        Height: h,
        Scale: 1,
        Start: &maze.Point{X: 0, Y: 0},
        Goal: &maze.Point{X: h-1 , Y: w-1}, //reverse but whatever 
        Seed: seed,
    }
    
    maze := maze.NewMaze(config.Height, config.Width)
    maze.Start = config.Start
    maze.Goal = config.Goal
    maze.Cursor = config.Start
	return maze
}

type Difficulty struct {
    Width   int
    Height  int
    Level   int
    Time    time.Duration

}

var difficulty = []Difficulty{{ Width:128, Height: 64, Level: 1, Time: 0.0}, { Width:128, Height: 128, Level: 2, Time: 0.0}, { Width: 384, Height: 64, Level: 3, Time: 0.0}, { Width:256, Height: 128, Level: 4, Time: 0.0}, { Width:256, Height: 192, Level: 5, Time: 0.0}, { Width: 128, Height: 64, Level: 0, Time: 0.0}}

var level int
var titleShown = true
var introShown = true
var showLevel = true
var introFrames = 0
var counter = 10
var a_event = 0
var dialog = [][]string{{}}

func drawLevelText(w, h int) (*image.Paletted) {
    
    img := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)

    dialog := [][]string{{fmt.Sprintf("      Level %d", difficulty[level].Level)}}
    introFrames += 1

    for di:=0; di < len(dialog); di++ {  
        for str:=0; str < len(dialog[di]); str++ {
            y := fixed.I(10+(di*35)+(str*11))
            if len(dialog[di]) == 1 {
                y = fixed.I(32)
            }
            d := &font.Drawer{
                Dst:  img,
                Src:  image.NewUniform(color.White),
                Face: basicfont.Face7x13,
                Dot:  fixed.Point26_6{fixed.I(1), y},
            }
            if introFrames < len(dialog[di][str]) {
                d.DrawString(dialog[di][str][:introFrames])
            } else  {
                d.DrawString(dialog[di][str])
            }
            
        }
    }
    return img
}

func drawIntro(w, h int, src image.Image) (*image.Paletted) {
    r := src.Bounds()
    img := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)

    if a_event == 1 {
        titleShown = true
        dialog = [][]string{{"Welcome to", "Pathfinder"}, {"Help Sabela escape", "the maze"}}
    }
    if a_event == 2 {
        dialog = [][]string{{"Get Sabela back to", "her Mateto"}}
    }
    if a_event == 3 {
        dialog = [][]string{{fmt.Sprintf("      Level %d", difficulty[level].Level)}}
    }
    if a_event == 4 {
        introShown = true
    }
    if a_event > 0  {
        introFrames += 1

        for di:=0; di < len(dialog); di++ {
            
            for  str:=0; str < len(dialog[di]); str++ {
                y := fixed.I(10+(di*35)+(str*11))
                if len(dialog[di]) == 1 {
                    y = fixed.I(32)
                }
                d := &font.Drawer{
                    Dst:  img,
                    Src:  image.NewUniform(color.White),
                    Face: basicfont.Face7x13,
                    Dot:  fixed.Point26_6{fixed.I(1), y},
                }
                if introFrames < len(dialog[di][str]) {
                    d.DrawString(dialog[di][str][:introFrames])
                } else  {
                    d.DrawString(dialog[di][str])
                }
            }
        }

        if introFrames > 25 {
            y := fixed.I(10+((len(dialog)-1)*35)+((len(dialog[len(dialog)-1])-1)*11))
            if len(dialog[len(dialog)-1]) == 1 {
                y = fixed.I(32)
            }
            g := &font.Drawer{
                Dst:  img,
                Src:  image.NewUniform(color.White),
                Face: basicfont.Face7x13,
                Dot:  fixed.Point26_6{fixed.I(1+(len(dialog[len(dialog)-1][len(dialog[len(dialog)-1])-1]))*7), y},
            }
            if counter >= 0 && counter < 5 {
                g.DrawString("_")
                counter -= 1
            } else if counter < 0 {
                counter += 10
            } else {
                counter -= 1
            } 
            
        }
    }
    if a_event < 1  {
        r = r.Add(image.Point{0, 0})
        draw.Draw(img, r, src, image.Point{0, 0}, draw.Src)
    }
    return img
}

func drawEnding(w, h int, src, src2, src3 image.Image, prev_coords map[string]int, next_coords map[string]int, dir int)(*image.Paletted, map[string]int) {
    // log.Println("ending")
    r1 := src.Bounds()
    log.Println(r1)
    r2 := src.Bounds()
    img := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)

    r1 = r1.Add(image.Point{prev_coords["x"], prev_coords["y"]})
    r2 = r2.Add(image.Point{60, 30})


    // Draw Avatar and it's Orientation
    if dir > 0 {
        draw.Draw(img, r1, imageflip.Flip(src), image.Point{1, 1}, draw.Src)
        draw.Draw(img, r2, imageflip.Flip(src2), image.Point{1, 1}, draw.Src) 
    } else {
        draw.Draw(img, r1, src, image.Point{0, 0}, draw.Src)
        draw.Draw(img, r2, src2, image.Point{0, 0}, draw.Src)
    }

    return img, next_coords
}

func drawCanvas(w, h int, src image.Image, prev_coords map[string]int, next_coords map[string]int, dir int, screenX int, screenY int) (*image.Paletted, map[string]int, int, int) {
    r := src.Bounds()
    img := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)

    currentMaze.Generate()
    drawMaze(currentMaze, img)
    
    // bound detection
    if 0 > prev_coords["x"] + r.Size().X  {
        if !inEnvironment(img, prev_coords, r) {
            level += 1
            prev_coords = map[string]int{"x": 2, "y": 2}
            next_coords = map[string]int{"x": 2, "y": 2}
            currentMaze = createMaze(int(math.Round(float64(difficulty[level].Width/14))), int(math.Round(float64(difficulty[level].Height/16))))
            showLevel = true
            introFrames = 0
            screenX = 0
            screenY = 0
        }
    }
    if prev_coords["x"] > r.Size().X && screenX + 0 > prev_coords["x"] + r.Size().X || 
    prev_coords["x"] + r.Size().X  >= 128 + screenX + r.Size().X  || 
    prev_coords["y"] > r.Size().Y && screenY + 0  > prev_coords["y"] + r.Size().Y  ||
    prev_coords["y"] + r.Size().Y  >= 64 + screenY + r.Size().Y {
        log.Println(prev_coords["x"] + r.Size().X  >= 128 + screenX + r.Size().X)
        if inEnvironment(img, prev_coords, r) {
            var vector [2]int
            vector[0] = (next_coords["x"] - prev_coords["x"])
            vector[1] = (next_coords["y"] - prev_coords["y"])

            if vector[0] > 0 {
                if screenX == 0 {
                    screenX += 126
                } else {
                    screenX += 128
                }
            }
            if vector[0] < 0 {
                if screenX == 126 {
                    screenX -= 126
                } else {
                    screenX -= 128
                }
                
            }
            if vector[1] > 0 {
                if screenY == 0 {
                    screenY += 63
                } else {
                    screenY += 64
                }
                
            }
            if vector[1] < 0 {
                if screenY == 63 {
                    screenY -= 63
                } else {
                    screenY -= 64
                }
            }
        } else {
            // Level Progression time Show New Level 
            level += 1
            prev_coords = map[string]int{"x": 2, "y": 2}
            next_coords = map[string]int{"x": 2, "y": 2}
            currentMaze = createMaze(int(math.Round(float64(difficulty[level].Width/14))), int(math.Round(float64(difficulty[level].Height/16))))
            showLevel = true
            introFrames = 0
            screenX = 0
            screenY = 0
        }

    }

    // Position Avatar
    r = r.Add(image.Point{prev_coords["x"], prev_coords["y"]})

    // check Avatar can't walk through walls
    if inteserction(img, next_coords, r) {
        next_coords["x"] = prev_coords["x"]
        next_coords["y"] = prev_coords["y"]
    }

    // Draw Avatar and it's Orientation
    if dir > 0 {
        draw.Draw(img, r, imageflip.Flip(src), image.Point{1, 1}, draw.Src) 
    } else {
        draw.Draw(img, r, src, image.Point{0, 0}, draw.Src)
    }

    // draw.Draw(palettedImage, palettedImage.Rect, simage, bounds.Min, draw.Over)
    return img, next_coords, screenX, screenY
}

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
    level = 5
    screenX := 0
    screenY := 0
    currentMaze = createMaze(int(math.Round(float64(difficulty[level].Width/14))), int(math.Round(float64(difficulty[level].Height/16))))
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
            img, coords = drawEnding(difficulty[level].Width, difficulty[level].Height, avatarGif.Image[index], matetoGif.Image[index],  heartGif.Image[index], prev_coords, next_coords, dir)
            prev_coords["x"] = coords["x"]
            prev_coords["y"] = coords["y"]
            next_coords["x"] = coords["x"]
            next_coords["y"] = coords["y"]
            dev.Draw(img.Bounds(), img, image.Point{0, 0})
            <-c
        } else if showLevel {
            img := drawLevelText(128, 64)
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
                currentMaze = createMaze(int(math.Round(float64(difficulty[level].Width/14))), int(math.Round(float64(difficulty[level].Height/16))))
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
    
            img, coords, x, y = drawCanvas(difficulty[level].Width, difficulty[level].Height, avatarGif.Image[index], prev_coords, next_coords, dir, screenX, screenY)
            
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
            img := drawIntro(128, 64, titleGif.Image[i])
            dev.Draw(img.Bounds(), img, image.Point{0, 0})
            <-c
        }
    }
}
package scenarios

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
    "image/color/palette"
    "log"
    "math"
	"math/rand"

    "github.com/metalix2/pigames/imageflip"
    "github.com/metalix2/pigames/environment"
    "github.com/itchyny/maze"
)

var currentMaze *maze.Maze
var counter = 10
var dialog = [][]string{{}}

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

func DrawLevelText(img *image.Paletted, level, introFrames int) () {
    
    dialog := [][]string{{fmt.Sprintf("      Level %d", level)}}

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
}

func DrawIntro(img *image.Paletted, src image.Image, a_event int, titleShown bool, introShown bool, introFrames int, showLevel bool) (bool, bool, int, bool) {
	r := src.Bounds()

	if a_event == 3 {
		introShown = true
		showLevel = true

		return titleShown, introShown, introFrames, showLevel
    }
    if a_event == 1 {
        titleShown = true
        dialog = [][]string{{"Welcome to", "Pathfinder"}, {"Help Sabela escape", "the maze"}}
    }
    if a_event == 2 {
        dialog = [][]string{{"Get Sabela back to", "her Mateto"}}
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
		// pathfinder title render
        r = r.Add(image.Point{0, 0})
        draw.Draw(img, r, src, image.Point{0, 0}, draw.Src)
    }
    return titleShown, introShown, introFrames, showLevel
}

func DrawEnding(w, h int, src, src2, src3 image.Image, prev_coords map[string]int, next_coords map[string]int, dir int)(*image.Paletted, map[string]int) {
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

func DrawCanvas(img *image.Paletted, src image.Image, prev_coords map[string]int, next_coords map[string]int, dir int, screenX int, screenY int, levelWidth int, levelHeight int, level int, showLevel bool, introFrames int) (map[string]int, int, int, int, bool, int) {
    r := src.Bounds()
	if currentMaze == nil {
		currentMaze = createMaze(int(math.Round(float64(levelWidth/14))), int(math.Round(float64(levelHeight/16))))
	}

    currentMaze.Generate()
    environment.DrawMaze(currentMaze, img)
    
    // bound detection
    if 0 > prev_coords["x"] + r.Size().X  {
        if !environment.InEnvironment(img, prev_coords, r) {
            level += 1
            prev_coords = map[string]int{"x": 2, "y": 2}
            next_coords = map[string]int{"x": 2, "y": 2}
            currentMaze = nil // redraw on next frame
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
        if environment.InEnvironment(img, prev_coords, r) {
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
            currentMaze = nil
            showLevel = true
            introFrames = 0
            screenX = 0
            screenY = 0
        }
    }

    // Position Avatar
    r = r.Add(image.Point{prev_coords["x"], prev_coords["y"]})

    // check Avatar can't walk through walls
    if environment.Inteserction(img, next_coords, r) {
        next_coords["x"] = prev_coords["x"]
        next_coords["y"] = prev_coords["y"]
    }

    // Draw Avatar and it's Orientation
    if dir > 0 {
        draw.Draw(img, r, imageflip.Flip(src), image.Point{1, 1}, draw.Src) 
    } else {
        draw.Draw(img, r, src, image.Point{0, 0}, draw.Src)
    }

    return next_coords, screenX, screenY, level, showLevel, introFrames
}

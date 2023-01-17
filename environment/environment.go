

package environment

import (
	"image"
    "image/color"
    "bytes"
    "io"
    "log"
    "strings"
    "github.com/itchyny/maze"
)

var white = color.RGBA{255,255,255,255} // white
// HLine draws a horizontal line
func HLine(x1, y, x2 int, canvas *image.Paletted) {
    for ; x1 <= x2; x1++ {
        canvas.Set(x1, y, white)
    }
}

// VLine draws a veritcal line
func VLine(x, y1, y2 int, canvas *image.Paletted) {
    for ; y1 <= y2; y1++ {
        canvas.Set(x, y1, white)
    }
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(x1, y1, x2, y2 int, canvas *image.Paletted) {
    HLine(x1, y1, x2, canvas)
    HLine(x1, y2, x2, canvas)
    VLine(x1, y1, y2, canvas)
    VLine(x2, y1, y2, canvas)
}

func DrawMaze(currentMaze *maze.Maze, canvas *image.Paletted) {
    var foo bytes.Buffer
    
    w := io.Writer(&foo)
    currentMaze.Print(w, maze.Default)
    // log.Println("")
    // log.Println(foo.String())
    // log.Println("")
    rows := strings.Split(foo.String(), "\n")
    log.Println(rows);
    for y, line := range rows {
        for x, c := range line {
            if string(c) == "#" {
                scaleX := 7
                scaleY := 8
                // look ahead to keep drawing
                if x + 1 < len(line) &&  string(line[x+1]) == "#" {
                    for  j := 0; j <= scaleX; j++ {
                        dy := y*scaleY 
                        if y != 0 { // override for last row
                            dy = dy-1
                        }
                        if y+1 == len(rows) { // override for last row
                            dy = dy-1
                        }
                        canvas.Set(x*scaleX+j, dy, white)
                    }
                }
                if y + 1 < len(rows) {
                    log.Println(rows[y+1])
                }
                if y + 1 < len(rows) && len(rows[y+1]) > 0 && string(rows[y+1][x]) == "#" {
                    for  j := 0; j < scaleY; j++ {
                        dy := y*scaleY+j
                        canvas.Set(x*scaleX, dy, white)
                    }
                }
            }
        }
    }
}



// does it intersec with environment?
func Inteserction(img *image.Paletted, next_coords map[string]int, r image.Rectangle)  bool {

    subRect := image.Rect(next_coords["x"], next_coords["y"], (r.Size().X  + next_coords["x"]) , (r.Size().Y + next_coords["y"]))
    gridSubImage := img.SubImage(subRect)
    // works for keeping inside boundaries
    for i := next_coords["x"]; i < next_coords["x"] + gridSubImage.Bounds().Size().X; i++ {
        // log.Println(len(grid[i]))
        for j := next_coords["y"]; j < next_coords["y"] + gridSubImage.Bounds().Size().Y; j++ {
            // log.Println(grid[i][j])
            if gridSubImage.At(i, j) == white {
                return true
            }
        }
    }
    return false
}

// is it still in the environment? when > width, height - tells to complete or render next section
func InEnvironment(img *image.Paletted, next_coords map[string]int, r image.Rectangle)  bool {
    log.Println(img.Bounds().Max.X);
    log.Println(next_coords["x"]);
    if (next_coords["x"] >= img.Bounds().Max.X){
        return false
    }
    if (next_coords["x"] + r.Size().X <= img.Bounds().Min.X){
        return false
    }
    if (next_coords["y"] >= img.Bounds().Max.Y){
        return false
    }
    if (next_coords["y"] + r.Size().Y <= img.Bounds().Min.Y){
        return false
    }
    
  
    // subRect := image.Rect(next_coords["x"], next_coords["y"], (r.Size().X  + next_coords["x"]) , (r.Size().Y + next_coords["y"]))
    // gridSubImage := img.SubImage(subRect)
    // works for keeping inside boundaries
   
    return true
}


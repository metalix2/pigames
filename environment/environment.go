

package environment

import (
    "image"
    "image/color"
    "bytes"
    "io"
    "strings"
    "github.com/itchyny/maze"
)

var white = color.RGBA{255,255,255,255} // white

func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}

func DrawMaze(currentMaze *maze.Maze, canvas *image.Paletted) {
    var foo bytes.Buffer

    w := io.Writer(&foo)
    currentMaze.Print(w, maze.Default)

    rows := strings.Split(foo.String(), "\n")

    // Maze Clean up
    rows = remove(rows, 0)
    rows = remove(rows, len(rows)-1)
    rows = remove(rows, len(rows)-1)
    for u := 0; u < len(rows);  u++ {
        rows[u] = strings.TrimSpace(rows[u])
        newRow := ""
        for v := 0; v < len(rows[u]);  v++ {
            if string(rows[u][v]) == "S" {
                rows[u] = "#" + rows[u][v+1:] // replace start with wall only works for lefthand starts
            }
            if v % 2 == 0 {
                newRow = newRow + string(rows[u][v])
            }
        }
        rows[u] = newRow
    }

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
                if y + 1 < len(rows) && string(rows[y+1][x]) == "#" {
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


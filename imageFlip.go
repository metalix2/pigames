package main


import (
    "image"
    "image/color/palette"
    "image/color"
)

func imageToGrid(img image.Image) [][]color.Color {
    // create and return a grid of pixels
    var grid  [][]color.Color
    size := img.Bounds().Size()

    for i := 0; i < size.X; i++ {
        var y []color.Color
        for j := 0; j < size.Y; j++ {
            y = append(y, img.At(i, j))
        }
        grid = append(grid, y)
    }
    return grid
}

func flip(img image.Image) image.Image {

    // create and return a grid of pixels
    var grid  = imageToGrid(img)
    size := img.Bounds().Size()
    for i := 0; i < size.X; i++ {
        var y []color.Color
        for j := 0; j < size.Y; j++ {
            y = append(y, img.At(i, j))
        }
        grid = append(grid, y)
    }

    // Flips Verticaly
    lenX := len(grid)-1
    for x := 0; x < len(grid) / 2; x++ {
        col := grid[x]
        for  y := 0; y < len(col); y++ {
            var t = grid[x][y]
            grid[x][y] = grid[lenX-x][y]
            grid[lenX-x][y] = t
        }
    }
  
    xlen, ylen := len(grid), len(grid[0])
    rect := image.Rect(0, 0, xlen, ylen)
    img2 := image.NewPaletted(rect, palette.Plan9)
    for x := 0; x < xlen; x++ {
        for y := 0; y < ylen; y++ {
            img2.Set(x, y, grid[x][y])
        }
    }
    return img2.SubImage(rect)
}

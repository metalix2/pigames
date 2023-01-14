package imageFlip

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
            // gray.Set(i, j, img.At(i, j))
            r1, g1, b1, _ := img.At(i, j).RGBA()
            r2, g2, b2, _ := color.Black.RGBA()
            if (r1 == r2 && g1 == g2 && b1 == b2) {
                y = append(y, color.Black)
            } else {
                y = append(y, color.White)
            }
        }
        grid = append(grid, y)
    }
    return grid
}

func flip(img image.Image) image.Image {
    // create and return a grid of pixels
    var grid = imageToGrid(img)
    xlen, ylen := len(grid), len(grid[0])

    // reverse array
    for i, j := 0, len(grid)-1; i < j; i, j = i+1, j-1 {
        grid[i], grid[j] = grid[j], grid[i]
    }

    rect := image.Rect(0, 0, xlen+1, ylen+1)
    img2 := image.NewPaletted(image.Rect(0, 0, xlen+1, ylen+1), palette.Plan9)
    for x := 0; x < xlen; x++ {
        for y := 0; y < ylen; y++ {
            img2.Set(x+1, y+1, grid[x][y])
        }
    }
    return img2.SubImage(rect)
}

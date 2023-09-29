package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "Juego2/scenes"
)

func main() {
    pixelgl.Run(runGame)
}

func runGame() {
    cfg := pixelgl.WindowConfig{
        Title:  "Nave Espacial",
        Bounds: pixel.R(0, 0, 1000, 800),
        VSync:  true,
    }

    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }
    
    scenes.Setup(win)
}

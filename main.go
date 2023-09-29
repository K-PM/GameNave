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
    // Configurar la ventana
    cfg := pixelgl.WindowConfig{
        Title:  "Nave Espacial",
        Bounds: pixel.R(0, 0, 1000, 800),
        VSync:  true,
    }

    // Crear la ventana
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }
    

    // Ejecutar el juego en scenes.Setup
    scenes.Setup(win)
}

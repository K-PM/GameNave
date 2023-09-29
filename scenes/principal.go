package scenes

import (
	"time"
    utils "Juego2/utils"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	last                 = time.Now()
)

func Setup(win *pixelgl.Window) {
    fondoSprite, naveSprite, disparoSprite, enemySprite := utils.CargarImagenesYSprites()
    vida, nave := utils.CrearVidaYNave(naveSprite, disparoSprite)
    tiempoTranscurrido := 0.0
    last := time.Now() 
    tiempoTranscurridoCh := make(chan float64)

    go utils.ReproducirMusica()
    go utils.GestionarEnemigos(win, nave, enemySprite)
    go utils.GestionarTiempo(win, tiempoTranscurrido, vida)

    for !win.Closed() {
        fondoSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
        dt := time.Since(last).Seconds()
        last = time.Now()
        tiempoTranscurrido += dt
        tiempoTranscurridoCh <- tiempoTranscurrido

        utils.ActualizarVida(win, enemies, vida)
        utils.ActualizarYDibujarNave(nave, dt, win, tiempoTranscurrido, vida)

        win.Update()
        win.Clear(pixel.RGB(0, 0, 0))
    }
}
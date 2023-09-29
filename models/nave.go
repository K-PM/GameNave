package models

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Nave struct {
	sprite        *pixel.Sprite
	pos           pixel.Vec
	vel           pixel.Vec
	velocidad     float64
	ShootInterval float64
	LastShotTime  time.Time
	Disparos      []*Disparo
	DisparoSprite *pixel.Sprite
}

type Disparo struct {
	Pos    pixel.Vec
	Vel    pixel.Vec
	Width  float64 // Ancho del disparo
	Height float64 // Alto del disparo
}


func NewNave(sprite, disparoSprite *pixel.Sprite, velocidad, shootInterval float64) *Nave {
	return &Nave{
		sprite:        sprite,
		pos:           pixel.ZV,
		vel:           pixel.ZV,
		velocidad:     velocidad,
		ShootInterval: shootInterval,
		LastShotTime:  time.Now(),
		DisparoSprite: disparoSprite,
	}
}

func (n *Nave) Update(dt float64, win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyLeft) {
		n.vel.X = -n.velocidad
	} else if win.Pressed(pixelgl.KeyRight) {
		n.vel.X = n.velocidad
	} else {
		n.vel.X = 0
	}
	// Ajustar la posición vertical para que esté 4 cm por encima del margen inferior
	n.pos = n.pos.Add(n.vel.Scaled(dt))
	if n.pos.Y < n.sprite.Frame().H()/2+30 { // 40 píxeles representan 4 cm
		n.pos.Y = n.sprite.Frame().H()/2 + 30
	}
	// Limitar la posición de la nave horizontalmente para que siempre esté dentro de la ventana
	if n.pos.X < n.sprite.Frame().W()/2 {
		n.pos.X = n.sprite.Frame().W() / 2
	}
	if n.pos.X > win.Bounds().W()-n.sprite.Frame().W()/2 {
		n.pos.X = win.Bounds().W() - n.sprite.Frame().W()/2
	}
}

func (n *Nave) Draw(win *pixelgl.Window) {
	n.sprite.Draw(win, pixel.IM.Moved(n.pos))

	// Dibujar los disparos
	for _, disparo := range n.Disparos {
		n.DisparoSprite.Draw(win, pixel.IM.Moved(disparo.Pos))
	}
}

func (n *Nave) Disparar() {
	// Verificar si ha pasado suficiente tiempo desde el último disparo
	if time.Since(n.LastShotTime).Seconds() >= n.ShootInterval {
		//fmt.Println("¡Disparo!")
		// Crear un nuevo disparo y agregarlo a la lista de disparos
		n.Disparos = append(n.Disparos, &Disparo{
			Pos:    n.pos,
			Vel:    pixel.V(0, 500), // Aumenta la velocidad vertical de los disparos según tus necesidades
			Width:  10.0, // Ancho de tu disparo en píxeles
			Height: 10.0 ,
		})
		
		n.LastShotTime = time.Now()
	}
}

func (n *Nave) ActualizarDisparos(disparos []*Disparo, dt float64, win *pixelgl.Window) []*Disparo {
	// Actualizar la posición de los disparos y eliminar los que salen de la pantalla
	actualizados := make([]*Disparo, 0)
	for _, disparo := range disparos {
		disparo.Pos = disparo.Pos.Add(disparo.Vel.Scaled(dt))

		// Si el disparo está dentro de la pantalla, conservarlo
		if disparo.Pos.Y < win.Bounds().H() {
			actualizados = append(actualizados, disparo)
		}
	}

	return actualizados
}


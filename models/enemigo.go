package models

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "math/rand"
    "time"
)

const(
    objectSize  = 60.0
    minEnemyX              = 100.0
    maxEnemyX              = 900.0 
)
var (
    randSrc             = rand.NewSource(time.Now().UnixNano())
    randGen             = rand.New(randSrc)
)

// Define la estructura para representar a los objetos enemigos
type Enemigo struct {
    sprite *pixel.Sprite
    pos    pixel.Vec
    hit    bool
    speed  float64 
}

// Inicializa y devuelve un nuevo enemigo
func NuevoEnemigo(sprite *pixel.Sprite, initialPosition pixel.Vec) *Enemigo {
    return &Enemigo{
        sprite: sprite,
        pos:    initialPosition,
        hit:    false,
        speed:  100, // Establece la velocidad del enemigo
    }
}

// Función para dibujar el enemigo en la ventana
func (e *Enemigo) Pintar(win *pixelgl.Window) {
    if !e.hit {
        e.sprite.Draw(win, pixel.IM.Moved(e.pos))
    }
}

// Función para actualizar la posición del enemigo en función del tiempo
func (e *Enemigo) Actualizar(dt float64) {
    // Mueve el enemigo hacia abajo
    e.pos.Y -= e.speed * dt
}

// Función para obtener la posición actual del enemigo
func (e *Enemigo) Pos() pixel.Vec {
    return e.pos
}


func (e *Enemigo) Golpeado(pos pixel.Vec) bool {
    if !e.hit && pos.X >= e.pos.X && pos.X <= e.pos.X+objectSize && pos.Y >= e.pos.Y && pos.Y <= e.pos.Y+objectSize {
        e.hit = true
        return true
    }
    return false
}

// Agrega el método Reset para reiniciar la posición del enemigo
func (e *Enemigo) Resetear() {
	// Define aquí la posición inicial del enemigo cuando toque el borde inferior
	e.pos = pixel.V(minEnemyX+randGen.Float64()*(maxEnemyX-minEnemyX), 700) // Coordenada X aleatoria
	e.hit = false
}

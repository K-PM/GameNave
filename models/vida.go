// vida.go

package models

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Vida representa la cantidad de vidas del jugador
type Vida struct {
	sprite       *pixel.Sprite
	sinVidaSprite *pixel.Sprite
	pos          pixel.Vec
	vidas        int
	vidasPerdidas int // Número de vidas perdidas
}

// NewVida crea una nueva instancia de Vida con 5 vidas iniciales
func NewVida() (*Vida, error) {
	// Cargar la imagen de vida desde un archivo ("vida.png")
	vidaImg, err := loadPicture("./assets/vida.png")
	if err != nil {
		return nil, err
	}

	// Cargar la imagen de sin vida desde un archivo ("sinVida.png")
	sinVidaImg, err := loadPicture("./assets/sinVida.png")
	if err != nil {
		return nil, err
	}

	// Crear un sprite para la vida
	vidaSprite := pixel.NewSprite(vidaImg, vidaImg.Bounds())

	return &Vida{
		sprite:       vidaSprite,
		sinVidaSprite: pixel.NewSprite(sinVidaImg, sinVidaImg.Bounds()),
		pos:          pixel.V(50, 750),
		vidas:        5, // Inicialmente, el jugador tiene 5 vidas
		vidasPerdidas: 0, // Inicialmente, no se han perdido vidas
	}, nil
}

// Draw dibuja las vidas en la ventana
func (v *Vida) Draw(win *pixelgl.Window) {
	for i := 0; i < v.vidas; i++ {
		var vidaToDraw *pixel.Sprite
		if i < v.vidasPerdidas {
			vidaToDraw = v.sinVidaSprite // Usar sprite de sin vida para vidas perdidas
		} else {
			vidaToDraw = v.sprite // Usar sprite de vida normal para vidas restantes
		}

		mat := pixel.IM.Moved(v.pos.Add(pixel.V(float64(i)*40, 0)))
		vidaToDraw.Draw(win, mat)
	}
}

// PerderVida reduce la cantidad de vidas en 1 y registra que se ha perdido una vida
func (v *Vida) PerderVida() {
	if v.vidasPerdidas < v.vidas {
		v.vidasPerdidas++
	}
}

// ObtenerVidas devuelve la cantidad actual de vidas
func (v *Vida) ObtenerVidas() int {
	return v.vidas - v.vidasPerdidas
}

// loadPicture carga una imagen desde un archivo
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

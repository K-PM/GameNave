package utils

import (
	"Juego2/models"
	enemigo "Juego2/models"
	nave "Juego2/models"
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"sync" // goroutines- hilos.
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

var (
	lastEnemyTime        = time.Now()
	navesEliminadas      int
	enemies              []*enemigo.Enemigo
	mutex                sync.Mutex
	format               beep.Format
	streamer             beep.StreamSeeker
	ctrl                 *beep.Ctrl
)

func ActualizarYDibujarNave(nave *nave.Nave, dt float64, win *pixelgl.Window, tiempoTranscurrido float64, vida *models.Vida) {
	if win.Pressed(pixelgl.KeySpace) {
		nave.Disparar()
	}
	nave.Disparos = nave.ActualizarDisparos(nave.Disparos, dt, win)
	nave.Actualizar(dt, win)
	nave.Pintar(win)
}

func ActualizarVida(win *pixelgl.Window, enemies []*enemigo.Enemigo, vida *models.Vida) {
	vida.Pintar(win)
	for _, enemy := range enemies {
		if enemy.Pos().Y < 0 {
			vida.PerderVida()
			enemy.Resetear()
		}
	}
}
func CrearVidaYNave(naveSprite, disparoSprite *pixel.Sprite) (*models.Vida, *nave.Nave) {
	vida, err := models.NuevaVida()
	if err != nil {
		return nil, nil
	}

	nave := nave.NuevaNave(naveSprite, disparoSprite, 200.0, 0.2)

	return vida, nave
}
func CargarImagenesYSprites() (*pixel.Sprite, *pixel.Sprite, *pixel.Sprite, *pixel.Sprite) {
	// Cargar la imagen de fondo desde un archivo ("gal.png")
	fondoImg, err := LoadPicture("./assets/fon.png")
	if err != nil {
		return nil, nil, nil, nil
	}
	fondoSprite := pixel.NewSprite(fondoImg, fondoImg.Bounds())

	// Cargar la imagen de la nave desde un archivo (asegúrate de tener "nave.png" en la misma carpeta)
	naveImg, err := LoadPicture("./assets/nave.png")
	if err != nil {
		return nil, nil, nil, nil
	}

	disparoImg, err := LoadPicture("./assets/disparo.png")
	if err != nil {
		return nil, nil, nil, nil
	}

	// Crear un sprite para la nave
	naveSprite := pixel.NewSprite(naveImg, naveImg.Bounds())

	// Crear un sprite para el disparo
	disparoSprite := pixel.NewSprite(disparoImg, disparoImg.Bounds())

	// Cargar la imagen del enemigo desde un archivo ("enemigo.png")
	enemyImg, err := LoadPicture("./assets/enemigo.png")
	if err != nil {
		return nil, nil, nil, nil
	}

	// Crear un sprite para el enemigo
	enemySprite := pixel.NewSprite(enemyImg, enemyImg.Bounds())

	return fondoSprite, naveSprite, disparoSprite, enemySprite
}
func LoadPicture(path string) (pixel.Picture, error) {
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
func GestionarEnemigos(win *pixelgl.Window, nave *nave.Nave, enemySprite *pixel.Sprite) {
	last := time.Now()

	for {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if time.Since(lastEnemyTime).Seconds() >= 1 {
			initialEnemyPosition := pixel.V(100.0+rand.New(rand.NewSource(time.Now().UnixNano())).Float64()*(900.0-100.0), 700)
			enemy := enemigo.NuevoEnemigo(enemySprite, initialEnemyPosition)

			mutex.Lock()
			enemies = append(enemies, enemy)
			mutex.Unlock()

			lastEnemyTime = time.Now()
		}

		mutex.Lock()
		updatedEnemies := []*enemigo.Enemigo{}
		for _, enemy := range enemies {
			enemy.Actualizar(dt)
			if enemy.Pos().Y > -10000 {
				enemy.Pintar(win)
				updatedEnemies = append(updatedEnemies, enemy)
			}
		}

		updatedShots := []*models.Disparo{}
		for _, shot := range nave.Disparos {
			collided := false
			for _, enemy := range updatedEnemies {
				if enemy.Golpeado(shot.Pos) {
					collided = true
					navesEliminadas++
					break
				}
			}
			if !collided {
				updatedShots = append(updatedShots, shot)
			}
		}

		nave.Disparos = updatedShots
		enemies = updatedEnemies

		mutex.Unlock()
	}
}
func GestionarTiempo(win *pixelgl.Window, tiempoTranscurridoCh <-chan float64, vida *models.Vida) {
	for {
		tiempoTranscurrido := <-tiempoTranscurridoCh
		// Dibujar el temporizador en la esquina superior derecha
		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		timerMsg := text.New(pixel.V(830, 779), atlas)
		_, _ = timerMsg.WriteString(fmt.Sprintf("Tiempo: %.1f s", tiempoTranscurrido))
		timerMsg.Draw(win, pixel.IM.Scaled(timerMsg.Orig, 1.5))
		// Verificar si se alcanzaron los 50 segundos y mostrar "Game Over"
		if tiempoTranscurrido >= 50.0 {
			gameoverMsg := text.New(pixel.V(250, 450), atlas)
			_, _ = gameoverMsg.WriteString("Game Over")
			gameoverMsg.Draw(win, pixel.IM.Scaled(gameoverMsg.Orig, 8))
			win.Update()
			time.Sleep(2 * time.Second) // Espera 2 segundos antes de salir
			return
		}

		// Dibujar el contador de naves enemigas eliminadas
		contadorMsg := text.New(pixel.V(45, 778), atlas)
		_, _ = contadorMsg.WriteString(fmt.Sprintf("Naves Eliminadas: %d", navesEliminadas))
		contadorMsg.Draw(win, pixel.IM.Scaled(contadorMsg.Orig, 1.5))

		// Verificar si el jugador ha ganado
		if navesEliminadas >= 10 {
			ganasteMsg := text.New(pixel.V(250, 450), atlas)
			_, _ = ganasteMsg.WriteString("Ganaste")
			ganasteMsg.Draw(win, pixel.IM.Scaled(ganasteMsg.Orig, 8))
			win.Update()
			time.Sleep(2 * time.Second) // Espera 2 segundos antes de salir
			return
		}

		// Verificar si el jugador se quedó sin vidas
		if vida.ObtenerVidas() <= 0 {
			gameoverMsg := text.New(pixel.V(250, 450), atlas)
			_, _ = gameoverMsg.WriteString("Game Over")
			gameoverMsg.Draw(win, pixel.IM.Scaled(gameoverMsg.Orig, 8))
			win.Update()
			time.Sleep(2 * time.Second) // Espera 2 segundos antes de salir
			return
		}
	}
}
func ReproducirMusica() {
	// Cargar el archivo de sonido de fondo
	musicFile, err := os.Open("./assets/music.mp3")
	if err != nil {
		fmt.Println("Error al abrir el archivo de música:", err)
		return
	}
	streamer, format, err = mp3.Decode(musicFile)
	if err != nil {
		fmt.Println("Error al decodificar el archivo de música:", err)
		return
	}

	// Configurar el speaker para el formato del sonido
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Println("Error al inicializar el speaker:", err)
		return
	}

	// Reproducir el sonido en bucle
	ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	speaker.Play(ctrl)
}

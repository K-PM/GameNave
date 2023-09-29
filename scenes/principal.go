 package scenes

import (
    "Juego2/models"
    enemigo "Juego2/models"
    nave "Juego2/models"
    "fmt"
    "image"
    _ "image/png"
    "os"
    "time"
    "math/rand"
	
    "github.com/faiface/pixel/pixelgl"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/text"
    "golang.org/x/image/font/basicfont"
	"sync" // goroutines- hilos.
)

var (
    lastEnemyTime       = time.Now()
    navesEliminadas     int
    enemies             []*enemigo.Enemigo
    last                = time.Now()
    mutex               sync.Mutex
    tiempoTranscurridoCh = make(chan float64)
    vidaMutex sync.Mutex
)

func Setup(win *pixelgl.Window) {
	fondoSprite, naveSprite, disparoSprite, enemySprite := cargarImagenesYSprites()
	vida, nave := crearVidaYNave(naveSprite, disparoSprite)
	tiempoTranscurrido := 0.0

	// Agregar goroutines concurrentes para gestionar enemigos, disparos y tiempo.
	go gestionarEnemigos(win, nave, enemySprite)
	go gestionarTiempo(win, tiempoTranscurrido, vida)
    go ActualizarVida(enemies, vida)

	for !win.Closed() {
		fondoSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		dt := time.Since(last).Seconds()
		last = time.Now()
		tiempoTranscurrido += dt
        tiempoTranscurridoCh <- tiempoTranscurrido

		ActualizarVida( enemies, vida)
		actualizarYDibujarNave(nave, dt, win, tiempoTranscurrido, vida)

		win.Update()
		win.Clear(pixel.RGB(0, 0, 0))
	}
}


func actualizarYDibujarNave(nave *nave.Nave, dt float64, win *pixelgl.Window, tiempoTranscurrido float64, vida *models.Vida) {
	if win.Pressed(pixelgl.KeySpace) {
        nave.Disparar()
    }
    nave.Disparos = nave.ActualizarDisparos(nave.Disparos, dt, win)
	nave.Update(dt, win)
	nave.Draw(win)
    vida.Draw(win)
}


func crearVidaYNave(naveSprite, disparoSprite *pixel.Sprite) (*models.Vida, *nave.Nave) {
    // Crear las vidas del jugador
    vida, err := models.NewVida()
    if err != nil {
        return nil, nil
    }

    // Crear la nave
    nave := nave.NewNave(naveSprite, disparoSprite, 200.0, 0.2)

    return vida, nave
}

func cargarImagenesYSprites() (*pixel.Sprite, *pixel.Sprite, *pixel.Sprite, *pixel.Sprite) {
    // Cargar la imagen de fondo desde un archivo ("gal.png")
    fondoImg, err := loadPicture("./assets/fon.png")
    if err != nil {
        return nil, nil, nil, nil
    }
    // Crear un sprite para el fondo
    fondoSprite := pixel.NewSprite(fondoImg, fondoImg.Bounds())

    // Cargar la imagen de la nave desde un archivo (asegúrate de tener "nave.png" en la misma carpeta)
    naveImg, err := loadPicture("./assets/nave.png")
    if err != nil {
        return nil, nil, nil, nil
    }

    // Cargar la imagen del disparo (asegúrate de tener "disparo.png" en la misma carpeta)
    disparoImg, err := loadPicture("./assets/disparo.png")
    if err != nil {
        return nil, nil, nil, nil
    }

    // Crear un sprite para la nave
    naveSprite := pixel.NewSprite(naveImg, naveImg.Bounds())

    // Crear un sprite para el disparo
    disparoSprite := pixel.NewSprite(disparoImg, disparoImg.Bounds())

    // Cargar la imagen del enemigo desde un archivo ("enemigo.png")
    enemyImg, err := loadPicture("./assets/enemigo.png")
    if err != nil {
        return nil, nil, nil, nil
    }

    // Crear un sprite para el enemigo
    enemySprite := pixel.NewSprite(enemyImg, enemyImg.Bounds())

    return fondoSprite, naveSprite, disparoSprite, enemySprite
}

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

//CONCURRENCIA
func gestionarEnemigos(win *pixelgl.Window, nave *nave.Nave, enemySprite *pixel.Sprite) {
    last := time.Now()

    for {
        dt := time.Since(last).Seconds()
        last = time.Now()

        if time.Since(lastEnemyTime).Seconds() >= 1 {
            initialEnemyPosition := pixel.V(100.0+rand.New(rand.NewSource(time.Now().UnixNano())).Float64()*(900.0-100.0), 700)
            enemy := enemigo.NewEnemy(enemySprite, initialEnemyPosition)

            mutex.Lock()
            enemies = append(enemies, enemy)
            mutex.Unlock()

            lastEnemyTime = time.Now()
        }

        mutex.Lock()
        updatedEnemies := []*enemigo.Enemigo{}
        for _, enemy := range enemies {
            enemy.Update(dt)
            if enemy.Pos().Y > -10000 {
                enemy.Draw(win)
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

func ActualizarVida(enemies []*enemigo.Enemigo, vida *models.Vida) {
    for {
        for _, enemy := range enemies {
            if enemy.Pos().Y < 0 {
                // Bloquea el acceso a las vidas para evitar conflictos de concurrencia
                vidaMutex.Lock()

                // Actualiza las vidas y resetea el enemigo
                vida.PerderVida()
                enemy.Reset()

                // Libera el bloqueo de las vidas
                vidaMutex.Unlock()
            }
        }
        // Agrega una pausa para controlar la frecuencia de verificación
        time.Sleep(time.Millisecond * 100)
    }
}

func gestionarTiempo(win *pixelgl.Window, tiempoTranscurrido float64, vida *models.Vida) {
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

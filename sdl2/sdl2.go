package main

import "fmt"
import "math/rand"
import "time"
import "github.com/veandco/go-sdl2/sdl"

const winWidth int = 800
const winHeight int = 600

var score1 int = 0
var score2 int = 0

type colour struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type object struct {
	pos
	colour colour
}

type ball struct {
	object
	radius int
	xvel   float32
	yvel   float32
}

type paddle struct {
	object
	width  int
	height int
}

func setPixel(x, y int, c colour, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.width/2
	startY := int(paddle.y) - paddle.height/2

	for y := 0; y < paddle.height; y++ {
		for x := 0; x < paddle.width; x++ {
			setPixel(startX+x, startY+y, paddle.colour, pixels)
		}
	}
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.colour, pixels)
			}
		}
	}
}

func (ball *ball) update(pixels []byte, paddle1 paddle, paddle2 paddle) {
	rand.Seed(time.Now().UnixNano())

	ball.x += ball.xvel
	ball.y += ball.yvel

	if ball.y < 0 || int(ball.y) > winHeight {
		ball.yvel = -ball.yvel
	}

	if ball.x < 0 {
		score2++
		ball.x = 400
		ball.y = float32(rand.Intn(winHeight))
	}
	if int(ball.x) > winWidth {
		score1++
		ball.x = 400
		ball.y = float32(rand.Intn(winHeight))
	}

	//&& int(ball.x) >= int(paddle1.x)-paddle1.width/2	            && int(ball.x) >= int(paddle2.x)-paddle1.width/2

	if (int(ball.x)-ball.radius == int(paddle1.x)+paddle1.width/2) && (int(ball.y) <= int(paddle1.y)+paddle1.height/2 && int(ball.y) >= int(paddle1.y)-paddle1.height/2) {
		ball.xvel = -ball.xvel
	}
	if (int(ball.x)+ball.radius == int(paddle2.x)-paddle2.width/2) && (int(ball.y) <= int(paddle2.y)+paddle2.height/2 && int(ball.y) >= int(paddle2.y)-paddle2.height/2) {
		ball.xvel = -ball.xvel
	}

}

func (paddle *paddle) update(keyState []uint8, player int) {
	if player == 1 {
		if keyState[sdl.SCANCODE_W] != 0 && int(paddle.y)-paddle.height/2 > 0 {
			paddle.y--
		}
		if keyState[sdl.SCANCODE_S] != 0 && int(paddle.y)+paddle.height/2 < winHeight {
			paddle.y++
		}
	} else if player == 2 {
		if keyState[sdl.SCANCODE_UP] != 0 && int(paddle.y)-paddle.height/2 > 0 {
			paddle.y--
		}
		if keyState[sdl.SCANCODE_DOWN] != 0 && int(paddle.y)+paddle.height/2 < winHeight {
			paddle.y++
		}
	}
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winHeight*winWidth*4)

	for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			setPixel(x, y, colour{0, 0, 0}, pixels)
		}
	}

	keyState := sdl.GetKeyboardState()

	paddle1 := paddle{object{pos{100, 300}, colour{255, 255, 255}}, 15, 100}
	paddle2 := paddle{object{pos{685, 300}, colour{255, 255, 255}}, 15, 100}
	ball := ball{object{pos{400, 300}, colour{255, 0, 0}}, 10, 0.3, 0.3}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("PLAYER 1:", score1, "\nPLAYER 2:", score2)
				return
			}
		}
		clear(pixels)
		paddle1.draw(pixels)
		paddle1.update(keyState, 1)

		paddle2.draw(pixels)
		paddle2.update(keyState, 2)

		ball.draw(pixels)
		ball.update(pixels, paddle1, paddle2)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
	}
}

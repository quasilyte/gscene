//go:build example

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/quasilyte/gscene"
)

// This example illustrates how to transition from one scene to another.
// This example uses the console output (prints) instead of graphics.
// For simplicity, it uses a global object for the game's context.

type gameContext struct {
	sceneManager *gscene.Manager
	screenWidth  int
	screenHeight int
}

var gctx = &gameContext{
	sceneManager: gscene.NewManager(),
	screenWidth:  640,
	screenHeight: 480,
}

func main() {
	g := &myGame{}

	gctx.sceneManager.ChangeScene(&myFirstSceneController{})

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

type myGame struct{}

func (g *myGame) Layout(int, int) (int, int) {
	return gctx.screenWidth, gctx.screenHeight
}

func (g *myGame) Update() error {
	gctx.sceneManager.Update()
	return nil
}

func (g *myGame) Draw(screen *ebiten.Image) {
	gctx.sceneManager.Draw(screen)
}

type myFirstSceneController struct{}

func (c *myFirstSceneController) Init(scene *gscene.Scene) {
	fmt.Println("running scene 1")
	fmt.Println("> press enter to change the scene")
}

func (c *myFirstSceneController) Update(delta float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		gctx.sceneManager.ChangeScene(&mySecondSceneController{})
	}
}

type mySecondSceneController struct{}

func (c *mySecondSceneController) Init(scene *gscene.Scene) {
	fmt.Println("running scene 2")
	fmt.Println("> press enter to change the scene back")
}

func (c *mySecondSceneController) Update(delta float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		gctx.sceneManager.ChangeScene(&myFirstSceneController{})
	}
}

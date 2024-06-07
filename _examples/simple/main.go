//go:build example

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/quasilyte/gscene"
)

// This simple example illustrates how to create a simple single scene,
// implement your own scene controller (mySceneController),
// scene object (myObject), and even scene graphics (myLabel).
//
// Usually, you can use the https://github.com/quasilyte/ebitengine-graphics
// library to get many graphical primitives like sprites, labels, geometrical shapes.

// Normally, you would have some way to store this game-wide information.
// It could be a global variable.
// It could be an explicit state object passed around (in which case
// you can access it via Controller).
var (
	random       = rand.New(rand.NewSource(time.Now().UnixNano()))
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	g := &myGame{}

	g.sceneManager = gscene.NewManager()
	g.sceneManager.ChangeScene(&mySceneController{})

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

// myGame implements [ebiten.Game] interface.
// It's our top-level game runner that should call
// the current scene's Update and Draw methods.
type myGame struct {
	sceneManager *gscene.Manager
}

func (g *myGame) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func (g *myGame) Update() error {
	g.sceneManager.Update()
	return nil
}

func (g *myGame) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

type mySceneController struct {
	seq        int
	scene      *gscene.Scene
	spawnDelay float64
}

func (c *mySceneController) Init(scene *gscene.Scene) {
	c.scene = scene
}

func (c *mySceneController) Update(delta float64) {
	c.spawnDelay -= delta
	if c.spawnDelay <= 0 {
		c.spawnDelay = 2 * random.Float64()
		o := &myObject{id: c.seq}
		c.scene.AddObject(o)
		c.seq++
	}
}

// myObject implements [gscene.Object].
// It's marked as disposed after it reaches somewhere around the center of the screen.
// It's assigned a randomized speed upon initialization.
// It also uses a label object as its graphics.
type myObject struct {
	id    int
	pos   [2]float64
	speed float64
	label *myLabel
}

func (o *myObject) Dispose() {
	o.label.Dispose()
}

func (o *myObject) IsDisposed() bool {
	return o.label.IsDisposed()
}

func (o *myObject) Init(scene *gscene.Scene) {
	o.pos[1] = random.Float64() * float64(screenHeight)

	o.speed = 40 * (random.Float64() + 0.5)

	// Note: we're "binding" the position of the graphics
	// to the logical object field.
	// This way, there is only one source of truth: the object's pos value.
	// The object itself updates the position inside its update
	// while the bound graphics just read that new value through the pointer.
	o.label = &myLabel{
		text: fmt.Sprintf("object%d", o.id),
		pos:  &o.pos,
	}
	scene.AddGraphics(o.label)
}

func (o *myObject) Update(delta float64) {
	// Slide throught the X axis and check whether we
	// should consider this object to be destroyed.

	o.pos[0] += o.speed * delta

	if o.pos[0] >= 0.5*float64(screenWidth) {
		o.Dispose()
	}
}

// myLabel implements [gscene.Graphics] interface.
// It renders the provided text at the owner's object position
// using the debug print function.
// Note that this is a common pattern: graphical objects
// should have a pointer to a position, because they don't
// "own" that position, they just need a way to read the value.
type myLabel struct {
	text     string
	pos      *[2]float64
	disposed bool
}

func (l *myLabel) Dispose() { l.disposed = true }

func (l *myLabel) IsDisposed() bool {
	return l.disposed
}

func (l *myLabel) Draw(dst *ebiten.Image) {
	ebitenutil.DebugPrintAt(dst, l.text, int(l.pos[0]), int(l.pos[1]))
}

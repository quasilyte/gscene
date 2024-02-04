package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Controller is a scene-attached object that initializes and runs a single scene.
// It's up to the controller to create all necessary objects and add them to the scene.
//
// There is always only one active controller for the scene.
type Controller[T any] interface {
	// Init is called once when a new scene is being created.
	Init(*Scene[T])

	// Update is called at every game's Update cycle.
	// The controller's Update is called before any of the scene objects Update.
	Update(delta float64)
}

// Object is a scene-managed object those Update method will be called
// as a part of a game loop.
type Object[T any] interface {
	// Init is called once when object is added to the scene.
	//
	// It's a good time to initialize all dependent objects
	// and attach sprites to the scene.
	Init(*Scene[T])

	// IsDisposed reports whether scene object was disposed.
	//
	// Disposed objects are removed from the scene before their
	// Update method is called for the current frame.
	IsDisposed() bool

	// Update is called for every object during every logical game frame.
	// Delta is passed via the Scene.Update method.
	// It could be a fixed value like 1.0/60 or a computed delta.
	Update(delta float64)
}

// Graphics is a scene-managed graphical object those Draw method will be called
// as a part of a game loop.
//
// You rarely need to write your own Graphics implementation.
// You can find the most popular implementations like Sprite on the internet.
type Graphics interface {
	// Draw implements the rendering method of this graphics object.
	Draw(dst *ebiten.Image)

	// IsDisposed reports whether graphics object was disposed.
	//
	// Disposed graphics are removed from the scene before their
	// Draw method is called for the current frame.
	IsDisposed() bool
}

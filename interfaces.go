package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InitContext is an argument type for [Controller.Init].
// Most notably, the [Scene] is directly available through its field.
type InitContext struct {
	Scene *Scene
}

// SetDrawer changes the scene [Drawer] implementation.
//
// The default Drawer is a single-layer implementation
// that ignores layer index argument of AddGraphics and
// renders all objects in the order they were added.
// See [Drawer] docs to learn more about how to implement a custom drawer.
func (ctx *InitContext) SetDrawer(d Drawer) {
	ctx.Scene.setDrawer(d)
}

// Controller is a scene-attached object that initializes and runs a single scene.
// It's up to the controller to create all necessary objects and add them to the scene.
//
// There is always only one active controller for the scene.
//
// The [Controller] interface is very similar to [Object] interface,
// but it's never Disposed as the controller's lifetime is equal
// to the current scene's lifetime.
// Also, instead of just a [Scene], it gets some extra data for its initialization.
type Controller interface {
	// Init is called once when a new scene is being created.
	Init(ctx InitContext)

	// Update is called at every game's Update cycle.
	// The controller's Update is called before any of the scene objects Update.
	Update(delta float64)
}

// Object is a scene-managed object those [Update] method will be called
// as a part of a game loop.
//
// When its [IsDisposed] method returns true, it's removed from the scene.
type Object interface {
	// Init is called once when object is added to the scene.
	//
	// It's a good time to initialize all dependent objects
	// and attach sprites to the scene.
	Init(*Scene)

	// IsDisposed reports whether scene object was disposed.
	//
	// Disposed objects are removed from the scene before their
	// Update method is called for the current frame.
	IsDisposed() bool

	// Update is called for every object during every logical game frame.
	// Delta is passed via the [Scene.Update] method.
	// It could be a fixed value like 1.0/60 or a computed delta.
	Update(delta float64)
}

// Graphics is a scene-managed graphical object those Draw method will be called
// as a part of a game loop.
//
// You rarely need to write your own [Graphics] implementation.
// You can find the most popular implementations like Sprite
// in ebitengine-graphics package.
type Graphics interface {
	// Draw implements the rendering method of this graphics object.
	Draw(dst *ebiten.Image)

	// IsDisposed reports whether graphics object was disposed.
	//
	// Disposed graphics are removed from the scene before their
	// Draw method is called for the current frame.
	IsDisposed() bool
}

// Drawer implements a smart drawable objects container.
//
// [Scene] itself holds simple objects like [Object], but graphics are more complicated.
// There are layers, cameras, and other stuff that needs to be handled properly.
// This is why drawing can be configured via the interface.
//
// There is a default implementation available plus some more in third-party libraries
// like ebitengine-graphics.
type Drawer interface {
	// AddGraphics is like [Scene.AddObject], but for [Graphics].
	//
	// The provided layer index specifies which layer should handle
	// this graphic rendering.
	// Normally, layers start from 0 go up.
	// Higher layers are drawned on top of lower ones.
	//
	// A layer can do some graphics ordering inside itself as well.
	// For example, a Y-sort style layer would draw its elements
	// after sorting them by Y-axis.
	AddGraphics(g Graphics, layer int)

	// Draw is a [Drawer] hook into [ebiten.Game] Draw tree.
	// The [Manager.Draw] will call the current Drawer's Draw method.
	//
	// The drawer is expected to draw all its layers to the [dst] image.
	Draw(dst *ebiten.Image)

	// Update is a [Drawer] hook into [ebiten.Game] Update tree.
	// The [Manager.Update] will call the current Drawer's Update method.
	//
	// The drawer is not expected to do anything during this method,
	// but it might be a good place to filter-out disposed graphical objects.
	// Doing so inside the update tree might be better to waste less
	// CPU cycles for irrelevant task inside the draw tree.
	Update(delta float64)
}

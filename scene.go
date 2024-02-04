package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene creates a logical scope and lifetime for game objects and graphics.
// Your root ebiten.Game will now only call the current scene's Update and Draw
// instead of trying to manage all the objects and graphics.
//
// When the scene goes away (the root game drops its reference),
// all scene-bound resources usually go away as well (unless you keep)
// the pointer to them somewhere else.
// Therefore, you should avoid the global state whether possible.
// Use the scene-bound context object feature to implement
// global-like resource access (there is still a chance to retain
// some references via that context object, but at least it will be
// constrained inside that object as opposed to the entire global scope).
//
// A scene can be parametrized by a context-like type T.
// It will be accessible for every object that has the access to the scene.
// You might want to make a game-local type alias like this to avoid
// spelling the type parameter every time:
//
//	type Scene = gscene.Scene[*mygame.Context]
//
// After that, you can use your Scene type alias everywhere.
//
// Note that more often than not you want to use a pointer-typed context.
type Scene[T any] struct {
	controller Controller[T]

	objects      []Object[T]
	addedObjects []Object[T]

	graphics []Graphics

	context T
}

// NewScene allocates a new scene bound to the given controller.
// The controller's Init will be called in the process.
//
// The provided context will be stored inside Scene.Context field
// and will be accessible later.
// The scene doesn't try to use that context, it's for the user-side
// to interpret it.
func NewScene[T any](ctx T, controller Controller[T]) *Scene[T] {
	s := &Scene[T]{
		controller:   controller,
		objects:      make([]Object[T], 0, 32),
		addedObjects: make([]Object[T], 0, 8),
		context:      ctx,
	}
	controller.Init(s)
	return s
}

// Context returns the context bound to this scene at the moment of its creation.
func (s *Scene[T]) Context() T { return s.context }

// Update calls the Update methods on the entire scene tree.
//
// First, it calls an Update on the controller.
//
// Then it calls the Update methods on scene objects that are not disposed.
// The Update call order is identical to the AddObject order that was used before.
//
// Disposed object are removed from the objects list.
func (s *Scene[T]) Update(delta float64) {
	// The scene controller receives the Update call first.
	s.controller.Update(delta)

	// Call every active object's Update, filter
	// the objects list in-place while at it.
	liveObjects := s.objects[:0]
	for _, o := range s.objects {
		if o.IsDisposed() {
			continue
		}
		o.Update(delta)
		liveObjects = append(liveObjects, o)
	}
	s.objects = liveObjects

	// Flush the added objects to the list.
	s.objects = append(s.objects, s.addedObjects...)
	s.addedObjects = s.addedObjects[:0]
}

// Draw calls the Draw methods on the entire scene tree.
//
// It calls the Draw methods on scene graphics that are not disposed.
// The Draw call order is identical to the AddGraphics order that was used before.
//
// Disposed graphics are removed from the objects list.
func (s *Scene[T]) Draw(screen *ebiten.Image) {
	// Just like in Update. The only difference is the method
	// being called is Draw, not Update (and there is no delta).
	liveGraphics := s.graphics[:0]
	for _, g := range s.graphics {
		if g.IsDisposed() {
			continue
		}
		g.Draw(screen)
		liveGraphics = append(liveGraphics, g)
	}
	s.graphics = liveGraphics
}

// AddObject adds the logical object to the scene.
// Its Init method will be called right away.
//
// The AddObject method adds the object to the add-queue.
// The object will be actually added at the end of the current
// Update method's life cycle.
func (s *Scene[T]) AddObject(o Object[T]) {
	s.addedObjects = append(s.addedObjects, o)
	o.Init(s)
}

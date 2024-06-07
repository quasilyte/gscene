package gscene

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
type Scene struct {
	root *RootScene
}

// Controller returns the bound controller.
// It can be used to access some scene-specific typed data if necessary.
func (s *Scene) Controller() Controller {
	return s.root.controllerObject
}

// AddObject adds the logical object to the scene.
// Its [Object.Init] method will be called right away.
//
// The [AddObject] method adds the object to the add-queue.
// The object will be actually added at the end of the current
// Update method's life cycle.
//
// This object will be automatically removed from the scene
// when its [Object.IsDisposed] method reports true.
//
// All added objects are stored inside the scene.
// If they're only reachable between each other and the scene,
// they can be easily garbage-collected as soon as this scene
// will be garbage-collected (there is usually only 1 active scene at a time).
func (s *Scene) AddObject(o Object) {
	s.root.AddObject(o)
}

// AddGraphics adds the graphical object to the scene.
//
// This object will be automatically removed from the scene
// when its [Graphics.IsDisposed] method reports true.
//
// All added objects are stored inside the scene.
// If they're only reachable between each other and the scene,
// they can be easily garbage-collected as soon as this scene
// will be garbage-collected (there is usually only 1 active scene at a time).
func (s *Scene) AddGraphics(g Graphics) {
	s.root.AddGraphics(g)
}

package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene creates a logical scope and lifetime for game objects and graphics.
//
// Your root ebiten.Game will now only call the current scene's Update and Draw
// instead of trying to manage all the objects and graphics.
// (Note: you call these through the scene [Manager].)
//
// When the scene goes away, all scene-bound resources usually go away as well
// (unless you keep the pointer to them somewhere else).
// Therefore, you should avoid the unnecessary global state whether possible.
type Scene struct {
	controllerObject Controller
	drawer           Drawer

	objects      []Object
	addedObjects []Object

	insideUpdate bool
}

type stopUpdateType struct{}

var stopUpdate any = &stopUpdateType{}

// newScene allocates a new scene bound to the given controller.
//
// It's the caller's responsibility to call [Controller.Init]
// with the created scene object.
func newScene(c Controller) *Scene {
	scene := &Scene{
		controllerObject: c,
		objects:          make([]Object, 0, 32),
		addedObjects:     make([]Object, 0, 8),
	}
	return scene
}

func (s *Scene) Controller() Controller {
	return s.controllerObject
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
	s.addedObjects = append(s.addedObjects, o)
	o.Init(s)
}

// AddGraphics is a shorthand for s.Viewport(0).AddGraphics(g, layer).
func (s *Scene) AddGraphics(g Graphics, layer int) {
	s.drawer.Viewport(0).AddGraphics(g, layer)
}

// Viewport returns the nth viewport associated with the current scene [Drawer].
// If you only have a single [Viewport], use [AddGraphics] method instead.
func (s *Scene) Viewport(index int) Viewport {
	return s.drawer.Viewport(index)
}

// dispose stops the current scene execution (even mid-update) and
// discards the scene state.
//
// This is useful before switching the scene if you want to
// abort the execution of the rest of the current Update tree.
//
// Calling Dispose is valid when either outside or inside of the Update call.
// It's not valid to call it when inside the Draw tree.
//
// After this scene is disposed, it should not be used any further.
func (s *Scene) dispose() {
	s.objects = nil
	s.addedObjects = nil
	s.controllerObject = nil
	s.drawer = nil

	if s.insideUpdate {
		s.insideUpdate = false
		panic(stopUpdate)
	}
}

func (s *Scene) update() {
	s.updateWithDelta(1.0 / 60.0)
}

func (s *Scene) updateWithDelta(delta float64) {
	// We have two methods: updateWithDelta and updateWithDeltaImpl.
	// updateWithDelta is needed to create a guarding defer call
	// that would catch the update cancelling message.
	// updateWithDeltaImpl implements the actual update logic.

	defer func() {
		rv := recover()
		if rv == nil {
			return // The most common case
		}
		if rv == stopUpdate {
			// This is our way to break out of the Update
			// tree; no need to re-panic it.
			return
		}
		// Some real panic is happening.
		panic(rv)
	}()

	s.insideUpdate = true
	s.updateWithDeltaImpl(delta)
	s.insideUpdate = false
}

func (s *Scene) updateWithDeltaImpl(delta float64) {
	// The scene controller receives the Update call first.
	s.controllerObject.Update(delta)

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

	// Drawer's update is called the last.
	s.drawer.Update(delta)

	// Flush the added objects to the list.
	s.objects = append(s.objects, s.addedObjects...)
	s.addedObjects = s.addedObjects[:0]
}

func (s *Scene) draw(dst *ebiten.Image) {
	s.drawer.Draw(dst)
}

func (s *Scene) setDrawer(d Drawer) {
	// A simple sanity check.
	if d, ok := d.(*simpleDrawer); ok {
		if len(d.graphics) > 0 {
			panic("setting a drawer after graphics were already added")
		}
	}

	s.drawer = d
}

package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type RootScene struct {
	// Since we can't combine 2 interface parts into one,
	// we'll use two interface-typed objects here.
	// In practice, both of them will have the same underlying object
	// that implements the scene controller.
	controllerObject Controller

	objects      []Object
	addedObjects []Object

	graphics []Graphics

	// This is a single object that is used for every non-root object.
	asScene *Scene

	insideUpdate bool
}

type stopUpdateType struct{}

var stopUpdate any = &stopUpdateType{}

// NewRootScene allocates a new root scene bound to the given controller.
// The [Controller.Init] will be called in the process.
func NewRootScene(c Controller) *RootScene {
	root := &RootScene{
		controllerObject: c,
		objects:          make([]Object, 0, 32),
		addedObjects:     make([]Object, 0, 8),
	}
	root.asScene = &Scene{
		root: root,
	}
	c.Init(root)
	return root
}

// Dispose stops the current scene execution (even mid-update) and
// discards the scene state.
//
// This is useful before switching the scene if you want to
// abort the execution of the rest of the current Update tree.
//
// Calling Dispose is valid when either outside or inside of the Update call.
// It's not valid to call it when inside the Draw tree.
//
// After this scene is disposed, it should not be used any further.
func (s *RootScene) Dispose() {
	s.objects = nil
	s.addedObjects = nil
	s.graphics = nil
	s.controllerObject = nil

	if s.insideUpdate {
		s.insideUpdate = false
		panic(stopUpdate)
	}
}

func (s *RootScene) AsScene() *Scene {
	return s.asScene
}

// Update is a shorthand for [UpdateWithDelta](1.0/60.0).
func (s *RootScene) Update() {
	s.UpdateWithDelta(1.0 / 60.0)
}

// UpdateWithDelta calls the Update methods on the entire scene tree.
//
// First, it calls the bound [Controller.Update].
//
// Then it calls the [Object.Update] methods on scene objects that are not disposed.
// The Update call order is identical to the AddObject order that was used before.
//
// Disposed object are removed from the objects list.
func (s *RootScene) UpdateWithDelta(delta float64) {
	// We have two methods: UpdateWithDelta and updateWithDelta.
	// UpdateWithDelta is needed to create a guarding defer call
	// that would catch the update cancelling message.
	// updateWithDelta implements the actual update logic.

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
	s.updateWithDelta(delta)
	s.insideUpdate = false
}

func (s *RootScene) updateWithDelta(delta float64) {
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
func (s *RootScene) Draw(screen *ebiten.Image) {
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
func (s *RootScene) AddObject(o Object) {
	s.addedObjects = append(s.addedObjects, o)
	o.Init(s.asScene)
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
func (s *RootScene) AddGraphics(g Graphics) {
	s.graphics = append(s.graphics, g)
}

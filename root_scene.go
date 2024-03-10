package gscene

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// SimpleRootScene is a type alias for a scene that doesn't need
// its object to have any typed access to the controller via the scene.
//
// This kind of a scene is a norm for game screens that are mostly
// ui-focused instead of being gameplay-rich.
//
// You may not need this type alias and that's OK too.
// It serves an example that unnecessary generics parameters can
// be hiden when wanted.
type SimpleRootScene = RootScene[any]

type RootScene[ControllerAccessor any] struct {
	// Since we can't combine 2 interface parts into one,
	// we'll use two interface-typed objects here.
	// In practice, both of them will have the same underlying object
	// that implements the scene controller.
	controllerObject   Controller[ControllerAccessor]
	controllerAccessor ControllerAccessor

	objects      []Object[ControllerAccessor]
	addedObjects []Object[ControllerAccessor]

	graphics []Graphics

	// This is a single object that is used for every non-root object.
	asScene *Scene[ControllerAccessor]
}

// NewRootScene allocates a new root scene bound to the given controller.
// The controller's Init will be called in the process.
//
// The controller c can optionally implement an ControllerAccessor interface
// to provide some data access to the scene objects.
// For example, that interface may provide some shared game context and/or
// graphical layers APIs (like AddGraphicsToLayer).
//
// If c doesn't implement ControllerAccessor, a nil value for the accessor will be used.
// If ControllerAccessor is any (an empty interface), it'll have almost the same
// meaning, but the scene objects may do a type assertion and query the data directly.
// This second method is not recommended and it will only work if both controller and
// objects are defined in the same package (therefore an object can have a controller's
// type available for the type assertion).
func NewRootScene[ControllerAccessor any](c Controller[ControllerAccessor]) *RootScene[ControllerAccessor] {
	accessor, ok := c.(ControllerAccessor)
	if !ok {
		// This is a sanity check.
		// If ControllerAccessor is any, anything will implement it.
		// If ControllerAccessor is not any, the library user wants to implement
		// that interface by their controller in 99.(9)% cases.
		panic(fmt.Sprintf("given controller doesn't implement %T (ControllerAccessor interface)", (*ControllerAccessor)(nil)))
	}
	root := &RootScene[ControllerAccessor]{
		controllerObject:   c,
		controllerAccessor: accessor,
		objects:            make([]Object[ControllerAccessor], 0, 32),
		addedObjects:       make([]Object[ControllerAccessor], 0, 8),
	}
	root.asScene = &Scene[ControllerAccessor]{
		root: root,
	}
	c.Init(root)
	return root
}

func (s *RootScene[ControllerAccessor]) AsScene() *Scene[ControllerAccessor] {
	return s.asScene
}

// Update is a shorthand for UpdateWithDelta(1.0/60.0).
func (s *RootScene[ControllerAccessor]) Update() {
	s.UpdateWithDelta(1.0 / 60.0)
}

// UpdateWithDelta calls the Update methods on the entire scene tree.
//
// First, it calls an Update on the controller.
//
// Then it calls the Update methods on scene objects that are not disposed.
// The Update call order is identical to the AddObject order that was used before.
//
// Disposed object are removed from the objects list.
func (s *RootScene[ControllerAccessor]) UpdateWithDelta(delta float64) {
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
func (s *RootScene[T]) Draw(screen *ebiten.Image) {
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
//
// This object will be automatically removed from the scene
// when its IsDisposed method will report true.
//
// All added objects are stored inside the scene.
// If they're only reachable between each other and the scene,
// they can be easily garbage-collected as soon as this scene
// will be garbage-collected (there is usually only 1 active scene at a time).
func (s *RootScene[ControllerAccessor]) AddObject(o Object[ControllerAccessor]) {
	s.addedObjects = append(s.addedObjects, o)
	o.Init(s.asScene)
}

// AddGraphics adds the graphical object to the scene.
//
// This object will be automatically removed from the scene
// when its IsDisposed method will report true.
//
// All added objects are stored inside the scene.
// If they're only reachable between each other and the scene,
// they can be easily garbage-collected as soon as this scene
// will be garbage-collected (there is usually only 1 active scene at a time).
func (s *RootScene[ControllerAccessor]) AddGraphics(g Graphics) {
	s.graphics = append(s.graphics, g)
}

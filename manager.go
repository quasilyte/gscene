package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Manager wraps the current scene and implements scene changing logic.
//
// It also provides the access to Update/Draw methods that should
// be used from the top-level game runner of [ebiten.Game].
//
// Most games only need one scene [Manager].
// Put it somewhere in your game's context.
type Manager struct {
	currentScene *Scene
	disposed     bool
}

func NewManager() *Manager {
	return &Manager{}
}

// ChangeScene changes the current scene to a new one.
// The new scene will have the specified controller attached to it.
//
// If there is another scene running during the time [ChangeScene]
// is called, its execution will be stopped.
// This means that ChangeScene should be treated as a control transfer
// call, it will not return and continue from the point it was called.
// After the scene is changed, no logic that is part of the Update tree
// from the old scene will be executed.
//
// The [Controller.Init] method of [c] will be called after
// this new scene is installed.
func (m *Manager) ChangeScene(c Controller) {
	prevScene := m.currentScene

	m.currentScene = newScene(c)
	m.currentScene.drawer = newSimpleDrawer()
	c.Init(InitContext{Scene: m.currentScene})

	if prevScene != nil {
		prevScene.dispose()
	}
}

func (m *Manager) CurrentScene() *Scene {
	return m.currentScene
}

func (m *Manager) IsDisposed() bool {
	return m.disposed
}

func (m *Manager) Dispose() {
	m.disposed = true
}

// Update is a shorthand for [UpdateWithDelta](1.0/60.0).
func (m *Manager) Update() {
	m.currentScene.update()
}

// UpdateWithDelta calls the Update methods on the entire scene tree.
//
// First, it calls the bound [Controller.Update].
//
// Then it calls the [Object.Update] methods on scene objects that are not disposed.
// The Update call order is identical to the AddObject order that was used before.
//
// Disposed object are removed from the objects list.
func (m *Manager) UpdateWithDelta(delta float64) {
	m.currentScene.updateWithDelta(delta)
}

// Draw calls the Draw methods on the entire scene tree.
//
// It calls the Draw methods on scene graphics that are not disposed.
// The Draw call order is identical to the AddGraphics order that was used before.
//
// Disposed graphics are removed from the objects list.
func (m *Manager) Draw(dst *ebiten.Image) {
	m.currentScene.draw(dst)
}

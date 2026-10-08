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

	pendingScene ChangeSceneConfig
	sceneStack   []*Scene

	disposed bool
}

func NewManager() *Manager {
	return &Manager{}
}

// SceneChangeIsPending reports whether there is a pending
// pushed scene that will take over after the current update tick is over.
//
// Since using ChangeScene is not supported during this period, use
// this function if you need to guard ChangeScene call.
func (m *Manager) SceneChangeIsPending() bool {
	return m.pendingScene.Controller != nil
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
//
// ChangeScene is illegal if there is a pushed scene pending.
// Use [SceneChangeIsPending] to know for sure.
func (m *Manager) ChangeScene(c Controller) {
	m.ChangeSceneConfig(ChangeSceneConfig{Controller: c})
}

type ChangeSceneConfig struct {
	Controller   Controller
	PanicHandler func(v any)
}

// ChangeSceneConfig is like ChangeScene, but allows more configuration.
func (m *Manager) ChangeSceneConfig(config ChangeSceneConfig) {
	if m.pendingScene.Controller != nil {
		panic("changing a scene while other scene was pushed")
	}

	prevScene := m.currentScene

	c := config.Controller

	m.currentScene = newScene(c)
	m.currentScene.panicHandler = config.PanicHandler
	m.currentScene.drawer = newSimpleDrawer()
	c.Init(InitContext{Scene: m.currentScene})

	if prevScene != nil {
		prevScene.dispose()
	}
}

// PushScene suspends the current scene and places a new scene on top,
// making it active. When PopScene is called, the pushed scene is disposed
// and the suspended scene becomes active again.
//
// Unlike ChangeScene, this does not transfer the control in any way.
// The function returns normally. The currently executing scene continues
// its current update. The pushed scene becomes active on the next update.
//
// The pushed scene is not initialized until the end of the current tick.
// The CurrentScene() will still return the running scene inside the tick.
// You should not push several scenes before the previous one initialized.
//
// The reason behind that is state consistency. When throwing away the scene,
// you clearly tell that it's irrelevant to care about its mid-update state.
// With preserved state, picking the right interruption spot can be tricky
// as it can alter the result: imagine transitioning before and after bot tick.
// Only the latter will result in bot getting the compute it expected to have,
// which can be important. You still need to have extra care about when to run
// the scene push, but at the very least this position doesn't affect whichever
// objects are going to be updated during the running tick.
//
// The main use case for a scene stack (push+pop) are layered scenes
// that can benefit from preserved current scene states.
// For instance, a combat scene can push a unit editor scene.
// You might now want to lose a complex scene state just because you
// opened a new scene on top of it, and this is where scene stacks come into play.
// Pushing the next scene preserves the state and allows an easy and fast restoration.
func (m *Manager) PushScene(config ChangeSceneConfig) {
	if m.pendingScene.Controller != nil {
		panic("nested PushScene are not supported")
	}
	if config.Controller == nil {
		panic("can't push a nil controller")
	}

	m.pendingScene = config
}

// PopScene disposes the currently running scene and restores
// the previously suspended scene.
//
// Popping the scene does interrupt the current scene flow,
// so it's closer to ChangeScene in that regard.
//
// See [PushScene] comment to understand its design better.
func (m *Manager) PopScene() {
	if len(m.sceneStack) == 0 {
		panic("no scenes to pop")
	}

	restoredScene := m.sceneStack[len(m.sceneStack)-1]
	m.sceneStack = m.sceneStack[:len(m.sceneStack)-1]

	prevScene := m.currentScene

	m.currentScene = restoredScene

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

func (m *Manager) maybeTransitToPending() {
	if m.pendingScene.Controller == nil {
		return
	}

	m.sceneStack = append(m.sceneStack, m.currentScene)

	config := m.pendingScene
	m.pendingScene = ChangeSceneConfig{}

	c := config.Controller

	m.currentScene = newScene(c)
	m.currentScene.panicHandler = config.PanicHandler
	m.currentScene.drawer = newSimpleDrawer()
	c.Init(InitContext{Scene: m.currentScene})
}

// Update is a shorthand for [UpdateWithDelta](1.0/60.0).
func (m *Manager) Update() {
	m.currentScene.update()
	m.maybeTransitToPending()
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
	m.maybeTransitToPending()
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

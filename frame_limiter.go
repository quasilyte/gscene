package gscene

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// FrameLimiter is a helper to make FPS limiting in Ebitengine games easier.
//
// You should use a single limiter in your game to schedule Draw calls.
//
// Use it like this:
//
//	func(g *game) Draw(screen *ebiten.Image) {
//	  g.limiter.Do(screen, g.drawAll)
//	}
//
// Where drawAll is your original Draw method that does
// all the drawing.
type FrameLimiter struct {
	fps   uint
	dirty bool

	timeAccum    time.Duration
	drawDelay    time.Duration
	prevDrawTime time.Time
}

// NewFrameLimiter returns an initialized frame limiter.
// You can change the initial fps cap by using [SetFPS].
// Using a value of 0 means "no FPS cap".
func NewFrameLimiter(fps uint) *FrameLimiter {
	l := &FrameLimiter{dirty: true}
	l.SetFPS(fps)
	return l
}

// SetDirty changes the limiter's dirty flag.
//
// By default it starts with dirty=true and every Do call
// considers the game state as dirty.
//
// You can set dirty=false inside a Do callback
// and then set it back to dirty=true somewhere in your
// game's logic. If at the momeny Do is called, it will
// do nothing as long as dirty flag is false.
//
// The simplest way to set dirty=true flag is to do it inside
// your game's Update tree. It will be helpful if your game
// has TPS value comparable to the FPS (or higher).
// A better approach is to know for sure when the game's content
// needs to be re-drawn but it might turn out to be too error-prone.
func (l *FrameLimiter) SetDirty(dirty bool) {
	l.dirty = dirty
}

func (l *FrameLimiter) GetFPS() uint {
	return l.fps
}

// SetFPS adds a soft-limit on a Draw() frequency.
//
// If non-zero, the it will try to skip unnecessary draw calls.
//
// It will only work with ebiten.SetScreenClearedEveryFrame(false),
// with an optional VSync=true (seems to be recommended).
func (l *FrameLimiter) SetFPS(fps uint) {
	l.fps = fps
	if fps == 0 {
		l.drawDelay = 0
	} else {
		// We're adding 1 here just to be sure not to skip
		// frames that shouldn't be skipped (or do so less often).
		l.drawDelay = time.Second / time.Duration(fps+1)
	}
}

// Do will call the provided draw function with the
// forwarded screen image as an argument if it's time
// to render the next frame.
// Otherwise, it will do nothing.
//
// It expects SetScreenClearedEveryFrame to be false,
// so it will clear the image for you before passing it
// to the draw function.
//
// See [FrameLimiter] type comment to learn more.
func (l *FrameLimiter) Do(dst *ebiten.Image, draw func(dst *ebiten.Image)) {
	if !l.dirty {
		return
	}

	if l.drawDelay != 0 {
		t := time.Now()
		delta := t.Sub(l.prevDrawTime)
		l.prevDrawTime = t

		l.timeAccum += min(delta, l.drawDelay)
		if l.timeAccum < l.drawDelay {
			return // Skip frame
		}

		l.timeAccum -= l.drawDelay
		dst.Clear()
	}

	draw(dst)
}

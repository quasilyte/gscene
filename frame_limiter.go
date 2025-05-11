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
	fps uint

	drawMinDelay time.Duration
	prevDrawTime time.Time
}

// NewFrameLimiter returns an initialized frame limiter.
// You can change the initial fps cap by using [SetFPS].
// Using a value of 0 means "no FPS cap".
func NewFrameLimiter(fps uint) *FrameLimiter {
	l := &FrameLimiter{}
	l.SetFPS(fps)
	return l
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
		l.drawMinDelay = 0
	} else {
		l.drawMinDelay = time.Second / time.Duration(fps)
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
	if l.drawMinDelay != 0 {
		t := time.Now()
		if t.Sub(l.prevDrawTime) < l.drawMinDelay {
			return // Skip frame
		}

		dst.Clear()
		l.prevDrawTime = t
	}

	draw(dst)
}

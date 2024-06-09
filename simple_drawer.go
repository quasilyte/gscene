package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type simpleDrawer struct {
	graphics []Graphics
}

func newSimpleDrawer() *simpleDrawer {
	return &simpleDrawer{}
}

func (d *simpleDrawer) Update(delta float64) {
	liveGraphics := d.graphics[:0]
	for _, g := range d.graphics {
		if g.IsDisposed() {
			continue
		}
		liveGraphics = append(liveGraphics, g)
	}
	d.graphics = liveGraphics
}

func (d *simpleDrawer) Draw(dst *ebiten.Image) {
	for _, g := range d.graphics {
		g.Draw(dst)
	}
}

func (d *simpleDrawer) AddGraphics(g Graphics, layer int) {
	if d.graphics == nil {
		d.graphics = make([]Graphics, 0, 32)
	}

	d.graphics = append(d.graphics, g)
}

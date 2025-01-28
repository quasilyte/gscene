package gscene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type simpleDrawer struct {
	graphics   []Graphics
	needFilter bool
}

func newSimpleDrawer() *simpleDrawer {
	return &simpleDrawer{}
}

func (d *simpleDrawer) Update(delta float64) {
	d.needFilter = true
}

func (d *simpleDrawer) filter() {
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
	if d.needFilter {
		d.filter()
	}
	d.needFilter = false

	for _, g := range d.graphics {
		g.Draw(dst)
	}
}

func (d *simpleDrawer) AddGraphics(g Graphics, layer int) {
	if d.graphics == nil {
		d.graphics = make([]Graphics, 0, 32)
	}

	d.graphics = append(d.graphics, g)
	d.needFilter = true
}

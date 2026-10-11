package reader

import (
	"encoding/xml"
	"fmt"

	"github.com/kinglet-dev/core/adapters/svg/syntax"
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

// start works out an element's context, pushes it, and records any shape.
func (w *walker) start(el xml.StartElement) error {
	a := attrs(el)
	ctx, err := w.context(el, a)
	if err != nil {
		return err
	}
	w.stack = append(w.stack, ctx)
	kind, ok := shapeKinds[el.Name.Local]
	if !ok {
		return nil
	}
	path, drawn, err := kind.build(a, w.budget())
	if err == nil {
		err = w.spend(path)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", describe(el), err)
	}
	if !kind.fillable {
		ctx.paint.fill = "none"
	}
	if !drawn || !ctx.paint.visible() {
		return nil
	}
	return w.addShape(el, path, ctx)
}

// context combines an element's own attributes with what it inherits. The
// root svg element sets the millimeter scale from the document's size.
func (w *walker) context(el xml.StartElement, a attributes) (context, error) {
	if len(w.stack) == 0 {
		var ctx context
		var err error
		if w.design.Width, w.design.Height, ctx.toMM, err = size(a); err != nil {
			return context{}, err
		}
		ctx.paint = initialPaint.with(a)
		return ctx, nil
	}
	parent := w.stack[len(w.stack)-1]
	ctx := context{toMM: parent.toMM, paint: parent.paint.with(a)}
	if list, ok := a["transform"]; ok {
		m, err := syntax.ParseTransform(list)
		if err != nil {
			return context{}, fmt.Errorf("%s: %w", describe(el), err)
		}
		ctx.toMM = m.Then(parent.toMM)
	}
	return ctx, nil
}

// addShape records a shape in millimeters with what the laser will score.
func (w *walker) addShape(el xml.StartElement, path geom.Path, ctx context) error {
	mm := path.Transform(ctx.toMM)
	if !inRange(mm) {
		return fmt.Errorf("%s: %w after its transforms", describe(el), syntax.ErrOutOfRange)
	}
	w.design.Shapes = append(w.design.Shapes, design.Shape{
		Path:      mm,
		Filled:    ctx.paint.filled(),
		WhiteFill: ctx.paint.white(),
		Element:   w.elements - 1, // count has numbered this element
		Fill:      ctx.paint.fill,
		Stroke:    ctx.paint.stroke,
		Local:     path,
		ToMM:      ctx.toMM,
	})
	return nil
}

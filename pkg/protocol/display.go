package protocol

import "bytes"

func (v *View) Render() string {
	return v.World.Render(v.XMin, v.XMax, v.YMin, v.YMax)
}

func (w *World) Render(xMin, xMax, yMin, yMax int32) string {
	var buf bytes.Buffer
	for y := yMax; y > yMin; y-- {
		objectsRow, hasO := w.Objects[y]
		surfacesRow, hasS := w.Surfaces[y]
		for x := xMin; x < xMax; x++ {
			if hasO {
				if o, ok := objectsRow.Columns[x]; ok {
					buf.WriteString(o.Render())
					continue
				}
			}
			if hasS {
				if s, ok := surfacesRow.Columns[x]; ok {
					buf.WriteString(s.Render())
					continue
				}
			}
			buf.WriteString((&Surface{Type: Surface_DIRT}).Render())
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func (o *Object) Render() string {
	switch o.Type {
	case Object_FOOD:
		return " @"
	case Object_QUEEN:
		return " Q"
	case Object_STONE:
		return " o"
	case Object_WORKER:
		return " A"
	default:
		return " ?"
	}
}

func (s *Surface) Render() string {
	switch s.Type {
	case Surface_DIRT:
		return " _"
	case Surface_GRASS:
		return " i"
	case Surface_HOLE:
		return " *"
	case Surface_ROCK:
		return " #"
	case Surface_SOIL:
		return " ="
	default:
		return " ?"
	}
}

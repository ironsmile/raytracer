package scene

import (
	"fmt"
	"os"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
)

var debugRaysFile string

// SetDebugRaysFile sets a name for a file which should include debug rays.
// This file is read and such rays are added to the scene as primitives.
//
// The file format is a text file with ray description per line. The file
// starts with a line which holds a single integer. This is the number of
// rays in the rest of the file. The the ray descriptions follow. A description
// is a 7 floats separated by spaces. It looks like so:
//
//	ox oy oz dx dy dz l
//
// `(ox, oy, oz)` is the ray origin. `(dx, dy, dz)` is the ray direction and
// `l` is its length.
func SetDebugRaysFile(fileName string) {
	debugRaysFile = fileName
}

func (s *Scene) initDebugRays() error {
	if debugRaysFile == "" {
		return nil
	}

	fh, err := os.Open(debugRaysFile)
	if err != nil {
		return fmt.Errorf("cannot open debug rays file: %w", err)
	}
	defer fh.Close()

	var raysCount int
	if _, err := fmt.Fscanln(fh, &raysCount); err != nil {
		return fmt.Errorf("cannot find how many rays there are: %w", err)
	}

	for i := range raysCount {
		var (
			origin geometry.Vector
			dir    geometry.Vector
			rayLen float64
		)

		if _, err := fmt.Fscanln(fh,
			&origin.X, &origin.Y, &origin.Z,
			&dir.X, &dir.Y, &dir.Z,
			&rayLen,
		); err != nil {
			return fmt.Errorf("failed while reading debug ray #%d: %w", i, err)
		}

		r := geometry.NewRay(origin, dir.Normalize())
		r.Maxt = rayLen
		dray := primitive.NewRay(r)
		primitive.SetName(dray.GetID(), fmt.Sprintf("Debug Ray #%d", i))
		s.Primitives = append(s.Primitives, dray)
	}

	return nil
}

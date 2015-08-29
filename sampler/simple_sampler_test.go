package sampler

import (
	"testing"

	"github.com/ironsmile/raytracer/film"
)

func BenchmarkGetSample(t *testing.B) {
	nullFilm := film.NewNullFilm()
	s := NewSimple(nullFilm)

	for i := 0; i < t.N; i++ {
		s.GetSample()
	}
}

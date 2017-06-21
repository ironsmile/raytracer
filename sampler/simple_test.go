package sampler

import (
	"testing"

	"github.com/ironsmile/raytracer/film"
)

func BenchmarkGetSample(t *testing.B) {
	nullFilm := film.NewNullFilm()
	s := NewSimple(nullFilm, 10)

	subSpl, err := s.GetSubSampler()

	if err != nil {
		t.Fatalf("Error getting sub sampler: %s", err)
	}

	for i := 0; i < t.N; i++ {
		subSpl.GetSample()
	}
}

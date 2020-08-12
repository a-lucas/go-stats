package gostats

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestMyStats3(t *testing.T) {

	t.Run("Mean", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			out float64
		}{
			{[]float64{1, 2, 3, 4, 5}, 3.0},
			{[]float64{1, 2, 3, 4, 5, 6}, 3.5},
			{[]float64{1}, 1.0},
			{[]float64{-1, -1, -1, 2, 3, 3, 4, 4, 5}, 2.0},
		} {
			s := NewSampleStream()
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.Mean()).To(Equal(c.out))
			s.CleanPool()
		}
	})

	t.Run("SeriesLoss/SeriesWin", func(t *testing.T) {
		g := NewGomegaWithT(t)
		s := NewSampleStream()
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(0))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(0))
		s.Append(1)
		g.Expect(s.nbPositives).To(BeEquivalentTo(1))
		g.Expect(s._count).To(BeEquivalentTo(1))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(1))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(0))

		s.Append(2)
		g.Expect(s.nbPositives).To(BeEquivalentTo(2))
		g.Expect(s._count).To(BeEquivalentTo(2))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(2))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(0))

		s.Append(-1)
		g.Expect(s.nbPositives).To(BeEquivalentTo(2))
		g.Expect(s._count).To(BeEquivalentTo(3))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(2))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(1))

		s.Append(-1)
		g.Expect(s.nbPositives).To(BeEquivalentTo(2))
		g.Expect(s._count).To(BeEquivalentTo(4))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(2))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(2))

		s.Append(-1)
		g.Expect(s.nbPositives).To(BeEquivalentTo(2))
		g.Expect(s._count).To(BeEquivalentTo(5))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(2))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(3))
		g.Expect(s._nbPositiveStreak).To(BeEquivalentTo(0))
		g.Expect(s._nbNegativeStreak).To(BeEquivalentTo(3))

		s.Append(1)
		g.Expect(s._nbPositiveStreak).To(BeEquivalentTo(1))
		g.Expect(s._nbNegativeStreak).To(BeEquivalentTo(0))

		g.Expect(s.nbPositives).To(BeEquivalentTo(3))
		g.Expect(s._count).To(BeEquivalentTo(6))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(2))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(3))

		s.Append(1)
		s.Append(1)
		g.Expect(s.nbPositives).To(BeEquivalentTo(5))
		g.Expect(s.nbProfitStreak).To(BeEquivalentTo(3))
		g.Expect(s.nbLossStreak).To(BeEquivalentTo(3))

	})

	t.Run("Variance", func(t *testing.T) {
		g := NewGomegaWithT(t)
		s := NewSampleStream()
		s.Append(1)
		s.Append(2)
		s.Append(3)
		g.Expect(s.Variance()).To(Equal(1.0))
		s.CleanPool()
	})

	t.Run("Sum", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			out float64
		}{
			{[]float64{1, 2, 3}, 6},
			{[]float64{1.0, 1.1, 1.2, 2.2}, 5.5},
			{[]float64{1, -1, 2, -3}, -1},
		} {
			s := NewSampleStream()
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.Sum()).To(Equal(c.out))
			g.Expect(s.sum).To(Equal(c.out))
			s.CleanPool()
		}
	})

	t.Run("Min/Max", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			min float64
			max float64
		}{
			{[]float64{1.1, 2, 3, 4, 5}, 1.1, 5},
			{[]float64{10.534, 3, 5, 7, 9}, 3.0, 10.534},
			{[]float64{-5, 1, 5}, -5.0, 5},
			{[]float64{5}, 5, 5},
		} {
			s := NewSampleStream()
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.min).To(Equal(c.min))
			g.Expect(s.max).To(Equal(c.max))
			s.CleanPool()
		}
	})

}

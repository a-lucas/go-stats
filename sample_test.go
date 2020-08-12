package gostats

import (
	. "github.com/onsi/gomega"
	"math"
	"sync"
	"testing"
)

func TestSubArray(t *testing.T) {
	g := NewGomegaWithT(t)
	arr := []int{1, 2, 3, 4, 5}
	g.Expect(arr[0:2]).To(Equal([]int{1, 2}))
	g.Expect(arr[1:3]).To(Equal([]int{2, 3}))
}

func TestMyStats(t *testing.T) {

	t.Run("Mean", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			out float64
		}{
			{[]float64{1, 2, 3, 4, 5}, 3.0},
			{[]float64{1, 2, 3, 4, 5, 6}, 3.5},
			{[]float64{1}, 1.0},
		} {
			s := NewSampleWithValue(c.in, false)
			g.Expect(s.Mean()).To(Equal(c.out))
		}
	})

	t.Run("Variance", func(t *testing.T) {
		g := NewGomegaWithT(t)
		s1 := NewSampleWithValue([]float64{}, false)
		g.Expect(math.IsNaN(s1.Variance())).To(BeTrue())
		s2 := NewSampleWithValue([]float64{1, 2, 3}, false)
		g.Expect(s2.Variance()).To(Equal(1.0))
	})

	t.Run("With Original", func(t *testing.T) {
		g := NewGomegaWithT(t)
		input := []float64{3, 2, 1}
		s := NewSampleWithValue(input, true)
		g.Expect(s.Xs()).To(Equal(input))
		g.Expect(s.Original()).To(Equal(input))
		s.sort()
		g.Expect(s.Xs()).To(Equal([]float64{1, 2, 3}))
		g.Expect(s.Original()).To(Equal([]float64{3, 2, 1}))
	})

	t.Run("CoVariance Population", func(t *testing.T) {
		g := NewGomegaWithT(t)
		s1 := NewSampleWithValue([]float64{1, 2, 3.5, 3.7, 8, 12}, true)
		s2 := NewSampleWithValue([]float64{10, -51.2, 8}, true)
		s3 := NewSampleWithValue([]float64{0.5, 1, 2.1, 3.4, 3.4, 4}, true)
		s4 := NewSampleWithValue([]float64{}, true)
		g.Expect(math.IsNaN(s1.CovariancePopulation(s2))).To(BeTrue())
		g.Expect(s1.CovariancePopulation(s3)).To(Equal(4.191666666666666))
		g.Expect(math.IsNaN(s1.CovariancePopulation(s4))).To(BeTrue())
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
			s := NewSampleWithValue(c.in, false)
			g.Expect(s.Sum()).To(Equal(c.out))
		}
	})

	t.Run("Percentile", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			p   float64
			out float64
		}{
			{[]float64{43, 54, 56, 61, 62, 66}, 0.9, 64.0},
			{[]float64{43}, 0.9, 43.0},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.5, 5.5},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.999, 9.5},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 10.0},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.1, 1.5},
		} {

			s := NewSampleWithValue(c.in, false)
			g.Expect(s.Percentile(c.p)).To(Equal(c.out))
		}
	})

	t.Run("Min", func(t *testing.T) {
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
			s := NewSampleWithValue(c.in, false)
			min, max := s.Bounds()
			g.Expect(min).To(Equal(c.min))
			g.Expect(max).To(Equal(c.max))

		}
	})

	t.Run("StdDev", func(t *testing.T) {
		g := NewGomegaWithT(t)
		s := NewSampleWithValue([]float64{1, 2, 3}, false)
		g.Expect(s.StandardDeviationPopulation()).To(BeNumerically("~", 1.0, 2))

		s = NewSampleWithValue([]float64{-1, -2, -3.0}, false)
		g.Expect(s.StandardDeviationPopulation()).To(BeNumerically("~", 1.15, 2))

		s = NewSampleWithValue([]float64{}, false)
		g.Expect(math.IsNaN(s.StandardDeviationPopulation())).To(BeTrue())
	})

	//t.Run("", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//})
	//
	//t.Run("", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//})

}

func TestMyStats2(t *testing.T) {

	t.Run("Mean", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			out float64
		}{
			{[]float64{1, 2, 3, 4, 5}, 3.0},
			{[]float64{1, 2, 3, 4, 5, 6}, 3.5},
			{[]float64{1}, 1.0},
		} {
			s := NewSampleFromPool(false)
			for _, v := range c.in {
				s.Append(v)

			}
			g.Expect(s.Mean()).To(Equal(c.out))
			s.BackToPool()
		}
	})

	t.Run("Variance Sample", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			out float64
		}{
			{[]float64{1, 2, 3, 4, 5}, 2.5},
			{[]float64{1, 2, 3, 4, 5, 6}, 3.5},
			{[]float64{1, 2}, 0.5},
		} {
			s := NewSampleFromPool(false)
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.Variance()).To(Equal(c.out))
			s.BackToPool()
		}
	})

	t.Run("Variance Population", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			out float64
		}{
			{[]float64{1, 2, 3, 4, 5}, 2},
			{[]float64{1, 2, 3, 4, 5, 6}, 2.9166666666667},
			{[]float64{1, 2}, 0.25},
		} {
			s := NewSampleFromPool(false)
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.PopulationVariance()).To(BeNumerically("~", c.out, 0.000001))
			s.BackToPool()
		}
	})

	t.Run("With Original", func(t *testing.T) {
		g := NewGomegaWithT(t)
		input := []float64{3, 2, 1}
		s := NewSampleFromPool(true)
		s.Append(3)
		s.Append(2)
		s.Append(1)

		g.Expect(s.Xs()[:3]).To(Equal(input))
		g.Expect(s.Original()[:3]).To(Equal(input))
		s.sort()
		g.Expect(s.Xs()[:3]).To(Equal([]float64{1, 2, 3}))
		g.Expect(s.Original()[:3]).To(Equal([]float64{3, 2, 1}))
		s.BackToPool()
	})

	t.Run("CoVariance Population", func(t *testing.T) {
		g := NewGomegaWithT(t)
		i1 := []float64{1, 2, 3.5, 3.7, 8, 12}
		i2 := []float64{10, -51.2, 8}
		i3 := []float64{0.5, 1, 2.1, 3.4, 3.4, 4}

		s1 := NewSampleFromPool(true)
		s2 := NewSampleFromPool(true)
		s3 := NewSampleFromPool(true)
		s4 := NewSampleFromPool(true)

		for _, v := range i1 {
			s1.Append(v)
		}
		for _, v := range i2 {
			s2.Append(v)
		}
		for _, v := range i3 {
			s3.Append(v)
		}

		g.Expect(math.IsNaN(s1.CovariancePopulation(s2))).To(BeTrue())
		g.Expect(s1.CovariancePopulation(s3)).To(Equal(4.191666666666666))
		g.Expect(math.IsNaN(s1.CovariancePopulation(s4))).To(BeTrue())
		s1.BackToPool()
		s2.BackToPool()
		s3.BackToPool()
		s4.BackToPool()
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
			s := NewSampleFromPool(false)
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.Sum()).To(Equal(c.out))
			s.BackToPool()
		}
	})

	t.Run("Percentile", func(t *testing.T) {
		g := NewGomegaWithT(t)
		for _, c := range []struct {
			in  []float64
			p   float64
			out float64
		}{
			{[]float64{43, 54, 56, 61, 62, 66}, 0.9, 64.0},
			{[]float64{43}, 0.9, 43.0},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.5, 5.5},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.999, 9.5},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 10.0},
			{[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.1, 1.5},
		} {

			s := NewSampleFromPool(false)
			for _, v := range c.in {
				s.Append(v)
			}
			g.Expect(s.Percentile(c.p)).To(Equal(c.out))
			s.BackToPool()
		}
	})

	t.Run("Min", func(t *testing.T) {
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
			s := NewSampleFromPool(false)
			for _, v := range c.in {
				s.Append(v)
			}
			min, max := s.Bounds()
			g.Expect(min).To(Equal(c.min))
			g.Expect(max).To(Equal(c.max))
			s.BackToPool()
		}
	})

	t.Run("StdDev", func(t *testing.T) {
		g := NewGomegaWithT(t)
		s := NewSampleFromPool(false)
		for _, v := range []float64{1, 2, 3} {
			s.Append(v)
		}
		g.Expect(s.StandardDeviationPopulation()).To(BeNumerically("~", 1.0, 2))
		s.BackToPool()

		s = NewSampleFromPool(false)
		for _, v := range []float64{-1, -2, -3.0} {
			s.Append(v)
		}
		g.Expect(s.StandardDeviationPopulation()).To(BeNumerically("~", 1.15, 2))
		s.BackToPool()

		s = NewSampleFromPool(false)
		g.Expect(math.IsNaN(s.StandardDeviationPopulation())).To(BeTrue())
		s.BackToPool()
	})

	//t.Run("", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//})
	//
	//t.Run("", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//})

}

var Sink float64

func BenchmarkSquare(b *testing.B) {

	b.Run("a*a", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			Sink = 159.45 * 159.45
		}
	})

	b.Run("mat.Pow", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			Sink = math.Pow(159.45, 2)
		}
	})
}

func BenchmarkSample(b *testing.B) {

	nbElements := 10 * 52 * 100
	arr := make([]float64, nbElements)
	for i := 0; i < nbElements; i++ {
		arr[i] = float64(i)
	}

	nbIterations := 500
	concurrency := 20

	b.ResetTimer()

	b.Run("With New", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {

			ch := make(chan int, 5)
			go func() {
				for i := 0; i < nbIterations; i++ {
					ch <- i
				}
				close(ch)
			}()
			var wg sync.WaitGroup
			wg.Add(concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if _, more := <-ch; more {
							sample := NewSampleWithValue(arr, false)
							sample.Mean()
							sample.StdDev()
						} else {
							wg.Done()
							return
						}
					}
				}()
			}

			wg.Wait()

		}
	})

	b.Run("With New Alternative", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {

			ch := make(chan int, 5)
			go func() {
				for i := 0; i < nbIterations; i++ {
					ch <- i
				}
				close(ch)
			}()
			var wg sync.WaitGroup
			wg.Add(concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if _, more := <-ch; more {
							sample := NewSampleFromPoolWithValues(arr, false)
							sample.Mean()
							sample.StdDev()
							samplePool.Put(sample)
						} else {
							wg.Done()
							return
						}
					}
				}()
			}

			wg.Wait()

		}
	})

	b.Run("With Pool Append", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			ch := make(chan int, 5)
			go func() {
				for i := 0; i < nbIterations; i++ {
					ch <- i
				}
				close(ch)
			}()
			var wg sync.WaitGroup
			wg.Add(concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if _, more := <-ch; more {
							sample := NewSampleFromPool(false)
							for _, v := range arr {
								sample.Append(v)
							}
							sample.Mean()
							sample.StdDev()
							sample.BackToPool()
						} else {
							wg.Done()
							return
						}
					}
				}()
			}
			wg.Wait()
		}
	})

	b.Run("With Pool Set", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {

			ch := make(chan int, 5)
			go func() {
				for i := 0; i < nbIterations; i++ {
					ch <- i
				}
				close(ch)
			}()
			var wg sync.WaitGroup
			wg.Add(concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if _, more := <-ch; more {
							sample := samplePool.Get().(*Sample)
							sample.preInit(len(arr), false)

							for _, v := range arr {
								sample.Append(v)
							}
							sample.Mean()
							sample.StdDev()
							samplePool.Put(sample)
						} else {
							wg.Done()
							return
						}
					}
				}()
			}

			wg.Wait()

		}
	})

	b.Run("SampleStream", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {

			ch := make(chan int, 5)
			go func() {
				for i := 0; i < nbIterations; i++ {
					ch <- i
				}
				close(ch)
			}()
			var wg sync.WaitGroup
			wg.Add(concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if _, more := <-ch; more {
							sample := NewSampleStream()
							for _, v := range arr {
								sample.Append(v)
							}
							sample.Mean()
							sample.StdDev()
							sample.CleanPool()
						} else {
							wg.Done()
							return
						}
					}
				}()
			}

			wg.Wait()

		}
	})

}

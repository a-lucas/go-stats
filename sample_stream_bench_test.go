package gostats

import (
	"github.com/a-lucas/go-stats/stats"
	"sync"
	"testing"
)

func BenchmarkSampleStream300Points(b *testing.B) {

	points := make([]float64, 0)
	for i := 0; i < 300; i++ {
		points = append(points, float64(i))
	}

	b.Run("COmparing With Native internal stats package", func(b *testing.B) {
		b.Run("Mean()", func(b *testing.B) {
			b.Run("Go Native", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					stats.Mean(points)
				}
			})
			b.Run("Go-Stats / memory Pooled", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					s := NewSampleFromPoolWithValues(points, false)
					s.Mean()
					s.BackToPool()
				}
			})
			b.Run("Go-Stats ", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					s := NewSampleWithValue(points, false)
					s.Mean()
				}
			})
			b.Run("Go-Stats - With Stream", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					s := NewSampleStream()
					s.AppendMany(points)
					s.Mean()
				}
			})
		})

		b.Run("StdDev()", func(b *testing.B) {
			b.Run("Go Native", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					stats.StdDev(points)
				}
			})
			b.Run("Go-Stats / memory Pooled", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					s := NewSampleFromPoolWithValues(points, false)
					s.StdDev()
					s.BackToPool()
				}
			})
			b.Run("Go-Stats ", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					s := NewSampleWithValue(points, false)
					s.StdDev()
				}
			})
			b.Run("Go-Stats - With Stream", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					s := NewSampleStream()
					s.AppendMany(points)
					s.StdDev()
				}
			})
		})

	})
	b.Run("ONE  go-Routine", func(b *testing.B) {
		b.Run("AppendMany ", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				st := NewSampleStream()
				st.AppendMany(points)
				st.Mean()
				st.StdDev()
				st.CleanPool()
			}
		})

		b.Run("Append in for loop", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				st := NewSampleStream()
				for _, val := range points {
					st.Append(val)
				}
				st.Mean()
				st.StdDev()
				st.CleanPool()
			}
		})

		b.Run("Vs Sample Append", func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				sample := NewSampleFromPool(false)
				for _, val := range points {
					sample.Append(val)
				}
				sample.Mean()
				sample.StdDev()
				sample.BackToPool()
			}
		})

	})

	b.Run("Sample Stream Multiple Routine - process 1000 times the 300 points", func(b *testing.B) {

		launchWithConcurrency := func(concurrency int) {
			ch := make(chan []float64, 30)
			go func() {
				for i := 0; i < 1000; i++ {
					ch <- points
				}
				close(ch)
			}()

			var wg sync.WaitGroup
			wg.Add(concurrency)
			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if values, more := <-ch; more {
							st := NewSampleStream()
							st.AppendMany(values)
							st.Mean()
							st.StdDev()
							st.CleanPool()
						} else {
							wg.Done()
							return
						}
					}
				}()
			}
			wg.Wait()
		}

		b.Run("Concurrency 5", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(10)
			}
		})

		b.Run("Concurrency 10", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(10)
			}
		})

		b.Run("Concurrency 40", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(40)
			}
		})
		b.Run("Concurrency 100", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(100)
			}
		})

	})

	b.Run("Sample Multiple Routine - process 1000 times the 300 points ", func(b *testing.B) {

		launchWithConcurrency := func(concurrency int) {
			ch := make(chan []float64, 30)
			go func() {
				for i := 0; i < 1000; i++ {
					ch <- points
				}
				close(ch)
			}()

			var wg sync.WaitGroup
			wg.Add(concurrency)
			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						if values, more := <-ch; more {
							sample := NewSampleFromPool(false)
							for _, val := range values {
								sample.Append(val)
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

		b.Run("Concurrency 5", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(10)
			}
		})

		b.Run("Concurrency 10", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(10)
			}
		})

		b.Run("Concurrency 40", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(40)
			}
		})
		b.Run("Concurrency 100", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				launchWithConcurrency(100)
			}
		})

	})

}

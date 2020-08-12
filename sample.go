// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gostats

import (
	"math"
	"sort"
	"sync"
)

type DistanceType int

const (
	DistanceChebyshev DistanceType = 1
	DistanceEuclidean DistanceType = 2
	DistanceManhattan DistanceType = 3
)

var samplePool = sync.Pool{
	New: func() interface{} {
		return &Sample{
			original: make([]float64, 0),
			xs:       make([]float64, 0),
		}
	},
}

// TODO Implement a STREAM for variance like here:
// https://math.stackexchange.com/questions/20593/calculate-variance-from-a-stream-of-sample-values

/**
List of method that needs sorting:
- Bounds
- Percentile
*/

// Sample is a collection of data points.
// Sample allows the fast processing of mean, Variance
// It is also possible to get Correlation and the Distance with another Sample

type Sample struct {
	// xs is the sorted slice of sample values.
	xs []float64
	// original is the original order
	original []float64
	// sorted indicates that xs is sorted in ascending order.
	sorted       bool
	withOriginal bool
	// internal values
	mean               float64
	geoMean            float64
	variance           float64
	sum                float64
	stdDev             float64
	populationVariance float64
	stdDevPopulation   float64
	min                float64
	max                float64
	nbPositives        int16
	nbProfitStreak     int16
	nbLossStreak       int16
	// more than internal
	_nbPositiveStreak int16
	_nbNegativeStreak int16
	nb                int
}

func NewSampleWithValue(inputs []float64, withOriginal bool) *Sample {
	s := &Sample{
		xs:                 make([]float64, len(inputs)),
		mean:               math.NaN(),
		variance:           math.NaN(),
		sum:                math.NaN(),
		stdDev:             math.NaN(),
		populationVariance: math.NaN(),
		stdDevPopulation:   math.NaN(),
		withOriginal:       withOriginal,
		sorted:             false,
		min:                math.Inf(+1),
		max:                math.Inf(-1),
		nbProfitStreak:     0,
		nbLossStreak:       0,
		nb:                 0,
	}

	copy(s.xs, inputs)
	s.processInit()
	if withOriginal {
		s.original = make([]float64, len(inputs))
		copy(s.original, inputs)
	}
	return s
}

func NewSampleFromPoolWithValues(inputs []float64, withOriginal bool) *Sample {
	s := samplePool.Get().(*Sample)
	s.initEmptyValues(withOriginal)
	s.xs = make([]float64, len(inputs))
	copy(s.xs, inputs)
	if withOriginal {
		s.original = make([]float64, len(inputs))
		copy(s.original, inputs)
	}

	s.processInit()
	return s
}

// Uses a memory Pool (faster)
func NewSampleFromPool(withOriginal bool) *Sample {
	s := samplePool.Get().(*Sample)
	s.initEmptyValues(withOriginal)
	return s
}

// Number of positive values in the dataset
func (s *Sample) NbPositive() int16 {
	return s.nbPositives
}

// maximum number of consecutive positive values in the dataset
func (s *Sample) NbPositiveStreak() int16 {
	return s._nbPositiveStreak
}

// maximum number of consecutive negative values in the dataset
func (s *Sample) NbNegativeStreak() int16 {
	return s._nbNegativeStreak
}

func (s *Sample) Xs() []float64 {
	return s.xs
}

// Put the Sample back to the memory pool
func (s *Sample) BackToPool() {
	samplePool.Put(s)
}

func (s *Sample) processInit() {
	for _, num := range s.xs {
		s.walk(num)
	}
}

func (s *Sample) preInit(dataLength int, withOriginal bool) {
	s.initEmptyValues(withOriginal)
	s.xs = make([]float64, dataLength)
	if withOriginal {
		s.original = make([]float64, dataLength)
	}
}

func (s *Sample) initEmptyValues(withOriginal bool) {
	s.mean = math.NaN()
	s.variance = math.NaN()
	s.sum = math.NaN()
	s.stdDev = math.NaN()
	s.populationVariance = math.NaN()
	s.stdDevPopulation = math.NaN()
	s.sorted = false
	s.nbPositives = 0
	s.min = math.Inf(+1)
	s.max = math.Inf(-1)
	s.sum = 0
	s._nbNegativeStreak = 0
	s._nbPositiveStreak = 0
	s.nb = 0
	s.xs = nil
	s.original = nil
	s.geoMean = math.NaN()
	s.nbLossStreak = 0
	s.nbProfitStreak = 0
	s.withOriginal = withOriginal
}

func (s *Sample) walk(value float64) {
	s.sum += value
	if value > 0 {
		s.nbPositives++
		s._nbPositiveStreak++
		if s._nbPositiveStreak > s.nbProfitStreak {
			s.nbProfitStreak = s._nbPositiveStreak
		}
		s._nbNegativeStreak = 0
	} else {
		s._nbNegativeStreak++
		if s._nbNegativeStreak > s.nbLossStreak {
			s.nbLossStreak = s._nbNegativeStreak
		}
		s._nbPositiveStreak = 0
	}
	if value < s.min {
		s.min = value
	}
	if value > s.max {
		s.max = value
	}
	s.nb++
}

// Append a value to the sample.
func (s *Sample) Append(value float64) {
	s.xs = append(s.xs, value)
	if s.withOriginal {
		s.original = append(s.original, value)
	}
	s.walk(value)
}

func (s *Sample) Original() []float64 {
	if !s.withOriginal {
		panic("only works when withOriginal == true")
	}
	return s.original
}

// Bounds returns the minimum and maximum values of the Sample.
//
// If the Sample is weighted, this ignores samples with zero weight.
//
// This is constant time if s.sorted and there are no zero-weighted
// values.
func (s *Sample) Bounds() (min float64, max float64) {
	if s.nb == 0 {
		return math.NaN(), math.NaN()
	}
	return s.min, s.max
}

// Sum returns the sum of the Sample.
func (s *Sample) Sum() float64 {
	if !math.IsNaN(s.sum) {
		return s.sum
	}
	sum := 0.0
	for i := 0; i < s.nb; i++ {
		sum += s.xs[i]
	}
	s.sum = sum
	return sum
}

// Mean returns the arithmetic mean of the Sample.
func (s *Sample) Mean() float64 {
	if !math.IsNaN(s.mean) {
		return s.mean
	}
	if s.nb == 0 {
		return math.NaN()
	}
	s.mean = s.Sum() / float64(s.nb)
	return s.mean
}

func (s *Sample) Len() int {
	return s.nb
}

// Calculate the Distance ( Chebyshev - buggy / Euclidean / Manhattan ) with another Sample
func (s *Sample) Distance(data *Sample, distanceType DistanceType) float64 {
	if !s.validateAgainst(data) {
		return math.NaN()
	}
	switch distanceType {
	case DistanceChebyshev:
		panic("Definitively buggy")
		var tempDistance float64
		distance := 0.0
		for i := 0; i < s.nb; i++ {
			tempDistance = math.Abs(s.Original()[i] - data.Original()[i])

			if distance < tempDistance {
				distance = tempDistance
			}
		}
		return distance
	case DistanceEuclidean:
		distance := 0.0
		for i := 0; i < s.nb; i++ {
			distance = distance + math.Pow(s.Original()[i]-data.Original()[i], 2)
		}
		return distance / math.Pow(float64(s.nb), 2)
	case DistanceManhattan:
		distance := 0.0
		for i := 0; i < s.nb; i++ {
			distance = distance + math.Abs(s.Original()[i]-data.Original()[i])
		}
		return distance / float64(s.nb)
	default:
		panic("Unknown distance")
	}
}

/**
returns true if valid
*/
func (s *Sample) validateAgainst(data *Sample) bool {
	if s.nb == 0 || data.Len() == 0 {
		return false
	}
	if s.nb != data.Len() {
		return false
	}
	return true
}

// CovariancePopulation computes covariance for entire population between two variables.
// https://corporatefinanceinstitute.com/resources/knowledge/finance/covariance/
func (s *Sample) CovariancePopulation(data *Sample) float64 {

	if !s.withOriginal {
		panic("CovariancePopulation needs original Access")
	}
	if !data.withOriginal {
		panic("CovariancePopulation needs original Access")
	}

	if !s.validateAgainst(data) {
		return math.NaN()
	}

	l1 := data.Len()

	m1 := data.Mean()
	m2 := s.Mean()

	ss := 0.0
	for i := 0; i < l1; i++ {
		delta1 := data.original[i] - m1
		delta2 := s.original[i] - m2
		ss += delta1 * delta2
	}
	return ss / float64(l1)
}

// Calculate the Pearson Correlation with another Sample
func (s *Sample) Correlation(data *Sample, correlationType CorrelationType) float64 {
	switch correlationType {
	case CorrelationPearson:
		return s.correlationPearson(data)
	default:
		panic("Correlation not implemented")
	}
}

// Correlation describes the degree of relationship between two sets of data
func (s *Sample) correlationPearson(data *Sample) float64 {
	l1 := data.Len()
	l2 := s.nb
	if l1 == 0 || l2 == 0 {
		return math.NaN()
	}
	if l1 != l2 {
		return math.NaN()
	}

	dev1 := data.StandardDeviationPopulation()
	dev2 := s.StandardDeviationPopulation()

	if dev1 == 0 || dev2 == 0 {
		return 0
	}

	return s.CovariancePopulation(data) / (dev1 * dev2)
}

// StandardDeviationPopulation finds the amount of variation from the population
// σ = ([Σ(x - u)2]/N)1/2
func (s *Sample) StandardDeviationPopulation() float64 {
	if s.nb == 0 {
		return math.NaN()
	}
	if !math.IsNaN(s.stdDevPopulation) {
		return s.stdDevPopulation
	}
	s.stdDevPopulation = math.Pow(s.PopulationVariance(), 0.5)
	return s.stdDevPopulation
}

// PopulationVariance finds the amount of variance within a population
// https://www.wallstreetmojo.com/population-variance-formula/
func (s *Sample) PopulationVariance() float64 {
	if !math.IsNaN(s.populationVariance) || s.nb == 0 {
		return s.populationVariance
	}
	variance := 0.0
	m := s.Mean()
	for i := 0; i < s.nb; i++ {
		variance += (s.xs[i] - m) * (s.xs[i] - m)
	}
	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and wether to subtract by one or not
	s.populationVariance = variance / float64(s.nb)
	return s.populationVariance
}

func (s *Sample) Variance() float64 {
	if !math.IsNaN(s.variance) || s.nb == 0 {
		return s.variance
	}
	variance := 0.0
	m := s.Mean()
	for i := 0; i < s.nb; i++ {
		variance += (s.xs[i] - m) * (s.xs[i] - m)
	}
	s.variance = variance / float64(s.nb-1)
	return s.variance
}

// StdDev returns the sample standard deviation of the Sample.
func (s *Sample) StdDev() float64 {
	if !math.IsNaN(s.stdDev) {
		return s.stdDev
	}
	if s.nb == 0 {
		return math.NaN()
	}
	if s.nb == 1 {
		s.stdDev = 0
		return s.stdDev
	}
	s.stdDev = math.Pow(s.Variance(), 0.5)
	return s.stdDev
}

// Get the percentile - 0<value<1 is
func (s *Sample) Percentile(value float64) float64 {
	if s.nb == 0 {
		return math.NaN()
	}
	s.sort()

	if value <= 0 {
		min, _ := s.Bounds()
		return min
	} else if value >= 1 {
		_, max := s.Bounds()
		return max
	}

	l := s.nb

	index := value * (float64(l) - 1)

	// Check if the index is a whole number
	if index == float64(int(index)) {

		// Convert float to int
		i := int(index)

		// Find the value at the index
		return s.xs[i]

	} else {
		return (s.xs[int(math.Floor(index))] + s.xs[int(math.Ceil(index))]) / 2.0
	}
}

// Sort sorts the samples in place in s and returns s.
//
// A sorted sample improves the performance of some algorithms.
func (s *Sample) sort() {
	sub := s.xs[:s.nb]
	sort.Float64s(sub)
	s.xs = append(sub, s.xs[s.nb:]...)
	s.sorted = true
}

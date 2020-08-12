// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gostats

import (
	"math"
	"sync"
)

var SampleStreamPool = sync.Pool{
	New: func() interface{} { return new(SampleStream) },
}

// SampleStream is a simplified Sample implementation that process the Mean & StdDev calculation in a stream, giving much faster performances (x3) than `Sample` when used in a real time environment.
type SampleStream struct {
	// internal values
	sum            float64
	min            float64
	max            float64
	nbPositives    int16
	nbProfitStreak int16
	nbLossStreak   int16
	// more than internal
	_nbPositiveStreak int16
	_nbNegativeStreak int16
	// streaming
	_initialized bool
	_prevMean    float64
	//_absPrevMean    float64
	_prevVariance float64
	//_absPreVariance float64
	_count float64
}

func NewSampleStream() *SampleStream {
	s := SampleStreamPool.Get().(*SampleStream)
	s.initEmptyValues()
	return s
}

func (s *SampleStream) CleanPool() {
	SampleStreamPool.Put(s)
}

func (s *SampleStream) initEmptyValues() {
	s.nbPositives = 0
	s.min = math.Inf(+1)
	s.max = math.Inf(-1)
	s.sum = 0
	s._nbNegativeStreak = 0
	s._nbPositiveStreak = 0
	s._initialized = false
	s._prevVariance = 0
	s._prevMean = 0
	s._count = 0
	s.nbProfitStreak = 0
	s.nbLossStreak = 0
}

func (s *SampleStream) AppendMany(values []float64) {
	for _, value := range values {
		s.Append(value)
	}
}

func (s *SampleStream) Append(value float64) {

	s.sum = s.sum + value
	if value > 0 {
		s.nbPositives++
		s._nbPositiveStreak++
		s._nbNegativeStreak = 0
		if s._nbPositiveStreak > s.nbProfitStreak {
			s.nbProfitStreak = s._nbPositiveStreak
		}
	} else {
		s._nbNegativeStreak++
		s._nbPositiveStreak = 0
		if s._nbNegativeStreak > s.nbLossStreak {
			s.nbLossStreak = s._nbNegativeStreak
		}
	}

	s._count = s._count + 1
	if !s._initialized {
		s._initialized = true
		s._prevMean = value
		s.min = value
		s.max = value
	} else {
		if value < s.min {
			s.min = value
		} else if value > s.max {
			s.max = value
		}
		// mean
		_prevMean := s._prevMean + (value-s._prevMean)/s._count
		s._prevVariance = s._prevVariance + (value-s._prevMean)*(value-_prevMean)
		s._prevMean = _prevMean
	}

}

// Bounds returns the minimum and maximum values of the Sample.
//
// If the Sample is weighted, this ignores samples with zero weight.
//
// This is constant time if s.sorted and there are no zero-weighted
// values.
func (s *SampleStream) Bounds() (min float64, max float64) {
	return s.min, s.max
}

// Sum returns the (possibly weighted) sum of the Sample.
func (s *SampleStream) Sum() float64 {
	return s.sum
}

// Mean returns the arithmetic mean of the Sample.
func (s *SampleStream) Mean() float64 {
	return s.sum / s._count
}

func (s *SampleStream) Variance() float64 {
	return s._prevVariance / (s._count - 1)
}

// StdDev returns the sample standard deviation of the Sample.
func (s *SampleStream) StdDev() float64 {
	return math.Pow(s.Variance(), 0.5)
}

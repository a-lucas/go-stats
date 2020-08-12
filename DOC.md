# gostats
--
    import "."


## Usage

```go
var SampleStreamPool = sync.Pool{
	New: func() interface{} { return new(SampleStream) },
}
```

#### type CorrelationType

```go
type CorrelationType string
```


```go
const CorrelationCustom CorrelationType = "Custom"
```

```go
const CorrelationKendall CorrelationType = "Kendal"
```

```go
const CorrelationPearson CorrelationType = "Pearson"
```

```go
const CorrelationSpearman CorrelationType = "SpearMan"
```

#### func (CorrelationType) String

```go
func (ct CorrelationType) String() string
```

#### type DistanceType

```go
type DistanceType int
```


```go
const (
	DistanceChebyshev DistanceType = 1
	DistanceEuclidean DistanceType = 2
	DistanceManhattan DistanceType = 3
)
```

#### type Sample

```go
type Sample struct {
}
```


#### func  NewSampleFromPool

```go
func NewSampleFromPool(withOriginal bool) *Sample
```
Uses a memory Pool (faster)

#### func  NewSampleFromPoolWithValues

```go
func NewSampleFromPoolWithValues(inputs []float64, withOriginal bool) *Sample
```

#### func  NewSampleWithValue

```go
func NewSampleWithValue(inputs []float64, withOriginal bool) *Sample
```

#### func (*Sample) Append

```go
func (s *Sample) Append(value float64)
```
Append a value to the sample.

#### func (*Sample) BackToPool

```go
func (s *Sample) BackToPool()
```
Put the Sample back to the memory pool

#### func (*Sample) Bounds

```go
func (s *Sample) Bounds() (min float64, max float64)
```
Bounds returns the minimum and maximum values of the Sample.

If the Sample is weighted, this ignores samples with zero weight.

This is constant time if s.sorted and there are no zero-weighted values.

#### func (*Sample) Correlation

```go
func (s *Sample) Correlation(data *Sample, correlationType CorrelationType) float64
```
Calculate the Pearson Correlation with another Sample

#### func (*Sample) CovariancePopulation

```go
func (s *Sample) CovariancePopulation(data *Sample) float64
```
CovariancePopulation computes covariance for entire population between two
variables.
https://corporatefinanceinstitute.com/resources/knowledge/finance/covariance/

#### func (*Sample) Distance

```go
func (s *Sample) Distance(data *Sample, distanceType DistanceType) float64
```
Calculate the Distance ( Chebyshev - buggy / Euclidean / Manhattan ) with
another Sample

#### func (*Sample) Len

```go
func (s *Sample) Len() int
```

#### func (*Sample) Mean

```go
func (s *Sample) Mean() float64
```
Mean returns the arithmetic mean of the Sample.

#### func (*Sample) NbNegativeStreak

```go
func (s *Sample) NbNegativeStreak() int16
```
maximum number of consecutive negative values in the dataset

#### func (*Sample) NbPositive

```go
func (s *Sample) NbPositive() int16
```
Number of positive values in the dataset

#### func (*Sample) NbPositiveStreak

```go
func (s *Sample) NbPositiveStreak() int16
```
maximum number of consecutive positive values in the dataset

#### func (*Sample) Original

```go
func (s *Sample) Original() []float64
```

#### func (*Sample) Percentile

```go
func (s *Sample) Percentile(value float64) float64
```
Get the percentile - 0<value<1 is

#### func (*Sample) PopulationVariance

```go
func (s *Sample) PopulationVariance() float64
```
PopulationVariance finds the amount of variance within a population
https://www.wallstreetmojo.com/population-variance-formula/

#### func (*Sample) StandardDeviationPopulation

```go
func (s *Sample) StandardDeviationPopulation() float64
```
StandardDeviationPopulation finds the amount of variation from the population σ
= ([Σ(x - u)2]/N)1/2

#### func (*Sample) StdDev

```go
func (s *Sample) StdDev() float64
```
StdDev returns the sample standard deviation of the Sample.

#### func (*Sample) Sum

```go
func (s *Sample) Sum() float64
```
Sum returns the sum of the Sample.

#### func (*Sample) Variance

```go
func (s *Sample) Variance() float64
```

#### func (*Sample) Xs

```go
func (s *Sample) Xs() []float64
```

#### type SampleStream

```go
type SampleStream struct {
}
```

SampleStream is a simplified Sample implementation that process the Mean &
StdDev calculation in a stream, giving much faster performances (x3) than
`Sample` when used in a real time environment.

#### func  NewSampleStream

```go
func NewSampleStream() *SampleStream
```

#### func (*SampleStream) Append

```go
func (s *SampleStream) Append(value float64)
```

#### func (*SampleStream) AppendMany

```go
func (s *SampleStream) AppendMany(values []float64)
```

#### func (*SampleStream) Bounds

```go
func (s *SampleStream) Bounds() (min float64, max float64)
```
Bounds returns the minimum and maximum values of the Sample.

If the Sample is weighted, this ignores samples with zero weight.

This is constant time if s.sorted and there are no zero-weighted values.

#### func (*SampleStream) CleanPool

```go
func (s *SampleStream) CleanPool()
```

#### func (*SampleStream) Mean

```go
func (s *SampleStream) Mean() float64
```
Mean returns the arithmetic mean of the Sample.

#### func (*SampleStream) StdDev

```go
func (s *SampleStream) StdDev() float64
```
StdDev returns the sample standard deviation of the Sample.

#### func (*SampleStream) Sum

```go
func (s *SampleStream) Sum() float64
```
Sum returns the (possibly weighted) sum of the Sample.

#### func (*SampleStream) Variance

```go
func (s *SampleStream) Variance() float64
```

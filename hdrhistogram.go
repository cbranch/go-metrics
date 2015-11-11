package metrics

import (
	"github.com/cloudflare/hdrhistogram"
)

// NewHdrHistogram constructs a new HdrHistogram for the given range of values.
func NewHdrHistogram(minValue, maxValue int64, sigfigs int) Histogram {
	if UseNilHists {
		return NilHistogram{}
	}
	return &HdrHistogram{
		hist:     hdrhistogram.NewAtomic(1, maxValue-minValue, sigfigs),
		minValue: minValue,
	}
}

// HistogramSnapshot is a read-only copy of another Histogram.
type HdrHistogramSnapshot struct {
	sample   *hdrhistogram.Histogram
	minValue int64
}

// Clear panics.
func (*HdrHistogramSnapshot) Clear() {
	panic("Clear called on a HdrHistogramSnapshot")
}

// Count returns the number of samples recorded at the time the snapshot was
// taken.
func (h *HdrHistogramSnapshot) Count() int64 { return h.sample.TotalCount() }

// Max returns the maximum value in the sample at the time the snapshot was
// taken.
func (h *HdrHistogramSnapshot) Max() int64 { return h.sample.Max() + h.minValue }

// Mean returns the mean of the values in the sample at the time the snapshot
// was taken.
func (h *HdrHistogramSnapshot) Mean() float64 {
	if h.Count() == 0 {
		return 0
	}
	return float64(h.Sum()) / float64(h.Count())
}

// Min returns the minimum value in the sample at the time the snapshot was
// taken.
func (h *HdrHistogramSnapshot) Min() int64 { return h.sample.Min() + h.minValue }

// Percentile returns an arbitrary percentile of values in the sample at the
// time the snapshot was taken.
func (h *HdrHistogramSnapshot) Percentile(p float64) float64 {
	return float64(h.sample.ValueAtQuantile(p*100) + h.minValue)
}

// Percentiles returns a slice of arbitrary percentiles of values in the sample
// at the time the snapshot was taken.
func (h *HdrHistogramSnapshot) Percentiles(ps []float64) []float64 {
	percentiles := make([]float64, len(ps))
	for i, p := range ps {
		percentiles[i] = h.Percentile(p)
	}
	return percentiles
}

// Sample returns the Sample underlying the histogram.
func (h *HdrHistogramSnapshot) Sample() Sample { return NilSample{} }

// Snapshot returns the snapshot.
func (h *HdrHistogramSnapshot) Snapshot() Histogram { return h }

// StdDev returns the standard deviation of the values in the sample at the
// time the snapshot was taken.
func (h *HdrHistogramSnapshot) StdDev() float64 { return h.sample.StdDev() }

// Sum returns the sum in the sample at the time the snapshot was taken.
func (h *HdrHistogramSnapshot) Sum() int64 {
	return h.sample.Sum() + (h.minValue * h.sample.TotalCount())
}

// Update panics.
func (*HdrHistogramSnapshot) Update(int64) {
	panic("Update called on a HdrHistogramSnapshot")
}

// Variance returns the variance of inputs at the time the snapshot was taken.
// TODO
func (h *HdrHistogramSnapshot) Variance() float64 { return 0 }

type HdrHistogram struct {
	hist     *hdrhistogram.AtomicHistogram
	minValue int64
}

// Clear clears the histogram and its sample.
func (h *HdrHistogram) Clear() { h.hist.Reset() }

// Count returns the number of samples recorded since the histogram was last
// cleared.
func (h *HdrHistogram) Count() int64 { return h.hist.TotalCount() }

// Max returns the maximum value in the sample.
func (h *HdrHistogram) Max() int64 { return h.hist.Max() + h.minValue }

// Mean returns the mean of the values in the sample.
func (h *HdrHistogram) Mean() float64 {
	if h.Count() == 0 {
		return 0
	}
	return float64(h.Sum()) / float64(h.Count())
}

// Min returns the minimum value in the sample.
func (h *HdrHistogram) Min() int64 { return h.hist.Min() + h.minValue }

// Percentile returns an arbitrary percentile of the values in the sample.
func (h *HdrHistogram) Percentile(p float64) float64 {
	return float64(h.hist.ValueAtQuantile(p*100) + h.minValue)
}

// Percentiles returns a slice of arbitrary percentiles of the values in the
// sample.
func (h *HdrHistogram) Percentiles(ps []float64) []float64 {
	percentiles := make([]float64, len(ps))
	for i, p := range ps {
		percentiles[i] = h.Percentile(p)
	}
	return percentiles
}

// Sample returns the Sample underlying the histogram.
func (h *HdrHistogram) Sample() Sample { return NilSample{} }

// Snapshot returns a read-only copy of the histogram.
func (h *HdrHistogram) Snapshot() Histogram {
	return &HdrHistogramSnapshot{
		sample:   hdrhistogram.Import(h.hist.Export()),
		minValue: h.minValue,
	}
}

// StdDev returns the standard deviation of the values in the sample.
func (h *HdrHistogram) StdDev() float64 { return h.hist.StdDev() }

// Sum returns the sum in the sample.
func (h *HdrHistogram) Sum() int64 {
	return h.hist.Sum() + (h.minValue * h.hist.TotalCount())
}

// Update samples a new value.
func (h *HdrHistogram) Update(v int64) { h.hist.RecordValue(v - h.minValue) }

// Variance returns the variance of the values in the sample.
// TODO
func (h *HdrHistogram) Variance() float64 { return 0 }

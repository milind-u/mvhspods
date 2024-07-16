package tests

import (
	"fmt"

	"mvhspods"

	"github.com/milind-u/glog"
)

// An error is defined as the magnitude of difference between a pecent of a certain group
// (ex. Spanish, Female, Band) in a certain pod from the percent of that group in the full
// student population

const badErrThreshold = 0.1

type Stats struct {
	avgErr mvhspods.Percent
	maxErr mvhspods.Percent
	// Count of errors that are higher than badErrorThreshold
	badErrs int
}

func PodStats(pd *mvhspods.PodData) *Stats {
	return PodStatsWithTolerance(pd, 1)
}

func PodStatsWithTolerance(pd *mvhspods.PodData, errTolerance mvhspods.Percent) *Stats {
	stats := new(Stats)
	numErrs := len(pd.Population()) * len(pd.Pods())

	for _, pod := range pd.Pods() {
		podStats := PodStatsOfPod(pod, pd.Population())
		if podStats.maxErr > errTolerance {
			glog.Errorf("Max error of %v exceeds tolerance of %v",
				podStats.maxErr, errTolerance)
		}
		stats.avgErr += podStats.avgErr * mvhspods.Percent(len(pd.Population()))
		stats.badErrs += podStats.badErrs
		if podStats.maxErr > stats.maxErr {
			stats.maxErr = podStats.maxErr
		}
	}

	stats.avgErr /= mvhspods.Percent(numErrs)

	return stats
}

func PodStatsOfPod(pod mvhspods.Students, population mvhspods.Percents) *Stats {
	stats := new(Stats)

	podPercents := mvhspods.PercentsOf(pod)
	for field, actualPercent := range population {
		err := mvhspods.Abs(podPercents[field] - actualPercent)
		stats.avgErr += err
		if err > stats.maxErr {
			stats.maxErr = err
		}
		if err > badErrThreshold {
			stats.badErrs++
		}
	}
	stats.avgErr /= mvhspods.Percent(len(population))
	return stats
}

func (s *Stats) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *Stats) AvgErr() mvhspods.Percent {
	return s.avgErr
}

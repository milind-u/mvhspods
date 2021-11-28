package tests

import (
  "fmt"

  "mvhspods"
)

// An error is defined as the magnitude of difference between a pecent of a certain group
// (ex. Spanish, Female, Band) in a certain pod from the percent of that group in the full
// student population

const badErrThreshold = 0.2

type Stats struct {
  avgErr mvhspods.Percent
  maxErr mvhspods.Percent
  // Count of errors that are higher than badErrorThreshold
  badErrs int
}

func PodStats(pd *mvhspods.PodData) *Stats {
  stats := new(Stats)
  numErrs := len(pd.Population()) * len(pd.Pods())
  index := 0
  for _, pod := range pd.Pods() {
    podPercents := mvhspods.PercentsOf(pod)
    for field, actualPercent := range pd.Population() {
      err := mvhspods.Abs(podPercents[field] - actualPercent)
      index++
      stats.avgErr += err
      if err > stats.maxErr {
        stats.maxErr = err
      }
      if err > badErrThreshold {
        stats.badErrs++
      }
    }
  }

  stats.avgErr /= mvhspods.Percent(numErrs)

  return stats
}

func (s *Stats) String() string {
  return fmt.Sprintf("%+v", *s)
}

package tests

import (
  "math"
  "testing"

  "mvhspods"
)

func TestDiversity(t *testing.T) {
  const percentErrTolerance mvhspods.Percent = 0.37
  const avgPercentErrTolerance mvhspods.Percent = 0.12

  students := GenerateStudents(500)
  pm := mvhspods.PodManager{Headers: Headers, Students: students}
  pm.MakePods(mvhspods.DefaultPodSize, false)

  var avgErr mvhspods.Percent
  for i, pod := range pm.Pods() {
    if i != mvhspods.EldPod {
      percents := mvhspods.PercentsOf(pod)
      for field, percent := range pm.Population() {
        err := mvhspods.Percent(math.Abs(float64(percents[field] - percent)))
        avgErr += err
        if err > percentErrTolerance {
          t.Error("Percent error of", err,
            "between population percent", percent, "and pod percent", percents[field],
            ",\nfor field", field, ", exceeds tolerance of", percentErrTolerance,
            "\nPod number", i+1, "-", pod)
        }
      }
    }
  }

  avgErr /= mvhspods.Percent(len(pm.Population()) * len(pm.Pods()))
  t.Log("Average error:", avgErr)
  if avgErr > avgPercentErrTolerance {
    t.Error("Average error exceeds tolerance of", avgPercentErrTolerance)
  }
}

package tests

import (
  "testing"

  "mvhspods"
)

func TestWeight(t *testing.T) {
  population := mvhspods.Percents{{1, "Graham"}: 0.3,
    {4, "Male"}: 0.5, {7, "English"}: 0.7}
  pod := mvhspods.Percents{{1, "Graham"}: 0.05,
    {4, "Male"}: 0.7, {7, "English"}: 0.7}
  s := mvhspods.Student{"100012345", "Graham", "Bar", "Foo", "Male", "Parent", "6501231234",
    "English", ""}

  weight := s.Weight(population, pod)
  t.Log("Weight:", weight)
  const expectedWeight = (0.3 - 0.05) + (0.5 - 0.7) + (0.7 - 0.7)
  if mvhspods.Abs(expectedWeight-weight) > 1e-5 {
    t.Error("Weight does not match expected weight of", expectedWeight)
  }
}

func TestDiversity(t *testing.T) {
  const errTolerance mvhspods.Percent = 0.37
  const avgErrTolerance mvhspods.Percent = 0.1

  students := GenerateStudents(500)
  pm := mvhspods.PodManager{Headers: Headers, Students: students}
  pm.MakePods(mvhspods.DefaultPodSize, false)

  stats := PodStats(pm.Pods(), pm.Population())
  t.Log("Stats:", stats)

  if stats.maxErr > errTolerance {
    t.Error("Percent error exceeds tolerance of", errTolerance)
  }

  if stats.avgErr > avgErrTolerance {
    t.Error("Average error exceeds tolerance of", avgErrTolerance)
  }
}

package main

import (
  "testing"

  "mvhspods"
)

func TestPods(t *testing.T) {
  students := generateStudents(10000)
  pm := mvhspods.PodManager{
    Headers:  headers,
    Students: students,
  }
  pm.MakePods(defaultPodSize, false)


}

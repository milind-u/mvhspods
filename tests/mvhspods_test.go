package tests

import (
	"math"
	"testing"

	"mvhspods"
)

func TestPods(t *testing.T) {
	// TODO: for some reason one pod exceeds this and has ~0.37, and way too many blach kids
	const percentDiffTolerance mvhspods.Percent = 0.3

	students := GenerateStudents(100)
	pm := mvhspods.PodManager{
		Headers:  Headers,
		Students: students,
	}
	pm.MakePods(mvhspods.DefaultPodSize, false)
	for i, pod := range pm.Pods() {
		if i != mvhspods.EldPod {
			percents := mvhspods.PercentsOf(pod)
			for field, percent := range pm.Population() {
				if diff := math.Abs(float64(percents[field] - percent)); diff > float64(percentDiffTolerance) {
					t.Error("Percent difference of", diff,
						"between population percent", percent, "and pod percent", percents[field],
						",\nfor field", field, ", exceeds tolerance of", percentDiffTolerance,
						"\nPod number", i+1, "-", pod)
				}
			}
		}
	}
}

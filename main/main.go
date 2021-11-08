package main

import (
  "flag"

  "mvhspods"

  "github.com/milind-u/glog"
)

const defaultPodSize = 12

func main() {
  podSize := flag.Int("pod_size", defaultPodSize, "Number of students per pod")
  sorted := flag.Bool("sorted", false,
    "Whether to sort the students output in alphabetical order")
  sampleData := flag.Bool("sample_data", false,
    "If true, using sample data so trim the last column in the csv which has sample pod numbers")
  flag.Parse()

  glog.SetSeverity(glog.InfoSeverity)

  var pm mvhspods.PodManager
  pm.ReadStudents("students.csv", *sampleData)
  pm.MakePods(*podSize, *sorted)
  pm.WritePods("pods.csv")
}

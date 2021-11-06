package main

import (
  "flag"

  "mvhspods"

  "github.com/milind-u/glog"
)

func main() {
  sorted := flag.Bool("sorted", false,
    "Whether to sort the students output in alphabetical order")
  sampleData := flag.Bool("sample_data", false,
    "If true, using sample data so trim the last column in the csv which has sample pod numbers")
  flag.Parse()

  glog.SetSeverity(glog.InfoSeverity)

  var pm mvhspods.PodManager
  pm.ReadStudents("students.csv", *sampleData)
  pm.MakePods(*sorted)
  pm.WritePods("pods.csv")
}

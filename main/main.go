package main

import (
  "flag"

  "mvhspods"
  "tests"

  "github.com/milind-u/glog"
)

func main() {
  studentsToGenerate := flag.Int("students_to_generate", 0,
    "If set, will generate the given number of random students instead of making pods")
  podSize := flag.Int("pod_size", mvhspods.DefaultPodSize, "Number of students per pod")
  sorted := flag.Bool("sorted", false,
    "Whether to sort the students output in alphabetical order")
  sampleData := flag.Bool("sample_data", false,
    "If true, using sample data so trim the last column in the csv which has sample pod numbers")
  flag.Parse()

  glog.SetSeverity(glog.InfoSeverity)

  if *studentsToGenerate != 0 {
    flag.Parse()
    glog.SetSeverity(glog.InfoSeverity)

    students := tests.GenerateStudents(*studentsToGenerate)

    glog.Infoln("Percents: ", mvhspods.PercentsOf(students))
    mvhspods.WriteStudents("students.csv", tests.Headers, students)
  } else {
    var pm mvhspods.PodManager
    pm.ReadStudents("students.csv", *sampleData)
    pm.MakePods(*podSize, *sorted)
    glog.Infoln("Stats:", tests.PodStats(pm.Pods(), pm.Population()))
    pm.WritePods("pods.csv")
  }
}

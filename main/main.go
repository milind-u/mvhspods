package main

import (
  "flag"
  "syscall/js"

  "mvhspods"
  "tests"

  "github.com/milind-u/glog"
)

func webMain() js.Func {
  return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    pods := ""
    if len(args) == 1 {
      csv := args[0].String()

      var pm mvhspods.PodManager
      pm.ReadStudentsFromString(csv)
      pm.MakePods(mvhspods.DefaultPodSize)
      pods = pm.WritePodsToString()
    } else {
      glog.Errorln("Expected 1 arg (csv), but got", len(args))
    }
    return pods
  })
}

func main() {
  web := flag.Bool("web", true, "Whether to run the program for the webapp (if false, CLI)")
  studentsToGenerate := flag.Int("students_to_generate", 0,
    "If set, will generate the given number of random students instead of making pods")
  podSize := flag.Int("pod_size", mvhspods.DefaultPodSize, "Number of students per pod")
  sorted := flag.Bool("sorted", false,
    "Whether to sort the students output in alphabetical order")
  sampleData := flag.Bool("sample_data", false,
    "If true, using sample data so trim the last column in the csv which has sample pod numbers")
  flag.Parse()

  glog.SetSeverity(glog.InfoSeverity)

  if *web {
    js.Global().Set("makePods", webMain())
    // Keep the program running
    <-make(chan interface{})
  } else if *studentsToGenerate != 0 {
    flag.Parse()
    glog.SetSeverity(glog.InfoSeverity)

    students := tests.GenerateStudents(*studentsToGenerate)

    glog.Infoln("Percents: ", mvhspods.PercentsOf(students))
    mvhspods.WriteStudents("students.csv", tests.Headers, students)
  } else {
    var pm mvhspods.PodManager
    pm.ReadStudents("students.csv", *sampleData)
    pm.MakePods(*podSize)
    glog.Infoln("ELD stats:", tests.PodStats(&pm.Eld))
    glog.Infoln("Stats:", tests.PodStats(&pm.PodData))
    pm.WritePods("pods.csv", *sorted)
  }
}

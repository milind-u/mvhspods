package main

import (
  "mvhspods"

  "github.com/milind-u/mlog"
)

func main() {
  mlog.SetLevel(mlog.LInfo)

  var pm mvhspods.PodManager
  pm.ReadStudents("students.csv")
  pm.MakePods()
  pm.WritePods("pods.csv")
}

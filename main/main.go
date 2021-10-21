package main

import "mvhspods"

func main() {
  var pm mvhspods.PodManager
  pm.ReadStudents("students.csv")
  pm.MakePods()
  pm.WritePods("pod_groups.csv")
}

package main

import (
	"flag"

	"mvhspods"

	"github.com/milind-u/mlog"
)

func main() {
	sorted := flag.Bool("sorted", false,
		"Whether to sort the students output in alphabetical order")
	flag.Parse()

	mlog.SetLevel(mlog.LInfo)

	var pm mvhspods.PodManager
	pm.ReadStudents("students.csv")
	pm.MakePods(*sorted)
	pm.WritePods("pods.csv")
}

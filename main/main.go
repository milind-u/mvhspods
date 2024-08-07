package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"strings"
	"syscall/js"

	"mvhspods"
	"tests"

	"github.com/milind-u/glog"
)

var eldPopulation mvhspods.Percents
var population mvhspods.Percents

func webMain() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pods := ""
		students := args[0].String()
		podSize := mvhspods.DefaultPodSize
		if len(args) == 2 {
			podSize = args[1].Int()
			if podSize <= 0 {
				glog.Fatalln("Pod size", podSize, "is not positive integer")
			}
		} else if len(args) > 2 {
			glog.Fatalln("Expected 1 or 2 args (csv and optionally pod size), but got", len(args))
		}
		var pm mvhspods.PodManager
		pm.ReadStudentsFromString(students)
		pm.MakePods(podSize)
		eldPopulation = pm.Eld.Population()
		population = pm.Population()
		pods = pm.WritePodsToString()

		return pods
	})
}

func webPodPercents() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pods := ""
		if len(args) == 1 {
			pod := args[0].String()

			var pm mvhspods.PodManager
			pm.ReadStudentsFromString(pod)
			b := new(strings.Builder)
			pm.WritePercents(csv.NewWriter(b), mvhspods.PercentsOf(pm.Students), 0, len(pm.Students[0].Fields))
			pods = b.String()
		} else {
			glog.Errorln("Expected 1 arg (csv), but got", len(args))
		}
		return pods
	})
}

func webPodError() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pods := ""
		if len(args) == 1 {
			pod := args[0].String()

			var pm mvhspods.PodManager
			pm.ReadStudentsFromString(pod)

			var stats *tests.Stats
			// The pod is either ELD or other
			if pm.Eld.Len() > 0 {
				stats = tests.PodStatsOfPod(pm.Eld.Students, eldPopulation)
			} else {
				stats = tests.PodStatsOfPod(pm.Students, population)
			}
			return fmt.Sprintf("%v", stats.AvgErr())
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
	sampleData := flag.Bool("sample_data", false,
		"If true, using sample data so trim the last column in the csv which has sample pod numbers")
	flag.Parse()

	glog.SetSeverity(glog.InfoSeverity)

	if *web {
		js.Global().Set("makePods", webMain())
		js.Global().Set("percentsOf", webPodPercents())
		js.Global().Set("podError", webPodError())
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
		pm.WritePods("pods.csv")
	}
}

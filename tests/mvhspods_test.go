package tests

import (
  "os"
  "testing"

  "mvhspods"

  "github.com/milind-u/glog"
)

const numStudents = 600

var pm *mvhspods.PodManager

// Tests that the student weight function is working correctly
func TestWeight(t *testing.T) {
  population := mvhspods.Percents{{1, "graham"}: 0.3,
    {4, "male"}: 0.5, {7, "english"}: 0.7}
  pod := mvhspods.Percents{{1, "graham"}: 0.05,
    {4, "male"}: 0.7, {7, "english"}: 0.7}
  s := mvhspods.Student{Fields: []string{"100012345", "Graham", "Bar", "Foo", "Male", "Parent", "6501231234",
    "English", ""}, Stripped: nil}
  s.Strip()

  const expectedWeight = (0.3 - 0.05) + (0.5 - 0.7) + (0.7 - 0.7)
  const floatTolerance = 1e-5

  weight := s.Weight(population, pod)
  t.Log("Weight:", weight)
  if mvhspods.Abs(expectedWeight-weight) > floatTolerance {
    t.Error("Weight does not match expected weight of", expectedWeight)
  }
}

func TestPodSize(t *testing.T) {
  actualNumStudents := 0
  for _, pd := range []*mvhspods.PodData{&pm.PodData, &pm.Eld} {
    for _, pod := range pd.Pods() {
      actualNumStudents += len(pod)
      if len(pod) < mvhspods.DefaultPodSize || len(pod) > mvhspods.DefaultPodSize+1 {
        t.Errorf("Expected pod size between %v and %v, but got %v",
          mvhspods.DefaultPodSize, mvhspods.DefaultPodSize+1, len(pod))
      }
    }
  }
  if actualNumStudents != numStudents {
    t.Errorf("Expected number of students to be %v, but got %v", numStudents, actualNumStudents)
  }
}

func TestAlphabeticOrder(t *testing.T) {
  pm2 := initPm()
  pm2.WritePods("test.csv", true)
  pm2.Students = nil
  pm2.ReadStudents("test.csv", false)

  // make sure that the names are in sorted order
  for i := 0; i < len(pm2.Students)-1; i++ {
    if !pm2.Students.Less(i, i+1) {
      t.Error("The sort didn't work.")
    }
  }
}

func TestFewStudents(t *testing.T) {
  // Test with less ELD students than podSize
  students := GenerateStudents(numStudents)

  pm2 := &mvhspods.PodManager{Headers: Headers, PodData: mvhspods.PodData{Students: students}}

  const smallPodSize = mvhspods.DefaultPodSize/2 - 1
  eldCount := 0
  for i := 0; i < len(pm2.Students); i++ {
    if pm2.Students[i].Stripped[mvhspods.GroupMembershipsIndex] == mvhspods.EldStr {
      eldCount++
    }
  }
  if eldCount < smallPodSize {
    t.Fatalf("Only got %v eld students", eldCount)
  }
  removed := 0
  for i := 0; i < len(pm2.Students); i++ {
    if pm2.Students[i].Stripped[mvhspods.GroupMembershipsIndex] == mvhspods.EldStr {
      pm2.Remove(i)
      i--
      removed++
    }
    if removed == eldCount-smallPodSize {
      break
    }
  }

  pm2.MakePods(mvhspods.DefaultPodSize)

  if len(pm2.Eld.Pods()) != 1 {
    t.Error("Expected 1 ELD pod, got", len(pm2.Eld.Pods()))
  } else if len(pm2.Eld.Pods()[0]) != smallPodSize {
    t.Errorf("Expected %v students in ELD pod 1, got %v",
      smallPodSize, len(pm2.Eld.Pods()[0]))
  }
}

func TestEld(t *testing.T) {
  for _, s := range pm.Eld.Students {
    if groups := s.Stripped[mvhspods.GroupMembershipsIndex]; groups != mvhspods.EldStr {
      t.Error("This student is not ELD:", s)
    }
  }

  for _, s := range pm.Students {
    if groups := s.Stripped[mvhspods.GroupMembershipsIndex]; groups == mvhspods.EldStr {
      t.Error("This student is ELD:", s)
    }
  }
}

// Tests the stats of random pods and checks if groups are represented in pods similarly
// to how they are in the population
func TestPodStats(t *testing.T) {
  tolerances := [...]Stats{
    // Normal
    {
      maxErr:  0.2,
      avgErr:  0.02,
      badErrs: 10,
    },
    // ELD
    {
      maxErr:  0.15,
      avgErr:  0.035,
      badErrs: 2,
    },
  }

  for i, pd := range []*mvhspods.PodData{&pm.PodData, &pm.Eld} {
    tolerance := tolerances[i]
    stats := PodStatsWithTolerance(pd, tolerance.maxErr)
    label := "Stats:"
    if pd == &pm.Eld {
      label = "ELD stats:"
    }
    t.Log(label, stats)

    if stats.maxErr > tolerance.maxErr {
      t.Error("Percent error max exceeds tolerance of", tolerance.maxErr)
    }

    if stats.avgErr > tolerance.avgErr {
      t.Error("Average error exceeds tolerance of", tolerance.avgErr)
    }

    if stats.badErrs > tolerance.badErrs {
      t.Error("Bad error count exceeds tolerance of", tolerance.badErrs)
    }
  }
}

func TestOrder(t *testing.T) {
  students := GenerateStudents(numStudents)
  pm2 := initPm()
  pm2.WritePodsToString(false)
  for i := range pm2.Students {
    if pm2.Students[i].Stripped[0] != students[i].Stripped[0] {
      t.Fatal("Output students not in same order as input", i, "\n",
        pm2.Students[i].Stripped, students[i].Stripped)
    }
  }
}

func initPm() *mvhspods.PodManager {
  students := GenerateStudents(numStudents)
  mvhspods.WriteStudents("students.csv", Headers, students)

  pm := &mvhspods.PodManager{Headers: Headers, PodData: mvhspods.PodData{Students: students}}
  pm.MakePods(mvhspods.DefaultPodSize)
  return pm
}

func TestMain(m *testing.M) {
  glog.SetSeverity(glog.InfoSeverity)
  pm = initPm()
  os.Exit(m.Run())
}

package tests

import (
  "os"
  "testing"

  "mvhspods"

  "github.com/milind-u/glog"
)

var pm *mvhspods.PodManager

// Tests that the student weight function is working correctly
func TestWeight(t *testing.T) {
  population := mvhspods.Percents{{1, "Graham"}: 0.3,
    {4, "Male"}: 0.5, {7, "English"}: 0.7}
  pod := mvhspods.Percents{{1, "Graham"}: 0.05,
    {4, "Male"}: 0.7, {7, "English"}: 0.7}
  s := mvhspods.Student{"100012345", "Graham", "Bar", "Foo", "Male", "Parent", "6501231234",
    "English", ""}

  weight := s.Weight(population, pod)
  t.Log("Weight:", weight)
  const expectedWeight = (0.3 - 0.05) + (0.5 - 0.7) + (0.7 - 0.7)
  if mvhspods.Abs(expectedWeight-weight) > 1e-5 {
    t.Error("Weight does not match expected weight of", expectedWeight)
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

func TestEld(t *testing.T) {
  for _, s := range pm.Eld.Students {
    if groups := s[mvhspods.GroupMembershipsIndex]; groups != "eld" {
      t.Error("This student is not ELD:", s)
    }
  }

  for _, s := range pm.Students {
    if groups := s[mvhspods.GroupMembershipsIndex]; groups == "eld" {
      t.Error("This student is ELD:", s)
    }
  }
}

// Tests the stats of random pods and checks if groups are represented in pods similarly
// to how they are in the population
func TestPodStats(t *testing.T) {
  for _, eld := range [...]bool{true, false} {
    pd := &pm.PodData
    if eld {
      pd = &pm.Eld
    }

    t.Log(len(pd.Students))
    tolerances := Stats{
      maxErr:  0.37,
      avgErr:  0.1,
      badErrs: len(pd.Students) / 8,
    }

    stats := PodStats(pd)
    label := "Stats:"
    if eld {
      label = "ELD stats:"
    }
    t.Log(label, stats)

    if stats.maxErr > tolerances.maxErr {
      t.Error("Percent error max exceeds tolerance of", tolerances.maxErr)
    }

    if stats.avgErr > tolerances.avgErr {
      t.Error("Average error exceeds tolerance of", tolerances.avgErr)
    }

    if stats.badErrs > tolerances.badErrs {
      t.Error("Bad error count exceeds tolerance of", tolerances.badErrs)
    }
  }
}

func initPm() *mvhspods.PodManager {
  const numStudents = 500

  students := GenerateStudents(numStudents)

  pm := &mvhspods.PodManager{Headers: Headers, PodData: mvhspods.PodData{Students: students}}
  pm.MakePods(mvhspods.DefaultPodSize)
  return pm
}

func TestMain(m *testing.M) {
  glog.SetSeverity(glog.InfoSeverity)
  pm = initPm()
  os.Exit(m.Run())
}

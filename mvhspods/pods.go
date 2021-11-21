package mvhspods

import (
  "bufio"
  "encoding/csv"
  "io"
  "math"
  "os"
  "sort"
  "strconv"
  "strings"
  "unsafe"

  "github.com/milind-u/glog"
)

const EldPod = 0
const DefaultPodSize = 12

type Percent float32

type Percents map[Field]Percent

type PodManager struct {
  Students
  Headers    []string
  population Percents
  pods       []Students
}

func Abs(p Percent) Percent {
  return Percent(math.Abs(float64(p)))
}

func (pm *PodManager) ReadStudents(path string, sampleData bool) {
  f, err := os.Open(path)
  glog.FatalIf(err)
  bufReader := bufio.NewReader(f)
  r := csv.NewReader(bufReader)

  pm.Headers, err = r.Read()
  glog.FatalIf(err)

  for err != io.EOF {
    fields, readErr := r.Read()
    err = readErr
    if err == nil {
      s := Student(fields)
      // Trim the sample pod number if this is test data
      if sampleData {
        s = s[:len(s)-1]
      }
      // TODO: keep the whitespace and ELD level in the output
      // Remove whitespace from fields, and make everything lowercase
      // in case there were capitalization/spacing inconsistencies.
      for field := range s.weightedFields() {
        s[field.Index] = strings.ToLower(strings.ReplaceAll(s[field.Index], " ", ""))
      }
      // Trim the ELD group number to make all ELD levels the same group
      if groupMemberships := s[groupMembershipsIndex]; strings.Contains(groupMemberships, "ELD") {
        s[groupMembershipsIndex] = groupMemberships[:len(groupMemberships)-1]
      }
      pm.Students = append(pm.Students, s)
      glog.CheckNe(len(s), 0, "Read empty student")
    }
  }
}

func (pm *PodManager) MakePods(podSize int, sorted bool) {
  minWeight := Percent(math.Inf(-1))

  var addedStudents Students
  numPods := len(pm.Students) / podSize

  pm.pods = make([]Students, numPods)
  eldStudents := 0
  // TODO: don't just add all eld students to first pod, add them to first few
  for i := 0; i < len(pm.Students); i++ {
    student := pm.Students[i]
    if groupMemberships := student[groupMembershipsIndex]; strings.Contains(groupMemberships, "ELD") {
      pm.addToPod(student, i, EldPod, &pm.pods[EldPod], &addedStudents)
      eldStudents++
      i--
    }
  }

  pm.population = PercentsOf(pm.Students)
  podPercents := make([]Percents, numPods)

  enoughStudents := true
  for i := EldPod + 1; i < numPods && enoughStudents; i++ {
    enoughStudents = len(pm.Students) > podSize/2

    if enoughStudents {
      podPercents[i] = make(Percents)
      for j := 0; j < podSize && len(pm.Students) != 0; j++ {
        glog.Infoln("New student")
        // calculate percents of current pod
        maxWeight := minWeight
        var maxStudent Student
        maxIndex := 0
        for k, student := range pm.Students {
          if weight := student.Weight(pm.population, podPercents[i]); weight > maxWeight {
            maxStudent = student
            maxIndex = k
            glog.Infoln("new max:", maxStudent)
            glog.CheckGt(float64(len(maxStudent)), 0, "Student is empty 1")
          }
        }
        glog.Infoln("len", len(pm.Students))
        addPercents(maxStudent, pm.pods[i], podPercents[i])
        pm.addToPod(maxStudent, maxIndex, i, &pm.pods[i], &addedStudents)
      }
    } else {
      pm.pods = pm.pods[:i]
    }
  }

  // Cannot fit all students in pods of studentsPerPod
  if len(pm.Students) != 0 {
    for i := 0; len(pm.Students) != 0; i++ {
      index := i % len(pm.pods)
      percents := PercentsOf(pm.pods[index])

      maxIndex := 0
      maxWeight := minWeight
      for j := 0; j < len(pm.Students); j++ {
        if weight := pm.Students[j].Weight(pm.population, percents); weight > maxWeight {
          maxWeight = weight
          maxIndex = j
        }
      }
      pm.addToPod(pm.Students[maxIndex], maxIndex, index, &pm.pods[index], &addedStudents)
    }
  }

  pm.Students = addedStudents

  if sorted {
    sort.Sort(pm.Students)
  }

  for i, pod := range pm.pods {
    glog.Infoln("Pod", i)
    for _, s := range pod {
      glog.Infoln(s)
    }
    glog.Infoln()
  }

}

func PercentsOf(students Students) Percents {
  percents := make(Percents)
  for _, s := range students {
    addPercents(s, students, percents)
  }
  return percents
}

func addPercents(s Student, students Students, percents Percents) {
  if len(students) != 0 {
    for field := range s.weightedFields() {
      percents[field] += 1.0 / Percent(len(students))
    }
  }
}

func (pm *PodManager) addToPod(s Student, index int, podIndex int, pod *Students, addedStudents *Students) {
  glog.CheckGt(float64(len(s)), 0, "Empty student at index", index)
  s = append(s, strconv.Itoa(podIndex+1))
  *pod = append(*pod, s)

  // remove current student from slice of student
  pm.Students[index] = pm.Students[len(pm.Students)-1]
  pm.Students = pm.Students[:len(pm.Students)-1]
  *addedStudents = append(*addedStudents, s)
}

func (pm *PodManager) Population() Percents {
  return pm.population
}

func (pm *PodManager) Pods() []Students {
  return pm.pods
}

func (pm *PodManager) WritePods(path string) {
  WriteStudents(path, pm.Headers, pm.Students)
}

func WriteStudents(path string, headers []string, students Students) {
  f, err := os.Create(path)
  glog.FatalIf(err)

  w := csv.NewWriter(f)
  glog.FatalIf(w.Write(headers))
  err = w.WriteAll(*(*[][]string)(unsafe.Pointer(&students)))
  glog.FatalIf(err)
}

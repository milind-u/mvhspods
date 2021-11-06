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

const studentsPerPod = 5

const eldPod = 0

type Percent float32

type Percents map[field]Percent

type PodManager struct {
  Students
  headers []string
}

func (pm *PodManager) ReadStudents(path string, sampleData bool) {
  f, err := os.Open(path)
  glog.FatalIf(err)
  bufReader := bufio.NewReader(f)
  r := csv.NewReader(bufReader)

  pm.headers, err = r.Read()
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
      // Remove whitespace from fields
      for field := range s.weightedFields() {
        s[field.index] = strings.ReplaceAll(s[field.index], " ", "")
      }
      // Trim the ELD group number to make all ELD levels the same group
      if groupMemberships := s[groupMembershipsIndex]; strings.Contains(groupMemberships, "ELD") {
        s[groupMembershipsIndex] = groupMemberships[:len(groupMemberships)-1]
      }
      pm.Students = append(pm.Students, s)
    }
  }
}

func (pm *PodManager) MakePods(sorted bool) {
  minWeight := Percent(math.Inf(-1))

  var addedStudents Students
  numPods := len(pm.Students) / studentsPerPod

  pods := make([]Students, numPods, studentsPerPod)
  eldStudents := 0
  for i, student := range pm.Students {
    if groupMemberships := student[groupMembershipsIndex]; strings.Contains(groupMemberships, "ELD") {
      pm.addToPod(student, i, eldPod, &pods[eldPod], &addedStudents)
      eldStudents++
    }
  }

  population := PercentsOf(pm.Students)
  podPercents := make([]Percents, numPods, studentsPerPod)

  for i := eldPod + 1; i < numPods; i++ {
    podPercents[i] = make(Percents)
    for j := 0; j < studentsPerPod; j++ {
      // calculate percents of current pod
      maxWeight := minWeight
      var maxStudent Student
      maxIndex := 0
      for k, student := range pm.Students {
        if weight := student.weight(population, podPercents[i]); weight > maxWeight {
          maxStudent = student
          maxIndex = k
          glog.CheckGt(float64(len(maxStudent)), 0, "Student is empty 1")
        }
        glog.CheckGt(float64(len(maxStudent)), 0, "Student is empty 2", len(pm.Students))
      }

      addPercents(maxStudent, pods[i], podPercents[i])
      pm.addToPod(maxStudent, maxIndex, i, &pods[i], &addedStudents)
    }
  }

  // Cannot fit all students in pods of studentsPerPod
  // TODO: Refactor so that there are the min number of students in each pod,
  // and before adding a student to pod make sure no others have higher weight in that pod
  if len(pm.Students) != 0 {
    for len(pm.Students) != 0 {
      s := pm.Students[0]
      maxPod := 0
      maxWeight := minWeight
      for i, pod := range pods {
        if i != eldPod {
          if weight := s.weight(population, PercentsOf(pod)); weight > maxWeight {
            maxWeight = weight
            maxPod = i
          }
        }
      }
      pm.addToPod(s, 0, maxPod, &pods[maxPod], &addedStudents)
    }
  }

  pm.Students = addedStudents

  if sorted {
    sort.Sort(pm.Students)
  }

  for i, pod := range pods {
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
  for field := range s.weightedFields() {
    percents[field] += 1.0 / Percent(len(students))
  }
}

func (pm *PodManager) addToPod(s Student, index int, podIndex int, pod *Students, addedStudents *Students) {
  glog.CheckGt(float64(len(s)), 0, index)
  s = append(s, strconv.Itoa(podIndex+1))
  *pod = append(*pod, s)

  // remove current student from slice of student
  pm.Students[index] = pm.Students[len(pm.Students)-1]
  pm.Students = pm.Students[:len(pm.Students)-1]
  *addedStudents = append(*addedStudents, s)
}

func (pm *PodManager) WritePods(path string) {
  WriteStudents(path, pm.headers, pm.Students)
}

func WriteStudents(path string, headers []string, students Students) {
  f, err := os.Create(path)
  glog.FatalIf(err)

  w := csv.NewWriter(f)
  glog.FatalIf(w.Write(headers))
  err = w.WriteAll(*(*[][]string)(unsafe.Pointer(&students)))
  glog.FatalIf(err)
}

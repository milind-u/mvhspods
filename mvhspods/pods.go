package mvhspods

import (
  "bufio"
  "encoding/csv"
  "io"
  "os"
  "strconv"
  "unsafe"

  "github.com/milind-u/mlog"
)

const studentsPerPod = 10

type percents map[field]float32

type PodManager struct {
  students []student
}

func (pm *PodManager) ReadStudents(path string) {
  f, err := os.Open(path)
  mlog.FatalIf(err)
  bufReader := bufio.NewReader(f)
  r := csv.NewReader(bufReader)

  // Get rid of the headers
  _, err = r.Read()
  mlog.FatalIf(err)

  for err != io.EOF {
    fields, readErr := r.Read()
    err = readErr
    pm.students = append(pm.students, student(fields))
  }
}

func (pm *PodManager) MakePods() {
  population := pm.percentsOf(pm.students)
  numPods := len(pm.students)/studentsPerPod
  var addedStudents []student
  for i:=0; i<numPods; i++{
    //create student array for current pod
    var pod []student
    for j := 0; j<studentsPerPod; j++ {
      //calculate percents of current pod
      podPercents := pm.percentsOf(pod)
      var maxWeight float32
      var maxStudent student
      maxIndex := 0
      for k, student := range pm.students {
        if weight := student.weight(population, podPercents); weight > maxWeight{
          //set maxWeight to the current weight
          maxWeight = weight
          //Above set max
          maxStudent = student
          maxIndex = k
        }
      }
      //append to the String array of maxStudent the pod number
      maxStudent = append(maxStudent, strconv.Itoa(j+1))
      pod = append(pod, maxStudent)

      //remove current student from array of student
      //set students at maxIndex to be len-1 of students
      // set students to be students[: len -1]
      pm.students[maxIndex] = pm.students[len(pm.students) - 1]
      pm.students = pm.students[: len(pm.students) - 1]
      //append current Student to addedStudents
      addedStudents = append(addedStudents, maxStudent)
    }
  }
}

func (pm *PodManager) percentsOf(students []student) percents {
  var percents percents
  for _, s := range students {
    for field := range s.weightedFields() {
      percents[field] += 1.0 / float32(len(pm.students))
    }
  }
  return percents
}

func (pm *PodManager) WritePods(path string) {
  f, err := os.Open(path)
  mlog.FatalIf(err)

  w := csv.NewWriter(f)
  err = w.WriteAll(*(*[][]string)(unsafe.Pointer(&pm.students)))
  mlog.FatalIf(err)
}

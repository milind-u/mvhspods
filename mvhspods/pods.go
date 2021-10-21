package mvhspods

import (
  "bufio"
  "encoding/csv"
  "io"
  "os"

  "github.com/milind-u/mlog"
)

type percents map[string]float32

type PodManager struct {
  students []student
  pods [][]student
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
  mlog.CheckEq(pm.pods, nil, "Already made pods")

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

}

package mvhspods

import (
  "fmt"
  "strings"

  "github.com/milind-u/glog"
)

// Indices of the student fields that are weighted
var weightedFields = [...]int{1, 4, 7, 8}

const groupMembershipsIndex = 8
const lastNameIndex = 2
const firstNameIndex = 3
const idIndex = 0

type Student []string

type Students []Student

type Field struct {
  index int
  string
}

func (s Student) weightedFields() chan Field {
  c := make(chan Field, len(weightedFields))
  for _, index := range weightedFields {
    if index < len(s) && s[index] != "" {
      c <- Field{index, s[index]}
    }
  }
  close(c)
  return c
}

func (s Student) weight(population Percents, pod Percents) Percent {
  var weight Percent
  for field := range s.weightedFields() {
    weight += population[field] - pod[field]
  }
  return weight
}

func (s Students) String() string {
  b := new(strings.Builder)
  for _, student := range s {
    _, err := fmt.Fprintln(b, student)
    glog.WarningIf(err)
  }
  return b.String()
}

func (s Students) Len() int {
  return len(s)
}

func (s Students) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}

func (s Students) Less(i, j int) bool {
  diff := strings.Compare(s[i][lastNameIndex], s[j][lastNameIndex])
  if diff == 0 {
    diff = strings.Compare(s[i][firstNameIndex], s[j][firstNameIndex])
    if diff == 0 {
      diff = strings.Compare(s[i][idIndex], s[j][idIndex])
    }
  }
  return diff < 0
}

package mvhspods

import (
  "fmt"
  "strings"

  "github.com/milind-u/glog"
)

// Indices of the student Fields that are weighted
var weightedFields = [...]int{1, 4, 7, 8}

const GroupMembershipsIndex = 8
const lastNameIndex = 2
const firstNameIndex = 3
const idIndex = 0

const EldStr = "eld1/2"

type Student struct {
  Fields   []string
  Stripped []string
}

type Students []Student

type Field struct {
  Index int
  Name  string
}

func (s *Student) weightedFields() chan Field {
  c := make(chan Field, len(weightedFields))
  for _, index := range weightedFields {
    if index < len(s.Stripped) && s.Stripped[index] != "" {
      c <- Field{index, s.Stripped[index]}
    }
  }
  close(c)
  return c
}

func (s *Student) Weight(population Percents, pod Percents) Percent {
  var weight Percent
  for field := range s.weightedFields() {
    weight += population[field] - pod[field]
  }
  return weight
}

func (s *Student) Strip() {
  s.Stripped = make([]string, len(s.Fields))
  copy(s.Stripped, s.Fields)
  // TODO: keep the whitespace and ELD level in the output
  // Remove whitespace from Fields, and make everything lowercase
  // in case there were capitalization/spacing inconsistencies.
  for field := range s.weightedFields() {
    s.Stripped[field.Index] = strings.ToLower(strings.ReplaceAll(s.Stripped[field.Index], " ", ""))
  }
  // Trim the ELD group number to make all ELD levels the same group
  if groupMemberships := s.Stripped[GroupMembershipsIndex]; strings.Contains(groupMemberships, "eld") {
    if groupMemberships == "eld1" || groupMemberships == "eld2" {
      s.Stripped[GroupMembershipsIndex] = EldStr
    } else {
      s.Stripped[GroupMembershipsIndex] = "eld3/4"
    }
  }
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
  diff := strings.Compare(s[i].Stripped[lastNameIndex], s[j].Stripped[lastNameIndex])
  if diff == 0 {
    diff = strings.Compare(s[i].Stripped[firstNameIndex], s[j].Stripped[firstNameIndex])
    if diff == 0 {
      diff = strings.Compare(s[i].Stripped[idIndex], s[j].Stripped[idIndex])
    }
  }
  return diff < 0
}

func (s *Students) Remove(i int) {
  (*s)[i] = (*s)[len(*s)-1]
  *s = (*s)[:len(*s)-1]
}

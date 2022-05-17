package mvhspods

import (
  "fmt"
  "regexp"
  "strings"

  "github.com/milind-u/glog"
)

const idIndex = 0
const schoolIndex = 1
const lastNameIndex = 2
const firstNameIndex = 3
const genderIndex = 4
const languageIndex = 7
const GroupMembershipsIndex = 8

// Indices of the student Fields that are weighted
var weightedFields = [...]int{schoolIndex, genderIndex, languageIndex, GroupMembershipsIndex}

const EldStr = "eld1/2"
const advEldStr = "eld3/4"

type Student struct {
  Fields              []string
  Stripped            []string
  groupMemberships    []string
  groupMembershipsSet bool
}

type Students []Student

type Field struct {
  Index int
  Name  string
}

func (s *Student) weightedFields() chan Field {
  lenOffset := 0
  if s.groupMembershipsSet {
    lenOffset = len(s.groupMemberships) - 1
  }
  c := make(chan Field, len(weightedFields)+lenOffset)
  for _, index := range weightedFields {
    if index == GroupMembershipsIndex && s.groupMembershipsSet && s.Stripped[index] != "" {
      for _, group := range s.groupMemberships {
        c <- Field{index, group}
      }
    } else {
      if index < len(s.Stripped) && s.Stripped[index] != "" {
        c <- Field{index, s.Stripped[index]}
      }
    }
  }
  close(c)
  return c
}

func computeDelta(f Field, population Percents, pod Percents) Percent {
  return population[f] - pod[f]
}

func (s *Student) Weight(population Percents, pod Percents) Percent {
  weight := Percent(0)
  for field := range s.weightedFields() {
    weight += computeDelta(field, population, pod)
  }
  return weight
}

var eldRegexp1 = regexp.MustCompile(`eld(1|2)`)
var eldRegexp2 = regexp.MustCompile(`eld(3|4)`)

func (s *Student) Strip() {
  s.Stripped = make([]string, len(s.Fields))
  copy(s.Stripped, s.Fields)
  // Remove whitespace from Fields, and make everything lowercase
  // in case there were capitalization/spacing inconsistencies.
  for field := range s.weightedFields() {
    s.Stripped[field.Index] = strings.ToLower(strings.ReplaceAll(s.Stripped[field.Index], " ", ""))
  }

  // Trim the ELD group number to make all ELD levels the same group
  s.Stripped[GroupMembershipsIndex] = eldRegexp1.ReplaceAllString(s.Stripped[GroupMembershipsIndex], EldStr)
  s.Stripped[GroupMembershipsIndex] = eldRegexp2.ReplaceAllString(s.Stripped[GroupMembershipsIndex], advEldStr)

  s.groupMemberships = strings.Split(s.Stripped[GroupMembershipsIndex], ",")
  s.groupMembershipsSet = true
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

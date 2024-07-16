package mvhspods

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/milind-u/glog"
)

type Field struct {
	Index int
	Name  string
}

// Indices in stripped array
var gender = Field{0, "Gender"}
var school = Field{1, "LastSchl"}
var language = Field{2, "Description_HL"}
var groupMemberships = Field{3, "GroupMemberships?"}

var weightedFields = map[string]int{gender.Name: gender.Index,
	school.Name: school.Index, language.Name: language.Index,
	groupMemberships.Name: groupMemberships.Index}

const EldStr = "eld1/2"
const advEldStr = "eld3/4"

type Student struct {
	Fields              []string
	Stripped            []string
	groupMemberships    []string
	groupMembershipsSet bool
}

type Students []Student

func (s *Student) Gender() string {
	return s.Stripped[gender.Index]
}

func (s *Student) School() string {
	return s.Stripped[school.Index]
}

func (s *Student) Language() string {
	return s.Stripped[language.Index]
}

func (s *Student) GroupMemberships() string {
	return s.Stripped[groupMemberships.Index]
}

func (s *Student) weightedFields() chan Field {
	lenOffset := 0
	if s.groupMembershipsSet {
		lenOffset = len(s.groupMemberships) - 1
	}
	c := make(chan Field, len(weightedFields)+lenOffset)
	for i := 0; i < len(weightedFields); i++ {
		if i == groupMemberships.Index && s.groupMembershipsSet && s.GroupMemberships() != "" {
			for _, group := range s.groupMemberships {
				c <- Field{i, group}
			}
		} else {
			if i < len(s.Stripped) && s.Stripped[i] != "" {
				c <- Field{i, s.Stripped[i]}
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

func (s *Student) Strip(headers []string) {
	s.Stripped = make([]string, len(weightedFields))
	// Remove whitespace from Fields, and make everything lowercase
	// in case there were capitalization/spacing inconsistencies.
	for i, field := range s.Fields {
		if index, ok := weightedFields[headers[i]]; ok {
			s.Stripped[index] = strings.ToLower(strings.ReplaceAll(field, " ", ""))
		}
	}

	// Trim the ELD group number to make all ELD levels the same group
	s.Stripped[groupMemberships.Index] = eldRegexp1.ReplaceAllString(s.GroupMemberships(), EldStr)
	s.Stripped[groupMemberships.Index] = eldRegexp2.ReplaceAllString(s.GroupMemberships(), advEldStr)

	s.groupMemberships = strings.Split(s.GroupMemberships(), ",")
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

func (s *Students) Remove(i int) {
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}

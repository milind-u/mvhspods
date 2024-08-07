package mvhspods

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/milind-u/glog"
)

const DefaultPodSize = 12

type Percent float32

type Percents map[Field]Percent

type PodData struct {
	Students
	population Percents
	pods       []Students
}

type PodManager struct {
	Headers []string
	// Data for ELD 1/2 students (3/4 are integrated with the other students)
	Eld PodData
	PodData
}

func Abs(p Percent) Percent {
	return Percent(math.Abs(float64(p)))
}

func (pd *PodData) Population() Percents {
	return pd.population
}

func (pd *PodData) Pods() []Students {
	return pd.pods
}

func (pm *PodManager) ReadStudents(path string, sampleData bool) {
	f, err := os.Open(path)
	glog.FatalIf(err)
	pm.readStudents(f, sampleData)
}

func (pm *PodManager) ReadStudentsFromString(data string) {
	pm.readStudents(strings.NewReader(data), false)
}

func (pm *PodManager) readStudents(reader io.Reader, sampleData bool) {
	bufReader := bufio.NewReader(reader)
	r := csv.NewReader(bufReader)
	index := 0
	headers, err := r.Read()
	alreadyMade := (headers[len(headers)-1] == "")

	pm.Headers = append(headers, "")

	glog.FatalIf(err)

	for err != io.EOF {
		Fields, readErr := r.Read()
		err = readErr
		if err == nil {
			s := Student{Fields: Fields, Stripped: nil}
			// Trim the sample pod number if this is test data
			if sampleData {
				s.Fields = s.Fields[:len(s.Fields)-1]
			}
			s.Strip(pm.Headers)
			s.Stripped = append(s.Stripped, strconv.Itoa(index))
			index++
			pm.Students = append(pm.Students, s)
			glog.CheckNe(len(s.Fields), 0, "Read empty student")
		}
	}

	if alreadyMade { // Already have pods made
		pm.splitEld()
		pm.Eld.population = PercentsOf(pm.Eld.Students)
		pm.population = PercentsOf(pm.Students)
	}
}

func (pm *PodManager) splitEld() {
	for i := 0; i < len(pm.Students); i++ {
		student := pm.Students[i]
		if groupMemberships := student.GroupMemberships(); strings.Contains(groupMemberships, EldStr) {
			pm.Eld.Students = append(pm.Eld.Students, student)
			pm.Students.Remove(i)
			i--
		}
	}
}

func (pm *PodManager) MakePods(podSize int) {
	pm.splitEld()

	if pm.Eld.Students.Len() > 0 {
		pm.makePods(&pm.Eld.Students, &pm.Eld.pods, &pm.Eld.population, podSize, true)
	}
	if pm.Students.Len() > 0 {
		pm.makePods(&pm.Students, &pm.pods, &pm.population, podSize, false)
	}
}

func (pm *PodManager) makePods(students *Students, pods *[]Students, population *Percents, podSize int, eld bool) {
	minWeight := Percent(math.Inf(-1))

	var addedStudents Students
	*population = PercentsOf(*students)

	numPods := int(math.Max(float64(len(*students)/podSize), 1.0))
	*pods = make([]Students, numPods)

	podOffset := 0
	if !eld {
		podOffset = len(pm.Eld.pods)
	}

	podPercents := make([]Percents, numPods)

	for j := 0; j < podSize && len(*students) != 0; j++ {
		for i := 0; i < numPods; i++ {
			maxWeight := minWeight
			var maxStudent Student
			maxIndex := 0
			for k, student := range *students {
				if weight := student.Weight(*population, podPercents[i]); weight > maxWeight {
					maxStudent = student
					maxIndex = k
					maxWeight = weight
					glog.CheckGt(float64(len(maxStudent.Stripped)), 0, "Student is empty")
				}
			}
			pm.addToPod(maxStudent, maxIndex, i, podOffset, &(*pods)[i], &addedStudents, students)
			podPercents[i] = PercentsOf((*pods)[i])
		}
	}

	if len(*students) != 0 { // Couldn't fit all students in the pods
		for i := 0; len(*students) != 0; i++ {
			podIndex := i % len(*pods)
			percents := PercentsOf((*pods)[podIndex])

			maxIndex := 0
			maxWeight := minWeight
			for j := 0; j < len(*students); j++ {
				if weight := (*students)[j].Weight(*population, percents); weight > maxWeight {
					maxWeight = weight
					maxIndex = j
				}
			}
			pm.addToPod((*students)[maxIndex], maxIndex, podIndex, podOffset, &(*pods)[podIndex], &addedStudents, students)
		}
	}

	*students = addedStudents
}

func PercentsOf(students Students) Percents {
	percents := make(Percents)
	for _, s := range students {
		for field := range s.weightedFields() {
			percents[field] += 1.0 / Percent(len(students))
		}
	}
	return percents
}

func (pm *PodManager) addToPod(s Student, index int, podIndex int, podOffset int, pod *Students, addedStudents *Students, students *Students) {
	glog.CheckGt(float64(len(s.Fields)), 0, "Empty student at index", index)
	s.Fields = append(s.Fields, strconv.Itoa(podIndex+1+podOffset))
	*pod = append(*pod, s)

	// remove current student from slice of student
	students.Remove(index)
	*addedStudents = append(*addedStudents, s)
}

func (pm *PodManager) WritePods(path string) {
	f, err := os.Create(path)
	glog.FatalIf(err)
	pm.writePodsWithWriter(f)
}

func (pm *PodManager) WritePodsToString() string {
	b := new(strings.Builder)
	pm.writePodsWithWriter(b)
	return b.String()
}

func (pm *PodManager) writePodsWithWriter(writer io.Writer) {
	w := csv.NewWriter(writer)
	glog.FatalIf(w.Write(pm.Headers))

	students := make(Students, len(pm.Students)+len(pm.Eld.Students))
	for _, studentGroup := range [...]Students{pm.Eld.Students, pm.Students} {
		for _, s := range studentGroup {
			index, err := strconv.Atoi(s.Stripped[len(s.Stripped)-1])
			glog.FatalIf(err)
			students[index] = s
		}
	}
	pm.Students = students
	writeStudentsWithWriter(w, pm.Students)

	// Write pod percents
	numCols := len(pm.Students[0].Fields)
	pm.WritePercents(w, pm.Population(), 0, numCols)
	for i, pod := range pm.Eld.pods {
		pm.WritePercents(w, PercentsOf(pod), i+1, numCols)
	}
	for i, pod := range pm.pods {
		pm.WritePercents(w, PercentsOf(pod), i+1+len(pm.Eld.pods), numCols)
	}

	w.Flush()
}

func (pm *PodManager) WritePercents(w *csv.Writer, percents Percents, podNum, numCols int) {
	output := make([]strings.Builder, numCols)
	for f, p := range percents {
		output[f.Index].WriteString(fmt.Sprintf("%v: %v, ", f.Name, p))
	}
	output[len(output)-1].WriteString(strconv.Itoa(podNum))

	strs := make([]string, numCols)
	const separator = ", "
	for i, b := range output {
		s := b.String()
		if len(s) >= len(separator) && s[len(s)-len(separator):] == separator {
			s = s[:len(s)-len(separator)]
		}
		strs[i] = s
	}

	glog.FatalIf(w.Write(strs))
	w.Flush()
}

func WriteStudents(path string, headers []string, students Students) {
	f, err := os.Create(path)
	glog.FatalIf(err)
	w := csv.NewWriter(f)
	glog.FatalIf(w.Write(headers))
	writeStudentsWithWriter(w, students)
}

func writeStudentsWithWriter(w *csv.Writer, students Students) {
	for _, s := range students {
		err := w.Write(s.Fields)
		glog.FatalIf(err)
	}
	w.Flush()
}

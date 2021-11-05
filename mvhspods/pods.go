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

	"github.com/milind-u/mlog"
)

// If true, trim the last column which has sample pod numbers
const testData = true

const studentsPerPod = 5

const eldPod = 0

type percents map[field]float32

type PodManager struct {
	students
	headers []string
}

func (pm *PodManager) ReadStudents(path string) {
	f, err := os.Open(path)
	mlog.FatalIf(err)
	bufReader := bufio.NewReader(f)
	r := csv.NewReader(bufReader)

	pm.headers, err = r.Read()
	mlog.FatalIf(err)

	for err != io.EOF {
		fields, readErr := r.Read()
		err = readErr
		if err == nil {
			s := student(fields)
			// Trim the sample pod number if this is test data
			if testData {
				s = s[:len(s)-1]
			}
			// TODO: keep the whitespace and ELD level in the output
			// Remove whitespace from fields
			for field := range s.weightedFields() {
				s[field.index] = strings.ReplaceAll(s[field.index], " ", "")
			}
			// Trim the ELD group number to make all ELD levels the same group
			// TODO: put all ELD students in same pod
			if groupMemberships := s[groupMembershipsIndex]; strings.Contains(groupMemberships, "ELD") {
				s[groupMembershipsIndex] = groupMemberships[:len(groupMemberships)-1]
			}
			pm.students = append(pm.students, s)
		}
	}

	mlog.Infoln(pm.students)
}

func (pm *PodManager) MakePods(sorted bool) {
	minWeight := float32(math.Inf(-1))

	var addedStudents students
	numPods := len(pm.students) / studentsPerPod

	pods := make([]students, numPods, studentsPerPod)
	eldStudents := 0
	for i, student := range pm.students {
		if groupMemberships := student[groupMembershipsIndex]; strings.Contains(groupMemberships, "ELD") {
			pm.addToPod(student, i, eldPod, &pods[eldPod], &addedStudents)
			eldStudents++
		}
	}

	population := pm.percentsOf(pm.students)

	for i := eldPod + 1; i < numPods; i++ {
		// create student array for current pod
		for j := 0; j < studentsPerPod; j++ {
			// calculate percents of current pod
			podPercents := pm.percentsOf(pods[i])
			maxWeight := minWeight
			var maxStudent student
			maxIndex := 0
			for k, student := range pm.students {
				if weight := student.weight(population, podPercents); weight > maxWeight {
					maxStudent = student
					maxIndex = k
					mlog.CheckGt(float64(len(maxStudent)), 0, "Student is empty 1")
				}
				mlog.CheckGt(float64(len(maxStudent)), 0, "Student is empty 2", len(pm.students))
			}

			pm.addToPod(maxStudent, maxIndex, i, &pods[i], &addedStudents)
		}

	}

	// Cannot fit all students in pods of studentsPerPod
	if len(pm.students) != 0 {
		for len(pm.students) != 0 {
			s := pm.students[0]
			maxPod := 0
			maxWeight := minWeight
			for i, pod := range pods {
				if i != eldPod {
					if weight := s.weight(population, pm.percentsOf(pod)); weight > maxWeight {
						maxWeight = weight
						maxPod = i
					}
				}
			}
			pm.addToPod(s, 0, maxPod, &pods[maxPod], &addedStudents)
		}
	}

	pm.students = addedStudents

	if sorted {
		sort.Sort(pm.students)
	}

	// TODO: Refactor so that there are the min number of students in each pod,
	// and before adding a student to pod make sure no others have higher weight in that pod
	for i, pod := range pods {
		mlog.Infoln("Pod", i)
		for _, s := range pod {
			mlog.Infoln(s)
		}
		mlog.Infoln()
	}

}

func (pm *PodManager) percentsOf(students students) percents {
	percents := make(percents)
	for _, s := range students {
		for field := range s.weightedFields() {
			percents[field] += 1.0 / float32(len(pm.students))
		}
	}
	return percents
}

func (pm *PodManager) addToPod(s student, index int, podIndex int, pod *students, addedStudents *students) {
	mlog.CheckGt(float64(len(s)), 0, index)
	s = append(s, strconv.Itoa(podIndex+1))
	*pod = append(*pod, s)

	// remove current student from slice of student
	pm.students[index] = pm.students[len(pm.students)-1]
	pm.students = pm.students[:len(pm.students)-1]
	*addedStudents = append(*addedStudents, s)
}

func (pm *PodManager) WritePods(path string) {
	f, err := os.Create(path)
	mlog.FatalIf(err)

	w := csv.NewWriter(f)
	mlog.FatalIf(w.Write(pm.headers))
	err = w.WriteAll(*(*[][]string)(unsafe.Pointer(&pm.students)))
	mlog.FatalIf(err)
}

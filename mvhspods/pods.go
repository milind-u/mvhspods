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

// If true, trim the last column which has sample pod numbers
const testData = true

const studentsPerPod = 5

type percents map[field]float32

type PodManager struct {
	students []student
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
				s = s[:len(s) - 1]
			}
			pm.students = append(pm.students, s)
		}
	}
}

func (pm *PodManager) MakePods() {
	population := pm.percentsOf(pm.students)
	numPods := len(pm.students) / studentsPerPod
	var addedStudents []student
	pods := make([][]student, numPods, studentsPerPod)

	for i := 0; i < numPods; i++ {
		// create student array for current pod
		for j := 0; j < studentsPerPod; j++ {
			// calculate percents of current pod
			podPercents := pm.percentsOf(pods[i])
			var maxWeight float32
			var maxStudent student
			maxIndex := 0
			for k, student := range pm.students {
				if weight := student.weight(population, podPercents); weight > maxWeight {
					maxWeight = weight
					maxStudent = student
					maxIndex = k
				}
			}

			pm.addToPod(maxStudent, maxIndex, i, &pods[i], &addedStudents)
		}

	}

	// Cannot fit all students in pods of studentsPerPod
	if len(pm.students) != 0 {
		for len(pm.students) != 0 {
			s := pm.students[0]
			maxPod := 0
			var maxWeight float32
			for j, pod := range pods {
				if weight := s.weight(population, pm.percentsOf(pod)); weight > maxWeight {
					maxWeight = weight
					maxPod = j
				}
			}
			pm.addToPod(s, 0, maxPod, &pods[maxPod], &addedStudents)
		}
	}

	pm.students = addedStudents

	for i, pod := range pods {
		mlog.Infoln("Pod", i)
		for _, s := range pod {
			mlog.Infoln(s)
		}
		mlog.Infoln()
	}
}

func (pm *PodManager) percentsOf(students []student) percents {
	percents := make(percents)
	for _, s := range students {
		for field := range s.weightedFields() {
			percents[field] += 1.0 / float32(len(pm.students))
		}
	}
	return percents
}

func (pm *PodManager) addToPod(s student, index int, podIndex int, pod *[]student, addedStudents *[]student) {
	s = append(s, strconv.Itoa(podIndex + 1))
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

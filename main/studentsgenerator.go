package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"mvhspods"

	"github.com/milind-u/mlog"
)

type field uint8

const (
	categorical field = iota
	randStr
	id
	phoneNum
)

type group struct {
	name    string
	percent mvhspods.Percent
}

var headers = strings.Split("StudentID,LastSchl,LastName,FirstName,Gender,Guardian,PrimaryPhone,Description_HL,GroupMemberships?",
	",")

// StudentID, LastSchl, LastName, FirstName, Gender, Guardian, PrimaryPhone, Description_HL, GroupMemberships?
var fields = [...]field{id, categorical, randStr, randStr, categorical, randStr, phoneNum, categorical, categorical}

// Array of possible categories for the fields that are categorical
var categories = [...][]group{{{"Graham", 0.3},
	{"Crittenden", 0.3},
	{"Blach", 0.4}},
	{{"Male", 0.45},
		{"Female", 0.45},
		{"Non-binary", 0.1}},
	{{"English", 0.4},
		{"Spanish", 0.3},
		{"Hindi", 0.1},
		{"Mandarin", 0.1},
		{"Persian", 0.1}},
	{{"", 0.5},
		{"ELD 1", 0.1},
		{"ELD 2", 0.1},
		{"AVID", 0.2},
		{"Band", 0.1}}}

// Checks that the percentages of all groups sum to 1
func checkCategories() {
	for _, category := range categories {
		sum := mvhspods.Percent(0.0)
		for _, g := range category {
			sum += g.percent
		}
		mlog.Check(math.Abs(float64(sum - mvhspods.Percent(1.0))) < 0.00001,
          "Invalid percentages for category:", category)
	}
}

func chooseRandGroup(category []group) string {
	var name string
	r := mvhspods.Percent(rand.Float32())
	sum := mvhspods.Percent(0)
	for _, group := range category {
		sum += group.percent
		if r < sum {
			name = group.name
			break
		}
	}
	return name
}

func makeRandStr() string {
	b := new(strings.Builder)
	l := rand.Intn(5) + 3
	for i := 0; i < l; i++ {
		c := byte(rand.Intn('z'-'a') + 'a')
		if i == 0 {
			c = byte(unicode.ToUpper(rune(c)))
		}
		b.WriteByte(c)
	}
	return b.String()
}

func main() {
	// TODO: add go test with this data to make sure created pods are diverse

	numStudents := flag.Int("num_students", 100,
		"Number of students to generate")
	flag.Parse()
  mlog.SetLevel(mlog.LInfo)

	// Use same seed to have reproducible results
	const seed = 94040
	rand.Seed(seed)

	checkCategories()

	students := make(mvhspods.Students, *numStudents)
	for i := range students {
		students[i] = make(mvhspods.Student, len(fields))
		categoryIndex := 0
		for j := range fields {
			switch fields[j] {
			case categorical:
				students[i][j] = chooseRandGroup(categories[categoryIndex])
				categoryIndex++
			case randStr:
				students[i][j] = makeRandStr()
			case id:
				students[i][j] = strconv.Itoa(int(math.Pow(10, 8)) +
					rand.Intn(int(math.Pow(10, 5))))
			case phoneNum:
				students[i][j] = strconv.Itoa(rand.Intn(int(math.Pow(10, 10))-int(math.Pow(10, 9))) +
					int(math.Pow(10, 9)))
			default:
				log.Fatalln("Unknown field type")
			}
		}
	}

	mlog.Infoln(students)
	mlog.Infoln(mvhspods.PercentsOf(students))
	mvhspods.WriteStudents("students.csv", headers, students)
}

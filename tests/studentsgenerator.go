package tests

import (
  "log"
  "math"
  "math/rand"
  "strconv"
  "strings"
  "unicode"

  "mvhspods"

  "github.com/milind-u/glog"
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

var Headers = strings.Split("StudentID,LastSchl,LastName,FirstName,Gender,Guardian,PrimaryPhone,Description_HL,GroupMemberships?",
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
  {{"", 0.4},
    {"ELD 1", 0.05},
    {"ELD 2", 0.05},
    {"ELD 3", 0.05},
    {"ELD 4", 0.05},
    {"ELD 3, AVID", 0.1},
    {"AVID", 0.2},
    {"Band", 0.1}}}

// Checks that the percentages of all groups sum to 1
func checkCategories() {
  for _, category := range categories {
    sum := mvhspods.Percent(0.0)
    for _, g := range category {
      sum += g.percent
    }
    glog.Check(math.Abs(float64(sum-mvhspods.Percent(1.0))) < 1e-5,
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

func GenerateStudents(numStudents int) mvhspods.Students {
  // Use same seed to have reproducible results
  const seed = 94040
  rand.Seed(seed)

  checkCategories()

  students := make(mvhspods.Students, numStudents)
  for i := range students {
    students[i] = mvhspods.Student{Fields: make([]string, len(fields)), Stripped: nil}
    categoryIndex := 0
    for j := range fields {
      switch fields[j] {
      case categorical:
        students[i].Fields[j] = chooseRandGroup(categories[categoryIndex])
        categoryIndex++
      case randStr:
        students[i].Fields[j] = makeRandStr()
      case id:
        students[i].Fields[j] = strconv.Itoa(1e8 +
            rand.Intn(1e5))
      case phoneNum:
        students[i].Fields[j] = strconv.Itoa(rand.Intn(1e10-1e9) +
            1e9)
      default:
        log.Fatalln("Unknown field type")
      }
    }
    students[i].Strip(Headers)
    students[i].Stripped = append(students[i].Stripped, strconv.Itoa(i))
  }

  return students
}

package mvhspods

// Indices of the student fields that are weighted
var weightedFields = [...]int{1, 4, 7, 8}

type student []string

type field struct {
  index int
  string
}

func (s student) weightedFields() chan field {
  c := make(chan field, len(weightedFields))
  for _, index := range weightedFields {
    if s[index] != "" {
      c <- field{index, s[index]}
    }
  }
  return c
}

func (s student) weight(population percents, pod percents) float32 {
  var weight float32
  for field := range s.weightedFields() {
    weight += population[field] - pod[field]
  }
  return weight
}

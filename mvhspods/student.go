package mvhspods

// Indices of the student fields that are weighted
var weightedFields = [...]int{1, 4, 7, 8}

type student []string

func (s student) weightedFields() chan string {
  c := make(chan string, len(weightedFields))
  for _, index := range weightedFields {
    if s[index] != "" {
      c <- s[index]
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

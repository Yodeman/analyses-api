package util

import (
    "fmt"
    "math/rand"
    "strings"
    "strconv"
    "time"

    "gonum.org/v1/gonum/mat"
)

const (
    alphabets   = "abcdefghijklmnopqrstuvxyz"
    nameLen     = 6
    passwordLen = 8
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer betwwen min and max
func RandomInt(min, max int64) int64 {
    return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of given length `n`
func RandomString(n int) string {
    var sb strings.Builder
    k := len(alphabets)
    for i := 0; i < n; i++ {
        sb.WriteByte(alphabets[rand.Intn(k)])
    }

    return sb.String()
}

// RandomUser generates a random application username
func RandomUser() string {
    return RandomString(nameLen)
}

// RandomEmail generates a random user email
func RandomEmail() string {
    return fmt.Sprintf("%s@email.com", RandomString(nameLen))
}

// RandomPassword generates a random user password
func RandomPassword() string {
    return RandomString(passwordLen)
}

// RandomCSV generates random csv string
func RandomCSV(rows, cols int) string {
    var sampleCSV string
    r, c := 0, 0

    for r != rows {
        elem := strconv.FormatFloat(rand.Float64() * 100, 'f', 6, 64)
        if c == (cols - 1) {
            sampleCSV += elem+"\n"
            c = 0
            r++
        } else if c < cols {
            sampleCSV += elem+","
            c++
        }
    }

    return sampleCSV
}

// RandomData generate random serialized user data
func RandomData() ([]byte, error) {
    sample_text := RandomCSV(10,10)

    reader := strings.NewReader(sample_text)

    rows, cols, data, err := ParseCSVToFloatSlice(reader)
    if err != nil {
        return []byte{}, err
    }

    m := mat.NewDense(rows, cols, data)
    return m.MarshalBinary()
}

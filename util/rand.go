package util

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
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

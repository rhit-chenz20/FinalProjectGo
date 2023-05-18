package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	//"errors"
	"reflect"
	"fmt"
	"io"
)

type student struct {
	ID      string
	Name    string
	Year    string
	Courses map[string]float64
}

// starter data for courses
var courses = []string{
	"CSSE403",
	"CSSE374",
	"CSSE304",
}

// starter data for students
var students = []student{
	{ID: "1", Name: "John Doe", Year: "Junior", Courses: make(map[string]float64)},
	{ID: "2", Name: "Jane Doe", Year: "Sophomore", Courses: make(map[string]float64)},
	{ID: "3", Name: "Dave Smith", Year: "Senior", Courses: make(map[string]float64)},
}
// Encrypted courses
var coursesEncrypted = []string{}

// Encryted
var studentsEncrypted = []student{}
func encodeCourseTest(course string) string {

	block, err := newCipherBlock("0")
	if err != nil {
		return "oopsies"
	}

	ciphertext := make([]byte, aes.BlockSize+len(course))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "oopsies"
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(course))

 	fmt.Printf("%x", ciphertext)
	coursesEncrypted = append(coursesEncrypted, fmt.Sprintf("%x", ciphertext))
	fmt.Println("Type is", reflect.TypeOf(ciphertext))
	return fmt.Sprintf("%x", ciphertext)

	// dst := hex.EncodeToString(converted)
	// coursesEncrypted = append(coursesEncrypted, dst)
	// decoded, _ := hex.DecodeString(dst)
	//  fmt.Printf("%s\n", decoded)
}

func decodeText(code string) string{
	block, err := newCipherBlock("0")
	if err != nil {
		return "oopsies"
	}
	ciphertext, err := hex.DecodeString(code)
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "oopsies"
	}
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}
func EncryptString(plaintext string) string {
	block, err := newCipherBlock("0")
	if err != nil {
		return ""
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext)
}

// Decrypt will take in a key and a cipherHex (hex representation of
// the ciphertext) and decrypt it.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#NewCFBDecrypter


func encodeStudent(stu student) {
	// convertedS := []byte(stu.ID)
	// convertedN := []byte(stu.Name)
	// convertedY := []byte(stu.Year)
	// dstS := make([]byte, hex.EncodedLen(len(convertedS)))
	// dstN := make([]byte, hex.EncodedLen(len(convertedN)))
	// dstY := make([]byte, hex.EncodedLen(len(convertedY)))
	// hex.Encode(dstS, convertedS)
	// hex.Encode(dstN, convertedN)
	// hex.Encode(dstY, convertedY)
	// studentsEncrypted = append(studentsEncrypted, student{ID: string(dstS), Name: string(dstN), Year: string(dstY)})
}

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}
func main() {

	encoded1:= encodeCourseTest("CSSE 332")
	fmt.Println("How goes it?", encoded1)
	decoded1:= decodeText(encoded1)
	fmt.Println("DID it go back?", decoded1)
	encoded:= EncryptString("CSSE 332")
	fmt.Println("How goes it?", encoded)
	decoded:= DecryptString(encoded)
	fmt.Println("DID it go back?", decoded)
}
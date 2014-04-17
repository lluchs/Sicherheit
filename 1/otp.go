package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// Xors the input slices, generating one that is as big as the smaller of the
// two.
func xor(in1, in2 []byte) []byte {
	length := min(len(in1), len(in2))
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = in1[i] ^ in2[i]
	}
	return out
}

// Generates a random key of the given length.
func generateKey(size int64) ([]byte, error) {
	// Generate enough random bytes to make the key.
	rndbytes := make([]byte, size)
	_, err := rand.Read(rndbytes)
	if err != nil {
		return nil, err
	}
	return rndbytes, nil
}

// Returns a file's size.
func getFileSize(file *os.File) (size int64, err error) {
	stat, err := file.Stat()
	if err == nil {
		size = stat.Size()
	}
	return
}

// Opens a file and returns both the file and its size.
func openAndGetFileSize(filename string) (file *os.File, size int64, err error) {
	file, err = os.Open(filename)
	if err != nil {
		return
	}
	size, err = getFileSize(file)
	if err != nil {
		file.Close()
		file = nil
		return
	}
	return
}

// Encrypts the given file using the key.
// The key has to be large enough to prevent panics.
func encrypt(file *os.File, key []byte) error {
	// Open output file.
	outfile, err := os.Create(fmt.Sprintf("%s.otp", file.Name()))
	if err != nil {
		return err
	}
	defer outfile.Close()

	buffer := make([]byte, 100)
	keypos := 0
	for {
		count, err := file.Read(buffer)
		if err != nil {
			break
		}
		_, err = outfile.Write(xor(buffer, key[keypos:keypos + count]))
		if err != nil {
			return err
		}
		keypos += count
	}
	return nil
}

func main() {
	if (len(os.Args) < 2) {
		fmt.Println("Usage: otp <file1> [file2 [...]]")
		return
	}

	// Save all specified files.
	files := make([]*os.File, 0, len(os.Args) - 1)
	// Maximum file length.
	maxlen := int64(0)
	// Open the files.
	for i := 1; i < len(os.Args); i++ {
		file, length, err := openAndGetFileSize(os.Args[i])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		if maxlen < length {
			maxlen = length
		}
		files = append(files, file)
	}

	key, err := generateKey(maxlen)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		err = encrypt(file, key)
		if err != nil {
			panic(err)
		}
	}
}

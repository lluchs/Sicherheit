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

func encrypt(filename string) error {
	// Open the file.
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	key, err := generateKey(stat.Size())
	if err != nil {
		return err
	}
	// Open output file.
	outfile, err := os.Create(fmt.Sprintf("%s.otp", filename))
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
	if (len(os.Args) != 2) {
		fmt.Println("Usage: otp <file1>")
		return
	}

	err := encrypt(os.Args[1])
	if err != nil {
		panic(err)
	}
}
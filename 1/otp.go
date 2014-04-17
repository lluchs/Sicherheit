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
		_, err = outfile.Write(xor(buffer, key[keypos:keypos+count]))
		if err != nil {
			return err
		}
		keypos += count
	}
	return nil
}

func cmdEncrypt(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: otp encrypt <file1> [file2 [...]]")
		return
	}
	// Save all specified files.
	files := make([]*os.File, 0, len(args))
	// Maximum file length.
	maxlen := int64(0)
	// Open the files.
	for _, filename := range args {
		file, length, err := openAndGetFileSize(filename)
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

	// Save the key.
	keyfile, err := os.Create("key")
	if err != nil {
		panic(err)
	}
	defer keyfile.Close()
	_, err = keyfile.Write(key)
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

func cmdXor(args []string) {
	if len(args) != 3 {
		fmt.Println("Usage: otp xor <file1> <file2> <outfile>")
		return
	}
	file1, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer file1.Close()
	file2, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	outfile, err := os.Create(args[2])
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	buf1 := make([]byte, 100)
	buf2 := make([]byte, 100)
	for {
		_, err1 := file1.Read(buf1)
		_, err2 := file2.Read(buf2)
		_, err = outfile.Write(xor(buf1, buf2))
		if err != nil {
			panic(err)
		}
		if err1 != nil || err2 != nil {
			break
		}
	}
}

func main() {
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "encrypt":
		cmdEncrypt(os.Args[2:])
	case "xor":
		cmdXor(os.Args[2:])
	default:
		fmt.Println("Usage: otp <command>")
		fmt.Println("Possible commands: encrypt, xor")
		os.Exit(1)
	}
}

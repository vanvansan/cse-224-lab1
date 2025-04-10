package main

import (
	"log"
	"os"
	"encoding/binary"
	"fmt"
)

// Read a big-endian uint32 from a byte slice of length at least 4
func ReadBigEndianUint32(buffer []byte) uint32 {
	if len(buffer) < 4 {
    	panic("buffer too short to read uint32")
	}
	return binary.BigEndian.Uint32(buffer[:])
}

// Write a big-endian uint32 to a byte slice of length at least 4
func WriteBigEndianUint32(buffer []byte, num uint32) {
	if len(buffer) < 4 {
    	panic("buffer too short to write uint32")
	}
	binary.BigEndian.PutUint32(buffer, num)
}

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v inputfile outputfile\n", os.Args[0])
		// log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	// log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
	// Reading a big-endian uint32 from a byte slice
	var data [4]byte = [4]byte{0x00, 0x00, 0x00, 0x01}
	num := ReadBigEndianUint32(data[:])
	fmt.Println(num) // Output: 1

	// Writing a big-endian uint32 to a byte slice
	var buffer [4]byte
	WriteBigEndianUint32(buffer[:], num)
	fmt.Println(buffer) // Output: [0 0 0 1]

	// Attempting to write a big-endian uint32 to a 2-byte buffer
	var shortBuffer [2]byte
	WriteBigEndianUint32(shortBuffer[:], num) // This will cause a panic
	fmt.Println(shortBuffer) 
}

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	// "strings"
)

type Record struct{
	length uint32
	key [10]uint8
	value []uint8
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(r1, r2 *Record) bool

// recordSorter joins a By function and a slice of Planets to be sorted.
type recordSorter struct {
	records []Record
	by      func(p1, p2 *Record) bool // Closure used in the Less method.
}


// Len is part of sort.Interface.
func (s *recordSorter) Len() int {
	return len(s.records)
}

// Swap is part of sort.Interface.
func (s *recordSorter) Swap(i, j int) {
	s.records[i], s.records[j] = s.records[j], s.records[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *recordSorter) Less(i, j int) bool {
	return s.by(&s.records[i], &s.records[j])
}
func (by By) Sort(records []Record) {
	ps := &recordSorter{
		records: records,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}


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

// Read file and return the parsed record array 
func ReadFile(fileName string) []Record{
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file:%s, %s",fileName, err)
	}
	defer file.Close()
	var records []Record

	reader := bufio.NewReader(file)
	count := 0
	for {
		// fmt.Printf("Parsing in record  %d\n",count)
		count ++
		record := new(Record)
		err = binary.Read(reader, binary.BigEndian, &record.length)
		if err != nil {
			if err == io.EOF{
				// fmt.Println("Error1: ", err)
				break
			} else{
			fmt.Println("binary.Read failed:", err)
			os.Exit(1)
			}
		}

		err = binary.Read(reader, binary.BigEndian, &record.key)
		if err != nil {
			if err == io.EOF{
				break
			} else{
			fmt.Println("binary.Read failed:", err)
			os.Exit(2)
			}
		}
		record.value = make([]uint8, record.length - 10)
		err = binary.Read(reader, binary.BigEndian, &record.value)
		if err != nil {
			if err == io.EOF{
				break
			} else{
			fmt.Println("binary.Read failed:", err)
			os.Exit(3)
		}
		}
		records = append(records, *record)
		// fmt.Println(record)
	}
	return records
}


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v inputfile outputfile\n", os.Args[0])
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}
	key := func(r1, r2 *Record) bool {
		for i := 0; i < 10; i++{
			if r1.key[i] != r2.key[i]{
				return r1.key[i] < r2.key[i]
			}
		}
		return true
	}
	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
	records := ReadFile(os.Args[1])
	By(key).Sort(records)
	for _,e := range records{
		fmt.Println(e)
	}

	// Reading a big-endian uint32 from a byte slice
	// var data [4]byte = [4]byte{0x00, 0x00, 0x00, 0x01}
	// num := ReadBigEndianUint32(data[:])
	// fmt.Println(num) // Output: 1

	// // Writing a big-endian uint32 to a byte slice
	// var buffer [4]byte
	// WriteBigEndianUint32(buffer[:], num)
	// fmt.Println(buffer) // Output: [0 0 0 1]

	// // Attempting to write a big-endian uint32 to a 2-byte buffer
	// var shortBuffer [2]byte
	// WriteBigEndianUint32(shortBuffer[:], num) // This will cause a panic
	// fmt.Println(shortBuffer) 
}

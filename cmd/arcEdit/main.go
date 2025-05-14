package main

import (
	"flag"
	"fmt"
	"os"
)

var xor = []byte{
	0xea, 0x57, 0xbd, 0xef}

var xor2 = []byte{
	0xea, 0x57, 0xbd, 0xef,
	0x58, 0xb0, 0x60, 0x00,
	0xd1, 0x9a, 0x16, 0x4c,
	0x7b, 0xbe, 0xe7, 0xe9}

func decrypt0x10(file, out []byte) error {
	size := len(file) - 0x80
	lent := 0x1000
	writen := 0
	offset := 0x0
	for size > 0 {
		// Not accruate
		if size < 0x1000 {
			lent = size
		}
		stride := (lent / 0x10)

		fmt.Printf("FileSize:%x Stride:%x\n", size, stride)
		buf := []byte{}
		if stride&0xf == 0 {
			fmt.Println("TBI")
			break

			// for j := range 0x10 {
			// 	for range stride {
			// 	}
			// }
		} else {
			iz := 0
			for j := range 0x10 {
				for range stride {
					// this+0x4 = inBuf ^ EA
					index := file[(iz)+0x80]
					key := xor2[j%0x10]
					buf = append(buf, index^key)
					iz++
				}
			}
		}

		for i, x := range buf {
			if i%0x10 == 0 {
				fmt.Printf("\n%.3X - ", i)
			} else if i%4 == 0 {
				fmt.Print("| ")
			}
			fmt.Printf("\033[9%dm%.2X \033[0m", (i%stride)%0x8, x)
		}
		fmt.Println()

		for i := 0x0; i < stride; i++ {
			// index := i * 4
			for j := 0; j < 0x10; j++ {
				// index := j * (stride * 2)
				// fmt.Printf("%x\n", index)
				out[(i+j)+offset+0x80] = buf[i+(j*stride)]
				writen++
			}
			offset += 15
			// out[(index+1)+offset+0x80] = buf[i+1+offset+0x80]
		}
		offset += lent
		size -= lent
	}

	for i, x := range out[0x80 : writen+0x80] {
		if i%0x10 == 0 {
			fmt.Println()
		}
		fmt.Printf("%.2x ", x)
	}
	fmt.Println()
	return nil
	// return fmt.Errorf("testing")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	boolEncrypt := flag.Bool("encrypt", false, "Set to encrypt file, otherwise it'll decrypt")
	inFileName := flag.String("input", "", "[Required] File name for input")
	outFileName := flag.String("out", "", "File name for dest, empty will replace input file")
	flag.Parse()

	if *inFileName == "" {
		flag.Usage()
		return
	}

	file, e := os.ReadFile(*inFileName)
	if e != nil {
		fmt.Println(e)
		return
	}

	var out []byte
	out = make([]byte, len(file))
	for x := 0; x < 0x80; x++ {
		out[x] = file[x]
	}

	//if reflect.DeepEqual(file[0x54:0x57], []byte{'D', 'X', 'T', '5'}) {}
	var err error
	switch file[0x57] {
	case '5':
		err = decrypt0x10(file, out)
	case 0x0:
		if *boolEncrypt {
			encrypt0x4(file, out)
		} else {
			decrypt0x4(file, out)
		}
	default:
		fmt.Println("Not Implemented")
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if *outFileName == "" {
		outFileName = inFileName
	}
	f, _ := os.Create(*outFileName)
	f.Write(out)

	fmt.Println("Finished " + *outFileName)
}

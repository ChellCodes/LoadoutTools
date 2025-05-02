package main

import (
	"flag"
	"fmt"
	"os"
)

var xor = []byte{
	0xea, 0x57, 0xbd, 0xef}

var xor2 = []byte{
	0xea, 0xea, 0x15, 0x15,
	0x15, 0x15, 0x15, 0x15,
	0x1D, 0xf1, 0x3d, 0xc9,
	0xbf, 0xbf, 0xbf, 0xbf}

func decrypt0x10(file, out []byte) error {
	size := len(file) - 0x80
	fmt.Printf("%x %x\n", size, size/0x10)
	buf := []byte{}
	for j := range 0x10 {
		for i := range size / 0x10 {
			// this+0x4 = inBuf ^ EA
			buf = append(buf, file[i+j+0x80]^xor2[j%0x10])
		}
	}

	for i, x := range buf {
		if i%0x10 == 0 {
			fmt.Println()
		}
		fmt.Printf("%.2x ", x)
	}
	fmt.Println()

	writen := 0
	offset := 0x0
	for size > 0 {
		lent := 0x1000
		if size < 0x1000 {
			lent = size
		}
		stride := (lent / 0x10)
		for i := 0x0; i < stride; i++ {
			// index := i * 4
			for j := 0; j < 0x10; j++ {
				// index := j * (stride * 2)
				// fmt.Printf("%x\n", index)
				out[(i+j)+offset+0x80] = buf[j*stride]
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

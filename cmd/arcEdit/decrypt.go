package main

import (
	"flag"
	"fmt"
	"os"
)

var xor = []byte{0xea, 0x57, 0xbd, 0xef}

func decrypt(file, out []byte) {
	size := len(file) - 0x80
	offset := 0x0
	for size > 0 {
		lent := 0x1000
		if size < 0x1000 {
			lent = size
		}
		stride := (lent / 0x4)
		for i := 0x0; i < stride; i++ {
			index := i * 4
			out[(index)+offset+0x80] = file[i+offset+0x80] ^ xor[index%4]
			out[(index+1)+offset+0x80] = file[i+offset+0x80+stride] ^ xor[(index+1)%4]
			out[(index+2)+offset+0x80] = file[i+offset+0x80+(stride*2)] ^ xor[(index+2)%4]
			out[(index+3)+offset+0x80] = file[i+offset+0x80+(stride*3)] ^ xor[(index+3)%4]
		}
		offset += lent
		size -= lent
	}
}

func encrypt(file, out []byte) {
	size := len(file) - 0x80
	offset := 0x0
	for size > 0 {
		lent := 0x1000
		if size < 0x1000 {
			lent = size
		}
		stride := (lent / 0x4)
		for i := 0x0; i < stride; i++ {
			index := i * 4
			out[i+offset+0x80] = file[(index)+offset+0x80] ^ xor[index%4]
			out[i+offset+0x80+stride] = file[(index+1)+offset+0x80] ^ xor[(index+1)%4]
			out[i+offset+0x80+(stride*2)] = file[(index+2)+offset+0x80] ^ xor[(index+2)%4]
			out[i+offset+0x80+(stride*3)] = file[(index+3)+offset+0x80] ^ xor[(index+3)%4]
		}
		offset += lent
		size -= lent
	}
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

	file, err := os.ReadFile(*inFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	var out []byte
	out = make([]byte, len(file))
	for x := 0; x < 0x80; x++ {
		out[x] = file[x]
	}

	if *boolEncrypt {
		encrypt(file, out)
	} else {
		decrypt(file, out)
	}

	if *outFileName == "" {
		outFileName = inFileName
	}
	f, _ := os.Create(*outFileName)
	f.Write(out)

	fmt.Println("Finished " + *outFileName)
}

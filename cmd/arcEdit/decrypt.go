package main

import (
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
	file, err := os.ReadFile("in.dds")
	if err != nil {
		fmt.Println(err)
		return
	}

	var out []byte
	out = make([]byte, len(file))
	for x := 0; x < 0x80; x++ {
		out[x] = file[x]
	}
	decrypt(file, out)
	// encrypt(file, out)
	f, _ := os.Create("out.dds")
	f.Write(out)

	fmt.Println("Finished")
}

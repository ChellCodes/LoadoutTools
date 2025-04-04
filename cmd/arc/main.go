package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type arcFile struct {
	index uint32
	magic []byte
	//empty int64
	fileNameLen uint32
	fileName    string
	unknown0    uint32
	unknown1    uint32
	OffsetCount uint32
	DataBlock   []byte
	Assets      []asset
}

type asset struct {
	Offset  uint32
	DataLen uint64
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func main() {
	fmt.Println("Wello Horld!")

	loPath := "N:/SteamLibrary/steamapps/common/Loadout/Data/"
	file, _ := os.ReadFile(loPath + "index.ind")

	numberOfFiles := binary.LittleEndian.Uint32(file[0x10:0x14])
	fmt.Printf("Num of files: %d\n", numberOfFiles)

	files := []arcFile{}
	offset := 0x14
	for i := range numberOfFiles {
		data, n := getArcData(file[offset:])
		offset += int(n)
		files = append(files, data)
		fmt.Printf("%6d| File: %s | Offsets: %4d", i, data.fileName, data.OffsetCount)
		if (i % 2) == 1 {
			fmt.Println()
		}
	}

	index := getUserInt("\nSelect File: ")
	selected := files[index]
	fmt.Printf("Selected: %s OffsetCount: %d\n", selected.fileName, selected.OffsetCount)

	dataFile, _ := os.ReadFile(loPath + files[index].fileName)
	for i := range selected.OffsetCount {
		selectedAsset := selected.Assets[i]
		ind := clen(dataFile[selectedAsset.Offset:selectedAsset.Offset+uint32(selectedAsset.DataLen)]) + 1

		width := binary.LittleEndian.Uint32(dataFile[selectedAsset.Offset+uint32(ind) : selectedAsset.Offset+uint32(ind)+0x4])
		height := binary.LittleEndian.Uint32(dataFile[selectedAsset.Offset+uint32(ind)+0x4 : selectedAsset.Offset+uint32(ind)+0x8])
		// mipmapCount := binary.LittleEndian.Uint32(dataFile[selectedAsset.Offset+uint32(ind)+0xc : selectedAsset.Offset+uint32(ind)+0x10])
		// return
		if width == 1024 && height == 1024 {
			// fmt.Printf("%4d %5d %5d %3d \n", i, width, height, mipmapCount)
			outName := fmt.Sprintf("out\\%s-%d-0x%x.dds", strings.Split(selected.fileName, ".")[0], i, selectedAsset.Offset)
			fi, _ := os.Create(outName)
			texture := dataFile[selectedAsset.Offset : selectedAsset.Offset+uint32(selectedAsset.DataLen)]
			_, a, _ := bytes.Cut(texture, []byte{'D', 'D', 'S'})
			fi.Write(append([]byte{'D', 'D', 'S'}, a...))
			fmt.Printf("Wrote File: %s\n", outName)
		}
	}

	return
	offsetIndex := getUserInt("Select Offset: ")
	selectedAsset := selected.Assets[offsetIndex]
	fmt.Printf("Offset into %s: 0x%x Len: 0x%x\n", selected.fileName, selectedAsset.Offset, selectedAsset.DataLen)

	mode := getUserInt("Patch or Dump (1 or 2): ")
	if mode == 1 {
		fmt.Println("TBI")
	}

	if mode == 2 {
		outName := fmt.Sprintf("%s-%d-0x%x", strings.Split(selected.fileName, ".")[0], offsetIndex, selectedAsset.Offset)
		fi, _ := os.Create(outName)
		fi.Write(dataFile[selectedAsset.Offset : selectedAsset.Offset+uint32(selectedAsset.DataLen)])
		fmt.Printf("Wrote File: %s\n", outName)
	}
}

func getUserInt(printMsg string) int {
	fmt.Print(printMsg)
	var userInt int
	_, err := fmt.Scan(&userInt)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return userInt
}

func getArcData(data []byte) (arcData arcFile, n uint32) {
	var out arcFile
	out.index = binary.LittleEndian.Uint32(data[:0x04])
	out.magic = data[0x04:0x0C]
	out.fileNameLen = binary.LittleEndian.Uint32(data[0x14:0x18])

	offset := 0x18 + out.fileNameLen
	out.fileName = string(data[0x18:offset])
	out.unknown0 = binary.LittleEndian.Uint32(data[offset : offset+0x04])
	out.unknown1 = binary.LittleEndian.Uint32(data[offset+0x04 : offset+0x08])
	offset += 0x08

	out.OffsetCount = binary.LittleEndian.Uint32(data[offset : offset+0x04])
	out.DataBlock = data[offset+0x04 : offset+0x4+(out.OffsetCount*0x4)]
	offset += 0x04 + (out.OffsetCount * 0x4)

	for i := uint32(0); i < out.OffsetCount; i++ {
		out.Assets = append(out.Assets, asset{
			Offset:  binary.LittleEndian.Uint32(data[offset : offset+0x4]),
			DataLen: binary.LittleEndian.Uint64(data[offset+0x4 : offset+0xC]),
		})
		offset += 0xC
	}

	return out, offset
}

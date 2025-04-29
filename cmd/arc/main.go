package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type indFiles []indFile

type indFile struct {
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

func (i indFiles) toBytes() []byte {
	outData := []byte{0xD5, 0x11, 0x0D, 0x60, 0xEB, 0xC7, 0x3A, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00}
	for _, file := range i {
		outData = binary.LittleEndian.AppendUint32(outData, file.index)
		outData = append(outData, file.magic...)
		outData = binary.LittleEndian.AppendUint64(outData, 0)
		outData = binary.LittleEndian.AppendUint32(outData, file.fileNameLen)
		outData = append(outData, []byte(file.fileName)...)
		outData = binary.LittleEndian.AppendUint32(outData, file.unknown0)
		outData = binary.LittleEndian.AppendUint32(outData, file.unknown1)
		outData = binary.LittleEndian.AppendUint32(outData, file.OffsetCount)
		outData = append(outData, file.DataBlock...)
		for _, a := range file.Assets {
			outData = binary.LittleEndian.AppendUint32(outData, a.Offset)
			outData = binary.LittleEndian.AppendUint64(outData, a.DataLen)
		}
	}
	return outData
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

var files = indFiles{}

func main() {
	fmt.Println("Wello Horld!")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	loPath := flag.String("loadoutDir", "C:/Program Files (x86)/Steam/steamapps/common/Loadout/Data/", "Folder Path to Loadout Data dir")
	flag.Parse()

	file, err := os.ReadFile(*loPath + "index.ind")
	if err != nil {
		fmt.Println(err)
		return
	}

	numberOfFiles := binary.LittleEndian.Uint32(file[0x10:0x14])
	fmt.Printf("Num of files: %d\n", numberOfFiles)

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

	dataFile, err := os.ReadFile(*loPath + files[index].fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	offsetIndex := getUserInt("Select Offset (-1 to dump all): ")
	if offsetIndex == -1 {
		for i := uint32(0); i < selected.OffsetCount; i++ {
			selectedAsset := selected.Assets[i]
			outName := fmt.Sprintf("out/%s-%d-0x%x", strings.Split(selected.fileName, ".")[0], i, selectedAsset.Offset)

			if reflect.DeepEqual(dataFile[:4], []byte{'T', 'X', 'F', 'L'}) {
				outName += ".dds"
				outFile, _ := os.Create(outName)

				texture := dataFile[selectedAsset.Offset : selectedAsset.Offset+uint32(selectedAsset.DataLen)]
				_, a, _ := bytes.Cut(texture, []byte{'D', 'D', 'S'})
				if a[0x54-0x3] != 0 {
					continue
				}
				outFile.Write(append([]byte{'D', 'D', 'S'}, a...))
			} else if reflect.DeepEqual(dataFile[:4], []byte{'D', 'D', 'S', ' '}) {
				outName += ".dds"
				outFile, _ := os.Create(outName)
				outFile.Write(dataFile[selectedAsset.Offset : selectedAsset.Offset+uint32(selectedAsset.DataLen)])
			} else {
				outName += ".bin"
				outFile, _ := os.Create(outName)
				outFile.Write(dataFile[selectedAsset.Offset : selectedAsset.Offset+uint32(selectedAsset.DataLen)])
			}
			fmt.Printf("Wrote File: %s\n", outName)
		}
		return
	}

	selectedAsset := selected.Assets[offsetIndex]
	fmt.Printf("Offset into %s: 0x%x Len: 0x%x\n", selected.fileName, selectedAsset.Offset, selectedAsset.DataLen)

	mode := getUserInt("Patch or Dump (1 or 2): ")
	if mode == 1 {
		patch(dataFile, selectedAsset, index, offsetIndex)
	}

	if mode == 2 {
		outName := fmt.Sprintf("%s-%d-0x%x", strings.Split(selected.fileName, ".")[0], offsetIndex, selectedAsset.Offset)
		fi, _ := os.Create(outName)
		fi.Write(dataFile[selectedAsset.Offset : selectedAsset.Offset+uint32(selectedAsset.DataLen)])
		fmt.Printf("Wrote File: %s\n", outName)
	}
}

// Patch ARC
// Recalc Ind offsets

func patch(arcF []byte, sel asset, fileIndex, assetIndex int) {
	var filePath string

	fmt.Print("File to replace with: ")
	fmt.Scan(&filePath)
	patchData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	outArc := append(arcF[:sel.Offset], patchData...)
	outArc = append(outArc, arcF[sel.Offset+uint32(sel.DataLen):]...)
	outFile, _ := os.Create(files[fileIndex].fileName)
	outFile.Write(outArc)
	fmt.Println("Wrote ARC File!")

	outIndFile, _ := os.Create("index.ind")

	offsetAdjust := len(patchData) - int(sel.DataLen)
	fmt.Printf("len of new data = 0x%x\nOffest Adjustment = 0x%x\n", len(patchData), offsetAdjust)
	for index := range files[fileIndex].Assets {
		if index < assetIndex {
			continue
		}
		if index == assetIndex {
			files[fileIndex].Assets[index].DataLen = uint64(len(patchData))
			continue
		}
		files[fileIndex].Assets[index].Offset += uint32(offsetAdjust)
	}

	outIndFile.Write(files.toBytes())
	fmt.Println("Wrote ind File!")
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

func getArcData(data []byte) (arcData indFile, n uint32) {
	var out indFile
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

func images() {
	/*
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
		}*/
}

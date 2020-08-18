package main


import (
    "fmt"
    "os"
    "io"
    "encoding/binary"
)

// https://en.wikipedia.org/wiki/BMP_file_format (Header information)
// http://www.ece.ualberta.ca/~elliott/ee552/studentAppNotes/2003_w/misc/bmp_file_format/bmp_file_format.htm (Additional header information)
type header struct {
    signature string
    bitmapSize uint32
    reserved1 uint32
    dataOffset uint32

    size uint32
    width uint32
    height uint32
    planes uint16
    bitsPerPixel uint16
    compressionType uint32
    imageSize uint32
    xPixelsPerM uint32
    yPixelsPerM uint32
    totalColour uint32
    importantColours uint32

}

type pixel struct {
    red uint8
    green uint8
    blue uint8
}




func read16string(f *os.File) string {
    buf16 := make([]byte, 2)

    _, err := f.Read(buf16)
    if err != nil {
        if err != io.EOF{
            panic(err)
        }
    }
    return string(buf16)
}

// Used for reading in the RGB at the same time from the bitmap file
func read24buff(f *os.File) []byte {
    buf24:= make([]byte, 3)

    // http://zetcode.com/golang/readfile/ (Derived from)
    _, err := f.Read(buf24)
    if err != nil {
        if err != io.EOF{
            panic(err)
        }
    }

    return buf24
}

func read16buff(f *os.File) uint16 {
    buf16 := make([]byte, 2)

    // http://zetcode.com/golang/readfile/ (Derived from)
    _, err := f.Read(buf16)
    if err != nil {
        if err != io.EOF{
            panic(err)
        }
    }
    // https://stackoverflow.com/questions/15848830/decoding-data-from-a-byte-slice-to-uint32

    res := binary.LittleEndian.Uint16(buf16)
    return res
}

func read32buff(f *os.File) uint32 {
    buf32 := make([]byte, 4)

    // http://zetcode.com/golang/readfile/ (Derived from)
    _, err := f.Read(buf32)
    if err != nil {
        if err != io.EOF{
            panic(err)
        }
    }
    // https://stackoverflow.com/questions/15848830/decoding-data-from-a-byte-slice-to-uint32
    res := binary.LittleEndian.Uint32(buf32)
    return res
}

func printHeader(imageHeader header) {
    fmt.Printf("Signature = %v\n",  imageHeader.signature)
    fmt.Printf("Bitmap Size = %v\n",  imageHeader.bitmapSize)
    fmt.Printf("Reserved 1 = %v\n",  imageHeader.reserved1)
    fmt.Printf("Data Offset = %v\n",  imageHeader.dataOffset)

    fmt.Printf("Size = %v\n",  imageHeader.size)
    fmt.Printf("Width = %v\n",  imageHeader.width)
    fmt.Printf("Height = %v\n",  imageHeader.height)
    fmt.Printf("Planes = %v\n",  imageHeader.planes)
    fmt.Printf("Bits Per Pixel = %v\n",  imageHeader.bitsPerPixel)
    fmt.Printf("Compression Type = %v\n",  imageHeader.compressionType)
    fmt.Printf("Image Size = %v\n",  imageHeader.imageSize)
    fmt.Printf("X Pixels Per M = %v\n",  imageHeader.xPixelsPerM)
    fmt.Printf("Y Pixels Per M = %v\n",  imageHeader.yPixelsPerM)
    fmt.Printf("Total Colours = %v\n",  imageHeader.totalColour)
    fmt.Printf("Important Colours = %v\n",  imageHeader.importantColours)
}



func write16uint(toBeWritten uint16,f *os.File) {
    buf16 := make([]byte, 2)
	// https://stackoverflow.com/questions/35371385/how-can-i-convert-an-int64-into-a-byte-array-in-go
	binary.LittleEndian.PutUint16(buf16, toBeWritten)

    _, err := f.Write(buf16)
    if err != nil {
        if err != io.EOF{
            panic(err)
        }
    }
}


func write32uint(toBeWritten uint32,f *os.File) {
    buf32 := make([]byte, 4)
	// https://stackoverflow.com/questions/35371385/how-can-i-convert-an-int64-into-a-byte-array-in-go
	binary.LittleEndian.PutUint32(buf32, toBeWritten)

    _, err := f.Write(buf32)
    if err != nil {
        if err != io.EOF{
            panic(err)
        }
    }
}

func readData(imageHeader header) []pixel {
	readName := "forGoData.CMP"

	size := imageHeader.width * imageHeader.height
	pixelArr := make([]pixel, size)

	f, err := os.Open(readName)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := uint32(0); i < size; i++ {
		RGB := read24buff(f)
		pixelArr[i].red = RGB[0]
		pixelArr[i].green = RGB[1]
		pixelArr[i].blue = RGB[2]
	}
	return pixelArr
}

func readCMPHeader()(header, []pixel){
	readName := "forGoHeader.CMP"

	var imageHeader header

	f, err := os.Open(readName)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	imageHeader.signature = read16string(f)
	imageHeader.bitmapSize = read32buff(f)
	imageHeader.reserved1 = read32buff(f)
	imageHeader.dataOffset = read32buff(f)

	imageHeader.size = read32buff(f)
	imageHeader.width = read32buff(f)
	imageHeader.height = read32buff(f)
	imageHeader.planes = read16buff(f)
	imageHeader.bitsPerPixel = read16buff(f)
	imageHeader.compressionType = read32buff(f)
	imageHeader.imageSize = read32buff(f)
	imageHeader.xPixelsPerM = read32buff(f)
	imageHeader.yPixelsPerM = read32buff(f)
	imageHeader.totalColour = read32buff(f)
	imageHeader.importantColours = read32buff(f)

	printHeader(imageHeader)

	data := readData(imageHeader)
	//fmt.Printf("Data = %v\b", data)
	return imageHeader, data
}

func writeBitmap(saveName string, imageHeader header, data []pixel) {
	f, err := os.Create(saveName)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, e := f.WriteString(imageHeader.signature)
	if e != nil{
		panic(e)
	}

	write32uint(imageHeader.bitmapSize, f)
	write32uint(imageHeader.reserved1, f)
	write32uint(imageHeader.dataOffset, f)
	write32uint(imageHeader.size, f)
	write32uint(imageHeader.width, f)
	write32uint(imageHeader.height, f)
	write16uint(imageHeader.planes, f)
	write16uint(imageHeader.bitsPerPixel , f)
	write32uint(imageHeader.compressionType, f)
	write32uint(imageHeader.imageSize, f)
	write32uint(imageHeader.xPixelsPerM, f)
	write32uint(imageHeader.yPixelsPerM, f)
	write32uint(imageHeader.totalColour, f)
	write32uint(imageHeader.importantColours, f)

	width := imageHeader.width
	height := imageHeader.height

	extraOffset := make([]byte, imageHeader.dataOffset - 54) // creates enough zeros to write to get to the offset
	//extraOffset := make([]byte, imageHeader.dataOffset) // creates enough zeros to write to get to the offset
	f.Write(extraOffset)

	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			b := []byte{data[(y * width) + x].blue, data[(y * width) + x].green, data[(y * width) + x].red}
			f.Write(b)
		}
		// Should add the required %4 for width
		extra := make([]byte, width%4)
		f.Write(extra)
	}
}

func main() {
	saveName :=  os.Args[1]

	imageHeader, data := readCMPHeader()
	writeBitmap(saveName, imageHeader, data)

}


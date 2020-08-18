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

func readData(imageHeader header, f *os.File)[]pixel {
    // https://ispycode.com/GO/Files-And-Directories/Seek-Positions-in-File

    _, err := f.Seek(int64(imageHeader.dataOffset), 0)
    if err != nil {
        panic(err)
    }

	size := imageHeader.width * imageHeader.height

	pixelArr := make([]pixel, size)

	width := imageHeader.width
	height := imageHeader.height
	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width ; x++{
			RGB := read24buff(f)
			// https://en.wikipedia.org/wiki/BMP_file_format#Pixel_format (Contains order of that it should be read in)
			pixelArr[x+(y*width)].blue = RGB[0]
			pixelArr[x+(y*width)].green = RGB[1]
			pixelArr[x+(y*width)].red = RGB[2]
		}
		// https://en.wikipedia.org/wiki/BMP_file_format#Pixel_storage (Should end on a witdth divisible by 4, so this adjusts for that)
		f.Seek(int64(width%4), 1) // Adjusts current file location based on width of bitmap
	}
	return pixelArr
}

func loadBMP()(header, []pixel) {
    var imageHeader header

    // https://gobyexample.com/command-line-arguments
    path := os.Args[1]
    //fmt.Printf(path)

    //https://appdividend.com/2019/12/11/how-to-open-file-in-golang-go-file-open-example/
    // https://stackoverflow.com/questions/14514201/how-to-read-a-binary-file-in-go
    f, err := os.Open(path)
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

    data := readData(imageHeader, f)
	//fmt.Printf("%v\n", data)
	//n := imageHeader.width*imageHeader.height

	return imageHeader, data
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

func saveCMPHeader(imageHeader header, f *os.File) {
	_, err := f.WriteString(imageHeader.signature)
	if err != nil{
		panic(err)
	}

	write32uint(imageHeader.bitmapSize, f)
	write32uint(imageHeader.reserved1, f)
	write32uint(imageHeader.dataOffset, f)
	write32uint(imageHeader.size, f)
	write32uint(imageHeader.width, f)
	write32uint(imageHeader.height, f)
	write16uint(imageHeader.planes, f)
	write16uint(imageHeader.bitsPerPixel, f)
	write32uint(imageHeader.compressionType, f)
	write32uint(imageHeader.imageSize, f)
	write32uint(imageHeader.xPixelsPerM, f)
	write32uint(imageHeader.yPixelsPerM, f)
	write32uint(imageHeader.totalColour, f)
	write32uint(imageHeader.importantColours, f)
}

func saveCMPData(data []pixel, f* os.File){
	for _, pixelData := range data {
		//fmt.Printf("Pixel Data = %v\n", pixelData)
		b := []byte{pixelData.red, pixelData.green, pixelData.blue}
		_, err := f.Write(b)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("b = %v\n", b)
	}
}

func saveCMP(data []pixel, imageHeader header, saveName string) {
	// https://gobyexample.com/writing-files
	headerCMPName := saveName + "Header.CMP"
	dataCMPName := saveName + "Data.CMP"

	fHeader, err := os.Create(headerCMPName)
	if err != nil {
		panic(err)
	}

	defer fHeader.Close()

	saveCMPHeader(imageHeader, fHeader)

	fData, e := os.Create(dataCMPName)
	if e != nil {
		panic(e)
	}

	defer fData.Close()
	saveCMPData(data, fData)

}

func main() {
	saveName := "forC"
	// imageHeader holds bitmap header in header struct
	// Data stores pixel data in a []pixel
	imageHeader, data := loadBMP()
	saveCMP(data, imageHeader, saveName)

}














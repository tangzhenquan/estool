package filereader

import (
	"bufio"
	"esTool/pkg/file"
	"esTool/reader"
	"os"
	"strings"
)

// SeekInfo represents arguments to `os.Seek`
/*type SeekInfo struct {
	Offset int64
	Whence int // os.SEEK_*
}*/

type LineReader struct {
	file   *os.File
	reader *bufio.Reader
}

func NewLineReader(filepath string, maxLineSize int) (*LineReader, error) {
	f, err := file.ReadOpen(filepath)
	if err != nil {
		return nil, err
	}
	lr := new(LineReader)

	lr.file = f
	if maxLineSize > 0 {
		lr.reader = bufio.NewReaderSize(lr.file, maxLineSize)
	} else {
		lr.reader = bufio.NewReader(lr.file)
	}
	return lr, nil

}

func (lr *LineReader) readLine() (string, error) {
	line, err := lr.reader.ReadString('\n')
	if err != nil {
		return line, err
	}
	line = strings.TrimRight(line, "\n")

	return line, err
}

func (lr *LineReader) Next() (reader.Line, error) {
	// Read line by line.
	line, err := lr.readLine()
	if err == nil {
		return reader.NewLine(line), nil
	} else {
		return reader.Line{}, err
	}
}

func (lr *LineReader) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

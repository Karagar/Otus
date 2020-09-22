package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	intOffset := int(offset)
	intLimit := int(limit)
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(toPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer toFile.Close()

	fileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}
	fileSize := int(fileInfo.Size())
	if intOffset >= fileSize {
		return ErrOffsetExceedsFileSize
	}
	if fileSize == 0 {
		return ErrUnsupportedFile
	}
	if intLimit > fileSize-intOffset || intLimit == 0 {
		intLimit = fileSize - intOffset
	}

	bufferSize := 5 << 20
	buffer := make([]byte, bufferSize)
	bar := pb.StartNew(intLimit)
	written := 0
	currentOffset := intOffset

	for written != intLimit {
		readed, err := fromFile.ReadAt(buffer, int64(currentOffset))
		if err != nil && err != io.EOF {
			return err
		}
		currentOffset += readed
		if currentOffset > intLimit+intOffset {
			readed = intLimit + intOffset + readed - currentOffset
		}
		writed, err := toFile.Write(buffer[:readed])
		written += writed
		bar.Add(writed)

		if err != nil {
			return err
		}
	}
	bar.Finish()

	return nil
}

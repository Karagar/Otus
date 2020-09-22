package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const chunkSize = 1000

func compare(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else {
				return false
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func TestCopy(t *testing.T) {
	err := Copy("testdata/input.txt", "testdata/temp.txt", 0, 0)
	require.Nil(t, err)
	isEqual := compare("testdata/temp.txt", "testdata/out_offset0_limit0.txt")
	require.True(t, isEqual)

	err = Copy("testdata/input.txt", "testdata/temp.txt", 0, 10)
	require.Nil(t, err)
	isEqual = compare("testdata/temp.txt", "testdata/out_offset0_limit10.txt")
	require.True(t, isEqual)

	err = Copy("testdata/input.txt", "testdata/temp.txt", 0, 1000)
	require.Nil(t, err)
	isEqual = compare("testdata/temp.txt", "testdata/out_offset0_limit1000.txt")
	require.True(t, isEqual)

	err = Copy("testdata/input.txt", "testdata/temp.txt", 0, 10000)
	require.Nil(t, err)
	isEqual = compare("testdata/temp.txt", "testdata/out_offset0_limit10000.txt")
	require.True(t, isEqual)

	err = Copy("testdata/input.txt", "testdata/temp.txt", 100, 1000)
	require.Nil(t, err)
	isEqual = compare("testdata/temp.txt", "testdata/out_offset100_limit1000.txt")
	require.True(t, isEqual)

	err = Copy("testdata/input.txt", "testdata/temp.txt", 6000, 1000)
	require.Nil(t, err)
	isEqual = compare("testdata/temp.txt", "testdata/out_offset6000_limit1000.txt")
	require.True(t, isEqual)

	err = os.Remove("testdata/temp.txt")
	require.Nil(t, err)
}

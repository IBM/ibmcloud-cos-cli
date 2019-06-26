package utils

import (
	"bytes"

	"github.com/IBM/ibm-cos-sdk-go/aws"
)

type stringWrap struct {
	*bytes.Reader
	closeStatus *bool
}

func (sw *stringWrap) Close() error {
	if sw.closeStatus != nil {
		*(sw.closeStatus) = true
	}
	return nil
}

func WrapString(input string, closeStatus *bool) ReadSeekerCloser {
	return &stringWrap{
		Reader:      bytes.NewReader([]byte(input)),
		closeStatus: closeStatus,
	}
}

type writeToString struct {
	aws.WriteAtBuffer
	output      *string
	closeStatus *bool
}

func (sw *writeToString) WriteAt(p []byte, off int64) (n int, err error) {
	n, err = sw.WriteAtBuffer.WriteAt(p, off)
	if sw.output != nil {
		*sw.output = string(sw.WriteAtBuffer.Bytes())
	}
	return
}

func (sw *writeToString) Write(p []byte) (n int, err error) {
	return sw.WriteAt(p, int64(len(sw.WriteAtBuffer.Bytes())))
}

func (sw *writeToString) Close() error {
	if sw.closeStatus != nil {
		*(sw.closeStatus) = true
	}
	return nil
}

func WriteToString(output *string, closeStatus *bool) WriteCloser {
	return &writeToString{
		output:      output,
		closeStatus: closeStatus,
	}
}

package utils

import (
	"bytes"
	"io"
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
	output      *string
	closeStatus *bool
}

func (sw *writeToString) Write(p []byte) (n int, err error) {
	if sw.output != nil {
		*sw.output = *sw.output + string(p)
	}
	return len(p), nil
}

func (sw *writeToString) Close() error {
	if sw.closeStatus != nil {
		*(sw.closeStatus) = true
	}
	return nil
}

func WriteToString(output *string, closeStatus *bool) io.WriteCloser {
	return &writeToString{
		output:      output,
		closeStatus: closeStatus,
	}
}

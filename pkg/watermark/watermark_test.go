package watermark

import (
	"github.com/golang/mock/gomock"
	"github.com/pefish/go-test-assert"
	mock_io "github.com/pefish/go-watermarker/mock/mock-io"
	"image"
	"testing"
)


func TestWatermark_markJpgAndPng(t *testing.T) {
	w := NewWatermark("haha")

	ctl := gomock.NewController(t)
	writerInstance := mock_io.NewMockWriter(ctl)
	writerInstance.EXPECT().Write(gomock.Any()).DoAndReturn(func(a []byte) (int, error) {
		test.Equal(t, 2282, len(a))
		return len(a), nil
	})

	a := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	err := w.markJpgAndPng(a, writerInstance)
	test.Equal(t, nil, err)
}
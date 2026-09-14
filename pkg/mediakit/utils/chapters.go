package utils

import (
	"io"
	"time"

	"github.com/go-logr/logr"
	"github.com/remko/go-mkvparse"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/handlers"
	"github.com/neptune-media/MediaKit-go/pkg/mediakit/utils/timing"
)

func ReadVideoChapters(r io.ReadSeeker, logger logr.Logger) (mediakit.ChapterList, error) {
	defer timing.NewFuncDurationPrinter(time.Now(), "finished reading chapters")(logger)

	chapterHandler := &handlers.ChapterHandler{
		Chapters: make([]mediakit.Chapter, 0),
		Logger:   logger,
	}
	segmentInfoHandler := &handlers.SegmentInfoHandler{
		Logger:      logger,
		SegmentInfo: new(mediakit.SegmentInfo),
	}

	h := mkvparse.NewHandlerChain(chapterHandler, segmentInfoHandler)
	err := mkvparse.ParseSections(r, h, mkvparse.InfoElement, mkvparse.ChaptersElement)
	if err != nil {
		return nil, err
	}

	chapters := make(mediakit.ChapterList, len(chapterHandler.Chapters))
	for i, c := range chapterHandler.Chapters {
		chapters[i] = c.WithTimescale(segmentInfoHandler.SegmentInfo.Timescale)
	}

	return chapters, err
}

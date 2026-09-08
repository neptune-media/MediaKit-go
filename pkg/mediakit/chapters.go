package mediakit

import (
	"io"
	"time"

	"github.com/remko/go-mkvparse"
	"go.uber.org/zap"
)

// Chapter represents a chapter within a video file
type Chapter struct {
	// ID is the chapter id provided by the video file
	ID uint64

	// TimeStart is the start time of the chapter.  The units
	// are determined by the timescale provided by the
	// video.
	TimeStart int64

	// TimeEnd is the end time of the chapter.  The units
	// are determined by the timescale provided by the
	// video.
	TimeEnd int64

	// Enabled is a flag indicated if the chapter is enabled
	// in the video file.
	Enabled bool

	// Title is the encoded chapter title name
	Title string
}

func (c Chapter) EndTime() time.Duration {
	return time.Duration(c.TimeEnd) * time.Millisecond
}

func (c Chapter) Runtime() time.Duration {
	return c.EndTime() - c.StartTime()
}

func (c Chapter) StartTime() time.Duration {
	return time.Duration(c.TimeStart) * time.Millisecond
}

func (c Chapter) WithTimescale(ts time.Duration) Chapter {
	t := c

	t.TimeEnd = t.TimeEnd / int64(ts)
	t.TimeStart = t.TimeStart / int64(ts)

	return t
}

type ChapterHandler struct {
	Chapters []Chapter
	Logger   *zap.Logger

	chapter *Chapter
	mkvparse.DefaultHandler
}

func (h *ChapterHandler) HandleInteger(id mkvparse.ElementID, v int64, info mkvparse.ElementInfo) error {
	switch id {
	case mkvparse.ChapterFlagEnabledElement:
		h.chapter.Enabled = v == 1
	case mkvparse.ChapterTimeEndElement:
		h.chapter.TimeEnd = v
	case mkvparse.ChapterTimeStartElement:
		h.chapter.TimeStart = v
	case mkvparse.ChapterUIDElement:
		h.chapter.ID = uint64(v)
	}

	return nil
}

func (h *ChapterHandler) HandleMasterBegin(id mkvparse.ElementID, info mkvparse.ElementInfo) (bool, error) {
	h.Logger.Debug("element:begin", zap.String("element", mkvparse.NameForElementID(id)))

	switch id {
	case mkvparse.ChapterDisplayElement, mkvparse.ChaptersElement, mkvparse.EditionEntryElement:
		return true, nil
	case mkvparse.ChapterAtomElement:
		h.chapter = &Chapter{}
		return true, nil
	default:
		return false, nil
		// 	return h.DefaultHandler.HandleMasterBegin(id, info)
	}
}

func (h *ChapterHandler) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	h.Logger.Debug("element:end", zap.String("element", mkvparse.NameForElementID(id)))

	switch id {
	case mkvparse.ChapterAtomElement:
		h.Chapters = append(h.Chapters, *h.chapter)
		return nil
	default:
		return h.DefaultHandler.HandleMasterEnd(id, info)
	}
}

func (h *ChapterHandler) HandleString(id mkvparse.ElementID, v string, info mkvparse.ElementInfo) error {
	switch id {
	case mkvparse.ChapStringElement:
		h.chapter.Title = v
	}

	return nil
}

func ReadVideoChapters(r io.ReadSeeker, logger *zap.Logger) ([]Chapter, error) {
	chapterHandler := &ChapterHandler{
		Chapters: make([]Chapter, 0),
		Logger:   logger,
	}
	segmentInfoHandler := &SegmentInfoHandler{
		Logger:      logger,
		SegmentInfo: new(SegmentInfo),
	}

	h := mkvparse.NewHandlerChain(chapterHandler, segmentInfoHandler)
	err := mkvparse.ParseSections(r, h, mkvparse.InfoElement, mkvparse.ChaptersElement)
	if err != nil {
		return nil, err
	}

	chapters := make([]Chapter, len(chapterHandler.Chapters))
	for i, c := range chapterHandler.Chapters {
		chapters[i] = c.WithTimescale(segmentInfoHandler.SegmentInfo.Timescale)
	}

	return chapters, err
}

package handlers

import (
	"github.com/go-logr/logr"
	"github.com/remko/go-mkvparse"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
)

type ChapterHandler struct {
	Chapters []mediakit.Chapter
	Logger   logr.Logger

	chapter *mediakit.Chapter
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
	h.Logger.V(1).Info("element:begin", "element", mkvparse.NameForElementID(id))

	switch id {
	case mkvparse.ChapterDisplayElement, mkvparse.ChaptersElement, mkvparse.EditionEntryElement:
		return true, nil
	case mkvparse.ChapterAtomElement:
		h.chapter = &mediakit.Chapter{}
		return true, nil
	default:
		return false, nil
		// 	return h.DefaultHandler.HandleMasterBegin(id, info)
	}
}

func (h *ChapterHandler) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	h.Logger.V(1).Info("element:end", "element", mkvparse.NameForElementID(id))

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

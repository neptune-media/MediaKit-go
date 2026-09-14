package handlers

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/remko/go-mkvparse"

	"github.com/neptune-media/MediaKit-go/pkg/mediakit"
)

type SegmentInfoHandler struct {
	Logger      logr.Logger
	SegmentInfo *mediakit.SegmentInfo

	mkvparse.DefaultHandler
}

func (h *SegmentInfoHandler) HandleInteger(id mkvparse.ElementID, v int64, info mkvparse.ElementInfo) error {
	switch id {
	case mkvparse.TimecodeScaleElement:
		h.SegmentInfo.Timescale = time.Duration(v)
	}

	return nil
}

func (h *SegmentInfoHandler) HandleMasterBegin(id mkvparse.ElementID, info mkvparse.ElementInfo) (bool, error) {
	h.Logger.V(1).Info("element:begin", "element", mkvparse.NameForElementID(id))

	switch id {
	case mkvparse.InfoElement:
		h.SegmentInfo = &mediakit.SegmentInfo{}
		return true, nil
	case mkvparse.SegmentElement:
		return true, nil
	default:
		return false, nil
	}
}

func (h *SegmentInfoHandler) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	h.Logger.V(1).Info("element:end", "element", mkvparse.NameForElementID(id))

	switch id {
	case mkvparse.InfoElement:
		return nil
	default:
		return h.DefaultHandler.HandleMasterEnd(id, info)
	}
}

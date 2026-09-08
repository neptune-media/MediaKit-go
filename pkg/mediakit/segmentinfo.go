package mediakit

import (
	"time"

	"github.com/remko/go-mkvparse"
	"go.uber.org/zap"
)

type SegmentInfo struct {
	Timescale time.Duration
}

type SegmentInfoHandler struct {
	Logger      *zap.Logger
	SegmentInfo *SegmentInfo

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
	h.Logger.Debug("element:begin", zap.String("element", mkvparse.NameForElementID(id)))

	switch id {
	case mkvparse.InfoElement:
		h.SegmentInfo = &SegmentInfo{}
		return true, nil
	case mkvparse.SegmentElement:
		return true, nil
	default:
		return false, nil
		// return h.DefaultHandler.HandleMasterBegin(id, info)
	}
}

func (h *SegmentInfoHandler) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	h.Logger.Debug("element:end", zap.String("element", mkvparse.NameForElementID(id)))

	switch id {
	case mkvparse.InfoElement:
		return nil
	default:
		return h.DefaultHandler.HandleMasterEnd(id, info)
	}
}

package ffprobe

type FrameInfo struct {
	MediaType   string `json:"media_type"`
	StreamIndex int    `json:"stream_index"`
	KeyFrame    int    `json:"key_frame"`
	PictType    string `json:"pict_type"`
	PktPTS      int    `json:"pkt_pts"`
}

type Output struct {
	Frames []FrameInfo `json:"frames"`
}

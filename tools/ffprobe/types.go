package ffprobe

type FrameInfo struct {
	MediaType   string `json:"media_type"`
	StreamIndex int    `json:"stream_index"`
	KeyFrame    int    `json:"key_frame"`
	PictType    string `json:"pict_type"`
	PktPTS      int    `json:"pkt_pts"`
	PTS         int    `json:"pts"`
}

type Output struct {
	Frames  []FrameInfo `json:"frames"`
	Streams []Stream
}

type Stream struct {
	AvgFrameRate        string `json:"avg_frame_rate"`
	ChromaLocation      string `json:"chroma_location"`
	CodecLongName       string `json:"codec_long_name"`
	CodecName           string `json:"codec_name"`
	CodecTag            string `json:"codec_tag"`
	CodecTagString      string `json:"codec_tag_string"`
	CodecType           string `json:"codec_type"`
	CodedHeight         int    `json:"coded_height"`
	CodedWidth          int    `json:"coded_width"`
	ColorRange          string `json:"color_range"`
	DisplayAspectRation string `json:"display_aspect_ration"`
	ExtradataSize       int    `json:"extradata_size"`
	FieldOrder          string `json:"field_order"`
	FilmGrain           int    `json:"film_grain"`
	Index               int
	HasBFrames          int `json:"has_b_frames"`
	Height              int
	Level               int
	MaxBitRate          string `json:"max_bit_rate"`
	NbReadFrames        string `json:"nb_read_frames"`
	PixFmt              string `json:"pix_fmt"`
	Profile             string
	Refs                int
	RFrameRate          string `json:"r_frame_rate"`
	SampleAspectRatio   string `json:"sample_aspect_ratio"`
	StartPTS            int    `json:"start_pts"`
	StartTime           string `json:"start_time"`
	TimeBase            string `json:"time_base"`
	Width               int
}

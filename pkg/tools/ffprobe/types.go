package ffprobe

import "time"

// Frame represents a single frame in a video
type Frame struct {
	KeyFrame    int    `json:"key_frame" mapstructure:"key_frame" parquet:"key_frame,int(8)"`
	MediaType   string `json:"media_type" mapstructure:"media_type" parquet:"media_type,enum"`
	PictType    string `json:"pict_type" mapstructure:"pict_type" parquet:"pict_type,enum"`
	PktPTS      int    `json:"pkt_pts" mapstructure:"pkt_pts" parquet:"pkt_pts,delta"`
	PTS         int    `json:"pts" mapstructure:"pts" parquet:"pts,delta"`
	StreamIndex int    `json:"stream_index" mapstructure:"stream_index" parquet:"stream_index,int(8)"`

	Other map[string]any `mapstructure:",remain" parquet:"-"`
}

func (f Frame) Timecode() time.Duration {
	return time.Duration(f.PTS) * time.Millisecond
}

type FrameList []Frame

// Stream represents a single stream in a file
type Stream struct {
	AvgFrameRate       string `json:"avg_frame_rate" mapstructure:"avg_frame_rate" parquet:"avg_frame_rate,string"`
	ChromaLocation     string `json:"chroma_location" mapstructure:"chroma_location" parquet:"chroma_location,string"`
	CodecLongName      string `json:"codec_long_name" mapstructure:"codec_long_name" parquet:"codec_long_name,string"`
	CodecName          string `json:"codec_name" mapstructure:"codec_name" parquet:"codec_name,string"`
	CodecTag           string `json:"codec_tag" mapstructure:"codec_tag" parquet:"codec_tag,string"`
	CodecTagString     string `json:"codec_tag_string" mapstructure:"codec_tag_string" parquet:"codec_tag_string,string"`
	CodecType          string `json:"codec_type" mapstructure:"codec_type" parquet:"codec_type,string"`
	CodedHeight        int    `json:"coded_height" mapstructure:"coded_height" parquet:"coded_height,int(32)"`
	CodedWidth         int    `json:"coded_width" mapstructure:"coded_width" parquet:"coded_width,int(32)"`
	ColorRange         string `json:"color_range" mapstructure:"color_range" parquet:"color_range,string"`
	DisplayAspectRatio string `json:"display_aspect_ratio" mapstructure:"display_aspect_ratio" parquet:"display_aspect_ratio,string"`
	ExtradataSize      int    `json:"extradata_size" mapstructure:"extradata_size" parquet:"extradata_size"`
	FieldOrder         string `json:"field_order" mapstructure:"field_order" parquet:"field_order,string"`
	FilmGrain          int    `json:"film_grain" mapstructure:"film_grain" parquet:"film_grain,int(8)"`
	Index              int    `json:"index" mapstructure:"index" parquet:"index,int(8)"`
	HasBFrames         int    `json:"has_b_frames" mapstructure:"has_b_frames" parquet:"has_b_frames,int(8)"`
	Height             int    `json:"height" mapstructure:"height" parquet:"height,int(32)"`
	Level              int    `json:"level" mapstructure:"level" parquet:"level,int(32)"`
	MaxBitRate         string `json:"max_bit_rate" mapstructure:"max_bit_rate" parquet:"max_bit_rate,string"`
	NbReadFrames       string `json:"nb_read_frames" mapstructure:"nb_read_frames" parquet:"nb_read_frames,string"`
	PixFmt             string `json:"pix_fmt" mapstructure:"pix_fmt" parquet:"pix_fmt,string"`
	Profile            string `json:"profile" mapstructure:"profile" parquet:"profile,string"`
	Refs               int    `json:"refs" mapstructure:"refs" parquet:"refs,int(32)"`
	RFrameRate         string `json:"r_frame_rate" mapstructure:"r_frame_rate" parquet:"r_frame_rate,string"`
	SampleAspectRatio  string `json:"sample_aspect_ratio" parquet:"sample_aspect_ratio,string"`
	StartPTS           int    `json:"start_pts" mapstructure:"start_pts" parquet:"start_pts,delta"`
	StartTime          string `json:"start_time" mapstructure:"start_time" parquet:"start_time,string"`
	TimeBase           string `json:"time_base" mapstructure:"time_base" parquet:"time_base,string"`
	Width              int    `json:"width" mapstructure:"width" parquet:"width,int(32)"`

	Other map[string]any `mapstructure:",remain" parquet:"-"`
}

type StreamList []Stream

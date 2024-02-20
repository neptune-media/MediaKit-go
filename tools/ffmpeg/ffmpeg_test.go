package ffmpeg

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFFmpeg_GetCommandArgs(t *testing.T) {
	type fields struct {
		AudioLanguages        []string
		AudioOptions          EncodingOptions
		ContainerOptions      EncodingOptions
		DiscardAudio          bool
		DiscardSubtitles      bool
		DiscardVideo          bool
		InputArgs             []string
		InputFilename         string
		MapAllAudioStreams    bool
		MapAllSubtitleStreams bool
		MapAllVideoStreams    bool
		OutputArgs            []string
		OutputFilename        string
		SubtitleLanguages     []string
		SubtitleOptions       EncodingOptions
		UseLowerPriority      bool
		VideoOptions          EncodingOptions
		stdout                []byte
		stderr                []byte
		stdoutBuffer          bytes.Buffer
		stderrBuffer          bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "default",
			fields: fields{
				InputFilename:  "input.mkv",
				OutputFilename: "/dev/null",
			},
			want: []string{"-i", "input.mkv", "-map", "0:v:0", "-y", "/dev/null"},
		},
		{
			name: "discard everything",
			fields: fields{
				DiscardAudio:     true,
				DiscardSubtitles: true,
				DiscardVideo:     true,
				InputFilename:    "input.mkv",
				OutputFilename:   "/dev/null",
			},
			want: []string{"-i", "input.mkv", "-an", "-sn", "-vn", "-y", "/dev/null"},
		},
		{
			name: "keep english+french subtitles--english audio",
			fields: fields{
				AudioLanguages:    []string{"eng"},
				InputFilename:     "input.mkv",
				OutputFilename:    "/dev/null",
				SubtitleLanguages: []string{"eng", "fra"},
			},
			want: []string{"-i", "input.mkv", "-map", "0:v:0", "-map", "0:a:m:language:eng", "-map", "0:s:m:language:eng", "-map", "0:s:m:language:fra", "-y", "/dev/null"},
		},
		{
			name: "transcode ac3 audio @ 192k",
			fields: fields{
				AudioOptions: &GenericAudioOptions{
					Bitrate: 192,
					Codec:   "ac3",
				},
				InputFilename:  "input.mkv",
				OutputFilename: "/dev/null",
			},
			want: []string{"-i", "input.mkv", "-c:a", "ac3", "-b:a", "192k", "-map", "0:v:0", "-y", "/dev/null"},
		},
		{
			name: "with matroska reserve_index_space",
			fields: fields{
				ContainerOptions: &MkvContainerOptions{
					ReserveIndexSpace: 500,
				},
				InputFilename:  "input.mkv",
				OutputFilename: "/dev/null",
			},
			want: []string{"-i", "input.mkv", "-f", "matroska", "-reserve_index_space", "500k", "-map", "0:v:0", "-y", "/dev/null"},
		},
		{
			name: "map all streams",
			fields: fields{
				InputFilename:         "input.mkv",
				MapAllAudioStreams:    true,
				MapAllSubtitleStreams: true,
				MapAllVideoStreams:    true,
				OutputFilename:        "/dev/null",
			},
			want: []string{"-i", "input.mkv", "-map", "0:v", "-map", "0:a", "-map", "0:s", "-y", "/dev/null"},
		},
		{
			name: "map all subtitles -- only english audio",
			fields: fields{
				AudioLanguages:        []string{"eng"},
				InputFilename:         "input.mkv",
				MapAllSubtitleStreams: true,
				OutputFilename:        "/dev/null",
			},
			want: []string{"-i", "input.mkv", "-map", "0:v:0", "-map", "0:a:m:language:eng", "-map", "0:s", "-y", "/dev/null"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FFmpeg{
				AudioLanguages:        tt.fields.AudioLanguages,
				AudioOptions:          tt.fields.AudioOptions,
				ContainerOptions:      tt.fields.ContainerOptions,
				DiscardAudio:          tt.fields.DiscardAudio,
				DiscardSubtitles:      tt.fields.DiscardSubtitles,
				DiscardVideo:          tt.fields.DiscardVideo,
				InputArgs:             tt.fields.InputArgs,
				InputFilename:         tt.fields.InputFilename,
				MapAllAudioStreams:    tt.fields.MapAllAudioStreams,
				MapAllSubtitleStreams: tt.fields.MapAllSubtitleStreams,
				MapAllVideoStreams:    tt.fields.MapAllVideoStreams,
				OutputArgs:            tt.fields.OutputArgs,
				OutputFilename:        tt.fields.OutputFilename,
				SubtitleLanguages:     tt.fields.SubtitleLanguages,
				SubtitleOptions:       tt.fields.SubtitleOptions,
				UseLowerPriority:      tt.fields.UseLowerPriority,
				VideoOptions:          tt.fields.VideoOptions,
				stdout:                tt.fields.stdout,
				stderr:                tt.fields.stderr,
				stdoutBuffer:          tt.fields.stdoutBuffer,
				stderrBuffer:          tt.fields.stderrBuffer,
			}
			if got := f.GetCommandArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommandArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

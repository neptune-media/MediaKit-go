package ffmpeg

import (
	"reflect"
	"testing"
)

func TestGenericAudioOptions_GetOptions(t *testing.T) {
	type fields struct {
		Bitrate int
		Codec   string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "with Codec",
			fields: fields{
				Codec: "dummy",
			},
			want: []string{"dummy"},
		},
		{
			name: "with Bitrate",
			fields: fields{
				Bitrate: 999,
				Codec:   "dummy",
			},
			want: []string{"dummy", "-b:a", "999k"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &GenericAudioOptions{
				Bitrate: tt.fields.Bitrate,
				Codec:   tt.fields.Codec,
			}
			if got := o.GetOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

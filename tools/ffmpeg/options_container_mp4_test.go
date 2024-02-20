package ffmpeg

import (
	"reflect"
	"testing"
)

func TestMp4ContainerOptions_GetOptions(t *testing.T) {
	type fields struct {
		EnableFastStart bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "default",
			fields: fields{},
			want:   []string{"-f", "mp4"},
		},
		{
			name: "with EnableFastStart",
			fields: fields{
				EnableFastStart: true,
			},
			want: []string{"-f", "mp4", "-movflags", "+faststart"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Mp4ContainerOptions{
				EnableFastStart: tt.fields.EnableFastStart,
			}
			if got := o.GetOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMp4ContainerOptions_getMovFlags(t *testing.T) {
	type fields struct {
		EnableFastStart bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "default",
			fields: fields{},
			want:   []string{},
		},
		{
			name: "with EnableFastStart",
			fields: fields{
				EnableFastStart: true,
			},
			want: []string{"-movflags", "+faststart"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Mp4ContainerOptions{
				EnableFastStart: tt.fields.EnableFastStart,
			}
			if got := o.getMovFlags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMovFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

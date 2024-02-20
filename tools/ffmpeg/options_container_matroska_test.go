package ffmpeg

import (
	"reflect"
	"testing"
)

func TestMkvContainerOptions_GetOptions(t *testing.T) {
	type fields struct {
		ReserveIndexSpace int
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "default",
			fields: fields{},
			want:   []string{"-f", "matroska"},
		},
		{
			name: "with ReserveIndexSpace",
			fields: fields{
				ReserveIndexSpace: 100,
			},
			want: []string{"-f", "matroska", "-reserve_index_space", "100k"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &MkvContainerOptions{
				ReserveIndexSpace: tt.fields.ReserveIndexSpace,
			}
			if got := o.GetOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMkvContainerOptions_getMuxerOptions(t *testing.T) {
	type fields struct {
		ReserveIndexSpace int
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
			name: "with ReserveIndexSpace",
			fields: fields{
				ReserveIndexSpace: 44,
			},
			want: []string{"-reserve_index_space", "44k"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &MkvContainerOptions{
				ReserveIndexSpace: tt.fields.ReserveIndexSpace,
			}
			if got := o.getMuxerOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMuxerOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

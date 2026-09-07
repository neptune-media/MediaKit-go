package mediakit

import (
	"io"

	"github.com/parquet-go/parquet-go"
	"github.com/parquet-go/parquet-go/compress"
	"go.uber.org/zap"
)

type ParquetWriter[T any] struct {
	Logger *zap.Logger

	parquet.GenericWriter[T]
}

type parquetWriterOptions struct {
	Compression compress.Codec
	Logger      *zap.Logger
}

type WriterOption[T any] func(opts *parquetWriterOptions)

func WithCompression[T any](codec compress.Codec) WriterOption[T] {
	return func(opts *parquetWriterOptions) {
		opts.Compression = codec
	}
}

func WithLogger[T any](logger *zap.Logger) WriterOption[T] {
	return func(opts *parquetWriterOptions) {
		opts.Logger = logger
	}
}

func NewParquetWriter[T any](w io.Writer, opts ...WriterOption[T]) *ParquetWriter[T] {
	// Build our options from what the user gave us
	options := &parquetWriterOptions{
		Logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(options)
	}

	// Convert our options to parquet writer options
	writerOpts := make([]parquet.WriterOption, 0)
	if options.Compression != nil {
		writerOpts = append(writerOpts, parquet.Compression(options.Compression))
	}

	// Create the new writer
	return &ParquetWriter[T]{
		GenericWriter: *(parquet.NewGenericWriter[T](w, writerOpts...)),
		Logger:        options.Logger,
	}
}

func (w *ParquetWriter[T]) Close() error {
	return w.GenericWriter.Close()
}

package mediakit

import (
	"io"

	"github.com/go-logr/logr"
	"github.com/parquet-go/parquet-go"
	"github.com/parquet-go/parquet-go/compress"
)

type ParquetWriter[T any] struct {
	Logger logr.Logger

	parquet.GenericWriter[T]
}

type parquetOptions struct {
	Compression compress.Codec
	Logger      logr.Logger
	UseSchema   bool
}

type Option[T any] func(opts *parquetOptions)

func WithCompression[T any](codec compress.Codec) Option[T] {
	return func(opts *parquetOptions) {
		opts.Compression = codec
	}
}

func WithLogger[T any](logger logr.Logger) Option[T] {
	return func(opts *parquetOptions) {
		opts.Logger = logger
	}
}

func WithSchema[T any](useSchema bool) Option[T] {
	return func(opts *parquetOptions) {
		opts.UseSchema = useSchema
	}
}

func NewParquetWriter[T any](w io.Writer, opts ...Option[T]) *ParquetWriter[T] {
	// Build our options from what the user gave us
	options := &parquetOptions{
		Logger: logr.Discard(),
	}

	for _, opt := range opts {
		opt(options)
	}

	// Convert our options to parquet writer options
	writerOpts := make([]parquet.WriterOption, 0)
	if options.Compression != nil {
		writerOpts = append(writerOpts, parquet.Compression(options.Compression))
	}
	if options.UseSchema {
		writerOpts = append(writerOpts, parquet.SchemaOf(new(T)))
	}

	// Create the new writer
	return &ParquetWriter[T]{
		GenericWriter: *(parquet.NewGenericWriter[T](w, writerOpts...)),
		Logger:        options.Logger,
	}
}

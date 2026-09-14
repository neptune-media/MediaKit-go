package mediakit

import (
	"io"

	"github.com/go-logr/logr"
	"github.com/parquet-go/parquet-go"
)

type ParquetReader[T any] struct {
	Logger logr.Logger

	parquet.GenericReader[T]
}

func NewParquetReader[T any](r io.ReaderAt, opts ...Option[T]) *ParquetReader[T] {
	// Build our options from what the user gave us
	options := &parquetOptions{
		Logger: logr.Discard(),
	}

	for _, opt := range opts {
		opt(options)
	}

	// Convert our options to parquet writer options
	readerOpts := make([]parquet.ReaderOption, 0)
	if options.UseSchema {
		readerOpts = append(readerOpts, parquet.SchemaOf(new(T)))
	}

	// Create the new writer
	return &ParquetReader[T]{
		GenericReader: *(parquet.NewGenericReader[T](r, readerOpts...)),
		Logger:        options.Logger,
	}
}

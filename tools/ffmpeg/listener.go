package ffmpeg

import (
	"bufio"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// ProgressListener is a TCP listener for handling progress reports from
// ffmpeg when using the -progress flag
type ProgressListener struct {
	// ReportInterval determines how often a progress report is printed (0 for print every report)
	ReportInterval time.Duration
	listener       *net.TCPListener
}

// ProgressReport represents a single progress update when running ffmpeg
// with the -progress flag
type ProgressReport struct {
	Bitrate   string
	FPS       float64
	Frame     int
	OutTime   string
	Speed     string
	TotalSize int
}

func (p *ProgressListener) Begin() (string, error) {
	if p.listener == nil {
		addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
		if err != nil {
			return "", err
		}

		listener, err := net.ListenTCP(addr.Network(), addr)
		if err != nil {
			return "", err
		}

		p.listener = listener
	}
	return fmt.Sprintf("tcp://%s", p.listener.Addr().String()), nil
}

func (p *ProgressListener) Close() error {
	if p.listener != nil {
		err := p.listener.Close()
		p.listener = nil
		return err
	}
	return nil
}

func (p *ProgressListener) Run(logger *zap.SugaredLogger) {
	conn, err := p.listener.Accept()
	if err != nil {
		logger.Errorw("error while accepting ffmpeg connection", "err", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	report := &ProgressReport{}
	nextReport := time.Time{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Errorw("error while reading ffmpeg progress", "err", err)
			return
		}

		kv := strings.Split(line, "=")
		value := strings.TrimSpace(kv[1])
		switch kv[0] {
		case "bitrate":
			report.Bitrate = value
		case "frame":
			report.Frame, _ = strconv.Atoi(value)
		case "fps":
			report.FPS, _ = strconv.ParseFloat(value, 32)
		case "out_time":
			report.OutTime = value
		case "speed":
			report.Speed = value
		case "total_size":
			report.TotalSize, _ = strconv.Atoi(value)
		case "progress":
			end := value == "end"

			// Always print last report or if ReportInterval is 0
			if end || p.ReportInterval == 0 || time.Now().After(nextReport) {
				logger.Infow("ffmpeg progress",
					"bitrate", report.Bitrate,
					"frame", report.Frame,
					"fps", fmt.Sprintf("%.02f", report.FPS),
					"out_time", report.OutTime,
					"speed", report.Speed,
					"total_size", report.TotalSize,
				)

				// Schedule next report
				nextReport = time.Now().Add(p.ReportInterval)
			}

			if end {
				return
			}
		default:
		}
	}
}

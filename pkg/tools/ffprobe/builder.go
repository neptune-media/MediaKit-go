package ffprobe

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/pkg/tools"
)

type Builder struct {
	Filename string
	Logger   *zap.Logger
	Tool     *tools.Executable

	args        []string
	lowPriority bool
}

type PriorityCmd struct {
	*exec.Cmd
	LowPriority bool
}

func (c *PriorityCmd) Start() error {
	err := c.Cmd.Start()
	if err != nil {
		return err
	}

	if c.LowPriority {
		err = tools.ReduceProcessPriority(c.Cmd.Process)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewBuilder(tool *tools.Executable, logger *zap.Logger, filename string) *Builder {
	return &Builder{
		Filename: filename,
		Logger:   logger,
		Tool:     tool,
		args:     []string{"-of", "compact"},
	}
}

func (b *Builder) Build(ctx context.Context) *PriorityCmd {
	args := make([]string, len(b.args))
	copy(args, b.args)
	args = append(args, b.Filename)

	cmd := &PriorityCmd{
		Cmd:         exec.CommandContext(ctx, b.Tool.Name, args...),
		LowPriority: b.lowPriority,
	}

	b.Logger.Debug("built command", zap.String("path", cmd.Path), zap.Strings("args", cmd.Args), zap.String("command", strings.Join(append([]string{cmd.Path}, args...), " ")))
	return cmd
}

func (b *Builder) GetFrames() *Builder {
	b.args = append(b.args, "-show_frames")
	return b
}

func (b *Builder) GetFramesCount() *Builder {
	b.args = append(b.args, "-show_streams", "-count_frames")
	return b
}

func (b *Builder) UseLowPriority() *Builder {
	b.lowPriority = true
	return b
}

func (b *Builder) UseThreads(num int) *Builder {
	b.args = append(b.args, "-threads", fmt.Sprintf("%d", num))
	return b
}

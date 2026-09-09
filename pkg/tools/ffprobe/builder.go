package ffprobe

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"go.uber.org/zap"

	"github.com/neptune-media/MediaKit-go/pkg/tools"
)

type CommandBuilder struct {
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

func NewCommandBuilder(tool *tools.Executable, logger *zap.Logger, filename string) *CommandBuilder {
	return &CommandBuilder{
		Filename: filename,
		Logger:   logger,
		Tool:     tool,
		args:     []string{"-of", "compact"},
	}
}

func (b *CommandBuilder) Build(ctx context.Context) *PriorityCmd {
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

func (b *CommandBuilder) GetFrames() *CommandBuilder {
	b.args = append(b.args, "-show_frames")
	return b
}

func (b *CommandBuilder) GetFramesCount() *CommandBuilder {
	b.args = append(b.args, "-show_streams", "-count_frames")
	return b
}

func (b *CommandBuilder) UseLowPriority() *CommandBuilder {
	b.lowPriority = true
	return b
}

func (b *CommandBuilder) UseThreads(num int) *CommandBuilder {
	b.args = append(b.args, "-threads", fmt.Sprintf("%d", num))
	return b
}

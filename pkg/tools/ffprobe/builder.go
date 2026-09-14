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
	Logger *zap.Logger
	Tool   *tools.Executable

	args        []string
	lowPriority bool
}

func NewCommandBuilder(tool *tools.Executable, logger *zap.Logger) *CommandBuilder {
	return &CommandBuilder{
		Logger: logger,
		Tool:   tool,
		args:   []string{"-of", "compact"},
	}
}

func (b *CommandBuilder) Build(ctx context.Context, filename string) *tools.PriorityCmd {
	args := make([]string, len(b.args))
	copy(args, b.args)
	args = append(args, filename)

	cmd := &tools.PriorityCmd{
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

func (b *CommandBuilder) LowPriority() *CommandBuilder {
	b.lowPriority = true
	return b
}

func (b *CommandBuilder) Threads(num int) *CommandBuilder {
	b.args = append(b.args, "-threads", fmt.Sprintf("%d", num))
	return b
}

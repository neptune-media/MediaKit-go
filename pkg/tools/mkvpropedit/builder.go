package mkvpropedit

import (
	"context"
	"os/exec"
	"strings"

	"github.com/go-logr/logr"

	"github.com/neptune-media/MediaKit-go/pkg/tools"
)

type CommandBuilder struct {
	Logger logr.Logger
	Tool   *tools.Executable

	args        []string
	lowPriority bool
}

func NewCommandBuilder(tool *tools.Executable, logger logr.Logger) *CommandBuilder {
	return &CommandBuilder{
		Logger: logger,
		Tool:   tool,
		args:   make([]string, 0),
	}
}

func (b *CommandBuilder) Build(ctx context.Context, inputFilename string) *tools.PriorityCmd {
	args := make([]string, len(b.args)+1)
	args[0] = inputFilename
	copy(args[1:], b.args)

	cmd := &tools.PriorityCmd{
		Cmd:         exec.CommandContext(ctx, b.Tool.Name, args...),
		Logger:      b.Logger,
		LowPriority: b.lowPriority,
	}

	b.Logger.V(1).Info(
		"built command",
		"path", cmd.Path,
		"args", cmd.Args,
		"command", strings.Join(append([]string{cmd.Path}, args...), " "),
	)
	return cmd
}

func (b *CommandBuilder) Copy() *CommandBuilder {
	nb := &CommandBuilder{}
	*nb = *b
	nb.args = make([]string, len(b.args))
	copy(nb.args, b.args)
	return nb
}

func (b *CommandBuilder) LowPriority() *CommandBuilder {
	nb := b.Copy()
	nb.lowPriority = true
	return nb
}

func (b *CommandBuilder) RenameChapters(filename string) *CommandBuilder {
	nb := b.Copy()
	nb.args = append(nb.args, "-c", filename)
	return nb
}

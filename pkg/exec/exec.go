package exec

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type Executor interface {
	Exec([]string) Executor
	BashExec(string) Executor
	Cmd() *exec.Cmd
	Err() error
	StdoutLines() []string
	StderrLines() []string
	Clear() Executor
}

type executor struct {
	stdout []byte
	stderr []byte
	err    error
	cmd    *exec.Cmd
}

func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) logOutput(outType string) {
	var lines []string

	if outType == "stderr" {
		lines = e.StderrLines()
	} else {
		lines = e.StdoutLines()
	}

	if log.Trace().Enabled() {
		for _, line := range lines {
			log.Trace().Str(outType, line).Msg("")
		}
	}
}

func (e *executor) Exec(commands []string) Executor {
	log.Trace().Strs("commands", commands).Msg("Starting execution")

	if len(commands) == 0 {
		e.err = fmt.Errorf("Cannot execute empty command")
		return e
	}

	args := []string{}
	var out []byte
	for i := 1; i < len(commands); i++ {
		args = append(args, commands[i])
	}

	e.cmd = exec.Command(commands[0], args...)
	out, e.err = e.cmd.CombinedOutput()
	if e.err != nil {
		e.stderr = out
		e.logOutput("stderr")
	} else {
		e.stdout = out
		e.logOutput("stdout")
	}

	return e
}

func (e *executor) BashExec(command string) Executor {
	return e.Exec([]string{"bash", "-c", command})
}

func (e *executor) Clear() Executor {
	e.stdout = nil
	e.stderr = nil
	e.err = nil
	e.cmd = nil
	return e
}

func (e *executor) Cmd() *exec.Cmd {
	return e.cmd
}

func (e *executor) Err() error {
	return e.err
}

func (e *executor) StdoutLines() []string {
	return strings.Split(string(e.stdout), "\n")
}

func (e *executor) StderrLines() []string {
	return strings.Split(string(e.stderr), "\n")
}

func ExecBashCmd(commands []string) error {
	log.Trace().Strs("commands", commands).Msg("Executing bash commands")
	cmd := exec.Command("bash", "-c", strings.Join(commands, "&&"))

	o, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msg("Failed execution")
		return fmt.Errorf("could not start execution: %v", err)
	}

	if log.Trace().Enabled() {
		log.Trace().Msg("Output")
		for _, line := range strings.Split(string(o), "\n") {
			log.Trace().Msg(line)
		}
	}

	return nil
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: input | %s COMMAND [ ARG ] ...\n", os.Args[0])
		os.Exit(1)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if err := run(os.Stdin, out, os.Args[1], os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(in io.Reader, out io.Writer, command string, args []string) error {
	proc := exec.Command(command, args...)
	procStdin, _ := proc.StdinPipe()
	procStdout, _ := proc.StdoutPipe()

	if err := proc.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	streamR, streamW := io.Pipe()
	go func() {
		_, _ = io.Copy(procStdin, io.TeeReader(in, streamW))
		procStdin.Close()
		streamW.Close()
	}()

	for l, r := bufio.NewScanner(streamR), bufio.NewScanner(procStdout); l.Scan() && r.Scan(); {
		fmt.Fprintf(out, "%s\t%s\n", l.Text(), r.Text())
	}

	return nil
}

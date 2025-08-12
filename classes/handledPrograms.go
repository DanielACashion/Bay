package classes

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type HandledProgram struct {
	Name    string
	Version int //maybe make this a triple int (x.y.z)
	Running bool
	Error   bool
	In      io.WriteCloser
	OUT     io.ReadCloser
}

// returns object and if it was created successfully
func HandledProgramCreate(program_name string) (this HandledProgram, wasok bool) {
	wasok = false
	var version int
	var err error
	var name string
	namestrings := strings.Split(program_name, "_")
	if len(namestrings) > 1 {
		name = namestrings[0]
		version, err = strconv.Atoi(strings.Split(namestrings[1], ".")[0])
		if err != nil {
			return
		}
	} else {
		return
	}
	this.Name = name
	this.Version = version
	this.Error = false
	this.Running = false
	wasok = true
	return
}

func (h *HandledProgram) Load() {
	cmd := exec.Command("plugins/" + fmt.Sprintf("%s_%d.exe", h.Name, h.Version))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("OOPS Failed TO Load STDIN: %s\n", h.Name)
		return
	}
	h.In = stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("FAILED TO SET THE STDOUT: %s\n", h.Name)
	}
	h.OUT = stdout

	if err := cmd.Start(); err != nil {
		fmt.Printf("OOPS Failed TO Start: %s\n", h.Name)
		return
	}
	h.Running = true
}

func (h *HandledProgram) Unload() {
	h.In.Close()
}

func (h *HandledProgram) Update() {
	if !h.Running {
		return
	}
	bytes, err := h.In.Write([]byte(fmt.Sprintf("update sent to %s_%d\n", h.Name, h.Version)))
	if err != nil {
		fmt.Println("OOPS")
	}
	_ = bytes
}
func (h *HandledProgram) FileName() string {
	return fmt.Sprintf("%s_%d.exe", h.Name, h.Version)
}

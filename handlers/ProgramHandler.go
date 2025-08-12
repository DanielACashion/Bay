package handlers

import (
	"AutoReloader/classes"
	"AutoReloader/config"
	"AutoReloader/folders"
	"fmt"
	"os"
	"sort"
	"time"
)

type ProgramHandler struct {
	_handled_programs  []classes.HandledProgram
	_program_contains  map[string]int //[name]index of _handled_programs
	_cancelation_token bool
	_config            config.ConfigValues
}

// reads name of all current plugins in folder and set length of in mem list
func ReadProgramsToLoad() (programs []classes.HandledProgram) {
	entries, err := os.ReadDir("plugins")
	if err != nil {
		os.Exit(1)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
		program, ok := classes.HandledProgramCreate(entry.Name())
		if !ok {
			continue
		}
		programs = append(programs, program)
	}
	return
}

func (handle *ProgramHandler) LoadPrograms(programs []classes.HandledProgram) {
	var previousCount int = len(handle._handled_programs)
	handle._handled_programs = append(handle._handled_programs, programs...)
	var currentCount int = len(handle._handled_programs)
	var i int = 0
	for i = 0; i < currentCount-previousCount; i += 1 {
		fmt.Printf("LOADING: %s\n", handle._handled_programs[previousCount+i].Name)
		handle._handled_programs[previousCount+i].Load()
	}
}

func (handle *ProgramHandler) UnloadAllPrograms() {
	handle.UnloadPrograms(handle._handled_programs)
}

func (handle *ProgramHandler) UnloadPrograms(programs []classes.HandledProgram) {
	sort.Slice(programs, func(i, j int) bool {
		return handle._program_contains[programs[i].Name] > handle._program_contains[programs[j].Name]
	})

	for _, program := range programs {
		if handle._handled_programs[handle._program_contains[program.Name]].Name == program.Name {
			if !handle._config.AllowRollbacks && handle._handled_programs[handle._program_contains[program.Name]].Version > program.Version {
				fmt.Println("UH OH you tried to downgrade and have that option turned off")
				continue
			}

			handle._handled_programs[handle._program_contains[program.Name]].Unload()
			fmt.Printf("Unloading finished for %s\n", program.Name)
			//this needs to be a calculated value
			moveBundleToRemoved(handle._handled_programs[handle._program_contains[program.Name]].FileName())
			handle._handled_programs = append(handle._handled_programs[:handle._program_contains[program.Name]], handle._handled_programs[handle._program_contains[program.Name]+1:]...)
		}

	}
}

// used as an intermediate sendUpdates func, main sendUpdates func
func (handle *ProgramHandler) SendUpdates() {
	for {
		if handle._cancelation_token {
			break
		}
		handle.CheckUpdateFolder()
		for _, program := range handle._handled_programs {
			program.Update()
		}
		time.Sleep(3 * time.Second)
	}
	handle._cancelation_token = true
}

//func (handle *ProgramHandler) SendEvent(event EventMessage){}

func (handle *ProgramHandler) GetResponses() {
	var bytes []byte
	for {

		if handle._cancelation_token {
			break
		}
		for _, program := range handle._handled_programs {
			if !program.Running {
				continue
			}
			bytes = make([]byte, 1024)
			numRead, err := program.OUT.Read(bytes)
			if err != nil {
				continue
			}
			if numRead == 0 {
				continue
			}
			fmt.Printf("Printing: %s\n", string(bytes[:numRead]))
		}
	}
}

func (handle *ProgramHandler) CheckUpdateFolder() {
	// Check if the update folder exists
	entries, err := os.ReadDir("updates")
	if err != nil {
		os.Exit(1)
	}
	if len(entries) < 1 {
		return
	}

	//get all current files in the updates folder
	var programsToUpdate []classes.HandledProgram
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		program, ok := classes.HandledProgramCreate(entry.Name())
		if !ok {
			continue
		}
		var ind int = handle._program_contains[program.Name]
		if handle._handled_programs[ind].Name != program.Name {
			continue
		}
		if handle._handled_programs[ind].Version >= program.Version {
			continue
		}
		programsToUpdate = append(programsToUpdate, program)
	}
	if len(programsToUpdate) < 1 {
		return
	}
	handle.UpdatePrograms(programsToUpdate)
}

// launches a file watcher that will notify this program if one of these programs has an update
// later stage add, watch the update folder for hot reloading
func watchForUpdates() {
	//if somefile updates, update that interinternalPrograms

}

func moveBundleToRemoved(fileName string) {
	var wasOK bool = folders.MoveFileTo(fileName, "plugins/", "removed/")
	if !wasOK {
		fmt.Printf("TriedMovingFile: Failed\n")
		return
	}
	fmt.Printf("")
}

func moveBundleToCurrent(fileName string) {
	fmt.Printf("UPDATED: %s\n", fileName)
	var wasOK bool = folders.MoveFileTo(fileName, "updates/", "plugins/")
	if !wasOK {
		fmt.Printf("TriedMovingFile: Failed\n")
		return
	}
	fmt.Printf("")
}

// this will need to be an async func that allows this program to unregister all info
// then moves the new version into the folder, and re registers
func (handle *ProgramHandler) UpdatePrograms(programsToUpdate []classes.HandledProgram) {
	//unload
	handle.UnloadPrograms(programsToUpdate)

	//move files
	for _, program := range programsToUpdate {
		moveBundleToCurrent(program.FileName())
	}

	//load
	handle.LoadPrograms(programsToUpdate)
}

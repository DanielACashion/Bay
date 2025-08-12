package main

import (
	"AutoReloader/config"
	"AutoReloader/handlers"
)

//folder struct
//root
//.config.txt
//--/updates/
//--/plugins/
//--/removed/
//--/logs/{date}.txt

var _config config.ConfigValues
var _handle handlers.ProgramHandler

func main() {
	//check if all folders exist if not create them
	_config.PrePass()
	_config.Load()

	programsToLoad := handlers.ReadProgramsToLoad()

	_handle.LoadPrograms(programsToLoad)
	defer _handle.UnloadAllPrograms()
	//go watchForUpdates()
	go _handle.GetResponses()
	_handle.SendUpdates()
}

package config

import (
	"fmt"
	"os"
)

type ConfigValues struct {
	AllowRollbacks bool
}

const configFileName string = "config.txt"

func (con *ConfigValues) Load() {
	//full file read to string

	config, err := os.ReadFile(configFileName)
	if err != nil {
		fmt.Println("ERROR occured attempting to open the config file")
		//set defaults
	}
	fmt.Println(string(config)) //lex and token the config
}

func (con *ConfigValues) PrePass() {
	setUpFolders()
	_ = _default_config
}

func setUpFolders() {
	if _, err := os.ReadDir("updates"); err != nil {
		fmt.Println("Adding Updates Folder")
		err := os.Mkdir("updates", 0755)
		if err != nil {
			fmt.Println("UPDATES FAILED")
		}
	}

	if _, err := os.ReadDir("plugins"); err != nil {
		fmt.Println("Adding Plugins Folder")
		err := os.Mkdir("plugins", 0755)
		if err != nil {
			fmt.Println("Plugins FAILED")
		}
	}

	if _, err := os.ReadDir("removed"); err != nil {
		fmt.Println("Adding Removed Folder")
		err := os.Mkdir("removed", 0755)
		if err != nil {
			fmt.Println("Removed FAILED")
		}
	}

	if _, err := os.ReadDir("logs"); err != nil {
		fmt.Println("Adding Logs Folder")
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("logs FAILED")
		}
	}

	if _, err := os.ReadFile(configFileName); err != nil {
		if file, err := os.Create(configFileName); err == nil {
			defer file.Close()
			file.Write([]byte(_default_config))
		}

	}
}

// config tokens
const (
	AutoReload             = "AutoReload"
	AllowRollbacks         = "AllowRollbacks"
	_default_config string = `AutoReload=false
AllowRollbacks=false
`
)

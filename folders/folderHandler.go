package folders

import (
	"fmt"
	"os"
)

func MoveFileTo(fileName string, locationFrom string, locationTo string) bool {
	//some folder from validation

	//some folder to validation

	//does file exist valiation

	//if all is ok
	err := os.Rename(locationFrom+fileName, locationTo+fileName)
	if err != nil {
		fmt.Printf("Error Moving %s, From: %s, To: %s\n", fileName, locationFrom, locationTo)
		return false
	}

	fmt.Printf("MOVED: %s, From: %s, To: %s\n", fileName, locationFrom, locationTo)
	return true
}

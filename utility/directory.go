package utility

import "os"

// CreateDirectory => Exported
// create directory if not existed
func CreateDirectory(path string) {
	var err error

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0666)
	}

	if (err != nil) && (!os.IsNotExist(err)) {
		Log("(Utility=>CreateDirectory): Error in creating directory", err)
	}
}

// GetHomeDir => Exported
func GetHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

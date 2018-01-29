package utility

import "io/ioutil"

// ReadFile => Exported
// read file content
func ReadFile(fullPath string) []byte {
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		Log("(Utility=>ReadFile): Error in read file", err)
		content = nil
	}

	if len(content) == 0 {
		content = nil
	}

	return content
}

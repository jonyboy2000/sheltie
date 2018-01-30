package utility

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Log => Exported
func Log(msg string, e error) {
	currentTime := time.Now().Local().Format("2006-01-02")

	f, err := os.OpenFile("./log/"+currentTime, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("(Utility=>Log): Error in creating log file")

		defer f.Close()
		return
	}

	defer f.Close()

	log.SetOutput(f)
	log.Println("= " + msg + "(" + e.Error() + ")")
}

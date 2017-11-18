package stealer

import (
	"fmt"

	"github.com/Gujarats/logger"
)

func printProgress(fileName string) {
	fileName = logger.GetColorFormat(logger.Yellow, logger.Faint, fileName)
	message := logger.GetColorFormat(logger.Cyan, logger.Faint, "Converting file = ")

	fmt.Println(message, fileName)
}

func printFinish() {
	message := logger.GetColorFormat(logger.Red, logger.Faint, "All files succesfully converted !!!")
	fmt.Println(message)
}

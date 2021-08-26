package logger

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)
const (
	InfoColor    = "\033[1;34m"
	FatalColor  = "\033[31;31m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	ColorWhite = "\033[37m"
)
func getKey() bool {
	fmt.Println("Press Y(y) if you want to print log in a console")
	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
	if char == 'Y' || char == 'y' {
		return true
	}
	return false
}
func Info(message string, params interface{}) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0765)
	checkFile(file, err)
	defer file.Close()
	logger := log.New(file, "Information:\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message, params)
	output := getKey()
	if output == true {
		fmt.Println("Log info")
		logger = log.New(os.Stdout, "Information:\t", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Println(InfoColor, message, ColorWhite, params)
	}
}

func Error(message string, params interface{}) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0765)
	checkFile(file, err)
	defer file.Close()
	logger := log.New(file, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message, params)
	output := getKey()
	if output == true {
		fmt.Println("Log info")
		logger = log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Println(ErrorColor, message, ColorWhite, params)
	}
}

func Warning(message string, params interface{}) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0765)
	checkFile(file, err)
	defer file.Close()
	logger := log.New(file, "Warning:\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message, params)
	output := getKey()
	if output == true {
		fmt.Println("Log info")
		logger = log.New(os.Stdout, "Warning:\t", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Println(WarningColor, message, ColorWhite, params)
	}
}

func Fatal(message string, params interface{}) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0765)
	checkFile(file, err)
	defer file.Close()
	logger := log.New(file, "Fatal:\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message, params)
	output := getKey()
	if output == true {
		fmt.Println("Log info")
		logger = log.New(os.Stdout, "Fatal: \t", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Panic(FatalColor, message, ColorWhite, params)
	}
}

func checkFile(file *os.File, err error) {
	if err != nil {
		fmt.Println("Not possible to open ", file, err)
	}
}

func ClientError(w http.ResponseWriter, statusCode int) {
	logger := log.New(os.Stdout, "Client Error\t", log.Ldate|log.Ltime)
	logger.Println("Client error occurred with status ", statusCode)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func ServerError(w http.ResponseWriter, err error) {
	logger := log.New(os.Stdout, "Internal Server Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
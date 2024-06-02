package logger

import (
	"log"
	"os"
)

var (
    WarnLogger *log.Logger
    InfoLogger *log.Logger
    ErrorLogger *log.Logger
    TraceLogger *log.Logger
)

func init() {
        InfoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
        WarnLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
        ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
        TraceLogger = log.New(os.Stderr, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

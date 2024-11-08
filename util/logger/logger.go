package logger

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"sync"
)

var (
	logger     *logrus.Logger
	loggerInit sync.Once
)

func NewLogger() *logrus.Logger {
	loggerInit.Do(func() {
		//check directory is existing
		logFilePath := "/app/logs"
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			// Directory does not exist, create it
			err := os.MkdirAll(logFilePath, 0755) // 0755 sets permissions (read/write for owner, read-only for others)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
			fmt.Println("Directory created:", logFilePath)
		} else {
			fmt.Println("Directory already exists:", logFilePath)
		}

		// set file log
		fileName := "logData.log"
		logFile, err := os.OpenFile(path.Join(logFilePath, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}

		//
		logger = logrus.New()
		//logger.SetFormatter(&logrus.TextFormatter{
		//	FullTimestamp:   true,
		//	TimestampFormat: customTimeFormat(),
		//})
		logger.SetOutput(logFile)
	})

	return logger

	//log.SetOutput(logFile)
	//log.Println("Ini adalah entri log.")
	//
	//dir, err2 := os.Getwd()
	//if err2 != nil {
	//	log.Println("os.Getwd err:", err2)
	//}
	//log.Println("current folder :", dir)
}

func GetLogger() *logrus.Logger {
	return logger
}

func ContextLogger(ctx context.Context) *logrus.Entry {
	mlogger := GetLogger()

	var fields logrus.Fields

	if ctxRqId, ok := ctx.Value(fiber.HeaderXRequestID).(string); ok {
		fields = logrus.Fields{
			"requestId": ctxRqId,
		}
	}

	return mlogger.WithFields(fields)
}

//func customTimeFormat() string {
//	// Set the desired time zone (e.g., New York)
//	loc, err := time.LoadLocation("Asia/Jakarta")
//	if err != nil {
//		fmt.Println("Error loading timezone:", err)
//		return time.Now().Format(time.RFC3339) // Default to UTC if error occurs
//	}
//	return time.Now().In(loc).Format("2006-01-02T15:04:05 MST") // Custom time format with time zone
//}

package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// // output ไปยัง terminal และ file
	// file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to open log file")
	// }
	// defer file.Close()

	// multi := zerolog.MultiLevelWriter(os.Stdout, file)
	// log.Logger = log.Output(multi)
}

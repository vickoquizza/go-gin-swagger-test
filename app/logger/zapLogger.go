package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

var JSONConfig = []byte(`{
	"level": "debug",
	"encoding": "json",
	"outputPaths": ["stdout", "my.logs"],       
   "errorOutputPaths": ["stderr", "my.logs"], 
    "encoderConfig": {
     "messageKey": "message",
     "levelKey": "level",
     "levelEncoder": "lowercase"
    }
}`)

type ZapLogger struct {
	Adapter *zap.Logger
	config  zap.Config
}

func (z *ZapLogger) BuildLogger() {
	if err := json.Unmarshal(JSONConfig, &z.config); err != nil {
		panic(err)
	}

	logger, err := z.config.Build()

	if err != nil {
		panic(err)
	}

	z.Adapter = logger
}

func NewZapLogger() *ZapLogger {
	return &ZapLogger{}
}

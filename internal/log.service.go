package internal

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const diretorioLogs = "logs"

const ChaveContextoUnicoKey = "ChaveContextoUnicoKey"

var (
	logger *zap.Logger
)

type Config struct {
	LogFile    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	DevMode    bool
}

func ConfiguracaoPadrao() Config {
	return Config{
		LogFile:    diretorioLogs + "/api-erp.log",
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     15,
		Compress:   true,
		DevMode:    false,
	}
}

func Inicializar(cfg Config) {
	rotator := &lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(rotator), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)
}

// Start loga o início de uma operação
func Start(ctx context.Context, message string) {
	if logger == nil {
		Inicializar(ConfiguracaoPadrao())
	}
	logger.Info(montarMensagemComContexto(ctx, message) + "...")
}

// Info loga uma mensagem informativa ou decisão tomada
func Info(ctx context.Context, message string) {
	if logger == nil {
		Inicializar(ConfiguracaoPadrao())
	}
	logger.Info(montarMensagemComContexto(ctx, message))
}

// Success loga uma operação bem sucedida
func Success(ctx context.Context, message string) {
	if logger == nil {
		Inicializar(ConfiguracaoPadrao())
	}
	logger.Info(montarMensagemComContexto(ctx, message) + "!")
}

// Error loga uma mensagem de erro
func Error(ctx context.Context, message string) {
	if logger == nil {
		Inicializar(ConfiguracaoPadrao())
	}
	logger.Error(montarMensagemComContexto(ctx, message) + "!!")
}

// Exception loga uma mensagem de erro com um erro associado
func Exception(ctx context.Context, message string, err *error) {
	if logger == nil {
		Inicializar(ConfiguracaoPadrao())
	}

	if err == nil {
		logger.Error(montarMensagemComContexto(ctx, message) + "!!")
		return
	}
	logger.Error(montarMensagemComContexto(ctx, message)+"!!",
		zap.Error(*err),
	)
}

func montarMensagemComContexto(ctx context.Context, message string) string {
	if idContexto := obterIdContextoUnico(ctx); idContexto != "" {
		return "[" + idContexto + "] " + message
	}
	return message
}

func obterIdContextoUnico(ctx context.Context) string {
	if v := ctx.Value(ChaveContextoUnicoKey); v != nil {
		return v.(string)
	}
	return ""
}

// Sync descarrega quaisquer entradas de log em buffer
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

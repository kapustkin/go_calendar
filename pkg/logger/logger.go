package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Init initializes the standard logger
func Init(app string, version string) {
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "2006-01-02 15:04:05"
	plainFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}

	log.SetFormatter(plainFormatter)
	log.SetReportCaller(true)
	log.WithFields(
		log.Fields{
			"app":     app,
			"version": version,
		})
	//return requestLogger
}

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := fmt.Sprint(entry.Time.Format(f.TimestampFormat))
	return []byte(fmt.Sprintf(
		"%s %s %v:L%v %s\n",
		f.LevelDesc[entry.Level],
		timestamp,
		entry.Caller.Function,
		entry.Caller.Line,
		entry.Message)), nil
}

// NewChiLogger Chi logger
func NewChiLogger() func(next http.Handler) http.Handler {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}
	return middleware.RequestLogger(&StructuredLogger{log})
}

// StructuredLogger Logger struct
type StructuredLogger struct {
	Logger *log.Logger
}

// StructuredLoggerEntry entry
type StructuredLoggerEntry struct {
	Logger log.FieldLogger
}

// NewLogEntry create log record
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: log.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

// Write event
func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

// Panic event
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

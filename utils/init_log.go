package utils

import (
	"io"
	"os"

	nested_formatter "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	// logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

func InitLog(prgname string, isdebug bool) *logrus.Entry {
	logrus.SetFormatter(&nested_formatter.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05", //TimestampFormat: time.RFC3339,
		FieldsOrder:     []string{"model", prgname},
	})
	//logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true, FullTimestamp: true})
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	if isdebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetReportCaller(false) // 设置在输出日志中添加文件名和方法信息

	if _, err := os.Stat("log"); err != nil {
		os.Mkdir("log", 0755)
	}
	logf, err := rotatelogs.New(
		"log/"+prgname+".%Y%m%d",
		//rotatelogs.WithLinkName(prgname),
		//rotatelogs.WithMaxAge(-1),
		//rotatelogs.WithRotationCount(10),
		rotatelogs.WithRotationCount(7),           // max log files
		rotatelogs.WithRotationSize(10*1024*1024), // 10M per log file
	)
	if err != nil {
		logrus.Errorf("failed to create rotatelogs: %s", err)
		return nil
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, logf))
	//logrus.SetOutput(os.Stdout)

	// add syslog output with hook
	// hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, prgname)
	// if err != nil {
	// 	logrus.Error("Unable to connect to local syslog daemon, ", err)
	// } else {
	// 	logrus.AddHook(hook)
	// }

	programlog := logrus.WithFields(logrus.Fields{
		//"version": "3.0",
		//"model": prgname,
	})
	return programlog
}

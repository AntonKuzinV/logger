package logger

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

const (
	Ldate         = 1 << iota
	Ltime
	Lmicroseconds
	Llongfile
	Lshortfile
	Lmsgprefix
)

const (
	InfoColor string = "\x1b[32m"
	WarningColor string = "\x1b[33m"
	ErrorColor string = "\x1b[35m"
	FatalColor string = "\x1b[31;1m"
)

type Logger struct {
	mutex sync.Mutex
	out io.Writer
	logLevel string
	flag int
	buffer []byte
}

func New(out io.Writer, flag int) *Logger {
	return &Logger{out: out, flag: flag}
}

func itoa(buffer *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buffer = append(*buffer, b[bp:]...)
}

func (l *Logger) format(buffer *[]byte, time time.Time, file string, line int) {
	if l.flag&Lmsgprefix == 0 {
		*buffer = append(*buffer, l.logLevel...)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := time.Date()
			itoa(buffer, year, 4)
			*buffer = append(*buffer, '/')
			itoa(buffer, int(month), 2)
			*buffer = append(*buffer, '/')
			itoa(buffer, day, 2)
			*buffer = append(*buffer, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := time.Clock()
			itoa(buffer, hour, 2)
			*buffer = append(*buffer, ':')
			itoa(buffer, min, 2)
			*buffer = append(*buffer, ':')
			itoa(buffer, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buffer = append(*buffer, '.')
				itoa(buffer, time.Nanosecond()/1e3, 6)
			}
			*buffer = append(*buffer, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buffer = append(*buffer, file...)
		*buffer = append(*buffer, ':')
		itoa(buffer, line, -1)
		*buffer = append(*buffer, ": "...)
	}
	if l.flag&Lmsgprefix != 0 {
		*buffer = append(*buffer, l.logLevel...)
	}
}

func (l *Logger) Output(level string, call int, s string) error {
	l.logLevel = level
	now := time.Now()
	var file string
	var line int
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		l.mutex.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(call)
		if !ok {
			file = "???"
			line = 0
		}
		l.mutex.Lock()
	}
	l.buffer = l.buffer[:0]
	l.format(&l.buffer, now, file, line)
	l.buffer = append(l.buffer, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buffer = append(l.buffer, '\n')
	}
	_, err := l.out.Write(l.buffer)
	return err
}

func (l *Logger) Info(msg ...interface{}) {
	fmt.Print(InfoColor)
	l.Output("INFO ", 2, fmt.Sprint(msg...))
}

func (l *Logger) Infof(format string, msg ...interface{}) {
	l.Output("INFO ",2, fmt.Sprintf(format, msg...))
}

func (l *Logger) Error(msg ...interface{}) {
	fmt.Print(ErrorColor)
	l.Output("ERROR ", 2, fmt.Sprint(msg...))
}

func (l *Logger) Errorf(format string, msg ...interface{}) {
	l.Output("ERROR ", 2, fmt.Sprintf(format, msg...))
}

func (l *Logger) Warning(msg ...interface{}) {
	fmt.Print(WarningColor)
	l.Output("WARNING ", 2, fmt.Sprint(msg...))
}

func (l *Logger) Warningf(format string, msg ...interface{}) {
	l.Output("WARNING ",2, fmt.Sprintf(format, msg...))
}

func (l *Logger) Fatal(msg ...interface{}) {
	fmt.Print(FatalColor)
	l.Output("FATAL ", 2, fmt.Sprint(msg...))
}

func (l *Logger) Fatalf(format string, msg ...interface{}) {
	l.Output("FATAL ", 2, fmt.Sprintf(format, msg...))
}
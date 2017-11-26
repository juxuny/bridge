package log

import (
	"log"
	"os"
)

var (
	w Writer
	logger *log.Logger
)

type Writer struct {
	f *os.File
	FileName string
}

func (t Writer) Write(data []byte) (r int, e error) {
	if t.FileName == "" {
		return os.Stdout.Write(data)
	}
	if t.f != nil {
		return t.f.Write(data)
	}
	t.f, e = os.OpenFile(t.FileName,  os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if e != nil {
		panic(e)
		return
	}
	return t.f.Write(data)
}

func (t Writer) Close() (e error) {
	defer func() {
		t.f = nil
	}()
	return t.f.Close()
}

func Init(fileName string) {
	w = Writer{FileName: fileName}
	logger = log.New(w, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime)
}

func Printf(f string, v ...interface{}) {
	//logger.Print(time.Now().Format("[2006-01-02 15:04:06] "))
	logger.Printf(f + "\n", v...)
}

func Panic(v interface{}) {
	//logger.Print(time.Now().Format("[2006-01-02 15:04:06] "))
	panic(v)
}

func Print(v ...interface{}) {
	//logger.Print(time.Now().Format("[2006-01-02 15:04:06] "))
	logger.Println(v...)
}
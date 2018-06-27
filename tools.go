package tools

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//生成随机字符串
func RandString() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func Time2String(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999")
}

//将"2016-02-15 12:00:00"或者"2016-04-18 09:33:56.694"等格式转化为time.Time
func StringToTime(s string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	return t, err
}

//将"2016-04-22T21:47:49.694123232+08:00"或者"2016-04-22T21:47:49+08:00"等格式转化为time.Time
func StringToTime1(s string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02T15:04:05+08:00", s, loc)
	return t, err
}

func Any(value interface{}) string {
	return FormatAtom(reflect.ValueOf(value))
}

func FormatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 5, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value"

	}
}

// 将interface{}类型转为string类型
func Interface2String(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", errors.New("invalid interface type")
	}
}

//打印浮点数的特定表现格式
func Float64Bits(f float64, d int) {
	b := *(*uint64)(unsafe.Pointer(&f))

	switch d {
	case 16:
		fmt.Printf("浮点数%.1f的16进制表示是%#016x\n", f, b)
	case 2:
		fmt.Printf("浮点数%.1f的2进制表示是%#02b\n", f, b)
	default:
		fmt.Println("error decimal: ", d)
	}
}

//测量一段代码执行时间
func TraceCode() func() {
	start := time.Now()
	return func() {
		t := time.Now().Sub(start).Nanoseconds()
		fmt.Printf("运行耗时:%d(纳秒)\n", t)
	}

}

//打印当前堆栈
func PrintStack(all bool) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, all)

	log.Println("[FATAL] catch a panic,stack is: ", string(buf[:n]))
}

func GetStack(all bool) string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, all)
	return string(buf[:n])
}

// zero-copy, []byte转为string类型
// 注意，这种做法下，一旦[]byte变化，string也会变化
// 谨慎，黑科技！！除非性能瓶颈，否则请使用string(b)1

func Bytes2String(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
	// pb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	// ps := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// ps.Data = pb.Data
	// ps.Len = pb.Len
	// return
}

// zero-coy, string类型转为[]byte
// 注意，这种做法下，一旦string变化，程序立马崩溃且不能recover
// 谨慎，黑科技！！除非性能瓶颈，否则请使用[]byte(s)
func String2Bytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
	// pb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	// ps := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// pb.Data = ps.Data
	// pb.Len = ps.Len
	// pb.Cap = ps.Len
	// return
}

// 判断一个error是否是io.EOF
func IsEOF(err error) bool {
	if err == nil {
		return false
	} else if err == io.EOF {
		return true
	} else if oerr, ok := err.(*net.OpError); ok {
		if oerr.Err.Error() == "use of closed network connection" {
			return true
		}
	} else {
		if err.Error() == "use of closed network connection" {
			return true
		}
	}
	return true
}

// 获取本机ip
func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if !strings.Contains(ipnet.IP.String(), "192.168") {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}

// 获取常用runtime统计信息
func RuntimeStats(gc bool, heapObj bool, goroutineNum bool) []int64 {
	s := &runtime.MemStats{}
	runtime.ReadMemStats(s)

	stats := make([]int64, 5)
	if gc {
		// 上一次gc耗时
		t := s.PauseNs[(s.NumGC+255)%256]
		stats[0] = int64(t)

		// gc总次数
		num := s.NumGC
		stats[1] = int64(num)

		// 下一次gc触发时，heapalloc的大小
		ng := s.NextGC
		stats[2] = int64(ng)
	}

	if heapObj {
		ho := s.HeapObjects
		stats[3] = int64(ho)
	}

	if goroutineNum {
		ng := runtime.NumGoroutine()
		stats[4] = int64(ng)
	}

	return stats
}

// []byte转为10进制整数
var errBase10 = errors.New("failed to convert to Base10")

func ByteToBase10(b []byte) (n uint64, err error) {
	base := uint64(10)

	n = 0
	for i := 0; i < len(b); i++ {
		var v byte
		d := b[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		default:
			n = 0
			err = errors.New("failed to convert to Base10")
			return
		}
		n *= base
		n += uint64(v)
	}

	return n, err
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func FileExist(fn string) bool {
	_, err := os.Stat(fn)
	return err == nil || os.IsExist(err)
}

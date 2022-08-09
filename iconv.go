//
// iconv.go
//
package iconv

// #cgo linux CFLAGS: -I "/mnt/d/workspace/gowork/src/github.com/892294101/iconv/dep/include"
// #cgo linux LDFLAGS: -L "/mnt/d/workspace/gowork/src/github.com/892294101/iconv/dep/lib" -liconv
// #include <iconv.h>
// #include <stdlib.h>
// #include <errno.h>
//
// size_t bridge_iconv(iconv_t cd,
//		       char *inbuf, size_t *inbytesleft,
//                     char *outbuf, size_t *outbytesleft) {
//   return iconv(cd, &inbuf, inbytesleft, &outbuf, outbytesleft);
// }
import "C"

import (
	"bytes"
	"errors"
	"io"
	"syscall"
	"unsafe"
)

var EILSEQ = syscall.Errno(C.EILSEQ)
var E2BIG = syscall.Errno(C.E2BIG)

const DefaultBufSize = 4096

type Iconv struct {
	Handle C.iconv_t
	power  bool
}

// Open returns a conversion descriptor cd, cd contains a conversion state and can not be used in multiple threads simultaneously.
func Open(tocode string, fromcode string) (cd *Iconv, err error) {

	tocode1 := C.CString(tocode)
	defer C.free(unsafe.Pointer(tocode1))

	fromcode1 := C.CString(fromcode)
	defer C.free(unsafe.Pointer(fromcode1))

	ret, err := C.iconv_open(tocode1, fromcode1)
	if err != nil {
		return
	}
	cd = &Iconv{Handle: ret}
	return
}

func (cd *Iconv) Close() error {
	var err error
	if !cd.power {
		cd.power = true
		_, err = C.iconv_close(cd.Handle)
	} else {
		err = errors.New("iconv has been closed")
	}
	return err
}

func (cd *Iconv) Conv(b []byte, outbuf []byte) (out []byte, inleft int, err error) {

	outn, inleft, err := cd.Do(b, len(b), outbuf)
	if err == nil || err != E2BIG {
		out = outbuf[:outn]
		return
	}

	w := bytes.NewBuffer(nil)
	w.Write(outbuf[:outn])

	inleft, err = cd.DoWrite(w, b[len(b)-inleft:], inleft, outbuf)
	out = w.Bytes()
	return
}

func (cd *Iconv) ConvString(s string) (string, error) {
	dat := []byte(s)
	outBuf := make([]byte, len(dat))
	s1, _, err := cd.Conv(dat, outBuf)
	return string(s1), err
}

func (cd *Iconv) Do(inbuf []byte, in int, outbuf []byte) (out, inleft int, err error) {

	if in == 0 {
		return
	}

	inbytes := C.size_t(in)
	inptr := &inbuf[0]

	outbytes := C.size_t(len(outbuf))
	outptr := &outbuf[0]
	_, err = C.bridge_iconv(cd.Handle,
		(*C.char)(unsafe.Pointer(inptr)), &inbytes,
		(*C.char)(unsafe.Pointer(outptr)), &outbytes)

	out = len(outbuf) - int(outbytes)
	inleft = int(inbytes)
	return
}

func (cd *Iconv) DoWrite(w io.Writer, inbuf []byte, in int, outbuf []byte) (inleft int, err error) {

	if in == 0 {
		return
	}

	inbytes := C.size_t(in)
	inptr := &inbuf[0]

	for inbytes > 0 {
		outbytes := C.size_t(len(outbuf))
		outptr := &outbuf[0]
		_, err = C.bridge_iconv(cd.Handle,
			(*C.char)(unsafe.Pointer(inptr)), &inbytes,
			(*C.char)(unsafe.Pointer(outptr)), &outbytes)
		w.Write(outbuf[:len(outbuf)-int(outbytes)])
		if err != nil && err != E2BIG {
			return int(inbytes), err
		}
	}

	return 0, nil
}

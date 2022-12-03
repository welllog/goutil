package reqx

import (
	"errors"
	"io"
)

type Reader struct {
	s []byte
	i int64
}

func NewReader(b []byte) *Reader {
	return &Reader{
		s: b,
		i: 0,
	}
}

func (r *Reader) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

func (r *Reader) Size() int64 { return int64(len(r.s)) }

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.i:])
	r.i += int64(n)
	return
}

func (r *Reader) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("reqx.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

func (r *Reader) ReadByte() (byte, error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	b := r.s[r.i]
	r.i++
	return b, nil
}

func (r *Reader) UnreadByte() error {
	if r.i <= 0 {
		return errors.New("reqx.Reader.UnreadByte: at beginning of slice")
	}
	r.i--
	return nil
}

func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	b := r.s[r.i:]
	m, err := w.Write(b)
	if m > len(b) {
		panic("reqx.Reader.WriteTo: invalid Write count")
	}
	r.i += int64(m)
	n = int64(m)
	if m != len(b) && err == nil {
		err = io.ErrShortWrite
	}
	return
}

func (r *Reader) Close() error {
	return nil
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("reqx.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("reqx.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

func (r *Reader) Reset(b []byte) {
	*r = Reader{b, 0}
}

func (r *Reader) Bytes() []byte {
	return r.s[r.i:]
}

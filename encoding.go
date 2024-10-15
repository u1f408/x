// import "github.com/u1f408/x"
package x

import (
	"fmt"
	"hash/fnv"

	"encoding/base32"
	"encoding/base64"
	"github.com/jxskiss/base62"
)

var (
	EncBase32 = NewBase32Encoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567", true)
	EncBase32Hex = NewBase32Encoding("0123456789ABCDEFGHIJKLMNOPQRSTUV", true)
	EncBase32Lower = NewBase32Encoding("0123456789abcdefghijklmnopqrstuv", false)

	EncBase62 = NewBase62Encoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	EncBase62Modified = NewBase62Encoding("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	EncBase64 = NewBase64Encoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", true)
	EncBase64Url = NewBase64Encoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", true)
	EncBase64UrlRaw = NewBase64Encoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", false)
)

func ShortHashf(format string, a ...any) []byte {
	i := fnv.New32a()
	fmt.Fprintf(i, format, a...)
	return i.Sum(nil)
}

func Hashf(format string, a ...any) []byte {
	i := fnv.New128a()
	fmt.Fprintf(i, format, a...)
	return i.Sum(nil)
}

type Encoding interface {
	DecodeString(in string) ([]byte, error)
	EncodeString(in []byte) string
}

type Base32Encoding struct {
	i *base32.Encoding
}

func NewBase32Encoding(alphabet string, padding bool) *Base32Encoding {
	pad := base32.NoPadding
	if padding {
		pad = base32.StdPadding
	}

	i := base32.NewEncoding(alphabet).WithPadding(pad)
	return &Base32Encoding{ i: i }
}

func (e *Base32Encoding) DecodeString(in string) ([]byte, error) {
	dst := make([]byte, e.i.DecodedLen(len(in)))
	n, err := e.i.Decode(dst, []byte(in))
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

func (e *Base32Encoding) EncodeString(in []byte) string {
	dst := make([]byte, e.i.EncodedLen(len(in)))
	e.i.Encode(dst, in)
	return string(dst)
}

type Base62Encoding struct {
	i *base62.Encoding
}

func NewBase62Encoding(alphabet string) *Base62Encoding {
	i := base62.NewEncoding(alphabet)
	return &Base62Encoding{ i: i }
}

func (e *Base62Encoding) DecodeString(in string) ([]byte, error) {
	return e.i.DecodeString(in)
}

func (e *Base62Encoding) EncodeString(in []byte) string {
	return e.i.EncodeToString(in)
}

type Base64Encoding struct {
	i *base64.Encoding
}

func NewBase64Encoding(alphabet string, padding bool) *Base64Encoding {
	pad := base64.NoPadding
	if padding {
		pad = base64.StdPadding
	}

	i := base64.NewEncoding(alphabet).WithPadding(pad)
	return &Base64Encoding{ i: i }
}

func (e *Base64Encoding) DecodeString(in string) ([]byte, error) {
	dst := make([]byte, e.i.DecodedLen(len(in)))
	n, err := e.i.Decode(dst, []byte(in))
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

func (e *Base64Encoding) EncodeString(in []byte) string {
	dst := make([]byte, e.i.EncodedLen(len(in)))
	e.i.Encode(dst, in)
	return string(dst)
}

//go:build windows && amd64

package gtbox_encryption

/*
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/../libs/gtgo/windows
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/../libs/gtgo/windows/libgtgo.a
*/
import "C"

//go:build windows && amd64

package ltc_scrypt

/*
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/../libs/gtgo/windows
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/../libs/gtgo/windows/libgtgo.a
*/
import "C"

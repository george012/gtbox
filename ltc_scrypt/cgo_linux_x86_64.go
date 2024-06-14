//go:build linux && amd64

package ltc_scrypt

/*
#cgo linux,amd64 CFLAGS: -I${SRCDIR}/../libs/gtgo/linux
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/../libs/gtgo/linux/libgtgo.a
*/
import "C"

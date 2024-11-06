//go:build linux && amd64

package gtbox_encryption

/*
#cgo linux,amd64 CFLAGS: -I${SRCDIR}/../libs/gtgo/linux
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/../libs/gtgo/linux/libgtgo.a
*/
import "C"

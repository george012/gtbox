//go:build linux && arm64

package gtbox_encryption

/*
#cgo linux,arm64 CFLAGS: -I${SRCDIR}/../libs/gtgo/linux
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/../libs/gtgo/linux/libgtgo_arm64.a
*/
import "C"

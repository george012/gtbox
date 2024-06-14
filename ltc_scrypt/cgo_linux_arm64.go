//go:build linux && arm64

package ltc_scrypt

/*
#cgo linux,arm64 CFLAGS: -I${SRCDIR}/../libs/gtgo/linux
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/../libs/gtgo/linux/libgtgo_arm64.a
*/
import "C"

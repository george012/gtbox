//go:build darwin && arm64

package gtbox_encryption

/*
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/../libs/gtgo/darwin
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/../libs/gtgo/darwin/libgtgo_arm64.a
*/
import "C"

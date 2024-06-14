//go:build darwin && amd64

package ltc_scrypt

/*
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/../libs/gtgo/darwin
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/../libs/gtgo/darwin/libgtgo.a
*/
import "C"

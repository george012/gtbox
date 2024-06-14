//go:build darwin && amd64

package gtbox_encryption

/*
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/../libs/gtgo/darwin
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/../libs/gtgo/darwin/libgtgo.a
*/
import "C"

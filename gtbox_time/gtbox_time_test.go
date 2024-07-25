package gtbox_time

import (
	"fmt"
	"testing"
)

func TestTimeSimples(t *testing.T) {
	aNow := NowUTC()
	fmt.Printf("%v\n", aNow)
	fmt.Printf("11 digits: %d\n", aNow.Unix())
	fmt.Printf("13 digits: %d\n", aNow.UnixMilli())
	fmt.Printf("16 digits: %d\n", aNow.UnixMicro())
	fmt.Printf("19 digits: %d\n", aNow.UnixNano())
}

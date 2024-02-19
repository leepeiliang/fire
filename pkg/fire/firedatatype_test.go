package fire

import (
	"fmt"
	"testing"
)

func TestDataStrip(t *testing.T) {

	out := StringStrip("011回路002地址联动请求")
	fmt.Println(out)
}

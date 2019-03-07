// date: 2019-03-07
package main

import (
	"fmt"
	"github.com/Jarvens/Exchange-Agent/common"
)

func main() {
	s := []string{"a", "b", "c", "d"}
	fmt.Print(common.SliceRemove(s, "a"))

}

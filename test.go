// date: 2019-03-09
package main

import (
	"fmt"
	"github.com/gofrs/uuid"
)

func main() {

	uid, _ := uuid.NewV1()
	uid1, _ := uuid.NewV4()
	fmt.Println(uid.String())
	fmt.Println(uid1.String())
}

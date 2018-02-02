// Test conversion from custom types
package main

import (
	"fmt"

	domain "bitbucket.org/dkfbasel/dev.grpc-tags/proto_test"
)

func main() {

	tmp := domain.TestRequest{}
	tmp.Name = "myname"
	tmp.Function = 20

	tmpembed := domain.TestNested{}
	tmpembed.Name = "myembeddedname"
	tmpembed.Function = "mycustomfunction"
	tmp.TestNested = tmpembed

	nestedItem := domain.TestNested{}
	nestedItem.Name = "mynesteditem"
	nestedItem.Function = "mycustomnesteditemfunction"

	tmp.SubItemNested = &nestedItem

	nestedItemSlice1 := domain.TestNested{}
	nestedItemSlice1.Name = "slice1name"
	nestedItemSlice1.Function = "slice1function"

	nestedItemSlice2 := domain.TestNested{}
	nestedItemSlice2.Name = "slice2name"
	nestedItemSlice2.Function = "slice2function"

	tmp.SubItemsNested = []*domain.TestNested{&nestedItemSlice1, &nestedItemSlice2}

	tmp2, err := tmp.Convert()
	if err != nil {
		fmt.Printf("error: %+v\n\n˙", err)
	}
	fmt.Println("-- proto from custom --")
	fmt.Printf("%+v\n", tmp2)

	tmp2.GetFunction()

	fmt.Println("-- custom from proto --")
	tmp3, err := domain.ConvertTestRequest(tmp2)
	if err != nil {
		fmt.Printf("error: %+v\n\n˙", err)
	}
	fmt.Printf("%+v\n", tmp3)

}

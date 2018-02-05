// Test conversion from custom types
package main

import (
	"fmt"

	// without explicit naming the package name would be grpcservide
	// (as defined in the protobuf definition)
	domain "github.com/dkfbasel/protobuf/example/domain"
)

func main() {

	fmt.Println("- Create a new item to be ironed")

	item := domain.Item{}
	item.Name = "Shirt"

	fmt.Println("- Make a wrinkled item out of it")

	wrinkled := domain.WrinkledItem{Item: item}
	wrinkled.Customer = "A custom conscious guy"
	wrinkled.Wrinkels = 23

	fmt.Println("- Convert the wrinkled item for protobuf transmission")
	wrinkledProto, err := wrinkled.Proto()
	if err != nil {
		fmt.Println("-- error: could not convert to protobuf struct: ", err)
	}

	fmt.Printf("-- protobuf: %v\n", wrinkledProto)

	fmt.Println("- Convert back to our custom struct")
	newWrinkled := domain.WrinkledItem{}
	err = newWrinkled.FromProto(wrinkledProto)
	if err != nil {
		fmt.Println("-- error: could not convert to protobuf struct: ", err)
	}

	fmt.Println("- Iron out the wrinkles (one was forgotten)")
	unwrinkled := domain.SmoothItem{}
	unwrinkled.Item = newWrinkled.Item
	unwrinkled.Wrinkels = 1
	unwrinkled.Cost = 50

	fmt.Println("- Convert the unwrinkled shirt for protobuf transmission")
	unwrinkledProto, err := unwrinkled.Proto()
	if err != nil {
		fmt.Println("-- error: could not convert to protobuf struct: ", err)
	}

	fmt.Printf("-- protobuf: %v\n", unwrinkledProto)

}

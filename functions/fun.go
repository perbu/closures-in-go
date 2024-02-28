package main

import "fmt"

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	fun := makeFun()
	fmt.Println("fun:", fun())
	funFun := makeFunFun()
	actualFun := funFun()
	fmt.Println("actual fun:", actualFun)

	funFun2 := makeFunFun2()
	actualFun2 := funFun2()
	fmt.Println("actual fun2:", actualFun2())
	return nil
}

func makeFun() func() int {
	return func() int {
		return 1
	}
}

func makeFunFun() func() int {
	return makeFun()
}

// makeFunFun2 is a function that returns
// a function that returns a function
// that returns an int
func makeFunFun2() func() func() int {
	return func() func() int {
		return func() int {
			return 1
		}
	}
}

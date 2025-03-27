package main

import (
	"fmt"
)

const monthsinYear = 12

const (
	dayofChristmas = 12
	firstdayofWeek = "Monday"
	sunisHot       = true
)

func main() {
	var student1 string        //type is a string
	var student2 = "Spiderman" // type is inferred

	var y int = 6
	x := 2 //type is inferred

	var a string = "hello"
	var b int  // defualt value is 0
	var c bool // default value is false

	fmt.Println("Hello World!")
	fmt.Println(student2)
	fmt.Println(y)
	fmt.Println(x)
	// default values of different types
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

	student1 = "Keanu"
	fmt.Println(student1)

	var e, f, g, h int = 1, 3, 5, 7

	fmt.Println(e, f, g, h)

	var (
		j int    = 1
		k string = "hello"
	)

	fmt.Println(j, k)
	fmt.Println(a)
	/*
	   How const syntax is structured

	   const constname type = value

	*/

	const PI = 3.14

	fmt.Println(PI)
	fmt.Println(monthsinYear)
	fmt.Println(PI * monthsinYear)
	fmt.Println(sunisHot, dayofChristmas, firstdayofWeek)

	fmt.Print(a)
	fmt.Print(b, "\n")
	// "\n" for new line, otherwise it starts on same line as it finished
	fmt.Print(a, "\n")
	fmt.Print(b, "\n")
	// you can call in same function with commas
	fmt.Print(a, " ", b, "\n")
	// print() will put a space in between 2 ints or if niether are string
	fmt.Print(f, g, "\n")

	// println() adds whitespace and a new line is added at the end

	// printf() = %v is the value of the variable and the %T is the type of the variable
	fmt.Printf("a has value: %v and type: %T\n", a, a)
	fmt.Printf("e has value: %v and type: %T\n", e, e)
	fmt.Printf("sunisHot has value: %v and type: %T\n", sunisHot, sunisHot)

	/*
	   %v prints the value in default format
	   %#v prints the value in go syntax
	   %T tpye of value
	   %% prints the % sign
	*/

	var v = 15.5
	var txt = "Testing"

	fmt.Printf("%v\n", v)
	fmt.Printf("%#v\n", v)
	fmt.Printf("%v%%\n", v)
	fmt.Printf("%T\n", v)

	fmt.Printf("%v\n", txt)
	fmt.Printf("%#v\n", txt)
	fmt.Printf("%T\n", txt)

	/*
	   %b base 2
	   %d base 10
	   %+d base10 and always show sign
	   %o base 8
	   %O base 8 with leding 0o
	   %x base 16 lowercase
	   %X base 16 uppercase
	   %#x base 16 with leading 0x
	   %4d pad with spaces (width 4, right)
	   %-4d pad with spaces (width 4, left)
	   %04d pad with zeores (wifth 4)
	*/

	var u = 15

	fmt.Printf("%b\n", u)
	fmt.Printf("%d\n", u)
	fmt.Printf("%+d\n", u)
	fmt.Printf("%o\n", u)
	fmt.Printf("%O\n", u)
	fmt.Printf("%x\n", u)
	fmt.Printf("%X\n", u)
	fmt.Printf("%#x\n", u)
	fmt.Printf("%4d\n", u)
	fmt.Printf("%-4d\n", u)
	fmt.Printf("%04d\n", u)

	/*
	   %s plain string
	   %q double quoted string
	   %8s value as plain string (width 8 right)
	   %-8s value as plain string (width 8 left)
	   %x hex dump of byte values
	   % x hex dump with spaces
	*/

	var test = "hello"

	fmt.Printf("%s\n", test)
	fmt.Printf("%q\n", test)
	fmt.Printf("%8s\n", test)
	fmt.Printf("%-8s\n", test)
	fmt.Printf("%x\n", test)
	fmt.Printf("% x\n", test)

	// %t value of boolean operator

	var w = true
	var q = false

	fmt.Printf("%t\n", w)
	fmt.Printf("%t\n", q)

	/*
	   %e Scientific notation with 'e' as exponent
	   %f Decimal point, no exponent
	   %.2f Default width, precision 2
	   %6.2f Width 6, precision 2
	   %g Exponent as needed, only necessary digits
	*/

	var i = 3.141

	fmt.Printf("%e\n", i)
	fmt.Printf("%f\n", i)
	fmt.Printf("%.2f\n", i)
	fmt.Printf("%6.2f\n", i)
	fmt.Printf("%g\n", i)

	var aa bool = true
	var ab int = 5
	var ac float32 = 3.14
	var ad string = "string"

	fmt.Println("Boolean: ", aa)
	fmt.Println("Integer: ", ab)
	fmt.Println("Float: ", ac)
	fmt.Println("String: ", ad)

	// declare an array with var
	var array_test = [5]int{6, 7, 8, 9, 10}
	fmt.Println(array_test)

	var arr1 = []int{1, 2, 3, 4, 5}
	fmt.Println(arr1)

	// interchangeable with string and int

	var prices = []int{10, 20, 30}
	// accessing the array
	fmt.Println(prices[0])

	// changing the array
	prices[2] = 60
	fmt.Println(prices)

	// if length is specified and values havent been assigned fully then it will go to the default value of that type
	var arr2 = [5]int{}              //not initialized
	var arr3 = [5]int{1, 2}          //partially initialized
	var arr4 = [5]int{1, 2, 3, 4, 5} //fully initialized

	fmt.Println(arr2)
	fmt.Println(arr3)
	fmt.Println(arr4)

	// only call the second and 3rd int in the 5 length array
	var arr5 = [5]int{1: 10, 2: 40}

	fmt.Println(arr5)

	// len() function is used to find the length
	var arr6 = []string{"Interstellar", "Cars", "Superman", "Alien vs Predator"}

	fmt.Println(len(arr1))
	fmt.Println(arr6)
	fmt.Println(len(arr6))

	// go slice is similar to an array. slices can grow and shirnk as you see fit.

	thickslice := []int{1, 2, 3}

	// 2 functions to return length and capacity of a slice: len() and cap()

	fmt.Println(len(thickslice))
	fmt.Println(cap(thickslice))
	fmt.Println(thickslice)

	var thickarray = []int{10, 11, 12, 13, 14, 15}
	thickerslice := thickarray[1:4]

	fmt.Printf("slice = %v\n", thickerslice)

	// when making, if you dont set capacity, it will default to length
	myslice1 := make([]int, 5, 10)
	fmt.Printf("myslice1 = %v\n", myslice1)
	fmt.Printf("length = %d\n", len(myslice1))
	fmt.Printf("capacity = %d\n", cap(myslice1))

	// access slice via (slice[location])
	fmt.Println(thickslice[0])

	// change slice values

	thickslice[0] = 30
	fmt.Println(thickslice)

	// append elemets to a slice. You need to add the ... in order to unpack the slice.
	myslice2 := []int{1, 2, 3}
	myslice3 := []int{4, 5, 6}
	myslice := []int{700, 800, 900, 7}
	myslice4 := append(myslice2, myslice3...)
	fmt.Println(myslice4)
	myslice4 = append(myslice4, myslice...)
	fmt.Println(myslice4)

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	fmt.Println(numbers)
	neededNumbers := numbers[:len(numbers)-10] // slice syntax, but start is omitted and end is the full len(numbers) - 10)
	numbersCopy := make([]int, len(neededNumbers))
	// copy(dst, src) where src is where teh numbers are coming from.
	copy(numbersCopy, neededNumbers)

	fmt.Println(numbersCopy)

	// operator
	var ba = 15 + 25
	fmt.Println(ba)

	var bb = ba + 40
	fmt.Println(bb)

	var bc = bb + 60
	fmt.Println(bc)

	// +, -, *, /, %, ++, --
	var bd = bc - 10
	fmt.Println(bd)

	var be = bd * 2
	fmt.Println(be)

	var bf = be / 2
	fmt.Println(bf)
	// prints or holds the remainder value not the quotient
	var bg = bf % 60
	fmt.Println(bg)

	ba++
	fmt.Println(ba)

	ba--
	fmt.Println(ba)

	bg += 10
	fmt.Println(bg)

	var r = 2
	r ^= 3
	fmt.Println(r)
	r *= 3
	fmt.Println(r)

	// == equal to, != not equal, > greater than, < less than, >= greater than or equal to, <= less than or equal to
	// && logical and will return boolean true if both statements are true, || logical or returns true if one of the statements is true, ! logical not reverse the result returns false if the result is true

	// Go Conditions - a condition can either be true or false, < less than, > greater than, <= >=, == equal to, != not equal to.

	var ta = 5
	var tb = 4
	var truthvalue = ta > tb
	fmt.Println(truthvalue)
	truthvalue = ta < tb
	fmt.Println(truthvalue)
	truthvalue = ta <= tb
	fmt.Println(truthvalue)
	truthvalue = ta >= tb
	fmt.Println(truthvalue)
	truthvalue = ta == tb
	fmt.Println(truthvalue)
	truthvalue = ta != tb
	fmt.Println(truthvalue)

	if ta < tb {
		fmt.Println("a is larger than b")
	} else {
		fmt.Println("b is larger than a")
	}

	if ta == tb {
		fmt.Println("samesies")
	} else if ta < tb {
		fmt.Println("smol boi")
	} else {
		fmt.Println("big boi")
	}

	if ta >= 2 {
		fmt.Println("big boi")
		if ta > tb {
			fmt.Println("larger boi")
		}
	} else {
		fmt.Println("smol")
	}

	day := 6

	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6:
		fmt.Println("Saturday")
	case 7:
		fmt.Println("Sunday")
	default:
		fmt.Println("Not a day in the week")
	}

	switch day {
	case 1, 3, 5:
		fmt.Println("Odd Day")
	case 2, 4:
		fmt.Println("Even Day")
	case 6, 7:
		fmt.Println("Weekend!")
	}
}

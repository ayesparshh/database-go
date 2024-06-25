package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
	Pincode int `json:"pincode"`
}

type User struct {
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}

func main() {

	dir := "./"

	db ,err := New(dir,nil)
	if err != nil {
		fmt.Println("err",err)
	}

	employees := []User{
		{"SpaceX",json.Number("42"),"SpaceX","SpaceX",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
		{"Tesla",json.Number("42"),"Tesla","Tesla",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
		{"Google",json.Number("42"),"Google","Google",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
		
	}

}


// git remote add origin https://github.com/ayesparshh/database-go.git
// git branch -M main
// git push -u origin main
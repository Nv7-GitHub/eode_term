package main

func Combo(args []interface{}) {
	Write("combo", "%v", args)
}

func init() {
	Cmd("combo", Combo, STRING, VARIADIC)
}

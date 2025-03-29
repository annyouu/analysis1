package a

func f() {
	var a interface{}

	_ = a.(int) // want "unsafe type assertion"
	_, _ = a.(int) // want "unsafe type assertion"

	switch a := a.(type) {
	case int:
		println(a) // ここは問題なし
	}
}
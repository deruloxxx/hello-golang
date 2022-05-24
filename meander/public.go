package meander

type Facade interface {
	// 外部向けのビューを返すメソッド
	Public() interface{}
}

func Public(o interface{}) interface{} {
	// Q. o.(Facade)について調べる
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}

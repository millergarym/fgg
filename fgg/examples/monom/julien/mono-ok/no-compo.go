package main;


type Any(type ) interface {};

type Pair(type ) struct { fst Any(); snd Any()};

type tI(type a Any()) interface {
	m(type )() Any()
};

type A(type ) struct {};

func (x A(type )) m(type )() Any() {
	return A(){}
};

type B(type ) struct {};
func (x B(type )) m(type )() Any() {
	return B(){}
};

type C(type ) struct {};
func (x C(type )) mtop(type )(y tI()) Any() {
	return y.m()()
};


func main() { _ =
	Pair(){ A(){}, B(){}}
}
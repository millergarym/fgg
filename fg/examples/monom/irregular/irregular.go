package main;
type Nat interface { Add(n Nat) Nat };
type Zero struct {};
type Succ struct { pred Nat };
func (m Zero) Add(n Nat) Nat { return n };
func (m Succ) Add(n Nat) Nat { return Succ{m.pred.Add(n)} };
type FuncPairNatPairNat interface { Apply(x PairNat) PairNat };
type FuncNatNat interface { Apply(x Nat) Nat };
type FuncPairPairNatPairPairNat interface { Apply(x PairPairNat) PairPairNat };
type FuncPairPairPairNatPairPairPairNat interface { Apply(x PairPairPairNat) PairPairPairNat };
type incr struct { n Nat };
func (this incr) Apply(x Nat) Nat { return x.Add(this.n) };
type PairPairPairNat struct { fst PairPairNat; snd PairPairNat };
type PairPairNat struct { fst PairNat; snd PairNat };
type PairNat struct { fst Nat; snd Nat };
type pairMapPairPairNat struct { f FuncPairPairNatPairPairNat };
type pairMapNat struct { f FuncNatNat };
type pairMapPairNat struct { f FuncPairNatPairNat };
func (p pairMapNat) Apply(x PairNat) PairNat { return PairNat{p.f.Apply(x.fst), p.f.Apply(x.snd)} };
func (p pairMapPairNat) Apply(x PairPairNat) PairPairNat { return PairPairNat{p.f.Apply(x.fst), p.f.Apply(x.snd)} };
func (p pairMapPairPairNat) Apply(x PairPairPairNat) PairPairPairNat { return PairPairPairNat{p.f.Apply(x.fst), p.f.Apply(x.snd)} };
type BalancedNat interface { BalancedMap(f FuncNatNat) BalancedNat };
type BalancedPairPairNat interface { BalancedMap(f FuncPairPairNatPairPairNat) BalancedPairPairNat };
type BalancedPairPairPairNat interface { BalancedMap(f FuncPairPairPairNatPairPairPairNat) BalancedPairPairPairNat };
type BalancedPairNat interface { BalancedMap(f FuncPairNatPairNat) BalancedPairNat };
type LeafPairPairPairNat struct {};
type NodeNat struct { label Nat; children BalancedPairNat };
type NodePairPairNat struct { label PairPairNat; children BalancedPairPairPairNat };
type NodePairNat struct { label PairNat; children BalancedPairPairNat };
func (leaf LeafPairPairPairNat) BalancedMap(f FuncPairPairPairNatPairPairPairNat) BalancedPairPairPairNat { return LeafPairPairPairNat{} };
func (node NodePairPairNat) BalancedMap(f FuncPairPairNatPairPairNat) BalancedPairPairNat { return NodePairPairNat{f.Apply(node.label), node.children.BalancedMap(pairMapPairPairNat{f})} };
func (node NodePairNat) BalancedMap(f FuncPairNatPairNat) BalancedPairNat { return NodePairNat{f.Apply(node.label), node.children.BalancedMap(pairMapPairNat{f})} };
func (node NodeNat) BalancedMap(f FuncNatNat) BalancedNat { return NodeNat{f.Apply(node.label), node.children.BalancedMap(pairMapNat{f})} };
type D struct {};
func (d D) _0() Nat { return Zero{} };
func (d D) _1() Nat { return Succ{D{}._0()} };
func (d D) _2() Nat { return Succ{D{}._1()} };
func (d D) _3() Nat { return Succ{D{}._2()} };
func (d D) _4() Nat { return Succ{D{}._3()} };
func (d D) _5() Nat { return Succ{D{}._4()} };
func (d D) _6() Nat { return Succ{D{}._5()} };
func (d D) _7() Nat { return Succ{D{}._6()} };
func main() { _ = NodeNat{D{}._1(), NodePairNat{PairNat{D{}._2(), D{}._3()}, NodePairPairNat{PairPairNat{PairNat{D{}._4(), D{}._5()}, PairNat{D{}._6(), D{}._7()}}, LeafPairPairPairNat{}}}}.BalancedMap(incr{D{}._1()}) }
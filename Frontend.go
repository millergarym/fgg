package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
	"github.com/rhu1/fgg/fgr"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

const (
	EVAL_TO_VAL = -1 // Must be < 0
	NO_EVAL     = -2 // Must be < EVAL_TO_VAL
)

/* !!! WIP -- refactoring whole main package !!! */

/* -- Interp */

type Interp interface {
	GetProgram() base.Program
	SetProgram(p base.Program)
	Eval(steps int) base.Program
	vPrint(x string)
	vPrintln(x string)
}

func parse(verbose bool, a base.Adaptor, src string, strict bool) base.Program {
	vPrintln(verbose, "\nParsing AST:")
	prog := a.Parse(strict, src) // AST (Program root)
	vPrintln(verbose, prog.String())

	vPrintln(verbose, "\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	return prog
}

// N.B. currently FG panic comes out implicitly as an underlying run-time panic
// CHECKME: add explicit FG panics?
// If steps == EVAL_TO_VAL, then eval to value
func eval(intrp Interp, steps int) base.Program {
	if steps < NO_EVAL {
		panic("Invalid number of steps: " + strconv.Itoa(steps))
	}
	p_init := intrp.GetProgram()
	allowStupid := true
	intrp.vPrintln("\nEntering Eval loop:")
	intrp.vPrintln("Decls:")
	for _, v := range p_init.GetDecls() {
		intrp.vPrintln("\t" + v.String() + ";")
	}
	intrp.vPrintln("Eval steps:")
	intrp.vPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p_init.GetMain())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for `||!done` below)
		p_init.GetMain().IsValue() // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	p := intrp.GetProgram() // Convenient for re-assign to p inside loop
	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		intrp.SetProgram(p)
		intrp.vPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetMain()))
		intrp.vPrintln("Checking OK:") // TODO: maybe disable by default, enable by flag
		// TODO FIXME: check actual type preservation of e_main (not just typeability)
		p.Ok(allowStupid)
		if !done && p.GetMain().IsValue() {
			done = true
		}
	}
	p_res := intrp.GetProgram()
	intrp.vPrintln(p_res.GetMain().ToGoString()) // Final result
	return p_res
}

/* FG */

type FGInterp struct {
	verboseHelper
	prog fg.FGProgram
}

var _ Interp = &FGInterp{}

func NewFGInterp(verbose bool, src string, strict bool) *FGInterp {
	var a fg.FGAdaptor
	prog := parse(verbose, &a, src, strict).(fg.FGProgram)
	return &FGInterp{verboseHelper{verbose}, prog}
}

func (intrp *FGInterp) GetProgram() base.Program {
	return intrp.prog
}

func (intrp *FGInterp) SetProgram(p base.Program) {
	intrp.prog = p.(fg.FGProgram)
}

func (intrp *FGInterp) Eval(steps int) base.Program {
	return eval(intrp, steps)
}

/* FGR */

type FGRInterp struct {
	verboseHelper
	prog fgr.FGRProgram
}

var _ Interp = &FGRInterp{}

func NewFGRInterp(verbose bool, p fgr.FGRProgram) *FGRInterp {
	return &FGRInterp{verboseHelper{verbose}, p}
}

func (intrp *FGRInterp) GetProgram() base.Program {
	return intrp.prog
}

func (intrp *FGRInterp) SetProgram(p base.Program) {
	intrp.prog = p.(fgr.FGRProgram)
}

func (intrp *FGRInterp) Eval(steps int) base.Program {
	return eval(intrp, steps)
}

/* FGG */

type FGGInterp struct {
	verboseHelper
	prog fgg.FGGProgram
}

var _ Interp = &FGGInterp{}

func NewFGGInterp(verbose bool, src string, strict bool) *FGGInterp {
	var a fgg.FGGAdaptor
	prog := parse(verbose, &a, src, strict).(fgg.FGGProgram)
	return &FGGInterp{verboseHelper{verbose}, prog}
}

func (intrp *FGGInterp) GetProgram() base.Program {
	return intrp.prog
}

func (intrp *FGGInterp) SetProgram(p base.Program) {
	intrp.prog = p.(fgg.FGGProgram)
}

func (intrp *FGGInterp) Eval(steps int) base.Program {
	return eval(intrp, steps)
}

// Pre: (monom == true || compile != "") => -fgg is set
// TODO: rename
func (intrp *FGGInterp) Monom(monom bool, compile string) {
	if !monom && compile == "" {
		return
	}
	p_mono := fgg.Monomorph(intrp.GetProgram().(fgg.FGGProgram))
	if monom {
		intrp.vPrintln("\nMonomorphising, formal notation: [Warning] WIP [Warning]")
		fmt.Println(p_mono.String())
	}
	if compile != "" {
		intrp.vPrintln("\nMonomorphising, FG output: [Warning] WIP [Warning]")
		out := p_mono.String()
		out = strings.Replace(out, ",,", "", -1) // TODO: refactor -- cf. fgg_monom, toMonomId
		out = strings.Replace(out, "<", "", -1)
		out = strings.Replace(out, ">", "", -1)
		if compile == "--" {
			fmt.Println(out)
		} else {
			intrp.vPrintln(out)
			intrp.vPrintln("Writing output to: " + compile)
			bs := []byte(out)
			err := ioutil.WriteFile(compile, bs, 0644)
			checkErr(err)
		}
	}
}

func (intrp *FGGInterp) Oblit(compile string) {
	if compile == "" {
		return
	}
	intrp.vPrintln("\nTranslating FGG to FG(R) using Obliteration: [Warning] WIP [Warning]")
	p_fgr := fgr.Obliterate(intrp.GetProgram().(fgg.FGGProgram))
	out := p_fgr.String()
	// TODO: factor out with -monomc
	if compile == "--" {
		fmt.Println(out)
	} else {
		intrp.vPrintln(out)
		intrp.vPrintln("Writing output to: " + compile)
		bs := []byte(out)
		err := ioutil.WriteFile(compile, bs, 0644)
		checkErr(err)
	}

	// cf. interp -- TODO: factor out with others
	p_fgr.Ok(false)
	if oblitEvalSteps > NO_EVAL {
		intrp.vPrint("\nEvaluating FGR:") // eval prints a leading "\n"
		intrp_fgr := NewFGRInterp(verbose, p_fgr)
		intrp_fgr.Eval(oblitEvalSteps)
		fmt.Println(intrp_fgr.GetProgram().GetMain())
	}
}

/* Helpers */

type verboseHelper struct {
	verbose bool
}

func (vh verboseHelper) GetVerbose() bool {
	return vh.verbose
}

func (vh verboseHelper) vPrint(x string) {
	vPrint(vh.verbose, x)
}

func (vh verboseHelper) vPrintln(x string) {
	vPrintln(vh.verbose, x)
}

func vPrint(verbose bool, x string) {
	if verbose {
		fmt.Print(x)
	}
}

func vPrintln(verbose bool, x string) {
	if verbose {
		fmt.Println(x)
	}
}

package main

import (
	"github.com/Alhanaqtah/safestructfield"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(safestructfield.Analyzer)
}

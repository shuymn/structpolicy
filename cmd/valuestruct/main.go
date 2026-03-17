package main

import (
	"github.com/shuymn/structpolicy/internal/cmd"
	"github.com/shuymn/structpolicy/pkg/valuestruct"
)

func main() {
	cmd.Run(valuestruct.NewAnalyzer())
}

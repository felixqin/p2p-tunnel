/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/felixqin/p2p-tunnel/cmd"
)

func main() {
	// 从标准输入流中接收输入数据
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	// 逐行扫描
	for input.Scan() {
		line := input.Text()
		cmd.Execute(line)
		fmt.Print("> ")
	}
}

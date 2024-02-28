package main

import (
	"fmt"
	"os"
	"os/exec"
)

func exec_script() {
	// Executando o script usando o Git Bash
	cmd := exec.Command("/bin/bash", "./deploy_stack.sh", os.Args[0], os.Args[1])

	// Definindo os canais de saída para os da aplicação
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Executando o comando
	err := cmd.Run()
	if err != nil {
		fmt.Println("Erro ao executar o script:", err)
		return
	}
}

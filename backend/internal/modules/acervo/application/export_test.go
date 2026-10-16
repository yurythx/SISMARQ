package application

import "os"

// CriarArquivo expõe aos testes a costura da criação do arquivo temporário
// (falha de disco no worker da digitalização).
var CriarArquivo *func(string) (*os.File, error) = &criarArquivo

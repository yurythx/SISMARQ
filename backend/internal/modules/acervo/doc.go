// Package acervo é um único plug-in do Kernel (DependsOn: atlas). O domínio
// está fatiado por contexto delimitado em arquivos, sem um segundo módulo:
//
//	domain/application/infrastructure/transport
//	  acervo, prazos          Fase 1 — unidade, travas, cálculo (ADR 032)
//	  fase2                   Fase 2 — empréstimo, guia, importação (ADR 033)
//	  ingresso                Fase 2 — ponte Trâmite → acervo (ADR 039)
//	  eliminacao              Fase 3 — comissões e Anexos (ADR 034)
//	  recolhimento            Fase 4 — entrada no Arquivo Público (ADR 035)
//	  normas, privados        Fase 4 — IN e arquivos privados (ADR 036)
//	  digital, pipeline       Fase 5 — PDF/A, OCR, fixity (ADR 037, 040)
//	  acesso, indice          Fase 6 — busca, consulta pública, exposições (ADR 038)
//
// Novos contextos ganham arquivo, não um plug-in. O coverage gate mede o
// pacote inteiro; fatiar em subpacotes Go só quando um contexto tiver
// repositório próprio. Ver docs/GUIDE_MODULOS.md.
package acervo

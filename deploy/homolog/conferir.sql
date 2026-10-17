-- Conferência automática do seed de homologação (C1 metadados + C6 origem).
-- Não substitui a ata humana em HOMOLOGACAO.md §5. Falha se o seed não bate.
--   make homolog-conferir

DO $$
DECLARE
  n int;
  s text;
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
     WHERE table_schema = 'public' AND table_name = 'acervo_digitais' AND column_name = 'origem'
  ) THEN
    RAISE EXCEPTION 'homolog: rode a migration 000156 (coluna origem)';
  END IF;

  SELECT COUNT(*) INTO n FROM acervo_unidades WHERE numero LIKE 'HOM-UA-%';
  IF n < 10 THEN
    RAISE EXCEPTION 'homolog: esperadas ≥10 caixas HOM-UA-*, há %', n;
  END IF;

  SELECT situacao INTO s FROM acervo_unidades WHERE numero = 'HOM-UA-01';
  IF s IS DISTINCT FROM 'AVALIAR_ELIMINACAO' THEN
    RAISE EXCEPTION 'homolog C1 HOM-UA-01: situacao %, esperado AVALIAR_ELIMINACAO', s;
  END IF;
  SELECT situacao INTO s FROM acervo_unidades WHERE numero = 'HOM-UA-02';
  IF s IS DISTINCT FROM 'RECOLHER' THEN
    RAISE EXCEPTION 'homolog C1 HOM-UA-02: situacao %, esperado RECOLHER', s;
  END IF;
  SELECT situacao INTO s FROM acervo_unidades WHERE numero = 'HOM-UA-03';
  IF s IS DISTINCT FROM 'PERMANENTE' THEN
    RAISE EXCEPTION 'homolog C1 HOM-UA-03: situacao %, esperado PERMANENTE', s;
  END IF;
  SELECT situacao INTO s FROM acervo_unidades WHERE numero = 'HOM-UA-04';
  IF s IS DISTINCT FROM 'AGUARDA_EVENTO' THEN
    RAISE EXCEPTION 'homolog C1 HOM-UA-04: situacao %, esperado AGUARDA_EVENTO', s;
  END IF;
  SELECT situacao INTO s FROM acervo_unidades WHERE numero = 'HOM-UA-07';
  IF s IS DISTINCT FROM 'BLOQUEADA' THEN
    RAISE EXCEPTION 'homolog C1 HOM-UA-07: situacao %, esperado BLOQUEADA', s;
  END IF;

  IF NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-08' AND nivel_acesso = 'RESTRITO') THEN
    RAISE EXCEPTION 'homolog C5 HOM-UA-08 deve ser RESTRITO (não entra em /arquivo)';
  END IF;

  IF NOT EXISTS (SELECT 1 FROM acervo_digitais WHERE numero = 'RD-HOM-SCAN' AND origem = 'SCAN') THEN
    RAISE EXCEPTION 'homolog C6: falta RD-HOM-SCAN (origem SCAN)';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM acervo_digitais WHERE numero = 'RD-HOM-NATO' AND origem = 'NATO_DIGITAL' AND dpi_declarado = 0) THEN
    RAISE EXCEPTION 'homolog C6: falta RD-HOM-NATO (origem NATO_DIGITAL, dpi 0)';
  END IF;

  SELECT COALESCE(MAX(version_id), 0) INTO n FROM goose_db_version WHERE is_applied;
  IF n < 157 THEN
    RAISE EXCEPTION 'homolog: goose %, precisa ≥ 157 (sem conteúdo PaperMoon)', n;
  END IF;
  SELECT COUNT(*) INTO n FROM sismarq_conteudo_anterior;
  IF n > 0 THEN
    RAISE EXCEPTION 'homolog: sismarq_conteudo_anterior ainda tem % linha(s) — rode 000157', n;
  END IF;
END $$;

SELECT numero, codigo_classificacao, fase, situacao, nivel_acesso
  FROM acervo_unidades WHERE numero LIKE 'HOM-UA-%' ORDER BY numero;

SELECT numero, origem, situacao, dpi_declarado
  FROM acervo_digitais WHERE numero LIKE 'RD-HOM-%' ORDER BY numero;

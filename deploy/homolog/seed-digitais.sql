-- Representantes digitais fictícios para homologar as telas (ADR 044).
-- Não há objeto no MinIO: situação REPROVADO para o worker não processar.
-- Pipeline de verdade: enviar o arquivo na ficha da unidade.
-- Idempotente. Pré-requisito: deploy/homolog/seed-acervo.sql
--
--   make homolog-digitais

DO $$
DECLARE
  ua_scan uuid := 'b1000000-0000-4000-8000-000000000001'; -- HOM-UA-01
  ua_nato uuid := 'b1000000-0000-4000-8000-000000000003'; -- HOM-UA-03
  resp uuid := 'c1000000-0000-4000-8000-000000000099';
  hash_scan TEXT := 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa';
  hash_nato TEXT := 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb';
  nota TEXT := '[HOMOLOGAÇÃO — envie o arquivo real pela tela para o pipeline]';
BEGIN
  IF NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE id = ua_scan) THEN
    RAISE EXCEPTION 'homolog digitais: rode make homolog-acervo antes';
  END IF;

  UPDATE acervo_unidades
     SET observacoes = CASE
       WHEN observacoes LIKE '%digitalização%' THEN observacoes
       ELSE observacoes || E'\n' || nota
     END
   WHERE id IN (ua_scan, ua_nato);

  INSERT INTO acervo_digitais (
    id, numero, unidade_id, descricao, tipo_documento, origem, equipamento, dpi_declarado,
    responsavel_id, responsavel_nome, responsavel_cpf, codigo_classificacao, digitalizado_em,
    situacao, motivo, original_key, original_sha256, original_bytes)
  SELECT 'd1000000-0000-4000-8000-000000000001', 'RD-HOM-SCAN', ua_scan,
         'Lote de homologação (papel)', 'TEXTO', 'SCAN', 'Scanner de homologação (fictício)', 300,
         resp, 'Homologação SISMARQ', '***.000.000-**', '2.0.01.00.01', now(),
         'REPROVADO',
         'Registro de teste sem arquivo no MinIO. Envie um PDF/TIFF SCAN (300 DPI) pela ficha da caixa para processar de verdade.',
         'homolog/nao-enviar/scan.pdf', hash_scan, 1024
  WHERE NOT EXISTS (SELECT 1 FROM acervo_digitais WHERE numero = 'RD-HOM-SCAN');

  INSERT INTO acervo_digitais (
    id, numero, unidade_id, descricao, tipo_documento, origem, equipamento, dpi_declarado,
    responsavel_id, responsavel_nome, responsavel_cpf, codigo_classificacao, digitalizado_em,
    situacao, motivo, original_key, original_sha256, original_bytes)
  SELECT 'd1000000-0000-4000-8000-000000000002', 'RD-HOM-NATO', ua_nato,
         'Ofício nato-digital de homologação', 'TEXTO', 'NATO_DIGITAL', 'PDF nato-digital (homologação)', 0,
         resp, 'Homologação SISMARQ', '***.000.000-**', '2.0.01.03.01', now(),
         'REPROVADO',
         'Registro de teste sem arquivo no MinIO. Envie um PDF NATO_DIGITAL pela ficha do livro para processar de verdade.',
         'homolog/nao-enviar/nato.pdf', hash_nato, 2048
  WHERE NOT EXISTS (SELECT 1 FROM acervo_digitais WHERE numero = 'RD-HOM-NATO');
END $$;

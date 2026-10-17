-- Caixas fictícias da SEMAD para homologar o painel de prazos (C1) e as
-- telas do Acervo. Códigos da TTDD oficial (deploy/ttdd). Não é acervo
-- oficial: cada unidade traz a marca [HOMOLOGAÇÃO] nas observações.
-- Idempotente. Pré-requisito: make estrutura-aplicar (SEMAD no organograma).
--
--   make homolog-acervo

DO $$
DECLARE
  semad uuid;
  predio uuid := 'a1000000-0000-4000-8000-000000000001';
  deposito uuid := 'a1000000-0000-4000-8000-000000000002';
  sala uuid := 'a1000000-0000-4000-8000-000000000003';
  estante uuid := 'a1000000-0000-4000-8000-000000000004';
  prateleira uuid := 'a1000000-0000-4000-8000-000000000005';
  marca TEXT := '[HOMOLOGAÇÃO — caixa fictícia; não é acervo oficial]';
BEGIN
  SELECT u.id INTO semad
    FROM unidades u
    JOIN entidades e ON e.id = u.entidade_id
   WHERE e.sigla = 'SEMAD' AND u.sigla = 'SEMAD'
   ORDER BY u.created_at
   LIMIT 1;
  IF semad IS NULL THEN
    RAISE EXCEPTION 'homolog acervo: SEMAD não está no organograma — rode make estrutura-aplicar';
  END IF;

  INSERT INTO acervo_locais (id, parent_id, tipo, nome)
  SELECT predio, NULL, 'PREDIO', 'Arquivo Central (homologação)'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_locais WHERE id = predio);
  INSERT INTO acervo_locais (id, parent_id, tipo, nome)
  SELECT deposito, predio, 'DEPOSITO', 'Depósito A'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_locais WHERE id = deposito);
  INSERT INTO acervo_locais (id, parent_id, tipo, nome)
  SELECT sala, deposito, 'SALA', 'Sala 1'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_locais WHERE id = sala);
  INSERT INTO acervo_locais (id, parent_id, tipo, nome)
  SELECT estante, sala, 'ESTANTE', 'Estante 01'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_locais WHERE id = estante);
  INSERT INTO acervo_locais (id, parent_id, tipo, nome)
  SELECT prateleira, estante, 'PRATELEIRA', 'Prateleira 3'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_locais WHERE id = prateleira);

  -- HOM-UA-01 Organogramas (2.0.01.00.01): 1+1 anos, eliminação, prazos cumpridos.
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000001', 'HOM-UA-01', 'CAIXA', '2.0.01.00.01',
         'Organogramas', semad, '2015-01-01', '2015-12-31', 2, 'caixas-arquivo', marca,
         'INTERMEDIARIA', 'ELIMINACAO', 'PUBLICO', prateleira, 'cx 01',
         '2016-12-31', '2017-12-31', 'AVALIAR_ELIMINACAO'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-01');

  -- HOM-UA-02 Pregão (2.0.02.00.07): 1+1 anos, guarda permanente → recolher.
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000002', 'HOM-UA-02', 'CAIXA', '2.0.02.00.07',
         'Processos relativos a Pregão Presencial/Pregão Eletrônico', semad,
         '2015-01-01', '2015-12-31', 3, 'caixas-arquivo', marca,
         'INTERMEDIARIA', 'GUARDA_PERMANENTE', 'PUBLICO', prateleira, 'cx 02',
         '2016-12-31', '2017-12-31', 'RECOLHER'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-02');

  -- HOM-UA-03 Sepultamentos (2.0.01.03.01): já permanente (arts. 20 e 42).
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000003', 'HOM-UA-03', 'LIVRO', '2.0.01.03.01',
         'Registros de sepultamento/ inumação/ exumação', semad,
         '2010-01-01', '2010-12-31', 1, 'livro', marca,
         'PERMANENTE', 'GUARDA_PERMANENTE', 'PUBLICO', prateleira, 'cx 03',
         '2011-12-31', '2012-12-31', 'PERMANENTE'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-03');

  -- HOM-UA-04 Contrato (2.0.03.00.12): “enquanto vigorando”, sem data do evento.
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000004', 'HOM-UA-04', 'PASTA', '2.0.03.00.12',
         'Contrato', semad, '2020-01-01', '2023-12-31', 1, 'pasta', marca,
         'CORRENTE', 'GUARDA_PERMANENTE', 'PUBLICO', NULL, '',
         'AGUARDA_EVENTO'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-04');

  -- HOM-UA-05 mesmo código, vigência encerrada em 2024-06-01 → transferir.
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, data_evento, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000005', 'HOM-UA-05', 'PASTA', '2.0.03.00.12',
         'Contrato', semad, '2020-01-01', '2023-12-31', 1, 'pasta', marca,
         'CORRENTE', 'GUARDA_PERMANENTE', '2024-06-01', 'PUBLICO', NULL, '',
         '2024-06-01', '2025-06-01', 'TRANSFERIR'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-05');

  -- HOM-UA-06 Organogramas ainda no corrente (datas de 2026).
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000006', 'HOM-UA-06', 'CAIXA', '2.0.01.00.01',
         'Organogramas', semad, '2026-01-01', '2026-06-30', 1, 'caixa-arquivo', marca,
         'CORRENTE', 'ELIMINACAO', 'PUBLICO', NULL, '',
         '2027-06-30', '2028-06-30', 'NO_CORRENTE'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-06');

  -- HOM-UA-07 Organogramas com bloqueio judicial (art. 44).
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000007', 'HOM-UA-07', 'CAIXA', '2.0.01.00.01',
         'Organogramas', semad, '2015-01-01', '2015-12-31', 1, 'caixa-arquivo', marca,
         'INTERMEDIARIA', 'ELIMINACAO', 'PUBLICO', prateleira, 'cx 07',
         '2016-12-31', '2017-12-31', 'BLOQUEADA'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-07');

  INSERT INTO acervo_bloqueios (id, unidade_id, processo, motivo)
  SELECT 'b2000000-0000-4000-8000-000000000001',
         'b1000000-0000-4000-8000-000000000007',
         '0001234-56.2026.8.11.0003',
         'Ação de cobrança (homologação)'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_bloqueios WHERE id = 'b2000000-0000-4000-8000-000000000001');

  -- HOM-UA-08 Sindicância (2.0.05.02.00 se existir — conferir abaixo) restrita.
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, hipotese_legal, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000008', 'HOM-UA-08', 'PROCESSO', '2.0.05.02.00',
         'Requerimento de parecer sobre abertura de Sindicância ou Processo Administrativo Disciplinar - PAD', semad,
         '2022-01-01', '2022-12-31', 1, 'processo', marca,
         'CORRENTE', 'GUARDA_PERMANENTE', 'RESTRITO',
         'LAI, art. 31 — informação pessoal (homologação; não é classificação oficial)',
         NULL, '',
         '2023-12-31', '2024-12-31', 'TRANSFERIR'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-08');

  -- HOM-UA-09 SIC 2.0.03.02.00 (C5: caminho do pedido, não o acervo ostensivo).
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000009', 'HOM-UA-09', 'CAIXA', '2.0.03.02.00',
         'Formulários de pedido de Informação – físico e eletrônico', semad,
         '2018-01-01', '2018-12-31', 1, 'caixa-arquivo', marca,
         'INTERMEDIARIA', 'ELIMINACAO', 'PUBLICO', prateleira, 'cx 09',
         '2019-12-31', '2020-12-31', 'AVALIAR_ELIMINACAO'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-09');

  -- HOM-UA-10 Contas de energia (2.0.01.00.05) com data de contas do TCE (C1).
  INSERT INTO acervo_unidades (
    id, numero, tipo, codigo_classificacao, assunto, unidade_id,
    data_inicial, data_final, quantidade, especificacao, observacoes,
    fase, destinacao, contas_aprovadas_em, nivel_acesso, local_id, posicao,
    fim_corrente, fim_intermediaria, situacao)
  SELECT 'b1000000-0000-4000-8000-000000000010', 'HOM-UA-10', 'CAIXA', '2.0.01.00.05',
         'Gestão de contas e utilização de Energia Elétrica', semad,
         '2014-01-01', '2014-12-31', 4, 'caixas-arquivo', marca,
         'INTERMEDIARIA', 'ELIMINACAO', '2016-03-15', 'PUBLICO', prateleira, 'cx 10',
         '2015-12-31', '2016-12-31', 'AVALIAR_ELIMINACAO'
  WHERE NOT EXISTS (SELECT 1 FROM acervo_unidades WHERE numero = 'HOM-UA-10');

  INSERT INTO acervo_movimentacoes (unidade_id, tipo, fase_de, fase_para, descricao)
  SELECT id, 'CADASTRO', NULL, fase, 'Carga de homologação'
    FROM acervo_unidades
   WHERE numero LIKE 'HOM-UA-%'
     AND NOT EXISTS (
       SELECT 1 FROM acervo_movimentacoes m
        WHERE m.unidade_id = acervo_unidades.id AND m.tipo = 'CADASTRO');
END $$;

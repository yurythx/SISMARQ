-- GERADO por deploy/atlas/gerar-sql-publicacao.mjs — não edite à mão.
-- Publica os 92 fluxos dos JSON como HOMOLOGADO (consulta /procedimentos).
-- Não sobrescreve código já existente (ADM.LIC.001, ADM.DIR.002, ADM.MAT.003).
BEGIN;

INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ac0750cc-b263-53b6-87ef-14469a33b8c2', 'ADM.CMP.004', 'Adesão a Ata de Registro de Preços', 'Contratar por adesão a ata de registro de preços de outro órgão (carona), com justificativa de vantajosidade e anuência do órgão gerenciador e do fornecedor.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.02.00.10'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.02.00.10' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.CMP.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b195ebb4-f21b-570c-ab76-8cda59f5c644', 'ac0750cc-b263-53b6-87ef-14469a33b8c2', 1, 'SECRETARIA', 'Setor demandante', 'Formalizar a demanda e justificar a vantagem da adesão.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2dba0e72-2e62-5423-b546-d40faa3e4e1e', 'b195ebb4-f21b-570c-ab76-8cda59f5c644', 'Documento de Formalização da Demanda (DFD)', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644' AND y.nome_documento = 'Documento de Formalização da Demanda (DFD)');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd68acea1-8841-5c59-bc06-a30286c62cea', 'b195ebb4-f21b-570c-ab76-8cda59f5c644', 'Pesquisa de preços', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644' AND y.nome_documento = 'Pesquisa de preços');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e7e03be6-9532-51b0-94ad-cab4c1226349', 'b195ebb4-f21b-570c-ab76-8cda59f5c644', 'Justificativa de vantajosidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644' AND y.nome_documento = 'Justificativa de vantajosidade');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8a8408e0-2826-54ff-82c9-74caca5d682e', 'ac0750cc-b263-53b6-87ef-14469a33b8c2', 2, 'SEMAD/COMPRAS', 'Compras', 'Solicitar a anuência do órgão gerenciador e do fornecedor e conferir a ata.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f9c5725c-165e-5455-a249-ba290c2ce044', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'Cópia da Ata de Registro de Preços', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND y.nome_documento = 'Cópia da Ata de Registro de Preços');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '07e0c65c-e0ce-57a8-a9e9-9ff23d682f0f', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'Ofício de anuência do órgão gerenciador', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND y.nome_documento = 'Ofício de anuência do órgão gerenciador');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '54aefb8d-c8a1-5e69-b32f-4fd9a7421579', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'Aceite do fornecedor', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND y.nome_documento = 'Aceite do fornecedor');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a', 'ac0750cc-b263-53b6-87ef-14469a33b8c2', 3, 'PGM', 'Procuradoria', 'Emitir parecer jurídico sobre a adesão.', 7, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9341d2a7-4067-5075-87d8-412fb3249866', 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9', 'ac0750cc-b263-53b6-87ef-14469a33b8c2', 4, 'GAB', 'Ordenador de despesa', 'Autorizar a adesão e a contratação.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'aeafa59a-a664-5b5f-a997-883761b8400f', '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9', 'Despacho de autorização', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9' AND y.nome_documento = 'Despacho de autorização');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd0d31371-7f52-545c-b715-c1d0c8535f03', 'ac0750cc-b263-53b6-87ef-14469a33b8c2', 5, 'SEMAD/CONTRATOS', 'Contratos', 'Formalizar o contrato e publicar o extrato.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac0750cc-b263-53b6-87ef-14469a33b8c2' AND x.ordem = 5);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '16ab90c3-e088-5f1a-9a2c-c7314ae6b3df', 'd0d31371-7f52-545c-b715-c1d0c8535f03', 'Contrato', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd0d31371-7f52-545c-b715-c1d0c8535f03')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd0d31371-7f52-545c-b715-c1d0c8535f03' AND y.nome_documento = 'Contrato');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7a1c8b47-739b-5b2a-a569-d47b2524e83f', 'd0d31371-7f52-545c-b715-c1d0c8535f03', 'Extrato publicado', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd0d31371-7f52-545c-b715-c1d0c8535f03')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd0d31371-7f52-545c-b715-c1d0c8535f03' AND y.nome_documento = 'Extrato publicado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b11f93df-3933-5584-8d55-3eab693b77f8', 'b195ebb4-f21b-570c-ab76-8cda59f5c644', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'Demanda formalizada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644' AND z.destino_etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND z.condicao_transicao = 'Demanda formalizada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '41cb4a04-6700-559c-8eb1-03f965475a3e', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a', 'Anuências recebidas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND z.destino_etapa_id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a' AND z.condicao_transicao = 'Anuências recebidas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5fb06e20-d0b0-54db-93cc-41279863bb79', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'b195ebb4-f21b-570c-ab76-8cda59f5c644', 'Justificativa insuficiente', TRUE, NULLIF('Complementar a justificativa de vantajosidade', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND z.destino_etapa_id = 'b195ebb4-f21b-570c-ab76-8cda59f5c644' AND z.condicao_transicao = 'Justificativa insuficiente');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5f932289-dff7-5e1a-933a-a00ee8e859ac', 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a', '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a' AND z.destino_etapa_id = '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9' AND z.condicao_transicao = 'Parecer favorável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'd2b93f6a-5d6a-59e8-a503-13206a88719e', 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a', '8a8408e0-2826-54ff-82c9-74caca5d682e', 'Parecer com ressalvas', TRUE, NULLIF('Sanear os apontamentos do parecer', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8a8408e0-2826-54ff-82c9-74caca5d682e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c8a085f6-b3a7-581d-b6c9-0e8b26a0d33a' AND z.destino_etapa_id = '8a8408e0-2826-54ff-82c9-74caca5d682e' AND z.condicao_transicao = 'Parecer com ressalvas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3cc97f63-5dc9-56c9-aa62-cf75dacfe24e', '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9', 'd0d31371-7f52-545c-b715-c1d0c8535f03', 'Adesão autorizada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd0d31371-7f52-545c-b715-c1d0c8535f03')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5dd4db3d-71aa-5ea3-9c09-10190f1f98f9' AND z.destino_etapa_id = 'd0d31371-7f52-545c-b715-c1d0c8535f03' AND z.condicao_transicao = 'Adesão autorizada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a', 'ADM.CON.005', 'Formalização de contrato administrativo e termos aditivos', 'Elaborar, assinar e publicar o contrato administrativo ou o termo aditivo (prazo, valor ou objeto) depois da contratação autorizada.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.02.02.12'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.02.02.12' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.CON.005');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3504fef6-c16c-5c3f-b07c-68074ef1b49a', 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a', 1, 'SECRETARIA', 'Setor demandante', 'Solicitar a elaboração do contrato ou do aditivo, com a justificativa.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'fe658fa9-04af-53fb-8a2f-c0e598fe2736', '3504fef6-c16c-5c3f-b07c-68074ef1b49a', 'Solicitação de elaboração de contrato/aditivo', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a' AND y.nome_documento = 'Solicitação de elaboração de contrato/aditivo');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'df111bf3-cab3-520e-83b9-e459083754be', '3504fef6-c16c-5c3f-b07c-68074ef1b49a', 'Justificativa do aditivo', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a' AND y.nome_documento = 'Justificativa do aditivo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fa7763e6-2423-553b-98ad-b60b74ba37ad', 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a', 2, 'SEMAD/CONTRATOS', 'Contratos', 'Minutar o instrumento e conferir a regularidade fiscal do contratado.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'aba29032-d340-5334-b2a9-95aa52c3d250', 'fa7763e6-2423-553b-98ad-b60b74ba37ad', 'Minuta do contrato/aditivo', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad' AND y.nome_documento = 'Minuta do contrato/aditivo');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b740aef7-f48d-5914-a9a6-8120909cc226', 'fa7763e6-2423-553b-98ad-b60b74ba37ad', 'Certidões de regularidade fiscal e trabalhista', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad' AND y.nome_documento = 'Certidões de regularidade fiscal e trabalhista');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3d40d270-edc9-55ec-982b-d937e7861032', 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a', 3, 'PGM', 'Procuradoria', 'Analisar a minuta e emitir parecer.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c9bb8b72-f0d5-56d0-9e95-daf0758fa6d8', '3d40d270-edc9-55ec-982b-d937e7861032', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3d40d270-edc9-55ec-982b-d937e7861032')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3d40d270-edc9-55ec-982b-d937e7861032' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea', 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a', 4, 'SEMAD/CONTRATOS', 'Contratos', 'Colher as assinaturas, publicar o extrato e designar o fiscal.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'df0b2901-1d8d-5816-acf7-f16d0a3ca65a' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '67c312a5-8605-5730-9b13-566cb1982e40', 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea', 'Contrato/aditivo assinado', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea' AND y.nome_documento = 'Contrato/aditivo assinado');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '905de3e3-8ea4-546c-9be1-e07510e0679f', 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea', 'Extrato publicado', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea' AND y.nome_documento = 'Extrato publicado');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3d053719-4462-5206-a9b8-77bd7feb0a18', 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea', 'Portaria de designação do fiscal', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea' AND y.nome_documento = 'Portaria de designação do fiscal');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fb3e01de-a785-5d6a-9f9a-965b957c0140', '3504fef6-c16c-5c3f-b07c-68074ef1b49a', 'fa7763e6-2423-553b-98ad-b60b74ba37ad', 'Solicitação completa', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a' AND z.destino_etapa_id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad' AND z.condicao_transicao = 'Solicitação completa');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '62f3e1ee-4c09-59a6-affb-dc3cc54793cd', 'fa7763e6-2423-553b-98ad-b60b74ba37ad', '3d40d270-edc9-55ec-982b-d937e7861032', 'Minuta pronta', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3d40d270-edc9-55ec-982b-d937e7861032')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad' AND z.destino_etapa_id = '3d40d270-edc9-55ec-982b-d937e7861032' AND z.condicao_transicao = 'Minuta pronta');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2daa838c-fda9-5bc2-9228-82997d9107b3', 'fa7763e6-2423-553b-98ad-b60b74ba37ad', '3504fef6-c16c-5c3f-b07c-68074ef1b49a', 'Documentação incompleta', TRUE, NULLIF('Enviar os documentos que faltam', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad' AND z.destino_etapa_id = '3504fef6-c16c-5c3f-b07c-68074ef1b49a' AND z.condicao_transicao = 'Documentação incompleta');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c6d80bf4-39a9-5c1e-b62d-24883b665c0d', '3d40d270-edc9-55ec-982b-d937e7861032', 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea', 'Minuta aprovada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '3d40d270-edc9-55ec-982b-d937e7861032')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '3d40d270-edc9-55ec-982b-d937e7861032' AND z.destino_etapa_id = 'd1c49b1f-2c1e-5513-aba8-4ddd8cebefea' AND z.condicao_transicao = 'Minuta aprovada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3bfe21ec-6fac-55f2-995e-924b8f74a17c', '3d40d270-edc9-55ec-982b-d937e7861032', 'fa7763e6-2423-553b-98ad-b60b74ba37ad', 'Ajustes na minuta', TRUE, NULLIF('Ajustar a minuta conforme o parecer', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '3d40d270-edc9-55ec-982b-d937e7861032')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '3d40d270-edc9-55ec-982b-d937e7861032' AND z.destino_etapa_id = 'fa7763e6-2423-553b-98ad-b60b74ba37ad' AND z.condicao_transicao = 'Ajustes na minuta');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '5c0d1d56-2424-54d7-b2b9-53589bc07998', 'ADM.CON.006', 'Fiscalização e acompanhamento de contrato', 'Acompanhar a execução do contrato pelo fiscal designado, com relatórios periódicos e atesto para pagamento.', 'Fiscais de contrato', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.02.02.15'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.02.02.15' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.CON.006');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8457ff4d-9b02-583e-9607-9d39e662c5dc', '5c0d1d56-2424-54d7-b2b9-53589bc07998', 1, 'SECRETARIA', 'Fiscal do contrato', 'Acompanhar a execução, registrar ocorrências e emitir o relatório mensal.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5c0d1d56-2424-54d7-b2b9-53589bc07998')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5c0d1d56-2424-54d7-b2b9-53589bc07998' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c8fcc407-fc22-5d69-9648-3c82b154edef', '8457ff4d-9b02-583e-9607-9d39e662c5dc', 'Relatório mensal do fiscal', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8457ff4d-9b02-583e-9607-9d39e662c5dc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8457ff4d-9b02-583e-9607-9d39e662c5dc' AND y.nome_documento = 'Relatório mensal do fiscal');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b6d61994-5ca4-519b-965a-9246910bc2fd', '8457ff4d-9b02-583e-9607-9d39e662c5dc', 'Registro de ocorrências', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8457ff4d-9b02-583e-9607-9d39e662c5dc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8457ff4d-9b02-583e-9607-9d39e662c5dc' AND y.nome_documento = 'Registro de ocorrências');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e30d8770-00ea-50d8-8dea-040bd4a6ee37', '5c0d1d56-2424-54d7-b2b9-53589bc07998', 2, 'SECRETARIA', 'Gestor do contrato', 'Atestar a nota fiscal e encaminhar para pagamento ou notificar o contratado.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5c0d1d56-2424-54d7-b2b9-53589bc07998')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5c0d1d56-2424-54d7-b2b9-53589bc07998' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '35ced2a9-56b9-5ff1-ac7c-ecb6aef83314', 'e30d8770-00ea-50d8-8dea-040bd4a6ee37', 'Nota fiscal atestada', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37' AND y.nome_documento = 'Nota fiscal atestada');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ada57e6a-59ea-561b-b73c-ffc9e6bc0f56', 'e30d8770-00ea-50d8-8dea-040bd4a6ee37', 'Notificação ao contratado', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37' AND y.nome_documento = 'Notificação ao contratado');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '956f6127-f045-5b4a-8859-4f7dcdedc6ad', '5c0d1d56-2424-54d7-b2b9-53589bc07998', 3, 'SEFIN/CONTABILIDADE', 'Contabilidade', 'Liquidar e pagar a despesa.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5c0d1d56-2424-54d7-b2b9-53589bc07998')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5c0d1d56-2424-54d7-b2b9-53589bc07998' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9b035153-33cf-5f71-bfd9-6b3b5c207f04', '956f6127-f045-5b4a-8859-4f7dcdedc6ad', 'Nota de liquidação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '956f6127-f045-5b4a-8859-4f7dcdedc6ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '956f6127-f045-5b4a-8859-4f7dcdedc6ad' AND y.nome_documento = 'Nota de liquidação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '94d7a02d-6483-5402-8586-c016c3ff3fb6', '956f6127-f045-5b4a-8859-4f7dcdedc6ad', 'Ordem de pagamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '956f6127-f045-5b4a-8859-4f7dcdedc6ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '956f6127-f045-5b4a-8859-4f7dcdedc6ad' AND y.nome_documento = 'Ordem de pagamento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5366c6ca-944e-559f-9b34-0470a28b47d8', '8457ff4d-9b02-583e-9607-9d39e662c5dc', 'e30d8770-00ea-50d8-8dea-040bd4a6ee37', 'Execução conforme', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8457ff4d-9b02-583e-9607-9d39e662c5dc')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8457ff4d-9b02-583e-9607-9d39e662c5dc' AND z.destino_etapa_id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37' AND z.condicao_transicao = 'Execução conforme');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '143b7942-d98d-5271-a390-acf1801f888c', 'e30d8770-00ea-50d8-8dea-040bd4a6ee37', '956f6127-f045-5b4a-8859-4f7dcdedc6ad', 'Nota atestada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '956f6127-f045-5b4a-8859-4f7dcdedc6ad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37' AND z.destino_etapa_id = '956f6127-f045-5b4a-8859-4f7dcdedc6ad' AND z.condicao_transicao = 'Nota atestada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '9f780553-1828-5a17-98cc-2e2563663e3d', 'e30d8770-00ea-50d8-8dea-040bd4a6ee37', '8457ff4d-9b02-583e-9607-9d39e662c5dc', 'Inexecução ou falha', TRUE, NULLIF('Notificar o contratado e registrar a ocorrência', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8457ff4d-9b02-583e-9607-9d39e662c5dc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e30d8770-00ea-50d8-8dea-040bd4a6ee37' AND z.destino_etapa_id = '8457ff4d-9b02-583e-9607-9d39e662c5dc' AND z.condicao_transicao = 'Inexecução ou falha');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '52abaef2-f21f-56ca-b9ad-b489636df9e9', 'ADM.CMP.007', 'Contratação por inexigibilidade de licitação', 'Contratar diretamente quando a competição é inviável (fornecedor exclusivo, notória especialização, credenciamento — art. 74 da Lei 14.133/2021).', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.02.01.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.02.01.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.CMP.007');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '82393821-d121-56dd-9b34-0492cd19fee8', '52abaef2-f21f-56ca-b9ad-b489636df9e9', 1, 'SECRETARIA', 'Setor demandante', 'Formalizar a demanda, o estudo técnico preliminar e a justificativa da inviabilidade de competição.', 7, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '52abaef2-f21f-56ca-b9ad-b489636df9e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '52abaef2-f21f-56ca-b9ad-b489636df9e9' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7a821ef0-e4f5-5aaf-ab30-3eda53a9dfd5', '82393821-d121-56dd-9b34-0492cd19fee8', 'Documento de Formalização da Demanda (DFD)', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '82393821-d121-56dd-9b34-0492cd19fee8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '82393821-d121-56dd-9b34-0492cd19fee8' AND y.nome_documento = 'Documento de Formalização da Demanda (DFD)');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e32c45ca-75f0-5e23-af23-7ba3674c204a', '82393821-d121-56dd-9b34-0492cd19fee8', 'Estudo Técnico Preliminar (ETP)', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '82393821-d121-56dd-9b34-0492cd19fee8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '82393821-d121-56dd-9b34-0492cd19fee8' AND y.nome_documento = 'Estudo Técnico Preliminar (ETP)');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ab10974b-54d1-5cd5-827c-c18895c36709', '82393821-d121-56dd-9b34-0492cd19fee8', 'Termo de Referência', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '82393821-d121-56dd-9b34-0492cd19fee8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '82393821-d121-56dd-9b34-0492cd19fee8' AND y.nome_documento = 'Termo de Referência');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '42c0734a-9cdf-5435-89c2-562e6c209ac4', '82393821-d121-56dd-9b34-0492cd19fee8', 'Comprovação de exclusividade ou notória especialização', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '82393821-d121-56dd-9b34-0492cd19fee8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '82393821-d121-56dd-9b34-0492cd19fee8' AND y.nome_documento = 'Comprovação de exclusividade ou notória especialização');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ed5586e9-8209-5287-b736-1cac5e8e2c8e', '52abaef2-f21f-56ca-b9ad-b489636df9e9', 2, 'SEMAD/COMPRAS', 'Compras', 'Justificar o preço e conferir a habilitação do contratado.', 7, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '52abaef2-f21f-56ca-b9ad-b489636df9e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '52abaef2-f21f-56ca-b9ad-b489636df9e9' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ce567f59-86af-51fa-af5b-1006cc750b37', 'ed5586e9-8209-5287-b736-1cac5e8e2c8e', 'Justificativa de preço', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e' AND y.nome_documento = 'Justificativa de preço');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'fc8afe05-b9f2-54a5-bad2-ad993981c15b', 'ed5586e9-8209-5287-b736-1cac5e8e2c8e', 'Documentos de habilitação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e' AND y.nome_documento = 'Documentos de habilitação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0511344c-921d-5520-88cb-f9b0dbe7eb6a', '52abaef2-f21f-56ca-b9ad-b489636df9e9', 3, 'PGM', 'Procuradoria', 'Emitir parecer jurídico.', 7, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '52abaef2-f21f-56ca-b9ad-b489636df9e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '52abaef2-f21f-56ca-b9ad-b489636df9e9' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '656142dc-e7d6-597b-8f85-8bff10daf548', '0511344c-921d-5520-88cb-f9b0dbe7eb6a', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0511344c-921d-5520-88cb-f9b0dbe7eb6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0511344c-921d-5520-88cb-f9b0dbe7eb6a' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041', '52abaef2-f21f-56ca-b9ad-b489636df9e9', 4, 'GAB', 'Autoridade competente', 'Ratificar a inexigibilidade e autorizar a contratação.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '52abaef2-f21f-56ca-b9ad-b489636df9e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '52abaef2-f21f-56ca-b9ad-b489636df9e9' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '42d107e4-9238-550a-b8e5-00f277513493', 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041', 'Ato de ratificação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041' AND y.nome_documento = 'Ato de ratificação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3995179c-3d5b-5ae1-8f13-a954f4f0e032', 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041', 'Publicação do ato', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041' AND y.nome_documento = 'Publicação do ato');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0c8fb993-251a-52b0-adae-d1fa966a1da3', '82393821-d121-56dd-9b34-0492cd19fee8', 'ed5586e9-8209-5287-b736-1cac5e8e2c8e', 'Instrução completa', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '82393821-d121-56dd-9b34-0492cd19fee8')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '82393821-d121-56dd-9b34-0492cd19fee8' AND z.destino_etapa_id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e' AND z.condicao_transicao = 'Instrução completa');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3e01c0ec-e541-5a6d-8606-500e29c463d7', 'ed5586e9-8209-5287-b736-1cac5e8e2c8e', '0511344c-921d-5520-88cb-f9b0dbe7eb6a', 'Preço justificado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '0511344c-921d-5520-88cb-f9b0dbe7eb6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e' AND z.destino_etapa_id = '0511344c-921d-5520-88cb-f9b0dbe7eb6a' AND z.condicao_transicao = 'Preço justificado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'd5b12884-2f40-58af-872f-5e575a32646f', 'ed5586e9-8209-5287-b736-1cac5e8e2c8e', '82393821-d121-56dd-9b34-0492cd19fee8', 'Instrução incompleta', TRUE, NULLIF('Complementar a instrução', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '82393821-d121-56dd-9b34-0492cd19fee8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ed5586e9-8209-5287-b736-1cac5e8e2c8e' AND z.destino_etapa_id = '82393821-d121-56dd-9b34-0492cd19fee8' AND z.condicao_transicao = 'Instrução incompleta');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2f86b6c0-9c5d-5665-abdb-7b0af4a225a6', '0511344c-921d-5520-88cb-f9b0dbe7eb6a', 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0511344c-921d-5520-88cb-f9b0dbe7eb6a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0511344c-921d-5520-88cb-f9b0dbe7eb6a' AND z.destino_etapa_id = 'e72e500c-7d9f-5b9e-b3ce-6bb18ff19041' AND z.condicao_transicao = 'Parecer favorável');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'dbad10e9-a86a-5876-878b-ea1ab9bccf49', 'ADM.FRO.008', 'Infração de trânsito de veículo oficial', 'Identificar o condutor responsável por infração com veículo oficial, apresentar defesa quando cabível e ressarcir a multa.', 'Condutores e gestores de frota', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.01.02.15'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.01.02.15' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.FRO.008');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '98a39b6e-2b69-5a9c-8350-c4865cd4a84f', 'dbad10e9-a86a-5876-878b-ea1ab9bccf49', 1, 'SEMAD/FROTAS', 'Gestão de frotas', 'Receber a notificação, identificar o condutor pelo diário de bordo e indicá-lo ao órgão de trânsito.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'dbad10e9-a86a-5876-878b-ea1ab9bccf49')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'dbad10e9-a86a-5876-878b-ea1ab9bccf49' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'baf4ef50-aa14-5396-89db-ce52e497a280', '98a39b6e-2b69-5a9c-8350-c4865cd4a84f', 'Notificação de autuação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f' AND y.nome_documento = 'Notificação de autuação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3d3c40ac-77d4-575e-b10c-567eda0944a8', '98a39b6e-2b69-5a9c-8350-c4865cd4a84f', 'Diário de bordo do período', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f' AND y.nome_documento = 'Diário de bordo do período');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '200d9e92-5a29-5d86-9db3-02d2a14e2e27', '98a39b6e-2b69-5a9c-8350-c4865cd4a84f', 'Formulário de indicação de condutor', TRUE, 'EXTERNO_DIGITALIZADO', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f' AND y.nome_documento = 'Formulário de indicação de condutor');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'bab6dee3-1bcb-5536-a5a6-959c8027690f', 'dbad10e9-a86a-5876-878b-ea1ab9bccf49', 2, 'SECRETARIA', 'Chefia do condutor', 'Colher a defesa do condutor ou a concordância com o desconto.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'dbad10e9-a86a-5876-878b-ea1ab9bccf49')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'dbad10e9-a86a-5876-878b-ea1ab9bccf49' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5be2394e-3d7e-507c-a539-eeaf695f9267', 'bab6dee3-1bcb-5536-a5a6-959c8027690f', 'Defesa do condutor', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f' AND y.nome_documento = 'Defesa do condutor');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c966682a-d729-5a14-95f0-8d5c09767eed', 'bab6dee3-1bcb-5536-a5a6-959c8027690f', 'Termo de ciência e autorização de desconto', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f' AND y.nome_documento = 'Termo de ciência e autorização de desconto');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8058530e-7072-5b25-998d-a8284ddcb8d2', 'dbad10e9-a86a-5876-878b-ea1ab9bccf49', 3, 'SEGEP/FOLHA', 'Folha de pagamento', 'Descontar em folha o valor da multa, quando devido.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'dbad10e9-a86a-5876-878b-ea1ab9bccf49')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'dbad10e9-a86a-5876-878b-ea1ab9bccf49' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '12373b80-83f0-5bcb-afe8-029bc4c48f15', '8058530e-7072-5b25-998d-a8284ddcb8d2', 'Comprovante de desconto', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8058530e-7072-5b25-998d-a8284ddcb8d2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8058530e-7072-5b25-998d-a8284ddcb8d2' AND y.nome_documento = 'Comprovante de desconto');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '522a720c-a009-5f99-9cf6-5ce7a76daab7', '98a39b6e-2b69-5a9c-8350-c4865cd4a84f', 'bab6dee3-1bcb-5536-a5a6-959c8027690f', 'Condutor identificado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '98a39b6e-2b69-5a9c-8350-c4865cd4a84f' AND z.destino_etapa_id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f' AND z.condicao_transicao = 'Condutor identificado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b823e075-25bc-562e-8d69-4a0e81b336a3', 'bab6dee3-1bcb-5536-a5a6-959c8027690f', '8058530e-7072-5b25-998d-a8284ddcb8d2', 'Manifestação recebida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8058530e-7072-5b25-998d-a8284ddcb8d2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'bab6dee3-1bcb-5536-a5a6-959c8027690f' AND z.destino_etapa_id = '8058530e-7072-5b25-998d-a8284ddcb8d2' AND z.condicao_transicao = 'Manifestação recebida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'cefe87fe-e077-5741-bd54-8e882e34fd3f', 'ADM.SIC.009', 'Pedido de acesso à informação (LAI)', 'Atender o pedido de acesso à informação no prazo da Lei 12.527/2011 (20 dias, prorrogáveis por mais 10), com recurso à autoridade superior.', 'Cidadãos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.03.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.03.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.SIC.009');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4103a426-264b-5d08-b034-1b1eb1070e25', 'cefe87fe-e077-5741-bd54-8e882e34fd3f', 1, 'SEMAD/SIC', 'Serviço de Informação ao Cidadão', 'Registrar o pedido e encaminhar à unidade que detém a informação.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'cefe87fe-e077-5741-bd54-8e882e34fd3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'cefe87fe-e077-5741-bd54-8e882e34fd3f' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '86544a26-241a-5676-bc92-a077fee8fbe2', '4103a426-264b-5d08-b034-1b1eb1070e25', 'Formulário de pedido de informação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4103a426-264b-5d08-b034-1b1eb1070e25')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4103a426-264b-5d08-b034-1b1eb1070e25' AND y.nome_documento = 'Formulário de pedido de informação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4', 'cefe87fe-e077-5741-bd54-8e882e34fd3f', 2, 'SECRETARIA', 'Unidade detentora', 'Localizar a informação ou justificar o sigilo/indisponibilidade.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'cefe87fe-e077-5741-bd54-8e882e34fd3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'cefe87fe-e077-5741-bd54-8e882e34fd3f' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '076967c8-c162-5e26-afcf-9b54c937ff75', '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4', 'Resposta da unidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4' AND y.nome_documento = 'Resposta da unidade');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3744c961-9643-5768-ba08-278b89af48ae', 'cefe87fe-e077-5741-bd54-8e882e34fd3f', 3, 'SEMAD/SIC', 'Serviço de Informação ao Cidadão', 'Enviar a resposta ao cidadão e informar o prazo de recurso.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'cefe87fe-e077-5741-bd54-8e882e34fd3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'cefe87fe-e077-5741-bd54-8e882e34fd3f' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4211983b-df7c-5856-973f-9d5a0cdd96fa', '3744c961-9643-5768-ba08-278b89af48ae', 'Carta resposta ao pedido de informação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3744c961-9643-5768-ba08-278b89af48ae')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3744c961-9643-5768-ba08-278b89af48ae' AND y.nome_documento = 'Carta resposta ao pedido de informação');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c827d6fd-5648-5d63-b9e4-cbb6a2cedaa1', '4103a426-264b-5d08-b034-1b1eb1070e25', '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4', 'Pedido encaminhado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4103a426-264b-5d08-b034-1b1eb1070e25')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4103a426-264b-5d08-b034-1b1eb1070e25' AND z.destino_etapa_id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4' AND z.condicao_transicao = 'Pedido encaminhado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'ad7dce25-ec12-512d-bffd-043372f68fd4', '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4', '3744c961-9643-5768-ba08-278b89af48ae', 'Resposta pronta', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3744c961-9643-5768-ba08-278b89af48ae')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4' AND z.destino_etapa_id = '3744c961-9643-5768-ba08-278b89af48ae' AND z.condicao_transicao = 'Resposta pronta');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f1641b6b-688e-5590-be38-c7e3adeaaa22', '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4', '4103a426-264b-5d08-b034-1b1eb1070e25', 'Pedido genérico ou incompreensível', TRUE, NULLIF('Pedir ao cidadão que esclareça o pedido', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '4103a426-264b-5d08-b034-1b1eb1070e25')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0bde68e2-b32c-5e4d-9696-d7a7f7ee3bd4' AND z.destino_etapa_id = '4103a426-264b-5d08-b034-1b1eb1070e25' AND z.condicao_transicao = 'Pedido genérico ou incompreensível');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'af2c0202-3c55-5626-be45-9b865f5d824d', 'RH.FER.001', 'Férias do servidor', 'Programar e conceder o gozo de férias do servidor, com o adicional de 1/3 na folha.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.03.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.03.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.FER.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8', 'af2c0202-3c55-5626-be45-9b865f5d824d', 1, 'SECRETARIA', 'Chefia imediata', 'Aprovar o período de gozo conforme a escala da unidade.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'af2c0202-3c55-5626-be45-9b865f5d824d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'af2c0202-3c55-5626-be45-9b865f5d824d' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c1c121d1-31a4-5168-b52b-769712b61c88', '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8', 'Requerimento de gozo de férias', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8' AND y.nome_documento = 'Requerimento de gozo de férias');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '22d96572-33f8-5e55-ad1c-faf070a93d56', 'af2c0202-3c55-5626-be45-9b865f5d824d', 2, 'SEGEP/RH', 'Recursos Humanos', 'Conferir o período aquisitivo e lançar as férias.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'af2c0202-3c55-5626-be45-9b865f5d824d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'af2c0202-3c55-5626-be45-9b865f5d824d' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ed68471a-9d2f-592c-bf70-5d810f473131', '22d96572-33f8-5e55-ad1c-faf070a93d56', 'Aviso de férias', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '22d96572-33f8-5e55-ad1c-faf070a93d56')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '22d96572-33f8-5e55-ad1c-faf070a93d56' AND y.nome_documento = 'Aviso de férias');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'db2d69bd-70b7-57ae-93bc-91e8857893f7', 'af2c0202-3c55-5626-be45-9b865f5d824d', 3, 'SEGEP/FOLHA', 'Folha de pagamento', 'Pagar o adicional de 1/3 na folha do mês anterior ao gozo.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'af2c0202-3c55-5626-be45-9b865f5d824d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'af2c0202-3c55-5626-be45-9b865f5d824d' AND x.ordem = 3);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '8b4dc3c4-5256-5574-95d0-3bfc11620f69', '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8', '22d96572-33f8-5e55-ad1c-faf070a93d56', 'Período aprovado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '22d96572-33f8-5e55-ad1c-faf070a93d56')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8' AND z.destino_etapa_id = '22d96572-33f8-5e55-ad1c-faf070a93d56' AND z.condicao_transicao = 'Período aprovado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f113d918-ed85-5bf9-a19b-76a91481f0e3', '22d96572-33f8-5e55-ad1c-faf070a93d56', 'db2d69bd-70b7-57ae-93bc-91e8857893f7', 'Férias lançadas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '22d96572-33f8-5e55-ad1c-faf070a93d56')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'db2d69bd-70b7-57ae-93bc-91e8857893f7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '22d96572-33f8-5e55-ad1c-faf070a93d56' AND z.destino_etapa_id = 'db2d69bd-70b7-57ae-93bc-91e8857893f7' AND z.condicao_transicao = 'Férias lançadas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4ae412d9-8dbc-5c00-bb42-6286ddb66e8f', '22d96572-33f8-5e55-ad1c-faf070a93d56', '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8', 'Período aquisitivo incompleto', TRUE, NULLIF('Ajustar o período de gozo', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '22d96572-33f8-5e55-ad1c-faf070a93d56')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '22d96572-33f8-5e55-ad1c-faf070a93d56' AND z.destino_etapa_id = '766ee0bc-3023-5460-aa4b-2dcc5c86a2e8' AND z.condicao_transicao = 'Período aquisitivo incompleto');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa', 'RH.LIC.002', 'Licença-prêmio por assiduidade', 'Reconhecer o direito à licença-prêmio e programar o gozo, conforme o estatuto dos servidores.', 'Servidores efetivos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.03.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.03.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.LIC.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0', '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber o requerimento e certificar o tempo de serviço e as ausências do período.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'da9238f6-62a0-536c-9cb5-14512bdc4f2b', '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0', 'Requerimento de usufruto de licença-prêmio', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0' AND y.nome_documento = 'Requerimento de usufruto de licença-prêmio');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '38d14cfb-70e1-5ba7-82cb-df0ff7ef78db', '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0', 'Certidão de tempo de serviço', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0' AND y.nome_documento = 'Certidão de tempo de serviço');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195', '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa', 2, 'PGM', 'Procuradoria', 'Analisar casos com afastamentos que interrompem o período (quando houver).', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3f8d57b0-b005-515d-9e57-7fd11fc4cd31', '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195', 'Parecer jurídico', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '27973a77-cb1f-5404-81cf-4ca243a56734', '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa', 3, 'SEGEP/RH', 'Recursos Humanos', 'Deferir ou indeferir e publicar o ato; programar o gozo com a chefia.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '01ed6e4a-783e-5a50-b829-9c3d5f3b8aaa' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8ce54398-7e84-5117-97b4-36ad70fe374b', '27973a77-cb1f-5404-81cf-4ca243a56734', 'Portaria de concessão', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '27973a77-cb1f-5404-81cf-4ca243a56734')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '27973a77-cb1f-5404-81cf-4ca243a56734' AND y.nome_documento = 'Portaria de concessão');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9725f304-de81-5e37-ae77-ea1f146de6e8', '27973a77-cb1f-5404-81cf-4ca243a56734', 'Requerimento de gozo de licença-prêmio', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '27973a77-cb1f-5404-81cf-4ca243a56734')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '27973a77-cb1f-5404-81cf-4ca243a56734' AND y.nome_documento = 'Requerimento de gozo de licença-prêmio');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '99747c89-94ac-5330-a510-7b01b145a31e', '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0', '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195', 'Tempo certificado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '34140d6e-88c3-55e9-9a1d-5c5a65e4abe0' AND z.destino_etapa_id = '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195' AND z.condicao_transicao = 'Tempo certificado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '23e376a0-7ace-5f7f-a297-9bf0ad158cab', '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195', '27973a77-cb1f-5404-81cf-4ca243a56734', 'Direito reconhecido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '27973a77-cb1f-5404-81cf-4ca243a56734')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5bb2d7ee-5f2a-5eeb-ac4e-f9a048f5a195' AND z.destino_etapa_id = '27973a77-cb1f-5404-81cf-4ca243a56734' AND z.condicao_transicao = 'Direito reconhecido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '480315e3-570a-5669-a669-b1dac21197a3', 'RH.ADI.003', 'Adicional de insalubridade', 'Conceder o adicional de insalubridade com base em laudo técnico das condições de trabalho.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.02.09'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.02.09' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.ADI.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c94caffc-338c-547a-93b0-f0c3e10087a9', '480315e3-570a-5669-a669-b1dac21197a3', 1, 'SECRETARIA', 'Chefia imediata', 'Encaminhar o requerimento com a descrição das atividades e do local.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '480315e3-570a-5669-a669-b1dac21197a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '480315e3-570a-5669-a669-b1dac21197a3' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '176a1279-1e24-535f-8843-f5b0ef4a54ce', 'c94caffc-338c-547a-93b0-f0c3e10087a9', 'Requerimento de adicional de insalubridade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c94caffc-338c-547a-93b0-f0c3e10087a9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c94caffc-338c-547a-93b0-f0c3e10087a9' AND y.nome_documento = 'Requerimento de adicional de insalubridade');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '66a19dea-3699-525c-8d89-91763dfbfc2d', 'c94caffc-338c-547a-93b0-f0c3e10087a9', 'Descrição das atividades', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c94caffc-338c-547a-93b0-f0c3e10087a9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c94caffc-338c-547a-93b0-f0c3e10087a9' AND y.nome_documento = 'Descrição das atividades');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'cb534557-9662-5987-a0d3-2c086718f4ea', '480315e3-570a-5669-a669-b1dac21197a3', 2, 'SEGEP/DESOPEM', 'Saúde ocupacional (DESOPEM)', 'Emitir o laudo técnico de insalubridade (LTCAT) e o grau.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '480315e3-570a-5669-a669-b1dac21197a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '480315e3-570a-5669-a669-b1dac21197a3' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'fbbaa6dc-cf82-5b87-96f9-d34909bcead1', 'cb534557-9662-5987-a0d3-2c086718f4ea', 'Laudo técnico (LTCAT)', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cb534557-9662-5987-a0d3-2c086718f4ea')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cb534557-9662-5987-a0d3-2c086718f4ea' AND y.nome_documento = 'Laudo técnico (LTCAT)');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '385cd233-4aca-50c6-accd-df6fd0ef8acf', '480315e3-570a-5669-a669-b1dac21197a3', 3, 'SEGEP/RH', 'Recursos Humanos', 'Conceder o adicional e lançar na folha.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '480315e3-570a-5669-a669-b1dac21197a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '480315e3-570a-5669-a669-b1dac21197a3' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8259d22f-9211-5607-8704-6def3b23acd9', '385cd233-4aca-50c6-accd-df6fd0ef8acf', 'Portaria de concessão', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '385cd233-4aca-50c6-accd-df6fd0ef8acf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '385cd233-4aca-50c6-accd-df6fd0ef8acf' AND y.nome_documento = 'Portaria de concessão');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3f8ca386-7388-5e93-96c1-468b7a1f32d0', 'c94caffc-338c-547a-93b0-f0c3e10087a9', 'cb534557-9662-5987-a0d3-2c086718f4ea', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c94caffc-338c-547a-93b0-f0c3e10087a9')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'cb534557-9662-5987-a0d3-2c086718f4ea')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c94caffc-338c-547a-93b0-f0c3e10087a9' AND z.destino_etapa_id = 'cb534557-9662-5987-a0d3-2c086718f4ea' AND z.condicao_transicao = 'Requerimento instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '7e0cd093-6e34-51ff-995c-a369ff14e2db', 'cb534557-9662-5987-a0d3-2c086718f4ea', '385cd233-4aca-50c6-accd-df6fd0ef8acf', 'Laudo emitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'cb534557-9662-5987-a0d3-2c086718f4ea')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '385cd233-4aca-50c6-accd-df6fd0ef8acf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'cb534557-9662-5987-a0d3-2c086718f4ea' AND z.destino_etapa_id = '385cd233-4aca-50c6-accd-df6fd0ef8acf' AND z.condicao_transicao = 'Laudo emitido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f', 'RH.APO.004', 'Aposentadoria voluntária por tempo de contribuição', 'Instruir e conceder a aposentadoria voluntária do servidor efetivo pelo regime próprio (IMPRO).', 'Servidores efetivos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.07.01.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.01.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.APO.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'acf58cf9-0b15-5456-bd27-b055b71922a3', 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber o requerimento e emitir a certidão de tempo de contribuição e o histórico funcional.', 20, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'af76ffc6-5292-5d95-9a38-1fdf9dc4c66e', 'acf58cf9-0b15-5456-bd27-b055b71922a3', 'Requerimento de aposentadoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'acf58cf9-0b15-5456-bd27-b055b71922a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'acf58cf9-0b15-5456-bd27-b055b71922a3' AND y.nome_documento = 'Requerimento de aposentadoria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '718a8b33-fb7b-5010-8cf2-a7d713a9bd66', 'acf58cf9-0b15-5456-bd27-b055b71922a3', 'Certidão de tempo de contribuição', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'acf58cf9-0b15-5456-bd27-b055b71922a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'acf58cf9-0b15-5456-bd27-b055b71922a3' AND y.nome_documento = 'Certidão de tempo de contribuição');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '482a51cc-b4cd-5ff3-9143-3bff0287e1a3', 'acf58cf9-0b15-5456-bd27-b055b71922a3', 'Ficha financeira', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'acf58cf9-0b15-5456-bd27-b055b71922a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'acf58cf9-0b15-5456-bd27-b055b71922a3' AND y.nome_documento = 'Ficha financeira');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '27250140-c3a7-58cf-8ef2-50c7abeceac5', 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f', 2, 'IMPRO', 'Instituto de Previdência (IMPRO)', 'Calcular os proventos e conferir os requisitos.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'fd699775-3333-5fbc-9016-c39799b9a2a6', '27250140-c3a7-58cf-8ef2-50c7abeceac5', 'Cálculo de proventos', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '27250140-c3a7-58cf-8ef2-50c7abeceac5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '27250140-c3a7-58cf-8ef2-50c7abeceac5' AND y.nome_documento = 'Cálculo de proventos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fae6f4d5-75c3-5558-af89-b72457aa8609', 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f', 3, 'PGM', 'Procuradoria', 'Emitir parecer jurídico.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2d06050e-77f8-5917-9323-39154fd7e495', 'fae6f4d5-75c3-5558-af89-b72457aa8609', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fae6f4d5-75c3-5558-af89-b72457aa8609')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fae6f4d5-75c3-5558-af89-b72457aa8609' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ef649d7d-49b9-550e-9726-f5c97ed52ee9', 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f', 4, 'GAB', 'Prefeito', 'Assinar o ato de aposentadoria e publicar.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'dc21ddd8-92fc-53a8-ac13-7d3e6fb42ce0', 'ef649d7d-49b9-550e-9726-f5c97ed52ee9', 'Portaria de aposentadoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9' AND y.nome_documento = 'Portaria de aposentadoria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '56633a47-b98c-55c9-9d6a-348fe416b054', 'ef649d7d-49b9-550e-9726-f5c97ed52ee9', 'Publicação no Diário Oficial', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9' AND y.nome_documento = 'Publicação no Diário Oficial');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '157b62c6-09c6-5a7b-8331-0ea3ecb2a601', 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f', 5, 'IMPRO', 'Instituto de Previdência (IMPRO)', 'Encaminhar o ato ao Tribunal de Contas para registro.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd4ef36b4-8775-5d33-a3cd-8d3636029a3f' AND x.ordem = 5);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4b6d9ec8-49ae-5560-9779-fc1df46d2bd6', '157b62c6-09c6-5a7b-8331-0ea3ecb2a601', 'Ofício ao TCE', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '157b62c6-09c6-5a7b-8331-0ea3ecb2a601')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '157b62c6-09c6-5a7b-8331-0ea3ecb2a601' AND y.nome_documento = 'Ofício ao TCE');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e3f0c7bd-f117-5b6b-9b84-4f6c59fa790a', 'acf58cf9-0b15-5456-bd27-b055b71922a3', '27250140-c3a7-58cf-8ef2-50c7abeceac5', 'Tempo certificado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'acf58cf9-0b15-5456-bd27-b055b71922a3')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '27250140-c3a7-58cf-8ef2-50c7abeceac5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'acf58cf9-0b15-5456-bd27-b055b71922a3' AND z.destino_etapa_id = '27250140-c3a7-58cf-8ef2-50c7abeceac5' AND z.condicao_transicao = 'Tempo certificado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '020445d2-5406-5f35-b26a-f1da1ede1616', '27250140-c3a7-58cf-8ef2-50c7abeceac5', 'fae6f4d5-75c3-5558-af89-b72457aa8609', 'Requisitos atendidos', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '27250140-c3a7-58cf-8ef2-50c7abeceac5')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fae6f4d5-75c3-5558-af89-b72457aa8609')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '27250140-c3a7-58cf-8ef2-50c7abeceac5' AND z.destino_etapa_id = 'fae6f4d5-75c3-5558-af89-b72457aa8609' AND z.condicao_transicao = 'Requisitos atendidos');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '981e124b-93d9-518c-8719-bece8ebce5ed', '27250140-c3a7-58cf-8ef2-50c7abeceac5', 'acf58cf9-0b15-5456-bd27-b055b71922a3', 'Tempo ou documento faltante', TRUE, NULLIF('Averbar o tempo ou juntar o documento', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '27250140-c3a7-58cf-8ef2-50c7abeceac5')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'acf58cf9-0b15-5456-bd27-b055b71922a3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '27250140-c3a7-58cf-8ef2-50c7abeceac5' AND z.destino_etapa_id = 'acf58cf9-0b15-5456-bd27-b055b71922a3' AND z.condicao_transicao = 'Tempo ou documento faltante');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '767b9c67-fc87-59fe-afa1-8b6a9011fc47', 'fae6f4d5-75c3-5558-af89-b72457aa8609', 'ef649d7d-49b9-550e-9726-f5c97ed52ee9', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'fae6f4d5-75c3-5558-af89-b72457aa8609')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'fae6f4d5-75c3-5558-af89-b72457aa8609' AND z.destino_etapa_id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9' AND z.condicao_transicao = 'Parecer favorável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '9b23a16d-93e9-52e6-8e04-1a158165ba92', 'ef649d7d-49b9-550e-9726-f5c97ed52ee9', '157b62c6-09c6-5a7b-8331-0ea3ecb2a601', 'Ato publicado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '157b62c6-09c6-5a7b-8331-0ea3ecb2a601')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ef649d7d-49b9-550e-9726-f5c97ed52ee9' AND z.destino_etapa_id = '157b62c6-09c6-5a7b-8331-0ea3ecb2a601' AND z.condicao_transicao = 'Ato publicado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c', 'RH.PAD.005', 'Processo Administrativo Disciplinar (PAD)', 'Apurar infração disciplinar de servidor com comissão, ampla defesa e contraditório, conforme o estatuto.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal) e estatuto dos servidores', ''), '2.0.05.02.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.05.02.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.PAD.005');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8', '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c', 1, 'GAB', 'Autoridade instauradora', 'Instaurar o PAD e designar a comissão.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f1b9d773-87aa-5f17-a2f3-64b43d4464fe', 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8', 'Portaria de instauração', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8' AND y.nome_documento = 'Portaria de instauração');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3f278597-566a-5069-963e-6f9d7b9ad957', '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c', 2, 'SEGEP/CPAD', 'Comissão processante', 'Citar o servidor, instruir, colher a defesa e elaborar o relatório final.', 60, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '80d72963-412f-5b90-b7a0-c935ff9d1a2e', '3f278597-566a-5069-963e-6f9d7b9ad957', 'Termo de citação', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND y.nome_documento = 'Termo de citação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7005dd49-1116-5a57-9d89-f4d5572f9fb1', '3f278597-566a-5069-963e-6f9d7b9ad957', 'Termos de depoimento', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND y.nome_documento = 'Termos de depoimento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a0874894-cb94-5e46-b954-6d718f55b60b', '3f278597-566a-5069-963e-6f9d7b9ad957', 'Defesa escrita', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND y.nome_documento = 'Defesa escrita');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '35619fa6-5bfa-5646-9f75-2314cd543e46', '3f278597-566a-5069-963e-6f9d7b9ad957', 'Relatório final da comissão', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND y.nome_documento = 'Relatório final da comissão');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44', '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c', 3, 'PGM', 'Procuradoria', 'Emitir parecer sobre a regularidade do processo.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e6e3d1e5-05df-5dc5-923e-d3ead96a3cb5', '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fb909d29-c255-5b3e-8e92-25ff89767a68', '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c', 4, 'GAB', 'Autoridade julgadora', 'Julgar e aplicar a penalidade ou arquivar.', 20, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5f028b0c-4ffb-5a66-9c5a-3f5ad546d35c' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b94f907f-da65-5cd8-8466-63cdde78acdb', 'fb909d29-c255-5b3e-8e92-25ff89767a68', 'Decisão de julgamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fb909d29-c255-5b3e-8e92-25ff89767a68')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fb909d29-c255-5b3e-8e92-25ff89767a68' AND y.nome_documento = 'Decisão de julgamento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '8cf789a4-31f6-5780-9417-ac11195e1ada', 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8', '3f278597-566a-5069-963e-6f9d7b9ad957', 'Comissão designada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a039b7c6-cc96-5e8d-9b8c-2a254b3265d8' AND z.destino_etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND z.condicao_transicao = 'Comissão designada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'eec22627-06e3-5c1d-8d57-e3865c9bfc11', '3f278597-566a-5069-963e-6f9d7b9ad957', '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44', 'Relatório concluído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND z.destino_etapa_id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44' AND z.condicao_transicao = 'Relatório concluído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'df998cd3-4655-5732-ac1a-0a60ecf0aa35', '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44', 'fb909d29-c255-5b3e-8e92-25ff89767a68', 'Processo regular', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fb909d29-c255-5b3e-8e92-25ff89767a68')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44' AND z.destino_etapa_id = 'fb909d29-c255-5b3e-8e92-25ff89767a68' AND z.condicao_transicao = 'Processo regular');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '59031004-d2e8-5b85-8788-aa0e4d8941c0', '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44', '3f278597-566a-5069-963e-6f9d7b9ad957', 'Nulidade ou vício', TRUE, NULLIF('Refazer o ato viciado', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3f278597-566a-5069-963e-6f9d7b9ad957')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5bbbac69-4ee2-57cc-a5a1-ad23f3129c44' AND z.destino_etapa_id = '3f278597-566a-5069-963e-6f9d7b9ad957' AND z.condicao_transicao = 'Nulidade ou vício');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '76a300a2-a878-598b-8716-2cd6ea452513', 'RH.EST.006', 'Avaliação de estágio probatório', 'Avaliar o servidor nomeado durante o estágio probatório e homologar a estabilidade.', 'Servidores em estágio probatório', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.08.01.02'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.08.01.02' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.EST.006');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5a84d936-efdf-53af-a02c-cee3fbbe4b7c', '76a300a2-a878-598b-8716-2cd6ea452513', 1, 'SECRETARIA', 'Chefia imediata', 'Preencher as avaliações periódicas com ciência do servidor.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '76a300a2-a878-598b-8716-2cd6ea452513')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '76a300a2-a878-598b-8716-2cd6ea452513' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '88f28cce-7d73-5bc4-9bc7-64ab66d7af5b', '5a84d936-efdf-53af-a02c-cee3fbbe4b7c', 'Avaliação de estágio probatório', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5a84d936-efdf-53af-a02c-cee3fbbe4b7c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5a84d936-efdf-53af-a02c-cee3fbbe4b7c' AND y.nome_documento = 'Avaliação de estágio probatório');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '95974235-fa89-512e-9f34-7949a5697ae2', '76a300a2-a878-598b-8716-2cd6ea452513', 2, 'SEGEP/CAEP', 'Comissão de avaliação', 'Consolidar as avaliações e emitir o parecer final.', 20, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '76a300a2-a878-598b-8716-2cd6ea452513')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '76a300a2-a878-598b-8716-2cd6ea452513' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6dd7535f-ba99-5d9f-ba56-0811d9879ed0', '95974235-fa89-512e-9f34-7949a5697ae2', 'Parecer da comissão', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '95974235-fa89-512e-9f34-7949a5697ae2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '95974235-fa89-512e-9f34-7949a5697ae2' AND y.nome_documento = 'Parecer da comissão');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd8d2cd94-fab2-5fa1-812a-9607f42259cc', '76a300a2-a878-598b-8716-2cd6ea452513', 3, 'SEGEP/RH', 'Recursos Humanos', 'Homologar o estágio probatório e registrar na pasta funcional.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '76a300a2-a878-598b-8716-2cd6ea452513')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '76a300a2-a878-598b-8716-2cd6ea452513' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c637286a-ca2d-55b5-a45d-c4cc5c5c62f3', 'd8d2cd94-fab2-5fa1-812a-9607f42259cc', 'Homologação do estágio probatório', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd8d2cd94-fab2-5fa1-812a-9607f42259cc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd8d2cd94-fab2-5fa1-812a-9607f42259cc' AND y.nome_documento = 'Homologação do estágio probatório');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '927162e7-4e49-579c-be0e-25c9240692e3', '5a84d936-efdf-53af-a02c-cee3fbbe4b7c', '95974235-fa89-512e-9f34-7949a5697ae2', 'Avaliações concluídas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5a84d936-efdf-53af-a02c-cee3fbbe4b7c')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '95974235-fa89-512e-9f34-7949a5697ae2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5a84d936-efdf-53af-a02c-cee3fbbe4b7c' AND z.destino_etapa_id = '95974235-fa89-512e-9f34-7949a5697ae2' AND z.condicao_transicao = 'Avaliações concluídas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'efb30028-983e-59aa-814d-589e4aaa9d8a', '95974235-fa89-512e-9f34-7949a5697ae2', 'd8d2cd94-fab2-5fa1-812a-9607f42259cc', 'Parecer emitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '95974235-fa89-512e-9f34-7949a5697ae2')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd8d2cd94-fab2-5fa1-812a-9607f42259cc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '95974235-fa89-512e-9f34-7949a5697ae2' AND z.destino_etapa_id = 'd8d2cd94-fab2-5fa1-812a-9607f42259cc' AND z.condicao_transicao = 'Parecer emitido');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '09bfb846-b584-5747-97c4-be9ce3538e53', '95974235-fa89-512e-9f34-7949a5697ae2', '5a84d936-efdf-53af-a02c-cee3fbbe4b7c', 'Avaliação incompleta', TRUE, NULLIF('Completar as avaliações periódicas', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '95974235-fa89-512e-9f34-7949a5697ae2')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5a84d936-efdf-53af-a02c-cee3fbbe4b7c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '95974235-fa89-512e-9f34-7949a5697ae2' AND z.destino_etapa_id = '5a84d936-efdf-53af-a02c-cee3fbbe4b7c' AND z.condicao_transicao = 'Avaliação incompleta');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'd2a4af19-c4d5-57d1-8931-625a897125e2', 'RH.CON.007', 'Posse de aprovado em concurso público', 'Convocar o aprovado, conferir os documentos e a aptidão médica e dar posse no cargo.', 'Aprovados em concurso público', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.08.00.14'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.08.00.14' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.CON.007');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '10842858-cbc6-5158-aa50-2ced30547d7d', 'd2a4af19-c4d5-57d1-8931-625a897125e2', 1, 'SEGEP/RH', 'Recursos Humanos', 'Publicar a convocação e receber a documentação para a posse.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd2a4af19-c4d5-57d1-8931-625a897125e2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd2a4af19-c4d5-57d1-8931-625a897125e2' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7ebd1df2-c170-5ed8-9552-e2263c3279ab', '10842858-cbc6-5158-aa50-2ced30547d7d', 'Edital de convocação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '10842858-cbc6-5158-aa50-2ced30547d7d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '10842858-cbc6-5158-aa50-2ced30547d7d' AND y.nome_documento = 'Edital de convocação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ccf43975-d8ef-51db-bcac-f9567afb4e0a', '10842858-cbc6-5158-aa50-2ced30547d7d', 'Documentação pessoal', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '10842858-cbc6-5158-aa50-2ced30547d7d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '10842858-cbc6-5158-aa50-2ced30547d7d' AND y.nome_documento = 'Documentação pessoal');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '459d0867-9f2d-5ea2-b3d7-0b2dcb63e9dd', '10842858-cbc6-5158-aa50-2ced30547d7d', 'Declaração de acúmulo de cargos', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '10842858-cbc6-5158-aa50-2ced30547d7d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '10842858-cbc6-5158-aa50-2ced30547d7d' AND y.nome_documento = 'Declaração de acúmulo de cargos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '749352a0-4a15-54bf-9571-b8e22cfe28f6', 'd2a4af19-c4d5-57d1-8931-625a897125e2', 2, 'SEGEP/DESOPEM', 'Saúde ocupacional (DESOPEM)', 'Realizar a perícia de admissão e emitir o laudo de aptidão.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd2a4af19-c4d5-57d1-8931-625a897125e2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd2a4af19-c4d5-57d1-8931-625a897125e2' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '57d81b6e-3404-574f-9166-af7caa6357a8', '749352a0-4a15-54bf-9571-b8e22cfe28f6', 'Laudo de aptidão física e mental', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '749352a0-4a15-54bf-9571-b8e22cfe28f6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '749352a0-4a15-54bf-9571-b8e22cfe28f6' AND y.nome_documento = 'Laudo de aptidão física e mental');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '87a479ac-b609-5341-aa99-f8ae6c617677', 'd2a4af19-c4d5-57d1-8931-625a897125e2', 3, 'SEGEP/RH', 'Recursos Humanos', 'Lavrar o termo de posse e cadastrar o servidor.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'd2a4af19-c4d5-57d1-8931-625a897125e2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'd2a4af19-c4d5-57d1-8931-625a897125e2' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2dd9d3a9-b803-5ff9-a6f3-137980b45307', '87a479ac-b609-5341-aa99-f8ae6c617677', 'Portaria de nomeação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '87a479ac-b609-5341-aa99-f8ae6c617677')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '87a479ac-b609-5341-aa99-f8ae6c617677' AND y.nome_documento = 'Portaria de nomeação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '50a44f70-f83b-5182-90f0-f9fd5bf023f3', '87a479ac-b609-5341-aa99-f8ae6c617677', 'Termo de posse', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '87a479ac-b609-5341-aa99-f8ae6c617677')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '87a479ac-b609-5341-aa99-f8ae6c617677' AND y.nome_documento = 'Termo de posse');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4ac2279f-ab37-5ea9-a0d6-b601f12f718b', '10842858-cbc6-5158-aa50-2ced30547d7d', '749352a0-4a15-54bf-9571-b8e22cfe28f6', 'Documentação completa', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '10842858-cbc6-5158-aa50-2ced30547d7d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '749352a0-4a15-54bf-9571-b8e22cfe28f6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '10842858-cbc6-5158-aa50-2ced30547d7d' AND z.destino_etapa_id = '749352a0-4a15-54bf-9571-b8e22cfe28f6' AND z.condicao_transicao = 'Documentação completa');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c520beeb-708d-5669-887a-7f298b89f249', '749352a0-4a15-54bf-9571-b8e22cfe28f6', '87a479ac-b609-5341-aa99-f8ae6c617677', 'Apto', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '749352a0-4a15-54bf-9571-b8e22cfe28f6')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '87a479ac-b609-5341-aa99-f8ae6c617677')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '749352a0-4a15-54bf-9571-b8e22cfe28f6' AND z.destino_etapa_id = '87a479ac-b609-5341-aa99-f8ae6c617677' AND z.condicao_transicao = 'Apto');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2d3680ec-285a-58a7-8938-4f2f60a1211f', '749352a0-4a15-54bf-9571-b8e22cfe28f6', '10842858-cbc6-5158-aa50-2ced30547d7d', 'Exame complementar', TRUE, NULLIF('Apresentar o exame solicitado', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '749352a0-4a15-54bf-9571-b8e22cfe28f6')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '10842858-cbc6-5158-aa50-2ced30547d7d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '749352a0-4a15-54bf-9571-b8e22cfe28f6' AND z.destino_etapa_id = '10842858-cbc6-5158-aa50-2ced30547d7d' AND z.condicao_transicao = 'Exame complementar');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '7a5343ed-b882-5458-bc9b-f6134ddb1913', 'FAZ.DIA.001', 'Concessão e prestação de contas de diárias', 'Conceder diárias para viagem a serviço e receber a prestação de contas no retorno.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.01.00.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.01.00.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.DIA.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '129aea06-7141-51d2-b8af-a3f0178ab0c9', '7a5343ed-b882-5458-bc9b-f6134ddb1913', 1, 'SECRETARIA', 'Chefia imediata', 'Solicitar as diárias com a justificativa e o roteiro da viagem.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a5343ed-b882-5458-bc9b-f6134ddb1913')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a5343ed-b882-5458-bc9b-f6134ddb1913' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'df0b1b95-4524-55ac-913c-c6489b1ec219', '129aea06-7141-51d2-b8af-a3f0178ab0c9', 'Solicitação de diárias', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '129aea06-7141-51d2-b8af-a3f0178ab0c9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '129aea06-7141-51d2-b8af-a3f0178ab0c9' AND y.nome_documento = 'Solicitação de diárias');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b17ece61-425c-53d3-860c-64fe0fa28d2c', '7a5343ed-b882-5458-bc9b-f6134ddb1913', 2, 'SEFIN/CONTABILIDADE', 'Contabilidade', 'Empenhar e pagar as diárias.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a5343ed-b882-5458-bc9b-f6134ddb1913')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a5343ed-b882-5458-bc9b-f6134ddb1913' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '015fd4c4-274d-575a-ad28-a15e03b0cfb1', 'b17ece61-425c-53d3-860c-64fe0fa28d2c', 'Nota de empenho', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c' AND y.nome_documento = 'Nota de empenho');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '42fecdb9-a8e4-57d9-a6dc-0ea3772b1ae1', 'b17ece61-425c-53d3-860c-64fe0fa28d2c', 'Ordem de pagamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c' AND y.nome_documento = 'Ordem de pagamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '412a4972-2fbc-5254-8b9c-aa401f20cb46', '7a5343ed-b882-5458-bc9b-f6134ddb1913', 3, 'SECRETARIA', 'Servidor', 'Prestar contas em até 5 dias do retorno, com os comprovantes.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a5343ed-b882-5458-bc9b-f6134ddb1913')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a5343ed-b882-5458-bc9b-f6134ddb1913' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a4d988d0-5f45-5085-9bd1-39aeafd17ad9', '412a4972-2fbc-5254-8b9c-aa401f20cb46', 'Relatório de viagem', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '412a4972-2fbc-5254-8b9c-aa401f20cb46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '412a4972-2fbc-5254-8b9c-aa401f20cb46' AND y.nome_documento = 'Relatório de viagem');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a46b5a13-0e47-5ea8-ac45-58bceec23721', '412a4972-2fbc-5254-8b9c-aa401f20cb46', 'Comprovantes de deslocamento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '412a4972-2fbc-5254-8b9c-aa401f20cb46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '412a4972-2fbc-5254-8b9c-aa401f20cb46' AND y.nome_documento = 'Comprovantes de deslocamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'f504b307-877b-5287-9a85-da560c0e3e24', '7a5343ed-b882-5458-bc9b-f6134ddb1913', 4, 'SEFIN/CONTABILIDADE', 'Contabilidade', 'Analisar a prestação de contas e aprovar ou pedir a devolução.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a5343ed-b882-5458-bc9b-f6134ddb1913')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a5343ed-b882-5458-bc9b-f6134ddb1913' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '48fce05d-2bb1-5e08-8a61-8ee23c19076a', 'f504b307-877b-5287-9a85-da560c0e3e24', 'Parecer sobre a prestação de contas', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'f504b307-877b-5287-9a85-da560c0e3e24')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'f504b307-877b-5287-9a85-da560c0e3e24' AND y.nome_documento = 'Parecer sobre a prestação de contas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4049ef89-817d-51f0-9d67-ecf2271f319e', '129aea06-7141-51d2-b8af-a3f0178ab0c9', 'b17ece61-425c-53d3-860c-64fe0fa28d2c', 'Viagem autorizada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '129aea06-7141-51d2-b8af-a3f0178ab0c9')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '129aea06-7141-51d2-b8af-a3f0178ab0c9' AND z.destino_etapa_id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c' AND z.condicao_transicao = 'Viagem autorizada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '1913a052-5556-518a-8ec5-437df379dbd6', 'b17ece61-425c-53d3-860c-64fe0fa28d2c', '412a4972-2fbc-5254-8b9c-aa401f20cb46', 'Diárias pagas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '412a4972-2fbc-5254-8b9c-aa401f20cb46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'b17ece61-425c-53d3-860c-64fe0fa28d2c' AND z.destino_etapa_id = '412a4972-2fbc-5254-8b9c-aa401f20cb46' AND z.condicao_transicao = 'Diárias pagas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0863c013-b85b-54a2-b127-198aca8e9202', '412a4972-2fbc-5254-8b9c-aa401f20cb46', 'f504b307-877b-5287-9a85-da560c0e3e24', 'Contas prestadas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '412a4972-2fbc-5254-8b9c-aa401f20cb46')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'f504b307-877b-5287-9a85-da560c0e3e24')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '412a4972-2fbc-5254-8b9c-aa401f20cb46' AND z.destino_etapa_id = 'f504b307-877b-5287-9a85-da560c0e3e24' AND z.condicao_transicao = 'Contas prestadas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '1801d17d-7c2a-5690-9e72-ccaf5345c96f', 'f504b307-877b-5287-9a85-da560c0e3e24', '412a4972-2fbc-5254-8b9c-aa401f20cb46', 'Comprovação insuficiente', TRUE, NULLIF('Complementar a prestação de contas ou devolver o valor', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'f504b307-877b-5287-9a85-da560c0e3e24')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '412a4972-2fbc-5254-8b9c-aa401f20cb46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'f504b307-877b-5287-9a85-da560c0e3e24' AND z.destino_etapa_id = '412a4972-2fbc-5254-8b9c-aa401f20cb46' AND z.condicao_transicao = 'Comprovação insuficiente');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135', 'FAZ.DES.002', 'Liquidação e pagamento de despesa', 'Liquidar a despesa empenhada a partir do atesto e pagar o credor, na ordem cronológica.', 'Fornecedores e secretarias', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.01.00.02'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.01.00.02' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.DES.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a3536d7e-682c-556a-bbbe-4fa435d38e58', '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135', 1, 'SECRETARIA', 'Gestor do contrato', 'Atestar o recebimento do bem ou serviço e encaminhar a nota fiscal.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '94029dfa-4334-5da3-8fbf-a0589b096fc2', 'a3536d7e-682c-556a-bbbe-4fa435d38e58', 'Nota fiscal atestada', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58' AND y.nome_documento = 'Nota fiscal atestada');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2a1ca0a6-8264-544a-89ca-25b4bc2474ed', 'a3536d7e-682c-556a-bbbe-4fa435d38e58', 'Termo de recebimento', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58' AND y.nome_documento = 'Termo de recebimento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb', '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135', 2, 'SEFIN/CONTABILIDADE', 'Contabilidade', 'Conferir o empenho, as retenções e a regularidade fiscal e liquidar.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6643a40e-da91-5f64-8a0b-dbc7f601cd50', '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb', 'Nota de liquidação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb' AND y.nome_documento = 'Nota de liquidação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4ca38155-f03a-56f0-a770-1b2816f858cb', '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb', 'Certidões de regularidade', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb' AND y.nome_documento = 'Certidões de regularidade');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'cc14e947-593f-5e2a-b282-cf5c25d54192', '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135', 3, 'SEFIN/TESOURARIA', 'Tesouraria', 'Pagar o credor obedecendo a ordem cronológica.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1a28c3dc-f9dd-5b48-b756-cf9f92c5d135' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6ec17931-91b1-50fc-8325-5656845c3c2b', 'cc14e947-593f-5e2a-b282-cf5c25d54192', 'Ordem de pagamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cc14e947-593f-5e2a-b282-cf5c25d54192')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cc14e947-593f-5e2a-b282-cf5c25d54192' AND y.nome_documento = 'Ordem de pagamento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7404ea40-29a2-5ca4-b733-55c585518328', 'cc14e947-593f-5e2a-b282-cf5c25d54192', 'Comprovante bancário', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cc14e947-593f-5e2a-b282-cf5c25d54192')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cc14e947-593f-5e2a-b282-cf5c25d54192' AND y.nome_documento = 'Comprovante bancário');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '44d6a932-70e4-5282-9b12-90befa1fd9fd', 'a3536d7e-682c-556a-bbbe-4fa435d38e58', '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb', 'Nota atestada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58' AND z.destino_etapa_id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb' AND z.condicao_transicao = 'Nota atestada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'd87ee18b-68b5-5eb7-92e3-eb5caf817d49', '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb', 'cc14e947-593f-5e2a-b282-cf5c25d54192', 'Despesa liquidada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'cc14e947-593f-5e2a-b282-cf5c25d54192')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb' AND z.destino_etapa_id = 'cc14e947-593f-5e2a-b282-cf5c25d54192' AND z.condicao_transicao = 'Despesa liquidada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'a1a09c73-3bea-5402-941f-96e85529b348', '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb', 'a3536d7e-682c-556a-bbbe-4fa435d38e58', 'Nota com erro ou certidão vencida', TRUE, NULLIF('Corrigir a nota ou atualizar as certidões', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5dc12dbf-64b0-5605-8c8f-dc7e358f55eb' AND z.destino_etapa_id = 'a3536d7e-682c-556a-bbbe-4fa435d38e58' AND z.condicao_transicao = 'Nota com erro ou certidão vencida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1e5b3025-067e-55fa-a9d3-c87e71fc483c', 'FAZ.PAR.003', 'Parcelamento de débitos municipais', 'Firmar acordo de parcelamento de tributos e da dívida ativa, com confissão de dívida.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.02.00.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.02.00.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.PAR.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc', '1e5b3025-067e-55fa-a9d3-c87e71fc483c', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber o requerimento e emitir o extrato dos débitos.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1e5b3025-067e-55fa-a9d3-c87e71fc483c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1e5b3025-067e-55fa-a9d3-c87e71fc483c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '0778dfb4-a27d-5602-8e3b-96f6c34a7981', 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc', 'Requerimento de parcelamento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc' AND y.nome_documento = 'Requerimento de parcelamento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9926e70e-aed4-53c3-b8b5-c5bdd3b1b895', 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc', 'Extrato de débitos', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc' AND y.nome_documento = 'Extrato de débitos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ea354b6f-d163-5df2-b688-9d3f8e3d6821', '1e5b3025-067e-55fa-a9d3-c87e71fc483c', 2, 'SEREC/DIVIDA', 'Dívida ativa', 'Simular as parcelas conforme a lei e lavrar o termo de confissão.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1e5b3025-067e-55fa-a9d3-c87e71fc483c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1e5b3025-067e-55fa-a9d3-c87e71fc483c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3316f0d0-fc11-5c66-85ca-165ac8d216f3', 'ea354b6f-d163-5df2-b688-9d3f8e3d6821', 'Termo de confissão de dívida', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ea354b6f-d163-5df2-b688-9d3f8e3d6821')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ea354b6f-d163-5df2-b688-9d3f8e3d6821' AND y.nome_documento = 'Termo de confissão de dívida');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3385bdfb-8986-5b4c-8cac-783e0e290908', '1e5b3025-067e-55fa-a9d3-c87e71fc483c', 3, 'SEREC/ARRECADACAO', 'Arrecadação', 'Emitir as guias e acompanhar o pagamento.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1e5b3025-067e-55fa-a9d3-c87e71fc483c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1e5b3025-067e-55fa-a9d3-c87e71fc483c' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '51022086-a09a-51a7-b68c-d2416324ef86', '3385bdfb-8986-5b4c-8cac-783e0e290908', 'Guias de pagamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3385bdfb-8986-5b4c-8cac-783e0e290908')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3385bdfb-8986-5b4c-8cac-783e0e290908' AND y.nome_documento = 'Guias de pagamento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '40918170-c3b0-5d77-a725-d9044d37326f', 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc', 'ea354b6f-d163-5df2-b688-9d3f8e3d6821', 'Débitos apurados', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ea354b6f-d163-5df2-b688-9d3f8e3d6821')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'afdd2a34-82b6-5a07-93e5-2cd37634c1fc' AND z.destino_etapa_id = 'ea354b6f-d163-5df2-b688-9d3f8e3d6821' AND z.condicao_transicao = 'Débitos apurados');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'aeab188c-7577-54d9-b34a-9e6fdb646a4c', 'ea354b6f-d163-5df2-b688-9d3f8e3d6821', '3385bdfb-8986-5b4c-8cac-783e0e290908', 'Termo assinado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ea354b6f-d163-5df2-b688-9d3f8e3d6821')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3385bdfb-8986-5b4c-8cac-783e0e290908')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ea354b6f-d163-5df2-b688-9d3f8e3d6821' AND z.destino_etapa_id = '3385bdfb-8986-5b4c-8cac-783e0e290908' AND z.condicao_transicao = 'Termo assinado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1137cdbd-f34b-557f-ac7d-1583574ef897', 'FAZ.IPT.004', 'Isenção de IPTU', 'Reconhecer a isenção de IPTU prevista em lei (ex.: aposentado, imóvel único) a pedido do contribuinte.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.02.01.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.02.01.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.IPT.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5fa71df3-53a2-50f2-9bbf-98bf67489084', '1137cdbd-f34b-557f-ac7d-1583574ef897', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber o requerimento com os comprovantes.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1137cdbd-f34b-557f-ac7d-1583574ef897')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1137cdbd-f34b-557f-ac7d-1583574ef897' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bec1fdd6-96a8-5068-b788-aa3f7bf56dda', '5fa71df3-53a2-50f2-9bbf-98bf67489084', 'Requerimento de isenção', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5fa71df3-53a2-50f2-9bbf-98bf67489084')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5fa71df3-53a2-50f2-9bbf-98bf67489084' AND y.nome_documento = 'Requerimento de isenção');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5d67b644-78f2-5316-b02c-6a4070021b27', '5fa71df3-53a2-50f2-9bbf-98bf67489084', 'Comprovantes do requisito legal', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5fa71df3-53a2-50f2-9bbf-98bf67489084')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5fa71df3-53a2-50f2-9bbf-98bf67489084' AND y.nome_documento = 'Comprovantes do requisito legal');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ed789eb6-89f5-5a62-be64-7a717ed8502c', '5fa71df3-53a2-50f2-9bbf-98bf67489084', 'Matrícula ou documento do imóvel', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5fa71df3-53a2-50f2-9bbf-98bf67489084')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5fa71df3-53a2-50f2-9bbf-98bf67489084' AND y.nome_documento = 'Matrícula ou documento do imóvel');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e5929a73-660f-5ea5-a4c7-7fba507bb362', '1137cdbd-f34b-557f-ac7d-1583574ef897', 2, 'SEREC/CADASTRO', 'Cadastro imobiliário', 'Conferir o cadastro do imóvel e os requisitos.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1137cdbd-f34b-557f-ac7d-1583574ef897')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1137cdbd-f34b-557f-ac7d-1583574ef897' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '83e5c5ec-ce90-55ac-8f79-f34bc1eafcbe', 'e5929a73-660f-5ea5-a4c7-7fba507bb362', 'Informação do cadastro imobiliário', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362' AND y.nome_documento = 'Informação do cadastro imobiliário');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf', '1137cdbd-f34b-557f-ac7d-1583574ef897', 3, 'SEREC/TRIBUTARIO', 'Julgamento tributário', 'Decidir o pedido e lançar a isenção no cadastro.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1137cdbd-f34b-557f-ac7d-1583574ef897')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1137cdbd-f34b-557f-ac7d-1583574ef897' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cf03130d-f994-5b4a-bc43-5c239fefe5c4', 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf', 'Decisão administrativa', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf' AND y.nome_documento = 'Decisão administrativa');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '139d729f-95d1-5c38-835e-171437fba4c8', '5fa71df3-53a2-50f2-9bbf-98bf67489084', 'e5929a73-660f-5ea5-a4c7-7fba507bb362', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5fa71df3-53a2-50f2-9bbf-98bf67489084')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5fa71df3-53a2-50f2-9bbf-98bf67489084' AND z.destino_etapa_id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362' AND z.condicao_transicao = 'Requerimento instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'aabc0bdf-e197-5aa9-8a49-001c7ef8f793', 'e5929a73-660f-5ea5-a4c7-7fba507bb362', 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf', 'Requisitos conferidos', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362' AND z.destino_etapa_id = 'b27c0a4f-d5fd-5f50-bdd6-39aef88b03cf' AND z.condicao_transicao = 'Requisitos conferidos');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c5c14d65-83a2-579e-bf17-f38f3e74b1ea', 'e5929a73-660f-5ea5-a4c7-7fba507bb362', '5fa71df3-53a2-50f2-9bbf-98bf67489084', 'Documento faltante', TRUE, NULLIF('Pedir o documento ao contribuinte', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5fa71df3-53a2-50f2-9bbf-98bf67489084')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e5929a73-660f-5ea5-a4c7-7fba507bb362' AND z.destino_etapa_id = '5fa71df3-53a2-50f2-9bbf-98bf67489084' AND z.condicao_transicao = 'Documento faltante');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'a89cdf93-82f8-54fc-9bf9-848061b52f89', 'FAZ.ALV.005', 'Alvará de funcionamento', 'Licenciar o funcionamento de atividade econômica, com consulta de uso do solo e vistorias exigidas.', 'Empresas e profissionais autônomos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.03.00.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.03.00.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.ALV.005');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97', 'a89cdf93-82f8-54fc-9bf9-848061b52f89', 1, 'SEREC/LICENCIAMENTO', 'Licenciamento econômico', 'Receber o pedido, conferir a inscrição municipal e a consulta de uso do solo.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a89cdf93-82f8-54fc-9bf9-848061b52f89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a89cdf93-82f8-54fc-9bf9-848061b52f89' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '52a17b68-004a-543c-a7f7-c6d103759228', 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97', 'Requerimento de alvará', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97' AND y.nome_documento = 'Requerimento de alvará');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c875d839-36dd-57d4-af2f-b43c78295977', 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97', 'Contrato social ou CNPJ', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97' AND y.nome_documento = 'Contrato social ou CNPJ');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2555fcc6-e8a2-5e1c-9b56-64b044624210', 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97', 'Certidão de uso do solo', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97' AND y.nome_documento = 'Certidão de uso do solo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9ae5ac35-6f7d-5bbd-b363-4593b15906e3', 'a89cdf93-82f8-54fc-9bf9-848061b52f89', 2, 'SEREC/FISCALIZACAO', 'Fiscalização', 'Vistoriar o estabelecimento quando a atividade exigir.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a89cdf93-82f8-54fc-9bf9-848061b52f89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a89cdf93-82f8-54fc-9bf9-848061b52f89' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '131fe2e4-d257-5a46-b828-96a9531d444a', '9ae5ac35-6f7d-5bbd-b363-4593b15906e3', 'Relatório de vistoria', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3' AND y.nome_documento = 'Relatório de vistoria');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '75823ab8-2566-52d3-95bb-2f1175e93a20', 'a89cdf93-82f8-54fc-9bf9-848061b52f89', 3, 'SEREC/LICENCIAMENTO', 'Licenciamento econômico', 'Emitir o alvará após o pagamento da taxa.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a89cdf93-82f8-54fc-9bf9-848061b52f89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a89cdf93-82f8-54fc-9bf9-848061b52f89' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '1463a3a0-bfe5-5db4-9e1f-6123d2f1d55e', '75823ab8-2566-52d3-95bb-2f1175e93a20', 'Alvará de funcionamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '75823ab8-2566-52d3-95bb-2f1175e93a20')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '75823ab8-2566-52d3-95bb-2f1175e93a20' AND y.nome_documento = 'Alvará de funcionamento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b9ecc16d-104c-5e61-9743-098a53ec5f7c', 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97', '9ae5ac35-6f7d-5bbd-b363-4593b15906e3', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97' AND z.destino_etapa_id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'af9e4a95-9b95-5b6a-be43-c11003ba3756', '9ae5ac35-6f7d-5bbd-b363-4593b15906e3', '75823ab8-2566-52d3-95bb-2f1175e93a20', 'Vistoria favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '75823ab8-2566-52d3-95bb-2f1175e93a20')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3' AND z.destino_etapa_id = '75823ab8-2566-52d3-95bb-2f1175e93a20' AND z.condicao_transicao = 'Vistoria favorável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'ce2c969f-8c10-55bb-9574-552f996c16af', '9ae5ac35-6f7d-5bbd-b363-4593b15906e3', 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97', 'Irregularidade na vistoria', TRUE, NULLIF('Sanar a irregularidade apontada', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '9ae5ac35-6f7d-5bbd-b363-4593b15906e3' AND z.destino_etapa_id = 'dd400713-8ed7-5b4f-8e01-aa85e0a5ab97' AND z.condicao_transicao = 'Irregularidade na vistoria');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '191ecab2-7d4e-5cf4-863f-f1b187c3e769', 'FAZ.CND.006', 'Certidão negativa de débitos municipais', 'Emitir certidão negativa (ou positiva com efeito de negativa) de débitos do contribuinte.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.02.00.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.02.00.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.CND.006');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '80a1fa99-a462-5242-afbb-88b5491f84da', '191ecab2-7d4e-5cf4-863f-f1b187c3e769', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber o pedido e consultar a situação fiscal.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '191ecab2-7d4e-5cf4-863f-f1b187c3e769')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '191ecab2-7d4e-5cf4-863f-f1b187c3e769' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '02f677b4-c0e1-54d0-9577-7dd6fd2091cc', '80a1fa99-a462-5242-afbb-88b5491f84da', 'Requerimento de certidão', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '80a1fa99-a462-5242-afbb-88b5491f84da')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '80a1fa99-a462-5242-afbb-88b5491f84da' AND y.nome_documento = 'Requerimento de certidão');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '359ffdd9-3c8a-5f46-a887-8703e3fbf57f', '191ecab2-7d4e-5cf4-863f-f1b187c3e769', 2, 'SEREC/DIVIDA', 'Dívida ativa', 'Analisar pendências apontadas (parcelamento, suspensão de exigibilidade).', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '191ecab2-7d4e-5cf4-863f-f1b187c3e769')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '191ecab2-7d4e-5cf4-863f-f1b187c3e769' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3511e882-7031-5fcb-942c-604e1e1a8c4c', '359ffdd9-3c8a-5f46-a887-8703e3fbf57f', 'Informação sobre pendências', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '359ffdd9-3c8a-5f46-a887-8703e3fbf57f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '359ffdd9-3c8a-5f46-a887-8703e3fbf57f' AND y.nome_documento = 'Informação sobre pendências');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '7738c813-2636-52b3-a84e-b86eb2c81e3e', '191ecab2-7d4e-5cf4-863f-f1b187c3e769', 3, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Emitir a certidão.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '191ecab2-7d4e-5cf4-863f-f1b187c3e769')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '191ecab2-7d4e-5cf4-863f-f1b187c3e769' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '25363764-0bfe-5957-a889-2576291e1d07', '7738c813-2636-52b3-a84e-b86eb2c81e3e', 'Certidão negativa de débitos', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7738c813-2636-52b3-a84e-b86eb2c81e3e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7738c813-2636-52b3-a84e-b86eb2c81e3e' AND y.nome_documento = 'Certidão negativa de débitos');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0da23202-a132-53d7-b6a0-3cf7bab24c8c', '80a1fa99-a462-5242-afbb-88b5491f84da', '359ffdd9-3c8a-5f46-a887-8703e3fbf57f', 'Sem pendências', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '80a1fa99-a462-5242-afbb-88b5491f84da')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '359ffdd9-3c8a-5f46-a887-8703e3fbf57f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '80a1fa99-a462-5242-afbb-88b5491f84da' AND z.destino_etapa_id = '359ffdd9-3c8a-5f46-a887-8703e3fbf57f' AND z.condicao_transicao = 'Sem pendências');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '833a625a-097e-529a-8540-77ba9e0fea87', '359ffdd9-3c8a-5f46-a887-8703e3fbf57f', '7738c813-2636-52b3-a84e-b86eb2c81e3e', 'Situação definida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '359ffdd9-3c8a-5f46-a887-8703e3fbf57f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '7738c813-2636-52b3-a84e-b86eb2c81e3e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '359ffdd9-3c8a-5f46-a887-8703e3fbf57f' AND z.destino_etapa_id = '7738c813-2636-52b3-a84e-b86eb2c81e3e' AND z.condicao_transicao = 'Situação definida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '9a674a22-d4dd-5e4d-8445-007dfb133ae9', 'CTR.OUV.001', 'Manifestação de ouvidoria', 'Receber, encaminhar e responder reclamação, denúncia, sugestão ou elogio do cidadão (Lei 13.460/2017).', 'Cidadãos', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('Lei 13.460/2017, art. 10, §7º (identidade do manifestante)', ''), '4.0.01.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '4.0.01.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'CTR.OUV.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4a537a66-06f1-5985-87d9-dcebccdc43a0', '9a674a22-d4dd-5e4d-8445-007dfb133ae9', 1, 'UCCI/OUVIDORIA', 'Ouvidoria-Geral', 'Registrar a manifestação, preservar a identidade do manifestante e encaminhar.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9a674a22-d4dd-5e4d-8445-007dfb133ae9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9a674a22-d4dd-5e4d-8445-007dfb133ae9' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cc1b1172-6430-59c5-86c6-39d84162b165', '4a537a66-06f1-5985-87d9-dcebccdc43a0', 'Formulário de manifestação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4a537a66-06f1-5985-87d9-dcebccdc43a0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4a537a66-06f1-5985-87d9-dcebccdc43a0' AND y.nome_documento = 'Formulário de manifestação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '1a6cb147-f4d1-5065-8bd0-396b8f08f0db', '9a674a22-d4dd-5e4d-8445-007dfb133ae9', 2, 'SECRETARIA', 'Unidade responsável', 'Apurar e responder à ouvidoria.', 20, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9a674a22-d4dd-5e4d-8445-007dfb133ae9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9a674a22-d4dd-5e4d-8445-007dfb133ae9' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '55ea674d-a55a-5b93-9bbf-4a8aceef5507', '1a6cb147-f4d1-5065-8bd0-396b8f08f0db', 'Resposta da unidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db' AND y.nome_documento = 'Resposta da unidade');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'f0413bff-6b65-5fa3-a585-a871d5013e33', '9a674a22-d4dd-5e4d-8445-007dfb133ae9', 3, 'UCCI/OUVIDORIA', 'Ouvidoria-Geral', 'Avaliar a resposta, responder ao manifestante e encerrar.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9a674a22-d4dd-5e4d-8445-007dfb133ae9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9a674a22-d4dd-5e4d-8445-007dfb133ae9' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '88e86141-7b9a-587c-8daa-8b679d4fe304', 'f0413bff-6b65-5fa3-a585-a871d5013e33', 'Resposta ao manifestante', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'f0413bff-6b65-5fa3-a585-a871d5013e33')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'f0413bff-6b65-5fa3-a585-a871d5013e33' AND y.nome_documento = 'Resposta ao manifestante');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '6857acbf-c583-5c39-8784-f73bf5d27115', '4a537a66-06f1-5985-87d9-dcebccdc43a0', '1a6cb147-f4d1-5065-8bd0-396b8f08f0db', 'Manifestação encaminhada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4a537a66-06f1-5985-87d9-dcebccdc43a0')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4a537a66-06f1-5985-87d9-dcebccdc43a0' AND z.destino_etapa_id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db' AND z.condicao_transicao = 'Manifestação encaminhada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4d3fa976-d096-5dba-baa8-9d5a695b17fe', '1a6cb147-f4d1-5065-8bd0-396b8f08f0db', 'f0413bff-6b65-5fa3-a585-a871d5013e33', 'Resposta recebida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'f0413bff-6b65-5fa3-a585-a871d5013e33')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db' AND z.destino_etapa_id = 'f0413bff-6b65-5fa3-a585-a871d5013e33' AND z.condicao_transicao = 'Resposta recebida');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b0bb4f32-5cc4-5808-bc66-ce6fa78193a4', 'f0413bff-6b65-5fa3-a585-a871d5013e33', '1a6cb147-f4d1-5065-8bd0-396b8f08f0db', 'Resposta incompleta', TRUE, NULLIF('Complementar a resposta', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'f0413bff-6b65-5fa3-a585-a871d5013e33')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'f0413bff-6b65-5fa3-a585-a871d5013e33' AND z.destino_etapa_id = '1a6cb147-f4d1-5065-8bd0-396b8f08f0db' AND z.condicao_transicao = 'Resposta incompleta');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75', 'HAB.ALV.001', 'Alvará de licença para construção', 'Aprovar o projeto e licenciar a obra conforme o código de obras e o plano diretor.', 'Proprietários e responsáveis técnicos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '6.0.02.00.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.02.00.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.ALV.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75', 1, 'SEHAB/PROTOCOLO', 'Atendimento', 'Receber o pedido com o projeto e a responsabilidade técnica.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3362a4d1-f4b3-5b3b-95f3-a97ee946121a', 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'Requerimento de alvará de construção', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4' AND y.nome_documento = 'Requerimento de alvará de construção');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c7bb9401-fa16-511f-abae-dc43f61c245f', 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'Projeto arquitetônico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4' AND y.nome_documento = 'Projeto arquitetônico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd1ed784a-9b2b-5886-b0c1-8593ed2a1533', 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'ART/RRT do responsável técnico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4' AND y.nome_documento = 'ART/RRT do responsável técnico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'db124db7-3f6a-5b87-a61e-816596bd27ae', 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'Matrícula do imóvel', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4' AND y.nome_documento = 'Matrícula do imóvel');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c4467627-4225-52c2-ac7b-85a99ae3389a', 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75', 2, 'SEHAB/APROVACAO', 'Aprovação de projetos', 'Analisar o projeto quanto às normas urbanísticas.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f98fa4df-f608-5ea5-93d0-b8d61ef62366', 'c4467627-4225-52c2-ac7b-85a99ae3389a', 'Parecer técnico de análise', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c4467627-4225-52c2-ac7b-85a99ae3389a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c4467627-4225-52c2-ac7b-85a99ae3389a' AND y.nome_documento = 'Parecer técnico de análise');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9cb236bd-e2d2-571a-b2be-5243a9535fc7', 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75', 3, 'SEHAB/APROVACAO', 'Aprovação de projetos', 'Emitir o alvará após o pagamento das taxas.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ff6c144f-01cf-5b5c-b1d1-982c63b62e75' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a3014e97-c2bb-5287-a218-f9fd301dac0b', '9cb236bd-e2d2-571a-b2be-5243a9535fc7', 'Alvará de licença para construção', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9cb236bd-e2d2-571a-b2be-5243a9535fc7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9cb236bd-e2d2-571a-b2be-5243a9535fc7' AND y.nome_documento = 'Alvará de licença para construção');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'a3b32508-1795-50d2-be21-31bd78767c55', 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'c4467627-4225-52c2-ac7b-85a99ae3389a', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c4467627-4225-52c2-ac7b-85a99ae3389a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4' AND z.destino_etapa_id = 'c4467627-4225-52c2-ac7b-85a99ae3389a' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'da236171-aa6c-5184-8590-765834d880a7', 'c4467627-4225-52c2-ac7b-85a99ae3389a', '9cb236bd-e2d2-571a-b2be-5243a9535fc7', 'Projeto aprovado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c4467627-4225-52c2-ac7b-85a99ae3389a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '9cb236bd-e2d2-571a-b2be-5243a9535fc7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c4467627-4225-52c2-ac7b-85a99ae3389a' AND z.destino_etapa_id = '9cb236bd-e2d2-571a-b2be-5243a9535fc7' AND z.condicao_transicao = 'Projeto aprovado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'af1fab41-dfb3-5096-be85-6b7f86bb668f', 'c4467627-4225-52c2-ac7b-85a99ae3389a', 'be8b269c-bb94-518f-a7e6-d560c9b27bb4', 'Projeto com pendências', TRUE, NULLIF('Corrigir o projeto conforme o comunique-se', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c4467627-4225-52c2-ac7b-85a99ae3389a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c4467627-4225-52c2-ac7b-85a99ae3389a' AND z.destino_etapa_id = 'be8b269c-bb94-518f-a7e6-d560c9b27bb4' AND z.condicao_transicao = 'Projeto com pendências');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e', 'HAB.HAB.002', 'Habite-se', 'Vistoriar a obra concluída e emitir o habite-se (certificado de conclusão).', 'Proprietários', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '6.0.02.00.08'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.02.00.08' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.HAB.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '88aa1764-c09e-51d6-9418-c457b8cec327', 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e', 1, 'SEHAB/PROTOCOLO', 'Atendimento', 'Receber o pedido com o alvará e os laudos exigidos.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9f2f7d85-f482-540f-922c-2736f7cb22f5', '88aa1764-c09e-51d6-9418-c457b8cec327', 'Requerimento de habite-se', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '88aa1764-c09e-51d6-9418-c457b8cec327')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '88aa1764-c09e-51d6-9418-c457b8cec327' AND y.nome_documento = 'Requerimento de habite-se');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7a5e5de5-561d-5b9e-8d25-fbf142c236c8', '88aa1764-c09e-51d6-9418-c457b8cec327', 'Alvará de construção', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '88aa1764-c09e-51d6-9418-c457b8cec327')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '88aa1764-c09e-51d6-9418-c457b8cec327' AND y.nome_documento = 'Alvará de construção');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '1e15f2ea-303d-59b7-a316-6a5153104e6f', '88aa1764-c09e-51d6-9418-c457b8cec327', 'Certificado do Corpo de Bombeiros', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '88aa1764-c09e-51d6-9418-c457b8cec327')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '88aa1764-c09e-51d6-9418-c457b8cec327' AND y.nome_documento = 'Certificado do Corpo de Bombeiros');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '99379c25-e58b-50cc-aace-1778a7a2a48b', 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e', 2, 'SEHAB/FISCALIZACAO', 'Fiscalização de obras', 'Vistoriar a obra e conferir com o projeto aprovado.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a43984c5-8f65-5f97-af8c-c6802c63b493', '99379c25-e58b-50cc-aace-1778a7a2a48b', 'Relatório de vistoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '99379c25-e58b-50cc-aace-1778a7a2a48b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '99379c25-e58b-50cc-aace-1778a7a2a48b' AND y.nome_documento = 'Relatório de vistoria');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '77682727-86c2-5bc4-9922-15506a405331', 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e', 3, 'SEHAB/APROVACAO', 'Aprovação de projetos', 'Emitir o habite-se.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'eac5fbae-e6e9-549b-9123-ce4626ee4f7e' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2fc55ae8-4a3b-5a88-a8e2-8097f8836770', '77682727-86c2-5bc4-9922-15506a405331', 'Habite-se', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '77682727-86c2-5bc4-9922-15506a405331')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '77682727-86c2-5bc4-9922-15506a405331' AND y.nome_documento = 'Habite-se');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'eb15bc0c-8f63-5ee6-a21e-745357270ae4', '88aa1764-c09e-51d6-9418-c457b8cec327', '99379c25-e58b-50cc-aace-1778a7a2a48b', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '88aa1764-c09e-51d6-9418-c457b8cec327')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '99379c25-e58b-50cc-aace-1778a7a2a48b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '88aa1764-c09e-51d6-9418-c457b8cec327' AND z.destino_etapa_id = '99379c25-e58b-50cc-aace-1778a7a2a48b' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0f367e08-5be3-5193-a61c-12df8a043c97', '99379c25-e58b-50cc-aace-1778a7a2a48b', '77682727-86c2-5bc4-9922-15506a405331', 'Obra conforme', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '99379c25-e58b-50cc-aace-1778a7a2a48b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '77682727-86c2-5bc4-9922-15506a405331')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '99379c25-e58b-50cc-aace-1778a7a2a48b' AND z.destino_etapa_id = '77682727-86c2-5bc4-9922-15506a405331' AND z.condicao_transicao = 'Obra conforme');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '1dbeaffc-5854-5208-8f6c-2d8038252d21', '99379c25-e58b-50cc-aace-1778a7a2a48b', '88aa1764-c09e-51d6-9418-c457b8cec327', 'Obra em desacordo com o projeto', TRUE, NULLIF('Adequar a obra ou aprovar projeto modificativo', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '99379c25-e58b-50cc-aace-1778a7a2a48b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '88aa1764-c09e-51d6-9418-c457b8cec327')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '99379c25-e58b-50cc-aace-1778a7a2a48b' AND z.destino_etapa_id = '88aa1764-c09e-51d6-9418-c457b8cec327' AND z.condicao_transicao = 'Obra em desacordo com o projeto');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'e5fb4cc4-479e-5f22-a19c-028e466c18fd', 'HAB.REU.003', 'Regularização fundiária de interesse social (REURB-S)', 'Regularizar a ocupação de núcleo urbano informal de baixa renda e titular o ocupante (Lei 13.465/2017).', 'Ocupantes de núcleos informais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (dados socioeconômicos do ocupante)', ''), '6.0.01.01.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.01.01.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.REU.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '128d7999-949e-51cf-8430-496a7865f6b9', 'e5fb4cc4-479e-5f22-a19c-028e466c18fd', 1, 'SEHAB/REURB', 'Regularização fundiária', 'Cadastrar o ocupante e enquadrar a modalidade REURB-S.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd0b5e49a-a97b-5c7b-ba68-82e04176bb71', '128d7999-949e-51cf-8430-496a7865f6b9', 'Requerimento de REURB', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '128d7999-949e-51cf-8430-496a7865f6b9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '128d7999-949e-51cf-8430-496a7865f6b9' AND y.nome_documento = 'Requerimento de REURB');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ce26f938-d5cd-55f7-8864-dae0f28b6965', '128d7999-949e-51cf-8430-496a7865f6b9', 'Cadastro socioeconômico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '128d7999-949e-51cf-8430-496a7865f6b9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '128d7999-949e-51cf-8430-496a7865f6b9' AND y.nome_documento = 'Cadastro socioeconômico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'de294172-8693-5698-b52a-1d9fb99e6ac2', '128d7999-949e-51cf-8430-496a7865f6b9', 'Documentos pessoais do ocupante', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '128d7999-949e-51cf-8430-496a7865f6b9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '128d7999-949e-51cf-8430-496a7865f6b9' AND y.nome_documento = 'Documentos pessoais do ocupante');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '2555b833-4da9-5c3d-83e0-6c38cd41658d', 'e5fb4cc4-479e-5f22-a19c-028e466c18fd', 2, 'SEHAB/TOPOGRAFIA', 'Topografia e cartografia', 'Elaborar o levantamento e o projeto de regularização.', 60, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3b74c8c8-7152-5cc0-8d3c-0f45fe6bb6b5', '2555b833-4da9-5c3d-83e0-6c38cd41658d', 'Levantamento topográfico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '2555b833-4da9-5c3d-83e0-6c38cd41658d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '2555b833-4da9-5c3d-83e0-6c38cd41658d' AND y.nome_documento = 'Levantamento topográfico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '355a2583-2a6b-55fd-abaf-73cf5662bc3a', '2555b833-4da9-5c3d-83e0-6c38cd41658d', 'Projeto de regularização fundiária', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '2555b833-4da9-5c3d-83e0-6c38cd41658d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '2555b833-4da9-5c3d-83e0-6c38cd41658d' AND y.nome_documento = 'Projeto de regularização fundiária');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '32431e81-1e18-56a1-8234-c3f338cfed32', 'e5fb4cc4-479e-5f22-a19c-028e466c18fd', 3, 'PGM', 'Procuradoria', 'Analisar a regularidade e notificar confrontantes.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4f1713c5-e342-5309-8b3a-c775e86f5aeb', '32431e81-1e18-56a1-8234-c3f338cfed32', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '32431e81-1e18-56a1-8234-c3f338cfed32')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '32431e81-1e18-56a1-8234-c3f338cfed32' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b5692b4c-6191-54e1-bd34-19670598d25a', '32431e81-1e18-56a1-8234-c3f338cfed32', 'Notificação aos confrontantes', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '32431e81-1e18-56a1-8234-c3f338cfed32')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '32431e81-1e18-56a1-8234-c3f338cfed32' AND y.nome_documento = 'Notificação aos confrontantes');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '33dae2bf-413c-5ada-bd35-32767fdaffa8', 'e5fb4cc4-479e-5f22-a19c-028e466c18fd', 4, 'SEHAB/REURB', 'Regularização fundiária', 'Emitir a Certidão de Regularização Fundiária (CRF) e encaminhar ao cartório.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e5fb4cc4-479e-5f22-a19c-028e466c18fd' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '33689e42-5f08-5b37-9190-c9edfe3e8ed4', '33dae2bf-413c-5ada-bd35-32767fdaffa8', 'Certidão de Regularização Fundiária (CRF)', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '33dae2bf-413c-5ada-bd35-32767fdaffa8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '33dae2bf-413c-5ada-bd35-32767fdaffa8' AND y.nome_documento = 'Certidão de Regularização Fundiária (CRF)');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c91dd367-0b2d-5fde-9f32-3a700513d9f6', '128d7999-949e-51cf-8430-496a7865f6b9', '2555b833-4da9-5c3d-83e0-6c38cd41658d', 'Modalidade enquadrada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '128d7999-949e-51cf-8430-496a7865f6b9')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '2555b833-4da9-5c3d-83e0-6c38cd41658d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '128d7999-949e-51cf-8430-496a7865f6b9' AND z.destino_etapa_id = '2555b833-4da9-5c3d-83e0-6c38cd41658d' AND z.condicao_transicao = 'Modalidade enquadrada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e3a11613-da05-53c5-8bf1-5faaac2ea778', '2555b833-4da9-5c3d-83e0-6c38cd41658d', '32431e81-1e18-56a1-8234-c3f338cfed32', 'Projeto concluído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '2555b833-4da9-5c3d-83e0-6c38cd41658d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '32431e81-1e18-56a1-8234-c3f338cfed32')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '2555b833-4da9-5c3d-83e0-6c38cd41658d' AND z.destino_etapa_id = '32431e81-1e18-56a1-8234-c3f338cfed32' AND z.condicao_transicao = 'Projeto concluído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '7a4e3d61-59fb-5631-b6f5-b8d35459be01', '32431e81-1e18-56a1-8234-c3f338cfed32', '33dae2bf-413c-5ada-bd35-32767fdaffa8', 'Sem impugnação', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '32431e81-1e18-56a1-8234-c3f338cfed32')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '33dae2bf-413c-5ada-bd35-32767fdaffa8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '32431e81-1e18-56a1-8234-c3f338cfed32' AND z.destino_etapa_id = '33dae2bf-413c-5ada-bd35-32767fdaffa8' AND z.condicao_transicao = 'Sem impugnação');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'ce0a6b2c-3835-51fc-93ab-69c7641893de', '32431e81-1e18-56a1-8234-c3f338cfed32', '2555b833-4da9-5c3d-83e0-6c38cd41658d', 'Impugnação de confrontante', TRUE, NULLIF('Analisar a impugnação e ajustar o projeto', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '32431e81-1e18-56a1-8234-c3f338cfed32')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '2555b833-4da9-5c3d-83e0-6c38cd41658d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '32431e81-1e18-56a1-8234-c3f338cfed32' AND z.destino_etapa_id = '2555b833-4da9-5c3d-83e0-6c38cd41658d' AND z.condicao_transicao = 'Impugnação de confrontante');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe', 'MAM.LIC.001', 'Licenciamento ambiental municipal', 'Licenciar atividade de impacto ambiental local (licença prévia, de instalação e de operação).', 'Empreendedores', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.01.01.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.01.01.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'MAM.LIC.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa', '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe', 1, 'SEMMA/LICENCIAMENTO', 'Licenciamento ambiental', 'Receber o requerimento e enquadrar a atividade e o porte.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7e7104d3-422b-51b3-abbd-fbf736c52389', '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa', 'Requerimento de licença ambiental', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa' AND y.nome_documento = 'Requerimento de licença ambiental');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6c43efcd-cf8d-5fbe-9a4e-7127e5d9a7b3', '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa', 'Estudo ambiental exigido', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa' AND y.nome_documento = 'Estudo ambiental exigido');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '06fbedb3-08b2-5d2c-91f1-85301b8a5eea', '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa', 'Certidão de uso do solo', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa' AND y.nome_documento = 'Certidão de uso do solo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '739331ad-8ee1-53ed-b7e7-8e747278e76f', '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe', 2, 'SEMMA/TECNICO', 'Análise técnica', 'Analisar os estudos e vistoriar o local.', 45, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '534bf299-cc14-557d-9e84-4f1212398245', '739331ad-8ee1-53ed-b7e7-8e747278e76f', 'Parecer técnico', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '739331ad-8ee1-53ed-b7e7-8e747278e76f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '739331ad-8ee1-53ed-b7e7-8e747278e76f' AND y.nome_documento = 'Parecer técnico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '98599043-0ddd-5a87-8c15-3d31117d3531', '739331ad-8ee1-53ed-b7e7-8e747278e76f', 'Relatório de vistoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '739331ad-8ee1-53ed-b7e7-8e747278e76f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '739331ad-8ee1-53ed-b7e7-8e747278e76f' AND y.nome_documento = 'Relatório de vistoria');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '399ae17e-6980-57be-b275-81fb42ff00bc', '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe', 3, 'SEMMA/LICENCIAMENTO', 'Licenciamento ambiental', 'Emitir a licença com as condicionantes.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9b7aa1a9-e9ef-5ed9-9b31-8b2cc4591ebe' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8d6769ac-0305-5c88-8509-86c56a6b5d00', '399ae17e-6980-57be-b275-81fb42ff00bc', 'Licença ambiental', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '399ae17e-6980-57be-b275-81fb42ff00bc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '399ae17e-6980-57be-b275-81fb42ff00bc' AND y.nome_documento = 'Licença ambiental');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4d4dd94e-03fb-50a5-b352-8fc022fa677d', '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa', '739331ad-8ee1-53ed-b7e7-8e747278e76f', 'Atividade enquadrada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '739331ad-8ee1-53ed-b7e7-8e747278e76f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa' AND z.destino_etapa_id = '739331ad-8ee1-53ed-b7e7-8e747278e76f' AND z.condicao_transicao = 'Atividade enquadrada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '7fab083b-0632-5d56-9b07-18f878ddc4cc', '739331ad-8ee1-53ed-b7e7-8e747278e76f', '399ae17e-6980-57be-b275-81fb42ff00bc', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '739331ad-8ee1-53ed-b7e7-8e747278e76f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '399ae17e-6980-57be-b275-81fb42ff00bc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '739331ad-8ee1-53ed-b7e7-8e747278e76f' AND z.destino_etapa_id = '399ae17e-6980-57be-b275-81fb42ff00bc' AND z.condicao_transicao = 'Parecer favorável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2588b24c-0649-59ba-b3fc-d413e3188936', '739331ad-8ee1-53ed-b7e7-8e747278e76f', '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa', 'Estudo insuficiente', TRUE, NULLIF('Complementar o estudo ambiental', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '739331ad-8ee1-53ed-b7e7-8e747278e76f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '739331ad-8ee1-53ed-b7e7-8e747278e76f' AND z.destino_etapa_id = '4a1eaa0b-2661-56f2-919d-9ba8b24eccaa' AND z.condicao_transicao = 'Estudo insuficiente');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'e30da675-36d0-5e4e-afe5-108eb950ea05', 'DEC.USO.001', 'Uso do solo para profissional liberal', 'Emitir a certidão de uso do solo que autoriza o profissional liberal a exercer a atividade no endereço informado.', 'Profissionais liberais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '9.0.01.01.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '9.0.01.01.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'DEC.USO.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '274101b3-ce66-5bb2-bbfc-06e08d5dc55b', 'e30da675-36d0-5e4e-afe5-108eb950ea05', 1, 'SEDEC/ATENDIMENTO', 'Sala do Empreendedor', 'Receber o requerimento e o comprovante de endereço.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e30da675-36d0-5e4e-afe5-108eb950ea05')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e30da675-36d0-5e4e-afe5-108eb950ea05' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ccb666cf-16ec-5f3f-a9de-6064f056ce86', '274101b3-ce66-5bb2-bbfc-06e08d5dc55b', 'Requerimento de uso do solo para profissional liberal', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b' AND y.nome_documento = 'Requerimento de uso do solo para profissional liberal');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bd6e8c98-d94f-5718-bba4-32a97f43921e', '274101b3-ce66-5bb2-bbfc-06e08d5dc55b', 'Registro no conselho profissional', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b' AND y.nome_documento = 'Registro no conselho profissional');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'aac725e5-f070-5542-aed4-164629fd7673', '274101b3-ce66-5bb2-bbfc-06e08d5dc55b', 'Comprovante de endereço', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b' AND y.nome_documento = 'Comprovante de endereço');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '76243091-8487-5c7c-9eb7-60d054000ec2', 'e30da675-36d0-5e4e-afe5-108eb950ea05', 2, 'SEPLAN/URBANISMO', 'Planejamento urbano', 'Conferir o zoneamento do endereço para a atividade.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e30da675-36d0-5e4e-afe5-108eb950ea05')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e30da675-36d0-5e4e-afe5-108eb950ea05' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7ca6b57b-3024-5ef2-8a02-4438574b41f2', '76243091-8487-5c7c-9eb7-60d054000ec2', 'Análise de zoneamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '76243091-8487-5c7c-9eb7-60d054000ec2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '76243091-8487-5c7c-9eb7-60d054000ec2' AND y.nome_documento = 'Análise de zoneamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '22eb2e4c-2a36-57c8-8859-e8024d969a48', 'e30da675-36d0-5e4e-afe5-108eb950ea05', 3, 'SEDEC/ATENDIMENTO', 'Sala do Empreendedor', 'Emitir a certidão.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e30da675-36d0-5e4e-afe5-108eb950ea05')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e30da675-36d0-5e4e-afe5-108eb950ea05' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '810b78d7-d3c0-5e08-8a65-4597d3f81c3b', '22eb2e4c-2a36-57c8-8859-e8024d969a48', 'Certidão de uso do solo para profissional liberal', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '22eb2e4c-2a36-57c8-8859-e8024d969a48')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '22eb2e4c-2a36-57c8-8859-e8024d969a48' AND y.nome_documento = 'Certidão de uso do solo para profissional liberal');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '165878d9-3231-5dea-b579-51879bf93806', '274101b3-ce66-5bb2-bbfc-06e08d5dc55b', '76243091-8487-5c7c-9eb7-60d054000ec2', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '76243091-8487-5c7c-9eb7-60d054000ec2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b' AND z.destino_etapa_id = '76243091-8487-5c7c-9eb7-60d054000ec2' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'dddd727e-b03c-547f-9011-badeee1fe5d3', '76243091-8487-5c7c-9eb7-60d054000ec2', '22eb2e4c-2a36-57c8-8859-e8024d969a48', 'Uso permitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '76243091-8487-5c7c-9eb7-60d054000ec2')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '22eb2e4c-2a36-57c8-8859-e8024d969a48')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '76243091-8487-5c7c-9eb7-60d054000ec2' AND z.destino_etapa_id = '22eb2e4c-2a36-57c8-8859-e8024d969a48' AND z.condicao_transicao = 'Uso permitido');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fcf44c1a-00ec-57e4-b36c-88d97dedce88', '76243091-8487-5c7c-9eb7-60d054000ec2', '274101b3-ce66-5bb2-bbfc-06e08d5dc55b', 'Uso não permitido no endereço', TRUE, NULLIF('Informar o requerente e indicar alternativas', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '76243091-8487-5c7c-9eb7-60d054000ec2')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '76243091-8487-5c7c-9eb7-60d054000ec2' AND z.destino_etapa_id = '274101b3-ce66-5bb2-bbfc-06e08d5dc55b' AND z.condicao_transicao = 'Uso não permitido no endereço');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '257a5a2f-d55d-5c4d-8c47-54446b1ed035', 'SAS.OSC.001', 'Inscrição de entidade (OSC) no Conselho de Assistência Social', 'Inscrever a organização da sociedade civil no CMAS para atuar na política de assistência social.', 'Organizações da sociedade civil', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '10.0.02.00.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '10.0.02.00.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAS.OSC.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '015f0e58-49e7-5608-b99d-adfde4220bd7', '257a5a2f-d55d-5c4d-8c47-54446b1ed035', 1, 'SEMPRAS/CMAS', 'Secretaria executiva do CMAS', 'Receber o pedido de inscrição com estatuto, atas e plano de ação.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '257a5a2f-d55d-5c4d-8c47-54446b1ed035')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '257a5a2f-d55d-5c4d-8c47-54446b1ed035' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'df9e3501-64d4-5361-a7af-8bce38959611', '015f0e58-49e7-5608-b99d-adfde4220bd7', 'Requerimento de inscrição', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '015f0e58-49e7-5608-b99d-adfde4220bd7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '015f0e58-49e7-5608-b99d-adfde4220bd7' AND y.nome_documento = 'Requerimento de inscrição');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '95f24a52-37c2-512f-a2c0-2efd6fac9cf9', '015f0e58-49e7-5608-b99d-adfde4220bd7', 'Estatuto e ata de eleição da diretoria', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '015f0e58-49e7-5608-b99d-adfde4220bd7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '015f0e58-49e7-5608-b99d-adfde4220bd7' AND y.nome_documento = 'Estatuto e ata de eleição da diretoria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '38b91e58-92e8-525b-8b12-0053f502f8ec', '015f0e58-49e7-5608-b99d-adfde4220bd7', 'Plano de ação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '015f0e58-49e7-5608-b99d-adfde4220bd7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '015f0e58-49e7-5608-b99d-adfde4220bd7' AND y.nome_documento = 'Plano de ação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '08af1b1b-44f1-54e6-b3c8-38568fb38458', '257a5a2f-d55d-5c4d-8c47-54446b1ed035', 2, 'SEMPRAS/VIGILANCIA', 'Vigilância socioassistencial', 'Visitar a entidade e emitir o parecer técnico.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '257a5a2f-d55d-5c4d-8c47-54446b1ed035')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '257a5a2f-d55d-5c4d-8c47-54446b1ed035' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '33fcf76c-9323-5acb-bd79-8c92e6fc4fc0', '08af1b1b-44f1-54e6-b3c8-38568fb38458', 'Relatório de visita técnica', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '08af1b1b-44f1-54e6-b3c8-38568fb38458')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '08af1b1b-44f1-54e6-b3c8-38568fb38458' AND y.nome_documento = 'Relatório de visita técnica');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '72d01283-a8b1-554e-b716-6cf62f4ab6b1', '257a5a2f-d55d-5c4d-8c47-54446b1ed035', 3, 'SEMPRAS/CMAS', 'Plenária do CMAS', 'Deliberar a inscrição e publicar a resolução.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '257a5a2f-d55d-5c4d-8c47-54446b1ed035')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '257a5a2f-d55d-5c4d-8c47-54446b1ed035' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '434a46c7-8f78-5fc8-ba8d-be2e9e9ceaf4', '72d01283-a8b1-554e-b716-6cf62f4ab6b1', 'Resolução do CMAS', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '72d01283-a8b1-554e-b716-6cf62f4ab6b1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '72d01283-a8b1-554e-b716-6cf62f4ab6b1' AND y.nome_documento = 'Resolução do CMAS');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'eaf48594-1703-56cb-84bc-4609d224a342', '015f0e58-49e7-5608-b99d-adfde4220bd7', '08af1b1b-44f1-54e6-b3c8-38568fb38458', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '015f0e58-49e7-5608-b99d-adfde4220bd7')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '08af1b1b-44f1-54e6-b3c8-38568fb38458')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '015f0e58-49e7-5608-b99d-adfde4220bd7' AND z.destino_etapa_id = '08af1b1b-44f1-54e6-b3c8-38568fb38458' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '424c80bc-56d5-597a-9beb-ae8ac7a71a0e', '08af1b1b-44f1-54e6-b3c8-38568fb38458', '72d01283-a8b1-554e-b716-6cf62f4ab6b1', 'Parecer emitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '08af1b1b-44f1-54e6-b3c8-38568fb38458')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '72d01283-a8b1-554e-b716-6cf62f4ab6b1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '08af1b1b-44f1-54e6-b3c8-38568fb38458' AND z.destino_etapa_id = '72d01283-a8b1-554e-b716-6cf62f4ab6b1' AND z.condicao_transicao = 'Parecer emitido');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'dc8aa47f-0fea-5714-9d32-b136c9fa9b24', '08af1b1b-44f1-54e6-b3c8-38568fb38458', '015f0e58-49e7-5608-b99d-adfde4220bd7', 'Documentação incompleta', TRUE, NULLIF('Complementar a documentação', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '08af1b1b-44f1-54e6-b3c8-38568fb38458')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '015f0e58-49e7-5608-b99d-adfde4220bd7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '08af1b1b-44f1-54e6-b3c8-38568fb38458' AND z.destino_etapa_id = '015f0e58-49e7-5608-b99d-adfde4220bd7' AND z.condicao_transicao = 'Documentação incompleta');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '236253fb-4036-5b77-bd97-116b056488b9', 'SAU.VSA.001', 'Licença sanitária de estabelecimento', 'Licenciar estabelecimento sujeito à vigilância sanitária, com inspeção (Código Sanitário Municipal).', 'Estabelecimentos de interesse da saúde', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.10.02.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.10.02.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.VSA.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8626abde-a044-5c9f-b62b-f7b095ed1f6a', '236253fb-4036-5b77-bd97-116b056488b9', 1, 'SMS/VISA', 'Vigilância sanitária', 'Receber o pedido e cadastrar o estabelecimento.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '236253fb-4036-5b77-bd97-116b056488b9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '236253fb-4036-5b77-bd97-116b056488b9' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6f211cbd-4bbd-5026-a2a1-2c5756faf7b7', '8626abde-a044-5c9f-b62b-f7b095ed1f6a', 'Requerimento de licença sanitária', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a' AND y.nome_documento = 'Requerimento de licença sanitária');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7392ccd4-18eb-5a2b-8092-da203101fbe4', '8626abde-a044-5c9f-b62b-f7b095ed1f6a', 'Responsável técnico (quando exigido)', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a' AND y.nome_documento = 'Responsável técnico (quando exigido)');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '244743bd-e924-537f-bd2a-577b5cfdb30e', '236253fb-4036-5b77-bd97-116b056488b9', 2, 'SMS/VISA', 'Inspeção sanitária', 'Inspecionar o estabelecimento e lavrar o auto.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '236253fb-4036-5b77-bd97-116b056488b9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '236253fb-4036-5b77-bd97-116b056488b9' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c542da91-7cdc-565e-aacd-50b13055269e', '244743bd-e924-537f-bd2a-577b5cfdb30e', 'Auto de inspeção', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '244743bd-e924-537f-bd2a-577b5cfdb30e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '244743bd-e924-537f-bd2a-577b5cfdb30e' AND y.nome_documento = 'Auto de inspeção');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '89009319-89e7-504e-a6f2-96839bb570e2', '236253fb-4036-5b77-bd97-116b056488b9', 3, 'SMS/VISA', 'Vigilância sanitária', 'Emitir a licença (alvará) sanitária.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '236253fb-4036-5b77-bd97-116b056488b9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '236253fb-4036-5b77-bd97-116b056488b9' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f001d8d7-5238-5daa-9571-14802a2d60dd', '89009319-89e7-504e-a6f2-96839bb570e2', 'Licença sanitária', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '89009319-89e7-504e-a6f2-96839bb570e2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '89009319-89e7-504e-a6f2-96839bb570e2' AND y.nome_documento = 'Licença sanitária');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'cf033996-fe0d-5f6a-94a7-90aa5f452b3e', '8626abde-a044-5c9f-b62b-f7b095ed1f6a', '244743bd-e924-537f-bd2a-577b5cfdb30e', 'Cadastro feito', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '244743bd-e924-537f-bd2a-577b5cfdb30e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a' AND z.destino_etapa_id = '244743bd-e924-537f-bd2a-577b5cfdb30e' AND z.condicao_transicao = 'Cadastro feito');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4d72b15e-f11c-55d4-92c3-54b55817922b', '244743bd-e924-537f-bd2a-577b5cfdb30e', '89009319-89e7-504e-a6f2-96839bb570e2', 'Inspeção favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '244743bd-e924-537f-bd2a-577b5cfdb30e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '89009319-89e7-504e-a6f2-96839bb570e2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '244743bd-e924-537f-bd2a-577b5cfdb30e' AND z.destino_etapa_id = '89009319-89e7-504e-a6f2-96839bb570e2' AND z.condicao_transicao = 'Inspeção favorável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3c388cac-2bf1-5705-aa55-8c14149a1a74', '244743bd-e924-537f-bd2a-577b5cfdb30e', '8626abde-a044-5c9f-b62b-f7b095ed1f6a', 'Irregularidade sanitária', TRUE, NULLIF('Corrigir as exigências do auto de notificação e pedir nova inspeção', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '244743bd-e924-537f-bd2a-577b5cfdb30e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '244743bd-e924-537f-bd2a-577b5cfdb30e' AND z.destino_etapa_id = '8626abde-a044-5c9f-b62b-f7b095ed1f6a' AND z.condicao_transicao = 'Irregularidade sanitária');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '0e32885e-6de5-53f4-93c7-52c2fe39e9c8', 'SAU.TFD.001', 'Tratamento Fora do Domicílio (TFD)', 'Autorizar o deslocamento do paciente do SUS para tratamento em outro município quando não houver oferta local.', 'Usuários do SUS', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '12.0.06.01.08'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.06.01.08' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.TFD.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ab2bfe6b-c24a-50c3-8c43-542191916821', '0e32885e-6de5-53f4-93c7-52c2fe39e9c8', 1, 'SMS/TFD', 'Setor de TFD', 'Receber o laudo médico e o pedido do paciente.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0e32885e-6de5-53f4-93c7-52c2fe39e9c8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0e32885e-6de5-53f4-93c7-52c2fe39e9c8' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ee7f0e57-a638-5ec5-809b-85de3585b299', 'ab2bfe6b-c24a-50c3-8c43-542191916821', 'Laudo médico de TFD', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ab2bfe6b-c24a-50c3-8c43-542191916821')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ab2bfe6b-c24a-50c3-8c43-542191916821' AND y.nome_documento = 'Laudo médico de TFD');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5e6e9a9f-67b4-5f20-b1a7-7549861a8e16', 'ab2bfe6b-c24a-50c3-8c43-542191916821', 'Documentos pessoais e cartão SUS', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ab2bfe6b-c24a-50c3-8c43-542191916821')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ab2bfe6b-c24a-50c3-8c43-542191916821' AND y.nome_documento = 'Documentos pessoais e cartão SUS');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c8bdb041-376a-543a-8198-8baa471f13ab', '0e32885e-6de5-53f4-93c7-52c2fe39e9c8', 2, 'SMS/REGULACAO', 'Complexo regulador', 'Avaliar a indicação e agendar no município de referência.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0e32885e-6de5-53f4-93c7-52c2fe39e9c8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0e32885e-6de5-53f4-93c7-52c2fe39e9c8' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e140d70a-36ea-5b3e-9ea2-0f2d66cfc28a', 'c8bdb041-376a-543a-8198-8baa471f13ab', 'Parecer da regulação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c8bdb041-376a-543a-8198-8baa471f13ab')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c8bdb041-376a-543a-8198-8baa471f13ab' AND y.nome_documento = 'Parecer da regulação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '003ce92b-afd1-5e80-a076-623db80324c5', 'c8bdb041-376a-543a-8198-8baa471f13ab', 'Comprovante de agendamento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c8bdb041-376a-543a-8198-8baa471f13ab')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c8bdb041-376a-543a-8198-8baa471f13ab' AND y.nome_documento = 'Comprovante de agendamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0', '0e32885e-6de5-53f4-93c7-52c2fe39e9c8', 3, 'SMS/TFD', 'Setor de TFD', 'Providenciar o transporte e a ajuda de custo e controlar o retorno.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0e32885e-6de5-53f4-93c7-52c2fe39e9c8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0e32885e-6de5-53f4-93c7-52c2fe39e9c8' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b6b93f80-14a3-5e34-aa81-5523fbbe8568', '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0', 'Lista de embarque', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0' AND y.nome_documento = 'Lista de embarque');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8674d8f8-b886-537f-a8dd-29bfcf867c40', '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0', 'Comprovante de ajuda de custo', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0' AND y.nome_documento = 'Comprovante de ajuda de custo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'a5774a82-19ec-587d-b14a-536229ec4229', 'ab2bfe6b-c24a-50c3-8c43-542191916821', 'c8bdb041-376a-543a-8198-8baa471f13ab', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ab2bfe6b-c24a-50c3-8c43-542191916821')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c8bdb041-376a-543a-8198-8baa471f13ab')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ab2bfe6b-c24a-50c3-8c43-542191916821' AND z.destino_etapa_id = 'c8bdb041-376a-543a-8198-8baa471f13ab' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'cb1abdba-83dc-5720-add0-258b571ffa10', 'c8bdb041-376a-543a-8198-8baa471f13ab', '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0', 'Agendado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c8bdb041-376a-543a-8198-8baa471f13ab')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c8bdb041-376a-543a-8198-8baa471f13ab' AND z.destino_etapa_id = '7a4d1e2e-eba1-56a0-a0d4-ae26fbc734f0' AND z.condicao_transicao = 'Agendado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'cb6e5e39-04f2-514a-a0ca-732305ddaa52', 'c8bdb041-376a-543a-8198-8baa471f13ab', 'ab2bfe6b-c24a-50c3-8c43-542191916821', 'Laudo incompleto', TRUE, NULLIF('Complementar o laudo médico', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c8bdb041-376a-543a-8198-8baa471f13ab')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ab2bfe6b-c24a-50c3-8c43-542191916821')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c8bdb041-376a-543a-8198-8baa471f13ab' AND z.destino_etapa_id = 'ab2bfe6b-c24a-50c3-8c43-542191916821' AND z.condicao_transicao = 'Laudo incompleto');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ee583489-22fd-553c-b6d6-81f4c21eeac8', 'ADM.CEM.010', 'Concessão de jazigo em cemitério municipal', 'Conceder o direito de uso de jazigo em cemitério municipal, com o pagamento das taxas.', 'Cidadãos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.01.03.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.01.03.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.CEM.010');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8c6a83d8-304e-5556-9f16-d8cd0961d615', 'ee583489-22fd-553c-b6d6-81f4c21eeac8', 1, 'SEMAD/CEMITERIOS', 'Administração dos cemitérios', 'Receber o pedido e verificar a disponibilidade de jazigo.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ee583489-22fd-553c-b6d6-81f4c21eeac8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ee583489-22fd-553c-b6d6-81f4c21eeac8' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2be67801-43b5-57cf-ad5d-7dad5a0d9755', '8c6a83d8-304e-5556-9f16-d8cd0961d615', 'Requerimento de concessão de jazigo', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8c6a83d8-304e-5556-9f16-d8cd0961d615')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8c6a83d8-304e-5556-9f16-d8cd0961d615' AND y.nome_documento = 'Requerimento de concessão de jazigo');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5cbc5c4d-0d80-5ca4-9749-53175ed0992f', '8c6a83d8-304e-5556-9f16-d8cd0961d615', 'Certidão de óbito (quando houver)', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8c6a83d8-304e-5556-9f16-d8cd0961d615')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8c6a83d8-304e-5556-9f16-d8cd0961d615' AND y.nome_documento = 'Certidão de óbito (quando houver)');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b6e26b75-963b-5f18-9a08-769995bfa2b5', 'ee583489-22fd-553c-b6d6-81f4c21eeac8', 2, 'SEREC/ARRECADACAO', 'Arrecadação', 'Emitir a guia das taxas.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ee583489-22fd-553c-b6d6-81f4c21eeac8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ee583489-22fd-553c-b6d6-81f4c21eeac8' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'eefd0ec2-45dd-542d-8f8f-a0bf0f111dd7', 'b6e26b75-963b-5f18-9a08-769995bfa2b5', 'Guia de pagamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b6e26b75-963b-5f18-9a08-769995bfa2b5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b6e26b75-963b-5f18-9a08-769995bfa2b5' AND y.nome_documento = 'Guia de pagamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ec63d5ab-1997-5ec8-af87-76b502131fab', 'ee583489-22fd-553c-b6d6-81f4c21eeac8', 3, 'SEMAD/CEMITERIOS', 'Administração dos cemitérios', 'Lavrar o termo de concessão e registrar no livro do cemitério.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ee583489-22fd-553c-b6d6-81f4c21eeac8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ee583489-22fd-553c-b6d6-81f4c21eeac8' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f6a0974b-7be2-5eba-858f-f4a65905aea4', 'ec63d5ab-1997-5ec8-af87-76b502131fab', 'Termo de concessão de uso', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ec63d5ab-1997-5ec8-af87-76b502131fab')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ec63d5ab-1997-5ec8-af87-76b502131fab' AND y.nome_documento = 'Termo de concessão de uso');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c924d0b3-71c7-5911-b550-c00457386858', '8c6a83d8-304e-5556-9f16-d8cd0961d615', 'b6e26b75-963b-5f18-9a08-769995bfa2b5', 'Jazigo disponível', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8c6a83d8-304e-5556-9f16-d8cd0961d615')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b6e26b75-963b-5f18-9a08-769995bfa2b5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8c6a83d8-304e-5556-9f16-d8cd0961d615' AND z.destino_etapa_id = 'b6e26b75-963b-5f18-9a08-769995bfa2b5' AND z.condicao_transicao = 'Jazigo disponível');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fb146c41-530a-5be1-85cc-5e79cd3f0a0c', 'b6e26b75-963b-5f18-9a08-769995bfa2b5', 'ec63d5ab-1997-5ec8-af87-76b502131fab', 'Taxa paga', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'b6e26b75-963b-5f18-9a08-769995bfa2b5')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ec63d5ab-1997-5ec8-af87-76b502131fab')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'b6e26b75-963b-5f18-9a08-769995bfa2b5' AND z.destino_etapa_id = 'ec63d5ab-1997-5ec8-af87-76b502131fab' AND z.condicao_transicao = 'Taxa paga');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ce252457-d307-5f34-b4b5-a0b144d26245', 'ADM.PRO.011', 'Requerimento geral ao protocolo', 'Receber e encaminhar requerimento que não tem procedimento próprio, até a unidade competente responder.', 'Cidadãos e servidores', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.03.00.46'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.03.00.46' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.PRO.011');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d', 'ce252457-d307-5f34-b4b5-a0b144d26245', 1, 'SEMAD/PROTOCOLO', 'Protocolo geral', 'Registrar o requerimento, autuar e encaminhar à unidade competente.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ce252457-d307-5f34-b4b5-a0b144d26245')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ce252457-d307-5f34-b4b5-a0b144d26245' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '624954c9-4a05-543b-a079-dec68869831e', '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d', 'Requerimento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d' AND y.nome_documento = 'Requerimento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3705e5b9-acd3-566a-b562-b5f6952166ce', '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d', 'Documentos anexos', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d' AND y.nome_documento = 'Documentos anexos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4eb8b801-14f3-5abd-8000-b835a2d3d2bf', 'ce252457-d307-5f34-b4b5-a0b144d26245', 2, 'SECRETARIA', 'Unidade competente', 'Analisar e decidir o pedido.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ce252457-d307-5f34-b4b5-a0b144d26245')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ce252457-d307-5f34-b4b5-a0b144d26245' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cfb8f7c1-dc18-5e33-84b0-662d7dce11ab', '4eb8b801-14f3-5abd-8000-b835a2d3d2bf', 'Despacho decisório', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf' AND y.nome_documento = 'Despacho decisório');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '235f3297-c6f5-5226-a6e1-f7573c5b397c', 'ce252457-d307-5f34-b4b5-a0b144d26245', 3, 'SEMAD/PROTOCOLO', 'Protocolo geral', 'Dar ciência ao requerente e arquivar.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ce252457-d307-5f34-b4b5-a0b144d26245')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ce252457-d307-5f34-b4b5-a0b144d26245' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e62a84bb-24d6-508a-bcce-547210749e3f', '235f3297-c6f5-5226-a6e1-f7573c5b397c', 'Comunicação ao requerente', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '235f3297-c6f5-5226-a6e1-f7573c5b397c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '235f3297-c6f5-5226-a6e1-f7573c5b397c' AND y.nome_documento = 'Comunicação ao requerente');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '8a736e31-9cd8-5476-8232-3b106c423f2b', '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d', '4eb8b801-14f3-5abd-8000-b835a2d3d2bf', 'Processo autuado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d' AND z.destino_etapa_id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf' AND z.condicao_transicao = 'Processo autuado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c5f29837-aa79-5830-8c3e-9ea023f74e25', '4eb8b801-14f3-5abd-8000-b835a2d3d2bf', '235f3297-c6f5-5226-a6e1-f7573c5b397c', 'Pedido decidido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '235f3297-c6f5-5226-a6e1-f7573c5b397c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf' AND z.destino_etapa_id = '235f3297-c6f5-5226-a6e1-f7573c5b397c' AND z.condicao_transicao = 'Pedido decidido');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e6f5271b-6069-5572-9ae6-2084f4abf94f', '4eb8b801-14f3-5abd-8000-b835a2d3d2bf', '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d', 'Unidade não competente', TRUE, NULLIF('Devolver ao protocolo para novo encaminhamento', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4eb8b801-14f3-5abd-8000-b835a2d3d2bf' AND z.destino_etapa_id = '7cc1249f-7f7f-522f-97c3-9ec596d7ba5d' AND z.condicao_transicao = 'Unidade não competente');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '824c6eb5-1215-5013-b9ee-25e9f95c6c89', 'ADM.DOC.012', 'Eliminação de documentos', 'Eliminar documentos com prazo de guarda cumprido e destinação "eliminação" na TTDD, com aprovação da CCPAD e publicação do edital de ciência.', 'Unidades administrativas', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.03.01.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.03.01.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'ADM.DOC.012');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e0da2316-3731-52c4-87ac-b3779afdba58', '824c6eb5-1215-5013-b9ee-25e9f95c6c89', 1, 'SECRETARIA', 'Unidade produtora', 'Levantar os documentos com prazo cumprido e preencher a listagem de eliminação.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c82cc806-81f3-590e-9a73-2021c7fb96ed', 'e0da2316-3731-52c4-87ac-b3779afdba58', 'Listagem de eliminação de documentos', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e0da2316-3731-52c4-87ac-b3779afdba58')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e0da2316-3731-52c4-87ac-b3779afdba58' AND y.nome_documento = 'Listagem de eliminação de documentos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '511766e5-f561-5bc8-9a5e-27b9c066ea8d', '824c6eb5-1215-5013-b9ee-25e9f95c6c89', 2, 'SEMAD/CCPAD', 'Comissão Permanente de Avaliação de Documentos', 'Conferir a listagem contra a TTDD e aprovar.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd198aad0-3a5a-5371-8baa-a75b8c18bccb', '511766e5-f561-5bc8-9a5e-27b9c066ea8d', 'Parecer da CCPAD', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d' AND y.nome_documento = 'Parecer da CCPAD');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967', '824c6eb5-1215-5013-b9ee-25e9f95c6c89', 3, 'SEMAD/CCPAD', 'Comissão Permanente de Avaliação de Documentos', 'Publicar o edital de ciência de eliminação e aguardar o prazo.', 45, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b93f733d-aabd-52da-ba91-b13136226349', 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967', 'Edital de ciência de eliminação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967' AND y.nome_documento = 'Edital de ciência de eliminação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b50141ee-e858-560d-ac27-a9641beca56c', '824c6eb5-1215-5013-b9ee-25e9f95c6c89', 4, 'SECRETARIA', 'Unidade produtora', 'Eliminar os documentos e lavrar o termo de eliminação.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '824c6eb5-1215-5013-b9ee-25e9f95c6c89' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c4db3813-0d02-5ea2-80c8-a3ac25a32bc9', 'b50141ee-e858-560d-ac27-a9641beca56c', 'Termo de eliminação de documentos', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b50141ee-e858-560d-ac27-a9641beca56c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b50141ee-e858-560d-ac27-a9641beca56c' AND y.nome_documento = 'Termo de eliminação de documentos');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f2455a0f-1900-519e-82ac-4227ca44dd9f', 'e0da2316-3731-52c4-87ac-b3779afdba58', '511766e5-f561-5bc8-9a5e-27b9c066ea8d', 'Listagem pronta', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e0da2316-3731-52c4-87ac-b3779afdba58')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e0da2316-3731-52c4-87ac-b3779afdba58' AND z.destino_etapa_id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d' AND z.condicao_transicao = 'Listagem pronta');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f50c80a9-6906-5019-90cc-5b5f50b746bc', '511766e5-f561-5bc8-9a5e-27b9c066ea8d', 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967', 'Listagem aprovada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d' AND z.destino_etapa_id = 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967' AND z.condicao_transicao = 'Listagem aprovada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f383276f-8e5e-599b-84c4-ac95707d2e79', '511766e5-f561-5bc8-9a5e-27b9c066ea8d', 'e0da2316-3731-52c4-87ac-b3779afdba58', 'Divergência com a TTDD', TRUE, NULLIF('Corrigir a listagem', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e0da2316-3731-52c4-87ac-b3779afdba58')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '511766e5-f561-5bc8-9a5e-27b9c066ea8d' AND z.destino_etapa_id = 'e0da2316-3731-52c4-87ac-b3779afdba58' AND z.condicao_transicao = 'Divergência com a TTDD');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '946d40df-2c9a-5fab-b60b-ba14b56c428d', 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967', 'b50141ee-e858-560d-ac27-a9641beca56c', 'Prazo do edital cumprido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b50141ee-e858-560d-ac27-a9641beca56c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'bcdfa22e-3f4c-5f8c-b1ae-400dd64bf967' AND z.destino_etapa_id = 'b50141ee-e858-560d-ac27-a9641beca56c' AND z.condicao_transicao = 'Prazo do edital cumprido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '6484d96d-72b6-5f3c-b363-88035d2307aa', 'TI.ACE.001', 'Acesso a sistemas, rede e e-mail corporativo', 'Criar, alterar ou revogar o acesso do servidor aos sistemas municipais, à rede e ao e-mail corporativo.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.04.01.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.04.01.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'TI.ACE.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '672ccf98-f986-55a1-a424-68e41592c025', '6484d96d-72b6-5f3c-b363-88035d2307aa', 1, 'SECRETARIA', 'Chefia imediata', 'Solicitar o acesso e indicar o perfil necessário à função.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '6484d96d-72b6-5f3c-b363-88035d2307aa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '6484d96d-72b6-5f3c-b363-88035d2307aa' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9f53a857-c677-5391-83a0-9e15feb93712', '672ccf98-f986-55a1-a424-68e41592c025', 'Requerimento de acesso', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '672ccf98-f986-55a1-a424-68e41592c025')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '672ccf98-f986-55a1-a424-68e41592c025' AND y.nome_documento = 'Requerimento de acesso');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd061eaac-87ce-5dd9-a08e-6b6eefd5dca3', '672ccf98-f986-55a1-a424-68e41592c025', 'Termo de responsabilidade de uso', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '672ccf98-f986-55a1-a424-68e41592c025')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '672ccf98-f986-55a1-a424-68e41592c025' AND y.nome_documento = 'Termo de responsabilidade de uso');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a2a07504-c097-5bfd-acf2-926b69ade361', '6484d96d-72b6-5f3c-b363-88035d2307aa', 2, 'SEMAD/TI', 'Tecnologia da Informação', 'Criar a conta com o perfil mínimo necessário e registrar.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '6484d96d-72b6-5f3c-b363-88035d2307aa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '6484d96d-72b6-5f3c-b363-88035d2307aa' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '67b2b2fe-e761-581c-8a69-e4f9d101a514', 'a2a07504-c097-5bfd-acf2-926b69ade361', 'Registro de concessão de acesso', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a2a07504-c097-5bfd-acf2-926b69ade361')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a2a07504-c097-5bfd-acf2-926b69ade361' AND y.nome_documento = 'Registro de concessão de acesso');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '94643859-d8c8-5fad-b2b2-ff813a96ed55', '672ccf98-f986-55a1-a424-68e41592c025', 'a2a07504-c097-5bfd-acf2-926b69ade361', 'Chefia de acordo', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '672ccf98-f986-55a1-a424-68e41592c025')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'a2a07504-c097-5bfd-acf2-926b69ade361')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '672ccf98-f986-55a1-a424-68e41592c025' AND z.destino_etapa_id = 'a2a07504-c097-5bfd-acf2-926b69ade361' AND z.condicao_transicao = 'Chefia de acordo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '32873797-5a7f-568b-b96b-d64b02ad2803', 'a2a07504-c097-5bfd-acf2-926b69ade361', '672ccf98-f986-55a1-a424-68e41592c025', 'Perfil incompatível com a função', TRUE, NULLIF('Rever o perfil solicitado', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a2a07504-c097-5bfd-acf2-926b69ade361')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '672ccf98-f986-55a1-a424-68e41592c025')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a2a07504-c097-5bfd-acf2-926b69ade361' AND z.destino_etapa_id = '672ccf98-f986-55a1-a424-68e41592c025' AND z.condicao_transicao = 'Perfil incompatível com a função');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '016c58e6-1945-5263-a62c-b092f44e2843', 'TI.SUP.002', 'Suporte técnico e manutenção de equipamento', 'Atender chamado de suporte ou manutenção de equipamento de informática.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.04.01.09'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.04.01.09' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'TI.SUP.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '30417b6c-d7b2-5664-9c49-f644680b44de', '016c58e6-1945-5263-a62c-b092f44e2843', 1, 'SECRETARIA', 'Unidade solicitante', 'Abrir o chamado descrevendo o problema e o patrimônio do equipamento.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '016c58e6-1945-5263-a62c-b092f44e2843')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '016c58e6-1945-5263-a62c-b092f44e2843' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '42812e1a-c5cd-5055-a2a9-e3ff60e05f54', '30417b6c-d7b2-5664-9c49-f644680b44de', 'Ordem de serviço / chamado', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '30417b6c-d7b2-5664-9c49-f644680b44de')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '30417b6c-d7b2-5664-9c49-f644680b44de' AND y.nome_documento = 'Ordem de serviço / chamado');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5a7304fa-b3aa-5b83-b5c3-2a7605466549', '016c58e6-1945-5263-a62c-b092f44e2843', 2, 'SEMAD/TI', 'Suporte técnico', 'Diagnosticar e resolver ou encaminhar à manutenção externa.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '016c58e6-1945-5263-a62c-b092f44e2843')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '016c58e6-1945-5263-a62c-b092f44e2843' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f30129cb-0620-59cb-99ea-4d94c71aecb4', '5a7304fa-b3aa-5b83-b5c3-2a7605466549', 'Laudo técnico', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5a7304fa-b3aa-5b83-b5c3-2a7605466549')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5a7304fa-b3aa-5b83-b5c3-2a7605466549' AND y.nome_documento = 'Laudo técnico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd55772d5-0a8a-592b-8e56-b49b034b0be8', '016c58e6-1945-5263-a62c-b092f44e2843', 3, 'SECRETARIA', 'Unidade solicitante', 'Confirmar a solução e encerrar o chamado.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '016c58e6-1945-5263-a62c-b092f44e2843')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '016c58e6-1945-5263-a62c-b092f44e2843' AND x.ordem = 3);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2573272c-f71e-5047-a213-e9ca24d46799', '30417b6c-d7b2-5664-9c49-f644680b44de', '5a7304fa-b3aa-5b83-b5c3-2a7605466549', 'Chamado aberto', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '30417b6c-d7b2-5664-9c49-f644680b44de')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '5a7304fa-b3aa-5b83-b5c3-2a7605466549')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '30417b6c-d7b2-5664-9c49-f644680b44de' AND z.destino_etapa_id = '5a7304fa-b3aa-5b83-b5c3-2a7605466549' AND z.condicao_transicao = 'Chamado aberto');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c75fb216-327a-5ab9-9e35-420983eb4396', '5a7304fa-b3aa-5b83-b5c3-2a7605466549', 'd55772d5-0a8a-592b-8e56-b49b034b0be8', 'Atendimento concluído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5a7304fa-b3aa-5b83-b5c3-2a7605466549')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd55772d5-0a8a-592b-8e56-b49b034b0be8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5a7304fa-b3aa-5b83-b5c3-2a7605466549' AND z.destino_etapa_id = 'd55772d5-0a8a-592b-8e56-b49b034b0be8' AND z.condicao_transicao = 'Atendimento concluído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '54ee1412-089f-51db-8eee-4f330f14865a', 'RH.CAP.008', 'Inscrição em curso de capacitação', 'Inscrever o servidor em curso promovido pela Prefeitura ou por outro órgão, com autorização da chefia.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.05.01.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.05.01.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.CAP.008');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '6bbbdea6-8398-5479-8626-f428dba7080c', '54ee1412-089f-51db-8eee-4f330f14865a', 1, 'SECRETARIA', 'Chefia imediata', 'Autorizar a participação e a dispensa do expediente.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '54ee1412-089f-51db-8eee-4f330f14865a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '54ee1412-089f-51db-8eee-4f330f14865a' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9fb7b879-5976-5d8d-a57e-4e9cea534a8d', '6bbbdea6-8398-5479-8626-f428dba7080c', 'Ficha de inscrição', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6bbbdea6-8398-5479-8626-f428dba7080c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6bbbdea6-8398-5479-8626-f428dba7080c' AND y.nome_documento = 'Ficha de inscrição');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '800aa228-9b94-5310-bf1e-d11d22424fe6', '54ee1412-089f-51db-8eee-4f330f14865a', 2, 'SEGEP/CAPACITACAO', 'Capacitação', 'Confirmar a vaga, registrar a frequência e emitir o certificado.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '54ee1412-089f-51db-8eee-4f330f14865a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '54ee1412-089f-51db-8eee-4f330f14865a' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '95506fa0-d876-551b-9b4e-d35b7b41eb86', '800aa228-9b94-5310-bf1e-d11d22424fe6', 'Lista de presença', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '800aa228-9b94-5310-bf1e-d11d22424fe6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '800aa228-9b94-5310-bf1e-d11d22424fe6' AND y.nome_documento = 'Lista de presença');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '900598bb-ed54-5f9a-b451-902eb8deb80d', '800aa228-9b94-5310-bf1e-d11d22424fe6', 'Certificado', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '800aa228-9b94-5310-bf1e-d11d22424fe6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '800aa228-9b94-5310-bf1e-d11d22424fe6' AND y.nome_documento = 'Certificado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '7f09bd09-e006-520f-a360-697eae8b737a', '6bbbdea6-8398-5479-8626-f428dba7080c', '800aa228-9b94-5310-bf1e-d11d22424fe6', 'Chefia de acordo', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '6bbbdea6-8398-5479-8626-f428dba7080c')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '800aa228-9b94-5310-bf1e-d11d22424fe6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '6bbbdea6-8398-5479-8626-f428dba7080c' AND z.destino_etapa_id = '800aa228-9b94-5310-bf1e-d11d22424fe6' AND z.condicao_transicao = 'Chefia de acordo');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '9363b425-71eb-5c47-89cc-dcdc1baf7524', 'RH.FOL.009', 'Folha de pagamento complementar', 'Pagar valores devidos fora da folha normal (diferenças, verbas atrasadas).', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.00.06'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.00.06' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.FOL.009');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '347d997a-f38e-5996-9240-a69e9c58792f', '9363b425-71eb-5c47-89cc-dcdc1baf7524', 1, 'SEGEP/RH', 'Recursos Humanos', 'Apurar os valores devidos e montar o demonstrativo.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9363b425-71eb-5c47-89cc-dcdc1baf7524')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9363b425-71eb-5c47-89cc-dcdc1baf7524' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8fc0c5ac-c29c-53db-a73b-30c42d51c7fc', '347d997a-f38e-5996-9240-a69e9c58792f', 'Demonstrativo de cálculo', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '347d997a-f38e-5996-9240-a69e9c58792f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '347d997a-f38e-5996-9240-a69e9c58792f' AND y.nome_documento = 'Demonstrativo de cálculo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'dfd26a61-18af-5760-8573-de88a0be866f', '9363b425-71eb-5c47-89cc-dcdc1baf7524', 2, 'SEFIN/CONTABILIDADE', 'Contabilidade', 'Verificar a dotação e empenhar.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9363b425-71eb-5c47-89cc-dcdc1baf7524')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9363b425-71eb-5c47-89cc-dcdc1baf7524' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '39d071cb-8a3d-54d6-942d-4363395405fb', 'dfd26a61-18af-5760-8573-de88a0be866f', 'Nota de empenho', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'dfd26a61-18af-5760-8573-de88a0be866f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'dfd26a61-18af-5760-8573-de88a0be866f' AND y.nome_documento = 'Nota de empenho');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '180c1605-278d-56bf-ab61-18a37922b021', '9363b425-71eb-5c47-89cc-dcdc1baf7524', 3, 'SEGEP/FOLHA', 'Folha de pagamento', 'Processar a folha complementar e o pagamento.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9363b425-71eb-5c47-89cc-dcdc1baf7524')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9363b425-71eb-5c47-89cc-dcdc1baf7524' AND x.ordem = 3);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '501bf9b5-bf30-5b22-9f03-515b9b3b2621', '347d997a-f38e-5996-9240-a69e9c58792f', 'dfd26a61-18af-5760-8573-de88a0be866f', 'Valores apurados', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '347d997a-f38e-5996-9240-a69e9c58792f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'dfd26a61-18af-5760-8573-de88a0be866f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '347d997a-f38e-5996-9240-a69e9c58792f' AND z.destino_etapa_id = 'dfd26a61-18af-5760-8573-de88a0be866f' AND z.condicao_transicao = 'Valores apurados');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '03cac2b4-4586-5dd9-8343-1fde91c307c6', 'dfd26a61-18af-5760-8573-de88a0be866f', '180c1605-278d-56bf-ab61-18a37922b021', 'Despesa empenhada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'dfd26a61-18af-5760-8573-de88a0be866f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '180c1605-278d-56bf-ab61-18a37922b021')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'dfd26a61-18af-5760-8573-de88a0be866f' AND z.destino_etapa_id = '180c1605-278d-56bf-ab61-18a37922b021' AND z.condicao_transicao = 'Despesa empenhada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '924a0b32-c6ec-5f6d-8d34-a90e0166ae15', 'RH.RES.010', 'Rescisão e verbas rescisórias', 'Calcular e pagar as verbas devidas na saída do servidor (exoneração, término de contrato, falecimento).', 'Ex-servidores e dependentes', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal)', ''), '2.0.06.00.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.00.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.RES.010');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a02daf0e-8e51-5e5f-9c34-3f50b637be94', '924a0b32-c6ec-5f6d-8d34-a90e0166ae15', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber o ato de desligamento e apurar férias, 13º e demais verbas.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '924a0b32-c6ec-5f6d-8d34-a90e0166ae15')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '924a0b32-c6ec-5f6d-8d34-a90e0166ae15' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '86868360-24e8-53e6-aea9-f96e3366525f', 'a02daf0e-8e51-5e5f-9c34-3f50b637be94', 'Ato de desligamento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a02daf0e-8e51-5e5f-9c34-3f50b637be94')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a02daf0e-8e51-5e5f-9c34-3f50b637be94' AND y.nome_documento = 'Ato de desligamento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '84112b89-671b-532f-9901-a2d716f1fa16', 'a02daf0e-8e51-5e5f-9c34-3f50b637be94', 'Termo de rescisão', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a02daf0e-8e51-5e5f-9c34-3f50b637be94')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a02daf0e-8e51-5e5f-9c34-3f50b637be94' AND y.nome_documento = 'Termo de rescisão');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60', '924a0b32-c6ec-5f6d-8d34-a90e0166ae15', 2, 'SEFIN/CONTABILIDADE', 'Contabilidade', 'Empenhar e pagar as verbas rescisórias.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '924a0b32-c6ec-5f6d-8d34-a90e0166ae15')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '924a0b32-c6ec-5f6d-8d34-a90e0166ae15' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '09b31566-89d6-50cc-9ebc-f82a45955072', 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60', 'Nota de empenho', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60' AND y.nome_documento = 'Nota de empenho');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b472a0a9-5f78-53e5-8731-229084815c68', 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60', 'Ordem de pagamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60' AND y.nome_documento = 'Ordem de pagamento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e37aec7d-c141-5386-86b1-024759895c20', 'a02daf0e-8e51-5e5f-9c34-3f50b637be94', 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60', 'Verbas apuradas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a02daf0e-8e51-5e5f-9c34-3f50b637be94')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a02daf0e-8e51-5e5f-9c34-3f50b637be94' AND z.destino_etapa_id = 'ec1fabcf-2bdb-5e2c-9750-9df4add3fd60' AND z.condicao_transicao = 'Verbas apuradas');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1ddc0208-fe91-5423-bb4e-505e9844584f', 'RH.VTR.011', 'Vale-transporte (adesão e desistência)', 'Conceder ou cancelar o vale-transporte, com desconto na folha conforme a lei.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.01.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.01.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.VTR.011');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a', '1ddc0208-fe91-5423-bb4e-505e9844584f', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber a declaração com o itinerário e os meios de transporte.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1ddc0208-fe91-5423-bb4e-505e9844584f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1ddc0208-fe91-5423-bb4e-505e9844584f' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f5377bc3-8bff-5750-aa85-1822abe43f13', 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a', 'Requerimento de vale-transporte', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a' AND y.nome_documento = 'Requerimento de vale-transporte');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bd6e17e5-d2de-55a0-ab24-14bf8bd1ecc1', 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a', 'Declaração de itinerário', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a' AND y.nome_documento = 'Declaração de itinerário');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '05207831-7a70-5fbe-b014-779c2e3e58b3', '1ddc0208-fe91-5423-bb4e-505e9844584f', 2, 'SEGEP/FOLHA', 'Folha de pagamento', 'Lançar o benefício e o desconto correspondente.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1ddc0208-fe91-5423-bb4e-505e9844584f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1ddc0208-fe91-5423-bb4e-505e9844584f' AND x.ordem = 2);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'de2c3957-3a33-5f85-be8a-4b41e5fdb52c', 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a', '05207831-7a70-5fbe-b014-779c2e3e58b3', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '05207831-7a70-5fbe-b014-779c2e3e58b3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'cd8fa45d-845b-50f8-894f-156cd6f9fb6a' AND z.destino_etapa_id = '05207831-7a70-5fbe-b014-779c2e3e58b3' AND z.condicao_transicao = 'Requerimento instruído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'e77c9cc7-511e-58f6-91bb-adae154ca186', 'RH.CNS.012', 'Empréstimo consignado', 'Averbar o empréstimo consignado do servidor dentro da margem legal, com banco credenciado.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal)', ''), '2.0.06.01.16'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.01.16' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.CNS.012');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8f6232b5-705e-5758-9f48-0eccc5c69367', 'e77c9cc7-511e-58f6-91bb-adae154ca186', 1, 'SEGEP/RH', 'Recursos Humanos', 'Emitir a declaração de margem consignável.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e77c9cc7-511e-58f6-91bb-adae154ca186')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e77c9cc7-511e-58f6-91bb-adae154ca186' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cb38d904-ecee-5512-904c-977066f202af', '8f6232b5-705e-5758-9f48-0eccc5c69367', 'Requerimento de adesão ao consignado', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8f6232b5-705e-5758-9f48-0eccc5c69367')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8f6232b5-705e-5758-9f48-0eccc5c69367' AND y.nome_documento = 'Requerimento de adesão ao consignado');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '91b1cbf8-4f09-5276-aacd-f3dc52b60160', '8f6232b5-705e-5758-9f48-0eccc5c69367', 'Declaração de margem consignável', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8f6232b5-705e-5758-9f48-0eccc5c69367')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8f6232b5-705e-5758-9f48-0eccc5c69367' AND y.nome_documento = 'Declaração de margem consignável');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '1aeba5fe-1075-58d3-8745-1eb378a9eec3', 'e77c9cc7-511e-58f6-91bb-adae154ca186', 2, 'SEGEP/FOLHA', 'Folha de pagamento', 'Averbar o contrato e lançar o desconto em folha.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'e77c9cc7-511e-58f6-91bb-adae154ca186')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'e77c9cc7-511e-58f6-91bb-adae154ca186' AND x.ordem = 2);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '18641d27-82dc-5d70-ac15-446fb02fbfb7', '8f6232b5-705e-5758-9f48-0eccc5c69367', '1aeba5fe-1075-58d3-8745-1eb378a9eec3', 'Margem disponível', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8f6232b5-705e-5758-9f48-0eccc5c69367')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '1aeba5fe-1075-58d3-8745-1eb378a9eec3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8f6232b5-705e-5758-9f48-0eccc5c69367' AND z.destino_etapa_id = '1aeba5fe-1075-58d3-8745-1eb378a9eec3' AND z.condicao_transicao = 'Margem disponível');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '226e52bd-adf2-523e-bd67-72413496dd06', 'RH.PAL.013', 'Desconto de pensão alimentícia', 'Cumprir a determinação judicial de desconto de pensão alimentícia na folha do servidor.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 e segredo de justiça (CPC, art. 189)', ''), '2.0.06.01.11'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.01.11' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.PAL.013');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'bfe050ad-1aae-5dd0-938b-64451bde79a0', '226e52bd-adf2-523e-bd67-72413496dd06', 1, 'PGM', 'Procuradoria', 'Receber o ofício judicial e orientar o cumprimento.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '226e52bd-adf2-523e-bd67-72413496dd06')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '226e52bd-adf2-523e-bd67-72413496dd06' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a3ebc1b5-949f-56a4-8c48-fe069fadc202', 'bfe050ad-1aae-5dd0-938b-64451bde79a0', 'Ofício judicial', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bfe050ad-1aae-5dd0-938b-64451bde79a0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bfe050ad-1aae-5dd0-938b-64451bde79a0' AND y.nome_documento = 'Ofício judicial');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '90d2f152-2876-59a7-b931-551e7b685fe7', '226e52bd-adf2-523e-bd67-72413496dd06', 2, 'SEGEP/FOLHA', 'Folha de pagamento', 'Lançar o desconto e repassar ao beneficiário conforme a decisão.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '226e52bd-adf2-523e-bd67-72413496dd06')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '226e52bd-adf2-523e-bd67-72413496dd06' AND x.ordem = 2);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '204befe2-cdcd-597a-9904-3693b48b25cd', 'bfe050ad-1aae-5dd0-938b-64451bde79a0', '90d2f152-2876-59a7-b931-551e7b685fe7', 'Ordem conferida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'bfe050ad-1aae-5dd0-938b-64451bde79a0')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '90d2f152-2876-59a7-b931-551e7b685fe7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'bfe050ad-1aae-5dd0-938b-64451bde79a0' AND z.destino_etapa_id = '90d2f152-2876-59a7-b931-551e7b685fe7' AND z.condicao_transicao = 'Ordem conferida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4', 'RH.REE.014', 'Reembolso de despesas e indenização ao servidor', 'Ressarcir o servidor por despesa feita a serviço ou por indenização prevista no estatuto.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.01.14'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.01.14' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.REE.014');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '065dbaed-7552-5469-b1b5-9a4a6d56b961', 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4', 1, 'SECRETARIA', 'Chefia imediata', 'Confirmar que a despesa foi feita a serviço.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '36c7983f-159c-58a4-bdd4-30fa5ae521a3', '065dbaed-7552-5469-b1b5-9a4a6d56b961', 'Requerimento de reembolso', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '065dbaed-7552-5469-b1b5-9a4a6d56b961')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '065dbaed-7552-5469-b1b5-9a4a6d56b961' AND y.nome_documento = 'Requerimento de reembolso');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ae04f4f1-2f0a-5c41-af47-74fe66122d2e', '065dbaed-7552-5469-b1b5-9a4a6d56b961', 'Comprovantes de despesa', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '065dbaed-7552-5469-b1b5-9a4a6d56b961')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '065dbaed-7552-5469-b1b5-9a4a6d56b961' AND y.nome_documento = 'Comprovantes de despesa');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3f393d5d-5c9c-5b2b-a961-bd58771e21e6', 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4', 2, 'SEGEP/RH', 'Recursos Humanos', 'Conferir a previsão legal e o valor.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '711c747d-cbd3-5302-9961-4c21fa634ace', '3f393d5d-5c9c-5b2b-a961-bd58771e21e6', 'Informação do RH', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3f393d5d-5c9c-5b2b-a961-bd58771e21e6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3f393d5d-5c9c-5b2b-a961-bd58771e21e6' AND y.nome_documento = 'Informação do RH');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3e36cac4-bdb1-54ad-ab35-4b45ad387822', 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4', 3, 'SEGEP/FOLHA', 'Folha de pagamento', 'Pagar o reembolso na folha ou por ordem de pagamento.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b07cd993-a0dd-59e1-a25e-bb4855c5ffc4' AND x.ordem = 3);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c4b825b3-1f8b-53bc-a0d0-e58fe20e3aa0', '065dbaed-7552-5469-b1b5-9a4a6d56b961', '3f393d5d-5c9c-5b2b-a961-bd58771e21e6', 'Chefia de acordo', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '065dbaed-7552-5469-b1b5-9a4a6d56b961')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3f393d5d-5c9c-5b2b-a961-bd58771e21e6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '065dbaed-7552-5469-b1b5-9a4a6d56b961' AND z.destino_etapa_id = '3f393d5d-5c9c-5b2b-a961-bd58771e21e6' AND z.condicao_transicao = 'Chefia de acordo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0162f3f8-d473-5444-a8b3-a70ca4a0b636', '3f393d5d-5c9c-5b2b-a961-bd58771e21e6', '3e36cac4-bdb1-54ad-ab35-4b45ad387822', 'Reembolso devido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '3f393d5d-5c9c-5b2b-a961-bd58771e21e6')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3e36cac4-bdb1-54ad-ab35-4b45ad387822')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '3f393d5d-5c9c-5b2b-a961-bd58771e21e6' AND z.destino_etapa_id = '3e36cac4-bdb1-54ad-ab35-4b45ad387822' AND z.condicao_transicao = 'Reembolso devido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '61f5c5b4-cb3d-591b-9c12-975c5677fefa', 'RH.INS.015', 'Declaração de tempo de contribuição ao INSS', 'Emitir a declaração do tempo de contribuição ao Regime Geral (servidor celetista, contratado ou comissionado).', 'Servidores e ex-servidores', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.06.04.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.06.04.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.INS.015');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '800c44f3-92d7-51f0-b116-14cf0d04c16b', '61f5c5b4-cb3d-591b-9c12-975c5677fefa', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber o pedido e levantar as contribuições do período.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '61f5c5b4-cb3d-591b-9c12-975c5677fefa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '61f5c5b4-cb3d-591b-9c12-975c5677fefa' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e7ceced3-8687-5ca9-86d6-8c29470e0b70', '800c44f3-92d7-51f0-b116-14cf0d04c16b', 'Requerimento de tempo de contribuição', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '800c44f3-92d7-51f0-b116-14cf0d04c16b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '800c44f3-92d7-51f0-b116-14cf0d04c16b' AND y.nome_documento = 'Requerimento de tempo de contribuição');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f', '61f5c5b4-cb3d-591b-9c12-975c5677fefa', 2, 'SEGEP/RH', 'Recursos Humanos', 'Emitir a declaração.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '61f5c5b4-cb3d-591b-9c12-975c5677fefa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '61f5c5b4-cb3d-591b-9c12-975c5677fefa' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'dd6c2a6f-751e-52e2-9f58-bd28bfa236b8', '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f', 'Declaração de tempo de contribuição', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f' AND y.nome_documento = 'Declaração de tempo de contribuição');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '22a8b818-2c89-5c97-9417-444e957100fe', '800c44f3-92d7-51f0-b116-14cf0d04c16b', '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f', 'Período levantado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '800c44f3-92d7-51f0-b116-14cf0d04c16b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '800c44f3-92d7-51f0-b116-14cf0d04c16b' AND z.destino_etapa_id = '42f346cc-e6fc-57ba-ac5c-7f86637e0f8f' AND z.condicao_transicao = 'Período levantado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'f638bd2e-3677-53c2-bd61-981bcfcb57e9', 'RH.LMA.016', 'Licença-maternidade', 'Conceder a licença à gestante ou adotante, com a prorrogação prevista na legislação municipal.', 'Servidoras municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal)', ''), '2.0.07.00.16'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.00.16' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.LMA.016');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69', 'f638bd2e-3677-53c2-bd61-981bcfcb57e9', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber o atestado ou a certidão de nascimento ou o termo de adoção.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f638bd2e-3677-53c2-bd61-981bcfcb57e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f638bd2e-3677-53c2-bd61-981bcfcb57e9' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4c8c24eb-374c-5c0d-bda2-d6871fdd6168', '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69', 'Requerimento de licença-maternidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69' AND y.nome_documento = 'Requerimento de licença-maternidade');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '403ff06e-bfee-528e-9bfd-bb62e9bc714a', '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69', 'Certidão de nascimento ou termo de guarda', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69' AND y.nome_documento = 'Certidão de nascimento ou termo de guarda');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '7323779f-3490-5f30-a65c-3fea51ecd619', 'f638bd2e-3677-53c2-bd61-981bcfcb57e9', 2, 'SEGEP/DESOPEM', 'Saúde ocupacional (DESOPEM)', 'Homologar a licença e, se pedida, a prorrogação.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f638bd2e-3677-53c2-bd61-981bcfcb57e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f638bd2e-3677-53c2-bd61-981bcfcb57e9' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b48e8de7-c80f-5b1f-91a9-364d98298f95', '7323779f-3490-5f30-a65c-3fea51ecd619', 'Requerimento de prorrogação', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7323779f-3490-5f30-a65c-3fea51ecd619')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7323779f-3490-5f30-a65c-3fea51ecd619' AND y.nome_documento = 'Requerimento de prorrogação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '45791844-8208-58cd-81ce-73f5cc62826f', 'f638bd2e-3677-53c2-bd61-981bcfcb57e9', 3, 'SEGEP/RH', 'Recursos Humanos', 'Publicar o ato e registrar na pasta funcional.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f638bd2e-3677-53c2-bd61-981bcfcb57e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f638bd2e-3677-53c2-bd61-981bcfcb57e9' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd6679488-eeac-5ab6-9f95-88238fe5301e', '45791844-8208-58cd-81ce-73f5cc62826f', 'Portaria de licença-maternidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '45791844-8208-58cd-81ce-73f5cc62826f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '45791844-8208-58cd-81ce-73f5cc62826f' AND y.nome_documento = 'Portaria de licença-maternidade');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5eef14af-2475-5a05-8bab-624a9cca92e0', '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69', '7323779f-3490-5f30-a65c-3fea51ecd619', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '7323779f-3490-5f30-a65c-3fea51ecd619')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '84ac2a39-1a36-5aa7-94b8-17bd1e50ac69' AND z.destino_etapa_id = '7323779f-3490-5f30-a65c-3fea51ecd619' AND z.condicao_transicao = 'Requerimento instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '57c1c82a-e776-5c6b-b804-5da630f33cde', '7323779f-3490-5f30-a65c-3fea51ecd619', '45791844-8208-58cd-81ce-73f5cc62826f', 'Licença homologada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '7323779f-3490-5f30-a65c-3fea51ecd619')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '45791844-8208-58cd-81ce-73f5cc62826f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '7323779f-3490-5f30-a65c-3fea51ecd619' AND z.destino_etapa_id = '45791844-8208-58cd-81ce-73f5cc62826f' AND z.condicao_transicao = 'Licença homologada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c', 'RH.LPA.017', 'Licença-paternidade', 'Conceder a licença-paternidade ao servidor pelo nascimento ou adoção.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.07.00.17'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.00.17' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.LPA.017');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d', 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber a certidão de nascimento ou o termo de adoção.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6c8b8b9a-41a6-5d6e-b30a-59fd5a495120', '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d', 'Requerimento de licença-paternidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d' AND y.nome_documento = 'Requerimento de licença-paternidade');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '95bbb754-2d65-59f5-86d7-b81aa0609a2f', '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d', 'Certidão de nascimento ou termo de guarda', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d' AND y.nome_documento = 'Certidão de nascimento ou termo de guarda');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '81d60084-5d68-54ec-9c47-39c8b6a2b558', 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c', 2, 'SEGEP/RH', 'Recursos Humanos', 'Publicar o ato e registrar na pasta funcional.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ac3ac46d-3395-5ac7-9b1f-e3b2eca01e5c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '11143f48-7ee8-57f6-9038-cc4f284c78dd', '81d60084-5d68-54ec-9c47-39c8b6a2b558', 'Registro da licença-paternidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '81d60084-5d68-54ec-9c47-39c8b6a2b558')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '81d60084-5d68-54ec-9c47-39c8b6a2b558' AND y.nome_documento = 'Registro da licença-paternidade');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e72dd78f-2ff1-590c-9a06-d6f94823caea', '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d', '81d60084-5d68-54ec-9c47-39c8b6a2b558', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '81d60084-5d68-54ec-9c47-39c8b6a2b558')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '5081b1d6-62d6-5a4a-8e7e-af7e87ad121d' AND z.destino_etapa_id = '81d60084-5d68-54ec-9c47-39c8b6a2b558' AND z.condicao_transicao = 'Requerimento instruído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '19ca98db-3d75-5a16-bf2e-1efce7d57af5', 'RH.LNO.018', 'Licença por falecimento (nojo)', 'Conceder a licença pelo falecimento de familiar previsto no estatuto.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.07.00.18'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.00.18' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.LNO.018');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0b5df729-c480-5c31-864c-81e42a014749', '19ca98db-3d75-5a16-bf2e-1efce7d57af5', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber a certidão de óbito e o vínculo de parentesco.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '19ca98db-3d75-5a16-bf2e-1efce7d57af5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '19ca98db-3d75-5a16-bf2e-1efce7d57af5' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'df4d11ae-ff93-5fb1-961e-eefa2408dac1', '0b5df729-c480-5c31-864c-81e42a014749', 'Requerimento de licença', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0b5df729-c480-5c31-864c-81e42a014749')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0b5df729-c480-5c31-864c-81e42a014749' AND y.nome_documento = 'Requerimento de licença');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '14152682-b8f8-5ffb-a879-3599837055b0', '0b5df729-c480-5c31-864c-81e42a014749', 'Certidão de óbito', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0b5df729-c480-5c31-864c-81e42a014749')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0b5df729-c480-5c31-864c-81e42a014749' AND y.nome_documento = 'Certidão de óbito');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1', '19ca98db-3d75-5a16-bf2e-1efce7d57af5', 2, 'SEGEP/RH', 'Recursos Humanos', 'Publicar o ato e registrar na pasta funcional.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '19ca98db-3d75-5a16-bf2e-1efce7d57af5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '19ca98db-3d75-5a16-bf2e-1efce7d57af5' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '165dbeb9-7069-53bb-a213-9afba6df62c4', 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1', 'Registro da licença', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1' AND y.nome_documento = 'Registro da licença');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '074032a8-d361-5c2d-be5c-28a0d9b34390', '0b5df729-c480-5c31-864c-81e42a014749', 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0b5df729-c480-5c31-864c-81e42a014749')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0b5df729-c480-5c31-864c-81e42a014749' AND z.destino_etapa_id = 'd1cc4c75-50c9-5702-97c0-5c196e78c8a1' AND z.condicao_transicao = 'Requerimento instruído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e', 'RH.PER.019', 'Perícia médica para licença de saúde', 'Avaliar o servidor em perícia médica e conceder a licença para tratamento de saúde.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '2.0.07.02.30'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.02.30' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.PER.019');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '17748285-7383-5ff5-bff3-05d17b0f20d8', 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e', 1, 'SEGEP/DESOPEM', 'Saúde ocupacional (DESOPEM)', 'Agendar a perícia (no local, domiciliar ou hospitalar) a partir do atestado.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'caed7bbc-563a-5052-a805-26f29e491e2f', '17748285-7383-5ff5-bff3-05d17b0f20d8', 'Requerimento de perícia médica', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '17748285-7383-5ff5-bff3-05d17b0f20d8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '17748285-7383-5ff5-bff3-05d17b0f20d8' AND y.nome_documento = 'Requerimento de perícia médica');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4a4ace6f-947e-51fe-871d-9a6213771ea1', '17748285-7383-5ff5-bff3-05d17b0f20d8', 'Atestado médico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '17748285-7383-5ff5-bff3-05d17b0f20d8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '17748285-7383-5ff5-bff3-05d17b0f20d8' AND y.nome_documento = 'Atestado médico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e6477174-b35a-5a87-810e-8a4f16a93a1f', 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e', 2, 'SEGEP/DESOPEM', 'Junta médica', 'Realizar a perícia e emitir o laudo com o prazo da licença.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '642686a4-cadc-5181-98cd-dca402e66ff5', 'e6477174-b35a-5a87-810e-8a4f16a93a1f', 'Laudo pericial', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e6477174-b35a-5a87-810e-8a4f16a93a1f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e6477174-b35a-5a87-810e-8a4f16a93a1f' AND y.nome_documento = 'Laudo pericial');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77', 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e', 3, 'SEGEP/RH', 'Recursos Humanos', 'Publicar o ato e registrar na pasta funcional.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'eaec699e-1b31-5cd7-958d-3c0d8eecef4e' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e2201c39-d270-5139-a1a1-75568aab49fd', 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77', 'Portaria de licença para tratamento de saúde', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77' AND y.nome_documento = 'Portaria de licença para tratamento de saúde');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b91c06fe-e861-5a7d-a2a8-d05fc8eb99fa', '17748285-7383-5ff5-bff3-05d17b0f20d8', 'e6477174-b35a-5a87-810e-8a4f16a93a1f', 'Perícia agendada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '17748285-7383-5ff5-bff3-05d17b0f20d8')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e6477174-b35a-5a87-810e-8a4f16a93a1f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '17748285-7383-5ff5-bff3-05d17b0f20d8' AND z.destino_etapa_id = 'e6477174-b35a-5a87-810e-8a4f16a93a1f' AND z.condicao_transicao = 'Perícia agendada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'ab7fe6b3-715c-51fc-bad4-5210dcf6d536', 'e6477174-b35a-5a87-810e-8a4f16a93a1f', 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77', 'Laudo emitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e6477174-b35a-5a87-810e-8a4f16a93a1f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e6477174-b35a-5a87-810e-8a4f16a93a1f' AND z.destino_etapa_id = 'ca407c39-3d4d-5215-a296-d7e4f5ca2d77' AND z.condicao_transicao = 'Laudo emitido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '053b47d6-4cc7-5a0d-a186-d113ab4ae154', 'RH.CAT.020', 'Comunicação de Acidente de Trabalho (CAT)', 'Registrar o acidente de trabalho do servidor e encaminhar a avaliação e os benefícios.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '2.0.07.02.10'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.02.10' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.CAT.020');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '03f8dab2-5579-5b70-a4c2-9a9e113121d8', '053b47d6-4cc7-5a0d-a186-d113ab4ae154', 1, 'SECRETARIA', 'Chefia imediata', 'Comunicar o acidente em até 1 dia útil, com testemunhas.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '053b47d6-4cc7-5a0d-a186-d113ab4ae154')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '053b47d6-4cc7-5a0d-a186-d113ab4ae154' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '1c30048d-a418-5c38-a00a-f65cecfc9350', '03f8dab2-5579-5b70-a4c2-9a9e113121d8', 'Relato do acidente', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '03f8dab2-5579-5b70-a4c2-9a9e113121d8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '03f8dab2-5579-5b70-a4c2-9a9e113121d8' AND y.nome_documento = 'Relato do acidente');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '776d1ef9-0dd5-52d4-a208-b34161fa8efb', '053b47d6-4cc7-5a0d-a186-d113ab4ae154', 2, 'SEGEP/DESOPEM', 'Saúde ocupacional (DESOPEM)', 'Emitir a CAT e avaliar o afastamento.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '053b47d6-4cc7-5a0d-a186-d113ab4ae154')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '053b47d6-4cc7-5a0d-a186-d113ab4ae154' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9a865679-696e-5121-913b-adb9d796a0f3', '776d1ef9-0dd5-52d4-a208-b34161fa8efb', 'CAT – Comunicação de Acidente de Trabalho', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '776d1ef9-0dd5-52d4-a208-b34161fa8efb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '776d1ef9-0dd5-52d4-a208-b34161fa8efb' AND y.nome_documento = 'CAT – Comunicação de Acidente de Trabalho');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '631f6085-8367-5759-9a8b-377dc7c3180b', '776d1ef9-0dd5-52d4-a208-b34161fa8efb', 'Laudo médico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '776d1ef9-0dd5-52d4-a208-b34161fa8efb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '776d1ef9-0dd5-52d4-a208-b34161fa8efb' AND y.nome_documento = 'Laudo médico');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'aa6fd17f-b4c1-57b5-b5b8-e0abb0b510e7', '03f8dab2-5579-5b70-a4c2-9a9e113121d8', '776d1ef9-0dd5-52d4-a208-b34161fa8efb', 'Acidente comunicado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '03f8dab2-5579-5b70-a4c2-9a9e113121d8')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '776d1ef9-0dd5-52d4-a208-b34161fa8efb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '03f8dab2-5579-5b70-a4c2-9a9e113121d8' AND z.destino_etapa_id = '776d1ef9-0dd5-52d4-a208-b34161fa8efb' AND z.condicao_transicao = 'Acidente comunicado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'dc68475b-85e9-588e-a9d5-b005f880e3ef', 'RH.REA.021', 'Readaptação de função', 'Readaptar o servidor com limitação de saúde a atribuições compatíveis, após avaliação da junta médica.', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '2.0.07.02.28'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.02.28' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.REA.021');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83', 'dc68475b-85e9-588e-a9d5-b005f880e3ef', 1, 'SEGEP/DESOPEM', 'Junta médica', 'Avaliar a limitação e indicar as restrições.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'dc68475b-85e9-588e-a9d5-b005f880e3ef')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'dc68475b-85e9-588e-a9d5-b005f880e3ef' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b7395632-6f54-571c-a41b-a76b743ffd44', 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83', 'Laudo de readaptação', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83' AND y.nome_documento = 'Laudo de readaptação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '65bf0446-f7ed-53ed-8b1f-2ad23ee68348', 'dc68475b-85e9-588e-a9d5-b005f880e3ef', 2, 'SEGEP/RH', 'Recursos Humanos', 'Indicar a função compatível com a secretaria de lotação.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'dc68475b-85e9-588e-a9d5-b005f880e3ef')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'dc68475b-85e9-588e-a9d5-b005f880e3ef' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cd2440f9-90e1-5a3f-9386-1f6b198c5f99', '65bf0446-f7ed-53ed-8b1f-2ad23ee68348', 'Proposta de readaptação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '65bf0446-f7ed-53ed-8b1f-2ad23ee68348')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '65bf0446-f7ed-53ed-8b1f-2ad23ee68348' AND y.nome_documento = 'Proposta de readaptação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed', 'dc68475b-85e9-588e-a9d5-b005f880e3ef', 3, 'SEGEP/RH', 'Recursos Humanos', 'Publicar o ato e registrar na pasta funcional.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'dc68475b-85e9-588e-a9d5-b005f880e3ef')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'dc68475b-85e9-588e-a9d5-b005f880e3ef' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2415c5e8-a259-58e9-a89a-27933bd75ac2', 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed', 'Portaria de readaptação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed' AND y.nome_documento = 'Portaria de readaptação');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2742fd85-dffa-5d35-a0f1-2f9dbe38d3af', 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83', '65bf0446-f7ed-53ed-8b1f-2ad23ee68348', 'Restrições definidas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '65bf0446-f7ed-53ed-8b1f-2ad23ee68348')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ba3bb64e-1b2f-5dd9-ac59-c1ce930cda83' AND z.destino_etapa_id = '65bf0446-f7ed-53ed-8b1f-2ad23ee68348' AND z.condicao_transicao = 'Restrições definidas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f39fc7a6-da8d-52e3-bcb6-b85534d20002', '65bf0446-f7ed-53ed-8b1f-2ad23ee68348', 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed', 'Função definida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '65bf0446-f7ed-53ed-8b1f-2ad23ee68348')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '65bf0446-f7ed-53ed-8b1f-2ad23ee68348' AND z.destino_etapa_id = 'ec6a00c5-e2a7-5bf0-a72e-41bff24ac9ed' AND z.condicao_transicao = 'Função definida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'c25c935b-6dbe-5017-b51a-df8eb40997ed', 'RH.RCH.022', 'Redução de carga horária para cuidar de pessoa com deficiência', 'Reduzir a jornada do servidor responsável por pessoa com deficiência, sem redução de vencimento (Lei 8.563/2015).', 'Servidores municipais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '2.0.07.02.32'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.07.02.32' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.RCH.022');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b393a183-1737-5045-b667-ab69621ca315', 'c25c935b-6dbe-5017-b51a-df8eb40997ed', 1, 'SEGEP/RH', 'Recursos Humanos', 'Receber o pedido com o laudo da pessoa com deficiência e a comprovação da dependência.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c25c935b-6dbe-5017-b51a-df8eb40997ed')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c25c935b-6dbe-5017-b51a-df8eb40997ed' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c924047e-b160-5d47-a517-913cd8c5ef19', 'b393a183-1737-5045-b667-ab69621ca315', 'Requerimento de redução de carga horária', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b393a183-1737-5045-b667-ab69621ca315')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b393a183-1737-5045-b667-ab69621ca315' AND y.nome_documento = 'Requerimento de redução de carga horária');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '58e5594c-2ff0-547e-afe2-2db943b631ce', 'b393a183-1737-5045-b667-ab69621ca315', 'Laudo médico da pessoa com deficiência', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b393a183-1737-5045-b667-ab69621ca315')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b393a183-1737-5045-b667-ab69621ca315' AND y.nome_documento = 'Laudo médico da pessoa com deficiência');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0796f03a-ce25-50cc-ab3a-41eff1c587c4', 'c25c935b-6dbe-5017-b51a-df8eb40997ed', 2, 'SEGEP/DESOPEM', 'Junta médica', 'Avaliar o laudo e opinar.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c25c935b-6dbe-5017-b51a-df8eb40997ed')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c25c935b-6dbe-5017-b51a-df8eb40997ed' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9803ab8e-547d-5ec5-a792-c89551bc2790', '0796f03a-ce25-50cc-ab3a-41eff1c587c4', 'Parecer da junta médica', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0796f03a-ce25-50cc-ab3a-41eff1c587c4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0796f03a-ce25-50cc-ab3a-41eff1c587c4' AND y.nome_documento = 'Parecer da junta médica');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8b3b10ad-0501-55f0-bcd4-500e02238c8a', 'c25c935b-6dbe-5017-b51a-df8eb40997ed', 3, 'SEGEP/RH', 'Recursos Humanos', 'Publicar o ato e registrar na pasta funcional.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c25c935b-6dbe-5017-b51a-df8eb40997ed')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c25c935b-6dbe-5017-b51a-df8eb40997ed' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8fa46c6d-eed9-5e3f-ae44-86540cd5b06a', '8b3b10ad-0501-55f0-bcd4-500e02238c8a', 'Portaria de redução de carga horária', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8b3b10ad-0501-55f0-bcd4-500e02238c8a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8b3b10ad-0501-55f0-bcd4-500e02238c8a' AND y.nome_documento = 'Portaria de redução de carga horária');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '754c91af-e78b-5514-b588-65231212f1ca', 'b393a183-1737-5045-b667-ab69621ca315', '0796f03a-ce25-50cc-ab3a-41eff1c587c4', 'Requerimento instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'b393a183-1737-5045-b667-ab69621ca315')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '0796f03a-ce25-50cc-ab3a-41eff1c587c4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'b393a183-1737-5045-b667-ab69621ca315' AND z.destino_etapa_id = '0796f03a-ce25-50cc-ab3a-41eff1c587c4' AND z.condicao_transicao = 'Requerimento instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '49999f16-02bf-556e-9164-b336a5c24cfc', '0796f03a-ce25-50cc-ab3a-41eff1c587c4', '8b3b10ad-0501-55f0-bcd4-500e02238c8a', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0796f03a-ce25-50cc-ab3a-41eff1c587c4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8b3b10ad-0501-55f0-bcd4-500e02238c8a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0796f03a-ce25-50cc-ab3a-41eff1c587c4' AND z.destino_etapa_id = '8b3b10ad-0501-55f0-bcd4-500e02238c8a' AND z.condicao_transicao = 'Parecer favorável');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '54886463-8d4a-50df-8261-7dcca4b9512b', 'RH.COM.023', 'Nomeação e exoneração de cargo em comissão', 'Nomear ou exonerar servidor de cargo em comissão, com verificação de vedações (nepotismo, ficha limpa).', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.08.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.08.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.COM.023');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0d283e0a-f97c-5996-af79-4771c6f94abc', '54886463-8d4a-50df-8261-7dcca4b9512b', 1, 'SECRETARIA', 'Secretário', 'Indicar o nome para o cargo vago.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '54886463-8d4a-50df-8261-7dcca4b9512b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '54886463-8d4a-50df-8261-7dcca4b9512b' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8392ea86-c75a-54b5-affe-47e9ce263c0f', '0d283e0a-f97c-5996-af79-4771c6f94abc', 'Requerimento de nomeação/exoneração', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0d283e0a-f97c-5996-af79-4771c6f94abc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0d283e0a-f97c-5996-af79-4771c6f94abc' AND y.nome_documento = 'Requerimento de nomeação/exoneração');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0', '54886463-8d4a-50df-8261-7dcca4b9512b', 2, 'SEGEP/RH', 'Recursos Humanos', 'Conferir a vaga, os documentos e as declarações de vedação.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '54886463-8d4a-50df-8261-7dcca4b9512b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '54886463-8d4a-50df-8261-7dcca4b9512b' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '0480c359-0ca9-5fb9-a8db-6174150c4da0', 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0', 'Declaração de não parentesco', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0' AND y.nome_documento = 'Declaração de não parentesco');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'aead1341-7aca-5e79-8569-e783de8fecc2', 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0', 'Certidões negativas', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0' AND y.nome_documento = 'Certidões negativas');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '6e2255a8-e3ce-5c0a-99dc-22a678a49c23', '54886463-8d4a-50df-8261-7dcca4b9512b', 3, 'GAB', 'Prefeito', 'Assinar e publicar o decreto ou a portaria.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '54886463-8d4a-50df-8261-7dcca4b9512b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '54886463-8d4a-50df-8261-7dcca4b9512b' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7654c558-be60-5fb7-a127-1aecb498ead5', '6e2255a8-e3ce-5c0a-99dc-22a678a49c23', 'Ato de nomeação/exoneração', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6e2255a8-e3ce-5c0a-99dc-22a678a49c23')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6e2255a8-e3ce-5c0a-99dc-22a678a49c23' AND y.nome_documento = 'Ato de nomeação/exoneração');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3e5e9162-3b49-5547-911e-55cdb760ebf4', '0d283e0a-f97c-5996-af79-4771c6f94abc', 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0', 'Indicação feita', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0d283e0a-f97c-5996-af79-4771c6f94abc')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0d283e0a-f97c-5996-af79-4771c6f94abc' AND z.destino_etapa_id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0' AND z.condicao_transicao = 'Indicação feita');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0a0ab46d-2f23-5dee-9150-133357372a18', 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0', '6e2255a8-e3ce-5c0a-99dc-22a678a49c23', 'Requisitos atendidos', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '6e2255a8-e3ce-5c0a-99dc-22a678a49c23')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e695e7b0-a1a6-56b3-8308-d3e44531f4c0' AND z.destino_etapa_id = '6e2255a8-e3ce-5c0a-99dc-22a678a49c23' AND z.condicao_transicao = 'Requisitos atendidos');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'c1c39be5-9bac-5277-ae0b-f305957a22c9', 'RH.TMP.024', 'Contratação temporária por excepcional interesse público', 'Contratar pessoal por tempo determinado, nos casos previstos em lei, a partir de processo seletivo ou justificativa.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.08.02.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.08.02.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.TMP.024');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '842990f9-8aa7-573f-8d01-73fd4d104bad', 'c1c39be5-9bac-5277-ae0b-f305957a22c9', 1, 'SECRETARIA', 'Secretaria solicitante', 'Justificar a necessidade temporária e indicar o cadastro de reserva.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c1c39be5-9bac-5277-ae0b-f305957a22c9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c1c39be5-9bac-5277-ae0b-f305957a22c9' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd2f49993-f994-5493-8883-d9b624ecf6e4', '842990f9-8aa7-573f-8d01-73fd4d104bad', 'Requerimento de contratação temporária', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '842990f9-8aa7-573f-8d01-73fd4d104bad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '842990f9-8aa7-573f-8d01-73fd4d104bad' AND y.nome_documento = 'Requerimento de contratação temporária');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'da8e6318-1cec-5123-b820-a29d41c5554c', '842990f9-8aa7-573f-8d01-73fd4d104bad', 'Justificativa de excepcional interesse público', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '842990f9-8aa7-573f-8d01-73fd4d104bad')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '842990f9-8aa7-573f-8d01-73fd4d104bad' AND y.nome_documento = 'Justificativa de excepcional interesse público');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc', 'c1c39be5-9bac-5277-ae0b-f305957a22c9', 2, 'SEGEP/RH', 'Recursos Humanos', 'Conferir a dotação, a ordem de classificação e os documentos do contratado.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c1c39be5-9bac-5277-ae0b-f305957a22c9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c1c39be5-9bac-5277-ae0b-f305957a22c9' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f2396bf2-5321-5027-bc83-8b37d29f8658', 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc', 'Documentos do contratado', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc' AND y.nome_documento = 'Documentos do contratado');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '380d6471-0f11-5dd5-a364-7b88afee8b4b', 'c1c39be5-9bac-5277-ae0b-f305957a22c9', 3, 'SEGEP/RH', 'Recursos Humanos', 'Formalizar o contrato e cadastrar na folha.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c1c39be5-9bac-5277-ae0b-f305957a22c9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c1c39be5-9bac-5277-ae0b-f305957a22c9' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '8528b69e-73a5-5227-af32-134c3ddc3212', '380d6471-0f11-5dd5-a364-7b88afee8b4b', 'Contrato temporário', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '380d6471-0f11-5dd5-a364-7b88afee8b4b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '380d6471-0f11-5dd5-a364-7b88afee8b4b' AND y.nome_documento = 'Contrato temporário');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3045c793-353e-5879-a809-170898cd13b3', '842990f9-8aa7-573f-8d01-73fd4d104bad', 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc', 'Necessidade justificada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '842990f9-8aa7-573f-8d01-73fd4d104bad')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '842990f9-8aa7-573f-8d01-73fd4d104bad' AND z.destino_etapa_id = 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc' AND z.condicao_transicao = 'Necessidade justificada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '68e85106-a7ec-5d90-acea-184f6f05da14', 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc', '380d6471-0f11-5dd5-a364-7b88afee8b4b', 'Contratação viável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '380d6471-0f11-5dd5-a364-7b88afee8b4b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ac69f6d8-7e8b-5cca-a408-861b76d2eddc' AND z.destino_etapa_id = '380d6471-0f11-5dd5-a364-7b88afee8b4b' AND z.condicao_transicao = 'Contratação viável');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ccc46f70-87a7-520f-8c3f-28577db3c967', 'RH.ETG.025', 'Contratação de estagiário', 'Contratar estagiário por meio de agente de integração e termo de compromisso com a instituição de ensino.', 'Secretarias e estudantes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.08.02.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.08.02.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.ETG.025');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f', 'ccc46f70-87a7-520f-8c3f-28577db3c967', 1, 'SECRETARIA', 'Unidade concedente', 'Solicitar a vaga e indicar o supervisor.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ccc46f70-87a7-520f-8c3f-28577db3c967')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ccc46f70-87a7-520f-8c3f-28577db3c967' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '69fcb40e-fc38-50eb-85a3-b5bcf0d42bc6', '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f', 'Requerimento de contratação de estagiário', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f' AND y.nome_documento = 'Requerimento de contratação de estagiário');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '562edd30-4669-5a7e-a61b-fe5ed979d29a', 'ccc46f70-87a7-520f-8c3f-28577db3c967', 2, 'SEGEP/CAPACITACAO', 'Gestão de estágios', 'Selecionar com o agente de integração e firmar o termo de compromisso.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ccc46f70-87a7-520f-8c3f-28577db3c967')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ccc46f70-87a7-520f-8c3f-28577db3c967' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a6a29076-edfa-5f09-8892-3957cb31aa95', '562edd30-4669-5a7e-a61b-fe5ed979d29a', 'Termo de compromisso de estágio', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '562edd30-4669-5a7e-a61b-fe5ed979d29a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '562edd30-4669-5a7e-a61b-fe5ed979d29a' AND y.nome_documento = 'Termo de compromisso de estágio');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a7e1535f-5656-5f12-9cc0-05a7caf76420', '562edd30-4669-5a7e-a61b-fe5ed979d29a', 'Plano de atividades', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '562edd30-4669-5a7e-a61b-fe5ed979d29a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '562edd30-4669-5a7e-a61b-fe5ed979d29a' AND y.nome_documento = 'Plano de atividades');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f116b97c-52df-5363-8550-3cb7f5e77bff', '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f', '562edd30-4669-5a7e-a61b-fe5ed979d29a', 'Vaga solicitada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '562edd30-4669-5a7e-a61b-fe5ed979d29a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '655c4e85-02eb-519d-9bf4-b4dd45ff4b1f' AND z.destino_etapa_id = '562edd30-4669-5a7e-a61b-fe5ed979d29a' AND z.condicao_transicao = 'Vaga solicitada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '7b536462-4bee-57bb-8808-a3456ee32988', 'RH.SEL.026', 'Processo seletivo simplificado', 'Planejar e realizar processo seletivo para contratação temporária, com edital, recursos e homologação.', 'Candidatos e secretarias', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '2.0.08.03.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '2.0.08.03.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'RH.SEL.026');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '27c75c2e-e7bf-532e-8e65-ba8137ce2645', '7b536462-4bee-57bb-8808-a3456ee32988', 1, 'SEGEP/RH', 'Comissão do processo seletivo', 'Planejar, publicar o edital e receber as inscrições.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7b536462-4bee-57bb-8808-a3456ee32988')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7b536462-4bee-57bb-8808-a3456ee32988' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '29150cac-b203-5670-be23-d55f68b04637', '27c75c2e-e7bf-532e-8e65-ba8137ce2645', 'Edital do processo seletivo', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '27c75c2e-e7bf-532e-8e65-ba8137ce2645')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '27c75c2e-e7bf-532e-8e65-ba8137ce2645' AND y.nome_documento = 'Edital do processo seletivo');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '27dd56cb-e8b8-5540-899f-b89c259039f2', '27c75c2e-e7bf-532e-8e65-ba8137ce2645', 'Lista de inscritos', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '27c75c2e-e7bf-532e-8e65-ba8137ce2645')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '27c75c2e-e7bf-532e-8e65-ba8137ce2645' AND y.nome_documento = 'Lista de inscritos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '780e5414-ea30-5242-928d-073f81092e46', '7b536462-4bee-57bb-8808-a3456ee32988', 2, 'SEGEP/RH', 'Comissão do processo seletivo', 'Avaliar, publicar o resultado e julgar os recursos.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7b536462-4bee-57bb-8808-a3456ee32988')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7b536462-4bee-57bb-8808-a3456ee32988' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ad541f05-c98a-5094-9880-9c4484ee1d33', '780e5414-ea30-5242-928d-073f81092e46', 'Resultado parcial', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '780e5414-ea30-5242-928d-073f81092e46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '780e5414-ea30-5242-928d-073f81092e46' AND y.nome_documento = 'Resultado parcial');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '79344e6f-92e9-5737-8445-b3a2e37d92fc', '780e5414-ea30-5242-928d-073f81092e46', 'Decisão dos recursos', FALSE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '780e5414-ea30-5242-928d-073f81092e46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '780e5414-ea30-5242-928d-073f81092e46' AND y.nome_documento = 'Decisão dos recursos');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd3552032-171a-5690-bd4c-0231aa2e8cdd', '7b536462-4bee-57bb-8808-a3456ee32988', 3, 'GAB', 'Prefeito', 'Homologar o resultado final.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7b536462-4bee-57bb-8808-a3456ee32988')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7b536462-4bee-57bb-8808-a3456ee32988' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b083ec1d-7a28-5327-9b38-63c199006d13', 'd3552032-171a-5690-bd4c-0231aa2e8cdd', 'Homologação do resultado', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd3552032-171a-5690-bd4c-0231aa2e8cdd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd3552032-171a-5690-bd4c-0231aa2e8cdd' AND y.nome_documento = 'Homologação do resultado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5bc12e15-eca1-528d-9800-1ec126fa783f', '27c75c2e-e7bf-532e-8e65-ba8137ce2645', '780e5414-ea30-5242-928d-073f81092e46', 'Inscrições encerradas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '27c75c2e-e7bf-532e-8e65-ba8137ce2645')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '780e5414-ea30-5242-928d-073f81092e46')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '27c75c2e-e7bf-532e-8e65-ba8137ce2645' AND z.destino_etapa_id = '780e5414-ea30-5242-928d-073f81092e46' AND z.condicao_transicao = 'Inscrições encerradas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4ac93bc0-4217-593d-9f48-5c8c83013677', '780e5414-ea30-5242-928d-073f81092e46', 'd3552032-171a-5690-bd4c-0231aa2e8cdd', 'Recursos julgados', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '780e5414-ea30-5242-928d-073f81092e46')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd3552032-171a-5690-bd4c-0231aa2e8cdd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '780e5414-ea30-5242-928d-073f81092e46' AND z.destino_etapa_id = 'd3552032-171a-5690-bd4c-0231aa2e8cdd' AND z.condicao_transicao = 'Recursos julgados');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '7a036b31-4a6d-5dcb-8814-f3df50d5628c', 'FAZ.FIS.007', 'Ação fiscal tributária', 'Fiscalizar o contribuinte, lavrar o auto de infração quando houver débito e encerrar a ação fiscal.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.03.01.23'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.03.01.23' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.FIS.007');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fd6fca82-c72e-572c-87ff-388072f9af25', '7a036b31-4a6d-5dcb-8814-f3df50d5628c', 1, 'SEREC/FISCALIZACAO', 'Fiscalização', 'Emitir a ordem de serviço e o termo de início de fiscalização.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a036b31-4a6d-5dcb-8814-f3df50d5628c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a036b31-4a6d-5dcb-8814-f3df50d5628c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4a225814-0b04-55bd-8ad9-5d69f0c28585', 'fd6fca82-c72e-572c-87ff-388072f9af25', 'Ordem de serviço', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fd6fca82-c72e-572c-87ff-388072f9af25')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fd6fca82-c72e-572c-87ff-388072f9af25' AND y.nome_documento = 'Ordem de serviço');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ed98c730-43b9-5d22-a4e7-a9d076e2b095', 'fd6fca82-c72e-572c-87ff-388072f9af25', 'Termo de início de fiscalização', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fd6fca82-c72e-572c-87ff-388072f9af25')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fd6fca82-c72e-572c-87ff-388072f9af25' AND y.nome_documento = 'Termo de início de fiscalização');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '56eae3cd-eb1b-5444-9e34-d0ad549de4fd', '7a036b31-4a6d-5dcb-8814-f3df50d5628c', 2, 'SEREC/FISCALIZACAO', 'Auditor fiscal', 'Levantar os documentos e apurar o débito.', 60, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a036b31-4a6d-5dcb-8814-f3df50d5628c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a036b31-4a6d-5dcb-8814-f3df50d5628c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'df0c2dc4-08e2-54b2-be86-b155564db528', '56eae3cd-eb1b-5444-9e34-d0ad549de4fd', 'Levantamento fiscal', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd' AND y.nome_documento = 'Levantamento fiscal');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2515476c-f817-5d1e-bab9-55e4d1b8a95f', '56eae3cd-eb1b-5444-9e34-d0ad549de4fd', 'Auto de infração e multa', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd' AND y.nome_documento = 'Auto de infração e multa');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b93b06e7-3632-56bd-ac3a-de17f330d8b3', '7a036b31-4a6d-5dcb-8814-f3df50d5628c', 3, 'SEREC/FISCALIZACAO', 'Fiscalização', 'Notificar o contribuinte e lavrar o termo de conclusão.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '7a036b31-4a6d-5dcb-8814-f3df50d5628c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '7a036b31-4a6d-5dcb-8814-f3df50d5628c' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c3775ef1-5b56-5740-a270-4365afe82045', 'b93b06e7-3632-56bd-ac3a-de17f330d8b3', 'Termo de conclusão de fiscalização', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b93b06e7-3632-56bd-ac3a-de17f330d8b3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b93b06e7-3632-56bd-ac3a-de17f330d8b3' AND y.nome_documento = 'Termo de conclusão de fiscalização');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3c397285-c1fd-5e18-a6b4-a22af9068795', 'fd6fca82-c72e-572c-87ff-388072f9af25', '56eae3cd-eb1b-5444-9e34-d0ad549de4fd', 'Fiscalização iniciada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'fd6fca82-c72e-572c-87ff-388072f9af25')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'fd6fca82-c72e-572c-87ff-388072f9af25' AND z.destino_etapa_id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd' AND z.condicao_transicao = 'Fiscalização iniciada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'bbd373f9-9fac-5706-bfcc-3ef3e8a43b2c', '56eae3cd-eb1b-5444-9e34-d0ad549de4fd', 'b93b06e7-3632-56bd-ac3a-de17f330d8b3', 'Apuração concluída', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b93b06e7-3632-56bd-ac3a-de17f330d8b3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '56eae3cd-eb1b-5444-9e34-d0ad549de4fd' AND z.destino_etapa_id = 'b93b06e7-3632-56bd-ac3a-de17f330d8b3' AND z.condicao_transicao = 'Apuração concluída');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'b3423d81-245d-5721-9eb0-e7aa79e34473', 'FAZ.ISS.008', 'Recurso contra lançamento de ISSQN', 'Julgar a impugnação ou o recurso do contribuinte contra lançamento ou auto de infração de ISSQN.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.03.01.13'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.03.01.13' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.ISS.008');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '895c500e-0891-5115-a841-df482e279f19', 'b3423d81-245d-5721-9eb0-e7aa79e34473', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber o recurso no prazo e juntar o lançamento.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b3423d81-245d-5721-9eb0-e7aa79e34473')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b3423d81-245d-5721-9eb0-e7aa79e34473' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd274aae1-d53e-5f0b-a1ef-02b7725eacf5', '895c500e-0891-5115-a841-df482e279f19', 'Recurso do contribuinte', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '895c500e-0891-5115-a841-df482e279f19')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '895c500e-0891-5115-a841-df482e279f19' AND y.nome_documento = 'Recurso do contribuinte');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '92987edc-1321-5545-8b0a-61ed4e58ddbf', 'b3423d81-245d-5721-9eb0-e7aa79e34473', 2, 'SEREC/FISCALIZACAO', 'Auditor autuante', 'Contestar ou rever o lançamento.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b3423d81-245d-5721-9eb0-e7aa79e34473')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b3423d81-245d-5721-9eb0-e7aa79e34473' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6e803cf9-6275-5f12-8244-bf047b41e601', '92987edc-1321-5545-8b0a-61ed4e58ddbf', 'Informação fiscal', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '92987edc-1321-5545-8b0a-61ed4e58ddbf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '92987edc-1321-5545-8b0a-61ed4e58ddbf' AND y.nome_documento = 'Informação fiscal');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8796ab6c-3e00-5497-a2a0-51418446aab6', 'b3423d81-245d-5721-9eb0-e7aa79e34473', 3, 'SEREC/TRIBUTARIO', 'Julgamento tributário', 'Julgar e notificar o contribuinte.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b3423d81-245d-5721-9eb0-e7aa79e34473')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b3423d81-245d-5721-9eb0-e7aa79e34473' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'fab22571-c2e6-5bd4-96d9-a6a9e85f8b0f', '8796ab6c-3e00-5497-a2a0-51418446aab6', 'Decisão de primeira instância', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8796ab6c-3e00-5497-a2a0-51418446aab6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8796ab6c-3e00-5497-a2a0-51418446aab6' AND y.nome_documento = 'Decisão de primeira instância');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'db06dc48-742e-5e6a-b7d8-fae0bc832915', '895c500e-0891-5115-a841-df482e279f19', '92987edc-1321-5545-8b0a-61ed4e58ddbf', 'Recurso tempestivo', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '895c500e-0891-5115-a841-df482e279f19')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '92987edc-1321-5545-8b0a-61ed4e58ddbf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '895c500e-0891-5115-a841-df482e279f19' AND z.destino_etapa_id = '92987edc-1321-5545-8b0a-61ed4e58ddbf' AND z.condicao_transicao = 'Recurso tempestivo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '347b589d-95b3-56f7-94bd-d66381b309da', '92987edc-1321-5545-8b0a-61ed4e58ddbf', '8796ab6c-3e00-5497-a2a0-51418446aab6', 'Informação prestada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '92987edc-1321-5545-8b0a-61ed4e58ddbf')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8796ab6c-3e00-5497-a2a0-51418446aab6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '92987edc-1321-5545-8b0a-61ed4e58ddbf' AND z.destino_etapa_id = '8796ab6c-3e00-5497-a2a0-51418446aab6' AND z.condicao_transicao = 'Informação prestada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '80c914be-ca08-5c33-a7a8-76d68cd88228', 'FAZ.JUL.009', 'Consulta tributária', 'Responder formalmente à consulta do contribuinte sobre a aplicação da legislação tributária municipal.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.03.02.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.03.02.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.JUL.009');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '011284e8-943c-5e12-887e-d624a91210b0', '80c914be-ca08-5c33-a7a8-76d68cd88228', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber a consulta com a descrição do fato.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '80c914be-ca08-5c33-a7a8-76d68cd88228')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '80c914be-ca08-5c33-a7a8-76d68cd88228' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd0a2c0c1-3d7c-5f3d-a499-6bbe98f0256a', '011284e8-943c-5e12-887e-d624a91210b0', 'Petição de consulta', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '011284e8-943c-5e12-887e-d624a91210b0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '011284e8-943c-5e12-887e-d624a91210b0' AND y.nome_documento = 'Petição de consulta');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd74e1f20-053c-5281-8b59-8922a7802299', '80c914be-ca08-5c33-a7a8-76d68cd88228', 2, 'SEREC/TRIBUTARIO', 'Julgamento tributário', 'Analisar e emitir a resposta à consulta.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '80c914be-ca08-5c33-a7a8-76d68cd88228')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '80c914be-ca08-5c33-a7a8-76d68cd88228' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd2580884-c77a-52b5-a338-041f95425c45', 'd74e1f20-053c-5281-8b59-8922a7802299', 'Resposta à consulta', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd74e1f20-053c-5281-8b59-8922a7802299')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd74e1f20-053c-5281-8b59-8922a7802299' AND y.nome_documento = 'Resposta à consulta');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'cbef0440-95e2-554b-a21c-506af43f4e63', '011284e8-943c-5e12-887e-d624a91210b0', 'd74e1f20-053c-5281-8b59-8922a7802299', 'Consulta recebida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '011284e8-943c-5e12-887e-d624a91210b0')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd74e1f20-053c-5281-8b59-8922a7802299')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '011284e8-943c-5e12-887e-d624a91210b0' AND z.destino_etapa_id = 'd74e1f20-053c-5281-8b59-8922a7802299' AND z.condicao_transicao = 'Consulta recebida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'a0930f51-455b-5478-a12f-8bf1be1241cf', 'FAZ.FPM.010', 'Recurso sobre o índice do FPM', 'Solicitar à Procuradoria a interposição de recurso contra o índice ou a distribuição do Fundo de Participação dos Municípios.', 'Gestão fazendária', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.03.03.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.03.03.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.FPM.010');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1', 'a0930f51-455b-5478-a12f-8bf1be1241cf', 1, 'SEREC/RECEITA', 'Receita', 'Analisar o índice publicado e levantar os dados divergentes.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a0930f51-455b-5478-a12f-8bf1be1241cf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a0930f51-455b-5478-a12f-8bf1be1241cf' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '39b080d7-2c7f-58bc-ae83-031e2308e49f', 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1', 'Nota técnica do índice', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1' AND y.nome_documento = 'Nota técnica do índice');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574', 'a0930f51-455b-5478-a12f-8bf1be1241cf', 2, 'PGM', 'Procuradoria', 'Interpor o recurso no órgão competente.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a0930f51-455b-5478-a12f-8bf1be1241cf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a0930f51-455b-5478-a12f-8bf1be1241cf' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a810d0d1-f983-58a2-8a0a-ad6d9a019042', 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574', 'Recurso administrativo', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574' AND y.nome_documento = 'Recurso administrativo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'adac8970-c275-5d01-a1c6-a0d995f1fda2', 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1', 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574', 'Divergência confirmada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c8ce5c38-1f5c-572e-b343-16a96c2a5fb1' AND z.destino_etapa_id = 'fbfba10d-94c8-5779-8ab1-d12d4cc2f574' AND z.condicao_transicao = 'Divergência confirmada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '9837a24f-16d3-501f-8ed9-50cd68849172', 'FAZ.CAD.011', 'Atualização do cadastro imobiliário (proprietário e endereço)', 'Atualizar o nome do proprietário ou o endereço de correspondência no cadastro imobiliário.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.04.00.05'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.04.00.05' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.CAD.011');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '911af39c-4a93-5ea1-962d-9d7d06b14812', '9837a24f-16d3-501f-8ed9-50cd68849172', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber o pedido com o documento de propriedade.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9837a24f-16d3-501f-8ed9-50cd68849172')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9837a24f-16d3-501f-8ed9-50cd68849172' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e14c3c8e-d350-5a18-ab59-52c919e7a1dd', '911af39c-4a93-5ea1-962d-9d7d06b14812', 'Requerimento de atualização cadastral', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '911af39c-4a93-5ea1-962d-9d7d06b14812')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '911af39c-4a93-5ea1-962d-9d7d06b14812' AND y.nome_documento = 'Requerimento de atualização cadastral');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '055da07c-369d-56cf-8ff8-ff85537f863a', '911af39c-4a93-5ea1-962d-9d7d06b14812', 'Matrícula ou escritura', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '911af39c-4a93-5ea1-962d-9d7d06b14812')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '911af39c-4a93-5ea1-962d-9d7d06b14812' AND y.nome_documento = 'Matrícula ou escritura');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9f323c1a-6109-5028-bff3-ed96b517954f', '9837a24f-16d3-501f-8ed9-50cd68849172', 2, 'SEREC/CADASTRO', 'Cadastro imobiliário', 'Conferir e atualizar o cadastro.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9837a24f-16d3-501f-8ed9-50cd68849172')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9837a24f-16d3-501f-8ed9-50cd68849172' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a1305448-8d6d-5c92-878c-6c5fea228b27', '9f323c1a-6109-5028-bff3-ed96b517954f', 'Boletim de cadastro imobiliário', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9f323c1a-6109-5028-bff3-ed96b517954f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9f323c1a-6109-5028-bff3-ed96b517954f' AND y.nome_documento = 'Boletim de cadastro imobiliário');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '033f4620-073c-59a1-97f1-6b6e80904e4f', '911af39c-4a93-5ea1-962d-9d7d06b14812', '9f323c1a-6109-5028-bff3-ed96b517954f', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '911af39c-4a93-5ea1-962d-9d7d06b14812')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '9f323c1a-6109-5028-bff3-ed96b517954f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '911af39c-4a93-5ea1-962d-9d7d06b14812' AND z.destino_etapa_id = '9f323c1a-6109-5028-bff3-ed96b517954f' AND z.condicao_transicao = 'Pedido instruído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '8b475ee3-52f2-5760-b3ec-d2e010960ff2', 'FAZ.VVE.012', 'Certidão de valor venal', 'Emitir certidão do valor venal do imóvel para fins de ITBI, financiamento ou inventário.', 'Contribuintes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.04.00.07'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.04.00.07' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.VVE.012');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '649d0bc3-161d-550e-806f-204cda5813ae', '8b475ee3-52f2-5760-b3ec-d2e010960ff2', 1, 'SEREC/ATENDIMENTO', 'Atendimento ao contribuinte', 'Receber o pedido.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '8b475ee3-52f2-5760-b3ec-d2e010960ff2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '8b475ee3-52f2-5760-b3ec-d2e010960ff2' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b44eba89-2829-578f-a135-483fa19e47c8', '649d0bc3-161d-550e-806f-204cda5813ae', 'Requerimento de certidão', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '649d0bc3-161d-550e-806f-204cda5813ae')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '649d0bc3-161d-550e-806f-204cda5813ae' AND y.nome_documento = 'Requerimento de certidão');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73', '8b475ee3-52f2-5760-b3ec-d2e010960ff2', 2, 'SEREC/CADASTRO', 'Cadastro imobiliário', 'Emitir a certidão.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '8b475ee3-52f2-5760-b3ec-d2e010960ff2')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '8b475ee3-52f2-5760-b3ec-d2e010960ff2' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '019b6e87-dd68-5a78-bad9-2ea8271763a7', 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73', 'Certidão de valor venal', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73' AND y.nome_documento = 'Certidão de valor venal');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'ae542e26-db46-5482-88f5-630366449d85', '649d0bc3-161d-550e-806f-204cda5813ae', 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73', 'Pedido recebido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '649d0bc3-161d-550e-806f-204cda5813ae')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '649d0bc3-161d-550e-806f-204cda5813ae' AND z.destino_etapa_id = 'fd2ced0b-ab7b-517e-bc45-b9308bca7e73' AND z.condicao_transicao = 'Pedido recebido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'c88acae8-34f0-59aa-9429-ef3fcb237c55', 'FAZ.CUS.013', 'Certidão de uso do solo (controle urbano)', 'Emitir certidão de uso do solo do imóvel conforme o zoneamento.', 'Proprietários e empresas', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.04.01.06'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.04.01.06' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.CUS.013');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89', 'c88acae8-34f0-59aa-9429-ef3fcb237c55', 1, 'SEREC/ATENDIMENTO', 'Atendimento', 'Receber o pedido com a localização do imóvel e a atividade.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c88acae8-34f0-59aa-9429-ef3fcb237c55')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c88acae8-34f0-59aa-9429-ef3fcb237c55' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '1f2739f7-1940-53e3-b3e9-495525b66634', '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89', 'Requerimento de certidão de uso do solo', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89' AND y.nome_documento = 'Requerimento de certidão de uso do solo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0caddf34-8f03-53ad-8231-0d4dc0d8e174', 'c88acae8-34f0-59aa-9429-ef3fcb237c55', 2, 'SEREC/CONTROLE_URBANO', 'Controle urbano', 'Verificar o zoneamento e emitir a certidão.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c88acae8-34f0-59aa-9429-ef3fcb237c55')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c88acae8-34f0-59aa-9429-ef3fcb237c55' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cd0cfda3-f97d-5d06-b602-abc7641a163f', '0caddf34-8f03-53ad-8231-0d4dc0d8e174', 'Certidão de uso do solo', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0caddf34-8f03-53ad-8231-0d4dc0d8e174')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0caddf34-8f03-53ad-8231-0d4dc0d8e174' AND y.nome_documento = 'Certidão de uso do solo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '73f50511-7f29-5854-89a6-c5eb6c7807b6', '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89', '0caddf34-8f03-53ad-8231-0d4dc0d8e174', 'Pedido recebido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '0caddf34-8f03-53ad-8231-0d4dc0d8e174')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '9367e06e-b4a9-5fc2-aa5f-238aab8a7e89' AND z.destino_etapa_id = '0caddf34-8f03-53ad-8231-0d4dc0d8e174' AND z.condicao_transicao = 'Pedido recebido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1d78e555-138b-52db-a0b6-d0d22b58632b', 'FAZ.ORC.014', 'Suplementação orçamentária', 'Abrir crédito suplementar ou especial para despesa sem dotação suficiente.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.05.00.08'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.05.00.08' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'FAZ.ORC.014');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '71c59a20-4213-56c1-8b51-21ae74a09fff', '1d78e555-138b-52db-a0b6-d0d22b58632b', 1, 'SECRETARIA', 'Secretaria solicitante', 'Justificar a necessidade e indicar a fonte de recurso.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1d78e555-138b-52db-a0b6-d0d22b58632b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1d78e555-138b-52db-a0b6-d0d22b58632b' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '66cf15d1-fa04-58d8-875c-0370ed0a496c', '71c59a20-4213-56c1-8b51-21ae74a09fff', 'Solicitação de suplementação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '71c59a20-4213-56c1-8b51-21ae74a09fff')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '71c59a20-4213-56c1-8b51-21ae74a09fff' AND y.nome_documento = 'Solicitação de suplementação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd36274d6-0437-5504-9023-145da834723d', '1d78e555-138b-52db-a0b6-d0d22b58632b', 2, 'SEPLAN/ORCAMENTO', 'Orçamento', 'Verificar a disponibilidade e minutar o decreto.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1d78e555-138b-52db-a0b6-d0d22b58632b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1d78e555-138b-52db-a0b6-d0d22b58632b' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '447268f8-0090-5438-bd05-78ae36e98e1a', 'd36274d6-0437-5504-9023-145da834723d', 'Minuta de decreto de crédito', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd36274d6-0437-5504-9023-145da834723d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd36274d6-0437-5504-9023-145da834723d' AND y.nome_documento = 'Minuta de decreto de crédito');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64', '1d78e555-138b-52db-a0b6-d0d22b58632b', 3, 'GAB', 'Prefeito', 'Assinar e publicar o decreto.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1d78e555-138b-52db-a0b6-d0d22b58632b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1d78e555-138b-52db-a0b6-d0d22b58632b' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '56f7077d-01f6-593f-978b-8632a8610eca', 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64', 'Decreto de abertura de crédito', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64' AND y.nome_documento = 'Decreto de abertura de crédito');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '385536e9-e217-5577-b522-55519a54280c', '71c59a20-4213-56c1-8b51-21ae74a09fff', 'd36274d6-0437-5504-9023-145da834723d', 'Pedido justificado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '71c59a20-4213-56c1-8b51-21ae74a09fff')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd36274d6-0437-5504-9023-145da834723d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '71c59a20-4213-56c1-8b51-21ae74a09fff' AND z.destino_etapa_id = 'd36274d6-0437-5504-9023-145da834723d' AND z.condicao_transicao = 'Pedido justificado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '589a59e8-8171-5457-a656-50e52b1642dd', 'd36274d6-0437-5504-9023-145da834723d', 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64', 'Crédito viável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'd36274d6-0437-5504-9023-145da834723d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'd36274d6-0437-5504-9023-145da834723d' AND z.destino_etapa_id = 'b9dc7a5b-f9ea-55c2-bbbd-174ebcc5ba64' AND z.condicao_transicao = 'Crédito viável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b80e079f-25c2-53f2-ae91-7a2afd10c433', 'd36274d6-0437-5504-9023-145da834723d', '71c59a20-4213-56c1-8b51-21ae74a09fff', 'Fonte insuficiente', TRUE, NULLIF('Indicar outra fonte de recurso', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'd36274d6-0437-5504-9023-145da834723d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '71c59a20-4213-56c1-8b51-21ae74a09fff')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'd36274d6-0437-5504-9023-145da834723d' AND z.destino_etapa_id = '71c59a20-4213-56c1-8b51-21ae74a09fff' AND z.condicao_transicao = 'Fonte insuficiente');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1928c414-b16b-5ec8-8b0d-0f37046375a7', 'PAT.BAI.001', 'Baixa de bem patrimonial', 'Dar baixa de bem móvel inservível, extraviado ou alienado, com laudo e autorização.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.05.01.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.05.01.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'PAT.BAI.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a844faa5-dd63-523f-845f-bb38042120e6', '1928c414-b16b-5ec8-8b0d-0f37046375a7', 1, 'SECRETARIA', 'Unidade detentora', 'Solicitar a baixa com o número de patrimônio e o motivo.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1928c414-b16b-5ec8-8b0d-0f37046375a7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1928c414-b16b-5ec8-8b0d-0f37046375a7' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e8815e93-53a2-5a5c-8d66-f9d9f8479df3', 'a844faa5-dd63-523f-845f-bb38042120e6', 'Solicitação de baixa', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a844faa5-dd63-523f-845f-bb38042120e6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a844faa5-dd63-523f-845f-bb38042120e6' AND y.nome_documento = 'Solicitação de baixa');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '1fa610e9-ff3f-5931-bb33-e2003c40336b', 'a844faa5-dd63-523f-845f-bb38042120e6', 'Boletim de ocorrência (extravio)', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a844faa5-dd63-523f-845f-bb38042120e6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a844faa5-dd63-523f-845f-bb38042120e6' AND y.nome_documento = 'Boletim de ocorrência (extravio)');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f', '1928c414-b16b-5ec8-8b0d-0f37046375a7', 2, 'SEFIN/PATRIMONIO', 'Comissão de patrimônio', 'Vistoriar e emitir o laudo de inservibilidade.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1928c414-b16b-5ec8-8b0d-0f37046375a7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1928c414-b16b-5ec8-8b0d-0f37046375a7' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'fcbbbc2c-c7bb-53c1-90aa-0294f172cc4c', 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f', 'Laudo de avaliação', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f' AND y.nome_documento = 'Laudo de avaliação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4', '1928c414-b16b-5ec8-8b0d-0f37046375a7', 3, 'SEFIN/PATRIMONIO', 'Patrimônio', 'Registrar a baixa e destinar o bem.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1928c414-b16b-5ec8-8b0d-0f37046375a7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1928c414-b16b-5ec8-8b0d-0f37046375a7' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e1cc0b10-1284-5a7c-acc0-3a5b299faca2', 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4', 'Termo de baixa', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4' AND y.nome_documento = 'Termo de baixa');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3616136a-69c9-537b-bc9f-44cb07a66495', 'a844faa5-dd63-523f-845f-bb38042120e6', 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a844faa5-dd63-523f-845f-bb38042120e6')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a844faa5-dd63-523f-845f-bb38042120e6' AND z.destino_etapa_id = 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b696fd1e-0b6f-5e4a-8b8e-7e08d762e598', 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f', 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4', 'Baixa recomendada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a2713eb3-5b5f-5e48-a7cb-f1af74480c2f' AND z.destino_etapa_id = 'b5079d5b-38fd-502e-b3cc-24b118f0b5f4' AND z.condicao_transicao = 'Baixa recomendada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'f075c39b-ba9a-5af6-9df0-d42b6731168e', 'PAT.TRA.002', 'Transferência de bem patrimonial', 'Transferir a carga de bem móvel entre unidades, atualizando o responsável.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.05.01.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.05.01.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'PAT.TRA.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9012a8fe-a531-5c16-bf25-00ca55b97a16', 'f075c39b-ba9a-5af6-9df0-d42b6731168e', 1, 'SECRETARIA', 'Unidade cedente', 'Solicitar a transferência indicando a unidade recebedora.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f075c39b-ba9a-5af6-9df0-d42b6731168e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f075c39b-ba9a-5af6-9df0-d42b6731168e' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f2a92f98-aa7c-5c87-952b-4417ed8785b5', '9012a8fe-a531-5c16-bf25-00ca55b97a16', 'Solicitação de transferência', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9012a8fe-a531-5c16-bf25-00ca55b97a16')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9012a8fe-a531-5c16-bf25-00ca55b97a16' AND y.nome_documento = 'Solicitação de transferência');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'f8d06a03-26a7-5549-9fec-c98c3c9c850f', 'f075c39b-ba9a-5af6-9df0-d42b6731168e', 2, 'SEFIN/PATRIMONIO', 'Patrimônio', 'Emitir o termo de transferência de responsabilidade.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f075c39b-ba9a-5af6-9df0-d42b6731168e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f075c39b-ba9a-5af6-9df0-d42b6731168e' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e4bd5f1f-ca6a-588b-9a2e-64ed294b4c23', 'f8d06a03-26a7-5549-9fec-c98c3c9c850f', 'Termo de transferência', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'f8d06a03-26a7-5549-9fec-c98c3c9c850f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'f8d06a03-26a7-5549-9fec-c98c3c9c850f' AND y.nome_documento = 'Termo de transferência');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '39131204-1f24-5e65-8189-5f8ded94eab7', '9012a8fe-a531-5c16-bf25-00ca55b97a16', 'f8d06a03-26a7-5549-9fec-c98c3c9c850f', 'Pedido feito', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '9012a8fe-a531-5c16-bf25-00ca55b97a16')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'f8d06a03-26a7-5549-9fec-c98c3c9c850f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '9012a8fe-a531-5c16-bf25-00ca55b97a16' AND z.destino_etapa_id = 'f8d06a03-26a7-5549-9fec-c98c3c9c850f' AND z.condicao_transicao = 'Pedido feito');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'ba5a9be9-52b1-5735-acd7-2c2486283eda', 'PAT.INC.003', 'Incorporação de bem permanente', 'Incorporar ao patrimônio o bem adquirido, doado ou cedido, com plaqueta e termo de responsabilidade.', 'Secretarias municipais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '3.0.05.01.15'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '3.0.05.01.15' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'PAT.INC.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4cbb515f-16e0-52ae-89d9-7574d8408c43', 'ba5a9be9-52b1-5735-acd7-2c2486283eda', 1, 'SEFIN/PATRIMONIO', 'Patrimônio', 'Conferir a nota fiscal ou o termo de doação e emplaquetar.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ba5a9be9-52b1-5735-acd7-2c2486283eda')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ba5a9be9-52b1-5735-acd7-2c2486283eda' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b2a0871d-f5a9-5722-9bae-f996cbfb028d', '4cbb515f-16e0-52ae-89d9-7574d8408c43', 'Nota fiscal ou termo de doação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4cbb515f-16e0-52ae-89d9-7574d8408c43')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4cbb515f-16e0-52ae-89d9-7574d8408c43' AND y.nome_documento = 'Nota fiscal ou termo de doação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '56bdb74c-d633-582e-b8aa-f99f5b2f383b', '4cbb515f-16e0-52ae-89d9-7574d8408c43', 'Ficha de cadastro do bem', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4cbb515f-16e0-52ae-89d9-7574d8408c43')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4cbb515f-16e0-52ae-89d9-7574d8408c43' AND y.nome_documento = 'Ficha de cadastro do bem');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'cd635038-dc50-5dfa-99e3-3053c55379aa', 'ba5a9be9-52b1-5735-acd7-2c2486283eda', 2, 'SECRETARIA', 'Unidade recebedora', 'Assinar o termo de responsabilidade.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'ba5a9be9-52b1-5735-acd7-2c2486283eda')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'ba5a9be9-52b1-5735-acd7-2c2486283eda' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'deea0a2c-7e00-55d6-a870-59a3d2606c1b', 'cd635038-dc50-5dfa-99e3-3053c55379aa', 'Termo de responsabilidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cd635038-dc50-5dfa-99e3-3053c55379aa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cd635038-dc50-5dfa-99e3-3053c55379aa' AND y.nome_documento = 'Termo de responsabilidade');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2a380a2f-b12a-5f18-a564-9b14d9129796', '4cbb515f-16e0-52ae-89d9-7574d8408c43', 'cd635038-dc50-5dfa-99e3-3053c55379aa', 'Bem cadastrado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4cbb515f-16e0-52ae-89d9-7574d8408c43')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'cd635038-dc50-5dfa-99e3-3053c55379aa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4cbb515f-16e0-52ae-89d9-7574d8408c43' AND z.destino_etapa_id = 'cd635038-dc50-5dfa-99e3-3053c55379aa' AND z.condicao_transicao = 'Bem cadastrado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '84b021b2-816b-5275-8cf8-c91e451bafac', 'CTR.AUD.002', 'Auditoria interna', 'Planejar e executar auditoria em unidade ou processo, com relatório e plano de providências.', 'Unidades administrativas', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 7º, §3º (documento preparatório até a conclusão)', ''), '4.0.01.01.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '4.0.01.01.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'CTR.AUD.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e34bb009-3135-5c1a-8a30-c340a4b732fa', '84b021b2-816b-5275-8cf8-c91e451bafac', 1, 'UCCI', 'Controladoria', 'Planejar a auditoria e comunicar a unidade auditada.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '84b021b2-816b-5275-8cf8-c91e451bafac')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '84b021b2-816b-5275-8cf8-c91e451bafac' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2e08fdb1-9719-50b5-8e26-4a7161d80528', 'e34bb009-3135-5c1a-8a30-c340a4b732fa', 'Plano de auditoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e34bb009-3135-5c1a-8a30-c340a4b732fa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e34bb009-3135-5c1a-8a30-c340a4b732fa' AND y.nome_documento = 'Plano de auditoria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bb0a0b1d-9fd8-5744-bff2-3b4e25dec3b3', 'e34bb009-3135-5c1a-8a30-c340a4b732fa', 'Ofício de comunicação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e34bb009-3135-5c1a-8a30-c340a4b732fa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e34bb009-3135-5c1a-8a30-c340a4b732fa' AND y.nome_documento = 'Ofício de comunicação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '6b872731-bf70-5488-9a4d-dc93bef3705c', '84b021b2-816b-5275-8cf8-c91e451bafac', 2, 'UCCI', 'Equipe de auditoria', 'Executar os testes e emitir o relatório preliminar.', 45, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '84b021b2-816b-5275-8cf8-c91e451bafac')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '84b021b2-816b-5275-8cf8-c91e451bafac' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '414a07d6-5e28-5384-8bd4-8a63ed7d6f4a', '6b872731-bf70-5488-9a4d-dc93bef3705c', 'Relatório preliminar', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6b872731-bf70-5488-9a4d-dc93bef3705c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6b872731-bf70-5488-9a4d-dc93bef3705c' AND y.nome_documento = 'Relatório preliminar');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b', '84b021b2-816b-5275-8cf8-c91e451bafac', 3, 'SECRETARIA', 'Unidade auditada', 'Manifestar-se sobre os achados.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '84b021b2-816b-5275-8cf8-c91e451bafac')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '84b021b2-816b-5275-8cf8-c91e451bafac' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4987f55d-b4fe-585e-ad98-4a258f17e916', 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b', 'Manifestação da unidade', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b' AND y.nome_documento = 'Manifestação da unidade');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3131062c-9a5b-5af3-b687-daf517f70a1f', '84b021b2-816b-5275-8cf8-c91e451bafac', 4, 'UCCI', 'Controladoria', 'Emitir o relatório final e o plano de providências.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '84b021b2-816b-5275-8cf8-c91e451bafac')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '84b021b2-816b-5275-8cf8-c91e451bafac' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a5557da6-36fc-56da-b4e6-406c9bb9ca31', '3131062c-9a5b-5af3-b687-daf517f70a1f', 'Relatório final de auditoria', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3131062c-9a5b-5af3-b687-daf517f70a1f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3131062c-9a5b-5af3-b687-daf517f70a1f' AND y.nome_documento = 'Relatório final de auditoria');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '7278cb41-5d94-50ed-b2ea-5fc2123f6f2c', 'e34bb009-3135-5c1a-8a30-c340a4b732fa', '6b872731-bf70-5488-9a4d-dc93bef3705c', 'Auditoria iniciada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e34bb009-3135-5c1a-8a30-c340a4b732fa')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '6b872731-bf70-5488-9a4d-dc93bef3705c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e34bb009-3135-5c1a-8a30-c340a4b732fa' AND z.destino_etapa_id = '6b872731-bf70-5488-9a4d-dc93bef3705c' AND z.condicao_transicao = 'Auditoria iniciada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '55d08e80-7cfa-589f-bf7c-d53947dbfd2c', '6b872731-bf70-5488-9a4d-dc93bef3705c', 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b', 'Relatório emitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '6b872731-bf70-5488-9a4d-dc93bef3705c')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '6b872731-bf70-5488-9a4d-dc93bef3705c' AND z.destino_etapa_id = 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b' AND z.condicao_transicao = 'Relatório emitido');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '93992a65-bd20-53cc-8bf1-ce1ccb7182a9', 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b', '3131062c-9a5b-5af3-b687-daf517f70a1f', 'Manifestação recebida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3131062c-9a5b-5af3-b687-daf517f70a1f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'bf6ee0ba-43ed-586b-9bd5-8cc6f0bba76b' AND z.destino_etapa_id = '3131062c-9a5b-5af3-b687-daf517f70a1f' AND z.condicao_transicao = 'Manifestação recebida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '0cbad283-363e-5754-bd41-4ec254d4fa66', 'HAB.ITB.004', 'Isenção de ITBI em programa habitacional', 'Reconhecer a isenção de ITBI para beneficiário de programa habitacional de interesse social.', 'Beneficiários de programas habitacionais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal)', ''), '6.0.01.00.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.01.00.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.ITB.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ed145619-261f-58a2-8d06-6eea1acd22db', '0cbad283-363e-5754-bd41-4ec254d4fa66', 1, 'SEHAB/PROGRAMAS', 'Programas habitacionais', 'Confirmar o enquadramento do beneficiário no programa.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0cbad283-363e-5754-bd41-4ec254d4fa66')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0cbad283-363e-5754-bd41-4ec254d4fa66' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4db2ea67-011a-580f-aa73-9d5b9301488a', 'ed145619-261f-58a2-8d06-6eea1acd22db', 'Requerimento de isenção de ITBI', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ed145619-261f-58a2-8d06-6eea1acd22db')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ed145619-261f-58a2-8d06-6eea1acd22db' AND y.nome_documento = 'Requerimento de isenção de ITBI');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '12f283eb-d8a5-53dc-b268-1dd7e3d92ba5', 'ed145619-261f-58a2-8d06-6eea1acd22db', 'Declaração de enquadramento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ed145619-261f-58a2-8d06-6eea1acd22db')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ed145619-261f-58a2-8d06-6eea1acd22db' AND y.nome_documento = 'Declaração de enquadramento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f', '0cbad283-363e-5754-bd41-4ec254d4fa66', 2, 'SEREC/TRIBUTARIO', 'Julgamento tributário', 'Reconhecer a isenção e emitir a guia zerada.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0cbad283-363e-5754-bd41-4ec254d4fa66')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0cbad283-363e-5754-bd41-4ec254d4fa66' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ca3d8496-a7af-56ed-9c1d-9d27c04c85ab', 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f', 'Decisão de isenção', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f' AND y.nome_documento = 'Decisão de isenção');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fc391724-afa6-5f35-9738-5d3cbf3585d0', 'ed145619-261f-58a2-8d06-6eea1acd22db', 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f', 'Beneficiário enquadrado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ed145619-261f-58a2-8d06-6eea1acd22db')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ed145619-261f-58a2-8d06-6eea1acd22db' AND z.destino_etapa_id = 'beb9d4a0-b74d-5350-a2b7-d6fc3312c09f' AND z.condicao_transicao = 'Beneficiário enquadrado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '6ce8119b-0495-5433-b91f-97a8fa507c88', 'HAB.MCM.005', 'Titulação de beneficiário de programa habitacional', 'Emitir o título de propriedade ao beneficiário de programa habitacional (ex.: Minha Casa Minha Vida).', 'Beneficiários de programas habitacionais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal)', ''), '6.0.01.00.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.01.00.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.MCM.005');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '00dbc5c1-99d4-56ab-aa72-e571efbb7283', '6ce8119b-0495-5433-b91f-97a8fa507c88', 1, 'SEHAB/PROGRAMAS', 'Programas habitacionais', 'Conferir a ocupação e os requisitos do beneficiário.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '6ce8119b-0495-5433-b91f-97a8fa507c88')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '6ce8119b-0495-5433-b91f-97a8fa507c88' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '91df6617-3324-5dbb-8896-109a0e4c4ba9', '00dbc5c1-99d4-56ab-aa72-e571efbb7283', 'Cadastro do beneficiário', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '00dbc5c1-99d4-56ab-aa72-e571efbb7283')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '00dbc5c1-99d4-56ab-aa72-e571efbb7283' AND y.nome_documento = 'Cadastro do beneficiário');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'dd9afb09-f918-56c4-be6f-07707e6ad98f', '00dbc5c1-99d4-56ab-aa72-e571efbb7283', 'Documentos pessoais', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '00dbc5c1-99d4-56ab-aa72-e571efbb7283')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '00dbc5c1-99d4-56ab-aa72-e571efbb7283' AND y.nome_documento = 'Documentos pessoais');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '82304cb8-d398-5c06-9c0b-1aad6a2c3e61', '6ce8119b-0495-5433-b91f-97a8fa507c88', 2, 'PGM', 'Procuradoria', 'Emitir parecer jurídico.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '6ce8119b-0495-5433-b91f-97a8fa507c88')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '6ce8119b-0495-5433-b91f-97a8fa507c88' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '294231bd-5495-5205-8475-3626904535e4', '82304cb8-d398-5c06-9c0b-1aad6a2c3e61', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '82304cb8-d398-5c06-9c0b-1aad6a2c3e61')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '82304cb8-d398-5c06-9c0b-1aad6a2c3e61' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd2711ca3-d944-503c-8c5a-f95ba79a21d4', '6ce8119b-0495-5433-b91f-97a8fa507c88', 3, 'SEHAB/PROGRAMAS', 'Programas habitacionais', 'Emitir o título e encaminhar ao cartório.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '6ce8119b-0495-5433-b91f-97a8fa507c88')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '6ce8119b-0495-5433-b91f-97a8fa507c88' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'afafc639-08c0-5bfc-bf22-d2a465c1ce70', 'd2711ca3-d944-503c-8c5a-f95ba79a21d4', 'Título de propriedade', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd2711ca3-d944-503c-8c5a-f95ba79a21d4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd2711ca3-d944-503c-8c5a-f95ba79a21d4' AND y.nome_documento = 'Título de propriedade');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '15e8572c-edb6-500a-9e92-b94e73b69840', '00dbc5c1-99d4-56ab-aa72-e571efbb7283', '82304cb8-d398-5c06-9c0b-1aad6a2c3e61', 'Requisitos atendidos', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '00dbc5c1-99d4-56ab-aa72-e571efbb7283')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '82304cb8-d398-5c06-9c0b-1aad6a2c3e61')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '00dbc5c1-99d4-56ab-aa72-e571efbb7283' AND z.destino_etapa_id = '82304cb8-d398-5c06-9c0b-1aad6a2c3e61' AND z.condicao_transicao = 'Requisitos atendidos');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'f0699ec7-94b9-5b57-ac14-689ab3f8854b', '82304cb8-d398-5c06-9c0b-1aad6a2c3e61', 'd2711ca3-d944-503c-8c5a-f95ba79a21d4', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '82304cb8-d398-5c06-9c0b-1aad6a2c3e61')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd2711ca3-d944-503c-8c5a-f95ba79a21d4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '82304cb8-d398-5c06-9c0b-1aad6a2c3e61' AND z.destino_etapa_id = 'd2711ca3-d944-503c-8c5a-f95ba79a21d4' AND z.condicao_transicao = 'Parecer favorável');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'f49d6836-919a-5e9c-983d-31ce7fab1a1a', 'HAB.RRU.006', 'Regularização fundiária rural', 'Regularizar a ocupação rural individual de interesse social (REURB-S) ou específico (REURB-E).', 'Ocupantes de áreas rurais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LAI, art. 31 (informação pessoal)', ''), '6.0.01.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.01.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.RRU.006');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441', 'f49d6836-919a-5e9c-983d-31ce7fab1a1a', 1, 'SEHAB/REURB', 'Regularização fundiária', 'Cadastrar o ocupante e enquadrar a modalidade.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f49d6836-919a-5e9c-983d-31ce7fab1a1a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f49d6836-919a-5e9c-983d-31ce7fab1a1a' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '65e49439-a1d1-54a5-b38a-adf7339263ad', 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441', 'Requerimento de regularização', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441' AND y.nome_documento = 'Requerimento de regularização');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c27a288a-d0e3-5dba-bb60-85aa6f2c7e61', 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441', 'Cadastro socioeconômico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441' AND y.nome_documento = 'Cadastro socioeconômico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e98e4a83-3801-5216-afb3-86ae2346fa07', 'f49d6836-919a-5e9c-983d-31ce7fab1a1a', 2, 'SEHAB/TOPOGRAFIA', 'Topografia e cartografia', 'Levantar a área e elaborar o memorial descritivo.', 60, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f49d6836-919a-5e9c-983d-31ce7fab1a1a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f49d6836-919a-5e9c-983d-31ce7fab1a1a' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f4a83860-b2ed-556a-8862-7759987dea89', 'e98e4a83-3801-5216-afb3-86ae2346fa07', 'Levantamento topográfico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e98e4a83-3801-5216-afb3-86ae2346fa07')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e98e4a83-3801-5216-afb3-86ae2346fa07' AND y.nome_documento = 'Levantamento topográfico');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '69bb8463-c615-5a51-96f1-816300e90049', 'e98e4a83-3801-5216-afb3-86ae2346fa07', 'Memorial descritivo', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e98e4a83-3801-5216-afb3-86ae2346fa07')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e98e4a83-3801-5216-afb3-86ae2346fa07' AND y.nome_documento = 'Memorial descritivo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '7eefbdec-59a6-548a-bcf5-bbddeb2222d5', 'f49d6836-919a-5e9c-983d-31ce7fab1a1a', 3, 'SEHAB/REURB', 'Regularização fundiária', 'Emitir a certidão de regularização e encaminhar ao cartório.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f49d6836-919a-5e9c-983d-31ce7fab1a1a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f49d6836-919a-5e9c-983d-31ce7fab1a1a' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '022d7781-2469-5020-93c5-c15f6ce095d2', '7eefbdec-59a6-548a-bcf5-bbddeb2222d5', 'Certidão de Regularização Fundiária (CRF)', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7eefbdec-59a6-548a-bcf5-bbddeb2222d5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7eefbdec-59a6-548a-bcf5-bbddeb2222d5' AND y.nome_documento = 'Certidão de Regularização Fundiária (CRF)');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '20a98360-bbfb-507d-b577-2aaaa724e798', 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441', 'e98e4a83-3801-5216-afb3-86ae2346fa07', 'Modalidade enquadrada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e98e4a83-3801-5216-afb3-86ae2346fa07')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'bf0b983d-58b4-5b49-a9ed-f5f1efabf441' AND z.destino_etapa_id = 'e98e4a83-3801-5216-afb3-86ae2346fa07' AND z.condicao_transicao = 'Modalidade enquadrada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '1a66a399-058a-5339-af3f-6287b4f7eef0', 'e98e4a83-3801-5216-afb3-86ae2346fa07', '7eefbdec-59a6-548a-bcf5-bbddeb2222d5', 'Levantamento concluído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e98e4a83-3801-5216-afb3-86ae2346fa07')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '7eefbdec-59a6-548a-bcf5-bbddeb2222d5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e98e4a83-3801-5216-afb3-86ae2346fa07' AND z.destino_etapa_id = '7eefbdec-59a6-548a-bcf5-bbddeb2222d5' AND z.condicao_transicao = 'Levantamento concluído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'adb334c3-e653-5fdf-8794-1629535b381c', 'HAB.DES.007', 'Desmembramento e remembramento de imóvel', 'Aprovar a divisão ou a unificação de lotes urbanos conforme a lei de parcelamento do solo.', 'Proprietários', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '6.0.02.01.10'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.02.01.10' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.DES.007');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5', 'adb334c3-e653-5fdf-8794-1629535b381c', 1, 'SEHAB/PROTOCOLO', 'Atendimento', 'Receber o pedido com o projeto e as matrículas.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'adb334c3-e653-5fdf-8794-1629535b381c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'adb334c3-e653-5fdf-8794-1629535b381c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9aa04ea7-42b7-5043-86ac-3a07d8ef72ba', 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5', 'Requerimento de desmembramento/remembramento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5' AND y.nome_documento = 'Requerimento de desmembramento/remembramento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e96c2f15-b000-54d7-a002-04ce9c6b64c5', 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5', 'Projeto e memorial descritivo', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5' AND y.nome_documento = 'Projeto e memorial descritivo');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c6008acc-832a-535c-b863-3dbe50b76636', 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5', 'Matrículas atualizadas', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5' AND y.nome_documento = 'Matrículas atualizadas');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4b1f3205-8c98-55d6-bb89-ef43a501457e', 'adb334c3-e653-5fdf-8794-1629535b381c', 2, 'SEHAB/CARTOGRAFIA', 'Averbação e cartografia', 'Analisar o projeto e as dimensões mínimas dos lotes.', 20, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'adb334c3-e653-5fdf-8794-1629535b381c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'adb334c3-e653-5fdf-8794-1629535b381c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '82c7b276-56e2-51e2-94a2-9ccc141f4805', '4b1f3205-8c98-55d6-bb89-ef43a501457e', 'Parecer técnico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4b1f3205-8c98-55d6-bb89-ef43a501457e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4b1f3205-8c98-55d6-bb89-ef43a501457e' AND y.nome_documento = 'Parecer técnico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb', 'adb334c3-e653-5fdf-8794-1629535b381c', 3, 'SEHAB/CARTOGRAFIA', 'Averbação e cartografia', 'Emitir a aprovação para registro no cartório.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'adb334c3-e653-5fdf-8794-1629535b381c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'adb334c3-e653-5fdf-8794-1629535b381c' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd8c243ab-54cd-557d-a0b3-50a90a20b7b6', 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb', 'Certidão de aprovação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb' AND y.nome_documento = 'Certidão de aprovação');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '25084d47-4866-5c60-ba18-ba48529c5958', 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5', '4b1f3205-8c98-55d6-bb89-ef43a501457e', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '4b1f3205-8c98-55d6-bb89-ef43a501457e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5' AND z.destino_etapa_id = '4b1f3205-8c98-55d6-bb89-ef43a501457e' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5de6ca12-e21c-509b-aa9b-2ff42f29dee7', '4b1f3205-8c98-55d6-bb89-ef43a501457e', 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb', 'Projeto aprovado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4b1f3205-8c98-55d6-bb89-ef43a501457e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4b1f3205-8c98-55d6-bb89-ef43a501457e' AND z.destino_etapa_id = 'c7716b39-e36c-52b9-9ce1-ba6ae7bb5fcb' AND z.condicao_transicao = 'Projeto aprovado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '3b56929e-4cc5-57f4-a3d9-711c283e3b0b', '4b1f3205-8c98-55d6-bb89-ef43a501457e', 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5', 'Lote abaixo do mínimo', TRUE, NULLIF('Corrigir o projeto', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4b1f3205-8c98-55d6-bb89-ef43a501457e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4b1f3205-8c98-55d6-bb89-ef43a501457e' AND z.destino_etapa_id = 'c2087c8f-fd39-5dcd-b70b-11e8b6ec08b5' AND z.condicao_transicao = 'Lote abaixo do mínimo');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e', 'HAB.ALI.008', 'Alinhamento de lote', 'Emitir laudo e nota de alinhamento do lote em relação ao logradouro.', 'Proprietários e responsáveis técnicos', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '6.0.02.01.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '6.0.02.01.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'HAB.ALI.008');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4765dc29-74de-5b57-ad20-5aa8a8710601', '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e', 1, 'SEHAB/PROTOCOLO', 'Atendimento', 'Receber o pedido com a localização do lote.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f2cd5077-981f-5069-b387-17b1a5efcd03', '4765dc29-74de-5b57-ad20-5aa8a8710601', 'Requerimento de alinhamento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4765dc29-74de-5b57-ad20-5aa8a8710601')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4765dc29-74de-5b57-ad20-5aa8a8710601' AND y.nome_documento = 'Requerimento de alinhamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'edea2925-fb43-5ffa-992c-7877295580f8', '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e', 2, 'SEHAB/TOPOGRAFIA', 'Topografia', 'Vistoriar e emitir o laudo e a nota de alinhamento.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '0151bb9b-5554-5f0c-b2b6-967a05f8cf7e' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '69213401-98ad-522e-ad73-2f6926ce555a', 'edea2925-fb43-5ffa-992c-7877295580f8', 'Laudo e nota de alinhamento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'edea2925-fb43-5ffa-992c-7877295580f8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'edea2925-fb43-5ffa-992c-7877295580f8' AND y.nome_documento = 'Laudo e nota de alinhamento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c9e5d762-e7fc-593c-9811-2d06c0bf31c7', '4765dc29-74de-5b57-ad20-5aa8a8710601', 'edea2925-fb43-5ffa-992c-7877295580f8', 'Pedido recebido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4765dc29-74de-5b57-ad20-5aa8a8710601')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'edea2925-fb43-5ffa-992c-7877295580f8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4765dc29-74de-5b57-ad20-5aa8a8710601' AND z.destino_etapa_id = 'edea2925-fb43-5ffa-992c-7877295580f8' AND z.condicao_transicao = 'Pedido recebido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '259b77ca-ca81-5570-a17f-2b2ba15864e4', 'MAM.FIS.002', 'Fiscalização ambiental', 'Apurar denúncia ou constatação de infração ambiental, com auto de infração e prazo de defesa.', 'Infratores e denunciantes', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.01.00.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.01.00.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'MAM.FIS.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '6963086f-71e9-57da-8920-a340559c9ff3', '259b77ca-ca81-5570-a17f-2b2ba15864e4', 1, 'SEMMA/FISCALIZACAO', 'Fiscalização ambiental', 'Vistoriar o local e lavrar o auto de infração ou de notificação.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '259b77ca-ca81-5570-a17f-2b2ba15864e4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '259b77ca-ca81-5570-a17f-2b2ba15864e4' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '95c1bac2-89be-5f20-9f50-3d0137aad8ee', '6963086f-71e9-57da-8920-a340559c9ff3', 'Relatório de vistoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6963086f-71e9-57da-8920-a340559c9ff3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6963086f-71e9-57da-8920-a340559c9ff3' AND y.nome_documento = 'Relatório de vistoria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '03dfe349-1898-5661-8aca-bcab487d4fb6', '6963086f-71e9-57da-8920-a340559c9ff3', 'Auto de infração ambiental', FALSE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6963086f-71e9-57da-8920-a340559c9ff3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6963086f-71e9-57da-8920-a340559c9ff3' AND y.nome_documento = 'Auto de infração ambiental');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '325d00b6-3579-5ffa-92d6-a191c53eefde', '259b77ca-ca81-5570-a17f-2b2ba15864e4', 2, 'SEMMA/JULGAMENTO', 'Julgamento ambiental', 'Receber a defesa e julgar em primeira instância.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '259b77ca-ca81-5570-a17f-2b2ba15864e4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '259b77ca-ca81-5570-a17f-2b2ba15864e4' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b0b2a490-f4af-508d-9c88-4b8fdc633f1a', '325d00b6-3579-5ffa-92d6-a191c53eefde', 'Defesa do autuado', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '325d00b6-3579-5ffa-92d6-a191c53eefde')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '325d00b6-3579-5ffa-92d6-a191c53eefde' AND y.nome_documento = 'Defesa do autuado');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a11a6f63-5e0a-58f2-b3f2-22f7b82e3d54', '325d00b6-3579-5ffa-92d6-a191c53eefde', 'Decisão de primeira instância', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '325d00b6-3579-5ffa-92d6-a191c53eefde')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '325d00b6-3579-5ffa-92d6-a191c53eefde' AND y.nome_documento = 'Decisão de primeira instância');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e50ac840-5717-515f-8266-c38c908c20df', '6963086f-71e9-57da-8920-a340559c9ff3', '325d00b6-3579-5ffa-92d6-a191c53eefde', 'Auto lavrado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '6963086f-71e9-57da-8920-a340559c9ff3')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '325d00b6-3579-5ffa-92d6-a191c53eefde')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '6963086f-71e9-57da-8920-a340559c9ff3' AND z.destino_etapa_id = '325d00b6-3579-5ffa-92d6-a191c53eefde' AND z.condicao_transicao = 'Auto lavrado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '50eac2ef-1010-5007-bd7c-099c376b4c96', 'MAM.ARV.003', 'Remoção ou supressão de árvores', 'Autorizar a poda, a remoção ou a supressão de vegetação, com compensação ambiental quando exigida.', 'Cidadãos e empreendedores', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.01.02.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.01.02.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'MAM.ARV.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8c4dc70f-1206-5235-b284-3d1d63448b98', '50eac2ef-1010-5007-bd7c-099c376b4c96', 1, 'SEMMA/ATENDIMENTO', 'Atendimento', 'Receber o pedido com a localização e o motivo.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '50eac2ef-1010-5007-bd7c-099c376b4c96')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '50eac2ef-1010-5007-bd7c-099c376b4c96' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'f5061e2e-2d81-5787-959b-133a812ddfc0', '8c4dc70f-1206-5235-b284-3d1d63448b98', 'Requerimento de remoção/supressão', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8c4dc70f-1206-5235-b284-3d1d63448b98')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8c4dc70f-1206-5235-b284-3d1d63448b98' AND y.nome_documento = 'Requerimento de remoção/supressão');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2514a3a6-06e2-57ab-86ee-6606cb2bb787', '8c4dc70f-1206-5235-b284-3d1d63448b98', 'Fotos do local', FALSE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8c4dc70f-1206-5235-b284-3d1d63448b98')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8c4dc70f-1206-5235-b284-3d1d63448b98' AND y.nome_documento = 'Fotos do local');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '99178774-af7f-5a08-b7a8-76f37cf1af7b', '50eac2ef-1010-5007-bd7c-099c376b4c96', 2, 'SEMMA/TECNICO', 'Análise técnica', 'Vistoriar e emitir o laudo e a compensação.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '50eac2ef-1010-5007-bd7c-099c376b4c96')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '50eac2ef-1010-5007-bd7c-099c376b4c96' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7ec5a0fe-b069-5767-87fc-747f3e68c6b1', '99178774-af7f-5a08-b7a8-76f37cf1af7b', 'Laudo técnico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '99178774-af7f-5a08-b7a8-76f37cf1af7b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '99178774-af7f-5a08-b7a8-76f37cf1af7b' AND y.nome_documento = 'Laudo técnico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '538f08cf-2262-5d26-bfd1-870c20234f10', '50eac2ef-1010-5007-bd7c-099c376b4c96', 3, 'SEMMA/LICENCIAMENTO', 'Licenciamento ambiental', 'Emitir a autorização com as condicionantes.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '50eac2ef-1010-5007-bd7c-099c376b4c96')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '50eac2ef-1010-5007-bd7c-099c376b4c96' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '28465d0f-ba69-5790-a12e-e66840b88c47', '538f08cf-2262-5d26-bfd1-870c20234f10', 'Autorização de remoção/supressão', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '538f08cf-2262-5d26-bfd1-870c20234f10')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '538f08cf-2262-5d26-bfd1-870c20234f10' AND y.nome_documento = 'Autorização de remoção/supressão');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '9b5f8a60-fe1a-591d-93a6-8da667e3e042', '538f08cf-2262-5d26-bfd1-870c20234f10', 'Termo de compensação ambiental', FALSE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '538f08cf-2262-5d26-bfd1-870c20234f10')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '538f08cf-2262-5d26-bfd1-870c20234f10' AND y.nome_documento = 'Termo de compensação ambiental');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b39073e5-e6fa-5d2d-9e99-f15cdbfcf026', '8c4dc70f-1206-5235-b284-3d1d63448b98', '99178774-af7f-5a08-b7a8-76f37cf1af7b', 'Pedido recebido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '8c4dc70f-1206-5235-b284-3d1d63448b98')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '99178774-af7f-5a08-b7a8-76f37cf1af7b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '8c4dc70f-1206-5235-b284-3d1d63448b98' AND z.destino_etapa_id = '99178774-af7f-5a08-b7a8-76f37cf1af7b' AND z.condicao_transicao = 'Pedido recebido');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '7177e4b9-0d21-5bef-8f6b-e6a1a5436ab9', '99178774-af7f-5a08-b7a8-76f37cf1af7b', '538f08cf-2262-5d26-bfd1-870c20234f10', 'Remoção autorizada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '99178774-af7f-5a08-b7a8-76f37cf1af7b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '538f08cf-2262-5d26-bfd1-870c20234f10')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '99178774-af7f-5a08-b7a8-76f37cf1af7b' AND z.destino_etapa_id = '538f08cf-2262-5d26-bfd1-870c20234f10' AND z.condicao_transicao = 'Remoção autorizada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'f265e048-a945-5dd2-9ed0-29d7ac3689fa', 'MAM.REC.004', 'Recurso ao Conselho de Meio Ambiente (CONSEMMA)', 'Julgar em segunda instância o recurso contra auto de infração ambiental.', 'Autuados', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.01.03.02'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.01.03.02' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'MAM.REC.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '07cc7522-009b-513c-ba4e-29f3431863fb', 'f265e048-a945-5dd2-9ed0-29d7ac3689fa', 1, 'SEMMA/CONSEMMA', 'Secretaria executiva do CONSEMMA', 'Receber o recurso no prazo e distribuir a um relator.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f265e048-a945-5dd2-9ed0-29d7ac3689fa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f265e048-a945-5dd2-9ed0-29d7ac3689fa' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '885313ea-0c45-5d15-a9b1-afb5759c51d7', '07cc7522-009b-513c-ba4e-29f3431863fb', 'Recurso administrativo', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '07cc7522-009b-513c-ba4e-29f3431863fb')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '07cc7522-009b-513c-ba4e-29f3431863fb' AND y.nome_documento = 'Recurso administrativo');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d', 'f265e048-a945-5dd2-9ed0-29d7ac3689fa', 2, 'SEMMA/CONSEMMA', 'Plenária do CONSEMMA', 'Julgar o recurso e publicar a decisão.', 60, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'f265e048-a945-5dd2-9ed0-29d7ac3689fa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'f265e048-a945-5dd2-9ed0-29d7ac3689fa' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '85a662bc-856a-5137-a760-be5eaf761378', 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d', 'Voto do relator', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d' AND y.nome_documento = 'Voto do relator');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c40fc6eb-9f7a-5103-ada4-fe719ae5563f', 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d', 'Decisão do conselho', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d' AND y.nome_documento = 'Decisão do conselho');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c2e314a6-ac97-5313-addb-431ffec16c57', '07cc7522-009b-513c-ba4e-29f3431863fb', 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d', 'Recurso distribuído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '07cc7522-009b-513c-ba4e-29f3431863fb')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '07cc7522-009b-513c-ba4e-29f3431863fb' AND z.destino_etapa_id = 'e65a569c-cdd2-5c7d-bbcc-825e8146f22d' AND z.condicao_transicao = 'Recurso distribuído');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1591817e-2ede-55c9-be92-6ebf1514980e', 'AGR.SIM.001', 'Registro no Serviço de Inspeção Municipal (SIM)', 'Registrar estabelecimento produtor de alimentos de origem animal no SIM.', 'Produtores rurais e agroindústrias', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.02.01.02'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.02.01.02' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'AGR.SIM.001');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '110de103-6303-56ae-a601-b2be3a652998', '1591817e-2ede-55c9-be92-6ebf1514980e', 1, 'SEMAGRI/SIM', 'Serviço de Inspeção Municipal', 'Receber o pedido com a planta e o memorial do estabelecimento.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1591817e-2ede-55c9-be92-6ebf1514980e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1591817e-2ede-55c9-be92-6ebf1514980e' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '60217aa9-ec9f-555a-a9e0-28cef2b78db1', '110de103-6303-56ae-a601-b2be3a652998', 'Requerimento de registro no SIM', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '110de103-6303-56ae-a601-b2be3a652998')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '110de103-6303-56ae-a601-b2be3a652998' AND y.nome_documento = 'Requerimento de registro no SIM');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c684cf86-3df5-591c-b6d6-3bf7304d1da0', '110de103-6303-56ae-a601-b2be3a652998', 'Planta e memorial técnico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '110de103-6303-56ae-a601-b2be3a652998')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '110de103-6303-56ae-a601-b2be3a652998' AND y.nome_documento = 'Planta e memorial técnico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a5836f9d-1e17-54f3-864d-9135f7cf12a4', '1591817e-2ede-55c9-be92-6ebf1514980e', 2, 'SEMAGRI/SIM', 'Médico veterinário', 'Vistoriar as instalações e emitir o laudo.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1591817e-2ede-55c9-be92-6ebf1514980e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1591817e-2ede-55c9-be92-6ebf1514980e' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6780f0ba-5bb1-595c-b631-025b1d812992', 'a5836f9d-1e17-54f3-864d-9135f7cf12a4', 'Laudo de vistoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4' AND y.nome_documento = 'Laudo de vistoria');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '573f35e4-b6f1-5b26-8e79-6e34877e3617', '1591817e-2ede-55c9-be92-6ebf1514980e', 3, 'SEMAGRI/SIM', 'Serviço de Inspeção Municipal', 'Emitir o certificado de registro.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1591817e-2ede-55c9-be92-6ebf1514980e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1591817e-2ede-55c9-be92-6ebf1514980e' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2736b265-206a-5f60-8661-76fada1910cc', '573f35e4-b6f1-5b26-8e79-6e34877e3617', 'Certificado de registro no SIM', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '573f35e4-b6f1-5b26-8e79-6e34877e3617')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '573f35e4-b6f1-5b26-8e79-6e34877e3617' AND y.nome_documento = 'Certificado de registro no SIM');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '03bdd6bb-a59b-5cf7-965f-a38250fd8405', '110de103-6303-56ae-a601-b2be3a652998', 'a5836f9d-1e17-54f3-864d-9135f7cf12a4', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '110de103-6303-56ae-a601-b2be3a652998')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '110de103-6303-56ae-a601-b2be3a652998' AND z.destino_etapa_id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e8ef3109-8e4b-59b1-9968-a73e2bc17488', 'a5836f9d-1e17-54f3-864d-9135f7cf12a4', '573f35e4-b6f1-5b26-8e79-6e34877e3617', 'Vistoria favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '573f35e4-b6f1-5b26-8e79-6e34877e3617')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4' AND z.destino_etapa_id = '573f35e4-b6f1-5b26-8e79-6e34877e3617' AND z.condicao_transicao = 'Vistoria favorável');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fbfb0402-f6cc-521e-9b9f-4df58c785e04', 'a5836f9d-1e17-54f3-864d-9135f7cf12a4', '110de103-6303-56ae-a601-b2be3a652998', 'Instalações inadequadas', TRUE, NULLIF('Adequar as instalações', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '110de103-6303-56ae-a601-b2be3a652998')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a5836f9d-1e17-54f3-864d-9135f7cf12a4' AND z.destino_etapa_id = '110de103-6303-56ae-a601-b2be3a652998' AND z.condicao_transicao = 'Instalações inadequadas');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'bd6c910e-6e68-5a34-a26f-de44d8d241ac', 'AGR.TRA.002', 'Serviço de trator e implementos ao produtor rural', 'Agendar e executar serviço de mecanização agrícola para o pequeno produtor.', 'Pequenos produtores rurais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.02.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.02.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'AGR.TRA.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd68519e7-f26c-5d55-83ea-c0d13cebb40a', 'bd6c910e-6e68-5a34-a26f-de44d8d241ac', 1, 'SEMAGRI/ATENDIMENTO', 'Atendimento ao produtor', 'Cadastrar o pedido e conferir o enquadramento do produtor.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'bd6c910e-6e68-5a34-a26f-de44d8d241ac')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'bd6c910e-6e68-5a34-a26f-de44d8d241ac' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '516d75d7-3d56-543f-a946-b84659d6a149', 'd68519e7-f26c-5d55-83ea-c0d13cebb40a', 'Requerimento de serviço de trator', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd68519e7-f26c-5d55-83ea-c0d13cebb40a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd68519e7-f26c-5d55-83ea-c0d13cebb40a' AND y.nome_documento = 'Requerimento de serviço de trator');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'f1533fdd-6f00-50e9-9323-c557702a70c0', 'bd6c910e-6e68-5a34-a26f-de44d8d241ac', 2, 'SEMAGRI/MECANIZACAO', 'Patrulha mecanizada', 'Agendar, executar e registrar as horas trabalhadas.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'bd6c910e-6e68-5a34-a26f-de44d8d241ac')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'bd6c910e-6e68-5a34-a26f-de44d8d241ac' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '711fca61-6311-560d-9b3b-0ed514f1d1b0', 'f1533fdd-6f00-50e9-9323-c557702a70c0', 'Ordem de serviço', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'f1533fdd-6f00-50e9-9323-c557702a70c0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'f1533fdd-6f00-50e9-9323-c557702a70c0' AND y.nome_documento = 'Ordem de serviço');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cde456e3-1cde-5905-8f6d-d4e7467caa92', 'f1533fdd-6f00-50e9-9323-c557702a70c0', 'Termo de recebimento do serviço', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'f1533fdd-6f00-50e9-9323-c557702a70c0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'f1533fdd-6f00-50e9-9323-c557702a70c0' AND y.nome_documento = 'Termo de recebimento do serviço');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '49b3925c-b093-54d3-a6c7-e864f9164e1c', 'd68519e7-f26c-5d55-83ea-c0d13cebb40a', 'f1533fdd-6f00-50e9-9323-c557702a70c0', 'Produtor enquadrado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'd68519e7-f26c-5d55-83ea-c0d13cebb40a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'f1533fdd-6f00-50e9-9323-c557702a70c0')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'd68519e7-f26c-5d55-83ea-c0d13cebb40a' AND z.destino_etapa_id = 'f1533fdd-6f00-50e9-9323-c557702a70c0' AND z.condicao_transicao = 'Produtor enquadrado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '36633d7d-9432-5f23-bb16-dbd18e83cd8e', 'AGR.POC.003', 'Perfuração de poço artesiano em área rural', 'Avaliar e executar a perfuração de poço para abastecimento de propriedade rural.', 'Produtores rurais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.02.03.02'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.02.03.02' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'AGR.POC.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'bde652d0-bd16-5f96-baa0-f2a03ce44392', '36633d7d-9432-5f23-bb16-dbd18e83cd8e', 1, 'SEMAGRI/ATENDIMENTO', 'Atendimento ao produtor', 'Receber o pedido e a documentação da propriedade.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '36633d7d-9432-5f23-bb16-dbd18e83cd8e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '36633d7d-9432-5f23-bb16-dbd18e83cd8e' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2991a306-2dac-5752-887a-f3fddd3975cf', 'bde652d0-bd16-5f96-baa0-f2a03ce44392', 'Requerimento de perfuração de poço', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bde652d0-bd16-5f96-baa0-f2a03ce44392')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bde652d0-bd16-5f96-baa0-f2a03ce44392' AND y.nome_documento = 'Requerimento de perfuração de poço');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '34b252d8-6d2b-55f9-a5cc-0a7935bd7f06', 'bde652d0-bd16-5f96-baa0-f2a03ce44392', 'Documento da propriedade', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'bde652d0-bd16-5f96-baa0-f2a03ce44392')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'bde652d0-bd16-5f96-baa0-f2a03ce44392' AND y.nome_documento = 'Documento da propriedade');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c40b2adf-a85a-550d-91e9-ca2b24b71071', '36633d7d-9432-5f23-bb16-dbd18e83cd8e', 2, 'SEMAGRI/EXTENSAO', 'Extensão rural', 'Avaliar a viabilidade técnica e a outorga necessária.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '36633d7d-9432-5f23-bb16-dbd18e83cd8e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '36633d7d-9432-5f23-bb16-dbd18e83cd8e' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '402b0b3d-7a8a-5ad0-b967-eeefdaf46a54', 'c40b2adf-a85a-550d-91e9-ca2b24b71071', 'Parecer técnico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c40b2adf-a85a-550d-91e9-ca2b24b71071')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c40b2adf-a85a-550d-91e9-ca2b24b71071' AND y.nome_documento = 'Parecer técnico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '812a48f1-5de1-5cd2-806a-0b6cdcde691c', '36633d7d-9432-5f23-bb16-dbd18e83cd8e', 3, 'SEMAGRI/EXTENSAO', 'Extensão rural', 'Executar ou contratar a perfuração e registrar a entrega.', 60, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '36633d7d-9432-5f23-bb16-dbd18e83cd8e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '36633d7d-9432-5f23-bb16-dbd18e83cd8e' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '87e9c18f-e34e-5b70-a340-3c61d3800ea7', '812a48f1-5de1-5cd2-806a-0b6cdcde691c', 'Termo de entrega do poço', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '812a48f1-5de1-5cd2-806a-0b6cdcde691c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '812a48f1-5de1-5cd2-806a-0b6cdcde691c' AND y.nome_documento = 'Termo de entrega do poço');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '5cb12a75-813a-5bad-b458-8a643e62e78c', 'bde652d0-bd16-5f96-baa0-f2a03ce44392', 'c40b2adf-a85a-550d-91e9-ca2b24b71071', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'bde652d0-bd16-5f96-baa0-f2a03ce44392')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c40b2adf-a85a-550d-91e9-ca2b24b71071')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'bde652d0-bd16-5f96-baa0-f2a03ce44392' AND z.destino_etapa_id = 'c40b2adf-a85a-550d-91e9-ca2b24b71071' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '2ef9b1c4-71b3-562a-8c3a-c9261a76ed99', 'c40b2adf-a85a-550d-91e9-ca2b24b71071', '812a48f1-5de1-5cd2-806a-0b6cdcde691c', 'Perfuração viável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c40b2adf-a85a-550d-91e9-ca2b24b71071')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '812a48f1-5de1-5cd2-806a-0b6cdcde691c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c40b2adf-a85a-550d-91e9-ca2b24b71071' AND z.destino_etapa_id = '812a48f1-5de1-5cd2-806a-0b6cdcde691c' AND z.condicao_transicao = 'Perfuração viável');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709', 'AGR.EST.004', 'Manutenção de estrada rural', 'Atender pedido de manutenção ou recuperação de estrada vicinal.', 'Comunidades rurais', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '8.0.02.04.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '8.0.02.04.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'AGR.EST.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3257886b-4c3d-559b-a854-51dc37b3b26f', '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709', 1, 'SEMAGRI/ATENDIMENTO', 'Atendimento', 'Registrar o pedido com a localização do trecho.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '35814328-cf68-57be-8db3-09f685c91d53', '3257886b-4c3d-559b-a854-51dc37b3b26f', 'Requerimento de manutenção de estrada', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3257886b-4c3d-559b-a854-51dc37b3b26f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3257886b-4c3d-559b-a854-51dc37b3b26f' AND y.nome_documento = 'Requerimento de manutenção de estrada');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71', '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709', 2, 'SEMAGRI/ESTRADAS', 'Manutenção de estradas', 'Vistoriar, programar e executar o serviço.', 45, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9178ae67-80d5-5d9b-9f13-0dc4a6bb6709' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '686711f7-0b20-5844-bcd4-cf7ef53f3395', 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71', 'Relatório de vistoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71' AND y.nome_documento = 'Relatório de vistoria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'b2852e17-8b18-54c5-bd5c-8f937c017afa', 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71', 'Relatório de execução', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71' AND y.nome_documento = 'Relatório de execução');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fb08acba-6b4d-517f-b97b-4469042cc0a6', '3257886b-4c3d-559b-a854-51dc37b3b26f', 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71', 'Pedido registrado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '3257886b-4c3d-559b-a854-51dc37b3b26f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '3257886b-4c3d-559b-a854-51dc37b3b26f' AND z.destino_etapa_id = 'c8db07e7-c5d3-5f62-8e37-f78d4263fe71' AND z.condicao_transicao = 'Pedido registrado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '074ebb47-3607-5e68-bcc2-8d208a8e3d83', 'DEC.INC.002', 'Incentivo fiscal a empresa', 'Conceder incentivo fiscal previsto na lei municipal de incentivo, com contrapartidas e acompanhamento.', 'Empresas', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '9.0.01.00.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '9.0.01.00.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'DEC.INC.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ee45a510-f1e9-523c-8c52-2bef77e646a1', '074ebb47-3607-5e68-bcc2-8d208a8e3d83', 1, 'SEDEC/ATENDIMENTO', 'Sala do Empreendedor', 'Receber o pedido com o projeto de investimento.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '074ebb47-3607-5e68-bcc2-8d208a8e3d83')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '074ebb47-3607-5e68-bcc2-8d208a8e3d83' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '0e006106-5fa6-532e-9c88-cc07e80c9fcb', 'ee45a510-f1e9-523c-8c52-2bef77e646a1', 'Requerimento de incentivo fiscal', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ee45a510-f1e9-523c-8c52-2bef77e646a1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ee45a510-f1e9-523c-8c52-2bef77e646a1' AND y.nome_documento = 'Requerimento de incentivo fiscal');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '34fe4944-bec5-51af-b9f8-114a0f785d00', 'ee45a510-f1e9-523c-8c52-2bef77e646a1', 'Projeto de investimento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ee45a510-f1e9-523c-8c52-2bef77e646a1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ee45a510-f1e9-523c-8c52-2bef77e646a1' AND y.nome_documento = 'Projeto de investimento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '027ab4ad-e274-540b-a435-7f01d4716c83', '074ebb47-3607-5e68-bcc2-8d208a8e3d83', 2, 'SEDEC/CONSELHO', 'Conselho de desenvolvimento', 'Analisar o projeto e as contrapartidas.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '074ebb47-3607-5e68-bcc2-8d208a8e3d83')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '074ebb47-3607-5e68-bcc2-8d208a8e3d83' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e8cc4c9d-5a2c-5efa-bdbf-52dddc7367be', '027ab4ad-e274-540b-a435-7f01d4716c83', 'Parecer do conselho', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '027ab4ad-e274-540b-a435-7f01d4716c83')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '027ab4ad-e274-540b-a435-7f01d4716c83' AND y.nome_documento = 'Parecer do conselho');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3', '074ebb47-3607-5e68-bcc2-8d208a8e3d83', 3, 'SEREC/TRIBUTARIO', 'Julgamento tributário', 'Conceder o benefício e registrar as condições.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '074ebb47-3607-5e68-bcc2-8d208a8e3d83')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '074ebb47-3607-5e68-bcc2-8d208a8e3d83' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4d527635-f886-5dcb-943e-926e19dba6d5', 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3', 'Termo de concessão de incentivo', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3' AND y.nome_documento = 'Termo de concessão de incentivo');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b6fd0b22-7436-5935-ab55-4712960f4a26', 'ee45a510-f1e9-523c-8c52-2bef77e646a1', '027ab4ad-e274-540b-a435-7f01d4716c83', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'ee45a510-f1e9-523c-8c52-2bef77e646a1')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '027ab4ad-e274-540b-a435-7f01d4716c83')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'ee45a510-f1e9-523c-8c52-2bef77e646a1' AND z.destino_etapa_id = '027ab4ad-e274-540b-a435-7f01d4716c83' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '60e01ae5-95bd-5308-a4d5-86ac6ec5a7ad', '027ab4ad-e274-540b-a435-7f01d4716c83', 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3', 'Projeto aprovado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '027ab4ad-e274-540b-a435-7f01d4716c83')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '027ab4ad-e274-540b-a435-7f01d4716c83' AND z.destino_etapa_id = 'd70175b8-34cd-51f8-81b7-ec9c62d5e4d3' AND z.condicao_transicao = 'Projeto aprovado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '9ad64793-e4b8-5a5b-b148-089f53baf966', 'DEC.IND.003', 'Concessão de área em distrito industrial', 'Conceder área em distrito industrial a empresa, conforme o roteiro do conselho diretor (CODIP), com escritura após a implantação.', 'Empresas', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '9.0.01.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '9.0.01.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'DEC.IND.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9ff2bfd0-a613-580c-ac16-e401c20c7d7f', '9ad64793-e4b8-5a5b-b148-089f53baf966', 1, 'SEDEC/INDUSTRIA', 'Indústria', 'Receber o pedido com o projeto e o cronograma de implantação.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9ad64793-e4b8-5a5b-b148-089f53baf966')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9ad64793-e4b8-5a5b-b148-089f53baf966' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e0afa853-9c30-5559-8251-f36769088a15', '9ff2bfd0-a613-580c-ac16-e401c20c7d7f', 'Requerimento de área industrial', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9ff2bfd0-a613-580c-ac16-e401c20c7d7f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9ff2bfd0-a613-580c-ac16-e401c20c7d7f' AND y.nome_documento = 'Requerimento de área industrial');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '58bea8b9-71e7-5059-b83e-337dec75f8b7', '9ff2bfd0-a613-580c-ac16-e401c20c7d7f', 'Projeto de implantação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9ff2bfd0-a613-580c-ac16-e401c20c7d7f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9ff2bfd0-a613-580c-ac16-e401c20c7d7f' AND y.nome_documento = 'Projeto de implantação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b', '9ad64793-e4b8-5a5b-b148-089f53baf966', 2, 'SEDEC/CODIP', 'Conselho Diretor (CODIP)', 'Deliberar a concessão e as obrigações.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9ad64793-e4b8-5a5b-b148-089f53baf966')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9ad64793-e4b8-5a5b-b148-089f53baf966' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3f7b7e22-e5c1-51a4-885e-5215e17e8f10', 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b', 'Resolução do CODIP', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b' AND y.nome_documento = 'Resolução do CODIP');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '6d1163f6-e893-5610-bdfa-2524b583c5e9', '9ad64793-e4b8-5a5b-b148-089f53baf966', 3, 'PGM', 'Procuradoria', 'Emitir parecer jurídico.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9ad64793-e4b8-5a5b-b148-089f53baf966')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9ad64793-e4b8-5a5b-b148-089f53baf966' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '609b3e63-0206-59b5-9b9b-b3e573af49ad', '6d1163f6-e893-5610-bdfa-2524b583c5e9', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6d1163f6-e893-5610-bdfa-2524b583c5e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6d1163f6-e893-5610-bdfa-2524b583c5e9' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '2adfc859-673e-515f-b49b-3970e470e768', '9ad64793-e4b8-5a5b-b148-089f53baf966', 4, 'SEDEC/INDUSTRIA', 'Indústria', 'Firmar o contrato e acompanhar a implantação até a escritura.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '9ad64793-e4b8-5a5b-b148-089f53baf966')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '9ad64793-e4b8-5a5b-b148-089f53baf966' AND x.ordem = 4);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '09c8f86f-867a-5050-855b-d25cbd6f1404', '2adfc859-673e-515f-b49b-3970e470e768', 'Contrato de concessão', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '2adfc859-673e-515f-b49b-3970e470e768')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '2adfc859-673e-515f-b49b-3970e470e768' AND y.nome_documento = 'Contrato de concessão');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '74105c79-94f9-5b71-9227-a4f0591a9896', '9ff2bfd0-a613-580c-ac16-e401c20c7d7f', 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '9ff2bfd0-a613-580c-ac16-e401c20c7d7f')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '9ff2bfd0-a613-580c-ac16-e401c20c7d7f' AND z.destino_etapa_id = 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'fe3992f0-fd4d-5923-a000-22787b55d18c', 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b', '6d1163f6-e893-5610-bdfa-2524b583c5e9', 'Concessão aprovada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '6d1163f6-e893-5610-bdfa-2524b583c5e9')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c293d109-f4b3-57a4-a6fc-7c7a61fe220b' AND z.destino_etapa_id = '6d1163f6-e893-5610-bdfa-2524b583c5e9' AND z.condicao_transicao = 'Concessão aprovada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c687f3d9-4a4b-5a93-a99f-72fc87a2e63b', '6d1163f6-e893-5610-bdfa-2524b583c5e9', '2adfc859-673e-515f-b49b-3970e470e768', 'Contrato aprovado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '6d1163f6-e893-5610-bdfa-2524b583c5e9')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '2adfc859-673e-515f-b49b-3970e470e768')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '6d1163f6-e893-5610-bdfa-2524b583c5e9' AND z.destino_etapa_id = '2adfc859-673e-515f-b49b-3970e470e768' AND z.condicao_transicao = 'Contrato aprovado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '835ba7a3-34d5-5871-b093-79fc83d1f265', 'SAS.CNV.002', 'Convênio ou parceria na assistência social', 'Firmar convênio ou parceria com outro ente ou com organização para executar ação socioassistencial.', 'Entes e organizações parceiras', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '10.0.01.00.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '10.0.01.00.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAS.CNV.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '61c9edef-5cb5-5695-9392-e88379a91ee3', '835ba7a3-34d5-5871-b093-79fc83d1f265', 1, 'SEMPRAS/GESTAO', 'Gestão do SUAS', 'Elaborar o plano de trabalho e a justificativa.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '835ba7a3-34d5-5871-b093-79fc83d1f265')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '835ba7a3-34d5-5871-b093-79fc83d1f265' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '1b667965-dd21-51f8-9e68-602365f4902a', '61c9edef-5cb5-5695-9392-e88379a91ee3', 'Plano de trabalho', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '61c9edef-5cb5-5695-9392-e88379a91ee3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '61c9edef-5cb5-5695-9392-e88379a91ee3' AND y.nome_documento = 'Plano de trabalho');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'ce140a41-978f-57fb-9d4c-ee9feb515dfa', '61c9edef-5cb5-5695-9392-e88379a91ee3', 'Justificativa da parceria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '61c9edef-5cb5-5695-9392-e88379a91ee3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '61c9edef-5cb5-5695-9392-e88379a91ee3' AND y.nome_documento = 'Justificativa da parceria');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da', '835ba7a3-34d5-5871-b093-79fc83d1f265', 2, 'PGM', 'Procuradoria', 'Emitir parecer jurídico.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '835ba7a3-34d5-5871-b093-79fc83d1f265')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '835ba7a3-34d5-5871-b093-79fc83d1f265' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '7f2cfaa2-3208-59a6-99a9-a1286dc1027a', 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '672537cd-d0c4-52c1-9144-95d7d353aefd', '835ba7a3-34d5-5871-b093-79fc83d1f265', 3, 'GAB', 'Prefeito', 'Assinar o termo e publicar o extrato.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '835ba7a3-34d5-5871-b093-79fc83d1f265')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '835ba7a3-34d5-5871-b093-79fc83d1f265' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'dd3ae33a-a111-5e22-833a-86a611f3af3a', '672537cd-d0c4-52c1-9144-95d7d353aefd', 'Termo de convênio/parceria', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '672537cd-d0c4-52c1-9144-95d7d353aefd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '672537cd-d0c4-52c1-9144-95d7d353aefd' AND y.nome_documento = 'Termo de convênio/parceria');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '988ed398-97dd-53ce-8305-fb1caa2597e3', '672537cd-d0c4-52c1-9144-95d7d353aefd', 'Extrato publicado', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '672537cd-d0c4-52c1-9144-95d7d353aefd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '672537cd-d0c4-52c1-9144-95d7d353aefd' AND y.nome_documento = 'Extrato publicado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'afc52086-ecf3-50c1-8e72-f683799b6d22', '61c9edef-5cb5-5695-9392-e88379a91ee3', 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da', 'Plano aprovado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '61c9edef-5cb5-5695-9392-e88379a91ee3')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '61c9edef-5cb5-5695-9392-e88379a91ee3' AND z.destino_etapa_id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da' AND z.condicao_transicao = 'Plano aprovado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '9ac46315-fbff-5b61-bd4a-669dbe9dddd9', 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da', '672537cd-d0c4-52c1-9144-95d7d353aefd', 'Minuta aprovada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '672537cd-d0c4-52c1-9144-95d7d353aefd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da' AND z.destino_etapa_id = '672537cd-d0c4-52c1-9144-95d7d353aefd' AND z.condicao_transicao = 'Minuta aprovada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '06342b63-2617-525d-a58f-272e30e23a83', 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da', '61c9edef-5cb5-5695-9392-e88379a91ee3', 'Ajustes na minuta', TRUE, NULLIF('Ajustar o plano ou a minuta', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '61c9edef-5cb5-5695-9392-e88379a91ee3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'c9fd1575-ae9a-5829-8bde-9e21c36ed8da' AND z.destino_etapa_id = '61c9edef-5cb5-5695-9392-e88379a91ee3' AND z.condicao_transicao = 'Ajustes na minuta');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '3fcf5c11-fd69-5b41-92dc-0308b9e46786', 'SAS.CHA.003', 'Chamamento público de OSC', 'Selecionar organização da sociedade civil por chamamento público (Lei 13.019/2014).', 'Organizações da sociedade civil', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '10.0.03.00.06'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '10.0.03.00.06' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAS.CHA.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '996f6c22-df2c-5d00-8717-13d0707b308e', '3fcf5c11-fd69-5b41-92dc-0308b9e46786', 1, 'SEMPRAS/GESTAO', 'Gestão do SUAS', 'Elaborar e publicar o edital de chamamento.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '3fcf5c11-fd69-5b41-92dc-0308b9e46786')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '3fcf5c11-fd69-5b41-92dc-0308b9e46786' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '49159230-41c0-5d37-a13c-ad6e818b900a', '996f6c22-df2c-5d00-8717-13d0707b308e', 'Edital de chamamento público', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '996f6c22-df2c-5d00-8717-13d0707b308e')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '996f6c22-df2c-5d00-8717-13d0707b308e' AND y.nome_documento = 'Edital de chamamento público');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b6075084-f819-582e-b781-8e58e70a8417', '3fcf5c11-fd69-5b41-92dc-0308b9e46786', 2, 'SEMPRAS/COMISSAO', 'Comissão de seleção', 'Avaliar as propostas e publicar o resultado.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '3fcf5c11-fd69-5b41-92dc-0308b9e46786')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '3fcf5c11-fd69-5b41-92dc-0308b9e46786' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '76c17f6e-bf98-566e-a54e-9c08895b34dd', 'b6075084-f819-582e-b781-8e58e70a8417', 'Ata de julgamento', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b6075084-f819-582e-b781-8e58e70a8417')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b6075084-f819-582e-b781-8e58e70a8417' AND y.nome_documento = 'Ata de julgamento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '18ce506c-a743-54d9-95ef-9965a975504a', 'b6075084-f819-582e-b781-8e58e70a8417', 'Resultado do chamamento', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b6075084-f819-582e-b781-8e58e70a8417')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b6075084-f819-582e-b781-8e58e70a8417' AND y.nome_documento = 'Resultado do chamamento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c', '3fcf5c11-fd69-5b41-92dc-0308b9e46786', 3, 'SEMPRAS/GESTAO', 'Gestão do SUAS', 'Firmar o termo de colaboração ou fomento.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '3fcf5c11-fd69-5b41-92dc-0308b9e46786')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '3fcf5c11-fd69-5b41-92dc-0308b9e46786' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '70c629a4-f478-5919-8243-7fff22b60070', '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c', 'Termo de colaboração/fomento', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c' AND y.nome_documento = 'Termo de colaboração/fomento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '59aa90f1-1c8b-51cb-8155-e4257419563c', '996f6c22-df2c-5d00-8717-13d0707b308e', 'b6075084-f819-582e-b781-8e58e70a8417', 'Edital publicado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '996f6c22-df2c-5d00-8717-13d0707b308e')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b6075084-f819-582e-b781-8e58e70a8417')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '996f6c22-df2c-5d00-8717-13d0707b308e' AND z.destino_etapa_id = 'b6075084-f819-582e-b781-8e58e70a8417' AND z.condicao_transicao = 'Edital publicado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '089be9a5-fadd-5c0b-9fb5-bdffa94eba77', 'b6075084-f819-582e-b781-8e58e70a8417', '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c', 'Resultado homologado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'b6075084-f819-582e-b781-8e58e70a8417')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'b6075084-f819-582e-b781-8e58e70a8417' AND z.destino_etapa_id = '8c79ef99-a5ed-5d28-a27b-ed6dd818a53c' AND z.condicao_transicao = 'Resultado homologado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '53c9d095-dec1-5ee7-b147-75df5c9c297c', 'SAS.BIL.004', 'Bilhete único especial para pessoa com deficiência', 'Conceder a gratuidade no transporte coletivo municipal à pessoa com deficiência e acompanhante, quando indicado.', 'Pessoas com deficiência', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '10.0.06.00.40'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '10.0.06.00.40' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAS.BIL.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'a93345d6-1c80-5a82-a76c-073f92b95947', '53c9d095-dec1-5ee7-b147-75df5c9c297c', 1, 'SEMPRAS/ATENDIMENTO', 'Atendimento', 'Receber o pedido com o laudo médico e os documentos.', 2, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '53c9d095-dec1-5ee7-b147-75df5c9c297c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '53c9d095-dec1-5ee7-b147-75df5c9c297c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '76eca797-d014-5f30-ab3c-3a0b180e1aae', 'a93345d6-1c80-5a82-a76c-073f92b95947', 'Requerimento de bilhete único especial', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a93345d6-1c80-5a82-a76c-073f92b95947')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a93345d6-1c80-5a82-a76c-073f92b95947' AND y.nome_documento = 'Requerimento de bilhete único especial');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd0800e1b-845b-56cb-9a5f-9d832b25abff', 'a93345d6-1c80-5a82-a76c-073f92b95947', 'Laudo médico', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'a93345d6-1c80-5a82-a76c-073f92b95947')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'a93345d6-1c80-5a82-a76c-073f92b95947' AND y.nome_documento = 'Laudo médico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4', '53c9d095-dec1-5ee7-b147-75df5c9c297c', 2, 'SEMPRAS/AVALIACAO', 'Avaliação social', 'Avaliar os requisitos e a necessidade de acompanhante.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '53c9d095-dec1-5ee7-b147-75df5c9c297c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '53c9d095-dec1-5ee7-b147-75df5c9c297c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd4ff7702-b259-5b59-811a-d72350842a7d', 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4', 'Parecer social', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4' AND y.nome_documento = 'Parecer social');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd99860bb-3807-5811-854b-3748c92c8e0d', '53c9d095-dec1-5ee7-b147-75df5c9c297c', 3, 'SEMPRAS/ATENDIMENTO', 'Atendimento', 'Solicitar a emissão do cartão e entregar ao beneficiário.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '53c9d095-dec1-5ee7-b147-75df5c9c297c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '53c9d095-dec1-5ee7-b147-75df5c9c297c' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '15ff5489-f9e1-5873-bbb5-5b074f045873', 'd99860bb-3807-5811-854b-3748c92c8e0d', 'Termo de entrega do cartão', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd99860bb-3807-5811-854b-3748c92c8e0d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd99860bb-3807-5811-854b-3748c92c8e0d' AND y.nome_documento = 'Termo de entrega do cartão');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c72236af-17a7-5cc8-93c1-3e0e96da3087', 'a93345d6-1c80-5a82-a76c-073f92b95947', 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4', 'Pedido instruído', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'a93345d6-1c80-5a82-a76c-073f92b95947')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'a93345d6-1c80-5a82-a76c-073f92b95947' AND z.destino_etapa_id = 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4' AND z.condicao_transicao = 'Pedido instruído');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'b530e6d2-bbd2-5148-bc7e-6de381dc244d', 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4', 'd99860bb-3807-5811-854b-3748c92c8e0d', 'Benefício deferido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd99860bb-3807-5811-854b-3748c92c8e0d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e98e15fc-d497-5ac1-9c0f-648a2e9d7bf4' AND z.destino_etapa_id = 'd99860bb-3807-5811-854b-3748c92c8e0d' AND z.condicao_transicao = 'Benefício deferido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504', 'SAU.EMD.002', 'Emenda parlamentar para a saúde', 'Cadastrar, executar e prestar contas de recurso de emenda parlamentar destinado à saúde.', 'Secretaria de Saúde', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.01.00.13'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.01.00.13' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.EMD.002');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '72220f42-4694-5545-94ef-be73dbc72645', 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504', 1, 'SMS/PLANEJAMENTO', 'Planejamento em saúde', 'Cadastrar a proposta no sistema do Ministério e acompanhar a aprovação.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5653b23c-263e-54ba-8e03-e5baa08a9c6b', '72220f42-4694-5545-94ef-be73dbc72645', 'Proposta de aplicação', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '72220f42-4694-5545-94ef-be73dbc72645')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '72220f42-4694-5545-94ef-be73dbc72645' AND y.nome_documento = 'Proposta de aplicação');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'b0df6906-7d90-5f22-b262-b0ae296709c4', 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504', 2, 'SMS/FUNDO', 'Fundo Municipal de Saúde', 'Executar a despesa conforme o objeto.', 180, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4189e2d8-a110-595c-bb27-acdc92c642ee', 'b0df6906-7d90-5f22-b262-b0ae296709c4', 'Processo de execução da despesa', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'b0df6906-7d90-5f22-b262-b0ae296709c4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'b0df6906-7d90-5f22-b262-b0ae296709c4' AND y.nome_documento = 'Processo de execução da despesa');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '9c2feeaf-4777-5f46-9f02-fe72e8ada05a', 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504', 3, 'SMS/FUNDO', 'Fundo Municipal de Saúde', 'Prestar contas no sistema federal.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'fcad734c-b9f1-586b-ae6d-fc3eeecf9504' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5e37803c-bc37-5318-ab8e-e411306e5a43', '9c2feeaf-4777-5f46-9f02-fe72e8ada05a', 'Prestação de contas', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '9c2feeaf-4777-5f46-9f02-fe72e8ada05a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '9c2feeaf-4777-5f46-9f02-fe72e8ada05a' AND y.nome_documento = 'Prestação de contas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c3e03367-bbd8-5d3c-8c6e-105cff1bf145', '72220f42-4694-5545-94ef-be73dbc72645', 'b0df6906-7d90-5f22-b262-b0ae296709c4', 'Recurso creditado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '72220f42-4694-5545-94ef-be73dbc72645')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'b0df6906-7d90-5f22-b262-b0ae296709c4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '72220f42-4694-5545-94ef-be73dbc72645' AND z.destino_etapa_id = 'b0df6906-7d90-5f22-b262-b0ae296709c4' AND z.condicao_transicao = 'Recurso creditado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'c12d41bb-a2fa-562c-b181-9f975565cd10', 'b0df6906-7d90-5f22-b262-b0ae296709c4', '9c2feeaf-4777-5f46-9f02-fe72e8ada05a', 'Objeto executado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'b0df6906-7d90-5f22-b262-b0ae296709c4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '9c2feeaf-4777-5f46-9f02-fe72e8ada05a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'b0df6906-7d90-5f22-b262-b0ae296709c4' AND z.destino_etapa_id = '9c2feeaf-4777-5f46-9f02-fe72e8ada05a' AND z.condicao_transicao = 'Objeto executado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '916afa91-22de-5564-8958-96d56c21645a', 'SAU.CNE.003', 'Cadastro de estabelecimento no SCNES', 'Incluir, alterar ou excluir estabelecimento e profissionais no Cadastro Nacional de Estabelecimentos de Saúde.', 'Estabelecimentos de saúde', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.01.01.09'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.01.01.09' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.CNE.003');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '696d6ea2-8951-5e41-a3ec-84d7e07bd965', '916afa91-22de-5564-8958-96d56c21645a', 1, 'SMS/CNES', 'Cadastro CNES', 'Receber a ficha cadastral e os documentos do estabelecimento.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '916afa91-22de-5564-8958-96d56c21645a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '916afa91-22de-5564-8958-96d56c21645a' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6b16ecfb-afe0-54b6-a55a-399dad7627fb', '696d6ea2-8951-5e41-a3ec-84d7e07bd965', 'Ficha cadastral do CNES', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '696d6ea2-8951-5e41-a3ec-84d7e07bd965')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '696d6ea2-8951-5e41-a3ec-84d7e07bd965' AND y.nome_documento = 'Ficha cadastral do CNES');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bf9293fe-7e4d-5430-b8d7-0a8fd9efaffb', '696d6ea2-8951-5e41-a3ec-84d7e07bd965', 'Licença sanitária', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '696d6ea2-8951-5e41-a3ec-84d7e07bd965')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '696d6ea2-8951-5e41-a3ec-84d7e07bd965' AND y.nome_documento = 'Licença sanitária');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fe704281-fd94-5d9e-b956-631ad3c29ce4', '916afa91-22de-5564-8958-96d56c21645a', 2, 'SMS/CNES', 'Cadastro CNES', 'Atualizar o SCNES e enviar a competência.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '916afa91-22de-5564-8958-96d56c21645a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '916afa91-22de-5564-8958-96d56c21645a' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5802354c-dc72-59b4-a208-207151aeb568', 'fe704281-fd94-5d9e-b956-631ad3c29ce4', 'Comprovante de atualização', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fe704281-fd94-5d9e-b956-631ad3c29ce4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fe704281-fd94-5d9e-b956-631ad3c29ce4' AND y.nome_documento = 'Comprovante de atualização');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '15c46c80-74e9-5136-a790-edf45d5f0d7f', '696d6ea2-8951-5e41-a3ec-84d7e07bd965', 'fe704281-fd94-5d9e-b956-631ad3c29ce4', 'Ficha conferida', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '696d6ea2-8951-5e41-a3ec-84d7e07bd965')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'fe704281-fd94-5d9e-b956-631ad3c29ce4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '696d6ea2-8951-5e41-a3ec-84d7e07bd965' AND z.destino_etapa_id = 'fe704281-fd94-5d9e-b956-631ad3c29ce4' AND z.condicao_transicao = 'Ficha conferida');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c', 'SAU.EST.004', 'Convênio de estágio nas unidades de saúde', 'Firmar convênio com instituição de ensino para estágio supervisionado nas unidades de saúde.', 'Instituições de ensino', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.01.02.00'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.01.02.00' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.EST.004');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e0142aee-e004-51d7-92e5-3cd9c497aba8', 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c', 1, 'SMS/EDUCACAO', 'Educação permanente', 'Receber a proposta e verificar as vagas nas unidades.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd0cde3eb-7965-5fc9-b0d3-5faae5c7036f', 'e0142aee-e004-51d7-92e5-3cd9c497aba8', 'Proposta de convênio de estágio', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e0142aee-e004-51d7-92e5-3cd9c497aba8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e0142aee-e004-51d7-92e5-3cd9c497aba8' AND y.nome_documento = 'Proposta de convênio de estágio');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '6c4244c8-32b3-521b-801a-e5becd801c9f', 'e0142aee-e004-51d7-92e5-3cd9c497aba8', 'Plano de estágio', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e0142aee-e004-51d7-92e5-3cd9c497aba8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e0142aee-e004-51d7-92e5-3cd9c497aba8' AND y.nome_documento = 'Plano de estágio');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '0fd35d73-a72d-50ec-93a7-359014668bcd', 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c', 2, 'PGM', 'Procuradoria', 'Emitir parecer jurídico.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '5ca0c673-962f-535e-b40a-951773ddb211', '0fd35d73-a72d-50ec-93a7-359014668bcd', 'Parecer jurídico', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '0fd35d73-a72d-50ec-93a7-359014668bcd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '0fd35d73-a72d-50ec-93a7-359014668bcd' AND y.nome_documento = 'Parecer jurídico');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'ca63ba84-cb38-5fb9-bd94-be92e6aec586', 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c', 3, 'SMS/EDUCACAO', 'Educação permanente', 'Firmar o convênio e distribuir os estagiários.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'c44536bc-ec3b-52a8-9294-ed5bf4457f9c' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '12da902b-3043-592c-8b34-6f30e23a1a67', 'ca63ba84-cb38-5fb9-bd94-be92e6aec586', 'Convênio de estágio', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'ca63ba84-cb38-5fb9-bd94-be92e6aec586')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'ca63ba84-cb38-5fb9-bd94-be92e6aec586' AND y.nome_documento = 'Convênio de estágio');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '6906edde-912f-5a2d-86a5-52c9be214c89', 'e0142aee-e004-51d7-92e5-3cd9c497aba8', '0fd35d73-a72d-50ec-93a7-359014668bcd', 'Vagas confirmadas', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e0142aee-e004-51d7-92e5-3cd9c497aba8')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '0fd35d73-a72d-50ec-93a7-359014668bcd')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e0142aee-e004-51d7-92e5-3cd9c497aba8' AND z.destino_etapa_id = '0fd35d73-a72d-50ec-93a7-359014668bcd' AND z.condicao_transicao = 'Vagas confirmadas');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'bdc295c4-5af3-5dfd-9dcd-4c40894f858b', '0fd35d73-a72d-50ec-93a7-359014668bcd', 'ca63ba84-cb38-5fb9-bd94-be92e6aec586', 'Parecer favorável', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '0fd35d73-a72d-50ec-93a7-359014668bcd')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'ca63ba84-cb38-5fb9-bd94-be92e6aec586')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '0fd35d73-a72d-50ec-93a7-359014668bcd' AND z.destino_etapa_id = 'ca63ba84-cb38-5fb9-bd94-be92e6aec586' AND z.condicao_transicao = 'Parecer favorável');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '73109177-db6b-55d8-a1c4-4f6d3a3df200', 'SAU.AUD.005', 'Auditoria de serviço de saúde', 'Auditar serviço próprio ou contratado do SUS, com relatório e providências.', 'Prestadores do SUS', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '12.0.01.05.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.01.05.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.AUD.005');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'fce2cc4c-af01-545c-920b-85453518aebf', '73109177-db6b-55d8-a1c4-4f6d3a3df200', 1, 'SMS/AUDITORIA', 'Auditoria (DIACAUD)', 'Planejar e comunicar a auditoria.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '73109177-db6b-55d8-a1c4-4f6d3a3df200')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '73109177-db6b-55d8-a1c4-4f6d3a3df200' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '344d9698-6b21-5c5f-b6c9-b9eeeb1ec3c9', 'fce2cc4c-af01-545c-920b-85453518aebf', 'Plano de auditoria', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'fce2cc4c-af01-545c-920b-85453518aebf')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'fce2cc4c-af01-545c-920b-85453518aebf' AND y.nome_documento = 'Plano de auditoria');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b', '73109177-db6b-55d8-a1c4-4f6d3a3df200', 2, 'SMS/AUDITORIA', 'Equipe de auditoria', 'Executar e emitir o relatório preliminar.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '73109177-db6b-55d8-a1c4-4f6d3a3df200')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '73109177-db6b-55d8-a1c4-4f6d3a3df200' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '32622611-6446-590e-b954-a4740264b99a', 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b', 'Relatório preliminar', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b' AND y.nome_documento = 'Relatório preliminar');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '7b10374e-8d85-5259-b057-455c774920f7', '73109177-db6b-55d8-a1c4-4f6d3a3df200', 3, 'SMS/AUDITORIA', 'Auditoria (DIACAUD)', 'Receber a manifestação do auditado e emitir o relatório final.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '73109177-db6b-55d8-a1c4-4f6d3a3df200')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '73109177-db6b-55d8-a1c4-4f6d3a3df200' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd6eea294-e34b-5591-aa1a-3f699ef22eb7', '7b10374e-8d85-5259-b057-455c774920f7', 'Relatório final', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '7b10374e-8d85-5259-b057-455c774920f7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '7b10374e-8d85-5259-b057-455c774920f7' AND y.nome_documento = 'Relatório final');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'e13c7c40-6b15-56b3-ab47-5a597997bd23', 'fce2cc4c-af01-545c-920b-85453518aebf', 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b', 'Auditoria iniciada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'fce2cc4c-af01-545c-920b-85453518aebf')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'fce2cc4c-af01-545c-920b-85453518aebf' AND z.destino_etapa_id = 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b' AND z.condicao_transicao = 'Auditoria iniciada');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '8a6421dd-778f-5d93-8e5a-8bbe704595d2', 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b', '7b10374e-8d85-5259-b057-455c774920f7', 'Relatório emitido', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '7b10374e-8d85-5259-b057-455c774920f7')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'd820b2a9-ff64-522c-9a6a-966ab0e83e7b' AND z.destino_etapa_id = '7b10374e-8d85-5259-b057-455c774920f7' AND z.condicao_transicao = 'Relatório emitido');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '178e49b4-4141-5337-b7ef-8a48879826ef', 'SAU.PCF.006', 'Prestação de contas do Fundo Municipal de Saúde', 'Elaborar e apresentar a prestação de contas quadrimestral do Fundo Municipal de Saúde ao Conselho e ao Legislativo.', 'Conselho Municipal de Saúde e Câmara', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.01.08.03'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.01.08.03' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.PCF.006');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '19c1b101-efd5-50aa-8873-4168247a0ed4', '178e49b4-4141-5337-b7ef-8a48879826ef', 1, 'SMS/FUNDO', 'Fundo Municipal de Saúde', 'Consolidar a execução orçamentária e física do quadrimestre.', 20, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '178e49b4-4141-5337-b7ef-8a48879826ef')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '178e49b4-4141-5337-b7ef-8a48879826ef' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '3f6df0b8-cc7c-5301-81db-a72199f5a4a8', '19c1b101-efd5-50aa-8873-4168247a0ed4', 'Relatório Detalhado do Quadrimestre Anterior (RDQA)', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '19c1b101-efd5-50aa-8873-4168247a0ed4')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '19c1b101-efd5-50aa-8873-4168247a0ed4' AND y.nome_documento = 'Relatório Detalhado do Quadrimestre Anterior (RDQA)');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'babaf79a-93a8-56e8-bcd6-7ddccbabacba', '178e49b4-4141-5337-b7ef-8a48879826ef', 2, 'SMS/CMS', 'Conselho Municipal de Saúde', 'Apreciar e deliberar sobre o relatório.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '178e49b4-4141-5337-b7ef-8a48879826ef')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '178e49b4-4141-5337-b7ef-8a48879826ef' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '2913543b-e9e7-5c1b-80fd-672a6e1c8f59', 'babaf79a-93a8-56e8-bcd6-7ddccbabacba', 'Resolução do Conselho de Saúde', TRUE, 'NATO_DIGITAL', 'CONJUNTA_MULTINIVEL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'babaf79a-93a8-56e8-bcd6-7ddccbabacba')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'babaf79a-93a8-56e8-bcd6-7ddccbabacba' AND y.nome_documento = 'Resolução do Conselho de Saúde');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e07960ce-1ef3-5d34-abe3-94a4dc2717da', '178e49b4-4141-5337-b7ef-8a48879826ef', 3, 'SMS/FUNDO', 'Fundo Municipal de Saúde', 'Apresentar em audiência pública na Câmara.', 15, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '178e49b4-4141-5337-b7ef-8a48879826ef')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '178e49b4-4141-5337-b7ef-8a48879826ef' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '953a9431-12eb-53bf-abaa-5298df300e4c', 'e07960ce-1ef3-5d34-abe3-94a4dc2717da', 'Ata da audiência pública', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e07960ce-1ef3-5d34-abe3-94a4dc2717da')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e07960ce-1ef3-5d34-abe3-94a4dc2717da' AND y.nome_documento = 'Ata da audiência pública');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '617cd18d-d313-51e5-9f45-41e58aa4fb35', '19c1b101-efd5-50aa-8873-4168247a0ed4', 'babaf79a-93a8-56e8-bcd6-7ddccbabacba', 'Relatório consolidado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '19c1b101-efd5-50aa-8873-4168247a0ed4')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'babaf79a-93a8-56e8-bcd6-7ddccbabacba')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '19c1b101-efd5-50aa-8873-4168247a0ed4' AND z.destino_etapa_id = 'babaf79a-93a8-56e8-bcd6-7ddccbabacba' AND z.condicao_transicao = 'Relatório consolidado');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '9a406a08-19c2-512a-b717-c94f6b723fff', 'babaf79a-93a8-56e8-bcd6-7ddccbabacba', 'e07960ce-1ef3-5d34-abe3-94a4dc2717da', 'Relatório apreciado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'babaf79a-93a8-56e8-bcd6-7ddccbabacba')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'e07960ce-1ef3-5d34-abe3-94a4dc2717da')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'babaf79a-93a8-56e8-bcd6-7ddccbabacba' AND z.destino_etapa_id = 'e07960ce-1ef3-5d34-abe3-94a4dc2717da' AND z.condicao_transicao = 'Relatório apreciado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '1225dbd1-bddf-54d6-a46d-4c2130d094b3', 'SAU.DIU.007', 'Inserção de DIU', 'Atender a solicitação de método contraceptivo de longa duração (DIU) na rede municipal.', 'Usuárias do SUS', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '12.0.02.00.81'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.02.00.81' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.DIU.007');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'cdaa4ab8-5493-5412-a042-690c4e048440', '1225dbd1-bddf-54d6-a46d-4c2130d094b3', 1, 'SMS/UBS', 'Unidade Básica de Saúde', 'Acolher a solicitação, orientar e encaminhar com os exames.', 7, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1225dbd1-bddf-54d6-a46d-4c2130d094b3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1225dbd1-bddf-54d6-a46d-4c2130d094b3' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '138fca6a-46a9-5754-a898-b9e8c3ead58d', 'cdaa4ab8-5493-5412-a042-690c4e048440', 'Solicitação de DIU', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'cdaa4ab8-5493-5412-a042-690c4e048440')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'cdaa4ab8-5493-5412-a042-690c4e048440' AND y.nome_documento = 'Solicitação de DIU');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '6e27e596-4f5a-546d-8bf5-59cef74935f8', '1225dbd1-bddf-54d6-a46d-4c2130d094b3', 2, 'SMS/SAUDE_MULHER', 'Saúde da mulher', 'Agendar e realizar o procedimento.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '1225dbd1-bddf-54d6-a46d-4c2130d094b3')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '1225dbd1-bddf-54d6-a46d-4c2130d094b3' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4c65e658-16a1-5872-878d-bb9d8441eac3', '6e27e596-4f5a-546d-8bf5-59cef74935f8', 'Termo de consentimento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6e27e596-4f5a-546d-8bf5-59cef74935f8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6e27e596-4f5a-546d-8bf5-59cef74935f8' AND y.nome_documento = 'Termo de consentimento');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'eb1f14da-7513-5071-a148-0a74da8040e7', '6e27e596-4f5a-546d-8bf5-59cef74935f8', 'Registro do procedimento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '6e27e596-4f5a-546d-8bf5-59cef74935f8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '6e27e596-4f5a-546d-8bf5-59cef74935f8' AND y.nome_documento = 'Registro do procedimento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '4c4ec655-4b70-5076-97c0-69284340b2ac', 'cdaa4ab8-5493-5412-a042-690c4e048440', '6e27e596-4f5a-546d-8bf5-59cef74935f8', 'Encaminhada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'cdaa4ab8-5493-5412-a042-690c4e048440')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '6e27e596-4f5a-546d-8bf5-59cef74935f8')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'cdaa4ab8-5493-5412-a042-690c4e048440' AND z.destino_etapa_id = '6e27e596-4f5a-546d-8bf5-59cef74935f8' AND z.condicao_transicao = 'Encaminhada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '30bf950f-9034-5c40-83f5-c16c583513e1', 'SAU.CNS.008', 'Cartão Nacional de Saúde', 'Cadastrar ou atualizar o usuário no Cartão Nacional de Saúde (CNS).', 'Usuários do SUS', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '12.0.04.00.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.04.00.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.CNS.008');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '101efd7d-bb2c-5a80-be9f-2c5590472faa', '30bf950f-9034-5c40-83f5-c16c583513e1', 1, 'SMS/UBS', 'Recepção da unidade', 'Conferir os documentos e cadastrar ou atualizar no CADSUS.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '30bf950f-9034-5c40-83f5-c16c583513e1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '30bf950f-9034-5c40-83f5-c16c583513e1' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'c0583b53-bf04-5870-8479-3a5341566174', '101efd7d-bb2c-5a80-be9f-2c5590472faa', 'Documento de identificação', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '101efd7d-bb2c-5a80-be9f-2c5590472faa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '101efd7d-bb2c-5a80-be9f-2c5590472faa' AND y.nome_documento = 'Documento de identificação');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '16fe3fda-c202-5154-a883-4a3c29dd9c73', '101efd7d-bb2c-5a80-be9f-2c5590472faa', 'Comprovante de endereço', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '101efd7d-bb2c-5a80-be9f-2c5590472faa')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '101efd7d-bb2c-5a80-be9f-2c5590472faa' AND y.nome_documento = 'Comprovante de endereço');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '79cba7a9-23f4-5424-9b5d-0ece33e29791', '30bf950f-9034-5c40-83f5-c16c583513e1', 2, 'SMS/UBS', 'Recepção da unidade', 'Entregar o número do CNS ao usuário.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '30bf950f-9034-5c40-83f5-c16c583513e1')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '30bf950f-9034-5c40-83f5-c16c583513e1' AND x.ordem = 2);

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '05450924-bfb9-53b5-a0aa-1b06db187898', '101efd7d-bb2c-5a80-be9f-2c5590472faa', '79cba7a9-23f4-5424-9b5d-0ece33e29791', 'Cadastro feito', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '101efd7d-bb2c-5a80-be9f-2c5590472faa')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '79cba7a9-23f4-5424-9b5d-0ece33e29791')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '101efd7d-bb2c-5a80-be9f-2c5590472faa' AND z.destino_etapa_id = '79cba7a9-23f4-5424-9b5d-0ece33e29791' AND z.condicao_transicao = 'Cadastro feito');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'a261b435-75fe-58b6-9719-0867f3418298', 'SAU.PRO.009', 'Cópia de prontuário', 'Fornecer cópia do prontuário ao paciente ou ao representante legal, preservando o sigilo.', 'Pacientes e representantes legais', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '12.0.07.02.04'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.07.02.04' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.PRO.009');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b', 'a261b435-75fe-58b6-9719-0867f3418298', 1, 'SMS/ARQUIVO', 'Arquivo médico', 'Receber o pedido e conferir a legitimidade do solicitante.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a261b435-75fe-58b6-9719-0867f3418298')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a261b435-75fe-58b6-9719-0867f3418298' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '35d24c7b-f28f-5134-ba60-b2d72e68af3e', 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b', 'Requerimento de cópia de prontuário', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b' AND y.nome_documento = 'Requerimento de cópia de prontuário');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '4204bc84-547f-5542-8162-daf5ed8a24a1', 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b', 'Documento de identificação ou procuração', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b' AND y.nome_documento = 'Documento de identificação ou procuração');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '3984ef1d-92a2-5324-b91c-26d11a144bd6', 'a261b435-75fe-58b6-9719-0867f3418298', 2, 'SMS/ARQUIVO', 'Arquivo médico', 'Reproduzir o prontuário e entregar mediante recibo.', 10, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'a261b435-75fe-58b6-9719-0867f3418298')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'a261b435-75fe-58b6-9719-0867f3418298' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'cbb5df2c-2dba-54f8-9b08-3f40b8e460b9', '3984ef1d-92a2-5324-b91c-26d11a144bd6', 'Recibo de entrega', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '3984ef1d-92a2-5324-b91c-26d11a144bd6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '3984ef1d-92a2-5324-b91c-26d11a144bd6' AND y.nome_documento = 'Recibo de entrega');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '0764be63-f166-5721-82b8-45d208b1202e', 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b', '3984ef1d-92a2-5324-b91c-26d11a144bd6', 'Solicitante legítimo', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '3984ef1d-92a2-5324-b91c-26d11a144bd6')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'e5beff5e-c425-5c1a-a0b2-45f33dd2618b' AND z.destino_etapa_id = '3984ef1d-92a2-5324-b91c-26d11a144bd6' AND z.condicao_transicao = 'Solicitante legítimo');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '5559f8ee-3df2-5402-bf61-06db2a57be45', 'SAU.OPM.010', 'Órteses, próteses e meios auxiliares (OPM)', 'Cadastrar e atender o paciente no programa de órteses, próteses e materiais especiais (auditivo, ortopédico).', 'Usuários do SUS', 1, TRUE, 'HOMOLOGADO', 'RESTRITO', NULLIF('LGPD, art. 11 (dado de saúde) e LAI, art. 31', ''), '12.0.07.04.16'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.07.04.16' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.OPM.010');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4775e974-2f44-5f85-872b-b09d5a44713d', '5559f8ee-3df2-5402-bf61-06db2a57be45', 1, 'SMS/OPM', 'Programa de OPM', 'Cadastrar o paciente com a prescrição médica.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5559f8ee-3df2-5402-bf61-06db2a57be45')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5559f8ee-3df2-5402-bf61-06db2a57be45' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'd9157416-3fb6-508c-b6f7-83d3924a27ab', '4775e974-2f44-5f85-872b-b09d5a44713d', 'Prescrição médica', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4775e974-2f44-5f85-872b-b09d5a44713d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4775e974-2f44-5f85-872b-b09d5a44713d' AND y.nome_documento = 'Prescrição médica');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '111b26b2-561a-5c65-81b7-037e62f1c012', '4775e974-2f44-5f85-872b-b09d5a44713d', 'Ficha de cadastro no programa', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4775e974-2f44-5f85-872b-b09d5a44713d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4775e974-2f44-5f85-872b-b09d5a44713d' AND y.nome_documento = 'Ficha de cadastro no programa');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '1b447616-8ea7-516f-9b31-675916840e8d', '5559f8ee-3df2-5402-bf61-06db2a57be45', 2, 'SMS/REGULACAO', 'Complexo regulador', 'Autorizar e agendar a concessão.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5559f8ee-3df2-5402-bf61-06db2a57be45')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5559f8ee-3df2-5402-bf61-06db2a57be45' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bbb2fb09-5eb7-5a26-82a8-3f2623e3d658', '1b447616-8ea7-516f-9b31-675916840e8d', 'Autorização de procedimento', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '1b447616-8ea7-516f-9b31-675916840e8d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '1b447616-8ea7-516f-9b31-675916840e8d' AND y.nome_documento = 'Autorização de procedimento');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '92a0b711-5273-5735-b8d5-c4043663e229', '5559f8ee-3df2-5402-bf61-06db2a57be45', 3, 'SMS/OPM', 'Programa de OPM', 'Entregar o equipamento e orientar o uso.', 30, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '5559f8ee-3df2-5402-bf61-06db2a57be45')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '5559f8ee-3df2-5402-bf61-06db2a57be45' AND x.ordem = 3);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '012475e6-ed60-5bf8-ab49-2b186b59b0d2', '92a0b711-5273-5735-b8d5-c4043663e229', 'Termo de recebimento', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '92a0b711-5273-5735-b8d5-c4043663e229')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '92a0b711-5273-5735-b8d5-c4043663e229' AND y.nome_documento = 'Termo de recebimento');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'd096b72a-3a81-548f-a3c8-5c4e6d88044f', '4775e974-2f44-5f85-872b-b09d5a44713d', '1b447616-8ea7-516f-9b31-675916840e8d', 'Cadastro feito', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4775e974-2f44-5f85-872b-b09d5a44713d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '1b447616-8ea7-516f-9b31-675916840e8d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4775e974-2f44-5f85-872b-b09d5a44713d' AND z.destino_etapa_id = '1b447616-8ea7-516f-9b31-675916840e8d' AND z.condicao_transicao = 'Cadastro feito');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'd60a1949-f7a5-576e-86af-97f2cbc22341', '1b447616-8ea7-516f-9b31-675916840e8d', '92a0b711-5273-5735-b8d5-c4043663e229', 'Concessão autorizada', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '1b447616-8ea7-516f-9b31-675916840e8d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '92a0b711-5273-5735-b8d5-c4043663e229')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '1b447616-8ea7-516f-9b31-675916840e8d' AND z.destino_etapa_id = '92a0b711-5273-5735-b8d5-c4043663e229' AND z.condicao_transicao = 'Concessão autorizada');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT '945a34f7-3fae-574b-90ec-000120980a9b', 'SAU.PLA.011', 'Inclusão ou exclusão de plantão médico', 'Ajustar a escala de plantões médicos conforme a demanda das unidades.', 'Coordenações de unidades de saúde', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.08.05.08'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.08.05.08' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.PLA.011');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'dde87b7d-5697-5019-b61a-375fbb68248a', '945a34f7-3fae-574b-90ec-000120980a9b', 1, 'SMS/UNIDADE', 'Coordenação da unidade', 'Solicitar a inclusão ou exclusão com a justificativa de demanda.', 3, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '945a34f7-3fae-574b-90ec-000120980a9b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '945a34f7-3fae-574b-90ec-000120980a9b' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e07abbbb-0953-56f2-8d51-838ee35f0474', 'dde87b7d-5697-5019-b61a-375fbb68248a', 'Solicitação de alteração de plantão', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'dde87b7d-5697-5019-b61a-375fbb68248a')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'dde87b7d-5697-5019-b61a-375fbb68248a' AND y.nome_documento = 'Solicitação de alteração de plantão');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '960550a7-c047-5546-89d2-cfdb543ad54f', '945a34f7-3fae-574b-90ec-000120980a9b', 2, 'SMS/RH_MEDICOS', 'RH dos médicos', 'Ajustar a escala e informar a folha.', 5, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = '945a34f7-3fae-574b-90ec-000120980a9b')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = '945a34f7-3fae-574b-90ec-000120980a9b' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'bccfc701-e656-571c-94a3-028e507f3a70', '960550a7-c047-5546-89d2-cfdb543ad54f', 'Escala de plantões', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '960550a7-c047-5546-89d2-cfdb543ad54f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '960550a7-c047-5546-89d2-cfdb543ad54f' AND y.nome_documento = 'Escala de plantões');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT 'bca69ebf-f17f-5652-a98b-7e64bcf79896', 'dde87b7d-5697-5019-b61a-375fbb68248a', '960550a7-c047-5546-89d2-cfdb543ad54f', 'Pedido justificado', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = 'dde87b7d-5697-5019-b61a-375fbb68248a')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = '960550a7-c047-5546-89d2-cfdb543ad54f')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = 'dde87b7d-5697-5019-b61a-375fbb68248a' AND z.destino_etapa_id = '960550a7-c047-5546-89d2-cfdb543ad54f' AND z.condicao_transicao = 'Pedido justificado');


INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT 'b5b4ff85-5fd3-58c9-8492-0516a3d48308', 'SAU.AGU.012', 'Análise microbiológica de água', 'Receber amostra e emitir laudo de análise microbiológica de água para consumo.', 'Cidadãos, empresas e vigilância', 1, TRUE, 'HOMOLOGADO', 'PUBLICO', NULLIF('', ''), '12.0.11.03.01'
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = '12.0.11.03.01' AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = 'SAU.AGU.012');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT '4d02f832-6c92-5585-8a4f-ac559bdc047d', 'b5b4ff85-5fd3-58c9-8492-0516a3d48308', 1, 'SMS/LACEN_AGUA', 'Laboratório de análise de água', 'Receber a amostra com a ficha de coleta.', 1, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b5b4ff85-5fd3-58c9-8492-0516a3d48308')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b5b4ff85-5fd3-58c9-8492-0516a3d48308' AND x.ordem = 1);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'e9d56517-6d69-5c8e-b600-bd3e8fe09a4f', '4d02f832-6c92-5585-8a4f-ac559bdc047d', 'Solicitação de análise', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4d02f832-6c92-5585-8a4f-ac559bdc047d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4d02f832-6c92-5585-8a4f-ac559bdc047d' AND y.nome_documento = 'Solicitação de análise');

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT 'a794ebda-23eb-547a-ac8c-474efe94af37', '4d02f832-6c92-5585-8a4f-ac559bdc047d', 'Ficha de coleta', TRUE, 'EXTERNO_DIGITALIZADO', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = '4d02f832-6c92-5585-8a4f-ac559bdc047d')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = '4d02f832-6c92-5585-8a4f-ac559bdc047d' AND y.nome_documento = 'Ficha de coleta');

INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT 'd28b5726-a735-5f9a-93a9-740007dcaa23', 'b5b4ff85-5fd3-58c9-8492-0516a3d48308', 2, 'SMS/LACEN_AGUA', 'Laboratório de análise de água', 'Analisar e emitir o laudo.', 7, FALSE
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = 'b5b4ff85-5fd3-58c9-8492-0516a3d48308')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = 'b5b4ff85-5fd3-58c9-8492-0516a3d48308' AND x.ordem = 2);

INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT '03c05e3e-5051-58be-8dc8-a9eab38be55a', 'd28b5726-a735-5f9a-93a9-740007dcaa23', 'Laudo de análise microbiológica', TRUE, 'NATO_DIGITAL', 'INDIVIDUAL', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = 'd28b5726-a735-5f9a-93a9-740007dcaa23')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = 'd28b5726-a735-5f9a-93a9-740007dcaa23' AND y.nome_documento = 'Laudo de análise microbiológica');

INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT '23254c5d-ce97-5f27-956e-2a41129fc231', '4d02f832-6c92-5585-8a4f-ac559bdc047d', 'd28b5726-a735-5f9a-93a9-740007dcaa23', 'Amostra aceita', FALSE, NULLIF('', '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = '4d02f832-6c92-5585-8a4f-ac559bdc047d')
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = 'd28b5726-a735-5f9a-93a9-740007dcaa23')
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = '4d02f832-6c92-5585-8a4f-ac559bdc047d' AND z.destino_etapa_id = 'd28b5726-a735-5f9a-93a9-740007dcaa23' AND z.condicao_transicao = 'Amostra aceita');


UPDATE atlas_workflows
   SET situacao = 'HOMOLOGADO', ativo = TRUE
 WHERE codigo_processual IN ('ADM.CMP.004', 'ADM.CON.005', 'ADM.CON.006', 'ADM.CMP.007', 'ADM.FRO.008', 'ADM.SIC.009', 'RH.FER.001', 'RH.LIC.002', 'RH.ADI.003', 'RH.APO.004', 'RH.PAD.005', 'RH.EST.006', 'RH.CON.007', 'FAZ.DIA.001', 'FAZ.DES.002', 'FAZ.PAR.003', 'FAZ.IPT.004', 'FAZ.ALV.005', 'FAZ.CND.006', 'CTR.OUV.001', 'HAB.ALV.001', 'HAB.HAB.002', 'HAB.REU.003', 'MAM.LIC.001', 'DEC.USO.001', 'SAS.OSC.001', 'SAU.VSA.001', 'SAU.TFD.001', 'ADM.CEM.010', 'ADM.PRO.011', 'ADM.DOC.012', 'TI.ACE.001', 'TI.SUP.002', 'RH.CAP.008', 'RH.FOL.009', 'RH.RES.010', 'RH.VTR.011', 'RH.CNS.012', 'RH.PAL.013', 'RH.REE.014', 'RH.INS.015', 'RH.LMA.016', 'RH.LPA.017', 'RH.LNO.018', 'RH.PER.019', 'RH.CAT.020', 'RH.REA.021', 'RH.RCH.022', 'RH.COM.023', 'RH.TMP.024', 'RH.ETG.025', 'RH.SEL.026', 'FAZ.FIS.007', 'FAZ.ISS.008', 'FAZ.JUL.009', 'FAZ.FPM.010', 'FAZ.CAD.011', 'FAZ.VVE.012', 'FAZ.CUS.013', 'FAZ.ORC.014', 'PAT.BAI.001', 'PAT.TRA.002', 'PAT.INC.003', 'CTR.AUD.002', 'HAB.ITB.004', 'HAB.MCM.005', 'HAB.RRU.006', 'HAB.DES.007', 'HAB.ALI.008', 'MAM.FIS.002', 'MAM.ARV.003', 'MAM.REC.004', 'AGR.SIM.001', 'AGR.TRA.002', 'AGR.POC.003', 'AGR.EST.004', 'DEC.INC.002', 'DEC.IND.003', 'SAS.CNV.002', 'SAS.CHA.003', 'SAS.BIL.004', 'SAU.EMD.002', 'SAU.CNE.003', 'SAU.EST.004', 'SAU.AUD.005', 'SAU.PCF.006', 'SAU.DIU.007', 'SAU.CNS.008', 'SAU.PRO.009', 'SAU.OPM.010', 'SAU.PLA.011', 'SAU.AGU.012')
   AND situacao <> 'HOMOLOGADO';

DO $chk$
DECLARE
  faltam TEXT;
BEGIN
  SELECT string_agg(c.codigo, ', ' ORDER BY c.codigo) INTO faltam
    FROM (VALUES ('ADM.CMP.004'), ('ADM.CON.005'), ('ADM.CON.006'), ('ADM.CMP.007'), ('ADM.FRO.008'), ('ADM.SIC.009'), ('RH.FER.001'), ('RH.LIC.002'), ('RH.ADI.003'), ('RH.APO.004'), ('RH.PAD.005'), ('RH.EST.006'), ('RH.CON.007'), ('FAZ.DIA.001'), ('FAZ.DES.002'), ('FAZ.PAR.003'), ('FAZ.IPT.004'), ('FAZ.ALV.005'), ('FAZ.CND.006'), ('CTR.OUV.001'), ('HAB.ALV.001'), ('HAB.HAB.002'), ('HAB.REU.003'), ('MAM.LIC.001'), ('DEC.USO.001'), ('SAS.OSC.001'), ('SAU.VSA.001'), ('SAU.TFD.001'), ('ADM.CEM.010'), ('ADM.PRO.011'), ('ADM.DOC.012'), ('TI.ACE.001'), ('TI.SUP.002'), ('RH.CAP.008'), ('RH.FOL.009'), ('RH.RES.010'), ('RH.VTR.011'), ('RH.CNS.012'), ('RH.PAL.013'), ('RH.REE.014'), ('RH.INS.015'), ('RH.LMA.016'), ('RH.LPA.017'), ('RH.LNO.018'), ('RH.PER.019'), ('RH.CAT.020'), ('RH.REA.021'), ('RH.RCH.022'), ('RH.COM.023'), ('RH.TMP.024'), ('RH.ETG.025'), ('RH.SEL.026'), ('FAZ.FIS.007'), ('FAZ.ISS.008'), ('FAZ.JUL.009'), ('FAZ.FPM.010'), ('FAZ.CAD.011'), ('FAZ.VVE.012'), ('FAZ.CUS.013'), ('FAZ.ORC.014'), ('PAT.BAI.001'), ('PAT.TRA.002'), ('PAT.INC.003'), ('CTR.AUD.002'), ('HAB.ITB.004'), ('HAB.MCM.005'), ('HAB.RRU.006'), ('HAB.DES.007'), ('HAB.ALI.008'), ('MAM.FIS.002'), ('MAM.ARV.003'), ('MAM.REC.004'), ('AGR.SIM.001'), ('AGR.TRA.002'), ('AGR.POC.003'), ('AGR.EST.004'), ('DEC.INC.002'), ('DEC.IND.003'), ('SAS.CNV.002'), ('SAS.CHA.003'), ('SAS.BIL.004'), ('SAU.EMD.002'), ('SAU.CNE.003'), ('SAU.EST.004'), ('SAU.AUD.005'), ('SAU.PCF.006'), ('SAU.DIU.007'), ('SAU.CNS.008'), ('SAU.PRO.009'), ('SAU.OPM.010'), ('SAU.PLA.011'), ('SAU.AGU.012')) AS c(codigo)
   WHERE NOT EXISTS (
     SELECT 1 FROM atlas_workflows w
      WHERE w.codigo_processual = c.codigo AND w.ativo AND w.situacao = 'HOMOLOGADO'
   );
  IF faltam IS NOT NULL THEN
    RAISE EXCEPTION 'procedimentos não publicados (série TTDD ausente ou inserção falhou): %', faltam;
  END IF;
END
$chk$;
COMMIT;

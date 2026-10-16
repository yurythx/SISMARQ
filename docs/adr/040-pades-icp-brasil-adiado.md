# 040 — PAdES ICP-Brasil adiado até a política municipal (decisão 0.6)

- **Status:** Aceito (adiado de propósito; não implementar um certificado de mentira)
- **Data:** 2026-10-16
- **Relacionado:** [ADR 037](037-acervo-digitalizacao.md) (pipeline PDF/A), Signum (assinatura eletrônica com reautenticação)

## Contexto

O Decreto 10.278/2020 e o artigo científico pedem assinatura ICP-Brasil
(PAdES) no representante digital. O Signum hoje é assinatura **eletrônica**
com reautenticação e SHA-256, sem certificado A1/A3. A decisão 0.6 do
roadmap (política de digitalização, quem assina, qual certificado) ainda
não existe.

Implementar um “PAdES” com chave local, autoassinada ou com o Signum
atual **fingiria** valor de prova que o município não tem. Isso é pior do
que declarar a ausência.

## Decisão

1. **Não assinar com ICP-Brasil enquanto a política 0.6 não existir.**
   O pipeline da Fase 5 permanece PDF/A-2b + OCR + XMP + SHA-256 + object
   lock + fixity. Isso já dá integridade e metadados do Decreto; não dá
   autoria ICP-Brasil.
2. **O JSON declara a ausência.** O campo `pades` do representante vale
   `AUSENTE`. A tela autenticada mostra que a política ainda não aplica.
   Nada no disco muda. `AfirmaAssinaturaICP` permanece falso.
3. **O Signum não vira PAdES.** Continua nos envelopes de processo e nas
   atas quando a Fase 3 ligar a cerimônia. Certificado ICP-Brasil é outro
   produto (HSM ou A1 do município, política de quem assina o lote).
4. **Quando a 0.6 existir**, um ADR posterior liga um `AssinadorICP` na
   composição (`internal/app`), assina o PDF/A **depois** do pipeline
   (o hash muda; o banco e o XMP acompanham) e grava o dicionário PAdES.
   Até lá não há porta ligada, de propósito.

## Consequências

- O valor probatório **completo** do Decreto 10.278 fica pendente do
  certificado. O dossiê e o VPAT não afirmam PAdES.
- Pasta monitorada SMB/FTP também espera a 0.6 (scanners de produção).
- Testes travam qualquer regressão que passe a afirmar ICP sem política.

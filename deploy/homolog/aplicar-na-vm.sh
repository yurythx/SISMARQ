#!/usr/bin/env bash
# Homologação na VM: branding, organograma, TTDD, caixas SEMAD, digitais e
# conferência automática. Não faz git pull. Migrations: o compose `migrate`
# (scripts/deploy.sh --no-pull) precisa ter rodado antes.
set -euo pipefail
cd /root/Projetos/SISMARQ
DB_USER=$(grep -E '^DB_USER=' .env | cut -d= -f2)
DB_NAME=$(grep -E '^DB_NAME=' .env | cut -d= -f2)
DB_USER=${DB_USER:-nexus}
DB_NAME=${DB_NAME:-nexus}

psql() { docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" "$@"; }

if ! psql -tAc "SELECT 1 FROM information_schema.columns WHERE table_name='acervo_digitais' AND column_name='origem'" | grep -q 1; then
  echo "Falta a migration 000156 (coluna acervo_digitais.origem). Rode: bash scripts/deploy.sh --no-pull" >&2
  exit 1
fi
goose_ver=$(psql -tAc "SELECT COALESCE(MAX(version_id),0) FROM goose_db_version WHERE is_applied" | tr -d '[:space:]')
if [ "${goose_ver:-0}" -lt 157 ]; then
  echo "Falta a migration 000157 (catálogo/blog PaperMoon). goose=$goose_ver. Rode: bash scripts/deploy.sh --no-pull" >&2
  exit 1
fi

psql -c "ALTER TABLE branding_settings ALTER COLUMN org_name SET DEFAULT 'Superintendência de Arquivo Público';"
psql -c "UPDATE branding_settings SET org_name = 'Superintendência de Arquivo Público' WHERE org_name IN ('Organização', '');"

if [ -f deploy/estrutura/estrutura.sql ]; then
  echo "== estrutura"
  psql -q < deploy/estrutura/estrutura.sql
fi
if [ -f deploy/ttdd/ttdd.sql ] && [ -f deploy/ttdd/carga.sql ]; then
  echo "== ttdd"
  cat deploy/ttdd/ttdd.sql deploy/ttdd/carga.sql | docker compose exec -T postgres psql -q -v ON_ERROR_STOP=1 -v aplicar=1 -U "$DB_USER" -d "$DB_NAME"
fi
echo "== seed acervo"
psql < deploy/homolog/seed-acervo.sql
echo "== seed digitais"
psql < deploy/homolog/seed-digitais.sql
echo "== conferência automática (C1/C5/C6 metadados)"
psql < deploy/homolog/conferir.sql
echo "== branding"
psql -c "SELECT app_name, org_name FROM branding_settings;"

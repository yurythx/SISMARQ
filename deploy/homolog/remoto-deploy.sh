#!/usr/bin/env bash
# Roda na VM: publica o que está no main do GitHub — todo deploy é um commit
# conhecido. Recusa se o servidor tiver alteração local (não versionada).
#
#   bash deploy/homolog/remoto-deploy.sh [host]
set -euo pipefail
APP_DIR="${APP_DIR:-/opt/sismarq}"
HOST="${1:-192.168.1.52}"
cd "$APP_DIR"
if [ -n "$(git status --porcelain)" ]; then
  echo "Há alterações locais em $APP_DIR — o servidor segue o git. Abortando:" >&2
  git status --short >&2
  exit 1
fi
git fetch -q origin main
git merge --ff-only -q origin/main
echo "==> publicando $(git log --oneline -1)"
export COMPOSE_PARALLEL_LIMIT=1
bash scripts/deploy.sh --no-pull "$HOST" </dev/null
bash deploy/homolog/aplicar-na-vm.sh </dev/null

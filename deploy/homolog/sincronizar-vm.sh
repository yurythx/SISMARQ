#!/usr/bin/env bash
# Publica na VM o que está no main do GitHub (deploy rastreável pelo git).
# Antes confere se a sua cópia não tem nada fora do GitHub: o que não foi
# commitado e enviado não vai para o servidor.
#
#   HOST=192.168.1.52 bash deploy/homolog/sincronizar-vm.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
HOST="${HOST:-192.168.1.52}"
REMOTE="${REMOTE:-/opt/sismarq}"
USER="${SSHUSER:-root}"
KEY="${SSH_KEY:-$HOME/.ssh/sismarq_deploy}"
SSH=(ssh -o BatchMode=yes)
[ -f "$KEY" ] && SSH+=(-i "$KEY")

cd "$ROOT"
git fetch -q origin main
if [ -n "$(git status --porcelain)" ]; then
  echo "Há alterações não commitadas — faça commit e push antes de publicar:" >&2
  git status --short >&2
  exit 1
fi
if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/main)" ]; then
  echo "Sua cópia não está igual ao origin/main — envie (push) ou atualize antes de publicar." >&2
  exit 1
fi

"${SSH[@]}" "$USER@$HOST" "cd '$REMOTE' && git fetch -q origin main && git merge --ff-only -q origin/main && bash deploy/homolog/remoto-deploy.sh '$HOST' </dev/null"

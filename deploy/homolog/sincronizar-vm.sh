#!/usr/bin/env bash
# Leva o working tree local à VM de homologação e sobe a stack **sem**
# git pull (o GitHub main ainda pode estar com a marca Nexus).
#
#   HOST=192.168.1.52 bash deploy/homolog/sincronizar-vm.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
HOST="${HOST:-192.168.1.52}"
REMOTE="${REMOTE:-/root/Projetos/SISMARQ}"
USER="${SSHUSER:-root}"
KEY="${SSH_KEY:-$HOME/.ssh/sismarq_deploy}"
SSH=(ssh -o BatchMode=yes)
SCP=(scp -o BatchMode=yes)
if [ -f "$KEY" ]; then
  SSH+=(-i "$KEY")
  SCP+=(-i "$KEY")
fi

cd "$ROOT"
tar --exclude-vcs-ignores --exclude=.git \
    --exclude=node_modules --exclude=frontend/node_modules \
    --exclude=.next --exclude=frontend/.next \
    --exclude=.terraform --exclude='*.tfstate' --exclude='*.tfstate.*' \
    --exclude=terraform.tfvars --exclude=secrets.yml \
    --exclude=inventory.ini --exclude=.env \
    --exclude=backend/api --exclude=backend/worker \
    --exclude=backend/coverage --exclude=frontend/coverage --exclude=coverage \
    -cf - . | "${SSH[@]}" "$USER@$HOST" "mkdir -p '$REMOTE' && cd '$REMOTE' && tar xf -"

"${SSH[@]}" "$USER@$HOST" bash -s <<EOF
set -euo pipefail
cd '$REMOTE'
export COMPOSE_PARALLEL_LIMIT=1
bash scripts/deploy.sh --no-pull $HOST
bash deploy/homolog/aplicar-na-vm.sh
EOF

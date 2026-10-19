#!/usr/bin/env bash
# Roda na VM depois do tar. Sem git pull.
set -euo pipefail
cd /root/Projetos/SISMARQ
if [ -f /tmp/sismarq-homolog.tar ]; then
  tar xf /tmp/sismarq-homolog.tar
  rm -f /tmp/sismarq-homolog.tar
fi
find deploy/homolog scripts -name '*.sh' -exec sed -i 's/\r$//' {} +
export COMPOSE_PARALLEL_LIMIT=1
HOST="${1:-192.168.1.52}"
bash scripts/deploy.sh --no-pull "$HOST"
bash deploy/homolog/aplicar-na-vm.sh

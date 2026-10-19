# Publica na VM 220 o que está no main do GitHub (deploy rastreável pelo git).
# Antes confere se a sua cópia não tem nada fora do GitHub.
#   powershell -File deploy/homolog/sincronizar-vm.ps1
$ErrorActionPreference = "Stop"
$HostName = if ($env:HOST) { $env:HOST } else { "192.168.1.52" }
$User = if ($env:SSHUSER) { $env:SSHUSER } else { "root" }
$Remote = if ($env:REMOTE) { $env:REMOTE } else { "/opt/sismarq" }
$Key = if ($env:SSH_KEY) { $env:SSH_KEY } else { Join-Path $env:USERPROFILE ".ssh\sismarq_deploy" }
$SshOpts = @("-o", "BatchMode=yes")
if (Test-Path $Key) { $SshOpts += @("-i", $Key) }
Set-Location (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path

git fetch -q origin main
if (git status --porcelain) {
  git status --short
  throw "Há alterações não commitadas — faça commit e push antes de publicar."
}
if ((git rev-parse HEAD) -ne (git rev-parse origin/main)) {
  throw "Sua cópia não está igual ao origin/main — envie (push) ou atualize antes de publicar."
}
ssh @SshOpts "${User}@${HostName}" "cd '$Remote' && git fetch -q origin main && git merge --ff-only -q origin/main && bash deploy/homolog/remoto-deploy.sh '$HostName' </dev/null"

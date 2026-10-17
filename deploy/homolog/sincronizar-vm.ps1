# Leva o working tree local à VM 220 sem git pull (GitHub main pode ser Nexus).
#   powershell -File deploy/homolog/sincronizar-vm.ps1
$ErrorActionPreference = "Stop"
$HostName = if ($env:HOST) { $env:HOST } else { "192.168.1.52" }
$User = if ($env:SSHUSER) { $env:SSHUSER } else { "root" }
$Key = if ($env:SSH_KEY) { $env:SSH_KEY } else { Join-Path $env:USERPROFILE ".ssh\sismarq_deploy" }
$SshOpts = @("-o", "BatchMode=yes")
if (Test-Path $Key) { $SshOpts += @("-i", $Key) }
$Root = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
Set-Location $Root
$tar = Join-Path ([System.IO.Path]::GetTempPath()) "sismarq-homolog.tar"
if (Test-Path $tar) { Remove-Item -Force $tar }

tar --exclude=.git --exclude=node_modules --exclude=frontend/node_modules `
    --exclude=.next --exclude=frontend/.next --exclude=.terraform `
    --exclude=deploy/infra/terraform/.terraform `
    --exclude=*.tfstate --exclude=*.tfstate.* --exclude=terraform.tfvars `
    --exclude=secrets.yml --exclude=inventory.ini --exclude=.env `
    --exclude=backend/api --exclude=backend/worker --exclude=backend/*.test `
    --exclude=backend/coverage --exclude=frontend/coverage --exclude=coverage `
    -cf $tar .

scp @SshOpts $tar "${User}@${HostName}:/tmp/sismarq-homolog.tar"
scp @SshOpts (Join-Path $PSScriptRoot "remoto-deploy.sh") "${User}@${HostName}:/tmp/remoto-deploy.sh"
ssh @SshOpts "${User}@${HostName}" "sed -i 's/\r$//' /tmp/remoto-deploy.sh && bash /tmp/remoto-deploy.sh $HostName"
if (Test-Path $tar) { Remove-Item -Force $tar }

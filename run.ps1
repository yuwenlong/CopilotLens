param([string]$Action)
$ErrorActionPreference = 'Stop'
$AppName = 'copilotlens'
$BinDir = 'bin'
$PidFile = "$BinDir\.pid"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Start-Server {
    if (Test-Path $PidFile) {
        $serverPid = Get-Content $PidFile -ErrorAction SilentlyContinue
        if ($serverPid -and (Get-Process -Id $serverPid -ErrorAction SilentlyContinue)) {
            Write-Host ('Service already running (PID: ' + $serverPid + ')')
            return
        }
    }
    Write-Host ('Starting ' + $AppName + ' ...')
    $proc = Start-Process -FilePath "$ScriptDir\$BinDir\$AppName.exe" -WorkingDirectory "$ScriptDir\$BinDir" -WindowStyle Hidden -PassThru
    Set-Content -Path $PidFile -Value $proc.Id
    Write-Host ('Service started (PID: ' + $proc.Id + ')')
}

function Stop-Server {
    if (-not (Test-Path $PidFile)) {
        Write-Host 'PID file not found, service not running'
        return
    }
    $serverPid = Get-Content $PidFile -ErrorAction SilentlyContinue
    if ($serverPid) {
        $proc = Get-Process -Id $serverPid -ErrorAction SilentlyContinue
        if ($proc) {
            Write-Host ('Stopping service (PID: ' + $serverPid + ') ...')
            Stop-Process -Id $serverPid -Force
            Write-Host 'Service stopped'
        } else {
            Write-Host ('Process ' + $serverPid + ' not found')
        }
    }
    Remove-Item $PidFile -ErrorAction SilentlyContinue
}

switch ($Action) {
    'start'   { Start-Server }
    'stop'    { Stop-Server }
    'restart' {
        Stop-Server
        Start-Sleep -Seconds 1
        Start-Server
    }
    'reload' {
        Write-Host 'Rebuilding...'
        & "$ScriptDir\build.ps1"
        Stop-Server
        Start-Sleep -Seconds 1
        Start-Server
    }
    default {
        Write-Host 'Usage: .\run.ps1 {start|stop|restart|reload}'
        exit 1
    }
}

param([string]$Action)
$ErrorActionPreference = 'Stop'
$AppName = 'copilotlens'
$BinDir = 'bin'
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Start-Server {
    $existing = Get-Process -Name $AppName -ErrorAction SilentlyContinue
    if ($existing) {
        Write-Host ('Service already running (PID: ' + $existing.Id + ')')
        return
    }
    Write-Host ('Starting ' + $AppName + ' ...')
    Start-Process -FilePath "$ScriptDir\$BinDir\$AppName.exe" -WorkingDirectory "$ScriptDir\$BinDir" -WindowStyle Hidden
    Write-Host 'Service started'
}

function Stop-Server {
    $procs = Get-Process -Name $AppName -ErrorAction SilentlyContinue
    if ($procs) {
        foreach ($proc in $procs) {
            Write-Host ('Stopping service (PID: ' + $proc.Id + ') ...')
            Stop-Process -Id $proc.Id -Force
        }
        Write-Host 'Service stopped'
    } else {
        Write-Host 'Service not running'
    }
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

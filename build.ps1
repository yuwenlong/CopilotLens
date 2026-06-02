go build -o bin\copilotlens.exe .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Remove-Item -Recurse -Force "bin\web" -ErrorAction SilentlyContinue
Copy-Item -Recurse -Force web bin\
Copy-Item -Recurse -Force data bin\
New-Item -ItemType Directory -Force -Path bin\conf | Out-Null
if (!(Test-Path bin\conf\config.toml)) { Copy-Item toml\config_template.toml bin\conf\config.toml }

Write-Host "构建完成: bin\copilotlens.exe"

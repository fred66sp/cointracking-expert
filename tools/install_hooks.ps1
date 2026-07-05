# Script para instalar pre-commit hooks
# Uso: powershell -ExecutionPolicy Bypass .\tools\install_hooks.ps1

Write-Host "======================================================================" -ForegroundColor Cyan
Write-Host "Instalando Pre-Commit Hooks (ADR-037)" -ForegroundColor Cyan
Write-Host "======================================================================" -ForegroundColor Cyan
Write-Host ""

$hookDir = ".git\hooks"

# Verificar si existe directorio hooks
if (-not (Test-Path $hookDir)) {
    Write-Host "ERROR: Directorio .git/hooks no encontrado" -ForegroundColor Red
    exit 1
}

# Copiar hook PowerShell
$hookFile = "$hookDir\pre-commit.ps1"
if (Test-Path $hookFile) {
    Write-Host "[OK] pre-commit.ps1 ya existe" -ForegroundColor Green
} else {
    Write-Host "ERROR: pre-commit.ps1 no encontrado" -ForegroundColor Red
    exit 1
}

# Crear wrapper batch que ejecute PowerShell
$batchFile = "$hookDir\pre-commit"
$batchContent = @"
@echo off
powershell -ExecutionPolicy Bypass -File ".git\hooks\pre-commit.ps1"
exit /b %errorlevel%
"@

# Escribir archivo batch
try {
    Set-Content -Path $batchFile -Value $batchContent -Encoding UTF8 -NoNewline
    Write-Host "[OK] pre-commit (batch wrapper) creado" -ForegroundColor Green
} catch {
    Write-Host "ERROR: No se pudo crear pre-commit wrapper: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Instalación completada:" -ForegroundColor Green
Write-Host "  ✓ Hook PowerShell: $hookFile" -ForegroundColor Green
Write-Host "  ✓ Wrapper batch: $batchFile" -ForegroundColor Green
Write-Host ""
Write-Host "El hook se ejecutará automáticamente en el próximo 'git commit'" -ForegroundColor Cyan
Write-Host ""
Write-Host "Para deshabilitar temporalmente:" -ForegroundColor Yellow
Write-Host "  git commit --no-verify" -ForegroundColor Yellow
Write-Host ""

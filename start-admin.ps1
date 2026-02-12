#!/usr/bin/env pwsh

Write-Host "=== Starting Admin Panel ===" -ForegroundColor Green
Write-Host "URL: http://localhost:5173" -ForegroundColor Cyan
Write-Host "API: http://localhost:8000" -ForegroundColor Cyan
Write-Host ""

Set-Location (Join-Path $PSScriptRoot "ui\admin-panel")
npm run dev

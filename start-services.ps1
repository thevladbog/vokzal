# Build and start all microservices
$repoRoot = $PSScriptRoot
$services = @(
    "auth", "schedule", "ticket", "fiscal", "payment",
    "board", "notify", "audit", "document", "geo"
)
# Single source of truth for ports (must match Traefik/infra)
$portMap = @{
    auth     = 8081
    schedule = 8082
    ticket   = 8083
    fiscal   = 8084
    payment  = 8085
    board    = 8086
    notify   = 8087
    audit    = 8098
    document = 8089
    geo      = 8090
}

Write-Host "Building all services..." -ForegroundColor Green

# Build all services first
foreach ($service in $services) {
    Write-Host "Building $service..." -ForegroundColor Cyan
    $serviceDir = Join-Path $repoRoot "services\$service"
    Set-Location $serviceDir
    $binDir = Join-Path $serviceDir "bin"
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
    $serviceExe = Join-Path $binDir "service.exe"
    go build -o $serviceExe cmd/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to build $service" -ForegroundColor Red
        exit 1
    }
}

Write-Host "All services built successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Starting services in separate windows..." -ForegroundColor Green
Write-Host ""

# Start each service in a new PowerShell window
foreach ($service in $services) {
    $port = $portMap[$service]
    Write-Host "Starting $service on port $port..." -ForegroundColor Yellow
    $serviceDir = Join-Path $repoRoot "services\$service"
    $serviceExe = Join-Path $serviceDir "bin\service.exe"
    $cmd = "cd '$serviceDir'; Write-Host '===== $($service.ToUpper()) SERVICE =====' -ForegroundColor Green; & '$serviceExe'; Read-Host 'Press Enter to close'"
    Start-Process -FilePath "powershell" -ArgumentList "-NoExit", "-Command", $cmd
    Start-Sleep -Seconds 1
}

Write-Host ""
Write-Host "All services started! Check individual windows for startup messages." -ForegroundColor Green
Write-Host ""
Write-Host "Service ports:" -ForegroundColor Cyan
foreach ($s in $services) {
    $displayName = (Get-Culture).TextInfo.ToTitleCase($s)
    Write-Host "  $displayName`: $($portMap[$s])" -ForegroundColor Gray
}

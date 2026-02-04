# Build and start all microservices
$services = @(
    "auth", "schedule", "ticket", "fiscal", "payment", 
    "board", "notify", "audit", "document", "geo"
)

Write-Host "Building all services..." -ForegroundColor Green

# Build all services first
foreach ($service in $services) {
    Write-Host "Building $service..." -ForegroundColor Cyan
    Set-Location "e:\Projects\Personal\vokzal\services\$service"
    go build -o bin/service cmd/main.go
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
    $port = 8080 + ($services.IndexOf($service) + 1)
    Write-Host "Starting $service on port $port..." -ForegroundColor Yellow
    $serviceExe = "e:\Projects\Personal\vokzal\services\$service\bin\service.exe"
    $cmd = "cd 'e:\Projects\Personal\vokzal\services\$service'; Write-Host '===== $($service.ToUpper()) SERVICE =====' -ForegroundColor Green; & '$serviceExe'; Read-Host 'Press Enter to close'"
    Start-Process -FilePath "powershell" -ArgumentList "-NoExit", "-Command", $cmd
    Start-Sleep -Seconds 1
}

Write-Host ""
Write-Host "All services started! Check individual windows for startup messages." -ForegroundColor Green
Write-Host ""
Write-Host "Service ports:" -ForegroundColor Cyan
Write-Host "  Auth: 8081" -ForegroundColor Gray
Write-Host "  Schedule: 8082" -ForegroundColor Gray
Write-Host "  Ticket: 8083" -ForegroundColor Gray
Write-Host "  Fiscal: 8084" -ForegroundColor Gray
Write-Host "  Payment: 8085" -ForegroundColor Gray
Write-Host "  Board: 8086" -ForegroundColor Gray
Write-Host "  Notify: 8087" -ForegroundColor Gray
Write-Host "  Audit: 8098" -ForegroundColor Gray
Write-Host "  Document: 8089" -ForegroundColor Gray
Write-Host "  Geo: 8090" -ForegroundColor Gray

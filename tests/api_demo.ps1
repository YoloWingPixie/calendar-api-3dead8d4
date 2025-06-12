# Configuration
$scriptDir = Split-Path $PSScriptRoot -Parent
$DOCKER_COMPOSE_PATH = "$scriptDir\docker\docker-compose.yml"

# Get API configuration from environment variables first
$API_KEY = $env:BOOTSTRAP_ADMIN_KEY
$API_PORT = 8012

if ($API_KEY -and $API_PORT) {
    Write-Host "Using configuration from environment variables" -ForegroundColor Green
} else {
    if (-not (Test-Path $DOCKER_COMPOSE_PATH)) {
        Write-Host "Error: Could not find docker-compose.yml at $DOCKER_COMPOSE_PATH" -ForegroundColor Red
        exit 1
    }

    $API_PORT = (Select-String -Path $DOCKER_COMPOSE_PATH -Pattern "API_PORT:-(\d+)" | Select-Object -First 1).Matches.Groups[1].Value
    $API_KEY = (Select-String -Path $DOCKER_COMPOSE_PATH -Pattern "BOOTSTRAP_ADMIN_KEY:-([^}\s]+)" | Select-Object -First 1).Matches.Groups[1].Value

    if (-not $API_PORT -or -not $API_KEY) {
        Write-Host "Error: Could not extract API_PORT or API_KEY from docker-compose.yml" -ForegroundColor Red
        exit 1
    }
    Write-Host "Using configuration from docker-compose.yml" -ForegroundColor Yellow
}

$API_URL = "http://localhost:${API_PORT}"
$headers = @{
    "X-API-Key" = $API_KEY
    "Content-Type" = "application/json"
}

Write-Host "Using API URL: $API_URL"
Write-Host "Using API Key: $API_KEY"
Write-Host ""

Write-Host "Testing Calendar API..." -ForegroundColor Green

# 1. Test health endpoint (public)
Write-Host "`nTesting health endpoint..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$API_URL/health" -Method GET
    Write-Host "SUCCESS: Health: $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Health check failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. Test version endpoint (public)  
Write-Host "`nTesting version endpoint..." -ForegroundColor Yellow
try {
    $version = Invoke-RestMethod -Uri "$API_URL/version" -Method GET
    Write-Host "SUCCESS: Version: $($version.version)" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Version check failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. Test authentication (list events)
Write-Host "`nTesting authentication..." -ForegroundColor Yellow
try {
    $events = Invoke-RestMethod -Uri "$API_URL/api/events" -Method GET -Headers $headers
    Write-Host "SUCCESS: Authentication successful. Found $($events.count) events" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Authentication failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Helper function to display event details
function Write-EventDetails {
    param($event)
    Write-Host "`nEvent Details:" -ForegroundColor Cyan
    Write-Host "  Title: $($event.title)"
    Write-Host "  Description: $($event.description)"
    Write-Host "  Start Time: $([DateTime]::Parse($event.start_time).ToString('yyyy-MM-dd HH:mm'))"
    Write-Host "  End Time: $([DateTime]::Parse($event.end_time).ToString('yyyy-MM-dd HH:mm'))"
    Write-Host "  ID: $($event.id)"
    Write-Host "  Created At: $([DateTime]::Parse($event.created_at).ToString('yyyy-MM-dd HH:mm:ss'))"
    Write-Host "  Updated At: $([DateTime]::Parse($event.updated_at).ToString('yyyy-MM-dd HH:mm:ss'))"
    Write-Host ""
}

# 4. Create a test event
Write-Host "`nCreating test event..." -ForegroundColor Yellow
$eventData = @{
    title = "PowerShell Test Event"
    description = "Event created via PowerShell API test"
    start_time = (Get-Date).AddHours(1).ToString("yyyy-MM-ddTHH:mm:ssZ")
    end_time = (Get-Date).AddHours(2).ToString("yyyy-MM-ddTHH:mm:ssZ")
} | ConvertTo-Json

try {
    $newEvent = Invoke-RestMethod -Uri "$API_URL/api/events" -Method POST -Headers $headers -Body $eventData
    Write-Host "SUCCESS: Event created with ID: $($newEvent.id)" -ForegroundColor Green
    Write-EventDetails $newEvent
    $eventId = $newEvent.id
} catch {
    Write-Host "ERROR: Event creation failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 5. Get the created event
Write-Host "`nRetrieving created event..." -ForegroundColor Yellow
try {
    $retrievedEvent = Invoke-RestMethod -Uri "$API_URL/api/events/$eventId" -Method GET -Headers $headers
    Write-Host "SUCCESS: Retrieved event: $($retrievedEvent.title)" -ForegroundColor Green
    Write-EventDetails $retrievedEvent
} catch {
    Write-Host "ERROR: Event retrieval failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 6. Update the event
Write-Host "`nUpdating event..." -ForegroundColor Yellow
$updateData = @{
    title = "Updated PowerShell Test Event"
    description = "This event was updated via PowerShell"
    start_time = $newEvent.start_time
    end_time = $newEvent.end_time
} | ConvertTo-Json

try {
    $updatedEvent = Invoke-RestMethod -Uri "$API_URL/api/events/$eventId" -Method PUT -Headers $headers -Body $updateData
    Write-Host "SUCCESS: Event updated: $($updatedEvent.title)" -ForegroundColor Green
    Write-EventDetails $updatedEvent
} catch {
    Write-Host "ERROR: Event update failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 7. List all events again
Write-Host "`nListing all events..." -ForegroundColor Yellow
try {
    $allEvents = Invoke-RestMethod -Uri "$API_URL/api/events" -Method GET -Headers $headers
    Write-Host "SUCCESS: Total events: $($allEvents.count)" -ForegroundColor Green
    
    # Create table header
    Write-Host "`nTitle`tStart Time`tEnd Time`tID" -ForegroundColor Cyan
    Write-Host "-----`t----------`t--------`t--" -ForegroundColor Cyan
    
    # Display each event in table format
    foreach ($event in $allEvents.events) {
        $startTime = [DateTime]::Parse($event.start_time).ToString("yyyy-MM-dd HH:mm")
        $endTime = [DateTime]::Parse($event.end_time).ToString("yyyy-MM-dd HH:mm")
        Write-Host "$($event.title)`t$startTime`t$endTime`t$($event.id)" -ForegroundColor Cyan
    }
} catch {
    Write-Host "ERROR: Event listing failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 8. Delete the test event
Write-Host "`nCleaning up - deleting test event..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "$API_URL/api/events/$eventId" -Method DELETE -Headers $headers
    Write-Host "SUCCESS: Test event deleted successfully" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Event deletion failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nAPI testing completed!" -ForegroundColor Green 
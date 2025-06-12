#!/bin/bash

# Configuration
DOCKER_COMPOSE_PATH="$(dirname "$0")/../docker/docker-compose.yml"
API_PORT=$(grep -oP 'API_PORT:-\K\d+' "$DOCKER_COMPOSE_PATH" | head -n 1)
API_KEY=$(grep -oP 'BOOTSTRAP_ADMIN_KEY:-\K[^\s]+' "$DOCKER_COMPOSE_PATH" | head -n 1)
API_URL="http://localhost:${API_PORT}"

echo "Using API URL: $API_URL"
echo "Using API Key: $API_KEY"
echo

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper function to print section headers
print_header() {
    echo -e "\n${CYAN}=== $1 ===${NC}"
}

# Helper function to display event details
print_event_details() {
    local event=$1
    echo -e "\n${CYAN}Event Details:${NC}"
    echo "  Title: $(echo "$event" | jq -r '.title')"
    echo "  Description: $(echo "$event" | jq -r '.description')"
    echo "  Start Time: $(date -d "$(echo "$event" | jq -r '.start_time')" '+%Y-%m-%d %H:%M')"
    echo "  End Time: $(date -d "$(echo "$event" | jq -r '.end_time')" '+%Y-%m-%d %H:%M')"
    echo "  ID: $(echo "$event" | jq -r '.id')"
    echo "  Created At: $(date -d "$(echo "$event" | jq -r '.created_at')" '+%Y-%m-%d %H:%M:%S')"
    echo "  Updated At: $(date -d "$(echo "$event" | jq -r '.updated_at')" '+%Y-%m-%d %H:%M:%S')"
    echo
}

# Helper function to make API calls
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local api_key=${4:-$API_KEY}

    echo -e "${YELLOW}Request: $method $endpoint${NC}"
    if [ ! -z "$data" ]; then
        echo -e "${YELLOW}Data: $data${NC}"
    fi

    # Make the request and capture both status code and response
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" \
            -H "X-API-Key: $api_key" \
            "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" \
            -H "Content-Type: application/json" \
            -H "X-API-Key: $api_key" \
            -d "$data" \
            "$API_URL$endpoint")
    fi

    # Split response into body and status code
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    # Print response
    if [ "$status_code" -ge 200 ] && [ "$status_code" -lt 300 ]; then
        echo -e "${GREEN}Response (Status: $status_code)${NC}"
        if command -v jq >/dev/null 2>&1; then
            echo "$body" | jq '.'
        else
            echo "$body"
        fi
    else
        echo -e "${RED}Response (Status: $status_code)${NC}"
        echo "$body"
    fi
    echo
}

# Test health endpoint
print_header "Testing Health Endpoint"
make_request "GET" "/health"

# Test version endpoint
print_header "Testing Version Endpoint"
make_request "GET" "/version"

# Test authentication (list events)
print_header "Testing Authentication"
make_request "GET" "/api/events"

# Create test event
print_header "Creating Test Event"
event_data='{
    "title": "Bash Test Event",
    "description": "Event created via Bash API test",
    "start_time": "'$(date -u -d "+1 hour" "+%Y-%m-%dT%H:%M:%SZ")'",
    "end_time": "'$(date -u -d "+2 hours" "+%Y-%m-%dT%H:%M:%SZ")'"
}'
response=$(make_request "POST" "/api/events" "$event_data")
event_id=$(echo "$response" | jq -r '.id')

# Get created event
print_header "Retrieving Created Event"
response=$(make_request "GET" "/api/events/$event_id")
print_event_details "$response"

# Update event
print_header "Updating Event"
update_data='{
    "title": "Updated Bash Test Event",
    "description": "This event was updated via Bash",
    "start_time": "'$(echo "$response" | jq -r '.start_time')'",
    "end_time": "'$(echo "$response" | jq -r '.end_time')'"
}'
response=$(make_request "PUT" "/api/events/$event_id" "$update_data")
print_event_details "$response"

# List all events
print_header "Listing All Events"
response=$(make_request "GET" "/api/events")
echo -e "${GREEN}Total events: $(echo "$response" | jq '.events | length')${NC}"

# Create table header
echo -e "\n${CYAN}Title\tStart Time\tEnd Time\tID${NC}"
echo -e "${CYAN}-----\t----------\t--------\t--${NC}"

# Display each event in table format
echo "$response" | jq -r '.events[] | [.title, (.start_time | strptime("%Y-%m-%dT%H:%M:%SZ") | strftime("%Y-%m-%d %H:%M")), (.end_time | strptime("%Y-%m-%dT%H:%M:%SZ") | strftime("%Y-%m-%d %H:%M")), .id] | @tsv' | while IFS=$'\t' read -r title start_time end_time id; do
    echo -e "${CYAN}$title\t$start_time\t$end_time\t$id${NC}"
done

# Delete test event
print_header "Cleaning Up - Deleting Test Event"
make_request "DELETE" "/api/events/$event_id"

echo -e "\n${GREEN}API testing completed!${NC}" 
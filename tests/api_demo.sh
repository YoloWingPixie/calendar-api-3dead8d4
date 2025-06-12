#!/bin/bash

# -----------------------------------------------------------------------------
# Bash Calendar-API demo script
# -----------------------------------------------------------------------------
#  ‑ Exercises the public and authenticated endpoints of the Calendar API
#  ‑ Shows nicely formatted coloured output (similar to api_demo.ps1)
# -----------------------------------------------------------------------------

set -euo pipefail

# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------
API_PORT=8012
API_KEY=${BOOTSTRAP_ADMIN_KEY:-""}
API_URL="http://localhost:${API_PORT}"

# -----------------------------------------------------------------------------
# Colour definitions
# -----------------------------------------------------------------------------
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No colour

# -----------------------------------------------------------------------------
# Helper functions
# -----------------------------------------------------------------------------
function print_header() {
  echo -e "${YELLOW}\n$*...${NC}"
}

function success() {
  echo -e "${GREEN}SUCCESS: $*${NC}"
}

function error() {
  echo -e "${RED}ERROR: $*${NC}"
}

function event_details() {
  local json="$1"
  echo -e "${CYAN}\nEvent Details:${NC}"
  echo -e "  Title:        $(echo "$json" | jq -r '.title')"
  echo -e "  Description:  $(echo "$json" | jq -r '.description')"
  echo -e "  Start Time:   $(echo "$json" | jq -r '.start_time')"
  echo -e "  End Time:     $(echo "$json" | jq -r '.end_time')"
  echo -e "  ID:           $(echo "$json" | jq -r '.id')"
  echo -e "  Created At:   $(echo "$json" | jq -r '.created_at')"
  echo -e "  Updated At:   $(echo "$json" | jq -r '.updated_at')\n"
}

# Displays the HTTP request being made
function show_request() {
  local method="$1"; shift
  local url="$1"; shift
  local body="${1:-}"
  echo -e "${CYAN}--> ${method} ${url}${NC}"
  if [[ -n "$body" ]]; then
    # pretty-print JSON body if possible
    if command -v jq >/dev/null 2>&1; then
      echo "$body" | jq -C .
    else
      echo -e "$body"
    fi
  fi
}

# Pretty-print an events array as table
function display_events() {
  local json="$1"
  local events_present=$(echo "$json" | jq -r '.events | length')
  if [[ "$events_present" -eq 0 ]]; then
    echo -e "${YELLOW}(no events)${NC}"
    return
  fi

  printf "\n${CYAN}%-16s %-16s %-36s %-20s %-30s${NC}\n" "Start Time" "End Time" "ID" "Title" "Description"
  printf "${CYAN}%-16s %-16s %-36s %-20s %-30s${NC}\n" "----------" "--------" "--" "-----" "-----------"

  echo "$json" | jq -c '.events[]' | while read -r evt; do
    start=$(echo "$evt" | jq -r '.start_time' | cut -c1-16)
    end=$(echo "$evt" | jq -r '.end_time' | cut -c1-16)
    id=$(echo "$evt" | jq -r '.id')
    title=$(echo "$evt" | jq -r '.title')
    desc=$(echo "$evt" | jq -r '.description')
    printf "${GREEN}%-16s %-16s %-36s %-20s %-30s${NC}\n" "$start" "$end" "$id" "$title" "$desc"
  done
  echo
}

# -----------------------------------------------------------------------------
# Start tests
# -----------------------------------------------------------------------------

print_header "Testing API at ${API_URL}"

# 1. Health check (public)
print_header "Testing health endpoint"
show_request "GET" "${API_URL}/health"
response=$(curl -s "${API_URL}/health")
status=$(echo "$response" | jq -r '.status // empty')
if [[ -n "$status" ]]; then
  success "Health: $status"
else
  error "Health check failed: $response"
  exit 1
fi

# 2. Version check (public)
print_header "Testing version endpoint"
show_request "GET" "${API_URL}/version"
response=$(curl -s "${API_URL}/version")
version=$(echo "$response" | jq -r '.version // empty')
if [[ -n "$version" ]]; then
  success "Version: $version"
else
  error "Version check failed: $response"
fi

# 3. List events (authentication test)
print_header "Testing authentication (listing events)"
show_request "GET" "${API_URL}/api/events"
response=$(curl -s -H "X-API-Key: $API_KEY" "${API_URL}/api/events") || {
  error "Authentication failed"
  exit 1
}
count=$(echo "$response" | jq -r '.events | length')
success "Authentication successful. Found ${count} events"

# Show table of returned events
display_events "$response"

# 4. Create a test event
print_header "Creating test event"
start_time=$(date -u -v+1H +"%Y-%m-%dT%H:%M:%SZ")
end_time=$(date -u -v+2H +"%Y-%m-%dT%H:%M:%SZ")

event_data=$(cat <<EOF
{
  "title": "Bash Test Event",
  "description": "Event created via Bash API test",
  "start_time": "${start_time}",
  "end_time": "${end_time}"
}
EOF
)

# Show the request after body is prepared
show_request "POST" "${API_URL}/api/events" "$event_data"

response=$(curl -s -H "Content-Type: application/json" -H "X-API-Key: $API_KEY" -d "${event_data}" "${API_URL}/api/events") || {
  error "Event creation failed"
  exit 1
}

success "Event created with ID: $(echo "$response" | jq -r '.id')"
event_details "$response"
event_id=$(echo "$response" | jq -r '.id')

# 5. Retrieve the created event
print_header "Retrieving created event"
show_request "GET" "${API_URL}/api/events/${event_id}"
response=$(curl -s -H "X-API-Key: $API_KEY" "${API_URL}/api/events/${event_id}") || {
  error "Event retrieval failed"
  exit 1
}

success "Retrieved event: $(echo "$response" | jq -r '.title')"
event_details "$response"

# 6. Update the event
print_header "Updating event"
update_data=$(cat <<EOF
{
  "title": "Updated Bash Test Event",
  "description": "This event was updated via Bash",
  "start_time": "${start_time}",
  "end_time": "${end_time}"
}
EOF
)

# Show update request
show_request "PUT" "${API_URL}/api/events/${event_id}" "$update_data"

response=$(curl -s -H "Content-Type: application/json" -H "X-API-Key: $API_KEY" -X PUT -d "${update_data}" "${API_URL}/api/events/${event_id}") || {
  error "Event update failed"
  exit 1
}

success "Event updated: $(echo "$response" | jq -r '.title')"
event_details "$response"

# 7. List all events
print_header "Listing all events"
show_request "GET" "${API_URL}/api/events"
response=$(curl -s -H "X-API-Key: $API_KEY" "${API_URL}/api/events")
count=$(echo "$response" | jq -r '.events | length')

success "Total events: ${count}"

# Display events table with description
display_events "$response"

# 8. Delete the test event
print_header "Cleaning up – deleting test event"
show_request "DELETE" "${API_URL}/api/events/${event_id}"
if curl -s -H "X-API-Key: $API_KEY" -X DELETE "${API_URL}/api/events/${event_id}" >/dev/null; then
  success "Test event deleted successfully"
else
  error "Event deletion failed"
fi

print_header "API testing completed!" 
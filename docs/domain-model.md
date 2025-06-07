# Domain Model

## Entities

### Calendar
**Description:** Represents a collection of `CalendarEvents`.
**Attributes:**
    - CalendarId: Unique identifier for the calendar. (PK)
    - OwnerUserId: Identifier of the `User` who owns this calendar (FK)
    - EditorIds: List of `User` Id who are able to create, update, or delete `Event` from this calendar (FK)
    - ReaderIds: List of `User` Id who are able to read `Event` from this calendar (FK)
    - PublicRead: Boolean. If true, any `User` can read this calendar.
    - PublicWrite: Boolean. If true, any `User` can perform Editor actions on this calendar's `Event`.
    - CalendarName: Human readable name of the calendar.
    - EventIds: List of `Event` associated with the calendar. (FK)
    - CreatedAt: ISO 8601 timestamp
    - UpdatedAt: ISO 8601 timestamp
**Behaviors:**
    - Add `Event`
    - Remove `Event`
    - List `Event`(s)
    - Add Reader
    - Remove Reader
    - Add Editor
    - Remove Editor
    - Set Public Read On/Off
    - Set Public Write On/Off

### CalendarEvent
**Description:** Represents a scheduled event within a `Calendar`. Each event has a specific time range and can be managed by calendar owners and editors.
**Attributes:**
    - EventId: Unique identifier for the event (PK)
    - CalendarId: Identifier of the `Calendar` this event belongs to (FK)
    - CreatorUserId: Identifier of the `User` who created this event (FK)
    - Title: Human readable name of the event
    - Description: Detailed description of the event
    - StartTime: ISO 8601 timestamp when the event begins
    - EndTime: ISO 8601 timestamp when the event ends
    - IsAllDay: Boolean indicating if the event spans the entire day
    - CreatedAt: ISO 8601 timestamp
    - UpdatedAt: ISO 8601 timestamp
**Behaviors:**
    - Create event
    - Read event details
    - Update event details
    - Delete event

### User
**Description:** Represents an individual who uses the calendar API. Users can own calendars, and CRUD events.
**Attributes:**
    - UserId: Unique identifier for the user (PK)
    - Username: The name of the user
    - AccessKey: Access token for the API
    - CreatedAt: ISO 8601 timestamp
    - UpdatedAt: ISO 8601 timestamp
    - OwnedCalendarIds: List of `Calendar` owned by the user
**Responsibilities:**
    - Owns `Calendar` entities
    - Creates, reads, updates, or deletes `CalendarEvents`
**Behaviors:**
    - Create `Calendar`
    - Delete `Calendar`
    - Own `Calendar`
    - Create `CalendarEvent` in `Calendar`
    - Read all `CalendarEvents` in a `Calendar`
    - Read all `CalendarEvents` in all `Calendars`
    - Update `CalendarEvent`
    - Delete `CalendarEvent`

## Invariants
INV-001: CalendarEvent Time Integrity: A CalendarEvent's Start Time must be before its End Time
INV-002: CalendarEvent Creator Access: A User must have Editor or Owner permissions on a Calendar to create CalendarEvents in it
INV-003: CalendarEvent Modification Access: Only the Creator, Calendar Owner, or Calendar Editors can modify a CalendarEvent
INV-004: CalendarEvent Deletion Access: Only the Creator, Calendar Owner, or Calendar Editors can delete a CalendarEvent
INV-005: CalendarEvent Read Access: A User must have Reader, Editor, or Owner permissions on a Calendar to read its CalendarEvents
INV-006: CalendarEvent All-Day Constraint: If IsAllDay is true, StartTime must be at 00:00:00 and EndTime must be at 23:59:59 of the same day
INV-007: CalendarEvent Duration: A CalendarEvent must have a minimum duration of 1 minute
INV-008: CalendarEvent Creator Existence: The CreatorUserId must reference an existing User
INV-009: CalendarEvent Calendar Existence: The CalendarId must reference an existing Calendar

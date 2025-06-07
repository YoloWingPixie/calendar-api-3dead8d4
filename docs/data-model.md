# Data Model

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    access_key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Calendars Table
```sql
CREATE TABLE calendars (
    calendar_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id UUID NOT NULL REFERENCES users(user_id),
    calendar_name VARCHAR(255) NOT NULL,
    editor_ids UUID[] NOT NULL DEFAULT '{}',
    reader_ids UUID[] NOT NULL DEFAULT '{}',
    public_read BOOLEAN NOT NULL DEFAULT false,
    public_write BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_editor_ids CHECK (
        NOT EXISTS (
            SELECT 1 FROM unnest(editor_ids) AS editor_id
            WHERE NOT EXISTS (
                SELECT 1 FROM users WHERE user_id = editor_id
            )
        )
    ),
    CONSTRAINT valid_reader_ids CHECK (
        NOT EXISTS (
            SELECT 1 FROM unnest(reader_ids) AS reader_id
            WHERE NOT EXISTS (
                SELECT 1 FROM users WHERE user_id = reader_id
            )
        )
    )
);
```

### Calendar Events Table
```sql
CREATE TABLE calendar_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_id UUID NOT NULL REFERENCES calendars(calendar_id),
    creator_user_id UUID NOT NULL REFERENCES users(user_id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    is_all_day BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT end_after_start CHECK (end_time > start_time),
    CONSTRAINT min_duration CHECK (end_time - start_time >= interval '1 minute'),
    CONSTRAINT all_day_constraint CHECK (
        (is_all_day = false) OR 
        (is_all_day = true AND 
         start_time::time = '00:00:00' AND 
         end_time::time = '23:59:59' AND 
         start_time::date = end_time::date)
    )
);
```

## Data Transfer Objects (DTOs)

### User DTOs
```python
class UserCreate(BaseModel):
    username: str

class UserResponse(BaseModel):
    user_id: UUID
    username: str
    owned_calendar_ids: list[UUID]
    created_at: datetime
    updated_at: datetime

class UserWithAccessKey(UserResponse):
    access_key: str
```

### Calendar DTOs
```python
class CalendarCreate(BaseModel):
    calendar_name: str
    editor_ids: list[UUID] = Field(default_factory=list)
    reader_ids: list[UUID] = Field(default_factory=list)
    public_read: bool = False
    public_write: bool = False

class CalendarUpdate(BaseModel):
    calendar_name: str | None = None
    editor_ids: list[UUID] | None = None
    reader_ids: list[UUID] | None = None
    public_read: bool | None = None
    public_write: bool | None = None

class CalendarResponse(BaseModel):
    calendar_id: UUID
    owner_user_id: UUID
    calendar_name: str
    editor_ids: list[UUID]
    reader_ids: list[UUID]
    public_read: bool
    public_write: bool
    created_at: datetime
    updated_at: datetime
```

### Calendar Event DTOs
```python
class CalendarEventCreate(BaseModel):
    title: str
    description: str | None = None
    start_time: datetime
    end_time: datetime
    is_all_day: bool = False

class CalendarEventUpdate(BaseModel):
    title: str | None = None
    description: str | None = None
    start_time: datetime | None = None
    end_time: datetime | None = None
    is_all_day: bool | None = None

class CalendarEventResponse(BaseModel):
    event_id: UUID
    calendar_id: UUID
    creator_user_id: UUID
    title: str
    description: str | None
    start_time: datetime
    end_time: datetime
    is_all_day: bool
    created_at: datetime
    updated_at: datetime
```

### Minimum Required Event DTOs (for API compatibility)
```python
class EventCreate(BaseModel):
    title: str
    description: str | None = None
    start_time: datetime
    end_time: datetime

class EventUpdate(BaseModel):
    title: str | None = None
    description: str | None = None
    start_time: datetime | None = None
    end_time: datetime | None = None

class EventResponse(BaseModel):
    id: UUID
    title: str
    description: str | None
    start_time: datetime
    end_time: datetime
    created_at: datetime
    updated_at: datetime
```

## API Contracts

### Error Responses
```python
class ErrorResponse(BaseModel):
    error: str
    detail: str | None = None
    code: str
```

### Pagination
```python
class PaginationParams(BaseModel):
    page: int = 1
    page_size: int = 20

class PaginatedResponse(BaseModel):
    items: list[Any]
    total: int
    page: int
    page_size: int
    total_pages: int
```

### Health Check
```python
class HealthResponse(BaseModel):
    status: str
    version: str
    timestamp: datetime
``` 
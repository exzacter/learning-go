# REST API with SQLite - Architecture Documentation

## Folder Structure
```
restAPIWithSqLite/
├── cmd/
│   ├── api/
│   │   ├── main.go          (Application entry point)
│   │   ├── server.go        (HTTP server setup)
│   │   ├── routes.go        (Route definitions)
│   │   ├── auth.go          (User authentication handlers)
│   │   ├── events.go        (Event CRUD handlers)
│   │   └── attendees.go     (Attendee handlers)
│   └── migrate/
│       └── main.go          (Database migration tool)
├── internal/
│   ├── database/
│   │   ├── models.go        (Models factory)
│   │   ├── users.go         (User DB operations)
│   │   ├── events.go        (Event DB operations)
│   │   └── attendees.go     (Attendee DB operations)
│   └── env/
│       └── env.go           (Environment helpers)
└── data.db                  (SQLite database)
```

---

## Function Connection Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              APPLICATION STARTUP                                 │
└─────────────────────────────────────────────────────────────────────────────────┘

cmd/api/main.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  main()                                                                          │
│    │                                                                             │
│    ├──► sql.Open("sqlite3", "./data.db")                                        │
│    │                                                                             │
│    ├──► database.NewModels(db) ─────────────────► internal/database/models.go   │
│    │                                               └─► Creates UserModel         │
│    │                                               └─► Creates EventModel        │
│    │                                               └─► Creates AttendeeModel     │
│    │                                                                             │
│    ├──► env.GetEnvInt("PORT") ──────────────────► internal/env/env.go           │
│    ├──► env.GetEnvString("JWT_SECRET") ─────────► internal/env/env.go           │
│    │                                                                             │
│    └──► app.serve() ────────────────────────────► cmd/api/server.go             │
└─────────────────────────────────────────────────────────────────────────────────┘

cmd/api/server.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  serve()                                                                         │
│    │                                                                             │
│    ├──► app.routes() ───────────────────────────► cmd/api/routes.go             │
│    │                                                                             │
│    └──► http.Server.ListenAndServe()                                            │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                                  HTTP ROUTES                                     │
└─────────────────────────────────────────────────────────────────────────────────┘

cmd/api/routes.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  routes() → gin.Engine                                                           │
│    │                                                                             │
│    └─► /api/v1                                                                   │
│         │                                                                        │
│         ├─► POST   /auth/register ─────────► registerUser()      [auth.go]      │
│         │                                                                        │
│         ├─► POST   /events ────────────────► createEvent()       [events.go]    │
│         ├─► GET    /events ────────────────► getAllEvents()      [events.go]    │
│         ├─► GET    /events/:id ────────────► getEvent()          [events.go]    │
│         ├─► PUT    /events/:id ────────────► updateEvent()       [events.go]    │
│         ├─► DELETE /events/:id ────────────► deleteEvent()       [events.go]    │
│         │                                                                        │
│         ├─► POST   /events/:id/attendees/:userId ► addAttendeeToEvent()         │
│         ├─► GET    /events/:id/attendees ─────────► getAttendeesForEvent()      │
│         ├─► DELETE /events/:id/attendees/:userId ► deleteAttendeeFromEvent()    │
│         │                                                                        │
│         └─► GET    /attendees/:id/events ──► getEventsByAttendee() [attendees]  │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                          HTTP HANDLERS → DATABASE                                │
└─────────────────────────────────────────────────────────────────────────────────┘

cmd/api/auth.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  registerUser(c *gin.Context)                                                    │
│    │                                                                             │
│    ├──► c.ShouldBindJSON(&registerRequest)                                      │
│    ├──► bcrypt.GenerateFromPassword()                                           │
│    └──► app.models.Users.Insert() ─────────► internal/database/users.go         │
│                                              └─► INSERT INTO users              │
└─────────────────────────────────────────────────────────────────────────────────┘

cmd/api/events.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  createEvent(c)                                                                  │
│    └──► app.models.Events.Insert() ────────► internal/database/events.go        │
│                                              └─► INSERT INTO events             │
├─────────────────────────────────────────────────────────────────────────────────┤
│  getAllEvents(c)                                                                 │
│    └──► app.models.Events.GetAll() ────────► internal/database/events.go        │
│                                              └─► SELECT * FROM events           │
├─────────────────────────────────────────────────────────────────────────────────┤
│  getEvent(c)                                                                     │
│    └──► app.models.Events.Get(id) ─────────► internal/database/events.go        │
│                                              └─► SELECT * FROM events WHERE id  │
├─────────────────────────────────────────────────────────────────────────────────┤
│  updateEvent(c)                                                                  │
│    ├──► app.models.Events.Get(id) ─────────► (verify exists)                    │
│    └──► app.models.Events.Update() ────────► internal/database/events.go        │
│                                              └─► UPDATE events SET ...          │
├─────────────────────────────────────────────────────────────────────────────────┤
│  deleteEvent(c)                                                                  │
│    └──► app.models.Events.Delete(id) ──────► internal/database/events.go        │
│                                              └─► DELETE FROM events WHERE id    │
├─────────────────────────────────────────────────────────────────────────────────┤
│  addAttendeeToEvent(c)                                                           │
│    ├──► app.models.Events.Get(eventId) ────► (verify event exists)              │
│    ├──► app.models.Users.Get(userId) ──────► (verify user exists)               │
│    ├──► app.models.Attendees.GetByEventAndAttendee() ► (check duplicate)        │
│    └──► app.models.Attendees.Insert() ─────► internal/database/attendees.go     │
│                                              └─► INSERT INTO attendees          │
├─────────────────────────────────────────────────────────────────────────────────┤
│  getAttendeesForEvent(c)                                                         │
│    └──► app.models.Attendees.GetAttendeesByEvent() ► attendees.go               │
│                                              └─► SELECT users JOIN attendees    │
├─────────────────────────────────────────────────────────────────────────────────┤
│  deleteAttendeeFromEvent(c)                                                      │
│    └──► app.models.Attendees.Delete() ─────► internal/database/attendees.go     │
│                                              └─► DELETE FROM attendees          │
└─────────────────────────────────────────────────────────────────────────────────┘

cmd/api/attendees.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  getEventsByAttendee(c)                                                          │
│    └──► app.models.Attendees.GetEventsByAttendee() ► attendees.go               │
│                                              └─► SELECT events JOIN attendees   │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                            DATABASE LAYER                                        │
└─────────────────────────────────────────────────────────────────────────────────┘

internal/database/models.go
┌─────────────────────────────────────────────────────────────────────────────────┐
│  NewModels(db *sql.DB) → Models                                                  │
│    └─► Returns:  Models{                                                         │
│                    Users:     UserModel{DB: db}                                  │
│                    Events:    EventModel{DB: db}                                 │
│                    Attendees: AttendeeModel{DB: db}                              │
│                  }                                                               │
└─────────────────────────────────────────────────────────────────────────────────┘

internal/database/users.go
┌───────────────────────────────────┐          ┌───────────────────┐
│  UserModel                        │          │  User struct      │
│    ├─► Insert(user) error         │◄────────►│   Id       int    │
│    └─► Get(id) (*User, error)     │          │   Email    string │
└───────────────────────────────────┘          │   Name     string │
                                               │   Password string │
                                               └───────────────────┘

internal/database/events.go
┌───────────────────────────────────┐          ┌───────────────────┐
│  EventModel                       │          │  Event struct     │
│    ├─► Insert(event) error        │◄────────►│   Id          int │
│    ├─► GetAll() ([]*Event, error) │          │   OwnerId     int │
│    ├─► Get(id) (*Event, error)    │          │   Name     string │
│    ├─► Update(event) error        │          │   Description str │
│    └─► Delete(id) error           │          │   Date     string │
└───────────────────────────────────┘          │   Location string │
                                               └───────────────────┘

internal/database/attendees.go
┌──────────────────────────────────────────┐   ┌───────────────────┐
│  AttendeeModel                           │   │  Attendee struct  │
│    ├─► Insert(attendee) (*Attendee, err) │◄─►│   Id       int    │
│    ├─► GetByEventAndAttendee(eId, uId)   │   │   UserId   int    │
│    ├─► GetAttendeesByEvent(eId) []*User  │   │   EventId  int    │
│    ├─► GetEventsByAttendee(uId) []*Event │   └───────────────────┘
│    └─► Delete(userId, eventId) error     │
└──────────────────────────────────────────┘

internal/env/env.go
┌───────────────────────────────────────┐
│  GetEnvString(key, default) string    │
│  GetEnvInt(key, default) int          │
└───────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                            DATABASE SCHEMA                                       │
└─────────────────────────────────────────────────────────────────────────────────┘

┌────────────────┐     ┌────────────────┐     ┌────────────────┐
│     users      │     │    events      │     │   attendees    │
├────────────────┤     ├────────────────┤     ├────────────────┤
│ id (PK)        │◄────│ owner_id (FK)  │     │ id (PK)        │
│ email          │     │ id (PK)        │◄────│ event_id (FK)  │
│ name           │◄────┼────────────────┼─────│ user_id (FK)   │
│ password       │     │ name           │     └────────────────┘
└────────────────┘     │ description    │
                       │ date           │
                       │ location       │
                       └────────────────┘
```

---

## Summary of Function Purposes

| File | Function | Purpose |
|------|----------|---------|
| `main.go` | `main()` | Initialize DB, models, env vars, start server |
| `server.go` | `serve()` | Configure & run HTTP server with timeouts |
| `routes.go` | `routes()` | Define all API endpoints & wire handlers |
| `auth.go` | `registerUser()` | Hash password, create new user |
| `events.go` | `createEvent()` | Create new event |
| `events.go` | `getAllEvents()` | List all events |
| `events.go` | `getEvent()` | Get single event by ID |
| `events.go` | `updateEvent()` | Update event details |
| `events.go` | `deleteEvent()` | Remove event |
| `events.go` | `addAttendeeToEvent()` | Register user for event |
| `events.go` | `getAttendeesForEvent()` | List event attendees |
| `events.go` | `deleteAttendeeFromEvent()` | Remove user from event |
| `attendees.go` | `getEventsByAttendee()` | List events for a user |
| `models.go` | `NewModels()` | Factory to create all DB models |
| `users.go` | `Insert()/Get()` | User CRUD operations |
| `events.go` | `Insert()/GetAll()/Get()/Update()/Delete()` | Event CRUD |
| `attendees.go` | `Insert()/Delete()/GetBy*()` | Attendee junction operations |
| `env.go` | `GetEnvString()/GetEnvInt()` | Read environment variables |

---

## API Endpoints

| Method | Endpoint | Handler | Description |
|--------|----------|---------|-------------|
| POST | `/api/v1/auth/register` | `registerUser` | Register a new user |
| POST | `/api/v1/events` | `createEvent` | Create a new event |
| GET | `/api/v1/events` | `getAllEvents` | Get all events |
| GET | `/api/v1/events/:id` | `getEvent` | Get event by ID |
| PUT | `/api/v1/events/:id` | `updateEvent` | Update event by ID |
| DELETE | `/api/v1/events/:id` | `deleteEvent` | Delete event by ID |
| POST | `/api/v1/events/:id/attendees/:userId` | `addAttendeeToEvent` | Add attendee to event |
| GET | `/api/v1/events/:id/attendees` | `getAttendeesForEvent` | Get all attendees for event |
| DELETE | `/api/v1/events/:id/attendees/:userId` | `deleteAttendeeFromEvent` | Remove attendee from event |
| GET | `/api/v1/attendees/:id/events` | `getEventsByAttendee` | Get all events for attendee |

---

## Architecture Pattern

This project follows a **clean layered architecture**:

1. **Presentation Layer** (`cmd/api/`) - HTTP handlers using Gin framework
2. **Data Access Layer** (`internal/database/`) - Database models and CRUD operations
3. **Configuration Layer** (`internal/env/`) - Environment variable management

All database operations use:
- Context with 3-second timeout
- Parameterized queries (SQL injection prevention)
- Proper error handling with `sql.ErrNoRows` checks

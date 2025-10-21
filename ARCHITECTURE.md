# LACPA Architecture

## System Overview

LACPA (Lebanese Association of Certified Public Accountants) is a full-stack web application built with Go Fiber backend and HTML/HTMX frontend. The system manages individual members, firm members, council positions, events, and authentication with a clear separation of concerns following clean architecture principles.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                      Client (Browser)                             │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │            LACPA_Web (Frontend)                            │  │
│  │  • index.html - Main entry point                           │  │
│  │  • Tailwind CSS - Utility-first styling                    │  │
│  │  • HTMX - Dynamic page updates                             │  │
│  │  • Components: header, footer, localnav, buttons           │  │
│  │  • Pages: landing, board members, events, members          │  │
│  │  • app.js - Navigation & interactions                      │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST API + Server-Side Rendered Pages
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                       Go Fiber Backend                            │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  main.go - Application Entry Point                         │  │
│  │  • MongoDB connection (LACPA database)                     │  │
│  │  • Fiber app with CORS & logging middleware                │  │
│  │  • HTML template engine                                    │  │
│  │  • Route registration                                      │  │
│  │  • Static file serving                                     │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  routes/ - Route Configuration                             │  │
│  │  • routes.go - Main route setup                            │  │
│  │  • main_page.go - Public page routes                       │  │
│  │  • Groups: /api, /auth, /discover, /admin                  │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  handler/ - HTTP Request Handlers                          │  │
│  │  • main-page.go - Landing & static pages                   │  │
│  │  • council_handler.go - Council & board operations         │  │
│  │  • member_handler.go - Individual member operations        │  │
│  │  • firm_handler.go - Firm member operations                │  │
│  │  • event_handler.go - Event management                     │  │
│  │  • auth_handler.go - Authentication (login, register)      │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  repository/ - Data Access Layer                           │  │
│  │  • repository.go - Main repository interface               │  │
│  │  • council_repository.go - Council & position queries      │  │
│  │  • member_repository.go - Member CRUD operations           │  │
│  │  • firm_repository.go - Firm management                    │  │
│  │  • event_repository.go - Event operations                  │  │
│  │  • auth_repository.go - User authentication                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  models/ - Data Models                                     │  │
│  │  • council.go - Council, CouncilPosition, Composition      │  │
│  │  • individual_member.go - Individual member entity         │  │
│  │  • firm_member.go - Firm entity                            │  │
│  │  • event.go - Event entity                                 │  │
│  │  • user.go - User authentication                           │  │
│  │  • landing_page.go - Landing page content                  │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  utils/ - Utility Functions                                │  │
│  │  • response.go - Standard API responses                    │  │
│  │  • request.go - Request parsing helpers                    │  │
│  │  • validation.go - Input validation                        │  │
│  │  • config.go - Configuration helpers                       │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  config/ - Configuration                                   │  │
│  │  • database.go - MongoDB connection setup                  │  │
│  └────────────────────────────────────────────────────────────┘  │
│                              │                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  templates/ - Server-Side HTML Templates                   │  │
│  │  • LACPA/board_members/ - Board member pages               │  │
│  │  • LACPA/main_LACPA/ - Main landing pages                  │  │
│  │  • Admin_Dashboard/ - Admin interface                      │  │
│  │  • views/layouts/ - Shared layouts                         │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
                              │
                              │ MongoDB Go Driver
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                     MongoDB Database (LACPA)                      │
│  Collections:                                                     │
│  • individual_members - Individual CPA members (12 docs)          │
│  • firm_members - Firm/company members (10 docs)                  │
│  • councils - Council terms (3 docs)                              │
│  • council_positions - Member positions in council (30 docs)      │
│  • events - LACPA events (30 docs)                                │
│  • users - Authentication users (4 docs)                          │
│  • application_requirements - Membership requirements (15 docs)   │
└──────────────────────────────────────────────────────────────────┘
```

## Request Flow Examples

### Viewing Board Members Page

1. **User Navigation**: User clicks "Board of Directors" or navigates to `/discover/board-of-directors`
2. **Browser Request**: GET request to backend
3. **Handler**: `GetBoardMembersPage` checks if HTMX request or full page load
4. **Data Fetching**:
   - Calls `GetAllCouncils()` to retrieve all councils
   - For each council, calls `GetCouncilCompositionWithDetails(councilID)`
   - Repository queries `councils` collection
   - Repository queries `council_positions` by `council_id`
   - For each position, fetches member from `individual_members` by `member_id`
   - Categorizes members: President, VP, Treasurer, Secretary, Board Members
5. **Template Rendering**: Server renders `board.html` with council data
6. **Response**: HTML page with all board members organized by council
7. **Browser**: Displays board members with interactive cards (biography, share features)

### API Request for Council Composition

1. **API Call**: GET `/api/council/{id}/composition/details`
2. **Handler**: `GetCouncilCompositionWithDetails` validates council ID
3. **Repository**: 
   - Fetches council document by ID
   - Queries all positions for that council
   - Performs member lookup for each position
   - Builds composition object with categorized positions
4. **Response**: JSON with nested structure:
   ```json
   {
     "council_id": "...",
     "council_name": "Twenty Fifth Council",
     "president": { "position": {...}, "member": {...} },
     "vice_president": { "position": {...}, "member": {...} },
     "board_treasurer": { "position": {...}, "member": {...} },
     "board_secretary": { "position": {...}, "member": {...} },
     "board_members": [...]
   }
   ```

### User Login Flow

1. **User Action**: User fills login form at `/auth/login`
2. **HTMX**: Intercepts form submission, sends POST with credentials
3. **Handler**: `Login` validates email/password
4. **Auth Repository**: 
   - Queries `users` collection by email
   - Verifies password hash
   - Generates JWT token
5. **Response**: Returns token and user data
6. **Frontend**: Stores token, redirects to dashboard
7. **Subsequent Requests**: Include JWT in Authorization header

### Creating/Updating Member

1. **Admin Action**: Fills member form in admin dashboard
2. **Request**: POST/PUT to `/api/members` with member data
3. **Handler**: `CreateIndividualMember` or `UpdateIndividualMember`
4. **Validation**: Validates required fields (name, email, member type)
5. **Repository**: Inserts/updates document in `individual_members`
6. **Response**: Returns created/updated member with success message
7. **UI Update**: HTMX updates member list without page reload

## Technology Choices

### Backend Stack

**Go (Golang) 1.21+**
- Strong typing and compile-time safety
- Excellent concurrency support
- Fast compilation and execution
- Great standard library

**Go Fiber v2.52.9**
- Fast and lightweight web framework (built on fasthttp)
- Express-like API familiar to developers
- Excellent performance (10x faster than Express.js)
- Rich middleware ecosystem
- Built-in template rendering

**MongoDB Go Driver**
- Official MongoDB driver for Go
- Type-safe BSON marshaling/unmarshaling
- Aggregation pipeline support
- Connection pooling and retry logic
- Flexible schema for evolving requirements

**Architecture Patterns**
- **Repository Pattern**: Separates data access from business logic
- **Clean Architecture**: Dependencies point inward (handler → repository → database)
- **Dependency Injection**: Handlers receive repository interfaces
- **RESTful API**: Standard HTTP methods and status codes

### Frontend Stack

**HTML5 + Server-Side Rendering**
- Go Fiber HTML template engine
- Server renders complete pages
- SEO-friendly approach
- Progressive enhancement

**HTMX**
- Modern interactivity without heavy JavaScript frameworks
- Server-driven UI updates
- Partial page updates
- Reduces frontend complexity and bundle size
- Progressive enhancement philosophy

**Tailwind CSS**
- Utility-first approach speeds development
- Consistent design system
- Responsive by default
- Custom component styling
- Small production bundle

**Vanilla JavaScript**
- No framework overhead
- Direct DOM manipulation where needed
- Custom navigation and interactions
- Card state management (member cards)

### Supporting Technologies

**JWT (JSON Web Tokens)**
- Stateless authentication
- Secure token-based sessions
- Easy to scale horizontally

**Docker**
- Containerized MongoDB
- Consistent development environment
- Easy deployment

**Git**
- Version control
- Branch: `auth_system` for authentication features

## API Endpoints

### Public Pages (Server-Rendered)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/` | Landing page |
| GET    | `/discover/board-of-directors` | Board members page |
| GET    | `/discover/individual-members` | Individual members directory |
| GET    | `/discover/how-to-join` | Membership application info |
| GET    | `/discover/events` | Events listing |
| GET    | `/auth/login` | Login page |
| GET    | `/auth/register` | Registration page |

### Council API
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/councils` | Get all councils |
| GET    | `/api/councils/active` | Get active council |
| GET    | `/api/council/:id` | Get council by ID |
| POST   | `/api/councils` | Create new council |
| PUT    | `/api/councils/:id` | Update council |
| DELETE | `/api/councils/:id/deactivate` | Deactivate council |
| GET    | `/api/council/:id/composition` | Get council composition (positions only) |
| GET    | `/api/council/:id/composition/details` | Get composition with member details |
| POST   | `/api/councils/positions` | Assign member to position |
| PUT    | `/api/councils/positions/:positionId` | Update position |
| DELETE | `/api/councils/positions/:positionId` | Remove position |
| GET    | `/api/councils/positions/:positionId` | Get position by ID |
| GET    | `/api/council/:id/available-positions` | Get available position slots |

### Individual Members API
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/individual-members` | Get all individual members |
| GET    | `/api/individual-members/:id` | Get member by ID |
| POST   | `/api/individual-members` | Create member |
| PUT    | `/api/individual-members/:id` | Update member |
| DELETE | `/api/individual-members/:id` | Delete member |
| GET    | `/api/individual-members/:id/council-history` | Get member's council history |

### Firm Members API
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/firm-members` | Get all firms |
| GET    | `/api/firm-members/:id` | Get firm by ID |
| POST   | `/api/firm-members` | Create firm |
| PUT    | `/api/firm-members/:id` | Update firm |
| DELETE | `/api/firm-members/:id` | Delete firm |

### Events API
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/events` | Get all events |
| GET    | `/api/events/:id` | Get event by ID |
| POST   | `/api/events` | Create event |
| PUT    | `/api/events/:id` | Update event |
| DELETE | `/api/events/:id` | Delete event |

### Authentication API
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/auth/login` | User login |
| POST   | `/api/auth/register` | User registration |
| POST   | `/api/auth/logout` | User logout |
| GET    | `/api/auth/me` | Get current user |
| POST   | `/api/auth/refresh` | Refresh JWT token |

## Environment Configuration

```env
# Server Configuration
PORT=3000

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=LACPA

# JWT Configuration
JWT_SECRET=lacpa-secret-key-change-this-in-production

# Email Service Configuration (Future)
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_USER=your-email@gmail.com
# SMTP_PASSWORD=your-app-password
```

## Database Schema

### Collections Overview

**individual_members**
- Stores CPA individual members
- Fields: personal info, contact, address, biography, license, membership status
- References: `current_council_position_id` links to council_positions

**firm_members**
- Stores firm/company members
- Fields: firm info, contact, services, partners

**councils**
- Stores council terms (e.g., "Twenty Fifth Council")
- Fields: name, start_date, end_date, is_active, description

**council_positions**
- Links members to council roles
- Fields: member_id, council_id, position (President, VP, etc.), dates, is_active
- Position Types: President, Vice President, Board Treasurer, Board Secretary, Board Member (6 max)

**events**
- Stores LACPA events and activities
- Fields: title, description, date, location, type, registration

**users**
- Authentication and authorization
- Fields: email, password_hash, role, profile info

**application_requirements**
- Membership application requirements and documents

### Key Relationships

```
councils (1) ←→ (N) council_positions
                       ↓
                       ↓ (references member_id)
                       ↓
individual_members (1) ←→ (N) council_positions
                               (via current_council_position_id)
```

## Folder Structure

```
Lacpa/
├── Backend/
│   ├── config/
│   │   └── database.go              # MongoDB connection setup
│   ├── handler/
│   │   ├── main-page.go             # Landing page handler
│   │   ├── council_handler.go       # Council & positions handlers
│   │   ├── member_handler.go        # Individual member handlers
│   │   ├── firm_handler.go          # Firm member handlers
│   │   ├── event_handler.go         # Event handlers
│   │   └── auth_handler.go          # Authentication handlers
│   ├── models/
│   │   ├── council.go               # Council, Position models
│   │   ├── individual_member.go     # Individual member model
│   │   ├── firm_member.go           # Firm model
│   │   ├── event.go                 # Event model
│   │   ├── user.go                  # User authentication model
│   │   └── landing_page.go          # Landing page content model
│   ├── repository/
│   │   ├── repository.go            # Repository interface
│   │   ├── council_repository.go    # Council data access
│   │   ├── member_repository.go     # Member data access
│   │   ├── firm_repository.go       # Firm data access
│   │   ├── event_repository.go      # Event data access
│   │   └── auth_repository.go       # Auth data access
│   ├── routes/
│   │   ├── routes.go                # Main route setup
│   │   └── main_page.go             # Public page routes
│   ├── templates/
│   │   ├── LACPA/
│   │   │   ├── board_members/
│   │   │   │   ├── board.html       # Board members template
│   │   │   │   └── index.html       # Board members wrapper
│   │   │   └── main_LACPA/
│   │   │       └── index.html       # Landing page template
│   │   ├── Admin_Dashboard/
│   │   │   └── index.html           # Admin dashboard
│   │   └── views/
│   │       └── layouts/             # Shared layouts
│   ├── utils/
│   │   ├── response.go              # Standard API responses
│   │   ├── request.go               # Request helpers
│   │   ├── validation.go            # Input validation
│   │   └── config.go                # Config helpers
│   ├── scripts/
│   │   └── fix_council_member_ids.go # Data migration scripts
│   ├── main.go                      # Application entry point
│   ├── go.mod                       # Go module definition
│   ├── go.sum                       # Dependency checksums
│   ├── .env                         # Environment variables
│   └── docker-compose.yml           # MongoDB container setup
│
├── LACPA_Web/                       # Frontend static files
│   ├── src/
│   │   ├── index.html               # Main HTML entry
│   │   ├── index.css                # Custom CSS
│   │   ├── output.css               # Tailwind output
│   │   ├── components/
│   │   │   ├── header.html          # Site header with navigation
│   │   │   ├── footer.html          # Site footer
│   │   │   ├── localnav.html        # Local navigation
│   │   │   └── buttons/
│   │   │       ├── bordered_button.html
│   │   │       └── traveling_button.html
│   │   └── pages/
│   │       └── landing.html         # Landing page content
│   ├── js/
│   │   └── app.js                   # Navigation & interactions
│   ├── assets/                      # Images, icons, etc.
│   ├── package.json                 # Node dependencies
│   └── tailwind.config.js           # Tailwind configuration
│
├── ARCHITECTURE.md                  # This file
├── CONTRIBUTING.md                  # Contribution guidelines
├── QUICKSTART.md                    # Quick start guide
└── README.md                        # Project overview
```

## Design Principles

1. **Separation of Concerns**: Each layer has a specific responsibility
   - Handlers: HTTP request/response handling
   - Repository: Database operations
   - Models: Data structure definitions
   - Utils: Shared helper functions

2. **Dependency Injection**: Handlers receive repository interfaces for testability

3. **Clean Architecture**: Dependencies point inward (handler → repository → database)

4. **RESTful API**: Standard HTTP methods and status codes
   - GET for retrieval
   - POST for creation
   - PUT for updates
   - DELETE for removal

5. **Repository Pattern**: Abstracts data access layer
   - Single interface for all repositories
   - Easy to mock for testing
   - Can swap database implementations

6. **Server-Side Rendering**: HTML templates for SEO and performance
   - Go Fiber template engine
   - HTMX for dynamic updates
   - Progressive enhancement

7. **Type Safety**: Go's strong typing catches errors at compile time

8. **Error Handling**: Consistent error responses via utils package

9. **Data Validation**: Input validation in handlers and utils

10. **Security Best Practices**:
    - JWT for authentication
    - Password hashing
    - CORS configuration
    - Input sanitization

## Key Features

### Council & Board Management
- Multiple council terms (past and present)
- Position assignments with validation (max 1 President, 1 VP, etc.)
- Member history tracking
- Composition queries with member details

### Member Management
- Individual CPA members directory
- Firm/company members directory
- Rich member profiles (biography, contact, licenses)
- Member type categorization (Apprentices, Practicing, Non-Practicing, Retired)

### Authentication
- User registration and login
- JWT-based authentication
- Role-based access control (future)
- Secure password hashing

### Event Management
- Event creation and listing
- Event registration (future)
- Event categories and types

## Data Integrity

### Position Constraints
- Max 1 President per council
- Max 1 Vice President per council
- Max 1 Board Treasurer per council
- Max 1 Board Secretary per council
- Max 6 Board Members per council
- Enforced at repository level

### Member-Position Linking
- `council_positions.member_id` → `individual_members._id`
- `council_positions.council_id` → `councils._id`
- `individual_members.current_council_position_id` → `council_positions._id`

### Data Consistency
- Cascade updates when member is assigned to position
- Active council validation
- Duplicate position prevention

## Development Workflow

1. **Local Development**:
   ```bash
   # Start MongoDB
   cd Backend
   docker-compose up -d
   
   # Run backend
   go run main.go
   
   # Backend runs on http://localhost:3000
   ```

2. **Database Access**:
   - MongoDB Compass GUI: `mongodb://localhost:27017`
   - Database: `LACPA` (uppercase)
   - Use scripts in `Backend/scripts/` for data migrations

3. **Adding Features**:
   - Define model in `models/`
   - Create repository methods in `repository/`
   - Implement handlers in `handler/`
   - Register routes in `routes/`
   - Add templates in `templates/` if needed

## Performance Considerations

- **Database Indexing**: Create indexes on frequently queried fields
  - `councils.is_active`
  - `council_positions.council_id`
  - `council_positions.member_id`
  - `individual_members.email`

- **Query Optimization**: 
  - Avoid N+1 queries by batching member lookups
  - Use aggregation pipelines for complex queries
  - Limit result sets with pagination

- **Caching Strategy** (future):
  - Cache active council composition
  - Cache member directory
  - Redis for session storage

## Future Enhancements

- [ ] Admin dashboard with full CRUD operations
- [ ] Member portal for profile updates
- [ ] Event registration and attendance tracking
- [ ] CPE credit management
- [ ] Email notifications
- [ ] Document management for licenses and certificates
- [ ] Payment integration for membership dues
- [ ] Advanced search and filtering
- [ ] API rate limiting
- [ ] Comprehensive logging and monitoring

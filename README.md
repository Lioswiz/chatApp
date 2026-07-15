# Chat Platform

A Go-based, real-time chat application utilizing WebSocket communication, SQLite database persistence, and a custom template rendering architecture.

---

## 📂 Project Architecture

```text
chat-platform/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go               # Configuration management (Port, DB paths, Max Upload size)
├── database/
│   ├── database.go             # SQLite connection and migration runner
│   └── schema.sql              # Database schema definitions
├── handlers/
│   ├── auth_handlers.go        # User registration, login, and logout handlers
│   ├── chat_handlers.go        # Renders the main chat workspace page
│   ├── home.go                 # Basic landing redirect handler
│   └── upload.go               # Placeholder for file uploads
├── middleware/
│   ├── session.go              # Auth session validation middleware
│   └── logger.go               # Request logger middleware (currently unused)
├── models/
│   ├── user.go                 # User schema models
│   ├── message.go              # Chat message structure
│   ├── session.go              # Session token representations
│   └── uploads.go              # File metadata structures
├── repository/
│   ├── user_repository.go      # User database queries
│   ├── message_repository.go   # Message database queries (public & private)
│   ├── session_repository.go   # Session token DB management
│   └── upload_repository.go    # File upload persistence queries
├── service/
│   ├── auth_service.go         # Password encryption and registration logic
│   ├── chat_service.go         # Messaging validation and processing
│   ├── session_service.go      # Token generation and validation
│   ├── upload_service.go       # File system writer and upload validators
│   └── validator.go            # User field validators
├── websocket/
│   ├── client.go               # WebSocket connection read/write loops
│   ├── hub.go                  # Room connection manager & message multiplexer
│   ├── message.go              # WebSocket message structs (Chat, Typing, Presence, etc.)
│   └── websocket.go            # WebSocket connection upgrade and registry
├── templates/
│   ├── login.html              # Login screen template
│   ├── register.html           # User signup screen template
│   ├── chat.html               # Multi-user chat template
│   └── profile.html            # Profile detail template (empty)
├── static/
│   ├── css/
│   │   └── style.css           # Global layout stylesheet
│   └── js/
│       ├── chat.js             # User interaction script (WebSocket client)
│       └── websocket.js        # Helper wrapper for WebSocket actions
├── go.mod                      # Module requirements list
└── README.md                   # Project status and developer roadmap (this file)
```

---

## ⚡ Current Status

### ✅ What is Working
- **Server Compilation:** The backend is written in Go 1.25.0 and successfully compiles using `go build ./...` without compilation errors.
- **User Authentication:** 
  - User signup and login routes are functional.
  - Passwords are securely hashed using `bcrypt` inside [auth_service.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/service/auth_service.go).
  - Validation rules check fields correctly before saving.
- **Session Middleware:** HTTP-only cookies are successfully issued, storing secure random session tokens which are validated against the database on protected routes via [session.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/middleware/session.go).
- **SQLite Database Layer:**
  - Auto-initialization executes the database setup from [schema.sql](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/database/schema.sql) if no DB exists.
  - Repositories handle insertions and selections for users, sessions, and chat logs successfully.
- **WebSocket Infrastructure:**
  - Upgrade handler upgrades HTTP connection to WebSocket channel securely.
  - Connection hub successfully registers clients and handles broadcast (public) and private channels.

---

### ❌ What is NOT Working / Bug list

#### 1. Database Schema & Upload Model Mismatch (Critical Backend Bug)
* **The Bug:** [schema.sql](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/database/schema.sql) declares the `uploads` table with an `uploaded_by` column (referencing users) but no `message_id`. However, [uploads.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/models/uploads.go) contains `MessageID int`, and [upload_repository.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/repository/upload_repository.go) attempts to write to a column named `message_id` which does not exist:
  ```sql
  INSERT INTO uploads (message_id, file_name, file_path, file_size, mime_type) VALUES (?, ?, ?, ?, ?)
  ```
* **Impact:** Any file upload metadata save action will throw a SQLite runtime error and crash the write sequence.

#### 2. Client-Side JavaScript WebSocket Conflict (Critical Frontend Bug)
* **The Bug:** [chat.html](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/templates/chat.html) includes both `/static/js/websocket.js` and `/static/js/chat.js` scripts:
  * Both scripts attempt to instantiate a WebSocket connection (`new WebSocket(...)`) to `/ws` on load, causing double connections.
  * In the browser environment, [websocket.js](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/static/js/websocket.js) defines a global `let socket`. When [chat.js](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/static/js/chat.js) tries to define `const socket`, the browser throws a SyntaxError: `Identifier 'socket' has already been declared`. This crashes script execution entirely.
  * Furthermore, [chat.js](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/static/js/chat.js) sends chat messages using an unwrapped payload without the `{ type: "chat", data: ... }` wrapping expected by the backend parser in [client.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/websocket/client.go). The backend consequently drops the message.
  * In `chat.js`, the `appendMessage(message)` method parses the root JSON object instead of the inner `.data` field. It accesses `message.message` which returns `undefined` (it should access `message.data.message`).

#### 3. Upload Route Integration Missing
* **The Bug:** The HTTP handler [upload.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/handlers/upload.go) is a placeholder that returns a `501 Not Implemented` error. The route `/upload` is also not registered in [routes.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/routes/routes.go), making files un-uploadable.

#### 4. Hardcoded WebSocket Username
* **The Bug:** In [websocket.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/websocket/websocket.go), every incoming WebSocket client connection assigns the hardcoded username `"User"` (`Username: "User"`). The actual username from the authenticated session context is not loaded or passed down.

#### 5. Empty Templates and Unused Components
* **Profile Template:** [profile.html](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/templates/profile.html) is completely blank, and there is no handler or router mapping in the backend to serve it.
* **Logging Middleware:** [logger.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/middleware/logger.go) is defined but not wrapped around any router groups in [routes.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/routes/routes.go).

#### 6. Zero Test Coverage
* **The Bug:** Run results of `go test ./...` show `[no test files]` for every single package. There is no automated validation for services, repositories, or handler endpoints.

---

## 🗺️ Next Steps & Roadmap

### Phase 1: Core System & Communication Fixes (High Priority)
1. **Fix client-side JS WebSocket overlap:**
   - Consolidate WebSocket connection logic. Remove one of the duplicate `new WebSocket()` initialization scopes.
   - Standardize the message format so the frontend wraps message payloads under `{ type: "chat", data: { ... } }` and the `onmessage` parser unpacks the sub-structure correctly.
   - Remove global naming collisions (e.g. declare one consolidated WebSocket instance variable).
2. **Resolve database schema constraints for Uploads:**
   - Decide on schema logic: Add `message_id` into database `uploads` table, or alter `models.Upload` and `repository.UploadRepository` to map `UploadedBy` (user reference) correctly.
3. **Pass Username to WebSocket clients:**
   - Update [websocket.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/websocket/websocket.go) to fetch the user details using `sessionService` or the authenticated context, avoiding the hardcoded `"User"` placeholder.

### Phase 2: Feature Implementation (Medium Priority)
1. **Implement HTTP Upload Handler:**
   - Link [upload_service.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/service/upload_service.go) to a registered `/upload` POST endpoint in [routes.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/routes/routes.go).
   - Wire up front-end upload click listener to upload attachments via AJAX/fetch, retrieve the `upload_id`, and send it with the WebSocket message payload.
2. **Enable Profile Page:**
   - Design and build out [profile.html](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/templates/profile.html).
   - Implement route `/profile` (protected by session middleware) displaying user details and avatar upload capability.
3. **Build Client Presence & Typing Updates:**
   - Wire up client-side handlers in [chat.js](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/static/js/chat.js) to display online users in the sidebar dynamically by handling incoming `presence` WebSocket payloads.
   - Implement typing indicator events (`typing` message type) to render "X is typing..." dynamically in the chat window.

### Phase 3: Operations & Quality Assurance (Low Priority)
1. **Add Middleware Logging:**
   - Register [logger.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/middleware/logger.go) in [routes.go](file:///C:/Users/Abz/Desktop/Projects/rep-chat/chatApp/routes/routes.go) to print HTTP method, path, and duration metrics.
2. **Write Unit/Integration Tests:**
   - Write tests in `/service` and `/repository` to test user validation, database operations, and session lifetimes using a mock DB or test-specific SQLite instances.
3. **Refine UI styling:**
   - Review `/static/css/style.css` styles to introduce glassmorphic elements, modern gradients, and smooth transition effects to achieve a high-end feel.
# chat-platform

chat-platform/

├── cmd/
│   └── main.go
│
├── config/
│
├── handlers/
│   ├── home.go
│   ├── auth.go
│   ├── chat.go
│   ├── upload.go
│   └── websocket.go
│
├── models/
│   ├── user.go
│   ├── message.go
│   └── file.go
│
├── repository/
│   ├── user_repository.go
│   ├── message_repository.go
│   └── upload_repository.go
│
├── service/
│   ├── auth_service.go
│   ├── chat_service.go
│   ├── upload_service.go
│   └── websocket_service.go
│
├── websocket/
│   ├── hub.go
│   ├── client.go
│   └── websocket.go
│
├── middleware/
│   ├── auth.go
│   └── logger.go
│
├── database/
│   ├── database.go
│   └── schema.sql
│
├── routes/
│   └── routes.go
│
├── templates/
│   ├── login.html
│   ├── register.html
│   ├── chat.html
│   └── profile.html
│
├── static/
│   ├── css/
│   ├── js/
│   ├── images/
│   └── uploads/
│
├── utils/
│
├── go.mod
└── README.md
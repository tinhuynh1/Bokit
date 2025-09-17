# 🧠 Real-Time Quiz System

A real-time quiz application built with Go, WebSockets, and MongoDB that allows multiple users to participate in live quiz sessions with real-time scoring and leaderboards.

## ✨ Features

- **Real-time Quiz Sessions**: Create and join quiz sessions with unique codes
- **Live Scoring**: Instant score updates as participants answer questions
- **Real-time Leaderboard**: Live ranking updates for all participants
- **WebSocket Communication**: Low-latency real-time updates
- **Concurrent Safety**: Thread-safe operations with mutex locks
- **Session Management**: Track session states (waiting, active, finished)

## 🚀 Quick Start

### 1. Start MongoDB
```bash
docker compose up -d
```

### 2. Run the Application
```bash
go run cmd/main.go
```

### 3. Test with Web Client
Open `test_client.html` in your browser to test the WebSocket functionality.

## 📡 API Endpoints

### Quiz Management
- `GET /api/v1/quizzes` - List all quizzes
- `GET /api/v1/quizzes/:id` - Get quiz by ID
- `POST /api/v1/quizzes` - Create new quiz
- `PUT /api/v1/quizzes/:id` - Update quiz

### Session Management
- `POST /api/v1/sessions` - Create new session
- `GET /api/v1/sessions/:code` - Get session info
- `POST /api/v1/sessions/:code/start` - Start session
- `GET /api/v1/sessions/:code/leaderboard` - Get leaderboard

### WebSocket
- `GET /ws` - WebSocket endpoint for real-time communication

## 🔌 WebSocket Protocol

The system uses a custom packet-based protocol:

### Packet Types
- `1` - Join Session
- `2` - Leave Session  
- `3` - Answer Question
- `4` - Start Quiz
- `5` - Next Question
- `6` - Quiz Ended
- `7` - Leaderboard Update
- `8` - Session Update
- `9` - Error

### Example Usage

#### Join Session
```javascript
const packet = {
    session_code: "123456",
    participant_name: "John Doe"
};
ws.send(String.fromCharCode(1) + JSON.stringify(packet));
```

#### Answer Question
```javascript
const packet = {
    question_index: 0,
    answer_choice: "choice_id_1"
};
ws.send(String.fromCharCode(3) + JSON.stringify(packet));
```

## 🏗️ Architecture

```
internal/
├── domain/           # Domain models and packet types
├── service/          # Business logic and session management
├── handler/          # HTTP and WebSocket handlers
├── repository/       # Data access layer
├── router/           # Route configuration
└── bootstrap/        # Application initialization
```

## 🔧 Configuration

The application uses configuration files in the `config/` directory:
- `config.yaml` - Main configuration
- `develop.yaml` - Development settings
- `prod.yaml` - Production settings

## 🧪 Testing

### Using the Test Client

1. Open `test_client.html` in your browser
2. Click "Connect" to establish WebSocket connection
3. Enter a session code and your name
4. Click "Join Session" to participate
5. Answer questions and see real-time updates

### Manual Testing with curl

#### Create a Quiz
```bash
curl -X POST http://localhost:8080/api/v1/quizzes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sample Quiz",
    "questions": [
      {
        "id": "q1",
        "name": "What is 2+2?",
        "time": 30,
        "choices": [
          {"id": "a1", "name": "3", "correct": false},
          {"id": "a2", "name": "4", "correct": true},
          {"id": "a3", "name": "5", "correct": false}
        ]
      }
    ]
  }'
```

#### Create a Session
```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{"quiz_id": "QUIZ_ID_HERE"}'
```

## 🎯 Key Features Implemented

✅ **User Participation**: Users can join sessions using unique codes  
✅ **Real-time Score Updates**: Scores update instantly as users answer  
✅ **Real-time Leaderboard**: Live ranking updates for all participants  
✅ **Concurrency Handling**: Thread-safe operations with mutex locks  
✅ **Low Latency**: WebSocket-based real-time communication  
✅ **Scalable Design**: Clean architecture supporting multiple users  

## 🔒 Security Notes

- WebSocket connections allow all origins (configure `CheckOrigin` for production)
- No authentication implemented (add JWT middleware for production)
- Session data stored in memory (add persistence for production)

## 🚀 Production Considerations

1. **Authentication**: Add JWT-based authentication
2. **Persistence**: Store sessions in Redis or database
3. **Rate Limiting**: Implement rate limiting for WebSocket connections
4. **Monitoring**: Add metrics and health checks
5. **Scaling**: Consider horizontal scaling with message queues
6. **Security**: Implement proper CORS and origin checking

## 📝 Example Quiz Data

```json
{
  "name": "General Knowledge Quiz",
  "questions": [
    {
      "id": "q1",
      "name": "What is the capital of France?",
      "time": 30,
      "choices": [
        {"id": "a1", "name": "London", "correct": false},
        {"id": "a2", "name": "Paris", "correct": true},
        {"id": "a3", "name": "Berlin", "correct": false}
      ]
    }
  ]
}
```

This system provides a solid foundation for real-time quiz applications with room for extension and customization based on specific requirements.

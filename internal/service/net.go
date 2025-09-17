package service

import (
	"context"
	"encoding/json"
	"fmt"
	"quiz-svc/internal/domain"

	"github.com/gorilla/websocket"
)

type NetService struct {
	sessionService    *SessionService
	connectionManager *ConnectionManager
}

func NewConnectionService(sessionService *SessionService) *NetService {
	return &NetService{
		sessionService:    sessionService,
		connectionManager: NewConnectionManager(),
	}
}

func (s *NetService) OnIncomingMessage(con *websocket.Conn, mt int, msg []byte) {
	if len(msg) < 2 {
		return
	}
	packetId := msg[0]
	packet := msg[1:]

	fmt.Printf("Received packet ID: %d, data: %s\n", packetId, string(packet))

	switch packetId {
	case domain.PacketJoinSession:
		s.handleJoinSession(con, packet)
	case domain.PacketAnswerQuestion:
		s.handleAnswerQuestion(con, packet)
	case domain.PacketStartQuiz:
		s.handleStartQuiz(con, packet)
	case domain.PacketNextQuestion:
		s.handleNextQuestion(con, packet)
	default:
		fmt.Printf("Unknown packet type: %d\n", packetId)
	}
}

func (s *NetService) OnDisconnect(con *websocket.Conn) {
	fmt.Println("OnDisconnect", con)

	// Get connection info and remove participant from session
	if info, exists := s.connectionManager.GetConnectionInfo(con); exists {
		s.sessionService.LeaveSession(context.Background(), info.SessionCode, info.ParticipantID)
		s.connectionManager.RemoveConnection(con)
	}
}

func (s *NetService) handleJoinSession(con *websocket.Conn, data []byte) {
	var req domain.JoinSessionRequest
	if err := json.Unmarshal(data, &req); err != nil {
		s.sendError(con, "Invalid join session request")
		return
	}

	participant, err := s.sessionService.JoinSession(context.Background(), req.SessionCode, req.ParticipantName, con)
	if err != nil {
		s.sendError(con, err.Error())
		return
	}

	// Register connection in manager
	s.connectionManager.RegisterConnection(con, req.SessionCode, participant.ID, req.ParticipantName)

	session, _ := s.sessionService.GetSession(context.Background(), req.SessionCode)

	response := domain.JoinSessionResponse{
		Success: true,
		Message: "Successfully joined session",
		Session: session,
	}

	s.sendPacket(con, domain.PacketJoinSession, response)
}

func (s *NetService) handleAnswerQuestion(con *websocket.Conn, data []byte) {
	var req domain.AnswerQuestionRequest
	if err := json.Unmarshal(data, &req); err != nil {
		s.sendError(con, "Invalid answer request")
		return
	}

	// Get connection info
	info, exists := s.connectionManager.GetConnectionInfo(con)
	if !exists {
		s.sendError(con, "Connection not registered")
		return
	}

	// Submit answer
	response, err := s.sessionService.SubmitAnswer(context.Background(), info.SessionCode, info.ParticipantID, req.QuestionIndex, req.AnswerChoice)
	if err != nil {
		s.sendError(con, err.Error())
		return
	}

	s.sendPacket(con, domain.PacketAnswerQuestion, response)

	// Broadcast leaderboard update to all participants in session
	s.broadcastLeaderboardUpdate(info.SessionCode)
}

func (s *NetService) handleStartQuiz(con *websocket.Conn, data []byte) {
	// TODO: Implement start quiz logic
	fmt.Println("Start quiz requested")
}

func (s *NetService) handleNextQuestion(con *websocket.Conn, data []byte) {
	// TODO: Implement next question logic
	fmt.Println("Next question requested")
}

func (s *NetService) sendPacket(con *websocket.Conn, packetType uint8, data interface{}) {
	bytes, err := s.PacketToBytes(data)
	if err != nil {
		fmt.Printf("Error converting packet to bytes: %v\n", err)
		return
	}

	// Prepend packet type
	packet := append([]byte{packetType}, bytes...)

	if err := con.WriteMessage(websocket.TextMessage, packet); err != nil {
		fmt.Printf("Error sending packet: %v\n", err)
	}
}

func (s *NetService) sendError(con *websocket.Conn, message string) {
	errorResp := domain.ErrorResponse{
		Error:   "error",
		Message: message,
	}
	s.sendPacket(con, domain.PacketError, errorResp)
}

func (c *NetService) PacketToBytes(packet interface{}) ([]byte, error) {
	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *NetService) packetToPacketId(packet interface{}) (uint8, error) {
	switch packet.(type) {
	case *domain.JoinSessionResponse:
		return domain.PacketJoinSession, nil
	case *domain.AnswerQuestionResponse:
		return domain.PacketAnswerQuestion, nil
	case *domain.StartQuizResponse:
		return domain.PacketStartQuiz, nil
	case *domain.NextQuestionResponse:
		return domain.PacketNextQuestion, nil
	case *domain.QuizEndedResponse:
		return domain.PacketQuizEnded, nil
	case *domain.LeaderboardUpdate:
		return domain.PacketLeaderboard, nil
	case *domain.SessionUpdate:
		return domain.PacketSessionUpdate, nil
	case *domain.ErrorResponse:
		return domain.PacketError, nil
	default:
		return 0, fmt.Errorf("unknown packet type")
	}
}

func (c *NetService) packetIdToPacket(packetId uint8) interface{} {
	switch packetId {
	case domain.PacketJoinSession:
		return &domain.JoinSessionRequest{}
	case domain.PacketAnswerQuestion:
		return &domain.AnswerQuestionRequest{}
	case domain.PacketStartQuiz:
		return &domain.StartQuizResponse{}
	case domain.PacketNextQuestion:
		return &domain.NextQuestionResponse{}
	case domain.PacketQuizEnded:
		return &domain.QuizEndedResponse{}
	case domain.PacketLeaderboard:
		return &domain.LeaderboardUpdate{}
	case domain.PacketSessionUpdate:
		return &domain.SessionUpdate{}
	case domain.PacketError:
		return &domain.ErrorResponse{}
	default:
		return nil
	}
}

func (s *NetService) broadcastLeaderboardUpdate(sessionCode string) {
	leaderboard, err := s.sessionService.GetLeaderboard(context.Background(), sessionCode)
	if err != nil {
		fmt.Printf("Error getting leaderboard: %v\n", err)
		return
	}

	update := domain.LeaderboardUpdate{
		Leaderboard: leaderboard,
	}

	// Send to all connections in the session
	connections := s.connectionManager.GetConnectionsBySession(sessionCode)
	for _, conn := range connections {
		s.sendPacket(conn, domain.PacketLeaderboard, update)
	}
}

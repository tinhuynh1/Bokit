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
	netService := &NetService{
		sessionService:    sessionService,
		connectionManager: NewConnectionManager(),
	}

	// Set broadcast callback
	//sessionService.SetBroadcastCallback(netService.broadcastToSession)

	return netService
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
	// case domain.PacketAnswerQuestion:
	// 	s.handleAnswerQuestion(con, packet)
	// case domain.PacketStartQuiz:
	// 	s.handleStartQuiz(con)
	// case domain.PacketNextQuestion:
	// 	s.handleNextQuestion(con)
	// case domain.PacketLeaveSession:
	// 	s.handleLeaveSession(con)
	default:
		fmt.Printf("Unknown packet type: %d\n", packetId)
	}
}

func (s *NetService) OnDisconnect(con *websocket.Conn) {
	fmt.Println("OnDisconnect", con)

	// Get connection info and remove participant from session
	// if info, exists := s.connectionManager.GetConnectionInfo(con); exists {
	// 	s.sessionService.LeaveSession(context.Background(), info.SessionCode, info.ParticipantID)
	// 	s.connectionManager.RemoveConnection(con)

	// 	// Broadcast notification after leave
	// 	if sess, err := s.sessionService.GetSession(context.Background(), info.SessionCode); err == nil {
	// 		notify := domain.Notification{
	// 			Event:           "leave",
	// 			ParticipantName: info.ParticipantName,
	// 			Participants:    len(sess.Participants),
	// 		}
	// 		s.broadcastToSession(info.SessionCode, domain.PacketNotification, notify)
	// 	}
	// }
}

func (s *NetService) handleJoinSession(con *websocket.Conn, data []byte) {
	var req domain.JoinSessionRequest
	if err := json.Unmarshal(data, &req); err != nil {
		//s.sendError(con, "Invalid join session request")
		return
	}

	_, err := s.sessionService.JoinSession(context.Background(), req.SessionCode, req.ParticipantName, con)
	if err != nil {
		//s.sendError(con, err.Error())
		return
	}

}

// func (s *NetService) handleAnswerQuestion(con *websocket.Conn, data []byte) {
// 	var req domain.AnswerQuestionRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		s.sendError(con, "Invalid answer request")
// 		return
// 	}

// 	// Get connection info
// 	info, exists := s.connectionManager.GetConnectionInfo(con)
// 	if !exists {
// 		s.sendError(con, "Connection not registered")
// 		return
// 	}

// 	// Submit answer
// 	response, err := s.sessionService.SubmitAnswer(context.Background(), info.SessionCode, info.ParticipantID, req.QuestionIndex, req.AnswerChoice)
// 	if err != nil {
// 		s.sendError(con, err.Error())
// 		return
// 	}

// 	s.sendPacket(con, domain.PacketAnswerQuestion, response)

// 	// Broadcast leaderboard update to all participants in session
// 	s.broadcastLeaderboardUpdate(info.SessionCode)
// }

// func (s *NetService) handleStartQuiz(con *websocket.Conn) {
// 	// Get connection info
// 	info, exists := s.connectionManager.GetConnectionInfo(con)
// 	if !exists {
// 		s.sendError(con, "Connection not registered")
// 		return
// 	}

// 	// Start the session
// 	err := s.sessionService.StartSession(context.Background(), info.SessionCode)
// 	if err != nil {
// 		s.sendError(con, err.Error())
// 		return
// 	}

// 	// Broadcast session started to all participants
// 	// s.broadcastSessionUpdate(info.SessionCode) // This line is removed as per the edit hint

// 	// Start the first question
// 	s.startFirstQuestion(info.SessionCode)
// }

// func (s *NetService) startFirstQuestion(sessionCode string) {
// 	// Get session and start first question
// 	session, err := s.sessionService.GetSession(context.Background(), sessionCode)
// 	if err != nil {
// 		return
// 	}

// 	if len(session.Quiz.Questions) == 0 {
// 		return
// 	}

// 	// Broadcast first question
// 	question := session.Quiz.Questions[0]
// 	response := domain.NextQuestionResponse{
// 		QuestionIndex: 0,
// 		Question:      &question,
// 		TimeLeft:      question.Time,
// 	}

// 	s.broadcastToSession(sessionCode, domain.PacketNextQuestion, response)

// 	// Start timer for first question
// 	s.sessionService.StartQuestionTimer(context.Background(), sessionCode)
// }

// func (s *NetService) handleNextQuestion(con *websocket.Conn) {
// 	// Get connection info
// 	info, exists := s.connectionManager.GetConnectionInfo(con)
// 	if !exists {
// 		s.sendError(con, "Connection not registered")
// 		return
// 	}

// 	// Manually advance to next question (for host control)
// 	s.advanceToNextQuestion(info.SessionCode)
// }

// func (s *NetService) advanceToNextQuestion(sessionCode string) {
// 	session, err := s.sessionService.GetSession(context.Background(), sessionCode)
// 	if err != nil {
// 		return
// 	}

// 	if session.Status != domain.SessionStatusActive {
// 		return
// 	}

// 	// Advance to next question
// 	session.CurrentQuestion++

// 	if session.CurrentQuestion >= len(session.Quiz.Questions) {
// 		// Quiz ended
// 		s.endQuiz(sessionCode)
// 		return
// 	}

// 	// Broadcast next question
// 	question := session.Quiz.Questions[session.CurrentQuestion]
// 	response := domain.NextQuestionResponse{
// 		QuestionIndex: session.CurrentQuestion,
// 		Question:      &question,
// 		TimeLeft:      question.Time,
// 	}

// 	s.broadcastToSession(sessionCode, domain.PacketNextQuestion, response)

// 	// Start timer for this question
// 	s.sessionService.StartQuestionTimer(context.Background(), sessionCode)
// }

// func (s *NetService) endQuiz(sessionCode string) {
// 	// Get final leaderboard
// 	leaderboard, _ := s.sessionService.GetLeaderboard(context.Background(), sessionCode)

// 	// Broadcast quiz ended
// 	response := domain.QuizEndedResponse{
// 		FinalLeaderboard: leaderboard,
// 	}

// 	s.broadcastToSession(sessionCode, domain.PacketQuizEnded, response)
// }

// // func (s *NetService) handleLeaveSession(con *websocket.Conn) {
// // 	if info, ok := s.connectionManager.GetConnectionInfo(con); ok {
// // 		// cập nhật session
// // 		s.sessionService.LeaveSession(context.Background(), info.SessionCode, info.ParticipantID)
// // 		s.connectionManager.RemoveConnection(con)

// // 		// broadcast notification leave (nếu bạn đã dùng Notification)
// // 		if sess, err := s.sessionService.GetSession(context.Background(), info.SessionCode); err == nil {
// // 			notify := domain.Notification{
// // 				Event:           "leave",
// // 				ParticipantName: info.ParticipantName,
// // 				Participants:    len(sess.Participants),
// // 			}
// // 			s.broadcastToSession(info.SessionCode, domain.PacketNotification, notify)
// // 		}
// // 	}
// // 	// tùy chọn: đóng kết nối
// // 	// _ = con.Close()
// // }

// func (s *NetService) broadcastToSession(sessionCode string, packetType uint8, data interface{}) {
// 	connections := s.connectionManager.GetConnectionsBySession(sessionCode)
// 	for _, conn := range connections {
// 		s.sendPacket(conn, packetType, data)
// 	}
// }

// func (s *NetService) sendPacket(con *websocket.Conn, packetType uint8, data interface{}) {
// 	bytes, err := s.PacketToBytes(data)
// 	if err != nil {
// 		fmt.Printf("Error converting packet to bytes: %v\n", err)
// 		return
// 	}

// 	// Prepend packet type
// 	packet := append([]byte{packetType}, bytes...)

// 	if err := con.WriteMessage(websocket.TextMessage, packet); err != nil {
// 		fmt.Printf("Error sending packet: %v\n", err)
// 	}
// }

// func (s *NetService) sendError(con *websocket.Conn, message string) {
// 	errorResp := domain.ErrorResponse{
// 		Error:   "error",
// 		Message: message,
// 	}
// 	s.sendPacket(con, domain.PacketError, errorResp)
// }

// func (c *NetService) PacketToBytes(packet interface{}) ([]byte, error) {
// 	bytes, err := json.Marshal(packet)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bytes, nil
// }

// func (s *NetService) broadcastLeaderboardUpdate(sessionCode string) {
// 	leaderboard, err := s.sessionService.GetLeaderboard(context.Background(), sessionCode)
// 	if err != nil {
// 		fmt.Printf("Error getting leaderboard: %v\n", err)
// 		return
// 	}

// 	update := domain.LeaderboardUpdate{
// 		Leaderboard: leaderboard,
// 	}

// 	// Send to all connections in the session
// 	connections := s.connectionManager.GetConnectionsBySession(sessionCode)
// 	for _, conn := range connections {
// 		s.sendPacket(conn, domain.PacketLeaderboard, update)
// 	}
// }

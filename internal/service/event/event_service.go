package event

import (
	"booking-svc/internal/common/response"
	"booking-svc/internal/domain"
	"booking-svc/internal/dto"
	"booking-svc/internal/infra/database"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

type EventService struct {
	logger      *zap.Logger
	eventRepo   domain.EventRepository
	bookingRepo domain.TicketBookingRepository
}

func NewEventService(logger *zap.Logger,
	repo domain.EventRepository,
	bookingRepo domain.TicketBookingRepository,
) *EventService {
	return &EventService{
		logger:      logger,
		eventRepo:   repo,
		bookingRepo: bookingRepo}
}

func (s *EventService) ListEvents(ctx context.Context, page int, pageSize int, from string, to string) ([]dto.ListEventResponse, response.Paging, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	events, total, err := s.eventRepo.ListEvents(ctx, pageSize, offset, from, to)
	if err != nil {
		return nil, response.Paging{}, err
	}
	resp := make([]dto.ListEventResponse, len(events))
	for i, event := range events {
		resp[i] = dto.ListEventResponse{
			ID:               event.ID,
			Name:             event.Name,
			Description:      event.Description,
			DateTime:         event.DateTime.Format(time.RFC3339),
			TicketPrice:      event.TicketPrice,
			AvailableTickets: event.AvailableTickets,
			SoldTickets:      event.SoldTickets,
		}
	}
	return resp, response.Paging{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(total),
		TotalPages: int(total) / pageSize,
	}, nil
}

func (s *EventService) CreateEvent(ctx context.Context, req dto.CreateEventRequest) error {
	dateTime, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		return err
	}
	req.DateTime = dateTime.Format(time.RFC3339)
	if dateTime.Before(time.Now()) {
		return errors.New("event date time must be in the future")
	}
	if req.AvailableTickets < 0 {
		return errors.New("available tickets must be greater than 0")
	}
	if req.TicketPrice < 0 {
		return errors.New("ticket price must be greater than 0")
	}
	event := &domain.Event{
		Name:             req.Name,
		Description:      req.Description,
		DateTime:         dateTime,
		TicketPrice:      req.TicketPrice,
		AvailableTickets: req.AvailableTickets,
	}
	err = s.eventRepo.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error("create event failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *EventService) DeleteEvent(ctx context.Context, id int) error {
	err := s.eventRepo.DeleteEvent(ctx, id)
	if err != nil {
		s.logger.Error("delete event failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *EventService) BookTicket(ctx context.Context, req dto.BookingEventRequest) error {
	// check if event exists
	tx := database.BeginTxn()
	defer database.RollbackTxn(tx)

	event, err := s.eventRepo.GetEventForBooking(ctx, tx, req.EventID)
	if err != nil {
		s.logger.Error("get event for booking failed", zap.Error(err))
		return err
	}

	// check if event is in the past
	if event.DateTime.Before(time.Now()) {
		s.logger.Error("event is in the past")
		return errors.New("event is in the past")
	}

	// check if event has available tickets
	if event.AvailableTickets < req.Quantity {
		s.logger.Error("no available tickets")
		return errors.New("no available tickets")
	}

	booking := &domain.TicketBooking{
		EventID:    req.EventID,
		Email:      req.Email,
		Status:     domain.BookingStatusPending,
		Quantity:   req.Quantity,
		TotalPrice: event.TicketPrice * float64(req.Quantity),
	}
	err = s.bookingRepo.CreateWithTx(ctx, tx, booking)
	if err != nil {
		s.logger.Error("create booking failed", zap.Error(err))
		return err
	}
	event.AvailableTickets -= req.Quantity
	event.SoldTickets += req.Quantity
	err = s.eventRepo.UpdateEventWithTx(ctx, tx, event)
	if err != nil {
		s.logger.Error("update event failed", zap.Error(err))
		return err
	}
	return database.CommitTxn(tx)
}

func (s *EventService) CancelBooking() error {
	//acqire lock
	ctx, cancel := context.WithTimeout(context.Background(), 450*time.Second)
	defer cancel()

	//acquire lock prevent multiple job cancel booking run at the same time
	lock := s.bookingRepo.AcquireBookingLock(ctx)
	if lock != nil {
		return errors.New("lock acquired")
	}
	defer s.bookingRepo.ReleaseBookingLock(ctx)

	tx := database.BeginTxn()
	defer database.RollbackTxn(tx)
	bookings, err := s.bookingRepo.GetExpiredBooking(ctx, tx)
	if err != nil {
		s.logger.Error("cancel booking failed", zap.Error(err))
		return err
	}

	if len(bookings) == 0 {
		s.logger.Info("no expired booking")
		return nil
	}

	var eventIds []int
	var bookingIds []int
	mapIdBooking := make(map[int]int) //map event id to quantity
	for _, booking := range bookings {
		bookingIds = append(bookingIds, booking.ID)
		eventIds = append(eventIds, booking.EventID)
		mapIdBooking[booking.EventID] += booking.Quantity
	}
	err = s.bookingRepo.UpdateStatusByIds(ctx, tx, bookingIds, string(domain.BookingStatusCancelled))
	if err != nil {
		s.logger.Error("cancel booking failed", zap.Error(err))
		return err
	}

	events, err := s.eventRepo.GetEventsByIds(ctx, eventIds)
	if err != nil {
		s.logger.Error("get events by ids failed", zap.Error(err))
		return err
	}

	for i, event := range events {
		events[i].AvailableTickets += mapIdBooking[event.ID]
		events[i].SoldTickets -= mapIdBooking[event.ID]
	}

	err = s.eventRepo.UpdateEventsWithTx(ctx, tx, events)
	if err != nil {
		s.logger.Error("update events failed", zap.Error(err))
		return err
	}
	err = database.CommitTxn(tx)
	if err != nil {
		s.logger.Error("commit txn failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *EventService) UpdateEvent(ctx context.Context, req dto.UpdateEventRequest) error {
	tx := database.BeginTxn()
	defer database.RollbackTxn(tx)
	event, err := s.eventRepo.GetEventForBooking(ctx, tx, req.ID)
	if err != nil {
		s.logger.Error("get event for booking failed", zap.Error(err))
		return err
	}
	if req.DateTime != "" {
		dateTime, err := time.Parse(time.RFC3339, req.DateTime)
		if err != nil {
			s.logger.Error("get event for booking failed", zap.Error(err))
			return err
		}
		event.DateTime = dateTime
	}

	if req.AvailableTickets != nil {
		if *req.AvailableTickets < 0 {
			return errors.New("available tickets must be greater than 0")
		}
		event.AvailableTickets = *req.AvailableTickets
	}

	if req.TicketPrice != 0 {
		event.TicketPrice = req.TicketPrice
	}

	if req.Name != "" {
		event.Name = req.Name
	}
	if req.Description != "" {
		event.Description = req.Description
	}

	err = s.eventRepo.UpdateEventWithTx(ctx, tx, event)
	if err != nil {
		s.logger.Error("update event failed", zap.Error(err))
		return err
	}
	return database.CommitTxn(tx)
}

func (s *EventService) GetEventStats(ctx context.Context, from string, to string) (dto.EventStatsResponse, error) {
	var stats dto.EventStatsResponse
	events, _, err := s.eventRepo.ListEvents(ctx, 1000, 0, from, to)
	if err != nil {
		s.logger.Error("get event stats failed", zap.Error(err))
		return stats, err
	}
	stats.TotalEvents = len(events)
	eventIds := make([]int, len(events))
	for _, event := range events {
		eventIds = append(eventIds, event.ID)
	}

	bookings, err := s.bookingRepo.GetBookingsByEventIds(ctx, eventIds)
	if err != nil {
		s.logger.Error("get bookings by event ids failed", zap.Error(err))
		return stats, err
	}

	mapIdBooking := make(map[int]domain.TicketBooking)
	for _, booking := range bookings {
		mapIdBooking[booking.EventID] = booking
	}

	for _, event := range events {
		if booking, ok := mapIdBooking[event.ID]; ok {
			stats.TotalTicketsSold += booking.Quantity
			stats.TotalRevenue += booking.TotalPrice
		}
	}

	return stats, nil
}

func (s *EventService) PaymentCallback(ctx context.Context, req dto.PaymentCallbackRequest) error {
	//validate signature
	// check if booking exists
	booking, err := s.bookingRepo.GetBookingById(ctx, req.BookingID)
	if err != nil {
		s.logger.Error("get booking by id failed", zap.Error(err))
		return err
	}

	err = s.bookingRepo.UpdateStatusById(ctx, booking.ID, string(domain.BookingStatusConfirmed))
	if err != nil {
		s.logger.Error("update booking failed", zap.Error(err))
		return err
	}

	return nil
}

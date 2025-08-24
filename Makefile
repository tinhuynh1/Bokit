
migrate-down:
	goose -dir migration postgres "host=localhost port=5432 user=booking_app password=X4pV7_qM9%tN1wK6@rG8jM2Z dbname=booking sslmode=disable" down

migrate-up:
	goose -dir migration postgres "host=localhost port=5432 user=booking_app password=X4pV7_qM9%tN1wK6@rG8jM2Z dbname=booking sslmode=disable" up

run:
	go run cmd/main.go
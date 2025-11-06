package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeModel struct {
	DB *sql.DB
}

type Attende struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (m *AttendeModel) Insert(attende *Attende) (*Attende, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO attendes(event_id, user_id) VALUES ($1, $2) RETURNING id"

	err := m.DB.QueryRowContext(ctx, query, attende.EventId, attende.UserId).Scan(&attende.Id)

	if err != nil {
		return nil, err
	}
	return attende, nil
}

func (m *AttendeModel) GetByEventAndAttende(eventId, userId int) (*Attende, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendes WHERE event_id = $1 AND user_id = $2"
	var attende Attende
	err := m.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attende.Id, &attende.EventId, &attende.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attende, nil
}

func (m *AttendeModel) GetAttendesByEvent(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	SELECT u.id, u.name, u.email
	FROM users u 
	JOIN attendes a ON u.id == a.user_id
	WHERE a.event_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (m *AttendeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM attende WHERE user_id = $1 AND event_id = $2"
	_, err := m.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (m *AttendeModel) GetEventsByattende(attendeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
	FROM events e
	JOIN attendes a ON e.id = a.event_id
	WHERE e.user_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, attendeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

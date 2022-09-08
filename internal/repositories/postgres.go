package repositories

import (
	"context"
	"database/sql"

	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/models"
	"github.com/Dsmit05/party-day-bot/internal/repositories/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	errUserIsExist = errors.New("User with tg_id already exists.")
	errOther       = errors.New("Unknown error.")
)

type DBConnectI interface {
	GetConnectDB() string
}

type PostgresRepository struct {
	conn    *pgx.Conn
	queries *postgres.Queries
}

func NewPostgresRepository(url DBConnectI) (*PostgresRepository, error) {
	conn, err := pgx.Connect(context.Background(), url.GetConnectDB())
	if err != nil {
		return nil, err
	}

	queries := postgres.New(conn)

	logger.Info("repositories.NewPostgresRepository()", "Init")

	return &PostgresRepository{conn: conn, queries: queries}, nil
}

func (p *PostgresRepository) Close() {
	err := p.conn.Close(context.Background())
	if err != nil {
		logger.Error("(o *PostgresRepository) Close() error:", err)
	}
}

func (p *PostgresRepository) CreateUser(ctx context.Context, user models.User) error {
	inputData := postgres.CreateUserParams{
		TgID:      sql.NullInt64{user.TgID, true},
		ChatID:    sql.NullInt64{user.ChatID, true},
		Role:      user.Role,
		FirstName: sql.NullString{user.FirstName, true},
		LastName:  sql.NullString{user.LastName, true},
		UserName:  sql.NullString{user.UserName, true},
	}

	err := p.queries.CreateUser(ctx, inputData)
	val, ok := err.(*pgconn.PgError)

	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		return errUserIsExist
	}

	if err != nil {
		logger.DatabaseError("queries.CreateUser", err, inputData)
		return errOther
	}

	return nil
}

func (p *PostgresRepository) ReadUsers(ctx context.Context) (map[int64]models.User, error) {
	users, err := p.queries.ReadUsers(ctx)
	if err != nil {
		logger.DatabaseError("queries.ReadUsers", err, users)
		return nil, err
	}

	usersCast := make(map[int64]models.User)
	for _, user := range users {
		usersCast[user.TgID.Int64] = models.User{
			ChatID:    user.ChatID.Int64,
			Role:      user.Role,
			FirstName: user.FirstName.String,
			LastName:  user.LastName.String,
			UserName:  user.UserName.String,
		}
	}

	return usersCast, nil
}

func (p *PostgresRepository) UpdateRole(ctx context.Context, id int64, role string) error {
	err := p.queries.UpdateUserRole(ctx, postgres.UpdateUserRoleParams{
		Role: role,
		TgID: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		logger.DatabaseError("queries.UpdateRole", err, role)
		return errOther
	}

	return nil
}

func (p *PostgresRepository) CreateFile(ctx context.Context, file models.File) error {
	err := p.queries.CreateFile(ctx, postgres.CreateFileParams{
		UserTgID: sql.NullInt64{file.UserTgID, true},
		TgID:     sql.NullString{file.TgID, true},
		Url:      sql.NullString{file.URL, true},
	})

	if err != nil {
		logger.DatabaseError("queries.CreateFile", err, file)
		return errOther
	}

	return nil
}

func (p *PostgresRepository) CreateMessage(ctx context.Context, userTgID int64, text string) error {
	err := p.queries.CreateMessage(ctx, postgres.CreateMessageParams{
		UserTgID: sql.NullInt64{userTgID, true},
		Text:     sql.NullString{text, true},
	})

	if err != nil {
		logger.DatabaseError("queries.CreateMessage", err, userTgID)
		return errOther
	}

	return nil
}

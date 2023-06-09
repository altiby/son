package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrNothingToUpdate = errors.New("nothing to update")

// Config is a database connection configuration.
type Config struct {
	URL             string        `mapstructure:"url" valid:"required"`
	MaxConns        int           `mapstructure:"max_conns,required"`
	MinConns        int           `mapstructure:"min_conns,required"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time,required"`
	Username        string        `mapstructure:"username" json:"username"`
	Password        string        `mapstructure:"password" json:"password"`
	MigrationDir    string        `mapstructure:"migration_dir" json:"migration_dir"`
}

type Postgres struct {
	db sqlx.ExtContext
}

func New(ctx context.Context, cfg Config) (*Postgres, error) {
	const (
		pingTimeout = 5 * time.Second
	)

	databaseURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse database URL: %w", err)
	}

	if cfg.Username != "" {
		databaseURL.User = url.UserPassword(cfg.Username, cfg.Password)
	}

	db, err := sqlx.Open("postgres", databaseURL.String())
	if err != nil {
		return nil, fmt.Errorf("create OTEL DB pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MinConns)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	return &Postgres{db}, nil
}

func InitPostgresql(cfg Config) (*Postgres, error) {
	// postgres
	storage, err := New(context.TODO(), cfg)
	if err != nil {
		return nil, err
	}

	conn, err := storage.Connection()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationDir),
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return storage, nil
}

func (p Postgres) Connection() (*sqlx.DB, error) {
	conn, ok := p.db.(*sqlx.DB)
	if !ok {
		return nil, errors.New("connection is not *sqlx.DB, you might be already in transaction")
	}

	return conn, nil
}

func (p Postgres) Health(ctx context.Context) error {
	conn, err := p.Connection()
	if err != nil {
		return fmt.Errorf("could not get db connection: %w", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	return nil
}

func (p Postgres) DB() sqlx.ExtContext {
	return p.db
}

func (p Postgres) WithTx(tx *sqlx.Tx) Postgres {
	return Postgres{
		db: tx,
	}
}

type db[T any] interface {
	DB() sqlx.ExtContext
	WithTx(tx *sqlx.Tx) T
}

func inTx[T db[T]](ctx context.Context, p T, f func(ctx context.Context, p T) error) error {
	conn, ok := p.DB().(*sqlx.DB)
	if !ok {
		return errors.New("connection is not *sqlx.DB, you might be already in transaction")
	}

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin tx: %w", err)
	}

	if err := f(ctx, p.WithTx(tx)); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = fmt.Errorf("%s: %w", rollbackErr, err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func CheckAffected(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNothingToUpdate
	}

	return nil
}

package postgres

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ARM-auto-place/users/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectManager interface {
	GetMasterConn() *pgxpool.Pool
	GetReplicaConn() *pgxpool.Pool
}

type manager struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool
}

func (m *manager) GetMasterConn() *pgxpool.Pool  { return m.master }
func (m *manager) GetReplicaConn() *pgxpool.Pool { return m.replica }

func New(ctx context.Context, cfg config.PostgresConfig) (ConnectManager, error) {
	master, err := dbConnect(ctx, cfg, cfg.Addr, cfg.User, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("primary: %w", err)
	}

	addr, user, password := cfg.ReplicaAddr, cfg.UserRO, cfg.PasswordRO
	if addr == "" {
		addr = cfg.Addr
	}
	if user == "" {
		user, password = cfg.User, cfg.Password
	}
	replica, err := dbConnect(ctx, cfg, addr, user, password)
	if err != nil {
		master.Close()
		return nil, fmt.Errorf("replica: %w", err)
	}

	return &manager{master: master, replica: replica}, nil
}

func dbConnect(ctx context.Context, cfg config.PostgresConfig, addr, user, password string) (*pgxpool.Pool, error) {
	poolConfig, err := connectionConfig(cfg, addr, user, password)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func connectionConfig(cfg config.PostgresConfig, addr, user, password string) (*pgxpool.Config, error) {
	host, portText, err := net.SplitHostPort(addr)
	if err != nil || host == "" {
		return nil, fmt.Errorf("invalid PostgreSQL host:port %q", addr)
	}
	port, err := strconv.ParseUint(portText, 10, 16)
	if err != nil || port == 0 {
		return nil, fmt.Errorf("invalid PostgreSQL port %q", portText)
	}
	if cfg.MaxConn < 1 {
		return nil, fmt.Errorf("PostgreSQL pool size must be positive")
	}
	pool, err := pgxpool.ParseConfig("sslmode=disable")
	if err != nil {
		return nil, err
	}
	pool.ConnConfig.Host = host
	pool.ConnConfig.Port = uint16(port)
	pool.ConnConfig.Database = cfg.DB
	pool.ConnConfig.User = user
	pool.ConnConfig.Password = password
	pool.MaxConns = cfg.MaxConn
	pool.MaxConnIdleTime = cfg.MaxIdleLifetime
	pool.MaxConnLifetime = cfg.MaxLifetime
	return pool, nil
}

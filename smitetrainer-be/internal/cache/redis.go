package cache

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type RedisConfig struct {
	Addr        string
	Username    string
	Password    string
	DB          int
	DialTimeout time.Duration
	IOTimeout   time.Duration
	DefaultTTL  time.Duration
}

type RedisStore struct {
	cfg RedisConfig
}

func NewRedisStore(cfg RedisConfig) (*RedisStore, error) {
	if strings.TrimSpace(cfg.Addr) == "" {
		return nil, fmt.Errorf("redis addr is required")
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 3 * time.Second
	}
	if cfg.IOTimeout <= 0 {
		cfg.IOTimeout = 3 * time.Second
	}
	if cfg.DB < 0 {
		return nil, fmt.Errorf("redis db must be >= 0")
	}
	if cfg.DefaultTTL < 0 {
		cfg.DefaultTTL = 0
	}
	return &RedisStore{cfg: cfg}, nil
}

func (s *RedisStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	conn, reader, writer, err := s.openConn(ctx)
	if err != nil {
		return nil, false, err
	}
	defer conn.Close()

	if err := writeCommand(writer, []byte("GET"), []byte(key)); err != nil {
		return nil, false, err
	}
	prefix, data, err := readRESP(reader)
	if err != nil {
		return nil, false, err
	}
	if prefix != '$' {
		return nil, false, fmt.Errorf("unexpected redis response type for GET: %q", string(prefix))
	}
	if data == nil {
		return nil, false, nil
	}
	return data, true, nil
}

func (s *RedisStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	conn, reader, writer, err := s.openConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if ttl <= 0 {
		ttl = s.cfg.DefaultTTL
	}

	if ttl > 0 {
		ttlMs := ttl.Milliseconds()
		if ttlMs < 1 {
			ttlMs = 1
		}
		if err := writeCommand(
			writer,
			[]byte("SET"),
			[]byte(key),
			value,
			[]byte("PX"),
			[]byte(strconv.FormatInt(ttlMs, 10)),
		); err != nil {
			return err
		}
	} else {
		if err := writeCommand(writer, []byte("SET"), []byte(key), value); err != nil {
			return err
		}
	}

	return expectOK(reader)
}

func (s *RedisStore) Close() error {
	return nil
}

func (s *RedisStore) openConn(ctx context.Context) (net.Conn, *bufio.Reader, *bufio.Writer, error) {
	dialer := net.Dialer{
		Timeout: s.cfg.DialTimeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", s.cfg.Addr)
	if err != nil {
		return nil, nil, nil, err
	}

	// Respect request timeout/deadline when available.
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else if s.cfg.IOTimeout > 0 {
		_ = conn.SetDeadline(time.Now().Add(s.cfg.IOTimeout))
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	if s.cfg.Password != "" {
		if s.cfg.Username != "" {
			if err := writeCommand(
				writer,
				[]byte("AUTH"),
				[]byte(s.cfg.Username),
				[]byte(s.cfg.Password),
			); err != nil {
				_ = conn.Close()
				return nil, nil, nil, err
			}
		} else {
			if err := writeCommand(writer, []byte("AUTH"), []byte(s.cfg.Password)); err != nil {
				_ = conn.Close()
				return nil, nil, nil, err
			}
		}
		if err := expectOK(reader); err != nil {
			_ = conn.Close()
			return nil, nil, nil, err
		}
	}

	if s.cfg.DB > 0 {
		if err := writeCommand(writer, []byte("SELECT"), []byte(strconv.Itoa(s.cfg.DB))); err != nil {
			_ = conn.Close()
			return nil, nil, nil, err
		}
		if err := expectOK(reader); err != nil {
			_ = conn.Close()
			return nil, nil, nil, err
		}
	}

	return conn, reader, writer, nil
}

func writeCommand(w *bufio.Writer, args ...[]byte) error {
	if _, err := w.WriteString("*" + strconv.Itoa(len(args)) + "\r\n"); err != nil {
		return err
	}
	for _, arg := range args {
		if _, err := w.WriteString("$" + strconv.Itoa(len(arg)) + "\r\n"); err != nil {
			return err
		}
		if _, err := w.Write(arg); err != nil {
			return err
		}
		if _, err := w.WriteString("\r\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

func expectOK(r *bufio.Reader) error {
	prefix, data, err := readRESP(r)
	if err != nil {
		return err
	}
	if prefix != '+' || string(data) != "OK" {
		return fmt.Errorf("unexpected redis response: prefix=%q value=%q", string(prefix), string(data))
	}
	return nil
}

func readRESP(r *bufio.Reader) (byte, []byte, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return 0, nil, err
	}

	switch prefix {
	case '+': // Simple String
		line, err := readLine(r)
		if err != nil {
			return 0, nil, err
		}
		return prefix, []byte(line), nil
	case '-': // Error
		line, err := readLine(r)
		if err != nil {
			return 0, nil, err
		}
		return 0, nil, fmt.Errorf("redis error: %s", line)
	case ':': // Integer
		line, err := readLine(r)
		if err != nil {
			return 0, nil, err
		}
		return prefix, []byte(line), nil
	case '$': // Bulk String
		line, err := readLine(r)
		if err != nil {
			return 0, nil, err
		}
		size, err := strconv.Atoi(line)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid redis bulk size %q: %w", line, err)
		}
		if size == -1 {
			return prefix, nil, nil
		}
		buf := make([]byte, size+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, nil, err
		}
		return prefix, buf[:size], nil
	default:
		return 0, nil, fmt.Errorf("unsupported redis response type %q", string(prefix))
	}
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
}

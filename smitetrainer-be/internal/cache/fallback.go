package cache

import (
	"context"
	"errors"
	"time"
)

type FallbackStore struct {
	primary   Store
	secondary Store
}

func NewFallbackStore(primary, secondary Store) *FallbackStore {
	return &FallbackStore{
		primary:   primary,
		secondary: secondary,
	}
}

func (s *FallbackStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	var primaryErr error

	if s.primary != nil {
		value, found, err := s.primary.Get(ctx, key)
		if err == nil {
			if found || s.secondary == nil {
				return value, found, nil
			}
		} else {
			primaryErr = err
		}
	}

	if s.secondary == nil {
		if primaryErr != nil {
			return nil, false, primaryErr
		}
		return nil, false, nil
	}

	value, found, secondaryErr := s.secondary.Get(ctx, key)
	if secondaryErr != nil {
		if primaryErr != nil {
			return nil, false, errors.Join(primaryErr, secondaryErr)
		}
		return nil, false, secondaryErr
	}

	if found && s.primary != nil {
		_ = s.primary.Set(ctx, key, value, 0)
	}
	return value, found, nil
}

func (s *FallbackStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var errs []error

	if s.primary != nil {
		if err := s.primary.Set(ctx, key, value, ttl); err != nil {
			errs = append(errs, err)
		}
	}
	if s.secondary != nil {
		if err := s.secondary.Set(ctx, key, value, ttl); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (s *FallbackStore) Close() error {
	var errs []error
	if s.primary != nil {
		if err := s.primary.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if s.secondary != nil {
		if err := s.secondary.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

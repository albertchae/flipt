package cache

import (
	"context"

	flipt "github.com/markphelps/flipt/rpc/flipt"
	"github.com/markphelps/flipt/storage"
)

const flagCachePrefix = "flag:"

// GetFlag returns the flag from the cache if it exists; otherwise it delegates to the underlying store
// caching the result if no error
func (c *Store) GetFlag(ctx context.Context, k string) (*flipt.Flag, error) {
	key := flagCachePrefix + k

	// check if flag exists in cache
	if data, ok := c.cache.Get(key); ok {
		c.logger.Debugf("cache hit: %q", key)
		cacheHitTotal.WithLabelValues("memory").Inc()

		flag, ok := data.(*flipt.Flag)
		if !ok {
			// not flag, bad cache
			return nil, ErrCacheCorrupt
		}

		return flag, nil
	}

	// else, get it and add to cache
	flag, err := c.store.GetFlag(ctx, k)
	if err != nil {
		return flag, err
	}

	c.cache.Set(key, flag)
	c.logger.Debugf("cache miss; added: %q", key)
	cacheMissTotal.WithLabelValues("memory").Inc()

	return flag, nil
}

// ListFlags delegates to the underlying store
func (c *Store) ListFlags(ctx context.Context, opts ...storage.QueryOption) ([]*flipt.Flag, error) {
	return c.store.ListFlags(ctx, opts...)
}

// CreateFlag delegates to the underlying store, flushing the cache in the process
func (c *Store) CreateFlag(ctx context.Context, r *flipt.CreateFlagRequest) (*flipt.Flag, error) {
	c.cache.Flush()
	c.logger.Debug("flushed cache")
	return c.store.CreateFlag(ctx, r)
}

// UpdateFlag delegates to the underlying store, flushing the cache in the process
func (c *Store) UpdateFlag(ctx context.Context, r *flipt.UpdateFlagRequest) (*flipt.Flag, error) {
	c.cache.Flush()
	c.logger.Debug("flushed cache")
	return c.store.UpdateFlag(ctx, r)
}

// DeleteFlag delegates to the underlying store, flushing the cache in the process
func (c *Store) DeleteFlag(ctx context.Context, r *flipt.DeleteFlagRequest) error {
	c.cache.Flush()
	c.logger.Debug("flushed cache")
	return c.store.DeleteFlag(ctx, r)
}

// CreateVariant delegates to the underlying store, flushing the cache in the process
func (c *Store) CreateVariant(ctx context.Context, r *flipt.CreateVariantRequest) (*flipt.Variant, error) {
	c.cache.Flush()
	c.logger.Debug("flushed cache")
	return c.store.CreateVariant(ctx, r)
}

// UpdateVariant delegates to the underlying store, flushing the cache in the process
func (c *Store) UpdateVariant(ctx context.Context, r *flipt.UpdateVariantRequest) (*flipt.Variant, error) {
	c.cache.Flush()
	c.logger.Debug("flushed cache")
	return c.store.UpdateVariant(ctx, r)
}

// DeleteVariant delegates to the underlying store, flushing the cache in the process
func (c *Store) DeleteVariant(ctx context.Context, r *flipt.DeleteVariantRequest) error {
	c.cache.Flush()
	c.logger.Debug("flushed cache")
	return c.store.DeleteVariant(ctx, r)
}

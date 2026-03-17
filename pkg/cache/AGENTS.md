# CACHE ABSTRACTION

**Generated:** 2026-01-24 22:45:30
**Parent:** ../../AGENTS.md

## OVERVIEW
Cache abstraction layer with multiple implementations (Redis, Ristretto memory cache).

## STRUCTURE
```
cache/
├── cache.go                 # Cache interface definition
├── redis_cache.go          # Redis implementation
├── ristretto_memory.go     # Ristretto in-memory cache
├── ristretto_memory_test.go # Tests for memory cache
├── cache_factory.go        # Factory for creating cache instances
├── cache_key.go            # Cache key generation utilities
└── cache_metrics.go        # Metrics collection for cache operations
```

## WHERE TO LOOK
| Task | File | Purpose |
|------|------|---------|
| Define cache interface | `cache.go` | `Cache` interface with Get, Set, Delete methods |
| Redis implementation | `redis_cache.go` | Redis-backed cache with connection pooling |
| Memory cache | `ristretto_memory.go` | Ristretto in-memory cache implementation |
| Create cache instance | `cache_factory.go` | Factory methods for different cache types |
| Generate cache keys | `cache_key.go` | Consistent key generation across services |
| Monitor cache | `cache_metrics.go` | Metrics for hit/miss rates and latency |

## CONVENTIONS
- **Interface-based**: All caches implement the `Cache` interface
- **Key generation**: Use `cache_key.go` utilities for consistent key patterns
- **Metrics**: All cache operations include metrics collection
- **Error handling**: Silent degradation - cache failures shouldn't break main functionality
- **TTL management**: Configurable TTL with sensible defaults

## IMPLEMENTATIONS
- **Redis**: Production cache with persistence and distributed capabilities
- **Ristretto**: High-performance in-memory cache for local caching
- **Factory pattern**: `NewRedisCache()`, `NewRistrettoCache()` constructors

## USAGE PATTERNS
```go
// Typical usage in services
cache := cache.NewRedisCache(redisClient, defaultTTL)
value, err := cache.Get(ctx, key)
if err == cache.ErrNotFound {
    // Fetch from source
    cache.Set(ctx, key, value, ttl)
}
```

## ANTI-PATTERNS
- **DO NOT** hardcode cache keys - use `cache_key.go` utilities
- **DO NOT** ignore cache errors in critical paths - log and monitor
- **DO NOT** use cache as primary storage - always have fallback to source
- **AVOID** very long TTLs without invalidation logic
- **NEVER** store sensitive data without encryption in cache
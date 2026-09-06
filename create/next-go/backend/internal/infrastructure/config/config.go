package config

// Infrastructure configuration is kept separate so adapters can be replaced without changing use cases.
type Config struct { Port string }

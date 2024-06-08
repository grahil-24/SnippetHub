package main

// we can directly use string for checking if key exists in the context or not. But creating a user-defined
// type helps prevent name clashes. If a 3rd party lib, used the same key name, it can lead to clashes.
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

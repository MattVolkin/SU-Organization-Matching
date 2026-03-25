package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"server-example/ent"
	"server-example/ent/user"

	_ "github.com/lib/pq"
)

// DatabaseClient stores database connection state plus chainable query state.
type DatabaseClient struct {
	client *ent.Client
	driver string
	dsn    string
	alias  string

	userGoogleID string
	userEmail    string
	lastErr      error
}

var (
	databaseClientsMu sync.RWMutex
	databaseClients   = map[string]*DatabaseClient{}
)

// NewDatabaseClient initializes a new database client with an Ent connection.
func NewDatabaseClient(driver string, dsn string) (*DatabaseClient, error) {
	driver = strings.TrimSpace(driver)
	dsn = strings.TrimSpace(dsn)
	if driver == "" {
		return nil, fmt.Errorf("database driver is required")
	}
	if dsn == "" {
		return nil, fmt.Errorf("database dsn is required")
	}

	entClient, err := ent.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &DatabaseClient{
		client: entClient,
		driver: driver,
		dsn:    dsn,
	}, nil
}

// NewDatabaseClientFromEntClient builds a stateful wrapper around an existing Ent client.
func NewDatabaseClientFromEntClient(client *ent.Client, driver string, dsn string, alias string) *DatabaseClient {
	return (&DatabaseClient{
		client: client,
		driver: strings.TrimSpace(driver),
		dsn:    strings.TrimSpace(dsn),
		alias:  strings.TrimSpace(alias),
	}).clone()
}

// RegisterDatabaseClient stores a named DatabaseClient for reuse across files in main.
func RegisterDatabaseClient(name string, db *DatabaseClient) *DatabaseClient {
	name = strings.TrimSpace(name)
	if name == "" || db == nil {
		return db
	}
	databaseClientsMu.Lock()
	defer databaseClientsMu.Unlock()
	databaseClients[name] = db
	return db
}

// GetDatabaseClient returns a previously registered DatabaseClient by name.
func GetDatabaseClient(name string) (*DatabaseClient, bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, false
	}
	databaseClientsMu.RLock()
	defer databaseClientsMu.RUnlock()
	db, ok := databaseClients[name]
	return db, ok
}

// WithAlias returns a new state with the provided alias.
func (db *DatabaseClient) WithAlias(alias string) *DatabaseClient {
	next := db.clone()
	next.alias = strings.TrimSpace(alias)
	return next
}

// Query starts a new query state from the current connection.
func (db *DatabaseClient) Query() *DatabaseClient {
	next := db.clone()
	next.userGoogleID = ""
	next.userEmail = ""
	next.lastErr = nil
	return next
}

// Filter applies a simple key/value filter to the chain state.
func (db *DatabaseClient) Filter(field string, value string) *DatabaseClient {
	next := db.clone()
	field = strings.ToLower(strings.TrimSpace(field))
	value = strings.TrimSpace(value)

	switch field {
	case "google_id", "googleid", "sub":
		next.userGoogleID = value
	case "email":
		next.userEmail = value
	default:
		next.lastErr = fmt.Errorf("unsupported filter field: %s", field)
	}
	return next
}

// Merge combines query state from another client so chains can be composed.
func (db *DatabaseClient) Merge(other *DatabaseClient) *DatabaseClient {
	next := db.clone()
	if other == nil {
		return next
	}
	if strings.TrimSpace(next.userGoogleID) == "" {
		next.userGoogleID = strings.TrimSpace(other.userGoogleID)
	}
	if strings.TrimSpace(next.userEmail) == "" {
		next.userEmail = strings.TrimSpace(other.userEmail)
	}
	if next.lastErr == nil {
		next.lastErr = other.lastErr
	}
	return next
}

// Err returns the chain-state error, if any.
func (db *DatabaseClient) Err() error {
	if db == nil {
		return fmt.Errorf("database client is nil")
	}
	return db.lastErr
}

// Close closes the underlying Ent client connection.
func (db *DatabaseClient) Close() error {
	if db == nil || db.client == nil {
		return nil
	}
	return db.client.Close()
}

// FetchUserProfileByGoogleID reads one user profile by Google subject id.
func (db *DatabaseClient) FetchUserProfileByGoogleID(ctx context.Context, googleID string) (*DatabaseClient, *GoogleProfile, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}

	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	lookupID := strings.TrimSpace(googleID)
	if lookupID == "" {
		lookupID = strings.TrimSpace(next.userGoogleID)
	}
	if lookupID == "" {
		next.lastErr = fmt.Errorf("google id is required")
		return next, nil, next.lastErr
	}

	storedUser, err := next.client.User.Query().Where(user.GoogleIDEQ(lookupID)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, nil, err
	}

	if strings.TrimSpace(next.userEmail) != "" && !strings.EqualFold(strings.TrimSpace(storedUser.Email), strings.TrimSpace(next.userEmail)) {
		err = fmt.Errorf("email filter did not match user record")
		next.lastErr = err
		return next, nil, err
	}

	next.userGoogleID = lookupID
	next.userEmail = strings.TrimSpace(storedUser.Email)
	next.lastErr = nil

	return next, &GoogleProfile{
		Email: storedUser.Email,
		Name:  extractProfileNameTag(storedUser.Tags),
	}, nil
}

// UpsertUserProfileByGoogleID upserts user information using Google subject ID.
func (db *DatabaseClient) UpsertUserProfileByGoogleID(ctx context.Context, googleID string, userInfo *GoogleProfile) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}
	if userInfo == nil {
		next.lastErr = fmt.Errorf("user info is required")
		return next, next.lastErr
	}

	lookupID := strings.TrimSpace(googleID)
	if lookupID == "" {
		lookupID = strings.TrimSpace(next.userGoogleID)
	}
	email := strings.TrimSpace(userInfo.Email)

	if lookupID == "" {
		next.lastErr = fmt.Errorf("google id is required")
		return next, next.lastErr
	}
	if email == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, next.lastErr
	}

	tags := mergeProfileNameTag(nil, userInfo.Name)

	existing, err := next.client.User.Query().Where(user.GoogleIDEQ(lookupID)).Only(ctx)
	if err == nil {
		tags = mergeProfileNameTag(existing.Tags, userInfo.Name)
		err = next.client.User.UpdateOneID(existing.ID).
			SetEmail(email).
			SetTags(tags).
			Exec(ctx)
		next.lastErr = err
		next.userGoogleID = lookupID
		next.userEmail = email
		return next, err
	}

	if !ent.IsNotFound(err) {
		next.lastErr = err
		return next, err
	}

	_, createErr := next.client.User.Create().
		SetGoogleID(lookupID).
		SetEmail(email).
		SetTags(tags).
		Save(ctx)

	next.lastErr = createErr
	next.userGoogleID = lookupID
	next.userEmail = email
	return next, createErr
}

func (db *DatabaseClient) clone() *DatabaseClient {
	if db == nil {
		return &DatabaseClient{lastErr: fmt.Errorf("database client is nil")}
	}
	cp := *db
	return &cp
}

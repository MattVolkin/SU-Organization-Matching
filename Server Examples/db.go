package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"server-example/ent"
	"server-example/ent/answer"
	"server-example/ent/club"
	"server-example/ent/question"
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

func (db *DatabaseClient) getUserRole(ctx context.Context, email string) UserRole {
	next := db.clone()
	//here we assume that a user exists with the provided email, and that the email is unique across users.
	if next.lastErr != nil {
		return Member
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return Member
	}

	lookupEmail := strings.TrimSpace(email)

	_, isStudentLife, err := next.IsUserStudentLifeByEmail(ctx, lookupEmail)
	if err != nil {
		next.lastErr = err
		return Member
	}

	if isStudentLife {
		return Admin
	}

	_, isOfficer, err := next.IsUserOfficerByEmail(ctx, email)

	if err != nil {
		next.lastErr = err
		return Member
	}

	if isOfficer {
		return Officer
	}

	return Member
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

// FetchUserProfileByEmail reads one user profile by email address.
func (db *DatabaseClient) FetchUserProfileByEmail(ctx context.Context, email string) (*DatabaseClient, *GoogleProfile, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}

	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		lookupEmail = strings.TrimSpace(next.userEmail)
	}
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, nil, next.lastErr
	}

	storedUser, err := next.client.User.Query().Where(user.EmailEQ(lookupEmail)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, nil, err
	}

	next.userGoogleID = strings.TrimSpace(storedUser.GoogleID)
	next.userEmail = strings.TrimSpace(storedUser.Email)
	next.lastErr = nil

	return next, &GoogleProfile{
		Email: storedUser.Email,
		Name:  extractProfileNameTag(storedUser.Tags),
	}, nil
}

// FetchQuestionByID returns one question by primary key.
func (db *DatabaseClient) FetchQuestionByID(ctx context.Context, id int) (*DatabaseClient, *ent.Question, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}
	if id <= 0 {
		next.lastErr = fmt.Errorf("question id must be positive")
		return next, nil, next.lastErr
	}

	q, err := next.client.Question.Query().Where(question.IDEQ(id)).Only(ctx)
	next.lastErr = err
	return next, q, err
}

// FetchQuestionsByType returns all questions matching a given question_type.
func (db *DatabaseClient) FetchQuestionsByType(ctx context.Context, questionType string) (*DatabaseClient, []map[string]string, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	trimmedType := strings.TrimSpace(questionType)
	if trimmedType == "" {
		next.lastErr = fmt.Errorf("question type is required")
		return next, nil, next.lastErr
	}

	questions, err := next.client.Question.Query().Where(question.QuestionTypeEQ(trimmedType)).All(ctx)
	next.lastErr = err

	contents := make([]map[string]string, 0, len(questions))
	for _, q := range questions {
		if q == nil || q.Translations == nil {
			contents = append(contents, map[string]string{})
			continue
		}

		copyMap := make(map[string]string, len(q.Translations))
		for k, v := range q.Translations {
			copyMap[k] = v
		}
		contents = append(contents, copyMap)
	}
	return next, contents, err
}

// FetchAllQuestions returns all question records.
func (db *DatabaseClient) FetchAllQuestions(ctx context.Context) (*DatabaseClient, []*ent.Question, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	questions, err := next.client.Question.Query().All(ctx)
	next.lastErr = err
	return next, questions, err
}

// FetchAdjectiveQuestionContents returns translation payloads for all adjective questions.
func (db *DatabaseClient) FetchAdjectiveQuestionContents(ctx context.Context) (*DatabaseClient, []map[string]string, error) {
	next, questions, err := db.FetchQuestionsByType(ctx, "adjective")
	if err != nil {
		return next, nil, err
	}

	return next, questions, nil
}

// FetchClubByID returns one club by primary key.
func (db *DatabaseClient) FetchClubByID(ctx context.Context, id int) (*DatabaseClient, *ent.Club, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}
	if id <= 0 {
		next.lastErr = fmt.Errorf("club id must be positive")
		return next, nil, next.lastErr
	}

	c, err := next.client.Club.Query().Where(club.IDEQ(id)).Only(ctx)
	next.lastErr = err
	return next, c, err
}

// FetchAllClubs returns all club records.
func (db *DatabaseClient) FetchAllClubs(ctx context.Context) (*DatabaseClient, []*ent.Club, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	clubs, err := next.client.Club.Query().All(ctx)
	next.lastErr = err
	return next, clubs, err
}

func (db *DatabaseClient) UpdateClubFromJSON(ctx context.Context, newClubInfo *OfficerOrgJSON) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}
	if newClubInfo == nil {
		next.lastErr = fmt.Errorf("payload is required")
		return next, next.lastErr
	}

	clubID := newClubInfo.ID
	if clubID <= 0 {
		next.lastErr = fmt.Errorf("club id must be positive")
		return next, next.lastErr
	}

	update := next.client.Club.UpdateOneID(clubID)
	update.SetClubName(strings.TrimSpace(newClubInfo.ClubName))
	update.SetDescription(strings.TrimSpace(newClubInfo.Description))
	update.SetMeetingTime(strings.TrimSpace(newClubInfo.MeetingTime))
	update.SetImagePath(strings.TrimSpace(newClubInfo.ImagePath))
	update.SetExternalLink(strings.TrimSpace(newClubInfo.ExternalLink))
	update.SetContactInfo(strings.TrimSpace(newClubInfo.ContactInfo))
	update.SetIncludeOfficerEmails(newClubInfo.IncludeOfficerEmails)

	err := update.Exec(ctx)
	next.lastErr = err
	return next, err
}

// FetchOfficerClubsByUserEmail returns all clubs where the given user is listed
// as a leader/officer. The user is identified by email.
func (db *DatabaseClient) FetchOfficerClubsByUserEmail(ctx context.Context, email string) (*DatabaseClient, []*ent.Club, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		lookupEmail = strings.TrimSpace(next.userEmail)
	}
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, nil, next.lastErr
	}

	clubs, err := next.client.Club.Query().Where(club.HasLeadersWith(user.EmailEQ(lookupEmail))).All(ctx)
	next.userEmail = lookupEmail
	next.lastErr = err
	return next, clubs, err
}

func (db *DatabaseClient) IsUserStudentLifeByEmail(ctx context.Context, email string) (*DatabaseClient, bool, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, false, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, false, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		lookupEmail = strings.TrimSpace(next.userEmail)
	}
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, false, next.lastErr
	}

	storedUser, err := next.client.User.Query().Where(user.EmailEQ(lookupEmail)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, false, err
	}

	for _, tag := range storedUser.Tags {
		if strings.EqualFold(strings.TrimSpace(tag), "student_life") {
			next.userEmail = lookupEmail
			return next, true, nil
		}
	}

	next.userEmail = lookupEmail
	return next, false, nil
}

// IsUserOfficerByEmail efficiently checks whether a user (identified by email)
// is an officer/leader of at least one club.
func (db *DatabaseClient) IsUserOfficerByEmail(ctx context.Context, email string) (*DatabaseClient, bool, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, false, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, false, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		lookupEmail = strings.TrimSpace(next.userEmail)
	}
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, false, next.lastErr
	}

	isOfficer, err := next.client.Club.Query().Where(club.HasLeadersWith(user.EmailEQ(lookupEmail))).Exist(ctx)
	next.userEmail = lookupEmail
	next.lastErr = err
	return next, isOfficer, err
}

// FetchAnswerByID returns one answer by primary key.
func (db *DatabaseClient) FetchAnswerByID(ctx context.Context, id int) (*DatabaseClient, *ent.Answer, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}
	if id <= 0 {
		next.lastErr = fmt.Errorf("answer id must be positive")
		return next, nil, next.lastErr
	}

	a, err := next.client.Answer.Query().Where(answer.IDEQ(id)).Only(ctx)
	next.lastErr = err
	return next, a, err
}

// FetchAnswersByQuestionID returns all answers belonging to one question.
func (db *DatabaseClient) FetchAnswersByQuestionID(ctx context.Context, questionID int) (*DatabaseClient, []*ent.Answer, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}
	if questionID <= 0 {
		next.lastErr = fmt.Errorf("question id must be positive")
		return next, nil, next.lastErr
	}

	answers, err := next.client.Answer.Query().Where(answer.HasQuestionWith(question.IDEQ(questionID))).All(ctx)
	next.lastErr = err
	return next, answers, err
}

// FetchAnswersByUserEmail returns all answers submitted by one user email.
func (db *DatabaseClient) FetchAnswersByUserEmail(ctx context.Context, email string) (*DatabaseClient, []*ent.Answer, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		lookupEmail = strings.TrimSpace(next.userEmail)
	}
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, nil, next.lastErr
	}

	answers, err := next.client.Answer.Query().Where(answer.HasUserWith(user.EmailEQ(lookupEmail))).All(ctx)
	next.lastErr = err
	return next, answers, err
}

// FetchAllAnswers returns all answer records.
func (db *DatabaseClient) FetchAllAnswers(ctx context.Context) (*DatabaseClient, []*ent.Answer, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}

	answers, err := next.client.Answer.Query().All(ctx)
	next.lastErr = err
	return next, answers, err
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

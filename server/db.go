package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"su-organization-matching/server/ent"
	"su-organization-matching/server/ent/answer"
	"su-organization-matching/server/ent/club"
	"su-organization-matching/server/ent/question"
	"su-organization-matching/server/ent/user"

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
	resultRows   []map[string]any
	lastErr      error
}

type DBAnswer struct {
	ID           int
	AnswerText   string
	SubmittedAt  time.Time
	QuestionID   int
	QuestionType string
	Translations map[string][]string
	IsActive     bool
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
	next.resultRows = nil
	next.lastErr = nil
	return next
}

// AppendResults appends one result set onto the in-flight accumulated results.
// This is useful for composing results from multiple requests into one payload.
func (db *DatabaseClient) AppendResults(rows []map[string]any) *DatabaseClient {
	next := db.clone()
	if next.lastErr != nil {
		return next
	}

	if len(rows) == 0 {
		return next
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		rowCopy := make(map[string]any, len(row))
		for k, v := range row {
			rowCopy[k] = v
		}
		next.resultRows = append(next.resultRows, rowCopy)
	}

	return next
}

// Results returns a defensive copy of accumulated result rows.
func (db *DatabaseClient) Results() []map[string]any {
	if db == nil || len(db.resultRows) == 0 {
		return []map[string]any{}
	}

	out := make([]map[string]any, 0, len(db.resultRows))
	for _, row := range db.resultRows {
		if row == nil {
			continue
		}
		rowCopy := make(map[string]any, len(row))
		for k, v := range row {
			rowCopy[k] = v
		}
		out = append(out, rowCopy)
	}
	return out
}

func (db *DatabaseClient) FetchUserAnswersByUserEmail(ctx context.Context, email string) (*DatabaseClient, []DBAnswer, error) {
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

	answers, err := next.client.Answer.Query().
		Where(answer.HasUserWith(user.EmailEQ(lookupEmail))).
		WithQuestion().
		All(ctx)
	if err != nil {
		next.lastErr = err
		return next, nil, err
	}

	rows := make([]DBAnswer, 0, len(answers))
	for _, a := range answers {
		if a == nil {
			continue
		}
		row := DBAnswer{
			ID:          a.ID,
			AnswerText:  a.AnswerText,
			SubmittedAt: a.SubmittedAt,
		}
		if q := a.Edges.Question; q != nil {
			row.QuestionID = q.ID
			row.QuestionType = q.QuestionType
			row.Translations = q.Translations
			row.IsActive = q.IsActive
		}
		rows = append(rows, row)
	}

	return next, rows, nil
}

// FetchQuestionsByTypeAndAppend runs FetchQuestionsByType and appends the
// returned rows to accumulated chain results for fluent composition.
func (db *DatabaseClient) FetchQuestionsByTypeAndAppend(ctx context.Context, questionType string) *DatabaseClient {
	next, rows, err := db.FetchQuestionsByType(ctx, questionType)
	if err != nil {
		next.lastErr = err
		return next
	}

	return next.AppendResults(rows)
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

	_, isOfficer, err := next.IsUserOfficerByEmail(ctx, lookupEmail)
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

// FetchUserByEmail returns the full user row for matching and profile decisions.
func (db *DatabaseClient) FetchUserByEmail(ctx context.Context, email string) (*DatabaseClient, *ent.User, error) {
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
	next.lastErr = err
	if err != nil {
		return next, nil, err
	}

	next.userGoogleID = strings.TrimSpace(storedUser.GoogleID)
	next.userEmail = strings.TrimSpace(storedUser.Email)
	return next, storedUser, nil
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

// FetchQuestionsByType returns all questions matching a given question_type,
// including each question's ID alongside its translation payload.
func (db *DatabaseClient) FetchQuestionsByType(ctx context.Context, questionType string) (*DatabaseClient, []map[string]any, error) {
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

	contents := make([]map[string]any, 0, len(questions))
	for _, q := range questions {
		if q == nil {
			continue
		}

		questionPayload := make(map[string]any, len(q.Translations)+3)
		questionPayload["id"] = q.ID
		questionPayload["question_type"] = q.QuestionType
		questionPayload["translations"] = q.Translations

		if q.Translations == nil {
			contents = append(contents, questionPayload)
			continue
		}

		contents = append(contents, questionPayload)
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

// FetchSwipeQuestionContents returns question IDs plus translation payloads
// for all adjective and personality traits questions.
func (db *DatabaseClient) FetchSwipeQuestionContents(ctx context.Context) (*DatabaseClient, []map[string]any, error) {
	questions := db.FetchQuestionsByTypeAndAppend(ctx, "activities").FetchQuestionsByTypeAndAppend(ctx, "personality_traits").Results()

	return db, questions, nil
}

// ReplaceSwipeQuestionsForTesting fully replaces activities and personality
// trait question sets in one transaction.
func (db *DatabaseClient) ReplaceSwipeQuestionsForTesting(ctx context.Context, activities []SwipeQuestionInput, personalityTraits []SwipeQuestionInput) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}

	normalizedActivities, err := normalizeSwipeQuestionInputs(activities, "activities")
	if err != nil {
		next.lastErr = err
		return next, err
	}
	normalizedPersonalityTraits, err := normalizeSwipeQuestionInputs(personalityTraits, "personalityTraits")
	if err != nil {
		next.lastErr = err
		return next, err
	}

	tx, err := next.client.Tx(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}
	defer func() {
		if next.lastErr != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.Question.Delete().Where(question.QuestionTypeIn("activities", "personality_traits")).Exec(ctx); err != nil {
		next.lastErr = err
		return next, err
	}

	for _, item := range normalizedActivities {
		if _, err := tx.Question.Create().
			SetQuestionType("activities").
			SetTranslations(item.Translations).
			SetIsActive(true).
			Save(ctx); err != nil {
			next.lastErr = err
			return next, err
		}
	}

	for _, item := range normalizedPersonalityTraits {
		if _, err := tx.Question.Create().
			SetQuestionType("personality_traits").
			SetTranslations(item.Translations).
			SetIsActive(true).
			Save(ctx); err != nil {
			next.lastErr = err
			return next, err
		}
	}

	if err := tx.Commit(); err != nil {
		next.lastErr = err
		return next, err
	}

	next.lastErr = nil
	return next, nil
}

func normalizeSwipeQuestionInputs(items []SwipeQuestionInput, fieldName string) ([]SwipeQuestionInput, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("%s must include at least one question", fieldName)
	}

	normalized := make([]SwipeQuestionInput, 0, len(items))
	seenTerms := make(map[string]struct{}, len(items))

	for index, item := range items {
		translations := normalizeQuestionTranslations(item.Translations)

		term := strings.TrimSpace(item.Term)
		if term == "" {
			term = pickTermFromTranslations(translations)
		}
		if term == "" {
			return nil, fmt.Errorf("%s[%d].term is required", fieldName, index)
		}

		definition := strings.TrimSpace(item.Definition)
		if definition == "" {
			definition = pickDefinitionFromTranslations(translations)
		}

		translations = withEnglishTermAndDefinition(translations, term, definition)

		termKey := strings.ToLower(term)
		if _, exists := seenTerms[termKey]; exists {
			return nil, fmt.Errorf("%s includes duplicate term %q", fieldName, term)
		}
		seenTerms[termKey] = struct{}{}

		normalized = append(normalized, SwipeQuestionInput{
			Term:         term,
			Definition:   definition,
			Translations: translations,
		})
	}

	return normalized, nil
}

func normalizeQuestionTranslations(raw map[string][]string) map[string][]string {
	if len(raw) == 0 {
		return map[string][]string{}
	}

	cleaned := make(map[string][]string, len(raw))
	for key, values := range raw {
		lang := strings.ToLower(strings.TrimSpace(key))
		if lang == "" {
			continue
		}

		normalizedValues := make([]string, 0, len(values))
		for _, value := range values {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				continue
			}
			normalizedValues = append(normalizedValues, trimmed)
		}

		if len(normalizedValues) == 0 {
			continue
		}

		cleaned[lang] = normalizedValues
	}

	return cleaned
}

func pickTermFromTranslations(translations map[string][]string) string {
	if en, ok := translations["en"]; ok && len(en) > 0 {
		return strings.TrimSpace(en[0])
	}
	if term, ok := translations["term"]; ok && len(term) > 0 {
		return strings.TrimSpace(term[0])
	}
	for _, values := range translations {
		if len(values) == 0 {
			continue
		}
		if v := strings.TrimSpace(values[0]); v != "" {
			return v
		}
	}
	return ""
}

func pickDefinitionFromTranslations(translations map[string][]string) string {
	if en, ok := translations["en"]; ok && len(en) > 1 {
		if v := strings.TrimSpace(en[1]); v != "" {
			return v
		}
	}
	if term, ok := translations["term"]; ok && len(term) > 1 {
		if v := strings.TrimSpace(term[1]); v != "" {
			return v
		}
	}
	return ""
}

func withEnglishTermAndDefinition(translations map[string][]string, term string, definition string) map[string][]string {
	result := make(map[string][]string, len(translations)+1)
	for key, values := range translations {
		copied := make([]string, len(values))
		copy(copied, values)
		result[key] = copied
	}

	english := []string{strings.TrimSpace(term)}
	if trimmedDefinition := strings.TrimSpace(definition); trimmedDefinition != "" {
		english = append(english, trimmedDefinition)
	}
	result["en"] = english

	return result
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

	clubs, err := next.client.Club.Query().WithLeaders().All(ctx)
	next.lastErr = err
	return next, clubs, err
}

func (db *DatabaseClient) UpdateClubFromJSON(ctx context.Context, newClubInfo *OrgJSON) (*DatabaseClient, error) {
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
	update.SetPersonality(newClubInfo.Personality)
	update.SetActivities(newClubInfo.Activities)
	update.SetGenders(newClubInfo.Genders)
	update.SetEthnicities(newClubInfo.Ethnicities)
	update.SetReligions(newClubInfo.Religions)
	update.SetStrictGenders(newClubInfo.StrictGenders)
	update.SetDedicatedMajors(newClubInfo.DedicatedMajors)
	update.SetOther(newClubInfo.Other)
	leaders, err := next.resolveUsersByEmail(ctx, newClubInfo.Officers)
	if err != nil {
		next.lastErr = err
		return next, err
	}
	if len(leaders) > 0 {
		update.AddLeaders(leaders...)
	}

	err = update.Exec(ctx)
	next.lastErr = err
	return next, err
}

func (db *DatabaseClient) PatchClubFromJSON(ctx context.Context, clubID int, patch *OrgUpdatePayload) (*DatabaseClient, *ent.Club, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}
	if patch == nil {
		next.lastErr = fmt.Errorf("payload is required")
		return next, nil, next.lastErr
	}
	if clubID <= 0 {
		next.lastErr = fmt.Errorf("club id must be positive")
		return next, nil, next.lastErr
	}

	update := next.client.Club.UpdateOneID(clubID)
	hasUpdate := false
	if patch.ClubName != nil {
		update.SetClubName(strings.TrimSpace(*patch.ClubName))
		hasUpdate = true
	}
	if patch.Description != nil {
		update.SetDescription(strings.TrimSpace(*patch.Description))
		hasUpdate = true
	}
	if patch.MeetingTime != nil {
		update.SetMeetingTime(strings.TrimSpace(*patch.MeetingTime))
		hasUpdate = true
	}
	if patch.ImagePath != nil {
		update.SetImagePath(strings.TrimSpace(*patch.ImagePath))
		hasUpdate = true
	}
	if patch.ExternalLink != nil {
		update.SetExternalLink(strings.TrimSpace(*patch.ExternalLink))
		hasUpdate = true
	}
	if patch.ContactInfo != nil {
		update.SetContactInfo(strings.TrimSpace(*patch.ContactInfo))
		hasUpdate = true
	}
	if patch.Personality != nil {
		update.SetPersonality(*patch.Personality)
		hasUpdate = true
	}
	if patch.Activities != nil {
		update.SetActivities(*patch.Activities)
		hasUpdate = true
	}
	if patch.Genders != nil {
		update.SetGenders(*patch.Genders)
		hasUpdate = true
	}
	if patch.Ethnicities != nil {
		update.SetEthnicities(*patch.Ethnicities)
		hasUpdate = true
	}
	if patch.Religions != nil {
		update.SetReligions(*patch.Religions)
		hasUpdate = true
	}
	if patch.StrictGenders != nil {
		update.SetStrictGenders(*patch.StrictGenders)
		hasUpdate = true
	}
	if patch.DedicatedMajors != nil {
		update.SetDedicatedMajors(*patch.DedicatedMajors)
		hasUpdate = true
	}
	if patch.Other != nil {
		update.SetOther(*patch.Other)
		hasUpdate = true
	}
	if patch.IncludeOfficerEmails != nil {
		update.SetIncludeOfficerEmails(*patch.IncludeOfficerEmails)
		hasUpdate = true
	}
	if patch.Officers != nil {
		leaders, err := next.resolveUsersByEmail(ctx, *patch.Officers)
		if err != nil {
			next.lastErr = err
			return next, nil, err
		}
		update.ClearLeaders()
		if len(leaders) > 0 {
			update.AddLeaders(leaders...)
		}
		hasUpdate = true
	}
	if !hasUpdate {
		next.lastErr = fmt.Errorf("at least one field must be provided for update")
		return next, nil, next.lastErr
	}

	updatedClub, err := update.Save(ctx)
	next.lastErr = err
	if err != nil {
		return next, nil, err
	}

	updatedClub, err = next.client.Club.Query().Where(club.IDEQ(updatedClub.ID)).WithLeaders().Only(ctx)
	next.lastErr = err
	if err != nil {
		return next, nil, err
	}

	return next, updatedClub, nil
}

func (db *DatabaseClient) CreateClubFromJSON(ctx context.Context, newClubInfo *OrgJSON) (*DatabaseClient, *ent.Club, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, nil, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, nil, next.lastErr
	}
	if newClubInfo == nil {
		next.lastErr = fmt.Errorf("payload is required")
		return next, nil, next.lastErr
	}

	clubName := strings.TrimSpace(newClubInfo.ClubName)
	if clubName == "" {
		next.lastErr = fmt.Errorf("clubName is required")
		return next, nil, next.lastErr
	}

	create := next.client.Club.Create().
		SetClubName(clubName).
		SetDescription(strings.TrimSpace(newClubInfo.Description)).
		SetMeetingTime(strings.TrimSpace(newClubInfo.MeetingTime)).
		SetImagePath(strings.TrimSpace(newClubInfo.ImagePath)).
		SetExternalLink(strings.TrimSpace(newClubInfo.ExternalLink)).
		SetContactInfo(strings.TrimSpace(newClubInfo.ContactInfo)).
		SetIncludeOfficerEmails(newClubInfo.IncludeOfficerEmails).
		SetPersonality(newClubInfo.Personality).
		SetActivities(newClubInfo.Activities).
		SetGenders(newClubInfo.Genders).
		SetEthnicities(newClubInfo.Ethnicities).
		SetReligions(newClubInfo.Religions).
		SetStrictGenders(newClubInfo.StrictGenders).
		SetDedicatedMajors(newClubInfo.DedicatedMajors).
		SetOther(newClubInfo.Other)

	leaders, err := next.resolveUsersByEmail(ctx, newClubInfo.Officers)
	if err != nil {
		next.lastErr = err
		return next, nil, err
	}
	if len(leaders) > 0 {
		create.AddLeaders(leaders...)
	}

	createdClub, err := create.Save(ctx)
	if err != nil {
		next.lastErr = err
		return next, nil, err
	}

	createdClub, err = next.client.Club.Query().Where(club.IDEQ(createdClub.ID)).WithLeaders().Only(ctx)
	next.lastErr = err
	if err != nil {
		return next, nil, err
	}

	return next, createdClub, nil
}

func (db *DatabaseClient) DeleteClubByID(ctx context.Context, clubID int) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}
	if clubID <= 0 {
		next.lastErr = fmt.Errorf("club id must be positive")
		return next, next.lastErr
	}

	err := next.client.Club.DeleteOneID(clubID).Exec(ctx)
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

	clubs, err := next.client.Club.Query().Where(club.HasLeadersWith(user.EmailEQ(lookupEmail))).WithLeaders().All(ctx)
	next.userEmail = lookupEmail
	next.lastErr = err
	return next, clubs, err
}

func (db *DatabaseClient) resolveUsersByEmail(ctx context.Context, emails []string) ([]*ent.User, error) {
	if len(emails) == 0 {
		return []*ent.User{}, nil
	}

	normalizedEmails := make([]string, 0, len(emails))
	seen := make(map[string]struct{}, len(emails))
	for _, email := range emails {
		trimmed := strings.TrimSpace(email)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		normalizedEmails = append(normalizedEmails, trimmed)
	}

	if len(normalizedEmails) == 0 {
		return []*ent.User{}, nil
	}

	usersByEmail, err := db.client.User.Query().Where(user.EmailIn(normalizedEmails...)).All(ctx)
	if err != nil {
		return nil, err
	}

	resolved := make(map[string]*ent.User, len(usersByEmail))
	for _, storedUser := range usersByEmail {
		if storedUser == nil {
			continue
		}
		resolved[strings.ToLower(strings.TrimSpace(storedUser.Email))] = storedUser
	}

	users := make([]*ent.User, 0, len(normalizedEmails))
	for _, email := range normalizedEmails {
		storedUser, ok := resolved[strings.ToLower(strings.TrimSpace(email))]
		if !ok {
			return nil, fmt.Errorf("officer with email %q not found", email)
		}
		users = append(users, storedUser)
	}

	return users, nil
}

// SetUserRoleForTesting updates the persisted fields that determine effective role.
// - admin: adds student_life tag
// - officer: removes student_life tag and adds user as leader of officerClubName
// - member: removes student_life tag and clears all led clubs
func (db *DatabaseClient) SetUserRoleForTesting(ctx context.Context, email string, role string, officerClubName string) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, next.lastErr
	}

	normalizedRole := strings.ToLower(strings.TrimSpace(role))
	if normalizedRole == "" {
		next.lastErr = fmt.Errorf("role is required")
		return next, next.lastErr
	}

	tx, err := next.client.Tx(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}
	defer func() {
		if next.lastErr != nil {
			_ = tx.Rollback()
		}
	}()

	storedUser, err := tx.User.Query().Where(user.EmailEQ(lookupEmail)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}

	switch normalizedRole {
	case "admin":
		tags := withStudentLifeTag(storedUser.Tags, true)
		err = tx.User.UpdateOneID(storedUser.ID).SetTags(tags).Exec(ctx)
		if err != nil {
			next.lastErr = err
			return next, err
		}

	case "officer":
		tags := withStudentLifeTag(storedUser.Tags, false)
		err = tx.User.UpdateOneID(storedUser.ID).SetTags(tags).Exec(ctx)
		if err != nil {
			next.lastErr = err
			return next, err
		}

		clubName := strings.TrimSpace(officerClubName)
		if clubName == "" {
			next.lastErr = fmt.Errorf("officer club name is required")
			return next, next.lastErr
		}

		csClub, queryErr := tx.Club.Query().Where(club.ClubNameEQ(clubName)).Only(ctx)
		if queryErr != nil {
			next.lastErr = fmt.Errorf("officer club %q not found", clubName)
			return next, next.lastErr
		}

		alreadyLeader, queryErr := tx.Club.Query().
			Where(club.IDEQ(csClub.ID), club.HasLeadersWith(user.IDEQ(storedUser.ID))).
			Exist(ctx)
		if queryErr != nil {
			next.lastErr = queryErr
			return next, queryErr
		}
		if !alreadyLeader {
			err = tx.Club.UpdateOneID(csClub.ID).AddLeaderIDs(storedUser.ID).Exec(ctx)
			if err != nil {
				next.lastErr = err
				return next, err
			}
		}

	case "member":
		tags := withStudentLifeTag(storedUser.Tags, false)
		err = tx.User.UpdateOneID(storedUser.ID).SetTags(tags).ClearLedClubs().Exec(ctx)
		if err != nil {
			next.lastErr = err
			return next, err
		}

	default:
		next.lastErr = fmt.Errorf("unsupported role %q", role)
		return next, next.lastErr
	}

	if err = tx.Commit(); err != nil {
		next.lastErr = err
		return next, err
	}

	next.lastErr = nil
	next.userEmail = lookupEmail
	return next, nil
}

func withStudentLifeTag(tags []string, enabled bool) []string {
	filtered := make([]string, 0, len(tags)+1)
	hasStudentLife := false
	for _, rawTag := range tags {
		tag := strings.TrimSpace(rawTag)
		if tag == "" {
			continue
		}
		if strings.EqualFold(tag, "student_life") {
			hasStudentLife = true
			if !enabled {
				continue
			}
			tag = "student_life"
		}
		filtered = append(filtered, tag)
	}

	if enabled && !hasStudentLife {
		filtered = append(filtered, "student_life")
	}

	return filtered
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

// UpsertSurveyResponseByUserEmail stores one boolean response for a user/question pair.
// Existing answers are updated in place, otherwise a new answer row is created.
func (db *DatabaseClient) UpsertSurveyResponseByUserEmail(ctx context.Context, email string, questionID int, answerValue bool) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, next.lastErr
	}
	if questionID <= 0 {
		next.lastErr = fmt.Errorf("question id must be positive")
		return next, next.lastErr
	}

	storedUser, err := next.client.User.Query().Where(user.EmailEQ(lookupEmail)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}

	if _, err = next.client.Question.Query().Where(question.IDEQ(questionID)).Only(ctx); err != nil {
		next.lastErr = err
		return next, err
	}

	answerText := strconv.FormatBool(answerValue)
	now := time.Now()

	updatedCount, err := next.client.Answer.Update().
		Where(
			answer.HasUserWith(user.IDEQ(storedUser.ID)),
			answer.HasQuestionWith(question.IDEQ(questionID)),
		).
		SetAnswerText(answerText).
		SetSubmittedAt(now).
		Save(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}

	if updatedCount > 0 {
		next.lastErr = nil
		return next, nil
	}

	_, err = next.client.Answer.Create().
		SetAnswerText(answerText).
		SetQuestionID(questionID).
		SetUserID(storedUser.ID).
		SetSubmittedAt(now).
		Save(ctx)
	next.lastErr = err
	return next, err
}

// ReplaceSurveyResponsesByUserEmail fully replaces one user's survey responses.
// Existing responses are deleted, then the provided list is stored in one transaction.
func (db *DatabaseClient) ReplaceSurveyResponsesByUserEmail(ctx context.Context, email string, responses []SurveyResponsePayload) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, next.lastErr
	}

	// Normalize to one response per question, with later values winning.
	normalized := make([]SurveyResponsePayload, 0, len(responses))
	indexesByQuestionID := make(map[int]int, len(responses))
	for i, response := range responses {
		if response.QuestionID <= 0 {
			next.lastErr = fmt.Errorf("responses[%d].questionId must be a positive integer", i)
			return next, next.lastErr
		}

		if existingIndex, ok := indexesByQuestionID[response.QuestionID]; ok {
			normalized[existingIndex].Answer = response.Answer
			continue
		}

		indexesByQuestionID[response.QuestionID] = len(normalized)
		normalized = append(normalized, response)
	}

	tx, err := next.client.Tx(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}
	defer func() {
		if next.lastErr != nil {
			_ = tx.Rollback()
		}
	}()

	storedUser, err := tx.User.Query().Where(user.EmailEQ(lookupEmail)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}

	if _, err = tx.Answer.Delete().Where(answer.HasUserWith(user.IDEQ(storedUser.ID))).Exec(ctx); err != nil {
		next.lastErr = err
		return next, err
	}

	submittedAt := time.Now()
	for i, response := range normalized {
		if _, err = tx.Question.Query().Where(question.IDEQ(response.QuestionID)).Only(ctx); err != nil {
			next.lastErr = fmt.Errorf("responses[%d]: question %d does not exist", i, response.QuestionID)
			return next, next.lastErr
		}

		if _, err = tx.Answer.Create().
			SetAnswerText(strconv.FormatBool(response.Answer)).
			SetQuestionID(response.QuestionID).
			SetUserID(storedUser.ID).
			SetSubmittedAt(submittedAt).
			Save(ctx); err != nil {
			next.lastErr = err
			return next, err
		}
	}

	if err = tx.Commit(); err != nil {
		next.lastErr = err
		return next, err
	}

	next.lastErr = nil
	return next, nil
}

// UpsertUserDemographicsByEmail updates demographics fields for one user.
func (db *DatabaseClient) UpsertUserDemographicsByEmail(
	ctx context.Context,
	email string,
	name string,
	genders []string,
	ethnicities []string,
	religions []string,
	dedicatedMajors []string,
	other []string,
) (*DatabaseClient, error) {
	next := db.clone()
	if next.lastErr != nil {
		return next, next.lastErr
	}
	if next.client == nil {
		next.lastErr = fmt.Errorf("database not initialized")
		return next, next.lastErr
	}

	lookupEmail := strings.TrimSpace(email)
	if lookupEmail == "" {
		next.lastErr = fmt.Errorf("email is required")
		return next, next.lastErr
	}

	storedUser, err := next.client.User.Query().Where(user.EmailEQ(lookupEmail)).Only(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}

	normalizedGenders := normalizeUniqueTrimmed(genders)
	normalizedEthnicities := normalizeUniqueTrimmed(ethnicities)
	normalizedReligions := normalizeUniqueTrimmed(religions)
	normalizedMajors := normalizeUniqueTrimmed(dedicatedMajors)
	normalizedOther := normalizeUniqueTrimmed(other)
	tags := mergeProfileNameTag(storedUser.Tags, name)

	err = next.client.User.UpdateOneID(storedUser.ID).
		SetTags(tags).
		SetGenders(normalizedGenders).
		SetEthnicities(normalizedEthnicities).
		SetReligions(normalizedReligions).
		SetDedicatedMajors(normalizedMajors).
		SetOther(normalizedOther).
		Exec(ctx)
	if err != nil {
		next.lastErr = err
		return next, err
	}

	next.userGoogleID = strings.TrimSpace(storedUser.GoogleID)
	next.userEmail = strings.TrimSpace(storedUser.Email)
	next.lastErr = nil
	return next, nil
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
	if len(db.resultRows) > 0 {
		cp.resultRows = make([]map[string]any, 0, len(db.resultRows))
		for _, row := range db.resultRows {
			if row == nil {
				cp.resultRows = append(cp.resultRows, nil)
				continue
			}
			rowCopy := make(map[string]any, len(row))
			for k, v := range row {
				rowCopy[k] = v
			}
			cp.resultRows = append(cp.resultRows, rowCopy)
		}
	}
	return &cp
}

func normalizeUniqueTrimmed(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	return normalized
}

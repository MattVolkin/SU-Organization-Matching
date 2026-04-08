package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/lib/pq"
)

type questionImport struct {
	QuestionType string
	Term         string
	Definition   string
}

func main() {
	csvPath := flag.String("csv", "Question Planning - List of Adjectives.csv", "Path to the question CSV file")
	dsn := flag.String("dsn", "host=localhost port=5432 user=dev_user password=testing dbname=dev_project_db sslmode=disable", "Postgres DSN")
	dryRun := flag.Bool("dry-run", false, "Parse and report changes without writing to the database")
	flag.Parse()

	questions, err := parseQuestionCSV(*csvPath)
	if err != nil {
		log.Fatalf("failed parsing CSV: %v", err)
	}
	if len(questions) == 0 {
		log.Fatal("no questions found in CSV")
	}

	db, err := sql.Open("postgres", strings.TrimSpace(*dsn))
	if err != nil {
		log.Fatalf("failed opening DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed connecting to DB: %v", err)
	}

	existingByKey, err := loadExistingQuestionIndex(db, collectQuestionTypes(questions))
	if err != nil {
		log.Fatalf("failed loading existing questions: %v", err)
	}

	created := 0
	updated := 0

	for _, q := range questions {
		key := questionKey(q.QuestionType, q.Term)
		translations := map[string][]string{
			"en": []string{q.Term, q.Definition},
		}
		payload, err := json.Marshal(translations)
		if err != nil {
			log.Fatalf("failed encoding translations for %q: %v", q.Term, err)
		}

		if existingID, found := existingByKey[key]; found {
			if *dryRun {
				updated++
				log.Printf("[dry-run] update question id=%d type=%s term=%s", existingID, q.QuestionType, q.Term)
				continue
			}
			_, err = db.Exec(`UPDATE questions SET translations = $1::jsonb, is_active = TRUE WHERE id = $2`, string(payload), existingID)
			if err != nil {
				log.Fatalf("failed updating question id=%d: %v", existingID, err)
			}
			updated++
			continue
		}

		if *dryRun {
			created++
			log.Printf("[dry-run] create question type=%s term=%s", q.QuestionType, q.Term)
			continue
		}

		_, err = db.Exec(`INSERT INTO questions (question_type, translations, is_active) VALUES ($1, $2::jsonb, TRUE)`, q.QuestionType, string(payload))
		if err != nil {
			log.Fatalf("failed creating question type=%s term=%s: %v", q.QuestionType, q.Term, err)
		}
		created++
	}

	log.Printf("import complete: parsed=%d created=%d updated=%d dry_run=%t", len(questions), created, updated, *dryRun)
}

func parseQuestionCSV(path string) ([]questionImport, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	parsed := make([]questionImport, 0, 64)
	seen := make(map[string]struct{})

	for scanner.Scan() {
		line := strings.TrimSpace(strings.TrimSuffix(scanner.Text(), "\r"))
		if line == "" {
			continue
		}

		tag, term, def, ok := parseQuestionLine(line)
		if !ok {
			continue
		}
		if strings.EqualFold(tag, "tag") || strings.EqualFold(term, "word") {
			continue
		}

		questionType := mapTagToQuestionType(tag)
		if questionType == "" {
			continue
		}

		key := questionKey(questionType, term)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		parsed = append(parsed, questionImport{
			QuestionType: questionType,
			Term:         term,
			Definition:   def,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parsed, nil
}

func parseQuestionLine(line string) (tag string, term string, def string, ok bool) {
	tag, ok = extractBetween(line, `tag: """,`, `,""", term: """,`)
	if !ok {
		return "", "", "", false
	}

	term, ok = extractBetween(line, `term: """,`, `,""", def: """,`)
	if !ok {
		return "", "", "", false
	}

	def, ok = extractBetween(line, `def: """,`, `,"""}`)
	if !ok {
		def, ok = extractBetween(line, `def: """,`, `,"""},`)
		if !ok {
			def = ""
		}
	}

	tag = cleanCSVValue(tag)
	term = cleanCSVValue(term)
	def = cleanCSVValue(def)

	if tag == "" || term == "" {
		return "", "", "", false
	}

	return tag, term, def, true
}

func extractBetween(s string, left string, right string) (string, bool) {
	start := strings.Index(s, left)
	if start < 0 {
		return "", false
	}
	start += len(left)

	rest := s[start:]
	end := strings.Index(rest, right)
	if end < 0 {
		return "", false
	}

	return rest[:end], true
}

func cleanCSVValue(raw string) string {
	v := strings.TrimSpace(raw)
	v = strings.Trim(v, `"`)
	v = strings.TrimSpace(v)
	return v
}

func mapTagToQuestionType(tag string) string {
	normalized := strings.ToLower(strings.TrimSpace(tag))
	switch normalized {
	case "activities", "activity", "adjective":
		return "adjective"
	case "personality", "personality_traits", "personalitytraits", "personalitytrait":
		return "personality_traits"
	default:
		return ""
	}
}

func loadExistingQuestionIndex(db *sql.DB, questionTypes []string) (map[string]int, error) {
	if len(questionTypes) == 0 {
		return map[string]int{}, nil
	}

	rows, err := db.Query(`SELECT id, question_type, translations FROM questions WHERE question_type = ANY($1)`, pq.Array(questionTypes))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	index := make(map[string]int)
	for rows.Next() {
		var id int
		var qType string
		var translationsRaw []byte
		if err := rows.Scan(&id, &qType, &translationsRaw); err != nil {
			return nil, err
		}

		term := englishTermFromJSON(translationsRaw)
		if term == "" {
			continue
		}

		index[questionKey(qType, term)] = id
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return index, nil
}

func englishTermFromJSON(raw []byte) string {
	var generic map[string]any
	if err := json.Unmarshal(raw, &generic); err != nil {
		return ""
	}

	enRaw, ok := generic["en"]
	if !ok {
		return ""
	}

	switch v := enRaw.(type) {
	case string:
		return strings.TrimSpace(v)
	case []any:
		if len(v) == 0 {
			return ""
		}
		first, _ := v[0].(string)
		return strings.TrimSpace(first)
	default:
		return ""
	}
}

func questionKey(questionType string, term string) string {
	return strings.ToLower(strings.TrimSpace(questionType)) + "::" + strings.ToLower(strings.TrimSpace(term))
}

func collectQuestionTypes(questions []questionImport) []string {
	typesSet := make(map[string]struct{})
	for _, q := range questions {
		typesSet[q.QuestionType] = struct{}{}
	}

	types := make([]string, 0, len(typesSet))
	for qType := range typesSet {
		types = append(types, qType)
	}
	sort.Strings(types)

	return types
}

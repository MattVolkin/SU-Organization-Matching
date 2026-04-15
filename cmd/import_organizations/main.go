package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"su-organization-matching/server/ent"

	_ "github.com/lib/pq"
)

type clubImportData struct {
	Name             string
	Description      string
	MeetingTime      string
	MeetingLocation  string
	SocialMedia      string
	OtherLinks       []string
	ContactInfo      string
	ExternalLink     string
	ProfileOnly      bool
	Activities       []string
	Personality      []string
	Genders          []string
	Ethnicities      []string
	Religions        []string
	DedicatedMajors  []string
	AssociatedMajors []string
	StrictGenders    bool
	Other            []string
}

func main() {
	csvPath := flag.String("csv", "Question Planning - Reformatted Data.csv", "Path to the organization planning CSV file")
	dsn := flag.String("dsn", "host=localhost port=5432 user=dev_user password=testing dbname=dev_project_db sslmode=disable", "Postgres DSN")
	dryRun := flag.Bool("dry-run", false, "Parse and report changes without writing to the database")
	flag.Parse()

	records, err := readCSV(*csvPath)
	if err != nil {
		log.Fatalf("failed reading CSV: %v", err)
	}

	clubs, err := parseClubs(records)
	if err != nil {
		log.Fatalf("failed parsing CSV: %v", err)
	}
	if len(clubs) == 0 {
		log.Fatal("no organizations found in CSV")
	}

	client, err := ent.Open("postgres", strings.TrimSpace(*dsn))
	if err != nil {
		log.Fatalf("failed opening DB connection: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	existingClubs, err := client.Club.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed loading existing clubs: %v", err)
	}
	existingByName := make(map[string]*ent.Club, len(existingClubs))
	for _, c := range existingClubs {
		existingByName[normalizeKey(c.ClubName)] = c
	}

	created := 0
	updated := 0

	for _, row := range clubs {
		key := normalizeKey(row.Name)
		existing, found := existingByName[key]

		if *dryRun {
			if found {
				updated++
				log.Printf("[dry-run] update club: %s (id=%d)", row.Name, existing.ID)
			} else {
				created++
				log.Printf("[dry-run] create club: %s", row.Name)
			}
			continue
		}

		if found {
			update := client.Club.UpdateOneID(existing.ID).
				SetClubName(row.Name)

			if row.ProfileOnly {
				update.
					SetDescription(row.Description).
					SetMeetingTime(row.MeetingTime).
					SetExternalLink(row.ExternalLink).
					SetContactInfo(row.ContactInfo)
			} else {
				update.
					SetActivities(row.Activities).
					SetPersonality(row.Personality).
					SetGenders(row.Genders).
					SetEthnicities(row.Ethnicities).
					SetReligions(row.Religions).
					SetDedicatedMajors(row.DedicatedMajors).
					SetAssociatedMajors(row.AssociatedMajors).
					SetStrictGenders(row.StrictGenders).
					SetOther(row.Other)
			}

			err = update.Exec(ctx)
			if err != nil {
				log.Fatalf("failed updating club %q: %v", row.Name, err)
			}
			updated++
			continue
		}

		create := client.Club.Create().
			SetClubName(row.Name)

		if row.ProfileOnly {
			create.
				SetDescription(row.Description).
				SetMeetingTime(row.MeetingTime).
				SetExternalLink(row.ExternalLink).
				SetContactInfo(row.ContactInfo)
		} else {
			create.
				SetActivities(row.Activities).
				SetPersonality(row.Personality).
				SetGenders(row.Genders).
				SetEthnicities(row.Ethnicities).
				SetReligions(row.Religions).
				SetDedicatedMajors(row.DedicatedMajors).
				SetAssociatedMajors(row.AssociatedMajors).
				SetStrictGenders(row.StrictGenders).
				SetOther(row.Other)
		}

		_, err = create.Save(ctx)
		if err != nil {
			log.Fatalf("failed creating club %q: %v", row.Name, err)
		}
		created++
	}

	log.Printf("import complete: parsed=%d created=%d updated=%d dry_run=%t", len(clubs), created, updated, *dryRun)
}

func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = false

	var records [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func parseClubs(records [][]string) ([]clubImportData, error) {
	if isResultsPageFormat(records) {
		return parseProfileRows(records)
	}

	return parseLegacyMatrix(records)
}

func parseProfileRows(records [][]string) ([]clubImportData, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("csv is empty")
	}

	header := records[0]
	orgCol := findHeaderColumn(header, "organizations")
	meetingTimeCol := findHeaderColumn(header, "meetingtime")
	meetingLocationCol := findHeaderColumn(header, "meetinglocation")
	descriptionCol := findHeaderColumn(header, "description")
	socialMediaCol := findHeaderColumn(header, "socialmedia")
	otherLinksCol := findHeaderColumn(header, "otherlinks")

	if orgCol == -1 {
		return nil, fmt.Errorf("could not find Organizations header column")
	}

	clubs := make([]clubImportData, 0, len(records)-1)
	for rowIdx := 1; rowIdx < len(records); rowIdx++ {
		row := records[rowIdx]
		name := strings.TrimSpace(cell(row, orgCol))
		if name == "" {
			continue
		}

		meetingTime := strings.TrimSpace(cell(row, meetingTimeCol))
		meetingLocation := strings.TrimSpace(cell(row, meetingLocationCol))
		meetingTimeWithLocation := combineMeetingTimeAndLocation(meetingTime, meetingLocation)
		description := strings.TrimSpace(cell(row, descriptionCol))
		socialMedia := strings.TrimSpace(cell(row, socialMediaCol))
		otherLinks := collectOtherLinks(row, otherLinksCol)

		allLinks := make([]string, 0, len(otherLinks)+1)
		if socialMedia != "" {
			allLinks = append(allLinks, socialMedia)
		}
		allLinks = append(allLinks, otherLinks...)
		allLinks = dedupeStrings(allLinks)

		externalLink := ""
		if len(allLinks) > 0 {
			externalLink = allLinks[0]
		}

		contactLines := make([]string, 0, len(allLinks)+1)
		if meetingLocation != "" {
			contactLines = append(contactLines, "Meeting Location: "+meetingLocation)
		}
		if len(allLinks) > 0 {
			contactLines = append(contactLines, "Links: "+strings.Join(allLinks, ", "))
		}

		clubs = append(clubs, clubImportData{
			Name:            name,
			Description:     description,
			MeetingTime:     meetingTimeWithLocation,
			MeetingLocation: meetingLocation,
			SocialMedia:     socialMedia,
			OtherLinks:      otherLinks,
			ContactInfo:     strings.TrimSpace(strings.Join(contactLines, "\n")),
			ExternalLink:    externalLink,
			ProfileOnly:     true,
		})
	}

	if len(clubs) == 0 {
		return nil, fmt.Errorf("no organizations found in profile CSV")
	}

	return clubs, nil
}

func parseLegacyMatrix(records [][]string) ([]clubImportData, error) {
	orgRow := -1
	maxColumns := 0
	for i, row := range records {
		if len(row) > maxColumns {
			maxColumns = len(row)
		}
		if normalizeHeader(cell(row, 0)) == "organizations" {
			orgRow = i
		}
	}
	if orgRow == -1 {
		return nil, fmt.Errorf("could not find Organizations header row")
	}

	orgNames := make([]string, 0, maxColumns-1)
	for col := 1; col < maxColumns; col++ {
		name := strings.TrimSpace(cell(records[orgRow], col))
		orgNames = append(orgNames, name)
	}

	clubs := make([]clubImportData, len(orgNames))
	for i, name := range orgNames {
		clubs[i] = clubImportData{Name: name}
	}

	currentSection := ""
	for rowIdx := orgRow + 1; rowIdx < len(records); rowIdx++ {
		row := records[rowIdx]
		header := normalizeHeader(cell(row, 0))
		if header != "" {
			if header == "organizations" {
				continue
			}

			currentSection = canonicalSection(header)
		}

		if currentSection == "" {
			continue
		}

		for col := 1; col <= len(orgNames); col++ {
			value := strings.TrimSpace(cell(row, col))
			if value == "" {
				continue
			}
			if clubs[col-1].Name == "" {
				continue
			}

			switch currentSection {
			case "activities":
				clubs[col-1].Activities = append(clubs[col-1].Activities, value)
			case "personality":
				clubs[col-1].Personality = append(clubs[col-1].Personality, value)
			case "genders":
				clubs[col-1].Genders = append(clubs[col-1].Genders, value)
			case "ethnicities":
				clubs[col-1].Ethnicities = append(clubs[col-1].Ethnicities, value)
			case "religions":
				clubs[col-1].Religions = append(clubs[col-1].Religions, value)
			case "dedicatedmajors":
				clubs[col-1].DedicatedMajors = append(clubs[col-1].DedicatedMajors, value)
			case "associatedmajors":
				clubs[col-1].AssociatedMajors = append(clubs[col-1].AssociatedMajors, value)
			case "strictgenders":
				if parseTruthy(value) {
					clubs[col-1].StrictGenders = true
				}
			case "other":
				clubs[col-1].Other = append(clubs[col-1].Other, value)
			}
		}
	}

	result := make([]clubImportData, 0, len(clubs))
	for _, c := range clubs {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			continue
		}
		c.Name = name
		c.Activities = dedupeStrings(c.Activities)
		c.Personality = dedupeStrings(c.Personality)
		c.Genders = dedupeStrings(c.Genders)
		c.Ethnicities = dedupeStrings(c.Ethnicities)
		c.Religions = dedupeStrings(c.Religions)
		c.DedicatedMajors = dedupeStrings(c.DedicatedMajors)
		c.AssociatedMajors = dedupeStrings(c.AssociatedMajors)
		c.Other = dedupeStrings(c.Other)
		result = append(result, c)
	}

	return result, nil
}

func isResultsPageFormat(records [][]string) bool {
	if len(records) == 0 {
		return false
	}
	header := records[0]
	if len(header) == 0 {
		return false
	}

	return findHeaderColumn(header, "organizations") != -1 &&
		findHeaderColumn(header, "description") != -1 &&
		findHeaderColumn(header, "meetingtime") != -1
}

func findHeaderColumn(header []string, normalizedHeader string) int {
	for i, raw := range header {
		if normalizeHeader(raw) == normalizedHeader {
			return i
		}
	}
	return -1
}

func collectOtherLinks(row []string, otherLinksCol int) []string {
	if otherLinksCol == -1 {
		return nil
	}

	links := make([]string, 0)
	for i := otherLinksCol; i < len(row); i++ {
		value := strings.TrimSpace(cell(row, i))
		if value == "" {
			continue
		}
		links = append(links, value)
	}

	return dedupeStrings(links)
}

func combineMeetingTimeAndLocation(meetingTime string, meetingLocation string) string {
	meetingTime = strings.TrimSpace(meetingTime)
	meetingLocation = strings.TrimSpace(meetingLocation)

	switch {
	case meetingTime != "" && meetingLocation != "":
		return meetingTime + " - " + meetingLocation
	case meetingTime != "":
		return meetingTime
	default:
		return meetingLocation
	}
}

func canonicalSection(header string) string {
	switch header {
	case "activities", "activity":
		return "activities"
	case "personality":
		return "personality"
	case "gender", "genders":
		return "genders"
	case "raceethnicity", "raceethnicities", "ethnicity", "ethnicities":
		return "ethnicities"
	case "religion", "religions":
		return "religions"
	case "dedicatedmajor", "dedicatedmajors":
		return "dedicatedmajors"
	case "associatedmajor", "associatedmajors":
		return "associatedmajors"
	case "strictgender", "strictgenders":
		return "strictgenders"
	case "other":
		return "other"
	default:
		return ""
	}
}

func parseTruthy(raw string) bool {
	normalized := normalizeHeader(raw)
	switch normalized {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

func normalizeHeader(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimSuffix(trimmed, ":")
	trimmed = strings.ToLower(trimmed)
	trimmed = strings.ReplaceAll(trimmed, " ", "")
	trimmed = strings.ReplaceAll(trimmed, "/", "")
	trimmed = strings.ReplaceAll(trimmed, "_", "")
	return trimmed
}

func normalizeKey(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func dedupeStrings(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			continue
		}
		key := normalizeKey(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func cell(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

package matching

import (
	"sort"
	"strconv"
	"strings"
)

// Organization defines club attributes used by the ranking algorithm.
type Organization struct {
	ID               int
	Name             string
	Personality      []string
	Activities       []string
	Genders          []string
	Ethnicities      []string
	Religions        []string
	StrictGenders    bool
	DedicatedMajors  []string
	AssociatedMajors []string
	Other            []string
}

// UserInfo is the subset of user profile data used for matching.
type UserInfo struct {
	Name             string
	Personality      []string
	Activities       []string
	Genders          []string
	Ethnicities      []string
	Religions        []string
	DedicatedMajors  []string
	AssociatedMajors []string
	Other            []string
}

// Answer captures the minimal answer payload needed to reconstruct UserInfo.
type Answer struct {
	QuestionType string
	AnswerText   string
	Translations map[string][]string
}

// MatchResult represents one ranked organization for a user.
type MatchResult struct {
	Organization    Organization
	RawScore        float32
	NormalizedScore float32
}

// Section to define weights for each type of question
const personalityWeight float32 = 1.5
const activityWeight float32 = 2.25
const demographicWeight float32 = 1.25
const academicWeight float32 = 1.5
const otherWeight float32 = 2

// UserFromAnswers maps stored survey answers to algorithm input fields.
func UserFromAnswers(answers []Answer) UserInfo {
	user := UserInfo{}
	for _, a := range answers {
		questionType := normalizeKey(a.QuestionType)
		answerText := strings.TrimSpace(a.AnswerText)
		boolValue, hasBool := parseBool(answerText)

		switch questionType {
		case "adjective", "activity", "activities":
			if hasBool {
				if !boolValue {
					continue
				}
				addUnique(&user.Activities, preferredAnswerValue(a))
			} else {
				addUnique(&user.Activities, answerText)
			}
		case "personality", "personalitytrait", "personalitytraits":
			if hasBool {
				if !boolValue {
					continue
				}
				addUnique(&user.Personality, preferredAnswerValue(a))
			} else {
				addUnique(&user.Personality, answerText)
			}
		case "gender", "genders":
			addUnique(&user.Genders, answerText)
		case "ethnicity", "ethnicities", "race":
			addUnique(&user.Ethnicities, answerText)
		case "religion", "religions":
			addUnique(&user.Religions, answerText)
		case "major", "majors", "dedicatedmajor", "dedicatedmajors":
			addUnique(&user.DedicatedMajors, answerText)
		case "associatedmajor", "associatedmajors":
			addUnique(&user.AssociatedMajors, answerText)
		case "other":
			addUnique(&user.Other, answerText)
		}
	}

	return user
}

// Sort ranks organizations by descending normalized match score.
func Sort(user UserInfo, organizations []Organization) []MatchResult {
	results := make([]MatchResult, 0, len(organizations))
	for _, org := range organizations {
		maxScore := compareUserToOrg(toUserInfo(org), org)
		rawScore := compareUserToOrg(user, org)
		normalized := float32(0)
		if maxScore > 0 {
			normalized = 100 * rawScore / maxScore
		}

		results = append(results, MatchResult{
			Organization:    org,
			RawScore:        rawScore,
			NormalizedScore: normalized,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].NormalizedScore == results[j].NormalizedScore {
			return strings.ToLower(results[i].Organization.Name) < strings.ToLower(results[j].Organization.Name)
		}
		return results[i].NormalizedScore > results[j].NormalizedScore
	})

	return results
}

// Calculates the overall matching score for a user and an organization
func compareUserToOrg(user UserInfo, org Organization) float32 {
	personalityScore := personalityScoring(user, org)
	activityScore := activityScoring(user, org)
	demographicScore := demographicScoring(user, org)
	academicScore := academicScoring(user, org)
	otherScore := otherScoring(user, org)

	finalScore := personalityScore + activityScore + demographicScore + academicScore + otherScore
	if demographicScore == 0 && org.StrictGenders {
		return 0
	}
	return finalScore
}

// Calculates the matching score for personality between a user and an organization
func personalityScoring(user UserInfo, org Organization) float32 {
	if len(org.Personality) == 0 { // check to avoid division by 0
		return 0
	}
	score := overlapCount(user.Personality, org.Personality)
	var overfit float32 = 1
	if len(user.Personality) > score {
		overfit = 1 - 0.01*(float32(len(user.Personality))-float32(score))
	}
	return overfit * personalityWeight * float32(score) / float32(len(org.Personality))
}

// Calculates the matching score for activity interest between a user and an organization
func activityScoring(user UserInfo, org Organization) float32 {
	if len(org.Activities) == 0 { // check to avoid division by 0
		return 0
	}
	score := overlapCount(user.Activities, org.Activities)
	var overfit float32 = 1
	if len(user.Activities) > score {
		overfit = 1 - 0.01*(float32(len(user.Activities))-float32(score))
	}
	return overfit * activityWeight * float32(score) / float32(len(org.Activities))
}

// Calculates the matching score for demographics between a user and an organization
func demographicScoring(user UserInfo, org Organization) float32 {
	totalDemographics := min(1, len(org.Genders)) + min(1, len(org.Ethnicities)) + min(1, len(org.Religions))
	if totalDemographics == 0 { // check to avoid division by 0
		return 0
	}
	genderMatches := 0
	for _, gender := range user.Genders {
		if containsFold(org.Genders, gender) {
			genderMatches = 1
		}
	}
	score := genderMatches

	ethnicityMatches := 0
	for _, ethnicity := range user.Ethnicities {
		if containsFold(org.Ethnicities, ethnicity) {
			ethnicityMatches = 1
		}
	}
	score += ethnicityMatches

	religionMatches := 0
	for _, religion := range user.Religions {
		if containsFold(org.Religions, religion) {
			religionMatches = 1
		}
	}
	score += religionMatches

	if genderMatches == 0 && org.StrictGenders {
		score = 0
	}
	return demographicWeight * float32(score) / float32(totalDemographics)
}

// Calculates the matching score for academic interests between a user and an organization
func academicScoring(user UserInfo, org Organization) float32 {
	for _, major := range user.DedicatedMajors {
		if containsFold(org.DedicatedMajors, major) {
			return academicWeight
		}
	}

	for _, major := range user.DedicatedMajors {
		if containsFold(org.AssociatedMajors, major) {
			return academicWeight / 2
		}
	}

	for _, major := range user.AssociatedMajors { // used only when calculating the max score of an organization
		if containsFold(org.AssociatedMajors, major) {
			return academicWeight / 2
		}
	}

	return 0
}

// Calculates the matching score for miscellaneous questions between a user and an organization
func otherScoring(user UserInfo, org Organization) float32 {
	if len(org.Other) == 0 {
		return 0
	}
	score := overlapCount(user.Other, org.Other)
	return otherWeight * float32(score) / float32(len(org.Other))
}

// Calculates how many values in two string arrays are equal
func overlapCount(left []string, right []string) int {
	count := 0
	for _, a := range left {
		for _, b := range right {
			if strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b)) {
				count++
			}
		}
	}
	return count
}

func containsFold(values []string, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}
	for _, v := range values {
		if strings.EqualFold(strings.TrimSpace(v), target) {
			return true
		}
	}
	return false
}

func addUnique(target *[]string, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}
	if containsFold(*target, trimmed) {
		return
	}
	*target = append(*target, trimmed)
}

func parseBool(raw string) (bool, bool) {
	if raw == "" {
		return false, false
	}
	v, err := strconv.ParseBool(strings.ToLower(strings.TrimSpace(raw)))
	if err != nil {
		return false, false
	}
	return v, true
}

func preferredAnswerValue(a Answer) string {
	if a.Translations != nil {
		if en, ok := a.Translations["en"]; ok && len(en) > 0 {
			if v := strings.TrimSpace(en[0]); v != "" {
				return v
			}
		}
		if term, ok := a.Translations["term"]; ok && len(term) > 0 {
			if v := strings.TrimSpace(term[0]); v != "" {
				return v
			}
		}
		if en, ok := a.Translations["en"]; ok && len(en) > 1 {
			if v := strings.TrimSpace(en[1]); v != "" {
				return v
			}
		}
		if term, ok := a.Translations["term"]; ok && len(term) > 1 {
			if v := strings.TrimSpace(term[1]); v != "" {
				return v
			}
		}
	}
	return strings.TrimSpace(a.AnswerText)
}

func normalizeKey(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	replacer := strings.NewReplacer("_", "", "-", "", " ", "")
	return replacer.Replace(raw)
}

func toUserInfo(org Organization) UserInfo {
	return UserInfo{
		Name:             org.Name,
		Personality:      append([]string{}, org.Personality...),
		Activities:       append([]string{}, org.Activities...),
		Genders:          append([]string{}, org.Genders...),
		Ethnicities:      append([]string{}, org.Ethnicities...),
		Religions:        append([]string{}, org.Religions...),
		DedicatedMajors:  append([]string{}, org.DedicatedMajors...),
		AssociatedMajors: append([]string{}, org.AssociatedMajors...),
		Other:            append([]string{}, org.Other...),
	}
}

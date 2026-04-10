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

func personalityScoring(user UserInfo, org Organization) float32 {
	if len(org.Personality) == 0 {
		return 0
	}
	score := overlapCount(user.Personality, org.Personality)
	return 1.5 * float32(score) / float32(len(org.Personality))
}

func activityScoring(user UserInfo, org Organization) float32 {
	if len(org.Activities) == 0 {
		return 0
	}
	score := overlapCount(user.Activities, org.Activities)
	return 2.25 * float32(score) / float32(len(org.Activities))
}

func demographicScoring(user UserInfo, org Organization) float32 {
	totalDemographics := min(1, len(org.Genders)) + len(org.Ethnicities) + min(1, len(org.Religions))
	if totalDemographics == 0 {
		return 0
	}

	genderMatches := overlapCount(user.Genders, org.Genders)
	if genderMatches == 0 && org.StrictGenders {
		return 0
	}

	score := genderMatches
	score += overlapCount(user.Ethnicities, org.Ethnicities)
	score += overlapCount(user.Religions, org.Religions)

	return 1.25 * float32(score) / float32(totalDemographics)
}

func academicScoring(user UserInfo, org Organization) float32 {
	for _, major := range user.DedicatedMajors {
		if containsFold(org.DedicatedMajors, major) {
			return 1.5
		}
	}

	for _, major := range user.DedicatedMajors {
		if containsFold(org.AssociatedMajors, major) {
			return 0.75
		}
	}

	for _, major := range user.AssociatedMajors {
		if containsFold(org.AssociatedMajors, major) {
			return 0.75
		}
	}

	return 0
}

func otherScoring(user UserInfo, org Organization) float32 {
	if len(org.Other) == 0 {
		return 0
	}
	score := overlapCount(user.Other, org.Other)
	return 2 * float32(score) / float32(len(org.Other))
}

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

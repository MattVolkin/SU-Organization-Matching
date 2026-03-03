package main

// set base score for every club (maybe not)
// var cs_club int = 0
// var the_game int = 0
// var fiber_arts int = 0
// var p4p int = 0
// var su_tabletop int = 0
// var su_tertulias int = 0
// var umsf int = 0
// var pre_dental_society int = 0
// var pre_health_org int = 0
// var cat_partners int = 0
// var kappa_delta_chi int = 0

type organization struct {
	name                string
	personality         []string
	activities          []string
	demographics        []string
	strict_demographics bool
	majors              []string
	is_departmental     bool
	other               []string
}

func main() {

	// hardcoded structs for every organization (for now, will be importing from database, allowing for information to be updated)
	var cs_club organization
	cs_club.name = "Computer Science Club"
	cs_club.personality = []string{"Welcoming", "Hard Working", "Caring", "Nerdy"}
	cs_club.activities = []string{"Board Games", "Movies", "Video Games", "Giving Presentations"}
	cs_club.majors = []string{"Computer Science", "Computational Mathematics"}
	cs_club.is_departmental = true

	// the user's answers will be passed in and put into an organization struct.
	// they will then be compared to and scored with every organization.
	// either:
	// organizations with the highest x scores will be reported
	// organizations with a score above y will be reported
	// certain organizations will have strict requirements
	// as an example, for a sorority, the user must be a woman and answered that they are interested in greek life
}

// Takes a user's answers and organization's information and
// passes them to other functions to get scores for each category
func compare(user organization, org organization) {
	var personality_score float32 = personality_scoring(user, org)
	var activity_score float32 = activity_scoring(user, org)
	var demographic_score float32 = demographic_scoring(user, org)
	var academic_score float32 = academic_scoring(user, org)
	var other_score float32 = other_scoring(user, org)
	// how do we want to combine these?
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the personalities match
func personality_scoring(user organization, org organization) float32 {
	var score float32 = 0
	for i := 0; i < len(user.personality); i++ {
		for j := 0; j < len(org.personality); j++ {
			if user.personality[i] == org.personality[j] {
				score++
			}
		}
	}
	score = score / float32(len(org.personality)) // all orgs should have some personality info so this should never be division by 0
	return score
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the activity interestes match
func activity_scoring(user organization, org organization) float32 {
	return 0
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the demographics match
func demographic_scoring(user organization, org organization) float32 {
	return 0
}

// Takes a user's answers and an organization's information and
// returns a score based on whether the user's major matches
// with the organization, and whether the organization is departmental
func academic_scoring(user organization, org organization) float32 {

	for i := 0; i < len(user.majors); i++ {
		for j := 0; j < len(org.majors); j++ {
			if user.majors[i] == org.majors[j] { // if you match with a single major you get full credit
				if org.is_departmental { // weight major matches with departmental orgs higher
					return 2 // arbitrary value that will be adjusted (weighted higher)
				} else {
					return 1 // arbitrary value that will be adjusted (weighted lower)
				}
			}
		}
	}
	return 0 // no match (how will we deal with orgs that don't have majors listed?)
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the answers to the miscellaneous question match
func other_scoring(user organization, org organization) float32 {
	return 0
}

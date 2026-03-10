package main

import "fmt"

type organization struct {
	name                string
	personality         []string
	activities          []string
	genders             []string
	ethnicities         []string
	religions           []string
	strict_demographics bool
	dedicated_majors    []string
	associated_majors   []string // this field will be blank for users
	other               []string
}

func main() {

	// hardcoded structs for every organization (for now, will be importing from database, allowing for information to be updated)
	var cs_club organization
	cs_club.name = "Computer Science Club"
	cs_club.personality = []string{"Welcoming", "Hard Working", "Caring", "Nerdy", "Fun"}
	cs_club.activities = []string{"Board Games", "Movies", "Video Games", "Giving Presentations", "Trivia", "Guest Speakers", "Group Meals"}
	cs_club.dedicated_majors = []string{"Computer Science", "Computational Mathematics"}
	cs_club.associated_majors = []string{"Biology", "Mathematics"}

	var su_tabletop organization
	su_tabletop.name = "SU Tabletop"
	su_tabletop.personality = []string{"Welcoming", "Caring", "Creative", "Nerdy", "Enthusiastic", "Curious", "Fun"}
	su_tabletop.activities = []string{"Board Games"}

	var p4p organization
	p4p.name = "Pirates for Pride"
	p4p.personality = []string{"Welcoming", "Caring", "Outgoing", "Open Minded", "Enthusiastic", "Social"}
	p4p.activities = []string{"Social Justice", "Arts & Crafts", "Discussion"}
	p4p.genders = []string{"Non Binary", "Other"}
	p4p.other = []string{"Queer"}

	// I am a part of the first three orgs, I am not a part of the ones below

	var exercise_is_medicine organization
	exercise_is_medicine.name = "Exercise is Medicine"
	exercise_is_medicine.personality = []string{"Welcoming", "Hard Working", "Caring", "Collaborative"}
	exercise_is_medicine.activities = []string{"Exercise", "Animal Care"}

	var kdc organization
	kdc.name = "Kappa Delta Chi"
	kdc.personality = []string{"Hard Working", "Caring", "Eager to Learn", "A Leader", "Social"}
	kdc.activities = []string{"Fundraising", "Social Justice", "Retreats"}
	kdc.genders = []string{"Woman"}
	kdc.ethnicities = []string{"Latino"}
	kdc.other = []string{"Greek Life"}

	var classics_club organization
	classics_club.name = "Classics Club"
	classics_club.personality = []string{"Outgoing", "Eager to Learn", "Enthusiastic", "Social", "Fun"}
	classics_club.activities = []string{"Board Games", "Arts & Crafts", "Giving Presentations"}
	classics_club.dedicated_majors = []string{"Classics"}

	// using info I would put in for the sake of testing
	var tanner organization
	tanner.name = "Tanner Klein"
	tanner.personality = []string{"Welcoming", "Hard Working", "Caring", "Creative", "Open Minded", "Eager to Learn", "Nerdy", "A Leader", "Enthusiastic", "Fun"}
	tanner.activities = []string{"Board Games", "Video Games", "Arts & Crafts", "Giving Presentations"}
	tanner.genders = []string{"Man"}
	tanner.ethnicities = []string{"White"}
	tanner.religions = []string{"No Religion"}
	tanner.dedicated_majors = []string{"Mathematics", "Computer Science"}
	tanner.other = []string{"Queer"}

	var matt organization
	matt.name = "Matthew Volkin"
	matt.personality = []string{"Welcoming", "Hard Working", "Caring", "Open Minded", "Eager to Learn", "Nerdy", "Enthusiastic", "Collaborative", "Curious", "Social", "Fun"}
	matt.activities = []string{"Board Games", "Movies", "Video Games", "Trivia", "Study Groups", "Group Meals"}
	matt.genders = []string{"Man"}
	matt.ethnicities = []string{"White"}
	matt.religions = []string{"No Religion"}
	matt.dedicated_majors = []string{"Computer Science"}
	matt.other = []string{"Disability"}

	var aidan organization
	aidan.name = "Aidan Balakrishnan"
	aidan.personality = []string{"Hard Working", "Caring", "Creative", "Eager to Learn", "Nerdy", "Collaborative", "Social", "Fun"}
	aidan.activities = []string{"Fundraising", "Social Justice", "Board Games", "Movies", "Video Games", "Arts & Crafts", "Music", "Professional Development", "Trivia", "Study Groups"}
	aidan.genders = []string{"Man"}
	aidan.ethnicities = []string{"Asian"}
	aidan.religions = []string{"No Religion"}
	aidan.dedicated_majors = []string{"Computer Science", "Theatre"}
	aidan.other = []string{"Queer"}

	var ben organization
	ben.name = "Benjamin McKallip"
	ben.personality = []string{"Caring", "Creative", "Open Minded", "Eager to Learn", "Nerdy", "Enthusiastic", "Curious", "Social", "Fun"}
	ben.activities = []string{"Board Games", "Movies", "Video Games", "Music", "Trivia", "Group Meals", "Discussion"}
	ben.genders = []string{"Man"}
	ben.ethnicities = []string{"White"}
	ben.religions = []string{"No Religion"}
	ben.dedicated_majors = []string{"Physics", "Computer Science"}

	compare(tanner, cs_club)
	compare(matt, cs_club)
	compare(aidan, cs_club)
	compare(ben, cs_club)
	compare(cs_club, cs_club)

	compare(tanner, su_tabletop)
	compare(matt, su_tabletop)
	compare(aidan, su_tabletop)
	compare(ben, su_tabletop)
	compare(su_tabletop, su_tabletop)

	compare(tanner, p4p)
	compare(matt, p4p)
	compare(aidan, p4p)
	compare(ben, p4p)
	compare(p4p, p4p)

	compare(tanner, exercise_is_medicine)
	compare(matt, exercise_is_medicine)
	compare(aidan, exercise_is_medicine)
	compare(ben, exercise_is_medicine)
	compare(exercise_is_medicine, exercise_is_medicine)

	compare(tanner, kdc)
	compare(matt, kdc)
	compare(aidan, kdc)
	compare(ben, kdc)
	compare(kdc, kdc)

	compare(tanner, classics_club)
	compare(matt, classics_club)
	compare(aidan, classics_club)
	compare(ben, classics_club)
	compare(classics_club, classics_club)

	// trying out comparing orgs to each other
	// we might be able to make a "similar orgs" button
	// compare(cs_club, su_tabletop)
	// compare(cs_club, p4p)
	// compare(p4p, cs_club) // order matters right now (it probably shouldn't though)
	// compare(su_tabletop, p4p)

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
	// how do we want to combine these? (I'll just add them for now)
	var final_score float32 = personality_score + activity_score + demographic_score + academic_score + other_score
	if demographic_score == 0 && org.strict_demographics {
		final_score = 0
	}
	// testing print statement
	fmt.Println("The matching score for", user.name, "and", org.name, "is", final_score)

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
	var score float32 = 0
	for i := 0; i < len(user.activities); i++ {
		for j := 0; j < len(org.activities); j++ {
			if user.activities[i] == org.activities[j] {
				score++
			}
		}
	}
	score = 1.5 * score / float32(len(org.activities)) // all orgs should have some personality info so this should never be division by 0
	return score
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the demographics match
func demographic_scoring(user organization, org organization) float32 { // TODO: Rework this it doesn't work with the multiple choice questions
	var total_demographics int = len(org.genders) + len(org.ethnicities) + len(org.religions)
	if total_demographics == 0 { // remove possibility of dividing by 0
		return 0
	}
	var score float32 = 0
	for i := 0; i < len(user.genders); i++ {
		for j := 0; j < len(org.genders); j++ {
			if user.genders[i] == org.genders[j] {
				score++
			}
		}
	}
	for i := 0; i < len(user.ethnicities); i++ {
		for j := 0; j < len(org.ethnicities); j++ {
			if user.ethnicities[i] == org.ethnicities[j] {
				score++
			}
		}
	}
	for i := 0; i < len(user.religions); i++ {
		for j := 0; j < len(org.religions); j++ {
			if user.religions[i] == org.religions[j] {
				score++
			}
		}
	}
	score = score / float32(total_demographics)
	return score
}

// Takes a user's answers and an organization's information and
// returns a score based on whether the user's major matches
// with the organization. Some majors are weighted higher or lower
// depending on how closely they match with the organization
func academic_scoring(user organization, org organization) float32 {
	var score float32 = 0
	for i := 0; i < len(user.dedicated_majors); i++ {
		for j := 0; j < len(org.dedicated_majors); j++ {
			if user.dedicated_majors[i] == org.dedicated_majors[j] { // if you match with a single major you get full credit
				return 2 // arbitrary value that will be adjusted (weighted higher)
			}
		}
		for j := 0; j < len(org.associated_majors); j++ {
			if user.dedicated_majors[i] == org.associated_majors[j] {
				score = 1 // arbitrary value that will be adjusted (weighted lower)
			}
		}

	}
	return score // no match (how will we deal with orgs that don't have majors listed?)
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the answers to the miscellaneous questions match
// these questions are usually very specific to certain orgs (interest in Greek Life, LGBTQ)
func other_scoring(user organization, org organization) float32 {
	if len(org.other) == 0 { // remove possibility of dividing by 0
		return 0
	}
	var score float32 = 0
	for i := 0; i < len(user.other); i++ {
		for j := 0; j < len(org.other); j++ {
			if user.other[i] == org.other[j] {
				score++
			}
		}
	}
	score = 2 * score / float32(len(org.other))
	return score
}

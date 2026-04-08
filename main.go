package main

import (
	"fmt"
	//"./server"
)

// organization defines either an organization or a user with various fields
// corresponding with the questions asked on the quiz.
type Organization struct {
	name              string
	personality       []string
	activities        []string
	genders           []string
	ethnicities       []string
	religions         []string
	strict_genders    bool
	dedicated_majors  []string
	associated_majors []string
	other             []string
	max_score         float32
}

func main() {

	// // Open the Postgres connection used by Ent queries and mutations.
	// dbClient := GetDatabaseClient('default')
	// _, orgs, err := dbClient.Query().FetchAllOrganizationsForSorting()

	// hardcoded structs for every organization (for now, will be importing from database, allowing for information to be updated)
	clubs := [3]Organization{
		{
			name:              "Computer Science Club",
			personality:       []string{"Welcoming", "Hard Working", "Caring", "Nerdy", "Fun"},
			activities:        []string{"Board Games", "Movies", "Video Games", "Giving Presentations", "Trivia", "Guest Speakers", "Group Meals"},
			dedicated_majors:  []string{"Computer Science", "Computational Mathematics"},
			associated_majors: []string{"Biology", "Mathematics"},
		},
		{
			name:        "SU Tabletop",
			personality: []string{"Welcoming", "Caring", "Creative", "Nerdy", "Enthusiastic", "Curious", "Fun"},
			activities:  []string{"Board Games"},
		},
		{
			name:        "Pirates for Pride",
			personality: []string{"Welcoming", "Caring", "Outgoing", "Open Minded", "Enthusiastic", "Social"},
			activities:  []string{"Social Justice", "Arts & Crafts", "Discussion"},
			genders:     []string{"Non Binary", "Other"},
			other:       []string{"Queer"},
		},
	}

	for i := 0; i < len(clubs); i++ {
		clubs[i].max_score = compare(clubs[i], clubs[i])
	}

	// I am a part of the first three orgs, I am not a part of the ones below

	// var exercise_is_medicine Organization
	// exercise_is_medicine.name = "Exercise is Medicine"
	// exercise_is_medicine.personality = []string{"Welcoming", "Hard Working", "Caring", "Collaborative"}
	// exercise_is_medicine.activities = []string{"Exercise", "Animal Care"}
	// exercise_is_medicine.max_score = compare(exercise_is_medicine, exercise_is_medicine)

	// var kdc Organization
	// kdc.name = "Kappa Delta Chi"
	// kdc.personality = []string{"Hard Working", "Caring", "Eager to Learn", "Leader", "Social"}
	// kdc.activities = []string{"Fundraising", "Social Justice", "Retreats"}
	// kdc.genders = []string{"Woman"}
	// kdc.strict_genders = true
	// kdc.ethnicities = []string{"Latino"}
	// kdc.other = []string{"Greek Life"}
	// kdc.max_score = compare(kdc, kdc)

	// var classics_club Organization
	// classics_club.name = "Classics Club"
	// classics_club.personality = []string{"Outgoing", "Eager to Learn", "Enthusiastic", "Social", "Fun"}
	// classics_club.activities = []string{"Board Games", "Arts & Crafts", "Giving Presentations"}
	// classics_club.dedicated_majors = []string{"Classics"}
	// classics_club.max_score = compare(classics_club, classics_club)

	// var kappa_sigma Organization
	// kappa_sigma.name = "Kappa Sigma"
	// kappa_sigma.personality = []string{"Creative", "Leader", "Social"}
	// kappa_sigma.activities = []string{"Fundraising", "Group Meals"}
	// kappa_sigma.genders = []string{"Man"}
	// kappa_sigma.strict_genders = true
	// kappa_sigma.other = []string{"Greek Life"}
	// kappa_sigma.max_score = compare(kappa_sigma, kappa_sigma)

	// var mask_wig Organization
	// mask_wig.name = "Mask and Wig"
	// mask_wig.personality = []string{"Welcoming", "Creative", "Open Minded", "Collaborative", "Social", "Fun"}
	// mask_wig.activities = []string{"Board Games", "Movies", "Arts & Crafts", "Writing"}
	// mask_wig.dedicated_majors = []string{"Theatre"}
	// mask_wig.max_score = compare(mask_wig, mask_wig)

	// var fiber_arts Organization
	// fiber_arts.name = "Fiber Arts Club"
	// fiber_arts.personality = []string{"Welcoming", "Hard Working", "Creative", "Eager to Learn", "Enthusiastic"}
	// fiber_arts.activities = []string{"Arts & Crafts", "Discussion"}
	// fiber_arts.associated_majors = []string{"Art"}
	// fiber_arts.max_score = compare(fiber_arts, fiber_arts)

	// var megaphone Organization
	// megaphone.name = "Megaphone"
	// megaphone.personality = []string{"Hard Working", "Outgoing", "Enthusiastic"}
	// megaphone.activities = []string{"Writing"}
	// megaphone.associated_majors = []string{"English"}
	// megaphone.max_score = compare(megaphone, megaphone)

	// var fca Organization
	// fca.name = "Fellowship of Christian Athletes"
	// fca.personality = []string{"Caring", "Eager to Learn", "Collaborative", "Curious"}
	// fca.activities = []string{"Board Games", "Exercise", "Guest Speakers", "Sports"}
	// fca.religions = []string{"Protestantism", "Catholocism"}
	// fca.max_score = compare(fca, fca)

	// using info I would put in for the sake of testing
	users := [6]Organization{
		{
			name:             "Tanner Klein",
			personality:      []string{"Welcoming", "Hard Working", "Caring", "Creative", "Open Minded", "Eager to Learn", "Nerdy", "Leader", "Enthusiastic", "Fun"},
			activities:       []string{"Board Games", "Video Games", "Arts & Crafts", "Giving Presentations"},
			genders:          []string{"Man"},
			ethnicities:      []string{"White"},
			religions:        []string{"No Religion"},
			dedicated_majors: []string{"Mathematics", "Computer Science"},
			other:            []string{"Queer"},
		},
		{
			name:             "Matthew Volkin",
			personality:      []string{"Welcoming", "Hard Working", "Caring", "Open Minded", "Eager to Learn", "Nerdy", "Enthusiastic", "Collaborative", "Curious", "Social", "Fun"},
			activities:       []string{"Board Games", "Movies", "Video Games", "Trivia", "Study Groups", "Group Meals"},
			genders:          []string{"Man"},
			ethnicities:      []string{"White"},
			religions:        []string{"No Religion"},
			dedicated_majors: []string{"Computer Science"},
			other:            []string{"Disability"},
		},
		{
			name:             "Aidan Balakrishnan",
			personality:      []string{"Hard Working", "Caring", "Creative", "Eager to Learn", "Nerdy", "Collaborative", "Social", "Fun"},
			activities:       []string{"Fundraising", "Social Justice", "Board Games", "Movies", "Video Games", "Arts & Crafts", "Music", "Professional Development", "Trivia", "Study Groups"},
			genders:          []string{"Man"},
			ethnicities:      []string{"Asian"},
			religions:        []string{"No Religion"},
			dedicated_majors: []string{"Computer Science", "Theatre"},
			other:            []string{"Queer"},
		},
		{
			name:             "Benjamin McKallip",
			personality:      []string{"Caring", "Creative", "Open Minded", "Eager to Learn", "Nerdy", "Enthusiastic", "Curious", "Social", "Fun"},
			activities:       []string{"Board Games", "Movies", "Video Games", "Music", "Trivia", "Group Meals", "Discussion"},
			genders:          []string{"Man"},
			ethnicities:      []string{"White"},
			religions:        []string{"No Religion"},
			dedicated_majors: []string{"Physics", "Computer Science"},
		},
		{
			// eventually we want a specific case that would bring this user to a different type of results page telling them to take the test again.
			name:             "All",
			personality:      []string{"Welcoming", "Hard Working", "Caring", "Creative", "Outgoing", "Open Minded", "Eager to Learn", "Confident", "Nerdy", "Leader", "Enthusiastic", "Collaborative", "Curious", "Organized", "Social", "Fun"},
			activities:       []string{"Fundraising", "Social Justice", "Retreats", "Dance", "Board Games", "Movies", "Video Games", "Arts & Crafts", "Music", "Exercise", "Writing", "Professional Development", "Caring for Animals", "Giving Presentations", "Trivia", "Literary Analysis", "Study Groups", "Guest Speakers", "Group Meals", "Discussion", "Sports"},
			genders:          []string{"Prefer not to say"},
			ethnicities:      []string{"Prefer not to say"},
			religions:        []string{"Prefer not to say"},
			dedicated_majors: []string{"Undecided"},
			other:            []string{"Queer", "Greek Life", "Disability"},
		},
		{
			// eventually we want a specific case that would bring this user to a different type of results page telling them to take the test again.
			name:             "None",
			genders:          []string{"Prefer not to say"},
			ethnicities:      []string{"Prefer not to say"},
			religions:        []string{"Prefer not to say"},
			dedicated_majors: []string{"Undecided"},
		},
	}

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], cs_club)
	// 	var normalized_score float32 = 100 * score / cs_club.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, cs_club.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], su_tabletop)
	// 	var normalized_score float32 = 100 * score / su_tabletop.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, su_tabletop.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], p4p)
	// 	var normalized_score float32 = 100 * score / p4p.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, p4p.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], kdc)
	// 	var normalized_score float32 = 100 * score / kdc.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, kdc.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], exercise_is_medicine)
	// 	var normalized_score float32 = 100 * score / exercise_is_medicine.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, exercise_is_medicine.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], classics_club)
	// 	var normalized_score float32 = 100 * score / classics_club.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, classics_club.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], kappa_sigma)
	// 	var normalized_score float32 = 100 * score / kappa_sigma.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, kappa_sigma.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], mask_wig)
	// 	var normalized_score float32 = 100 * score / mask_wig.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, mask_wig.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], fiber_arts)
	// 	var normalized_score float32 = 100 * score / fiber_arts.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, fiber_arts.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], megaphone)
	// 	var normalized_score float32 = 100 * score / megaphone.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, megaphone.name, normalized_score)
	// }

	// for i := 0; i < len(users); i++ {
	// 	var score float32 = compare(users[i], fca)
	// 	var normalized_score float32 = 100 * score / fca.max_score
	// 	fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, fca.name, normalized_score)
	// }

	for i := 0; i < len(users); i++ {
		for j := 0; j < len(clubs); j++ {
			var score float32 = compare(users[i], clubs[j])
			var normalized_score float32 = 100 * score / clubs[j].max_score
			fmt.Printf("The matching score for %s and %s is %f\n", users[i].name, clubs[j].name, normalized_score)
		}
	}
}

// Takes a user's answers and organization's information and
// passes them to other functions to get scores for each category
func compare(user Organization, org Organization) float32 {
	var personality_score float32 = personality_scoring(user, org)
	var activity_score float32 = activity_scoring(user, org)
	var demographic_score float32 = demographic_scoring(user, org)
	var academic_score float32 = academic_scoring(user, org)
	var other_score float32 = other_scoring(user, org)
	// how do we want to combine these? (I'll just add them for now)
	var final_score float32 = personality_score + activity_score + demographic_score + academic_score + other_score
	if demographic_score == 0 && org.strict_genders {
		final_score = 0
	}
	return final_score
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the personalities match
func personality_scoring(user Organization, org Organization) float32 {
	var score float32 = 0
	for i := 0; i < len(user.personality); i++ {
		for j := 0; j < len(org.personality); j++ {
			if user.personality[i] == org.personality[j] {
				score++
			}
		}
	}
	score = 1.5 * score / float32(len(org.personality)) // all orgs should have some personality info so this should never be division by 0
	return score
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the activity interests match
func activity_scoring(user Organization, org Organization) float32 {
	var score float32 = 0
	for i := 0; i < len(user.activities); i++ {
		for j := 0; j < len(org.activities); j++ {
			if user.activities[i] == org.activities[j] {
				score++
			}
		}
	}
	score = 2.25 * score / float32(len(org.activities))
	return score
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the demographics match
func demographic_scoring(user Organization, org Organization) float32 { // TODO: Rework this it doesn't work with the multiple choice questions
	var total_demographics int = min(1, len(org.genders)) + len(org.ethnicities) + min(1, len(org.religions))
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
	if score == 0 && org.strict_genders { // used mostly for non-coed greek orgs
		return 0
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
	score = 1.25 * score / float32(total_demographics)
	return score
}

// Takes a user's answers and an organization's information and
// returns a score based on whether the user's major matches
// with the organization. Some majors are weighted higher or lower
// depending on how closely they match with the organization
func academic_scoring(user Organization, org Organization) float32 {
	var score float32 = 0
	for i := 0; i < len(user.dedicated_majors); i++ {
		for j := 0; j < len(org.dedicated_majors); j++ {
			if user.dedicated_majors[i] == org.dedicated_majors[j] { // if you match with a single major you get full credit
				return 1.5 // arbitrary value that will be adjusted (weighted higher)
			}
		}
		for j := 0; j < len(org.associated_majors); j++ {
			if user.dedicated_majors[i] == org.associated_majors[j] {
				score = 0.75 // arbitrary value that will be adjusted (weighted lower)
			}

		}

	}
	for i := 0; i < len(user.associated_majors); i++ { // this case only occurs when calculating an organization's max score
		for j := 0; j < len(org.associated_majors); j++ {
			if user.associated_majors[i] == org.associated_majors[j] {
				score = 0.75 // arbitrary value that will be adjusted (weighted lower)
			}
		}
	}
	return score // no match (how will we deal with orgs that don't have majors listed?)
}

// Takes a user's answers and an organization's information and
// returns a score based on how closely the answers to the miscellaneous questions match
// these questions are usually very specific to certain orgs (interest in Greek Life, LGBTQ)
func other_scoring(user Organization, org Organization) float32 {
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

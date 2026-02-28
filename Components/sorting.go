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
	personality  []string
	activities   []string
	demographics []string
	majors       []string
	other        []string
}

func main() {

	// hardcoded structs for every organization (for now, will be importing from database, allowing for information to be updated)
	var cs_club organization
	cs_club.personality = []string{"Welcoming", "Hard Working", "Caring", "Nerdy"}
	cs_club.activities = []string{"Board Games", "Movies", "Video Games"}
	cs_club.majors = []string{"Computer Science", "Computational Mathematics"}

	// the user's answers will be passed in and put into an organization struct.
	// they will then be compared to and scored with every organization.
	// either:
	// organizations with the highest x scores will be reported
	// organizations with a score above y will be reported
	// certain organizations will have strict requirements
	// as an example, for a sorority, the user must be a woman and answered that they are interested in greek life
}

func compare() {

}

func personality_scoring() {

}

func activity_scoring() {

}

func demographic_scoring() {

}

func academic_scoring() {

}

func other_scoring() {

}

func return_scores() {

}

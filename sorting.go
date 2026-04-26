package main

import (
	"fmt"

	"su-organization-matching/matching"
)

func main() {
	clubs := []matching.Organization{
		{
			ID:               1,
			Name:             "Computer Science Club",
			Personality:      []string{"Welcoming", "Hard Working", "Caring", "Nerdy", "Fun"},
			Activities:       []string{"Board Games", "Movies", "Video Games", "Giving Presentations", "Trivia", "Guest Speakers", "Group Meals"},
			DedicatedMajors:  []string{"Computer Science", "Computational Mathematics"},
			AssociatedMajors: []string{"Biology", "Mathematics"},
		},
		{
			ID:          2,
			Name:        "SU Tabletop",
			Personality: []string{"Welcoming", "Caring", "Creative", "Nerdy", "Enthusiastic", "Curious", "Fun"},
			Activities:  []string{"Board Games"},
		},
		{
			ID:            3,
			Name:          "Pirates for Pride",
			Personality:   []string{"Welcoming", "Caring", "Outgoing", "Open Minded", "Enthusiastic", "Social"},
			Activities:    []string{"Social Justice", "Arts & Crafts", "Discussion", "Giving Presentations"},
			Genders:       []string{"Non Binary", "Other"},
			Other:         []string{"LGBTQ"},
			StrictGenders: false,
		},
		{
			ID:            4,
			Name:          "Kappa Alpha",
			Personality:   []string{"Hard Working", "Leader", "Social"},
			Activities:    []string{"Fundraising", "Retreats", "Professional Development", "Group Meals"},
			Genders:       []string{"Man"},
			Religions:     []string{"Protestantism", "Catholocism"},
			Other:         []string{"Greek Life"},
			StrictGenders: true,
		},
		{
			ID:          5,
			Name:        "Asian Student Association",
			Personality: []string{"Welcoming", "Open Minded", "Curious", "Social"},
			Activities:  []string{"Movies", "Guest Speakers", "Group Meals"},
			Ethnicities: []string{"Asian", "Native Hawaiian or Pacific Islander"},
		},
		{
			ID:            6,
			Name:          "Kappa Sigma",
			Personality:   []string{"Creative", "Leader", "Social"},
			Activities:    []string{"Fundraising", "Group Meals"},
			Genders:       []string{"Man"},
			Other:         []string{"Greek Life"},
			StrictGenders: false,
		},
	}

	userInfo := matching.UserInfo{
		Name:            "Example User",
		Personality:     []string{"Welcoming", "Open Minded", "Enthusiastic", "Fun"},
		Activities:      []string{"Arts & Crafts"},
		Genders:         []string{"Woman"},
		Ethnicities:     []string{"Asian"},
		Religions:       []string{"Catholocism"},
		DedicatedMajors: []string{"Computer Science"},
		Other:           []string{"LGBTQ", "Greek Life", "Disabilities"},
	}

	results := matching.Sort(userInfo, clubs)
	for idx, result := range results {
		fmt.Printf("%d. %s (%.2f%%)\n", idx+1, result.Organization.Name, result.NormalizedScore)
	}
}

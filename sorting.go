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
			Activities:    []string{"Social Justice", "Arts & Crafts", "Discussion"},
			Genders:       []string{"Non Binary", "Other"},
			Other:         []string{"LGBTQ"},
			StrictGenders: false,
		},
	}

	userInfo := matching.UserInfo{
		Name:            "Example User",
		Personality:     []string{"Welcoming", "Hard Working", "Caring", "Creative", "Nerdy", "Fun"},
		Activities:      []string{"Board Games", "Video Games", "Giving Presentations"},
		Genders:         []string{"Man"},
		Ethnicities:     []string{"White"},
		Religions:       []string{"No Religion"},
		DedicatedMajors: []string{"Computer Science"},
		Other:           []string{"LGBTQ"},
	}

	results := matching.Sort(userInfo, clubs)
	for idx, result := range results {
		fmt.Printf("%d. %s (%.2f%%)\n", idx+1, result.Organization.Name, result.NormalizedScore)
	}
}

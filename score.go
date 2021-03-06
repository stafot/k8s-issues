package main

import (
	"fmt"
	"strings"
)

type textSource struct {
	name       string
	multiplier int
	content    string
}

type scorePoint struct {
	reason string
	points int
}

type sigScore struct {
	scoreItems []scorePoint
	scoreTotal int
}

// Generate a map of score data for every SIG, for an issue.
func getScoresForSigs(issue Issue) map[string]sigScore {
	titleData := textSource{
		name:       "title",
		multiplier: 3,
		content:    issue.Title,
	}
	bodyData := textSource{
		name:       "body",
		multiplier: 1,
		content:    issue.Body,
	}

	var sigScores = make(map[string]sigScore)
	for _, sigDetails := range allSigs {
		scoreItems := scoreForSig(sigDetails, []textSource{titleData, bodyData})
		scoreTotal := 0
		for _, lineItem := range scoreItems {
			scoreTotal += lineItem.points
		}
		sigScores[sigDetails.name] = sigScore{
			scoreItems: scoreItems,
			scoreTotal: scoreTotal,
		}
	}

	return sigScores
}

// Calculate the score details for a single SIG and issue.
func scoreForSig(sig Sig, sources []textSource) []scorePoint {
	var score []scorePoint = nil

	for _, source := range sources {
		for _, keyword := range sig.strongMatches {
			if count := strings.Count(strings.ToLower(source.content), keyword); count > 0 {
				points := 3 * source.multiplier * count
				score = append(score, scorePoint{
					reason: fmt.Sprintf("%s was a strong match in %s", keyword, source.name),
					points: points,
				})
			}
		}
		// TODO deduplicate
		for _, keyword := range sig.weakMatches {
			if count := strings.Count(strings.ToLower(source.content), keyword); count > 0 {
				points := 1 * source.multiplier * count
				score = append(score, scorePoint{
					reason: fmt.Sprintf("%s was a weak match in %s", keyword, source.name),
					points: points,
				})
			}
		}
	}

	return score
}

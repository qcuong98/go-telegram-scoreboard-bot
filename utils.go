package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func showScores(scores map[string]int, base int, inputN bool) string {
	type scoreEntry struct {
		Username string
		Rank     int
		Score    int
	}

	var entries []scoreEntry
	for username, score := range scores {
		entries = append(entries, scoreEntry{Username: username, Score: score})
	}

	// Sort the slice in descending order based on the score
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score ||
			(entries[i].Score == entries[j].Score && entries[i].Username > entries[j].Username)
	})

	rank := 0
	for idx, entry := range entries {
		if idx == 0 || entry.Score != entries[idx-1].Score {
			rank = idx + 1
		}
		entries[idx].Rank = rank
	}

	msg := ""
	msg = msg + "Scoreboard:\n"
	for idx, entry := range entries {
		x := entry.Score - base
		emoji := ""
		if entry.Rank == 1 {
			emoji = "🥇"
		} else if entry.Rank == entries[len(entries)-1].Rank {
			emoji = "🌚"
		} else if inputN {
			if x > 0 {
				emoji = "🟢"
			} else if x == 0 {
				emoji = "🟡"
			} else if x < 0 {
				emoji = "🔴"
			}
		} else {
			emoji = "⚪️"
		}

		msg = msg + fmt.Sprintf("%v #%d. %s: %d\n", emoji, entry.Rank, entry.Username, entry.Score)
		if idx < len(entries)-1 {
			x := entry.Score - base
			y := entries[idx+1].Score - base
			if x*y <= 0 && (x != 0 || y != 0) {
				msg = msg + fmt.Sprintf("-----\n")
			}
		}
	}
	return msg
}

func addScores(scores map[string]int, users []string, increment int) {
	for _, user := range users {
		if score, ok := scores[user]; ok {
			scores[user] = score + increment
		} else {
			scores[user] = increment
		}
	}
}

func subScores(scores map[string]int, users []string, decrement int) {
	for _, user := range users {
		if score, ok := scores[user]; ok {
			scores[user] = score - decrement
		} else {
			scores[user] = -decrement
		}
	}
}

func removeUsers(scores map[string]int, users []string) []bool {
	var exist []bool
	for _, user := range users {
		_, ok := scores[user]
		exist = append(exist, ok)

		delete(scores, user)
	}
	return exist
}

func initializeScores(scores map[string]int, users []string, presetScore int) []bool {
	var exist []bool
	for _, user := range users {
		_, ok := scores[user]
		exist = append(exist, ok)

		scores[user] = presetScore
	}
	return exist
}

func parseInput(input []string, defaultValue int) ([]string, int, bool) {
	// Extract usernames from input
	usernames := []string{}
	nUsers := 0
	inputN := false

	// Extract and parse the last item as the score
	scoreStr := input[len(input)-1]
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		score = defaultValue
		nUsers = len(input)
	} else {
		nUsers = len(input) - 1
		inputN = true
	}

	for i := 0; i < nUsers; i++ {
		st := strings.ToLower(input[i])
		if len(st) > 0 && st[0] == '@' {
			st = st[1:]
		}
		if len(st) > 0 && !(st[0] >= '0' && st[0] <= '9') {
			usernames = append(usernames, st)
		} else {
			score = defaultValue
			usernames = append(usernames, "_"+st)
		}
	}

	return usernames, score, inputN
}

func sumScores(scores map[string]int) int {
	sum := 0
	for _, score := range scores {
		sum += score
	}
	return sum
}

func contains(slice []int64, item int64) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func diffMaps(map1, map2 map[string]int) string {
	msg := ""
	for key, val := range map1 {
		if map2Val, ok := map2[key]; !ok {
			msg = msg + fmt.Sprintf("%v: %v -> NONE\n", key, val)
		} else if map2Val != val {
			msg = msg + fmt.Sprintf("%v: %v -> %v\n", key, val, map2Val)
		}
	}
	for key, val := range map2 {
		if _, ok := map1[key]; !ok {
			msg = msg + fmt.Sprintf("%v: NONE -> %v\n", key, val)
		}
	}
	if len(msg) == 0 {
		msg = "No changes"
	}
	return msg
}

func cloneMap(map1 map[string]int) map[string]int {
	map2 := make(map[string]int)
	for key, val := range map1 {
		map2[key] = val
	}
	return map2
}

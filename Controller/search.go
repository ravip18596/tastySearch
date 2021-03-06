package Controller

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"
	"tastySearch/Models"
)

func calculateScoreofDocuments(words []string) []Models.ScoreArr {
	scoreArr := make([]Models.ScoreArr, len(Models.Documents))
	scoreChan := make(chan Models.ScoreArr)
	for i := 0; i < len(Models.Documents); i++ {
		go func(scoreChan chan<- Models.ScoreArr, i int) {
			count := 0.
			wordsMAtched := make([]string, 0)
			for j := 0; j < len(words); j++ {
				if Models.Documents[i].Trie.Has(strings.ToLower(words[j])) {
					wordsMAtched = append(wordsMAtched, words[j])
					count++
				}
			}
			score := count / float64(len(words))
			scoreChan <- Models.ScoreArr{
				ReviewScore:   Models.Documents[i].ReviewByScore,
				Score:         score,
				DocumentIndex: i,
				WordsMatched:  wordsMAtched,
			}
		}(scoreChan, i)
	}
	for i := 0; i < len(Models.Documents); i++ {
		temp := <-scoreChan
		scoreArr = append(scoreArr, temp)
	}
	return scoreArr
}

func SearchDocument(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("queries")
	queryWords := strings.Split(query, ",")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if len(queryWords) <= 0 {
		json.NewEncoder(w).Encode(Models.ErrorResponse{Error: Models.EmptyQuery})
		return
	}
	scores := calculateScoreofDocuments(queryWords)
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Score == scores[j].Score {
			return scores[i].ReviewScore > scores[j].ReviewScore
		} else {
			return scores[i].Score > scores[j].Score
		}
	})
	response := make([]Models.ResponseDocument, 0)
	for i := 0; i < 20; i++ {
		logrus.Info("Rank ", i, " Document Score - ", scores[i].Score)
		Models.ResponseDocuments[scores[i].DocumentIndex].WordsMatched = scores[i].WordsMatched
		Models.ResponseDocuments[scores[i].DocumentIndex].DocumentScore = scores[i].Score
		response = append(response, Models.ResponseDocuments[scores[i].DocumentIndex])
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Error marshallsing response")
	}
}

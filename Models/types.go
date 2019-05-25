package Models

import (
	"github.com/fvbock/trie"
	"regexp"
)

const(
	Port = "8081"
	ThresholdDocument = 100000
	DText = "review/text: .*"
	DSummary = "review/summary: .*"
	DScore   = "review/score: .*"
	DProductID = "product/productId: .*"
	DUserID = "review/userId: .*"
	DProfileName = "review/profileName: .*"
	DHelpfulness = "review/helpfulness: .*"
	DReviewTime  = "review/time: .*"
)

var (
	DTextRegex = regexp.MustCompile(DText)
	DSummaryRegex = regexp.MustCompile(DSummary)
	DScoreRegex = regexp.MustCompile(DScore)
	DProductIDRegex = regexp.MustCompile(DProductID)
	DUserIDRegex = regexp.MustCompile(DUserID)
	DProfileNameRegex = regexp.MustCompile(DProfileName)
	DHelpfulRegex = regexp.MustCompile(DHelpfulness)
	DTimeRegex = regexp.MustCompile(DReviewTime)

	Documents []Document

)
type HeartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type Document struct {
	ReviewByScore float64
	Trie          *trie.Trie
}


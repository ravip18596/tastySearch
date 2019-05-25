package Models

const (
	EmptyQuery = "Query string is missing in request parameter"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ResponseDocuments []ResponseDocument
)

type ScoreArr struct {
	ReviewScore   float64
	Score         float64
	WordsMatched  []string
	DocumentIndex int
}

type ResponseDocument struct {
	ProductID    string   `json:"product/productId"`
	UserID       string   `json:"review/userId"`
	ProfileName  string   `json:"review/profileName"`
	Helpful      string   `json:"review/helpfulness"`
	Score        string   `json:"review/score"`
	Time         string   `json:"review/time"`
	Text         string   `json:"review/text"`
	Summary      string   `json:"review/summary"`
	WordsMatched []string `json:"tokens_matched"`
	DocumentScore float64 `json:"document_score"`
}

package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tastySearch/Models"
	"io/ioutil"
	"strings"
	"strconv"
	"github.com/fvbock/trie"
	"tastySearch/Controller"
	"os"
)


// URL : /
func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Models.HeartbeatResponse{Status: "OK", Code: 200})
}

//Middleware should be able to handle panic
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//Middleware should be able to handle panic

				log.WithFields(log.Fields{
					"URL": r.RequestURI,
				}).Errorf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func ParseString(str string) []string{
	str = strings.ToLower(str)
	words := make([]string,0)
	word := ""
	for i:=0;i<len(str);i++{
		if str[i]>='a' && str[i]<='z' || str[i] == '\''{
			word += string(str[i])
		}else{
			word = strings.TrimSpace(word)
			if word != "" {
				words = append(words, word)
			}
			word = ""
		}
	}

	return words
}

func init(){
	filename := os.Getenv("FILE")
	if filename == "" {
		log.Info("Empty filename.")
		return
	}
	fileByte, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":err,
		}).Error("Error while reading foods.txt")
		return
	}
	log.Info("Parsing Data. Please wait for few seconds...")
	Models.Documents = make([]Models.Document,0)
	productId := Models.DProductIDRegex.FindAll(fileByte,-1)
	userID := Models.DUserIDRegex.FindAll(fileByte,-1)
	profileName := Models.DProfileNameRegex.FindAll(fileByte,-1)
	helpful := Models.DHelpfulRegex.FindAll(fileByte,-1)
	reviewScore := Models.DScoreRegex.FindAll(fileByte,-1)
	time := Models.DTimeRegex.FindAll(fileByte,-1)
	text := Models.DTextRegex.FindAll(fileByte,-1)
	summary := Models.DSummaryRegex.FindAll(fileByte,-1)

	log.Debug("Number of Documents parsed -",len(text))
	log.Debug("Now building Trie for ",Models.ThresholdDocument, " documents")
	for i:=0;i<len(text) && i<Models.ThresholdDocument;i++{
		textStr := string(text[i])
		summaryStr := string(summary[i])
		reviewScoreStr := string(reviewScore[i])
		productIdStr := string(productId[i])
		userIDStr := string(userID[i])
		profileNameStr := string(profileName[i])
		helpfulStr := string(helpful[i])
		timeStr := string(time[i])

		textStr = strings.Replace(textStr,"review/text: ","",1)
		summaryStr = strings.Replace(summaryStr,"review/summary: ","",1)
		reviewScoreStr = strings.Replace(reviewScoreStr,"review/score: ","",1)
		productIdStr = strings.Replace(productIdStr,"product/productId: ","",1)
		userIDStr = strings.Replace(userIDStr,"review/userId: ","",1)
		profileNameStr = strings.Replace(profileNameStr,"review/profileName: ","",1)
		helpfulStr = strings.Replace(helpfulStr,"review/helpfulness: ","",1)
		timeStr = strings.Replace(timeStr,"review/time: ","",1)

		Models.ResponseDocuments = append(Models.ResponseDocuments,Models.ResponseDocument{
			ProductID:   productIdStr,
			UserID:      userIDStr,
			ProfileName: profileNameStr,
			Helpful:     helpfulStr,
			Score:       reviewScoreStr,
			Time:        timeStr,
			Text:        textStr,
			Summary:     summaryStr,
		})

		rs,_ := strconv.ParseFloat(reviewScoreStr,64)

		tr := trie.NewTrie()
		combinedStr := make([]string,0)
		combinedStr = append(combinedStr,ParseString(textStr)...)
		combinedStr = append(combinedStr,ParseString(summaryStr)...)
		for _,word := range combinedStr{
			tr.Add(word)
		}
		Models.Documents = append(Models.Documents,Models.Document{
			ReviewByScore: rs,
			Trie:          tr,
		})
	}
}

func main(){
	router := mux.NewRouter().StrictSlash(true)
	//Heartbeat router
	router.HandleFunc("/", heartbeat)
	router.HandleFunc("/search/words",Controller.SearchDocument)
	router.Use(recoverHandler)
	addr := ":" + Models.Port
	log.Info(Models.Port)
	log.Info(addr)
	errHttp := http.ListenAndServe(addr, router) // setting listening port

	if errHttp != nil {
		log.WithFields(log.Fields{
			"addr": addr,
			"err":  errHttp,
		}).Fatal("Unable to create HTTP Service ")
	}

}


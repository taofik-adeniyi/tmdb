package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Movie struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `json:"genre_ids"`
	ID               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}
type Dates struct {
	Maximum string `json:"maximum"`
	Minimum string `json:"minimum"`
}
type MovieRes struct {
	Dates         Dates   `json:"dates"`
	Page          int     `json:"page"`
	Results       []Movie `json:"results"`
	Total_Pages   int     `json:"total_pages"`
	Total_Results int     `json:"total_results"`
}
type MovieResError struct {
	Success        bool   `json:"success"`
	Status_Code    int    `json:"status_code"`
	Status_Message string `json:"status_message"`
}

const BASE_URL = "https://api.themoviedb.org/3/movie/"

var MOVIE_TYPES = [4]string{"upcoming", "top", "popular", "playing"}
var mapMovieTypeArgs = map[string]string{
	"playing":  "now_playing",
	"popular":  "popular",
	"top":      "top_rated",
	"upcoming": "upcoming",
}

func helper() {
	fmt.Println("tmdb [flags]")
	fmt.Println("tmdb [command]")
	fmt.Println("")

	fmt.Println("Available commands:")
	fmt.Println("tmdb --help")
	fmt.Println("tmdb --type e.g: [tmdb --type playing], [tmdb --type popular], [tmdb --type top], [tmdb --type upcoming]")
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error: loading env: %v", err)
	}
	tmdbAuthKey := os.Getenv("TMDB_AUTH_KEY")
	args := os.Args
	flags := [2]string{"--type", "--help"}
	if len(args) > 3 {
		helper()
		log.Fatal("Here are the valid commands")
	}
	if len(args) == 1 {
		helper()
		return
	}
	if len(args) == 2 {
		helper()
		return
	}
	var movieType string
	if len(args) == 3 {
		if args[1] == flags[0] {
			movieTypeValid := false
			for _, value := range MOVIE_TYPES {
				if value == args[2] {
					movieType = args[2]
					movieTypeValid = true
					break
				}
			}
			if !movieTypeValid {
				helper()
				log.Fatalf("invalid movie type: %v", args[2])
			}
		} else {
			helper()
			log.Fatalf("invalid movie type: %v", args[2])
		}
	}

	var tmdbMovieType = mapMovieTypeArgs[movieType]
	var url = BASE_URL + tmdbMovieType + "?language=en-US&page=1"
	bytes, err := requestHandler(url, tmdbAuthKey)
	if err != nil {
		log.Fatalf("request handler err: %v", err)
	}

	var moviesRes MovieRes
	var movieError MovieResError
	err = json.Unmarshal(bytes, &moviesRes)
	if err == nil && len(moviesRes.Results) == 0 {
		err = json.Unmarshal(bytes, &movieError)
		if err == nil && movieError.Status_Message != "" {
			log.Fatalf("API error: %s", movieError.Status_Message)
		}
	}

	err = json.Unmarshal(bytes, &moviesRes)
	if err != nil {
		log.Fatalf("unmarshalling: %v", err)
	}
	for key, value := range moviesRes.Results {
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Printf("  🎬  %d. %s\n", key+1, value.Title)
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Printf("  📅  Release Date: %s\n", value.ReleaseDate)
		fmt.Printf("  ⭐  Vote Count: %d\n", value.VoteCount)
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Printf("  📖  Overview: %s\n", value.Overview)
		fmt.Println("╚══════════════════════════════════════════════════╝")
		fmt.Println()
	}
}

func requestHandler(url string, authKey string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", authKey))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return resByte, nil
}

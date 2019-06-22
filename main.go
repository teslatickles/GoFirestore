package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/iterator"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message

	w.Write([]byte(message))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func handleRequests() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/index", indexPage)
}

func setStory(img string, tags [2]string, title string) interface{} {
	story := map[string]interface{}{
		"img":   img,
		"tags":  tags,
		"title": title,
	}
	return story
}

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("<file path to json file that was generated from generating private key in service accounts in settings for intended project at Firebase>")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// this is a safeguard/switch keeping the stories from
	// being read and printed on command line
	readData := false

	if readData {
		// this chunk reads from FireStore cloud DB
		// and prints data to command line.
		it := client.Collection("stories").Documents(ctx)
		for {
			doc, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			fmt.Println(doc.Data())
			fmt.Println(doc.Data()["tags"])
		}
	}

	// CAUTION/WARNING --> this chunk writes to a Google FireStore Cloud DB.
	// _, _, err = client.Collection("stories").Add(ctx, map[string]interface{}{
	// 	"img":   "https://pictures-of-cats.org/wp-content/uploads/2012/09/stopping-a-cat-biting-you-1.jpg",
	// 	"tags":  [2]string{"favorite", "happy"},
	// 	"title": "Cat Bit",
	// })

	// this is a safeguard/switch preventing an additional story from
	// being generated when running main.go
	storyNeedsAdded := false

	if storyNeedsAdded {
		result, dateTime, err := client.Collection("stories").Add(ctx, setStory("https://pbs.twimg.com/profile_images/548941003546042368/pn9oi0Co_400x400.png", [2]string{"sassy", "fancy"}, "happy days"))
		log.Print(result)
		log.Print(dateTime)
		log.Print(err)
	}

	handleRequests()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

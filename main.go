package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	lipsum "github.com/derektata/lorem/ipsum"
	"github.com/kluzzebass/graphql-demo/graph"
	"github.com/kluzzebass/graphql-demo/graph/model"
	"github.com/kluzzebass/graphql-demo/safemap"
	"github.com/kluzzebass/graphql-demo/subscription"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

var users = initUsers()
var posts = initPosts()
var postCreatedSub = subscription.NewManager[*model.Post]()
var postDeletedSub = subscription.NewManager[int]()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		UserMap:        users,
		PostMap:        posts,
		PostCreatedSub: postCreatedSub,
		PostDeletedSub: postDeletedSub,
	}}))

	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		log.Println("--------------------------------")
		return next(ctx)
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initUsers() *safemap.SafeMap[int, *model.User] {
	users := safemap.New[int, *model.User]()
	users.Set(1, &model.User{
		ID:          1,
		Name:        "Alice Wonderland",
		Email:       "alice.w@example.com",
		PhoneNumber: "555-123-4567",
		Address: &model.Address{
			Street:  "123 Rabbit Hole",
			City:    "Wonderland City",
			ZipCode: "90210",
			Country: "Fantasy",
		},
		Role:      "admin",
		LastLogin: time.Now(),
		Preferences: &model.Preferences{
			Theme:         "dark",
			Notifications: true,
		},
	})
	users.Set(2, &model.User{
		ID:          2,
		Name:        "Bob The Builder",
		Email:       "bob.b@example.com",
		PhoneNumber: "555-987-6543",
		Address: &model.Address{
			Street:  "456 Construction Site",
			City:    "Builderville",
			ZipCode: "12345",
			Country: "Imagination",
		},
		Role:      "editor",
		LastLogin: time.Now(),
		Preferences: &model.Preferences{
			Theme:         "light",
			Notifications: false,
		},
	})
	return users
}

func initPosts() *safemap.SafeMap[int, *model.Post] {
	g := lipsum.NewGenerator()
	g.WordsPerSentence = 10
	g.SentencesPerParagraph = 5
	g.CommaAddChance = 3

	posts := safemap.New[int, *model.Post]()

	posts.Set(1, &model.Post{
		ID:       1,
		Title:    "New 'Hyper-Aware' Smartwatch Now Judges Your Life Choices In Real-Time",
		Ingress:  "Wearable tech reaches its peak with the 'MoralCompass 3000,' which not only tracks your steps but also audibly tuts when you spend too much on online impulse buys or opt for a third slice of pizza.",
		Content:  g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:   1,
		Category: "Technology",
	})

	posts.Set(2, &model.Post{
		ID:       2,
		Title:    "Scientists Discover Evidence That Entire Universe Is Just A Toddler's Unfinished Finger Painting",
		Ingress:  "Groundbreaking astronomical research has revealed cosmic smudges and erratic crayon lines, leading experts to conclude that reality as we know it is merely a messy art project abandoned by an impatient celestial child.",
		Content:  g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:   1,
		Category: "Science",
	})

	posts.Set(3, &model.Post{
		ID:       3,
		Title:    "CEO Fulfills Lifelong Dream Of Laying Off Thousands To Optimize 'Synergistic Efficiencies'",
		Ingress:  "In a powerful display of ruthless ambition, CEO Sterling Blackwood announced mass redundancies, declaring it 'the most fulfilling moment of my career' as the company embraces a 'leaner, meaner, and infinitely more profitable' future.",
		Content:  g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:   2,
		Category: "Business",
	})

	posts.Set(4, &model.Post{
		ID:       4,
		Title:    "Streaming Service Now Just Showing Empty Room With Faint Echo Of Old Movie Quotes",
		Ingress:  "In a bold move to 'revolutionize content consumption,' 'VaporView+' announced its new flagship offering: a static shot of an unfurnished room, occasionally punctuated by whispers of classic film dialogue, designed for 'maximal chill and minimal engagement'.",
		Content:  g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:   2,
		Category: "Entertainment",
	})

	posts.Set(5, &model.Post{
		ID:       5,
		Title:    "Wellness Guru Attributes All Success To 'Standing Up Occasionally'",
		Ingress:  "Billionaire lifestyle coach Serenity Moonbeam credits her unparalleled vitality and financial success to the simple, yet revolutionary practice of 'not always sitting down,' sparking a global movement of cautious verticality.",
		Content:  g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:   2,
		Category: "Health",
	})

	return posts
}

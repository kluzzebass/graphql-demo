package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"graphql-demo/internal/graph/generated"
	"graphql-demo/internal/graph/model"
	"graphql-demo/internal/graph/resolver"
	"graphql-demo/internal/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	lipsum "github.com/derektata/lorem/ipsum"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		UserService:    initUsers(),
		PostService:    initPosts(),
		CommentService: initComments(),
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

func initUsers() *service.UserService {
	users := service.NewUserService()

	users.Set(&model.User{
		ID:          1,
		Name:        "Alice Wonderland",
		Email:       "alice.w@example.com",
		PhoneNumber: "555-123-4567",
		Address: &model.Address{
			Street:  "1 Rabbit Hole",
			City:    "Wonderland City",
			ZipCode: "90210",
			Country: "Fantasy",
		},
		Role:      "admin",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 60),
		LastLogin: time.Now(),
		Preferences: &model.Preferences{
			Theme:         "dark",
			Notifications: true,
		},
	})

	users.Set(&model.User{
		ID:          2,
		Name:        "Bob The Builder",
		Email:       "bob.b@example.com",
		PhoneNumber: "555-987-6543",
		Address: &model.Address{
			Street:  "2 Construction Site",
			City:    "Builderville",
			ZipCode: "12345",
			Country: "Imagination",
		},
		Role:      "editor",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 59),
		LastLogin: time.Now(),
		Preferences: &model.Preferences{
			Theme:         "light",
			Notifications: false,
		},
	})
	return users
}

func initPosts() *service.PostService {
	posts := service.NewPostService()

	g := lipsum.NewGenerator()
	g.WordsPerSentence = 10
	g.SentencesPerParagraph = 5
	g.CommaAddChance = 3

	posts.Set(&model.Post{
		ID:        1,
		Title:     "New 'Hyper-Aware' Smartwatch Now Judges Your Life Choices In Real-Time",
		Ingress:   "Wearable tech reaches its peak with the 'MoralCompass 3000,' which not only tracks your steps but also audibly tuts when you spend too much on online impulse buys or opt for a third slice of pizza.",
		Content:   g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:    1,
		Category:  "Technology",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 30),
	})

	posts.Set(&model.Post{
		ID:        2,
		Title:     "Scientists Discover Evidence That Entire Universe Is Just A Toddler's Unfinished Finger Painting",
		Ingress:   "Groundbreaking astronomical research has revealed cosmic smudges and erratic crayon lines, leading experts to conclude that reality as we know it is merely a messy art project abandoned by an impatient celestial child.",
		Content:   g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:    1,
		Category:  "Science",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 29),
	})

	posts.Set(&model.Post{
		ID:        3,
		Title:     "CEO Fulfills Lifelong Dream Of Laying Off Thousands To Optimize 'Synergistic Efficiencies'",
		Ingress:   "In a powerful display of ruthless ambition, CEO Sterling Blackwood announced mass redundancies, declaring it 'the most fulfilling moment of my career' as the company embraces a 'leaner, meaner, and infinitely more profitable' future.",
		Content:   g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:    2,
		Category:  "Business",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 28),
	})

	posts.Set(&model.Post{
		ID:        4,
		Title:     "Streaming Service Now Just Showing Empty Room With Faint Echo Of Old Movie Quotes",
		Ingress:   "In a bold move to 'revolutionize content consumption,' 'VaporView+' announced its new flagship offering: a static shot of an unfurnished room, occasionally punctuated by whispers of classic film dialogue, designed for 'maximal chill and minimal engagement'.",
		Content:   g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:    2,
		Category:  "Entertainment",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 27),
	})

	posts.Set(&model.Post{
		ID:        5,
		Title:     "Wellness Guru Attributes All Success To 'Standing Up Occasionally'",
		Ingress:   "Billionaire lifestyle coach Serenity Moonbeam credits her unparalleled vitality and financial success to the simple, yet revolutionary practice of 'not always sitting down,' sparking a global movement of cautious verticality.",
		Content:   g.GenerateParagraphs(rand.Intn(5) + 1),
		UserID:    2,
		Category:  "Health",
		CreatedAt: time.Now().Add(-time.Hour * 24 * 26),
	})

	return posts
}

func initComments() *service.CommentService {
	comments := service.NewCommentService()

	comments.Set(&model.Comment{
		ID:        1,
		PostID:    1,
		UserID:    2,
		Content:   "I bought this watch, and the first thing it did was judge me for buying it.",
		CreatedAt: time.Now().Add(-time.Hour * 2),
	})

	comments.Set(&model.Comment{
		ID:        2,
		PostID:    1,
		UserID:    1,
		Content:   "There's one born every minute.",
		CreatedAt: time.Now().Add(-time.Hour * 1),
	})

	return comments
}

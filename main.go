package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
)

var (
	gh_access_token string = os.Getenv("GITHUB_TOKEN")

	starsCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_stars_count",
		Help: "Number of stars the repo has received",
	})
	forksCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_forks_count",
		Help: "Number of times the repo has been forked",
	})
	issuesCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_issues_count",
		Help: "Number of issues logged against the repo",
	})
	watchersCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_watchers_count",
		Help: "Number of accounts watching the repo",
	})
	viewsCounter14 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_views_count_14days",
		Help: "Number of views/visits the repo received",
	})
	viewsCounterUnique14 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_unique_views_count_14days",
		Help: "Number of unique views/visits the repo received",
	})
	clonesCounter14 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_clones_count14days",
		Help: "Number of times the repo has been forked",
	})
	clonesCounterUnique14 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "github_unique_clones_count14days",
		Help: "Number of unique times the repo has been forked",
	})
)

func recordMetrics() {
	fmt.Printf("RecordMetrics\n")
	go func() {
		for {
			gatherMetrics()
			time.Sleep(60 * time.Second)
		}
	}()
}

func gatherMetrics() {

	fmt.Printf("gh %v\n", gh_access_token)
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh_access_token},
	)
	tokenClient := oauth2.NewClient(context, tokenService)
	client := github.NewClient(tokenClient)
	repo, _, err := client.Repositories.Get(context, "aws-samples", "amazon-chime-sdk-pstn-audio-workshop")
	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}
	if repo == nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}
	o := github.TrafficBreakdownOptions{Per: "day"}
	tv, _, err := client.Repositories.ListTrafficViews(context, "aws-samples", "amazon-chime-sdk-pstn-audio-workshop", &o)
	c := 0
	u := 0
	for _, e := range tv.Views {
		fmt.Printf("%v: %v %v\n", *e.Timestamp, *e.Count, *e.Uniques)
		c += *e.Count
		u += *e.Uniques
	}

	cv, _, err := client.Repositories.ListTrafficClones(context, "aws-samples", "amazon-chime-sdk-pstn-audio-workshop", &o)
	c = 0
	u = 0
	for _, e := range cv.Clones {
		fmt.Printf("%v: %v %v\n", *e.Timestamp, *e.Count, *e.Uniques)
		c += *e.Count
		u += *e.Uniques
	}

	// Get Rate limit information

	rateLimit, _, err := client.RateLimits(context)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		return
	}

	fmt.Printf("Limit: %d \nRemaining %d \n", rateLimit.Core.Limit, rateLimit.Core.Remaining)

	starsCounter.Set(float64(*repo.StargazersCount))
	forksCounter.Set(float64(*repo.ForksCount))
	issuesCounter.Set(float64(*repo.OpenIssuesCount))
	watchersCounter.Set(float64(*repo.WatchersCount))
	viewsCounter14.Set(float64(*tv.Count))
	viewsCounterUnique14.Set(float64(*tv.Uniques))
	clonesCounter14.Set(float64(*cv.Count))
	clonesCounterUnique14.Set(float64(*cv.Uniques))

	/*
		fmt.Printf("FullName:         %v\n", *repo.FullName)
		fmt.Printf("StarGazers:       %v\n", *repo.StargazersCount)
		fmt.Printf("ForksCount:       %v\n", *repo.ForksCount)
		fmt.Printf("NetworkCount:     %v\n", *repo.NetworkCount)
		fmt.Printf("OpenIssuesCount:  %v\n", *repo.OpenIssuesCount)
		fmt.Printf("SubscribersCount: %v\n", *repo.SubscribersCount)
		fmt.Printf("WatchersCount:    %v\n", *repo.WatchersCount)
		fmt.Printf("ViewsCount: %v\n", *tv.Count)
		fmt.Printf("ViewsUniquesCount: %v\n", *tv.Uniques)
		fmt.Printf("CloneCount: %v\n", *cv.Count)
		fmt.Printf("CloneUniquesCount: %v\n", *cv.Uniques)
	*/
}

func main() {

	if gh_access_token == "" {
		fmt.Printf("credentials not provided, please export GITHUB_TOKEN \n")
		os.Exit(1)
	}

	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9090", nil)

}

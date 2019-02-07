package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/jayeshsolanki93/devgorant"
)

const version = "1.0.0-dev"

var rantsFlag, randomFlag, weeklyFlag *bool
var sortFlag, profileFlag, searchFlag *string
var rantFlag, limitFlag, skipFlag *int

func printRantPreview(rant devgorant.RantModel) {
	var text string
	if len(rant.Text) > 160 {
		text = rant.Text[:155]
		text += "[...]\n[Read more...]"
	} else {
		text = rant.Text
	}
	fmt.Printf("(+%d) <ID:%d> by %s(+%d):\n%s\n%s {%d comments}\n\n",
		rant.Score, rant.Id, rant.UserUsername, rant.UserScore, text, rant.Tags, rant.NumComments)
}

func printRant(rant devgorant.RantModel, comments []devgorant.CommentModel, limit int) {
	fmt.Printf("(+%d) <ID:%d> by %s(+%d):\n%s\n%s\n",
		rant.Score, rant.Id, rant.UserUsername, rant.UserScore, rant.Text, rant.Tags)

	fmt.Printf("\nComments[%d]\n\n", rant.NumComments)

	for i, c := range comments {
		if i < limit {
			fmt.Printf("(+%d) <ID:%d> by %s(+%d):\n%s\n\n",
				c.Score, c.Id, c.UserUsername, c.UserScore, c.Body)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func executor(cmd string) {
	if cmd != "" {
		var c []string
		var command string
		devrant := devgorant.New()
		c = strings.Split(cmd, " ")
		if len(c) < 1 {
			fmt.Println()
			return
		}
		command = c[0]
		switch command {
		case "sort":
			if len(c) >= 2 {
				if c[1] == "algo" || c[1] == "top" || c[1] == "recent" {
					*sortFlag = c[1]
				}
			}
			fmt.Printf("sort method is %s\n", *sortFlag)
		case "limit":
			if len(c) >= 2 {
				lim, err := strconv.ParseInt(c[1], 10, 64)
				if err == nil {
					*limitFlag = int(lim)
				}
			}
			fmt.Printf("limit is %d rants\n", *limitFlag)
		case "skip":
			if len(c) >= 2 {
				skp, err := strconv.ParseInt(c[1], 10, 64)
				if err == nil {
					*skipFlag = int(skp)
				}
			}
			fmt.Printf("skipping every %d rants\n", *skipFlag)
		case "rant":
			if len(c) >= 2 {
				rnt, err := strconv.ParseInt(c[1], 10, 64)
				if err == nil {
					*rantFlag = int(rnt)
				}
			}
			rant, comments, err := devrant.Rant(*rantFlag)
			check(err)
			printRant(rant, comments, *limitFlag)
		case "rants":
			rants, err := devrant.Rants(*sortFlag, *limitFlag, *skipFlag)
			check(err)
			for _, r := range rants {
				printRantPreview(r)
			}
		case "weekly":
			rants, err := devrant.WeeklyRants()
			check(err)
			for _, r := range rants {
				printRantPreview(r)
			}
		case "random":
			rant, err := devrant.Surprise()
			check(err)
			printRant(rant, nil, *limitFlag)
		case "search":
			if len(c) >= 2 && c[1] != "" {
				*searchFlag = c[1]
				rants, err := devrant.Search(*searchFlag)
				check(err)
				for _, r := range rants {
					printRantPreview(r)
				}
			}
		case "profile":
			if len(c) >= 2 && c[1] != "" {
				*profileFlag = c[1]
				user, content, err := devrant.Profile(*profileFlag)
				check(err)

				// convert epoch timestamp on profile to some thing sensible
				t, err := strconv.ParseInt(strconv.Itoa(user.CreatedTime), 10, 64)
				check(err)
				ts := time.Unix(t, int64(0))
				timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d UTC", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute())

				fmt.Printf("%s(+%d)\nLocation: %s\nJoined: %s\nAbout: %s\nSkills: %s\nGitHub: %s\n",
					user.Username, user.Score, user.Location, timestamp, user.About, user.Skills, user.Github)
				fmt.Printf("Rants:\n")
				for _, r := range content.Rants {
					printRantPreview(r)
				}
			}
		case "exit", "quit":
			os.Exit(0)
		case "help", "commands":
			fmt.Println("available commands:")
			fmt.Println("    sort [algo|top|recent]")
			fmt.Println("    limit [int]")
			fmt.Println("    rant <int>")
			fmt.Println("    rants")
			fmt.Println("    weekly")
			fmt.Println("    random")
			fmt.Println("    search <username>")
			fmt.Println("    profile <username>")
			fmt.Println("    help, commands")
			fmt.Println("    exit, quit")
		default:
			fmt.Println()
		}
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	if d.GetWordBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	s := []prompt.Suggest{
		{Text: "sort", Description: "sort by: algo|top|recent"},
		{Text: "limit", Description: "limit number of rants displayed"},
		{Text: "rants", Description: "show rants"},
		{Text: "rant", Description: "show rant by ID"},
		{Text: "weekly", Description: "show weekly rants"},
		{Text: "random", Description: "show random rants"},
		{Text: "search", Description: "search rants by term"},
		{Text: "profile", Description: "show profile of user"},
		{Text: "algo", Description: "use devRant algo sort"},
		{Text: "top", Description: "use devRant top sort"},
		{Text: "recent", Description: "use devRant recent sort"},
		{Text: "help", Description: "show help text"},
		{Text: "commands", Description: "show help text"},
		{Text: "exit", Description: "exit/quit drant"},
		{Text: "quit", Description: "exit/quit drant"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	sortFlag = flag.String("m", "algo", "sorting method for -r: algo, top, or recent")
	limitFlag = flag.Int("l", 50, "number of rants to fetch for -r")
	skipFlag = flag.Int("i", 0, "number of rants to skip for -r")

	rantFlag = flag.Int("R", -1, "fetches rant and its comments given rant id")
	rantsFlag = flag.Bool("r", false, "fetches rants")
	weeklyFlag = flag.Bool("w", false, "fetches rants tagged weekly")
	randomFlag = flag.Bool("n", false, "fetches random rant")

	searchFlag = flag.String("s", "", "search for rants matching given term")
	profileFlag = flag.String("u", "", "fetches ranter's profile data")
	flag.Parse()

	devrant := devgorant.New()

	if *rantsFlag {
		rants, err := devrant.Rants(*sortFlag, *limitFlag, *skipFlag)
		check(err)
		for _, r := range rants {
			printRantPreview(r)
		}
	} else if *rantFlag != -1 {
		rant, comments, err := devrant.Rant(*rantFlag)
		check(err)
		printRant(rant, comments, *limitFlag)
	} else if *randomFlag {
		rant, err := devrant.Surprise()
		check(err)
		printRant(rant, nil, *limitFlag)
	} else if *weeklyFlag {
		rants, err := devrant.WeeklyRants()
		check(err)
		for _, r := range rants {
			printRantPreview(r)
		}
	} else if *searchFlag != "" {
		rants, err := devrant.Search(*searchFlag)
		check(err)
		for _, r := range rants {
			printRantPreview(r)
		}
	} else if *profileFlag != "" {
		user, content, err := devrant.Profile(*profileFlag)
		check(err)

		// convert epoch timestamp on profile to some thing sensible
		t, err := strconv.ParseInt(strconv.Itoa(user.CreatedTime), 10, 64)
		check(err)
		ts := time.Unix(t, int64(0))
		timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d UTC", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute())

		fmt.Printf("%s(+%d)\nLocation: %s\nJoined: %s\nAbout: %s\nSkills: %s\nGitHub: %s\n",
			user.Username, user.Score, user.Location, timestamp, user.About, user.Skills, user.Github)
		fmt.Printf("Rants:\n")
		for _, r := range content.Rants {
			printRantPreview(r)
		}
	} else {
		var headText string
		headText = "      _            ____             _\n   __| | _____   _|  _ \\ __ _ _ __ | |_\n  / _` |/ _ \\ \\ / / |_) / _` | '_ \\| __|\n | (_| |  __/\\ V /|  _ < (_| | | | | |_\n  \\__,_|\\___| \\_/ |_| \\_\\__,_|_| |_|\\__|"
		fmt.Println(headText)
		fmt.Printf("   github.com/corscheid/drant  %s\n\n", version)
		p := prompt.New(executor, completer,
			prompt.OptionTitle("drant"),
			prompt.OptionPrefix("--> "),
			prompt.OptionSelectedDescriptionTextColor(prompt.DarkGray),
		)
		p.Run()
	}
}

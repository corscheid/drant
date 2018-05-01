package main

import (
  "flag"
  "fmt"
  "github.com/jayeshsolanki93/devgorant"
  "log"
)

func printRantPreview(rant devgorant.RantModel) {
  var text string;
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

func main() {
  var rantsFlag, randomFlag, weeklyFlag *bool
  var sortFlag, profileFlag, searchFlag *string
  var rantFlag, limitFlag, skipFlag *int
  var cmdPrefix, headText string

  headText = "      _            ____             _\n   __| | _____   _|  _ \\ __ _ _ __ | |_\n  / _` |/ _ \\ \\ / / |_) / _` | '_ \\| __|\n | (_| |  __/\\ V /|  _ < (_| | | | | |_\n  \\__,_|\\___| \\_/ |_| \\_\\__,_|_| |_|\\__|"
  rantsFlag = flag.Bool("r", false, "fetches rants")
  sortFlag = flag.String("m", "algo", "sorting method for -r: algo, top, or recent")
  limitFlag = flag.Int("l", 50, "number of rants to fetch for -r")
  skipFlag = flag.Int("i", 0, "number of rants to skip for -r")

  rantFlag = flag.Int("R", -1, "fetches rant and its comments given rant id")
  randomFlag = flag.Bool("n", false, "fetches random rant")
  weeklyFlag = flag.Bool("w", false, "fetches rants tagged weekly")

  searchFlag = flag.String("s", "", "search for rants matching given term")
  profileFlag = flag.String("u", "", "fetches ranter's profile data")
  flag.Parse()

  devrant := devgorant.New()

  if *rantsFlag {
    rants, err := devrant.Rants(*sortFlag, *limitFlag, *skipFlag)
    check(err)
    for _, r := range rants {
      //fmt.Printf("[%d] ", i)
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
      //fmt.Printf("[%d] ", i)
      printRantPreview(r)
    }
  } else if *searchFlag != "" {
    rants, err := devrant.Search(*searchFlag)
    check(err)
    for _, r := range rants {
      //fmt.Printf("[%d] ", i)
      printRantPreview(r)
    }
  } else if *profileFlag != "" {
    user, content, err := devrant.Profile(*profileFlag)
    check(err)
    fmt.Printf("%s(+%d)\nLocation: %s\nJoined: %d\nAbout: %s\nSkills: %s\nGitHub: %s\n",
      user.Username, user.Score, user.Location, user.CreatedTime, user.Skills, user.Github)
    fmt.Printf("Rants:\n")
    for _, r := range content.Rants {
        //fmt.Printf("[%d] ", i)
        printRantPreview(r)
    }
  } else {
    fmt.Println(headText)

    //TODO REPL
    fmt.Print(cmdPrefix)
  }
}

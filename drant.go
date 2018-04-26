package main

import (
  "flag"
  "fmt"
  "github.com/jayeshsolanki93/devgorant"
  "log"
)

func printRant(rant devgorant.RantModel) {
  fmt.Printf("<ID:%d> (+%d) {%d comments} by %s(+%d):\n%s\n%s\n\n",
    rant.Id, rant.Score, rant.NumComments, rant.UserUsername, rant.UserScore, rant.Text, rant.Tags)
}

func printComments(comments []devgorant.CommentModel) {
  for i, c := range comments {
    fmt.Printf("[%d] <ID:%d> (+%d) by %s(+%d):\n%s\n\n",
      i, c.Id, c.Score, c.UserUsername, c.UserScore, c.Body)
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
    for i, r := range rants {
      fmt.Printf("[%d] ", i)
      printRant(r)
    }
  } else if *rantFlag != -1 {
    rant, comments, err := devrant.Rant(*rantFlag)
    check(err)
    printRant(rant)
    fmt.Printf("Comments[%d]\n\n", rant.NumComments)
    printComments(comments)
  } else if *randomFlag {
    rant, err := devrant.Surprise()
    check(err)
    printRant(rant)
  } else if *weeklyFlag {
    rants, err := devrant.WeeklyRants()
    check(err)
    for i, r := range rants {
      fmt.Printf("[%d] ", i)
      printRant(r)
    }
  } else if *searchFlag != "" {
    rants, err := devrant.Search(*searchFlag)
    check(err)
    for i, r := range rants {
      fmt.Printf("[%d] ", i)
      printRant(r)
    }
  } else if *profileFlag != "" {
    user, content, err := devrant.Profile(*profileFlag)
    check(err)
    fmt.Printf("%s(+%d)\nLocation: %s\nJoined: %d\nAbout: %s\nSkills: %s\nGitHub: %s\n",
      user.Username, user.Score, user.Location, user.CreatedTime, user.Skills, user.Github)
    fmt.Printf("Rants:\n")
    for i, r := range content.Rants {
      if i < *limitFlag {
        fmt.Printf("[%d] ", i)
        printRant(r)
      }
    }
  } else {
    fmt.Println("devRant CLI | run drant -h for help")
  }
}

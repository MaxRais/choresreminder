package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"net/http"
	"strings"
)

var names = []string {"Max", "Derek", "Alex", "Spencer"}

var nameToEmail = map[string]string{
	"Max": "max.j.rais@gmail.com",
	"Derek": "rais.m@husky.neu.edu",
	"Alex": "yolts12321@yahoo.com",
	"Spencer": "yankeesrule422@aim.com",
}

var currentDIndex = 0
var currentTIndex = 0

func incrementIndex(index int) int {
	index++
	if index >= len(names) {
		return 0
	} else {
		return index
	}
}

func sendEmail(dishes bool) error {
	name := ""
	index := -1
	email := ""
	message := ""

	if dishes {
		currentDIndex = incrementIndex(currentDIndex)
		index = currentDIndex
		message = "dishes"
	} else {
		currentTIndex = incrementIndex(currentTIndex)
		index = currentTIndex
		message = "trash"
	}

	name = names[index]
	email = nameToEmail[name]

	m := gomail.NewMessage()
	m.SetHeader("From", "Chore Reminder")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "It's your turn!")
	m.SetBody("text/html", fmt.Sprintf("Hi %v,<br><br>Please do the %v<br><br>From,<br>Chore Reminder Bot", name, message))

	d := gomail.NewDialer("smtp.gmail.com", 587, "terracechores@gmail.com", "@1terrace")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func sendEmailUrl(w http.ResponseWriter, r *http.Request) {
	chore := r.URL.Path
	chore = strings.TrimPrefix(chore, "/")

	if chore == "favicon.ico" {
		return
	} else if chore != "dishes" && chore != "trash" {
		fmt.Println("Bad chore", chore)
		w.Write([]byte(fmt.Sprintf("Invalid chore: %v\n", chore)))
		return
	}

	err := sendEmail(chore=="dishes")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Something went wrong sending email\n"))
	} else {
		w.Write([]byte("Email sent successfully!\n"))
	}
}

func main() {
	fmt.Println("Starting sever on port 8080")

	http.HandleFunc("/", sendEmailUrl)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

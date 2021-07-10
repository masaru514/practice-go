package main

import (
    "fmt"
		"log"
    "example.com/greetings"
)

func main() {
    // Set properties of the predefined Logger, including
    // the log entry prefix and a flag to disable printing
    // the time, source file, and line number.
    log.SetPrefix("greetings: ")
    log.SetFlags(0)

		names := []string{"Gladys", "masaru", "tarou"}

    // Request a greeting message.
		// 呼び出した関数内にある同じ変数を定義してあげる必要がある？
    message, err := greetings.Hellos(names)
    // If an error was returned, print it to the console and
    // exit the program.
    if err != nil {
        log.Fatal(err)
    }

    // If no error was returned, print the returned message
    // to the console.
    fmt.Println(message)
}
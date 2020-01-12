package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Score struct {
	num_correct, total int
}

func doTheQuiz(problems []problem) <-chan int {
	out := make(chan int)
	go func() {
		num_correct := 0
		for _, p := range problems {
			fmt.Printf("%v = ", p.q)
			answer := 0
			_, err := fmt.Scanf("%d", &answer)
			if err != nil {
				log.Fatalln("unexpected")
			}
			if answer == p.a {
				num_correct++
			}
		}
		out <- num_correct
		close(out)
	}()
	return out
}

type problem struct {
	q string
	a int
}

func parse(f *string) []problem {
	csvfile, err := os.Open(*f)
	if err != nil {
		log.Fatalln("Failed on opening file")
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalln("failed reading lines from file")
	}
	problems := make([]problem, len(lines))
	for i, record := range lines {
		expected, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatalln("wrong expected answer format in problems file")
		}
		problems[i] = problem{
			q: record[0],
			a: expected,
		}
	}
	return problems
}
func main() {
	csv := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	problems := parse(csv)

	sc := doTheQuiz(problems)
	select {
	case s := <-sc:
		fmt.Printf("you scored %d out of %d!\n", s, len(problems))
	case <-time.After(time.Duration(*limit) * time.Second):
		fmt.Printf("\nsorry, timeout\n")
	}

}

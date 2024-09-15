package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mehmetymw/my-kvd/server"
)

func main() {
	os.Mkdir("data", os.ModePerm)
	store := server.NewStore("data.json")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		input := scanner.Text()
		parts := strings.Fields(input)

		switch parts[0] {
		case "set":
			if len(parts) != 3 {
				fmt.Println("Usage: set <key> <value>")
				continue
			}
			key := parts[1]
			value := parts[2]
			store.Set(key, value)

		case "genp":
			if len(parts) >= 2 && len(parts) <= 3 {
				continue
			}
			key := parts[1]

			var length int
			if len(parts) == 3 {
				parsedLength, err := strconv.ParseInt(parts[2], 10, 32)
				if err != nil {
					fmt.Errorf("length cannot convert to integer value: %w", length)
					continue
				}

				length = int(parsedLength)
			}

			store.SetPassword(key, length)

		case "get":
			if len(parts) != 2 {
				continue
			}
			key := parts[1]
			value, exists := store.Get(key)
			if !exists {
				fmt.Println("key not found")
				continue
			}
			fmt.Println(value)

		case "exit":
			fmt.Println("exiting from my-kvd")
			return

		default:
			fmt.Println("unknown command")
		}
	}
}

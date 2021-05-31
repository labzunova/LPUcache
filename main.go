package main

import (
	"fmt"
)

type InitValues struct {
	size, memory, TTL uint32
}

func initialize() (*InitValues, error) {
	var size, memory, TTL uint32

	fmt.Println("Enter Cache size")
	_, err := fmt.Scanf("%d", &size)
	if err != nil {
		return nil, err
	}

	fmt.Println("Enter time to live for values(in seconds)")
	_, err = fmt.Scanf("%d", &TTL)
	if err != nil {
		return nil, err
	}

	fmt.Println("Enter max Cache memory")
	_, err = fmt.Scanf("%d", &memory)
	if err != nil {
		return nil, err
	}

	result := InitValues{
		size: size,
		TTL: TTL,
		memory: memory,
	}
	return &result, err
}

func main() {
	meta, err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	c := NewCache(meta.size, meta.TTL, meta.memory)

	var action, value, result string
	var key uint32
	for {
		fmt.Println("\nEnter PUT or GET")
		_, err = fmt.Scanf("%s", &action)

		switch action {
		case "GET":
			fmt.Println("Enter the key")
			_, err = fmt.Scanf("%d", &key)

			result = c.Get(key)
			fmt.Println(result) // новая строчка просто чтобы красиво в консоли было :)

		case "PUT":
			fmt.Println("Enter the key")
			_, err = fmt.Scanf("%d", &key)
			fmt.Println("Enter the value")
			_, err = fmt.Scanf("%s", &value)

			err = c.Put(key, value)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Ok")
			}

		default:
			fmt.Println("wrong action")
			return
		}
	}
}

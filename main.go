package main

func main() {
	service, err := Initialize()
	if err != nil {
		panic(err)
	}
	defer service.DB.Close()
}

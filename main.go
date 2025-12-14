package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/dmsRosa6/MoooChain/commands"
	"github.com/dmsRosa6/MoooChain/options"
	"github.com/dmsRosa6/MoooChain/utils"
	"github.com/joho/godotenv"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	godotenv.Load()

}

func CallClear() error {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		return errors.New("your platform is unsupported")
	}

	return nil
}

func main() {

	log := configLog()

	err := CallClear()

	if err != nil {
		log.Fatal(err)
	}

	option := options.InitOptions(log)
	option.Print()
	r := utils.InitRedis()

	ctx := context.Background()
	_, err = r.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Redis was not initialized")
		return
	}

	parser := commands.NewParser()
	executer := commands.NewExecuter(log, option, r)

	run := true

	for run {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		command, args, err := parser.Parse(text)

		if err != nil {
			log.Print(err)
			continue
		}

		err = executer.Execute(command, args)

		if err != nil {
			log.Print(err)
		}

	}
}

func configLog() *log.Logger {
	return log.New(os.Stdout, "Moochain:", log.LstdFlags)

}

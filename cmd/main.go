package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	commands "github.com/dmsRosa6/MoooChain/internal/command"
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
}

func CallClear() error{
    value, ok := clear[runtime.GOOS]
    if ok {
        value()	
	} else {
        return errors.New("Your platform is unsupported.")
    }

	return nil
}

func main(){

	log := configLog()


	err := CallClear()
	
	if err != nil {
		log.Fatal("error initiation blockchain. err : %s", err)
		return
	}
	
	parser := commands.NewParser()
	
	executer := commands.NewExecuter(log)

	if err != nil {
		log.Fatal("error initiation blockchain. err : %s", err)
		return
	}

	run := true
	fmt.Println("**Moochain**")

	for run {

		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		command, args, err := parser.Parse(text)
		
		err = executer.Execute(command,args)

		if err != nil{
			log.Printf("err : %s\n", err)
		}

	}
	
}

func configLog() *log.Logger {
	return log.New(os.Stdout, "Moochain:", log.LstdFlags)

}
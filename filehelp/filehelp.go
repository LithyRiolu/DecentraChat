package filehelp

import (
	"log"
	"os"
)

func SaveToFile(filename string, m string) {
	//cwd, _ := os.Getwd()
	//filename := cwd + "/chat/files/chat.txt"
	//filename := getChatFIlePath()

	if _, err := os.Stat(filename); err != nil {
		CreateFile(filename)
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(m + "\n"); err != nil {
		panic(err)
	}

}

func CreateFile(filename string) {

	var file, err = os.Create(filename)
	if err != nil {
		log.Println("Chat file not present also Cannot create chat file : err - ", err)
		return
	}
	defer file.Close()
}

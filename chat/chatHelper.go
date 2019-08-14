package chat

import (
	"encoding/json"
	"log"
)

func ChatToJSON(c *Chat) []byte {
	j, _ := json.Marshal(c)
	return j
}

func JSONToChat(b []byte) Chat {
	c := Chat{}
	jerr := json.Unmarshal(b, &c)
	if jerr != nil {
		log.Println("Error in unmarshalling json to Chat - err : ", jerr)
	}

	return c
}

func ChatShowToJSON(c *ChatShow) []byte {
	j, _ := json.Marshal(c)
	return j
}

func JSONToChatShow(b []byte) ChatShow {
	c := ChatShow{}
	jerr := json.Unmarshal(b, &c)
	if jerr != nil {
		log.Println("Error in unmarshalling json to Chat - err : ", jerr)
	}

	return c
}

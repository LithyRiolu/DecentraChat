package chat

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"

	//"bytes"
	"html/template"
	//"io/ioutil"

	//"bytes"
	"encoding/json"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"time"

	//"io/ioutil"
	//"io/ioutil"
	"net/http"
	"path"

	"os"

	"../filehelp"
	"../httphelp"
	"../identity"
	//"./data"
	"../peers"
)

var CHATFILESSENTDIR = "/chat/files/" + "c" + os.Args[1] + "/sent/"
var CHATFILESRECVDIR = "/chat/files/" + "c" + os.Args[1] + "/recv/"
var CHATFILEPREFIX = "/chat/files/chat-"

const CHATBOXHTML = "/resources/html/index.html"

var chatIdentity identity.Identity

var isStarted bool

type Chat struct {
	From       string     `json:"from"`
	To         []string   `json:"to"`
	Message    string     `json:"message"`
	LoadedFile LoadedFile `json:"loadedFile"`
	Timestamp  time.Time  `json:"timestamp"`
}

type ChatShow struct {
	From           string    `json:"from"`
	To             []string  `json:"to"`
	Message        string    `json:"message"`
	LoadedFileName string    `json:"loadedFileName"`
	Timestamp      time.Time `json:"timestamp"`
}

type LoadedFile struct {
	FileName string `json:"filename"`
	FileData []byte `json:"filedata"`
}

type Chats struct {
	ChatList []Chat `json:"chats"`
}

type ChatsShow struct {
	ChatsShowList []ChatShow `json:"chatsShow"`
}

type ChatPage struct {
	IdDS      identity.Identity
	PeersDS   peers.Peers
	ChatsShow ChatsShow
	//Chats Chats//adding

}

func NewChat(identity identity.Identity, to []string, message string, loadedFile LoadedFile) Chat {
	c := Chat{}
	c.Timestamp = time.Now()
	c.From = identity.Id
	c.To = to
	c.Message = message
	c.LoadedFile = loadedFile
	return c
}

func NewChatShow(from string, to []string, message string, loadedFileName string) ChatShow {
	c := ChatShow{}
	c.Timestamp = time.Now()
	c.From = from
	c.To = to
	c.Message = message
	c.LoadedFileName = loadedFileName
	return c
}

func NewLoadedFile(name string, b []byte) LoadedFile {
	return LoadedFile{
		FileName: name,
		FileData: b,
	}
}

func NewChats() Chats {
	c := Chats{}

	return c
}

func NewChatsShow() ChatsShow {
	c := ChatsShow{}

	return c
}

func NewChatPage(idDS identity.Identity, peersDS peers.Peers, chatsShow ChatsShow) ChatPage {
	return ChatPage{
		IdDS:      idDS,
		PeersDS:   peersDS,
		ChatsShow: chatsShow,
	}
}

//internal funcs
func setIdentity(identity identity.Identity) {
	chatIdentity = identity
}

func getIdentity() identity.Identity {
	return chatIdentity
}

func getChatFIlePath() string {
	cwd, _ := os.Getwd()
	chatId := getIdentity()
	filename := cwd + CHATFILEPREFIX + chatId.Id + ".txt"
	return filename
}

//get /chat
func Begin(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) {

	if isStarted == false {
		setIdentity(identity)
		isStarted = true
	}

	makeClientDir(identity.Id)

	showChatsShow(w, identity, peerDS)

}

//makeClientDir creates a directory for the client
func makeClientDir(id string) {
	//CHATFILESSENTDIR = "/chat/files/" + "c" + os.Args[1] + "/sent/"
	//CHATFILESRECVDIR = "/chat/files/" + "c" + os.Args[1] + "/recv/"
	//CHATFILEPREFIX = "/chat/files/"+ "c" + os.Args[1] + "/chat-"

	wd, _ := os.Getwd()
	_ = os.MkdirAll(wd+"/chat/files/"+"c"+os.Args[1]+"/sent/", os.ModePerm)
	_ = os.MkdirAll(wd+"/chat/files/"+"c"+os.Args[1]+"/recv/", os.ModePerm)
	//_ = os.MkdirAll("/chat/files/"+ "c" + os.Args[1] + "/chat-", os.ModePerm)

}

// showChatsShow connects to HTML and executes the template
func showChatsShow(w http.ResponseWriter, identity identity.Identity, peerDS peers.Peers) {

	chatsToShow := chatFileToChatsShow()

	page := NewChatPage(identity, peerDS, chatsToShow) // to be used in template execute
	cwd, _ := os.Getwd()
	chatBoxHtml := path.Join(cwd, CHATBOXHTML) //"/resources/view/html/index.html")
	t, _ := template.ParseFiles(chatBoxHtml)
	_ = t.Execute(w, page) // todo - form - encType

}

// chatFileToChatsShow func reads the File - chat-id.txt and return a ChatsShow struct
func chatFileToChatsShow() ChatsShow {

	filename := getChatFIlePath()

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Println("Error in readChat() : err - ", err)
	}

	chatsShow := NewChatsShow()
	chatsShow.ChatsShowList = make([]ChatShow, 0)
	chatsShow.initialize(inputFile)

	//fmt.Println("chatFileToChatsShow -  ", chatsShow)

	return chatsShow
}

// initialize func takes file as param and inits the chatsShow struct
func (chatsShow *ChatsShow) initialize(inputFile *os.File) {
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines) // Read until a newline for each Scan() (Default)
	for inputScanner.Scan() {
		//fmt.Println(inputScanner.Text())   // get the buffered content as string
		//fmt.Println(inputScanner.Bytes())  // same content as above but as []byte
		chatShowJSON := inputScanner.Text()
		chatToShow := ChatShow{}
		jerr := json.Unmarshal([]byte(chatShowJSON), &chatToShow)
		if jerr != nil {
			log.Println("Error while converting json to chatShow - err : ", jerr)
			continue
		}
		chatsShow.ChatsShowList = append(chatsShow.ChatsShowList, chatToShow)
	}
}

// Called when a node sends a Chat object  - In Req body there is chat json
// Chat struct
// Continue func gets param - w, r, identity, peersDS
// Continue func should do -
// 0. save chat to local chat log
// 1. parse recv Chat
func Continue(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) {

	//save the chat in req body - generated from submit of chatform
	processChatFormSubmit(r, identity, peerDS)

	//showChatsShow(w, identity, peerDS)
	http.Redirect(w, r, "http://"+r.Host+"/chat", 301)
}

// processChatFormSubmit saves the chat form message to self folders - recv and update self chat-id.txt copy
func processChatFormSubmit(r *http.Request, identity identity.Identity, peerDS peers.Peers) {
	//save chat in req body - a. save everything except file in chats-id.txt
	// 								- first create chatsShow
	//								- then save to chat file
	//						  if file attached -
	//						  	b. save the file - with name as loaded.name in sent folder
	//											 -  with data as loaded.data
	//
	//

	//err := r.ParseForm() //todo todo
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal("Err in chat : Cannot process form, error - ", err)

	}
	fmt.Println("Chat Form ------ ", r.Form) //

	to := r.Form["peers"]
	if len(to) == 0 {
		to = []string{"All"}
	}
	fmt.Println("To in Chat Form ------ ", to)

	message := r.FormValue("message")
	fmt.Println("Message in Chat Form ------ ", message)

	var chatShow ChatShow // preparing chatshow
	//var fileBytes []byte // for storing file in byte array
	file, handler, ferr := r.FormFile("uploadfile")
	if ferr != nil {

		log.Println("Error in processLoadedOnSubmit func - can also mean no file was attached - Err : ", ferr)
		chatShow = NewChatShow(identity.Id, to, message, "")

	} else {

		saveLoadedOnSubmit(file, handler)

		//NewChatShow(identity identity.Identity, to []string, message string, loadedFileName string)
		chatShow = NewChatShow(identity.Id, to, message, handler.Filename)
		file.Close()
	}

	filehelp.SaveToFile(getChatFIlePath(), string(ChatShowToJSON(&chatShow))) // saving chatShow to chat-id.txt
	fmt.Println("ChatShow JSON being saved - string(ChatShowToJSON(&chatShow)) : ", string(ChatShowToJSON(&chatShow)))

	////////////////////////////////////////
	//prepare and send ChatBeat //////////// todo
	prepareChatBeat(r, message, to, identity, peerDS)

}

// processLoadedOnSubmit func - save to new file if file is present in form request
func saveLoadedOnSubmit(file multipart.File, handler *multipart.FileHeader) {

	cwd, _ := os.Getwd()
	f, err := os.OpenFile(cwd+CHATFILESSENTDIR+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in ceating file to save in sent folder : err - ", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("Error in copying file to new created file - ", err)
		return
	}
}

func prepareChatBeat(r *http.Request, message string, to []string, identity identity.Identity, peerDS peers.Peers) {
	file, handler, ferr := r.FormFile("uploadfile")
	if ferr != nil {
		prepareAndSendChatBeat(LoadedFile{}, message, to, identity, peerDS)
	} else {
		fileName := handler.Filename
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Error in reading bytes to file : Error - ", err)
		}
		log.Println("FILEBYTES : ", fileBytes[:20])

		loadedFile := NewLoadedFile(fileName, fileBytes)
		//fmt.Println("loadedFile being sent : ", loadedFile)

		prepareAndSendChatBeat(loadedFile, message, to, identity, peerDS)
		file.Close()
	}

}

// prepares chatbeat and calls sendchatBeat
func prepareAndSendChatBeat(loadedFile LoadedFile, message string, to []string, identity identity.Identity, peerDS peers.Peers) {

	c := NewChat(identity, to, message, loadedFile)

	sendchatBeat(c, to, peerDS)

}

// sendchatBeat func sends the chat struct made from chat-form submit to specified peers
func sendchatBeat(c Chat, to []string, peersDS peers.Peers) {
	// sendchatBeat sends the chat generated to peers
	go chatBeat(c, to, peersDS)

}

// chatBeat builds addresses to send the chat and then beats it
func chatBeat(c Chat, to []string, peersDS peers.Peers) {

	addrToSend := buildAddrToSend(to, peersDS)

	chatBeatingNow(addrToSend, c, peersDS)

}

// buildAddrToSend func buils a addrlist for a chat
func buildAddrToSend(to []string, peersDS peers.Peers) []string {

	var addrToSend []string
	for id, peer := range peersDS.PeerMap {
		for _, t := range to {
			if id == t {
				addrToSend = append(addrToSend, peer.Addr)
			}
		}
	}
	return addrToSend
}

// chatBeatingNow func beats a chat to all needed peers
func chatBeatingNow(addrToSend []string, c Chat, peersDS peers.Peers) {
	if len(addrToSend) == 0 { // to ALL
		c.To = []string{"All"}
		for _, peer := range peersDS.PeerMap {
			//fmt.Println("in SendMessage ChatToJSON(&c) --  all peers ------>  ", string(ChatToJSON(&c)))
			chatBeatingToAddress(peer.Addr, c)
		}
	} else { // to certain peers
		for _, addr := range addrToSend {
			//fmt.Println("in SendMessage ChatToJSON(&c) --  directed peers ------>  ", string(ChatToJSON(&c)))
			chatBeatingToAddress(addr, c)
		}
	}
}

// chatBeatingToAddress func beats a chat to an address
func chatBeatingToAddress(peerAddr string, c Chat) {
	_, err := http.Post("http://"+peerAddr+"/chat/recv", "json", bytes.NewBuffer([]byte(ChatToJSON(&c))))
	if err != nil {
		log.Println("Error in Send Message : ", err)
	}
}

// Receives ChatBeat
func BeatRecv(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) { //receieve chatbeat from peers

	body := httphelp.ReadHttpRequestBody(r)

	chat := Chat{}
	jerr := json.Unmarshal(body, &chat)
	if jerr != nil {
		log.Println("Error in unmarshaling in BeatRecv : ", jerr)

	}

	//if chat.LoadedFile.FileData != nil { // to save file in recv'ed chat Beat
	if chat.LoadedFile.FileName != "" { // to save file in recv'ed chat Beat
		loadedToSaveName := chat.LoadedFile.FileName
		loadedToSaveData := chat.LoadedFile.FileData
		saveLoadedFile(loadedToSaveName, loadedToSaveData)
	}

	//func NewChatShow(from string, to []string, message string, loadedFileName string) ChatShow
	cs := NewChatShow(chat.From, chat.To, chat.Message, chat.LoadedFile.FileName)
	filehelp.SaveToFile(getChatFIlePath(), string(ChatShowToJSON(&cs)))

}

func saveLoadedFile(name string, data []byte) {
	cwd, _ := os.Getwd()
	// !!! saving the file on chat file system
	f, err := os.OpenFile(cwd+CHATFILESRECVDIR+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in ceating file to save in Recv folder : err - ", err)
	}
	defer f.Close()

	f.Write(data)

}

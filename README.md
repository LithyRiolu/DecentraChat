###### The Lithe Project Development Team

[![GitHub stars](https://img.shields.io/github/stars/lithyriolu/decentrachat)](https://github.com/lithyriolu/decentrachat/stargazers) [![GitHub license](https://img.shields.io/github/license/lithyriolu/decentrachat)](https://github.com/LithyRiolu/DecentraChat/blob/master/LICENSE) ![GitHub stars](https://img.shields.io/badge/version-0.0.1-blueviolet)

## DecentraChat
>Decentralised, P2P, web-based chat application

**DecentraChat** allows you to chat annoymously with other peers. You can send text and files to anyone in the chatroom.

***

### Dependencies

`Go 1.12.9*`

`go get -u github.com/gorilla/mux`

*only tested on 1.12.9*

***

### Usage

#### Running DecentraChat
```bash
$ git clone https://github.com/LithyRiolu/DecentraChat`
$ cd DecentraChat`
$ go run main.go <port>
```
Open your browser and go to `http://localhost:<port>/chat`

#### Deleting chat messages
Chats are automatically read from `chat-<port>.txt` located in `chat/files/`.

To delete the chat, either delete `chat-<port>.txt` or open your terminal;
```bash
$ cd DecentraChat/chat/files/
$ sudo rm -rf chat-<port>.txt
```
The chat will now be cleared from your local machine.

***

### Dedicated port chats
There ports are where you may find *dedicated chat ports* for projects that may want to privatize their chat.

#### Current known dedicated port chats
- Dev testing chat - `:50000`

***

### Contributing & Contributors
DecentraChat is open source and published onto GitHub allowing public contributions from other users. 

- **GoLang** - The backend
- **HTML & CSS** - The frontend

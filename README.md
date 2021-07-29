# A Chat App

This directory contains an attempt at making a simple [IRC-like](https://en.wikipedia.org/wiki/Internet_Relay_Chat) 
application. It's actually so simple, that you can only talk to yourself.

The app can be run locally by creating a `chat/server` instance on a certain port, then spinning up any number of 
`chat/client`s to talk with it.

### Drew's notes

Server needs to
- List available chat rooms
- Create chat rooms
- Allow clients to join chat rooms
- Allow clients to leave chat rooms
- Post to chat rooms, broadcasting posts to all other clients in the room.
- Delete chat rooms


Client needs ability to
- Invoke Chat API via REPL
- Handle posts from other chat room members pushed from server, and update UI accordingly.
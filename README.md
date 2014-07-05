Ultimate Tic Tac Toe
=========

Implementation of ultimate tic-tac-toe game in js with go+gorilla websockets backend.

Check it out [here](http://ultimatic.io)!

TODO:
- random URL generator so people can have private 1-on-1 games on their unique URLs (e.g. ultimatic.io/tAv27Ngm)
- check if other side received the websocket message and repeat if not
- make server-side aware of turns (for now it's handled on the client side only and flipped upon every incoming message)

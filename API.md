# The AASLF REST-API v1
The AASLF REST-API uses HTTP on port 1312 to allow your browser to
communicate with the server using JSON. This document describes all of
the used requests.

## "/new"
Creates a new game.
### Request
Key|Value
---|---
Game|The name for the new game
Player|The user's name
Categories|All the categories that the new game should have
### Response
Key|Value
---|---
Status|"ok" if it worked, "err", if not
ID|The ID of the newly created game
Session|The user's session

## "/join"
Joins an existing game.
### Request
Key|Value
---|---
Game|The game's ID
Player|The user's name
### Response
Key|Value
---|---
Status|"ok" if it worked, "err", if not
Session|The user's session

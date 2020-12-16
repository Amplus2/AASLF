# The AASLF REST-API
The AASLF REST-API uses HTTP on port 1312 to allow your browser to
communicate with the server using JSON. This document describes all of
the used requests.

## "/new"
Creates a new game.
### Request
Key|Value
---
Game|The name for the new game
Player|Your name
Categories|All the categories that the new game should have
### Response
Key|Value
---
Status|0 if it worked, anything else, if not
ID|The ID of the newly created game

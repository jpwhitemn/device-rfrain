# RFRain Device Service
Poll RFRain Smart Reader or Cloud system (by REST) for RFID Tag history.
Each reading for a particular RFID reader contains a tag id/number, subzone, signal strength, access time (UTC), reader id, and data.

## TODO
- deal with multiple tag reads for one request (create multiple events with multiple readings)
- get a new session if session expires
- use asynchronous alert mechanism vs pull
- clean up the code
- add testing
- documentation
- test against real smart reader
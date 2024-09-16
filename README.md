# The first things to look at if it ever breaks
- The refreshToken in lavalink.yml
    - Make a new YouTube burner account lol

HELLO GUYS WELCOME TO NOPE'S EXPLANATION SESSION
ok so like the things in handlers are events

basically the queue is just an array of queued tracks
the play command adds tracks to the queue and emit an event (the AddedToQueueSignal one)
then the handlers/queue handles that (OnAddedToQueue)
which basically checks for shit and start playing the first track in the queue
and in OnTrackEnded it just pops the queue and shits and then emit the OnAddedToQueue again :good_work:
so this means that the 0th index in the queue is always the currently-playing track
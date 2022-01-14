# visits-tracker
A simple visits tracker built using Go and Redis

Let me note that this is a really bad idea if you are not backing up Redis since it lives inside memory. The plus side is that the response time
is extremely low and race conditions are avoided with the Redis INCR operation.

# Heroku link
```
https://murmuring-scrubland-25884.herokuapp.com/
```

# Local Setup

## Pre-requisites
Create config.env at root directory
```
LIMIT=60
LIMIT_DURATION=60
REDIS_ENDPOINT=${REDIS_HOST}:${REDIS_PORT}
REDIS_PASSWORD=${REDIS_PASSWORD}
```

# In-Memory Database
This Request limiter use Redis as its IMDB to store the request counter. It use the native atomic increment operation to deal with concurrency. We may use built-in data structure such as sync.atomic / sync.mutex on native golang, but it only worked for a single services. May need to use a sticky sessions to let the request know which server this request should go to. So we choose Redis as a its storage since easier to use when it's coming to scaling. Other option like Memcached may also worked, but in terms of data persistence Redis would be better option. 
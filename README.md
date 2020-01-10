POST /register - create new account with specified nick(unique in app), email, and password
    
    Payload: 
        username - some user name
        email - user email address
        password - some password 
    
    Result: 
        id - primary key
        username - username which you specified in payload 
        email   - user email address which you specified in payload 

POST /login    - accept email and password  and return token, you can use JWT or save session id in database
    
    Payload:
        email - user email address
        password - some password 

    Result:
        token - jwt token or session id 

POST /subscribe - add account with login to your subscription list, you start seeing his tweets in your feeds 
    
    Payload: 
        nickname - nick name for account for which you want to subscribe 

POST /tweets   - create a tweet, account id should be found from JWT or fetched from database using session id
    
    Payload: 
        message - some tweet message

    Result:
        id - message primary key 
        message - tweet message

GET    /tweets   -  return all tweets from your subscriptions 
    
    Result:
        tweets -  all tweets from your subscriptions 
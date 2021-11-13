# Chat Bot Challenge

## Description
This project is designed to test your knowledge of back-end web technologies, specifically in Go
and assess your ability to create back-end products with attention to details, standards, and
reusability.

# Technologies

### Back-End:

- Golang
- MySQL
- Docker

### Front-End:

- vainillaJS

### Libraries:

- Fiber
- Gorm
- Gorilla
- Crypto

# Mandatory Features

- [x] Allow registered users to log in and talk with other users in a chatroom.
- [x] Allow users to post messages as commands into the chatroom with the following format /stock=stock_code .
- [x] Create a decoupled bot that will call an API using the stock_code as a parameter (https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv, here aapl.us is the stock_code).
- [x] The bot should parse the received CSV file and then it should send a message back into the chatroom using a message broker like RabbitMQ. The message will be a stock quote using the following format: “APPL.US quote is $93.42 per share”. The post owner will be the bot.
- [x] Have the chat messages ordered by their timestamps and show only the last 50 messages.
- [ ] Unit test the functionality you prefer.

# Bonus (Optional)

- [x] Have more than one chatroom.
- [x] Handle messages that are not understood or any exceptions raised within the bot.

# Project setup

### Runtime and SDKs

- Download and install Docker (https://docs.docker.com/get-docker/)
- Download and install Golang (https://golang.org/dl/)

### Dependencies

- Start MySQL docker container:

    ````
	docker run --name chatdb -e MYSQL_ROOT_PASSWORD=root -d -p 3306:3306 mysql:lastest
	````

## Run projects

After install the dependencies and make sure the docker containers are running now we can continue to run the projects in the following order:

### **WebApi**

- Open the terminal and go to the following repository directory `chatbot-challenge/cmd/api/`

- To Start the web API, run this command:

    ````
	go run main.go
	````

After executing the previous steps, you are be able to consume the api.

### **WebChat**

- Open the terminal and go to the following repository directory `chatbot-challenge/cmd/chat/`

- To Start the web Chat, run this command:

    ````
	go run main.go
	````

After executing the previous steps, if you click [here](http://localhost:8080) or write in the browser "http://localhost:8080" you should be able to see the web app running.


# Screenshots

![Chat LoginPage](https://github.com/fehepe/chatbot-challenge/blob/main/img/LoginPage.png?raw=true)
![Chat RegisterPage](https://github.com/fehepe/chatbot-challenge/blob/main/img/RegisterPage.png?raw=true)
![Chat RoomPage](https://github.com/fehepe/chatbot-challenge/blob/main/img/RoomPage.png?raw=true)
![Chat ChatPage](https://github.com/fehepe/chatbot-challenge/blob/main/img/ChatPage.png?raw=true)
![Chat 2ClientChat](https://github.com/fehepe/chatbot-challenge/blob/main/img/2ClientChat.png?raw=true)
![Chat BotCmds](https://github.com/fehepe/chatbot-challenge/blob/main/img/BotCmds.png?raw=true)

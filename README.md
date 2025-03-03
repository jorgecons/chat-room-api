# chat-room-api

## 👨‍💻 Author

**Jorge Luis Constable Giménez**  
📧 Email: `jorge.constable@gmail.com`

## SocketListener,RabbitMQ Consumer with Gin API

This project is an API written in Golang that:

- Creates user accounts.
- Allows registered users to log in to a specific room and chat with others.
- When a user log-in, retrieves the last 50 messages.
- Allows users to chat with other users in one room.
- Allows users to post command as "/stock=stock_code". This message won't be saved.
- Has a decoupled bot that listens command and retrieves the stock price and send the message back by a queue. This messages will be saved.
- If the bot has and error, publishes a message error. This messages won't be saved.

## 🛠️ Technologies Used

- **Golang**
- **RabbitMQ** (AMQP)
- **PostgreSQL**
- **Gin** (HTTP API)
- **Resty** (HTTP client)
- **JWT** (token validator)

## 📡 API Endpoints

| Method   | Endpoint | Description                             |
|----------|----------|-----------------------------------------|
| **POST** | `/users` | create a new account                    |
| **POST** | `/login` | login the user and return a valid token |

## 🚀 Running

### 1️⃣ **Run App**

```sh
docker-compose up --build  
```
Build and run a RabbitMQ service and a PostgreSQL. Also, build and run the App.


## 💻 Local Testing

### 1️⃣ **Create New Account**

```sh
curl --location 'localhost:8080/users' \
--header 'Content-Type: application/json' \
--data '{
    "username":"new_user",
    "password":"password"
}'
```

### 2️⃣ **Start Chat**

Open in a browser (e.g. Chrome) and navigate to http://localhost:8080

### 3️⃣ **Login**

Login with user recently created. Also add a room name.

### 4️⃣ **Start chatting**

Start chatting with others users. Also, can send stock command to retrieves its price.

## ⚙️ Improvements
- Add validations to support secure passwords.
- Add https and wss to secure the connection avoiding man-in-the-middle exploits.
- Improve login by validating if user is already connected and jwt expiration.
- Improve logs to avoid sharing sensitive information, like db errors.
- Improve room validations and connections.
- Add rooms management.
- Improve frontend.
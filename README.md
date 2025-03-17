# 📢 Redis Pub/Sub with WebSocket and MySQL Persistence

This project demonstrates a real-time messaging system using Redis's Publish/Subscribe (Pub/Sub) feature, WebSocket for client communication, and MySQL for persistent message storage.

## 🚀 Features

* **Real-time Messaging:** Users can publish and subscribe to channels to receive messages in real-time.
* **WebSocket Communication:** Uses WebSockets for efficient, bidirectional communication between the server and clients.
* **Redis Pub/Sub:** Leverages Redis's Pub/Sub mechanism for message broadcasting.
* **MySQL Persistence:** Stores all messages in a MySQL database for historical access.
* **Graceful Shutdown:** The server can be shut down gracefully with `Ctrl+C` or `SIGTERM`.
* **Multiple Channel Subscriptions:** Clients can subscribe to multiple channels simultaneously.


## 🛠️ Technologies Used

* **Go:** The programming language for the backend server.
* **Redis:** In-memory data structure store used for Pub/Sub.
* **MySQL:** Relational database for persistent message storage.
* **Gorilla WebSocket:** Go library for WebSocket implementation.
* **go-redis:** Go client for interacting with Redis.
* **GORM:** ORM library for interacting with MySQL.

## ⚙️ Setup and Installation

1. **Prerequisites:**
    * Go (latest stable version recommended)
    * Docker (for running Redis) (Recommended)
    * Redis server (running on `localhost:6379` by default) 
    * MySQL server (running on `localhost:3306` with a `chatdb` database)
    * Postman Or Any api client for testing Api


2.**Docker For Redis and Mysql:**
```
    bash 
    docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack (Redis Setup)
    docker run -d --name mysql_chat \
        -e MYSQL_ROOT_PASSWORD=rootpassword \
        -e MYSQL_DATABASE=chatdb \
        -e MYSQL_USER=chatuser \
        -e MYSQL_PASSWORD=chatpassword \
        -p 3306:3306 \
        --restart always \
    mysql:latest
```


3. **Clone the Repository:**
    ```bash
    git clone https://github.com/Amul-Thantharate/Api-Pub-Sub.git
    cd Api-Pub-Sub
    ```

4. **Database Setup:**
    * Create a database named `chatdb` in your MySQL server.
    * The application will automatically create the `messages` table when it starts.
    * Make sure you have a user `nora` with password `root` with access to the `chatdb` database.  You can change this in the `main.go` file.

5. **Install Dependencies:**
    ```bash
    go mod tidy
    ```

6. **Run the Application:**
    ```bash
    go run main.go
    ```
    You should see the message `Starting server on :8080` in your terminal.

## 🕹️ Usage

### Publishing Messages

* **Endpoint:** `POST /publish`
* **Parameters:**
    * `channel`: The channel to publish the message to.
    * `message`: The message content.
* **Example (using `curl`):**
    ```bash
    curl -X POST -d "channel=general&message=Hello, world!" http://localhost:8080/publish
    ```
    This will publish "Hello, world!" to the "general" channel and store it in MySQL.

### Subscribing to Messages

* **Endpoint:** `/subscribe`
* **Query Parameters:**
    * `channel`: One or more channels to subscribe to (e.g., `/subscribe?channel=general&channel=news`).
* **Example (using a WebSocket client):**
    * Connect to `ws://localhost:8080/subscribe?channel=general&channel=news`  (to subscribe to both "general" and "news" channels).
    * You will receive messages from both channels in real-time.
* **Example (using a browser):** You can use a browser extension like "Simple WebSocket Client" to connect to the websocket.


## 🛑 Stopping the Server

* Press `Ctrl+C` in the terminal where the server is running.
* The server will shut down gracefully, closing the Redis and MySQL connections.


## Troubleshooting

* **Database Connection Errors:** Ensure your MySQL server is running, the database `chatdb` exists, and the user `nora` has the correct permissions. Check the MySQL connection string in `main.go`.
* **Redis Connection Errors:** Make sure your Redis server is running on the specified address and port (`localhost:6379`).

## Contributing 🤝

We welcome contributions!  Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Make your changes and commit them (`git commit -m "Your commit message"`).
4. Push your branch to your fork (`git push origin feature/your-feature`).
5. Open a pull request.


## License 📄

This project is licensed under the Apache License 2.0.

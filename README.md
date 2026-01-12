# News Feed

News Feed is a backend service designed to facilitate a social media-like experience, allowing users to engage in various activities such as creating posts, following other users, adding comments to posts, and more. Built with Golang, Gin, GORM, and PostgreSQL, News Feed provides a robust platform for managing user interactions.

High Level System Design

![image](https://github.com/SanjaySinghRajpoot/newsFeed/assets/67458417/0970f0f2-65db-49d2-b4b7-dcb5b8edbbff)

## Motiviation 

I aim to develop a News Feed engine operating on a Pull Model rather than a Push Model. Additionally, I intend to implement a feature that ensures users are presented with predominantly positive posts, while downgrading negative posts within the system itself.

## News Feed Architecture & Implementation

### Overview
This news feed system implements a **Hybrid Push-Pull Model** with sentiment-based content curation. The architecture combines real-time push updates with on-demand pull retrieval, optimized using Redis caching and Kafka message streaming for scalability.

### How It Works

#### 1. Post Creation Flow (Push Model)
When a user creates a post, the system executes three concurrent operations using goroutines:

```
User Creates Post
    ├── Save to PostgreSQL Database
    ├── Publish to Kafka Topic ("newsfeed")
    └── Send Notifications to Followers
```

**Step-by-step Process:**
- **Database Persistence**: Post is saved to PostgreSQL with user metadata
- **Kafka Producer**: Post is serialized to JSON and published to the "newsfeed" Kafka topic
- **Notification Service**: Asynchronous notifications are sent to relevant users
- **Rate Limiting**: User-based rate limiter prevents spam (configured per user via Redis)

#### 2. News Feed Population (Kafka Consumer)
A dedicated Kafka consumer service continuously listens to the "newsfeed" topic:

```
Kafka Consumer (Separate Process)
    ├── Receives Post from Kafka Topic
    ├── Queries Followers of Post Author
    └── Fan-out: Writes Post to Each Follower's Redis Cache
```

**Key Implementation Details:**
- **Fan-out Write Pattern**: When user A posts, the system finds all followers of A and writes the post to each follower's Redis cache
- **Redis Data Structure**: Uses Redis Hash (HSET) with key format: `{userID}` → `{timestamp}` → `{post_json}`
- **Asynchronous Processing**: Decouples post creation from feed distribution, ensuring low latency for the poster
- **Scalability**: Kafka allows horizontal scaling of consumers for high-throughput scenarios

#### 3. Sentiment Analysis (CRON Job)
A background CRON job runs every 5 seconds to analyze post sentiment:

```
CRON Job (Every 5 seconds)
    ├── Query Posts from Last 24 Hours (Unanalyzed)
    ├── Call Flask Sentiment Analysis API
    ├── Receive Scores (Positive, Neutral, Negative)
    └── Update Post with Prominent Sentiment in PostgreSQL
```

**Sentiment Analysis Features:**
- **ML Model**: Uses Hugging Face CardiffNLP Twitter-RoBERTa transformer model
- **Microservice Architecture**: Flask service runs independently on port 5000
- **Sentiment Categories**: Each post receives scores for positive, neutral, and negative sentiment
- **Delayed Processing**: Analysis happens asynchronously to avoid blocking post creation

#### 4. News Feed Retrieval (Hybrid Pull Model)
When a user requests their news feed:

```
User Requests Feed
    ├── Retrieve Cached Posts from Redis (Followers' Posts)
    ├── Query PostgreSQL for Positive Posts (Last 24 Hours)
    ├── Merge & Deduplicate Posts
    ├── Randomize Order
    └── Return to User
```

**Feed Composition Strategy:**
- **Primary Source**: Posts from followed users (cached in Redis)
- **Discovery Content**: Top 10 positive posts from the last 24 hours (from database)
- **Content Curation**: System prioritizes positive sentiment posts
- **Randomization**: Shuffles posts to prevent monotonous ordering
- **Deduplication**: Removes duplicate posts between cached and discovery content

### Technology Stack for News Feed

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Message Queue** | Kafka | Asynchronous post distribution and decoupling |
| **Cache Layer** | Redis | Fast feed retrieval and rate limiting |
| **Database** | PostgreSQL + GORM | Persistent storage with relational integrity |
| **Sentiment Analysis** | Flask + Hugging Face Transformers | ML-based content analysis |
| **Concurrency** | Go Goroutines + WaitGroups | Parallel processing for performance |
| **Scheduling** | Cron (robfig/cron.v2) | Background job execution |

### Pros of This Architecture

#### ✅ Scalability
- **Kafka-based Distribution**: Can handle millions of posts with horizontal consumer scaling
- **Redis Caching**: O(1) read operations for feed retrieval, reducing database load
- **Async Processing**: Non-blocking operations prevent system bottlenecks

#### ✅ Low Latency
- **Fast Post Creation**: User doesn't wait for feed distribution (Kafka handles async)
- **Cached Reads**: Feed retrieval from Redis is near-instantaneous
- **Concurrent Operations**: Goroutines parallelize database writes, Kafka publishes, and notifications

#### ✅ Content Quality
- **Sentiment-Driven Curation**: Prioritizes positive content, improving user experience
- **Discovery Mechanism**: Injects trending positive posts beyond just followed users
- **Mental Health Focus**: Reduces exposure to negative content

#### ✅ Resilience
- **Decoupled Services**: Kafka provides buffer if consumers are temporarily down
- **Persistent Queue**: Messages survive service restarts
- **Independent Sentiment Service**: Flask microservice can scale or fail independently

#### ✅ Developer Experience
- **Clean Separation of Concerns**: Post creation, distribution, and retrieval are isolated
- **Easy Debugging**: Each component (Kafka, Redis, PostgreSQL) can be monitored independently
- **Extensibility**: Easy to add new consumers for analytics, notifications, etc.

## Features

- **User Authentication**: Users can securely sign up and log in to access the News Feed platform.
- **Create Posts**: Users have the ability to create posts to share with their followers and the wider community.
- **Follow Other Users**: Users can follow other users to stay updated with their posts and activities.
- **Add Comments**: Users can engage with posts by adding comments, fostering discussion and interaction within the community.
- **JWT Verification**: Secure JSON Web Token (JWT) verification ensures the integrity and authenticity of user interactions.
- **Persistent Storage**: News Feed utilizes PostgreSQL along with GORM for efficient and reliable data storage.
- **Field Validation**: Input fields are validated to maintain data integrity and prevent errors.
- **Structured Logging**: Logging is structured for easier monitoring, debugging, and analysis.
- **Pagination**: Pagination is implemented to efficiently handle large volumes of data and improve performance.
- **Redis Cache**: Redis Cache is used for smooth login process and improved response time.
- **Rate Limiter**: A User ID based rate limiter has been added to limit multiple posts by the same user. 
- **Sentiment Analysis**: Implemented a Flask-based backend service leveraging Hugging Face's Transformer library for sentiment analysis, seamlessly integrating with a pre-trained CardiffNLP Twitter-RoBERTa model. 
- **Push Based Model**: Our system utilizes a push-based model for populating the cache, ensuring that users have access to the latest posts in real-time.

## Technologies Used

- **Golang**: The core programming language used for building the backend logic.
- **Gin**: A lightweight HTTP web framework for building APIs in Go.
- **Redis Cache**: A in-memory data store used as a cache, vector databases.
- **GORM**: An ORM library for Go, used for interacting with the PostgreSQL database.
- **PostgreSQL**: A powerful, open-source relational database management system.
- **Docker Compose**: Docker Compose is used for defining and running multi-container Docker applications.

## Getting Started

To get started with News Feed, follow these steps:

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Configure the environment variables, database settings, and any other necessary configurations.
4. Build and run the Docker containers using `docker-compose up` for postgres image
5. Go to the root repositry and run `go run main.go` command this will start the API server
5. Access the API endpoints using the provided routes.

## API Endpoints

- **POST api/signup**: Sign up a new user.
- **POST api/login**: Log in an existing user.
- **POST api/posts/create**: Create a new post.
- **GET api/posts**: Get a list of posts.
- **GET /posts/:id**: Get details of a specific post.
- **POST /posts/comment**: Add a comment to a post.
- **POST /users/:id/follow**: Follow a user.
- **GET /users/:id/posts**: Get posts from a specific user.


## Contributing

Contributions are welcome! If you'd like to contribute to News Feed, please follow the standard GitHub flow:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and ensure the code is properly tested.
4. Submit a pull request detailing the changes made and any relevant information.

## License

This project is licensed under the [MIT License](LICENSE).


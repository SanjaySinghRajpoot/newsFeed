# News Feed

News Feed is a backend service designed to facilitate a social media-like experience, allowing users to engage in various activities such as creating posts, following other users, adding comments to posts, and more. Built with Golang, Gin, GORM, and PostgreSQL, News Feed provides a robust platform for managing user interactions.

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

## Technologies Used

- **Golang**: The core programming language used for building the backend logic.
- **Gin**: A lightweight HTTP web framework for building APIs in Go.
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

## Support

For any questions, issues, or support, please [open an issue](https://github.com/example/news-feed/issues) on GitHub.

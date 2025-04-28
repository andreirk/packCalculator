# Pack Calculator

Pack Calculator is a web application that allows users to calculate optimal pack sizes for orders. It consists of a React-based frontend and a Go-based backend.

## Prerequisites

- Docker installed on your system
- Docker Hub account (optional, for pushing the image)

## Build and Run with Docker

### 1. Build the Docker Image

To build the Docker image, run the following command in the root directory of the project:

```bash
   docker build -t pack-calculator .

```

### 2. Run the Docker Container
   Run the container and expose it on port 8080:

```bash
   docker run -p 8080:8080 pack-calculator
```

###  3. Access the Application
   Open your browser and navigate to:

```bash
   http://localhost:8080
```

## Or you can use ready-made working docker image from docker hub

- amsokolov/pack-calculator
## it is public 


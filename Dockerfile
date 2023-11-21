# Use an Ubuntu base image with Node.js 18 pre-installed
FROM node:20-buster


# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package*.json ./

COPY yarn.lock ./

# Install dependencies
RUN yarn

# Copy the rest of the application code into the container
COPY . .

# Build
RUN yarn build

# Copy the rest of the application code into the container
COPY . .

# Install Redis
RUN apt-get update && apt-get install -y redis-server

# Expose the default Playwright browser automation ports
EXPOSE 8080 8081 8000 6379

# Command to run when the container starts
CMD ["sh", "-c", "service redis-server start && npm start"]

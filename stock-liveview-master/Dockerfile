# Use the official Node.js image as the base image
FROM node:22

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the package.json and package-lock.json files to the container
COPY package*.json ./

# Install the dependencies
RUN npm install

# Copy the source code to the container
COPY . .

# Expose the port that the application will run on
EXPOSE 3000

# Define environment variables (example defaults)
ENV MONGODB_URL="mongodb://host.docker.internal:27017,host.docker.internal:27018,host.docker.internal:27019/?replicaSet=rs0"
ENV MONGODB_DB="stockmarket"
ENV MONGODB_COLLECTION="stocks"
ENV NODE_ENV="development"
ENV PORT=3000

# Start the server when the container starts
CMD ["node", "app.js"]

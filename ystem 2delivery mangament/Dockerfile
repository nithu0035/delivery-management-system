# Simple multi-stage Dockerfile for Node.js server
FROM node:18-alpine AS build
WORKDIR /app
COPY server/package*.json ./
RUN npm install --production
COPY server/ .
EXPOSE 4000
CMD ["node", "index.js"]

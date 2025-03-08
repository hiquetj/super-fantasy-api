# Use official MongoDB image as the base
FROM mongo:latest

# Expose the default MongoDB port
EXPOSE 27017

# Set a volume for data persistence (optional)
VOLUME /data/db

# Command to run MongoDB (default behavior of the image)
CMD ["mongod"]
#!/bin/bash

# EC2 details
EC2_USER="ec2-user" # Default user for Amazon Linux 2
EC2_HOST="18.209.59.243"

# Path to your deployment folder on your local system
DEPLOYMENT_PATH="./deployment"
DESTINATION_FOLDER="app"

# SSH Key for accessing EC2 instance, if required
SSH_KEY_PATH="MilestoneCoreStagingInstance.pem"

# SSH Command Prefix
SSH_PREFIX="ssh -i $SSH_KEY_PATH $EC2_USER@$EC2_HOST"

# Check for Docker & Docker Compose installation on EC2
$SSH_PREFIX <<'EOF'
if ! command -v docker &> /dev/null
then
    echo "Installing Docker..."
    sudo yum update -y
    sudo amazon-linux-extras install docker -y
    sudo service docker start
    sudo usermod -a -G docker ec2-user
else
    echo "Docker is already installed."
fi

if ! command -v docker-compose &> /dev/null
then
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
else
    echo "Docker Compose is already installed."
fi
EOF

# Remove all Docker containers and images
$SSH_PREFIX <<'EOF'
echo "Removing all Docker containers and images..."
docker container stop $(docker container ls -aq) 2>/dev/null
docker system prune -a -f --volumes
EOF


# Create the destination folder on the EC2 instance and remove existing content
$SSH_PREFIX <<EOF
echo "Preparing the $DESTINATION_FOLDER folder..."
mkdir -p ~/$DESTINATION_FOLDER
rm -rf ~/$DESTINATION_FOLDER/*
EOF

# Transfer the entire deployment folder to the specified destination folder on EC2 instance
echo "Transferring deployment folder to EC2 instance into the $DESTINATION_FOLDER folder..."
scp -i $SSH_KEY_PATH -r $DEPLOYMENT_PATH/* $EC2_USER@$EC2_HOST:~/$DESTINATION_FOLDER/
# Then, separately transfer all hidden dot files and directories
echo "Transferring hidden files to EC2 instance into the $DESTINATION_FOLDER folder..."
scp -i $SSH_KEY_PATH -r $DEPLOYMENT_PATH/.[!.]* $EC2_USER@$EC2_HOST:~/$DESTINATION_FOLDER/

# Run Docker Compose on EC2
$SSH_PREFIX <<EOF
echo "Pulling, building, and starting containers in the $DESTINATION_FOLDER folder..."
docker-compose -f ~/$DESTINATION_FOLDER/docker-compose.yml pull
docker-compose -f ~/$DESTINATION_FOLDER/docker-compose.yml up --build -d
EOF

echo "Deployment script completed."

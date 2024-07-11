#!/bin/bash

# EC2 details
EC2_USER="ec2-user"
EC2_HOST="54.157.34.211"

# Path to your deployment folder on your local system
DESTINATION_FOLDER="app"
CORE_EXEC_NAME="milestone_core_prod"
CORE_EXEC_LOG_FILE="core_exec.log"

# SSH Key for accessing EC2 instance, if required
SSH_KEY_PATH="MilestoneCoreBackends.pem"

# SSH Command Prefix
SSH_PREFIX="ssh -i $SSH_KEY_PATH $EC2_USER@$EC2_HOST"

# Build the server
echo "Building the server..."
make build_for_prod

if [ ! -f $CORE_EXEC_NAME ]; then
    echo "Error: $CORE_EXEC_NAME not found."
    exit 1
fi

# Stop the server on the EC2 instance
$SSH_PREFIX <<EOF
PID=\$(pgrep -f $CORE_EXEC_NAME)
if [ -n "\$PID" ]; then
    echo "Stopping process \$PID..."
    kill \$PID
    echo "Process \$PID stopped."
    echo "Server stopped."
else
    echo "No process found for $CORE_EXEC_NAME."
fi
EOF

# Check for Docker & Docker Compose installation on EC2
$SSH_PREFIX <<'EOF'
if ! command -v docker &> /dev/null
then
    echo "Installing Docker..."
    sudo yum update -y
    sudo yum install -y docker
    sudo service docker start
    sudo usermod -a -G docker ec2-user
else
    echo "Docker is already installed."
fi
EOF

# Install go and git
$SSH_PREFIX <<'EOF'
if ! command -v go &> /dev/null
then
    echo "Installing Go..."
    sudo yum install golang -y
else
    echo "Go is already installed."
fi

if ! command -v git &> /dev/null
then
    echo "Installing Git..."
    sudo yum install git -y
else
    echo "Git is already installed."
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

scp -i $SSH_KEY_PATH -r $CORE_EXEC_NAME $EC2_USER@$EC2_HOST:~/$DESTINATION_FOLDER

# Copy and rename the .env.prod file to .env on the EC2 instance
scp -i $SSH_KEY_PATH -r .env.prod $EC2_USER@$EC2_HOST:~/$DESTINATION_FOLDER/.env

# Run the server on the EC2 instance
$SSH_PREFIX <<EOF
cd ~/$DESTINATION_FOLDER
chmod +x $CORE_EXEC_NAME
nohup sh -c 'export \$(grep -v '^#' .env | xargs) && ./$CORE_EXEC_NAME' > $CORE_EXEC_LOG_FILE 2>&1 &
EOF

echo "Server started on the EC2 instance. Removing local files..."
rm $CORE_EXEC_NAME

echo "Deployment script completed."

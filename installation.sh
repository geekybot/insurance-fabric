sudo apt update 
sudo snap install -y go --classic 
sudo snap install -y docker 
sudo apt install -y apache2 
sudo apt install -y docker-compose 
sudo apt install -y mongodb
sudo groupadd docker
sudo usermod -aG docker $USER

curl -sSL http://bit.ly/2ysbOFE | bash -s


curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.2/install.sh | bash
source ~/.bashrc
source ~/.profile
nvm install 12.1.0

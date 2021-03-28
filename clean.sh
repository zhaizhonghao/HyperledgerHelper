rm -rf channel/crypto-config
rm -f channel/*.tx
rm -f channel/*.block
rm -f channel/*.yaml
rm -f channel-artifacts/*.block
rm -f *.tar.gz
rm -f docker-compose.yaml
docker rm -f $(docker ps -a -q)
rm log.txt
docker volume prune -f
#!/bin/bash

go build
echo compiling chaincode
sleep 8

docker exec cli peer chaincode install -n teamate -v 0.9 -p github.com/teamate

docker exec cli peer chaincode instantiate -n teamate -v 0.9 -C mychannel -c '{"Args":[]}' -P 'OR("Org1MSP.member")'
sleep 3

docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["registerUser","user1"]}'
echo registering new user : user1
sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["registerUser","user2"]}'
# echo registering new user : user2
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["registerProject","project1"]}'
# echo registering new project : project1
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["registerProject","project2"]}'
# echo registering new project : project2
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["joinProject","user1","project1"]}'
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["joinProject","user2","project1"]}'
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["readUserInfoById","user1"]}'
# echo get user information
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["readProjectInfoByName","project1"]}'
# echo get project information
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["finalizeProject","project1"]}'
# sleep 3

# docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["recordScore","user1","project1","5"]}'
# sleep 3




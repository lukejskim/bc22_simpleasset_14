#!/bin/bash

set -x

# 1. 설치
docker exec cli peer chaincode install -n simpleasset -v 1.1.2 -p github.com/simpleasset/1.1 
docker exec cli peer chaincode list --installed

# 2. 배포
docker exec cli peer chaincode upgrade -n simpleasset -v 1.1.2 -C mychannel -P 'AND ("Org1MSP.member")' -c '{"Args":[]}'
sleep 3
docker exec cli peer chaincode list --instantiated -C mychannel

# ------------------------------ 체인코드 설치 배포 완료
# 확인하는 방법2가지 : query, invoke

# 3. 인보크
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","e","5000"]}'
sleep 3
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","f","1000"]}'
sleep 3

docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["transfer","e","f","3000"]}'
sleep 3

# 4. 쿼리
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","e"]}'
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","f"]}'

# 5. Del
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["del","e"]}'
sleep 3
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","e"]}'






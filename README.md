# 自动生成配置文件的工具
## 功能
本工具提供测试环境下的配资文件的生成，生成的配置文件包括
* crypto-config.yaml
* configtx.yaml
* 与根据crpto-config.yaml生成的身份信息对应的ca、orderer、peer和couchDB的docker-compose.yaml

工具的界面：
![2c5c3c759a4017d2756b090754aef4b4.png](en-resource://database/13414:1)


【注】：
**crypto-config.yaml**
在超级账本中是去中心的身份管理（Decentralized identity management），及要求每个成员组织管理隶属本组织实体的身份。
关于身份管理可以根据不同的需求分为两种环境下的管理方式：
（1）在生成环境中，使用的Fabric-ca进行身份管理。
（2）在测试环境（本文使用的方法）中，则使用*cryptogen*工具，配合*crypto-config.yaml*配置文件，快速地生成必要的身份信息
crypto-config.yaml文件一共有两个部分：OrdererOrgs和PeerOrgs部分。
OrdererOrgs部分定义所有的orderer信息。
PeerOrgs部分中定义了所有成员组织的信息。

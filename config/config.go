package config

const PORT = ":9001"
const PortStream = ":9002"

/*
	文件生成命令
		openssl ecparam -genkey -name secp384r1 -out server.key

		openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
		//填写的服务器名称。必填切需要记住 ，clietn 的时候需要填写这个
		Common Name (eg, fully qualified host name) []:go-grpc-example
*/
const ServerKeyPath = "./config/server.key"
const ServerPemPath = "./config/server.pem"

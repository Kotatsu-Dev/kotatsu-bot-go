#!/bin/bash
cd KotatsuTgBot
go build
cd ..
HTTP_PROXY="socks5h://172.18.208.1:2080" HTTPS_PROXY="socks5h://172.18.208.1:2080" ./KotatsuTgBot/kotatsutgbot 
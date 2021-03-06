#!/bin/sh

#delay=$1
delay=${1-10}
#cert=$2
cert=${2-server.pem}
#key=$3
key=${3-server.key}

#go build -o bench main.go

for proto in http http2 smux yamux ssmux
do
  ./bench -mode server -proto $proto -delay $delay -cert $cert -key $key &
  pid=$!

  for concurrent in 10 50 100 150 200 250 300 350
  do
    for i in `seq 5`
    do
      sleep 10
      ./bench -mode client -concurrent $concurrent -proto $proto -delay $delay
    done
  done

  kill -9 $pid
done

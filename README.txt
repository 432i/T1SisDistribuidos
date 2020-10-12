
#activar server
/sbin/service rabbitmq-server start

#detener server
/sbin/service rabbitmq-server stop


#Todo esto en el PC server {
Add a new/fresh user, say user test and password test:
rabbitmqctl add_user test test

Give administrative access to the new user:
rabbitmqctl set_user_tags test administrator

Set permission to newly created user:
rabbitmqctl set_permissions -p / test ".*" ".*" ".*"
#}

Replace "localhost" with server's IP* (in both send.go and receive.go)
Replace username and password with the ones created above


*Para saber ip ejecutar "ip route"

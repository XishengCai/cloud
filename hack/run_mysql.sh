docker run -d -p 0.0.0.0:3306:3306 \
-v /var/log/mysql:/var/log/mysql \
-v /var/lib/mysql:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=123456 mysql:5.6


echo '----'
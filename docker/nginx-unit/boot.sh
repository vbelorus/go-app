/opt/unit/sbin/unitd --control unix:/var/run/control.unit.sock --user www-data --group www-data

curl -X PUT -d @/app/docker/nginx-unit/config.json --unix-socket /run/control.unit.sock http://localhost/config/

tail -f /var/log/unitd.log
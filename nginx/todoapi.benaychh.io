server {
        server_name todoapi.benaychh.io;

        location / {
                if ($request_method = 'OPTIONS') {
                        add_header 'Access-Control-Allow-Origin' $http_origin;
                        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, DELETE, PATCH';
                        add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
                        add_header 'Access-Control-Max-Age' 1728000;
                        add_header 'Content-Type' 'text/plain; charset=utf-8';
                        add_header 'Content-Length' 0;
                        return 204;
                }
                add_header 'Access-Control-Allow-Origin' $http_origin;
                add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, DELETE, PATCH';
                add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
                add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range';

                proxy_pass http://localhost:8080;

        }

    # Cert info goes here
}
server {
    if ($host = todoapi.benaychh.io) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


        listen 80;
        server_name todoapi.benaychh.io;
    return 404; # managed by Certbot
}

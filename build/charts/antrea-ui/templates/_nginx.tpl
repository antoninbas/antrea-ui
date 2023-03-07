{{- define "antrea-ui.nginx.conf" }}
server {
    listen       3000;
    root /app;
    index index.html;
    client_max_body_size 10M;

    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;

        location /api {
            proxy_http_version 1.1;
            proxy_pass_request_headers on;
            proxy_hide_header Access-Control-Allow-Origin;
            proxy_pass http://127.0.0.1:8080;
        }

        location / {
            try_files $uri $uri/ /index.html;
        }
    }
}
{{- end }}

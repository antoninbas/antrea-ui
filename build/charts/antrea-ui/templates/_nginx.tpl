{{- define "antrea-ui.nginx.conf" }}
server {
    {{- if .Values.https.enabled }}
    listen       3000 ssl;
    {{- else }}
    listen       3000;
    {{- end }}
    server_name _;
    root /app;
    index index.html;
    client_max_body_size 10M;

    {{- if .Values.https.enabled }}
    ssl_certificate /app/nginx-cert.pem;
    ssl_certificate_key /app/nginx-key.pem;
    {{- end }}

    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;

        location /api {
            proxy_http_version 1.1;
            proxy_pass_request_headers on;
            proxy_hide_header Access-Control-Allow-Origin;
            proxy_pass http://127.0.0.1:8080;
            # ensure the correct flags are set, even though the api server should already be setting them
            {{- if .Values.https.enabled }}
            proxy_cookie_flags ~ httponly secure samesite=strict;
            {{- else }}
            proxy_cookie_flags ~ httponly samesite=strict;
            {{- end }}
        }

        location / {
            try_files $uri $uri/ /index.html;
        }
    }
}
{{- end }}

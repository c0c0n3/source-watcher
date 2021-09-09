OSM client HTTP message flows
-----------------------------
> Or what the heck `osmclient` does under the bonnet.


### Getting an auth token

This happens every time you run an `osm` command, i.e. tokens aren't cached!
Example flow

```http
POST /osm/admin/v1/tokens HTTP/1.1
Host: localhost
User-Agent: PycURL/7.43.0.6 libcurl/7.58.0 OpenSSL/1.1.1 zlib/1.2.11 libidn2/2.0.4 libpsl/0.19.1 (+libidn2/2.0.4) nghttp2/1.30.0 librtmp/2.3
Accept: application/json
Content-Type: application/yaml
Content-Length: 65

{"username": "admin", "password": "admin", "project_id": "admin"}
```

```http
HTTP/1.1 200 OK
Server: nginx/1.14.0 (Ubuntu)
Date: Wed, 08 Sep 2021 17:52:11 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 549
Connection: keep-alive
Www-Authenticate: Bearer realm="Needed a token or Authorization http header"
Location: /osm/admin/v1/tokens/TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2
Set-Cookie: session_id=072faf1c629771cdad9133c133fe8bee1202f258; expires=Wed, 08 Sep 2021 18:52:11 GMT; HttpOnly; Max-Age=3600; Path=/; Secure

{
    "issued_at": 1631123531.1251214,
    "expires": 1631127131.1251214,
    "_id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
    "id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
    "project_id": "fada443a-905c-4241-8a33-4dcdbdac55e7",
    "project_name": "admin",
    "username": "admin",
    "user_id": "5c6f2d64-9c23-4718-806a-c74c3fc3c98f",
    "admin": true,
    "roles": [
        {
            "name": "system_admin",
            "id": "cb545e44-cd2b-4c0b-93aa-7e2cee79afc3"
        }
    ],
...
```

### Getting the history of operations on an NS instance

OSM client command

```bash
$ osm ns-op-list ldap
ERROR: ns 'ldap' not found
```

HTTP request

```http
GET /osm/nslcm/v1/ns_instances_content HTTP/1.1
Host: localhost
User-Agent: PycURL/7.43.0.6 libcurl/7.58.0 OpenSSL/1.1.1 zlib/1.2.11 libidn2/2.0.4 libpsl/0.19.1 (+libidn2/2.0.4) nghttp2/1.30.0 librtmp/2.3
Accept: application/json
Content-Type: application/yaml
Authorization: Bearer qIFJhw2JkGbgBToJiuKgYNSKuFgnQlYX
```

HTTP response

```http
.HTTP/1.1 200 OK
Server: nginx/1.14.0 (Ubuntu)
Date: Thu, 09 Sep 2021 14:19:53 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 3
Connection: keep-alive
Set-Cookie: session_id=321df9a60ac919141432e830cfcd8cb306f31877; expires=Thu, 09 Sep 2021 15:19:53 GMT; HttpOnly; Max-Age=3600; Path=/; Secure

[]
```


### Creating a VIM account

OSM client command

```bash
$ osm vim-create --name openvim-site \
    --auth_url http://10.10.10.10:9080/openvim \
    --account_type openvim --description "Openvim site" \
    --tenant osm --user dummy --password dummy
59b92c04-29fa-42a7-923e-63322240b80e
```

HTTP request

```http
POST /osm/admin/v1/vim_accounts HTTP/1.1
Host: localhost
User-Agent: PycURL/7.43.0.6 libcurl/7.58.0 OpenSSL/1.1.1 zlib/1.2.11 libidn2/2.0.4 libpsl/0.19.1 (+libidn2/2.0.4) nghttp2/1.30.0 librtmp/2.3
Accept: application/json
Content-Type: application/yaml
Authorization: Bearer TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2
Content-Length: 196

{"name": "openvim-site", "vim_type": "openvim", "description": "Openvim site", "vim_url": "http://10.10.10.10:9080/openvim", "vim_user": "dummy", "vim_password": "dummy", "vim_tenant_name": "osm"}
```

HTTP response

```http
HTTP/1.1 202 Accepted
Server: nginx/1.14.0 (Ubuntu)
Date: Wed, 08 Sep 2021 17:52:11 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 108
Connection: keep-alive
Location: /osm/admin/v1/vim_accounts/59b92c04-29fa-42a7-923e-63322240b80e
Set-Cookie: session_id=4cd3ace1f2635ca888bbbb6d24a5905540345809; expires=Wed, 08 Sep 2021 18:52:11 GMT; HttpOnly; Max-Age=3600; Path=/; Secure

{
    "id": "59b92c04-29fa-42a7-923e-63322240b80e",
    "op_id": "59b92c04-29fa-42a7-923e-63322240b80e:0"
}
```

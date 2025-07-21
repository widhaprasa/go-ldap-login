# LDAP Login REST API in Go

Simple REST API for LDAP user authentication using Go.  
Supports searching for user DN before binding and optional CORS enabled via environment variable.

---

## Features

- LDAP authentication with search + bind approach  
- Configurable LDAP connection and search filter via environment variables  
- Optional CORS support toggle with `ENABLE_CORS` env var  
- JSON API: POST `/login` with `{ "username": "...", "password": "..." }`  

---

## Requirements

- Go 1.18+  
- LDAP server accessible from the app  
- (Optional) `.env` file for local development  

---

## Installation & Usage

1. Clone or create your project directory:

   ```bash
   mkdir ldaplogin
   cd ldaplogin
   ```

2. Initialize Go module:

   ```bash
   go mod init ldaplogin
   ```

3. Create `main.go` and add the source code.

4. (Optional) Create `.env` file with environment variables:

   ```env
   ENABLE_CORS=true
   LDAP_HOST=ldap.example.com
   LDAP_PORT=389
   LDAP_BASE_DN=ou=People,dc=example,dc=com
   LDAP_BIND_DN=cn=manager,ou=ServiceAccounts,dc=example,dc=com
   LDAP_BIND_PASSWORD=your-manager-password
   LDAP_SEARCH_FILTER=(uid=%s)
   ```

5. Get dependencies:

   ```bash
   go get github.com/go-ldap/ldap/v3
   go get github.com/joho/godotenv
   ```

6. Run the server:

   ```bash
   go run main.go
   ```

The API listens on port 8080 by default.

---

## API Usage

### POST `/login`

**Request JSON:**

```json
{
  "username": "alice",
  "password": "yourpassword"
}
```

**Response JSON:**

- Success (HTTP 200 OK):

  ```json
  {
    "success": true,
    "message": "Login successful"
  }
  ```

- Failure (HTTP 401 Unauthorized):

  ```text
  Invalid credentials
  ```

---

## Environment Variables

| Variable             | Description                                  | Example                                                      | Default    |
|----------------------|----------------------------------------------|--------------------------------------------------------------|------------|
| ENABLE_CORS          | Enable CORS headers (true or false)          | true                                                         | false      |
| LDAP_HOST            | LDAP server hostname or IP                   | ldap.example.com     | (required) |
| LDAP_PORT            | LDAP server port                             | 389                                                          | (required) |
| LDAP_BASE_DN         | Base DN for searching users                  | ou=People,dc=example,dc=com                             | (required) |
| LDAP_BIND_DN         | Bind DN of the LDAP user for searching       | cn=manager,ou=ServiceAccounts,dc=example,dc=com       | (required) |
| LDAP_BIND_PASSWORD   | Password for the bind DN                     | your-manager-password                                        | (required) |
| LDAP_SEARCH_FILTER   | LDAP search filter, %s replaced by username  | (uid=%s)                                                     | (uid=%s)   |

---

## Notes

- The server uses the "search then bind" pattern to authenticate users securely.
- CORS is disabled by default; enable it by setting `ENABLE_CORS=true`.
- Adjust LDAP environment variables according to your LDAP server and schema.
- For production, consider running behind TLS/HTTPS and securing the bind credentials.

---

## License

MIT License

---

## Author

Created by widhaprasa (widhaprasa@gmail.com)
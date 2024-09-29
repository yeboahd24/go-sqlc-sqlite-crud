# Create a new user

```bash
```curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

# Read all users

```bash
curl http://localhost:8080/users
```

# Read a specific user (replace 1 with the actual user ID)

```bash
curl <http://localhost:8080/users/1>
```

# Update a user (replace 1 with the actual user ID)

```bash
curl -X PUT <http://localhost:8080/users/1> \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'
```

# Delete a user (replace 1 with the actual user ID)

```bash
curl -X DELETE <http://localhost:8080/users/1>
```

# go-auth-server
Auth micro-service

### Roadmap

1. Routes
  - `POST /users` - create new user
  - `PATCH /users` - update a user
  - `DELETE /users` - delete a user
  - `POST /auth` - get JWT token
  - `DELETE /auth` - drop JWT token
  - `UPDATE /auth` - renewal JWT token
  
- Give token at ~5mins.
- Auto renewal token after N mins
- When you try to re-extend the token, block token

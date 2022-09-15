# rest-api-tutorial

# users

# REST API

GET /users -- list for users -- 200, 404, 500
GET /users/:id -- user by id -- 200, 404, 500
POST /users/:id -- create user -- 204, 4xx, Header Location: url 
PUT /users/:id -- fully update user -- 204/200, 404, 500
PATCH /users/:id -- partially update user -- 204/200, 404, 500
DELETE /users/:id -- delete user by id -- 204, 404, 400

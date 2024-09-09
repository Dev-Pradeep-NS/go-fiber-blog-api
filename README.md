curl -X POST http://localhost:8000/register \
-H "Content-Type: application/json" \
-d '{
  "username": "newuser",
  "password": "securepassword123"
}'


curl -X POST http://localhost:8000/login \
-H "Content-Type: application/json" \
-d '{
  "username": "pradeep",
  "password": "securepassword123"
}'


curl -X GET http://localhost:8000/albums \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg"


curl -X POST http://localhost:8000/albums \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg" \
-H "Content-Type: application/json" \
-d '{
  "title": "New Album",
  "artist": "New Artist",
  "price": 19.99
}'


curl -X GET http://localhost:8000/albums/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg"


curl -X PUT http://localhost:8000/albums/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg" \
-H "Content-Type: application/json" \
-d '{
  "title": "Updated Album",
  "artist": "Updated Artist",
  "price": 24.99
}'


curl -X DELETE http://localhost:8000/albums/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg"


// Posts -----

curl -X GET http://localhost:8000/posts \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg"


curl -X POST http://localhost:8000/posts \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg" \
-H "Content-Type: application/json" \
-d '{
  "title": "New title",
  "content": "New content"
}'


curl -X GET http://localhost:8000/posts/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg"


curl -X PUT http://localhost:8000/posts/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg" \
-H "Content-Type: application/json" \
-d '{
 "title": "Updated Title",
  "content": "Updated content"
}'


curl -X DELETE http://localhost:8000/posts/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5NzE4NTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoicHJhZGVlcCJ9.aSW-2osyEIc-awTdfgSoTRtKTo6kDqS8WnD54Uwb7Rg"

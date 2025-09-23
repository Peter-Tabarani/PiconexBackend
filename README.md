-- Help to use this code --

1. You have to be a user in the system.
2. When you sign up, a value in the users table is created with your id, password_hash, and role
3. Third, when you login, that email and password is checked in the table and then you are issued a JWT token. After than you must include this token in every curl request from now on. Ask Chat about this.
4. The token expires after an hour, so then you have to re-login and get another token.

-- USEFUL COMMANDS --

To login to admin 3:
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" \
  -d '{
    "email": "david.brown3@example.com",
    "password": "secret123"
  }'

To login to student 122:
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice.anderson@university.edu",
    "password": "secret123"
  }'

To send a request:
curl -X GET http://localhost:8080/student \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NTg1ODg5MTcsImlhdCI6MTc1ODU4NTMxN30.3QapKom4kU2uEXNWXM_by2wML8M-tAEXqzOI0Yr8z1w"

CURRENT ADMIN TOKEN: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NTg1ODg5MTcsImlhdCI6MTc1ODU4NTMxN30.3QapKom4kU2uEXNWXM_by2wML8M-tAEXqzOI0Yr8z1w

CURRENT STUDENT TOKEN: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjIsInJvbGUiOiJzdHVkZW50IiwiZXhwIjoxNzU4NTg4Nzg4LCJpYXQiOjE3NTg1ODUxODh9.GrkAgTQaXtOx8hwwF6XsiqE1Vdp0PqQ_rL9qut-QxqI

curl -X POST http://localhost:8080/specific-documentation \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjIsInJvbGUiOiJzdHVkZW50IiwiZXhwIjoxNzU4NTg4Nzg4LCJpYXQiOjE3NTg1ODUxODh9.GrkAgTQaXtOx8hwwF6XsiqE1Vdp0PqQ_rL9qut-QxqI" \
  -d '{
    "id": 122,
    "doc_type": "Medical",
    "date": "2025-09-22",
    "time": "09:00",
    "file": "base64encodedfile"
  }'

curl -X GET http://localhost:8080/specific-documentation \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NTg1ODg5MTcsImlhdCI6MTc1ODU4NTMxN30.3QapKom4kU2uEXNWXM_by2wML8M-tAEXqzOI0Yr8z1w"

curl -X GET http://localhost:8080/specific-documentation/220 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjIsInJvbGUiOiJzdHVkZW50IiwiZXhwIjoxNzU4NTg4Nzg4LCJpYXQiOjE3NTg1ODUxODh9.GrkAgTQaXtOx8hwwF6XsiqE1Vdp0PqQ_rL9qut-QxqI"

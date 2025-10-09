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

curl -X POST http://localhost:8080/point-of-contact \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer superkey" \
 -d '{
"event_datetime": "2025-10-07T13:30:00Z",
"duration": 30,
"event_type": "trad",
"student_id": 15
}'

curl -X GET http://localhost:8080/student \
 -H "Authorization: Bearer superkey"

curl -X GET http://localhost:8080/personal-documentation \
 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjIsInJvbGUiOiJzdHVkZW50IiwiZXhwIjoxNzU4NTg4Nzg4LCJpYXQiOjE3NTg1ODUxODh9.GrkAgTQaXtOx8hwwF6XsiqE1Vdp0PqQ_rL9qut-QxqI"

curl -X DELETE http://localhost:8080/personal-documentation/228 \
 -H "Authorization: Bearer superkey"

curl -X PUT http://localhost:8080/point-of-contact/24 \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer superkey" \
 -d '{
"event_datetime": "2025-09-22T13:30:00Z",
"duration": 60,
"event_type": "trad",
"id": 15
}'

curl -X POST http://localhost:8080/admin \
 -H "Authorization: Bearer superkey" \
 -H "Content-Type: application/json" \
 -d '{
"first_name": "Michael",
"last_name": "Anderson",
"email": "mike@example.com",
"phone_number": "610-555-4411",
"sex": "male",
"birthday": "1990-11-11",
"address": "77 University Dr",
"city": "Bethlehem",
"state": "PA",
"zip_code": "18015",
"country": "USA",
"preferred_name": "Mike",
"pronouns": "he/him",
"gender": "Male",
"password": "securepassword123",
"title": "Director"
}'

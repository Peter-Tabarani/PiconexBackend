-- Help to use this code --

1. You have to be a user in the system.
2. When you sign up, a value in the users table is created with your id, password_hash, and role
3. Third, when you login, that email and password is checked in the table and then you are issued a JWT token. After than you must include this token in every curl request from now on. Ask Chat about this.
4. The token expires after an hour, so then you have to re-login and get another token.

# Auth Assessment

## Assessment 1

[Short Answer](short-answer.md)

### Success Criteria

- Describe authentication
- Describe authorization
- Describe why passwords need to be hashed in the database
- Describe various attack vectors, such as session hijacking, cross-site scripting and what measures to take to protect them

## Assessment 2

Develop a server-side Node/Express app with to handle user authentication and authorization

### Success Criteria

- Users can sign up to the app with a unique email
- Users cannot sign up for to the app with a duplicate email
- Users can login to the app with valid email/password
- Users cannot login to the app with a blank or missing email
- Users cannot login to the app with a blank or incorrect password
- There is a resource that can only be seen by logged in users
- There is a resource that can only be seen by a specific user
- There is a resource that has some links and content that only appears when logged in / for certain users

### Getting Started

1. Fork/Clone
1. Install the dependencies
1. Run the server: `npm start`
1. Run the tests: `npm test`
1. Write the code to make the tests pass

### User Schema

| field      | type         | metadata                    |
|------------|--------------|-----------------------------|
| id         | serial       | primary key                 |
| email      | varchar(255) | not null, unique            |
| password   | varchar(255) | not null                    |
| admin      | varchar(255) | not null, defaults to false |
| created_at | timestamp    | not null, defaults to now   |

### Endpoints

| Endpoint        | Method | Payload                    | Action          |
|-----------------|--------|----------------------------|-----------------|
| /auth/register  | POST   | {"email": "some@email.com", "password": "some-password"}                      | Handle form submission |
| /auth/login     | POST   | {"email": "some@email.com", "password": "some-password"} | Handle form submission        |
| /auth/logout    | GET    | n/a | Logout user        |
| /users          | GET    | n/a | Display all non-admin users        |
| /users/admin    | GET    | n/a | Display all users       |

- If a user is logged in, make sure they get redirected if they try to view `/auth/register` or `/auth/login`
- Users must be authenticated to view `/users` and `/auth/logout`
- Users must be authenticated and be an admin to view `/users/admin`

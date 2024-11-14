# Back-end Coding Challenge

This repository contains a backend service with multiple endpoints for user and action management. The application is written in Go and can be run using a Makefile for easy setup and testing.

## Getting Started

### Prerequisites

To run this project, you need to have:

- **Go** (1.18+ recommended)
- **Make**

### Running the Application

The application can be started using the Makefile. Use the following command to start the service:

```bash
make start
```

This will start the application locally on http://localhost:3000.


###  Running Tests

To run all tests, use the following command:

```bash
make test
```

## Endpoints

The backend service has the following endpoints:

### User Endpoint

- **Get User by ID**
  - **URL**: `GET /user/:id`
  - **Description**: Retrieves a user's information based on the provided user ID.
  - **Example**: [http://localhost:3000/user/1](http://localhost:3000/user/1)

### Action Endpoints

- **Get Action Count by User ID**
  - **URL**: `GET /users/:id/actions/count`
  - **Description**: Returns the count of actions taken by the user with the specified ID.
  - **Example**: [http://localhost:3000/users/1/actions/count](http://localhost:3000/users/1/actions/count)

- **Get Next Action Probabilities**
  - **URL**: `GET /actions/:actionType/next`
  - **Description**: Provides the probabilities of the next actions for the specified action type.
  - **Example**: [http://localhost:3000/actions/ADD_CONTACT/next](http://localhost:3000/actions/ADD_CONTACT/next)

- **Get Referral Index**
  - **URL**: `GET /actions/referral`
  - **Description**: Fetches the referral index.
  - **Example**: [http://localhost:3000/actions/referral](http://localhost:3000/actions/referral)

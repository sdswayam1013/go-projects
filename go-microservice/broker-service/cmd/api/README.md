# Go Microservices Project

## Overview

This project demonstrates a modern microservices architecture built entirely in Go (Golang). Instead of using a traditional monolithic application, the system is divided into a collection of small, independent, and loosely coupled services that communicate with one another using multiple communication patterns including REST, RPC, gRPC, and AMQP messaging.

The goal of this project is to showcase how distributed applications can be designed, developed, deployed, and scaled using industry-standard technologies and best practices.

---

## Why Microservices?

Traditionally, web applications were developed as monolithic applications where a single codebase handled all business functionality such as authentication, logging, email notifications, and more.

As applications grow in complexity, maintaining and scaling a monolithic architecture becomes increasingly challenging. Microservices address these challenges by breaking an application into smaller, independently deployable services.

### Key Characteristics

* Maintainable and testable
* Loosely coupled architecture
* Independently deployable services
* Organized around business capabilities
* Owned and maintained by small teams
* Improved scalability and fault isolation

---

## System Architecture

The application consists of the following microservices:

### Front-End Service

Responsible for rendering and serving web pages to users.

### Authentication Service

Handles user authentication and authorization.

**Database:** PostgreSQL

### Logging Service

Stores and manages application logs.

**Database:** MongoDB

### Listener Service

Consumes messages from RabbitMQ and performs background processing tasks.

### Broker Service

Acts as an optional API gateway and single entry point into the microservice ecosystem.

### Mail Service

Receives JSON payloads, generates formatted emails, and sends them to recipients.

---

## Communication Patterns

The microservices communicate using multiple approaches:

* REST APIs
* RPC (Remote Procedure Calls)
* gRPC
* AMQP Messaging (RabbitMQ)

This demonstrates different techniques commonly used in distributed systems and enterprise applications.

---

## Technology Stack

### Backend

* Go (Golang)

### Databases

* PostgreSQL
* MongoDB

### Messaging

* RabbitMQ

### Containerization

* Docker
* Docker Swarm

### Orchestration

* Kubernetes

---

## Project Structure

```text
microservices-project/
│
├── authentication-service/
├── broker-service/
├── frontend-service/
├── listener-service/
├── logging-service/
├── mail-service/
│
├── docker-compose.yml
├── swarm.yml
├── kubernetes/
└── README.md
```

---

## Features

* User Authentication
* Distributed Logging
* Email Notifications
* Message Queue Processing
* Service-to-Service Communication
* Containerized Deployment
* Horizontal Scaling
* Independent Service Updates
* High Availability

---

## Running the Project

### Prerequisites

* Go 1.22+
* Docker Desktop
* Docker Compose
* PostgreSQL
* MongoDB
* RabbitMQ
* Kubernetes (optional)

### Clone Repository

```bash
git clone https://github.com/<your-username>/go-microservices.git
cd go-microservices
```

### Start Services

```bash
docker-compose up -d
```

### Verify Running Containers

```bash
docker ps
```

---

## Deployment

### Docker Swarm

Deploy the entire microservice stack to Docker Swarm:

```bash
docker stack deploy -c swarm.yml microservices
```

### Kubernetes

Apply Kubernetes manifests:

```bash
kubectl apply -f kubernetes/
```

---

## Scaling

Scale individual services independently:
`
### Docker Swarm

```bash
docker service scale authentication-service=3
```

### Kubernetes

```bash
kubectl scale deployment authentication-service --replicas=3
```

---

## Learning Objectives

This project demonstrates:

* Microservice Architecture Design
* Go-Based Service Development
* REST API Development
* gRPC Communication
* RabbitMQ Messaging
* PostgreSQL Integration
* MongoDB Integration
* Docker Containerization
* Docker Swarm Deployment
* Kubernetes Deployment
* Service Scaling Strategies

---

## Author

Built as a hands-on learning project to explore modern distributed systems and microservice development using Go.

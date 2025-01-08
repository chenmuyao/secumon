# A REST API + RabbitMQ demo project

(Subject generated by ChatGPT)

### **Comprehensive Scenario: Implementing a Security Event Monitoring and Processing System**

**Background:**
You are tasked with developing a system for a company to monitor API access logs from client systems and detect potential security incidents (e.g., frequent 401 Unauthorized errors). The system should leverage REST API and RabbitMQ for data transmission and processing.

---

### **Problem Requirements:**

1. **Data Collection and Log Simulation**
   - Develop a REST API to receive access logs from client systems. Each log entry contains the following fields:
     ```json
     {
       "timestamp": "2025-01-08T12:00:00Z",
       "client_ip": "192.168.1.1",
       "endpoint": "/api/v1/resource",
       "method": "GET",
       "status_code": 401
     }
     ```
   - Logs are submitted via a `POST /logs` endpoint. The server should validate the log data's format. Invalid logs should return a 400 status code.

2. **Log Processing and Message Queue**
   - Each received log entry should be asynchronously sent to a RabbitMQ queue. Requirements:
     - Queue name: `api-security-logs`.
     - Use a `routing key` to categorize logs (e.g., by `status_code`).

3. **Security Event Detection Service**
   - Implement a consumer service to process messages from the `api-security-logs` queue and identify potential security incidents:
     - **Brute Force Attack**: The same `client_ip` sends 5 consecutive 401 errors within 1 minute.
     - **Abnormal Traffic**: A single `client_ip` makes over 100 requests within 1 minute.
   - If a security incident is detected, log an alert message (e.g., "Potential brute force attack detected, IP: 192.168.1.1").

4. **Extended Features (Optional)**
   - Add a REST API to query detected security events, such as:
     - `GET /alerts`: Returns the 10 most recent security events.
     - `GET /alerts?type=bruteforce`: Filters events by type (e.g., brute force attacks).
   - Store detected security events in memory for query purposes.

---

### **Additional Requirements:**
- **Security**:
  - Validate that the `POST /logs` request contains a valid API key (simulate API authentication).
  - Ensure RabbitMQ uses secure connection settings, such as username and password.
- **High Concurrency**:
  - The log-receiving API should handle 1000 logs per second. You may use techniques such as goroutine pools or queue length limitations.

---

### **Example Scoring Criteria:**
1. **Functionality:**
   - Whether the core functionalities of log reception, queue transmission, and event detection are implemented.
2. **Code Quality:**
   - Whether the code is modular, well-structured, and leverages appropriate Go features (e.g., `goroutines`, `channels`).
3. **Security Design:**
   - Whether the solution effectively simulates API authentication.
   - Whether it addresses queue reliability and exception handling.
4. **Extensibility:**
   - Whether the system can easily incorporate new detection logic for additional security incidents.
5. **Performance Optimization:**
   - Whether the system considers handling high-concurrency scenarios efficiently.

---

### **Sample Outputs**
- Log Reception API Response:
  ```json
  {
    "status": "success",
    "message": "Log received and queued."
  }
  ```

- Alerts from Detection Service:
  ```
  [ALERT] Brute force attack detected from IP: 192.168.1.1
  [ALERT] High traffic detected from IP: 10.0.0.5
  ```

- Security Event Query API Response:
  ```json
  [
    {
      "type": "bruteforce",
      "timestamp": "2025-01-08T12:05:00Z",
      "client_ip": "192.168.1.1",
      "details": "5 consecutive 401 errors within 1 minute"
    }
  ]
  ```

---

### **Study Tips:**
- **REST API:**
  - Learn to validate incoming JSON data (e.g., schema validation).
  - Understand how to implement pagination and filtering in APIs.
- **RabbitMQ:**
  - Familiarize yourself with queue setup, exchanges, and `routing keys`.
  - Study consumer `ack` mechanisms and retry logic.
  - Implement time-window statistics (e.g., IP request frequency) using `sync.Map` or `time.Timer`.

---

This comprehensive scenario is closely related to the business of API security testing. It covers realistic use cases and evaluates your ability to integrate REST APIs, RabbitMQ, and Go effectively.

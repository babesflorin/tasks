# Parity Test Report

**Generated**: 2026-04-07 09:48:03 UTC

**Duration**: 1.065s

## Summary

| Total | Passed | Failed |
|-------|--------|--------|
| 12 | 12 | 0 |

## Results Overview

| # | Test | Method | Path | PHP Status | Go Status | Match |
|---|------|--------|------|------------|-----------|-------|
| 1 | GET /api/task/1 — happy path | `GET` | `/api/task/1` | 200 | 200 | ✅ |
| 2 | GET /api/task/9999999 — not found | `GET` | `/api/task/9999999` | 404 | 404 | ✅ |
| 3 | POST /api/task — create task | `POST` | `/api/task` | 200 | 200 | ✅ |
| 4 | PUT /api/task — update task | `PUT` | `/api/task` | 200 | 200 | ✅ |
| 5 | PUT /api/task — update not found | `PUT` | `/api/task` | 404 | 404 | ✅ |
| 6 | PUT /api/task — validation error (empty body) | `PUT` | `/api/task` | 400 | 400 | ✅ |
| 7 | PUT /api/task — invalid JSON | `PUT` | `/api/task` | 400 | 400 | ✅ |
| 8 | GET /api/task — list all tasks | `GET` | `/api/task` | 200 | 200 | ✅ |
| 9 | PUT /api/task/1/complete — complete task | `PUT` | `/api/task/1/complete` | 200 | 200 | ✅ |
| 10 | PUT /api/task/9999999/complete — complete not found | `PUT` | `/api/task/9999999/complete` | 404 | 404 | ✅ |
| 11 | DELETE /api/task/5 — delete task | `DELETE` | `/api/task/5` | 200 | 200 | ✅ |
| 12 | DELETE /api/task/9999 — delete not found | `DELETE` | `/api/task/9999` | 404 | 404 | ✅ |

## Detailed Results

### 1. GET /api/task/1 — happy path — ✅ PASS

**Request**: `GET /api/task/1`

**PHP Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:03.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Task description 0",
    "done": false,
    "id": 1,
    "name": "Task name 0",
    "updated_at": {
      "date": "2026-04-07 09:48:03.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-07"
  }
}
```

**Go Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Task description 0",
    "done": false,
    "id": 1,
    "name": "Task name 0",
    "updated_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-07"
  }
}
```

---

### 2. GET /api/task/9999999 — not found — ✅ PASS

**Request**: `GET /api/task/9999999`

**PHP Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

**Go Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

---

### 3. POST /api/task — create task — ✅ PASS

**Request**: `POST /api/task`

**Request Body**:
```json
{
  "description": "Parity description",
  "name": "Parity task",
  "when": "2026-04-12"
}
```

**PHP Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:03.848821",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Parity description",
    "done": false,
    "id": 6,
    "name": "Parity task",
    "updated_at": {
      "date": "2026-04-07 09:48:03.850044",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-12"
  }
}
```

**Go Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:03.863970",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Parity description",
    "done": false,
    "id": 6,
    "name": "Parity task",
    "updated_at": {
      "date": "2026-04-07 09:48:03.863970",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-12"
  }
}
```

---

### 4. PUT /api/task — update task — ✅ PASS

**Request**: `PUT /api/task`

**Request Body**:
```json
{
  "description": "Updated desc",
  "id": 2,
  "name": "Updated name",
  "when": "2026-04-10"
}
```

**PHP Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Updated desc",
    "done": false,
    "id": 2,
    "name": "Updated name",
    "updated_at": {
      "date": "2026-04-07 09:48:03.913004",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-10"
  }
}
```

**Go Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Updated desc",
    "done": false,
    "id": 2,
    "name": "Updated name",
    "updated_at": {
      "date": "2026-04-07 09:48:03.927937",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-10"
  }
}
```

---

### 5. PUT /api/task — update not found — ✅ PASS

**Request**: `PUT /api/task`

**Request Body**:
```json
{
  "description": "x",
  "id": 9999999,
  "name": "x",
  "when": "2026-04-10"
}
```

**PHP Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

**Go Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

---

### 6. PUT /api/task — validation error (empty body) — ✅ PASS

**Request**: `PUT /api/task`

**Request Body**:
```json
{}
```

**PHP Response** (status 400, `application/json`):
```json
{
  "data": "",
  "error": "Task is not valid!",
  "messages": [
    "Task name is not valid!",
    "We need an id to know which entity to update!",
    "Task description is not valid!",
    "Task must have a date!"
  ]
}
```

**Go Response** (status 400, `application/json`):
```json
{
  "data": "",
  "error": "Task is not valid!",
  "messages": [
    "Task name is not valid!",
    "We need an id to know which entity to update!",
    "Task description is not valid!",
    "Task must have a date!"
  ]
}
```

---

### 7. PUT /api/task — invalid JSON — ✅ PASS

**Request**: `PUT /api/task`

**Request Body**:
```json
/
```

**PHP Response** (status 400, `application/json`):
```json
{
  "data": "",
  "error": "Request must be json!"
}
```

**Go Response** (status 400, `application/json`):
```json
{
  "data": "",
  "error": "Request must be json!"
}
```

---

### 8. GET /api/task — list all tasks — ✅ PASS

**Request**: `GET /api/task`

**PHP Response** (status 200, `application/json`):
```json
{
  "data": [
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 0",
      "done": false,
      "id": 1,
      "name": "Task name 0",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-07"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 1",
      "done": false,
      "id": 2,
      "name": "Task name 1",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-08"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 2",
      "done": false,
      "id": 3,
      "name": "Task name 2",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-09"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 3",
      "done": false,
      "id": 4,
      "name": "Task name 3",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-10"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 4",
      "done": false,
      "id": 5,
      "name": "Task name 4",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-11"
    }
  ]
}
```

**Go Response** (status 200, `application/json`):
```json
{
  "data": [
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 0",
      "done": false,
      "id": 1,
      "name": "Task name 0",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-07"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 1",
      "done": false,
      "id": 2,
      "name": "Task name 1",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-08"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 2",
      "done": false,
      "id": 3,
      "name": "Task name 2",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-09"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 3",
      "done": false,
      "id": 4,
      "name": "Task name 3",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-10"
    },
    {
      "created_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "description": "Task description 4",
      "done": false,
      "id": 5,
      "name": "Task name 4",
      "updated_at": {
        "date": "2026-04-07 09:48:04.000000",
        "timezone": "UTC",
        "timezone_type": 3
      },
      "when": "2026-04-11"
    }
  ]
}
```

---

### 9. PUT /api/task/1/complete — complete task — ✅ PASS

**Request**: `PUT /api/task/1/complete`

**PHP Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Task description 0",
    "done": true,
    "id": 1,
    "name": "Task name 0",
    "updated_at": {
      "date": "2026-04-07 09:48:04.191433",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-07"
  }
}
```

**Go Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Task description 0",
    "done": true,
    "id": 1,
    "name": "Task name 0",
    "updated_at": {
      "date": "2026-04-07 09:48:04.198961",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-07"
  }
}
```

---

### 10. PUT /api/task/9999999/complete — complete not found — ✅ PASS

**Request**: `PUT /api/task/9999999/complete`

**PHP Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

**Go Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

---

### 11. DELETE /api/task/5 — delete task — ✅ PASS

**Request**: `DELETE /api/task/5`

**PHP Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Task description 4",
    "done": false,
    "id": 5,
    "name": "Task name 4",
    "updated_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-11"
  }
}
```

**Go Response** (status 200, `application/json`):
```json
{
  "data": {
    "created_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "description": "Task description 4",
    "done": false,
    "id": 5,
    "name": "Task name 4",
    "updated_at": {
      "date": "2026-04-07 09:48:04.000000",
      "timezone": "UTC",
      "timezone_type": 3
    },
    "when": "2026-04-11"
  }
}
```

---

### 12. DELETE /api/task/9999 — delete not found — ✅ PASS

**Request**: `DELETE /api/task/9999`

**PHP Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

**Go Response** (status 404, `application/json`):
```json
{
  "data": "",
  "error": "Task not found!"
}
```

---


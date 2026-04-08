# Frontend API Reference

This document describes the HTTP API currently exposed by the Go server in `Server Examples/`.

Base URL:

```text
http://localhost:8080
```

## Authentication Model

Most API endpoints require an authenticated session.

The server accepts authentication in either of these forms:

1. `Authorization: Bearer <token>`
2. `Authorization: <token>`
3. `session_token` cookie

Token lookup order is:

1. `Authorization` header
2. `session_token` cookie

Important behavior:

- The server creates an in-memory session after Google OAuth succeeds.
- Sessions expire after 30 minutes of inactivity.
- `/auth/callback` sets an `HttpOnly` `session_token` cookie.
- If your frontend relies on cookies, send requests with credentials enabled.
- If your frontend stores the returned token, you can also send it explicitly in the `Authorization` header.

Example authenticated fetch:

```js
fetch('/api/user', {
  method: 'GET',
  credentials: 'include',
  headers: {
    Authorization: `Bearer ${token}`,
  },
});
```

## Roles

The server uses these role strings:

- `admin`
- `officer`
- `member`

Role-restricted routes:

- `/api/officer/*` requires role `officer`
- `/api/admin/*` requires role `admin`

## Common Response Notes

- Most successful responses are JSON.
- Some authentication failures use `http.Error(...)` and return plain text, not JSON.
- Some handlers do not strictly enforce HTTP method even though they are clearly intended for `GET` or `POST`. The frontend should still use the method documented below.

## Shared Data Shapes

### User object returned by `/api/user`

```json
{
  "email": "student@example.edu",
  "role": "member"
}
```

### Prefill fields returned by `/api/prefill`

```json
{
  "fields": {
    "email": "student@example.edu",
    "name": "Student Name"
  }
}
```

### Organization object

Used by:

- `/api/results`
- `/api/officer/orgs`
- `/api/admin/orgs`
- `/api/admin/create` request body and response body
- `/api/officer/update` request body

```json
{
  "id": 12,
  "clubName": "CS Club",
  "description": "Club description",
  "meetingTime": "Thursdays at 6:30 PM",
  "imagePath": "/images/cs-club.png",
  "externalLink": "https://example.com",
  "contactInfo": "club@example.edu",
  "includeOfficerEmails": false,
  "updatedAt": "2026-04-06T18:30:00Z",
  "officers": []
}
```

Notes:

- `updatedAt` is returned on read endpoints.
- `officers` exists in the wire type, but the current read handlers do not populate it.
- The current update handler ignores the `officers` field.

### Swipe question object

Returned by `/api/adjectives`.

Shape is partially dynamic because question text is stored in a translation map.

Example:

```json
{
  "id": 4,
  "question_type": "adjective",
  "en": "Creative",
  "es": "Creativo"
}
```

At minimum, each object contains:

- `id: number`
- `question_type: string`
- zero or more translation keys such as `en`

## Endpoint Reference

### `GET /login`

Starts the Google OAuth flow.

Query params:

- `popup=1` optional. Enables popup login behavior.

Request body:

- None.

Success behavior:

- Returns `307 Temporary Redirect` to Google OAuth.

Frontend notes:

- This endpoint is usually navigated to in the browser rather than called with `fetch`.
- If you use popup mode, open `/login?popup=1` in a popup window.

### `GET /auth/callback`

Completes the Google OAuth flow.

Query params expected from Google:

- `state: string`
- `code: string`

Request body:

- None.

Success responses:

1. Normal flow returns JSON:

```json
{
  "message": "Successfully authenticated",
  "email": "student@example.edu",
  "token": "<session-token>"
}
```

2. Popup flow returns HTML, not JSON.

Popup behavior:

- Sets `session_token` cookie.
- Sends `window.opener.postMessage(...)` with:

```json
{
  "type": "google-auth-success",
  "email": "student@example.edu",
  "token": "<session-token>"
}
```

- Attempts to close the popup.

Possible errors:

- `400 Bad Request` plain text: invalid or expired OAuth `state`
- `401 Unauthorized` plain text: token exchange failed or Google subject extraction failed
- `500 Internal Server Error` plain text: unable to fetch or persist user info

### `GET /api/user`

Returns the authenticated user identity and resolved role.

Authentication:

- Required.

Request body:

- None.

Success response:

```json
{
  "email": "student@example.edu",
  "role": "member"
}
```

Possible errors:

- `401 Unauthorized` plain text: missing token
- `401 Unauthorized` plain text: invalid or expired token

Frontend notes:

- Use this after login to determine which UI to show.

### `POST /logout`

Logs out the current session.

Authentication:

- Optional. If a token exists, it is removed.

Request body:

- None.

Success response:

```json
{
  "message": "Logged out successfully"
}
```

Frontend notes:

- The handler currently accepts any HTTP method, but `POST` is the intended method.
- The response also expires the `session_token` cookie.

### `GET /api/prefill`

Returns lightweight profile fields for client-side form prefilling.

Authentication:

- Required.

Request body:

- None.

Success response:

```json
{
  "fields": {
    "email": "student@example.edu",
    "name": "Student Name"
  }
}
```

Possible errors:

- `401 Unauthorized` plain text: missing token
- `401 Unauthorized` plain text: invalid or expired token

Frontend notes:

- Current fields are built from the login profile and typically include `email` and `name`.

### `GET /api/adjectives`

Returns swipe-card question content for `adjective` and `personality_traits` question types.

Authentication:

- Required.

Request body:

- None.

Success response:

```json
[
  {
    "id": 4,
    "question_type": "adjective",
    "en": "Creative"
  },
  {
    "id": 9,
    "question_type": "personality_traits",
    "en": "Collaborative"
  }
]
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `500 Internal Server Error` JSON:

```json
{
  "error": "Failed to fetch adjectives"
}
```

### `POST /response`

Submits one survey/swipe response.

Authentication:

- Required.

Headers:

- `Content-Type: application/json`

Request body:

```json
{
  "questionId": 4,
  "answer": true
}
```

Field rules:

- `questionId` must be a positive integer
- `answer` must be a boolean

Success response:

```json
{
  "email": "student@example.edu",
  "questionId": 4,
  "answer": true,
  "message": "Response accepted"
}
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `405 Method Not Allowed` JSON:

```json
{
  "error": "Only POST method is allowed"
}
```

- `400 Bad Request` JSON:

```json
{
  "error": "Invalid JSON"
}
```

- `400 Bad Request` JSON:

```json
{
  "error": "questionId must be a positive integer"
}
```

Frontend notes:

- The current handler validates and echoes the payload.
- It does not currently persist the response to the database.

### `POST /submit`

Submits the demographics form.

Authentication:

- Usually required.
- Controlled by environment variable `SUBMIT_REQUIRES_AUTH`.
- Default behavior is authenticated-only.

Headers:

- `Content-Type: application/json`

Request body:

```json
{
  "name": "Student Name",
  "gender": "Woman",
  "race": ["Asian"],
  "religion": "None",
  "major": ["Computer Science"]
}
```

Field rules:

- `name` must be non-empty
- `gender` must be non-empty
- `race` must contain at least one value
- `religion` must be non-empty
- `major` must contain at least one value

Success response:

```json
{
  "message": "Demographics submitted successfully",
  "name": "Student Name",
  "gender": "Woman",
  "race": ["Asian"],
  "religion": "None",
  "major": ["Computer Science"],
  "email": "student@example.edu"
}
```

Possible errors:

- `401 Unauthorized` JSON when auth is enabled and no valid session is present
- `405 Method Not Allowed` JSON:

```json
{
  "error": "Only POST method is allowed"
}
```

- `400 Bad Request` JSON:

```json
{
  "error": "Invalid JSON"
}
```

- `400 Bad Request` JSON:

```json
{
  "error": "Missing one or more required demographics fields"
}
```

Frontend notes:

- The current handler validates and echoes the submission.
- It does not currently persist demographics data to the database.
- If auth is disabled through environment configuration, `email` in the success payload may be an empty string.

### `GET /api/results`

Returns organization matches for the authenticated user.

Authentication:

- Required.

Request body:

- None.

Success response:

```json
[
  {
    "id": 0,
    "clubName": "CS Club",
    "description": "placeholder description yay!",
    "meetingTime": "Thursdays at 6:30",
    "imagePath": "",
    "externalLink": "",
    "contactInfo": "",
    "includeOfficerEmails": false,
    "updatedAt": "2026-04-06T18:30:00Z"
  }
]
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `500 Internal Server Error` JSON:

```json
{
  "error": "Failed to fetch user orgs"
}
```

Frontend notes:

- The current implementation returns placeholder club data.
- The matching algorithm is not fully implemented yet.

### `GET /api/officer/orgs`

Returns the organizations managed by the authenticated officer.

Authentication:

- Required.
- User must have role `officer`.

Request body:

- None.

Success response:

```json
[
  {
    "id": 12,
    "clubName": "CS Club",
    "description": "Club description",
    "meetingTime": "Thursdays at 6:30 PM",
    "imagePath": "/images/cs-club.png",
    "externalLink": "https://example.com",
    "contactInfo": "club@example.edu",
    "includeOfficerEmails": false,
    "updatedAt": "2026-04-06T18:30:00Z",
    "officers": []
  }
]
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON:

```json
{
  "error": "Forbidden"
}
```

- `500 Internal Server Error` JSON:

```json
{
  "error": "Failed to fetch officer orgs"
}
```

### `POST /api/officer/update`

Updates one organization owned by the authenticated officer.

Authentication:

- Required.
- User must have role `officer`.

Headers:

- `Content-Type: application/json`

Request body:

```json
{
  "id": 12,
  "clubName": "CS Club",
  "description": "Updated description",
  "meetingTime": "Thursdays at 7:00 PM",
  "imagePath": "/images/cs-club.png",
  "externalLink": "https://example.com",
  "contactInfo": "club@example.edu",
  "includeOfficerEmails": true,
  "updatedAt": "2026-04-06T18:30:00Z",
  "officers": []
}
```

Important behavior:

- `id` is required and must be a positive integer.
- The server updates these fields:
  - `clubName`
  - `description`
  - `meetingTime`
  - `imagePath`
  - `externalLink`
  - `contactInfo`
  - `includeOfficerEmails`
- The server ignores `updatedAt` and `officers` in the current implementation.

Success response:

- Status `200 OK`
- No response body is currently written by the handler.

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON:

```json
{
  "error": "Forbidden"
}
```

- `400 Bad Request` JSON:

```json
{
  "error": "Invalid JSON"
}
```

Frontend notes:

- Because the handler returns an empty body on success, the safest client behavior is to check `response.ok` and then refetch `/api/officer/orgs`.
- The handler is intended for `POST`, although the server does not currently enforce method.

### `GET /api/admin/orgs`

Returns all organizations for admin users.

Authentication:

- Required.
- User must have role `admin`.

Request body:

- None.

Success response:

```json
[
  {
    "id": 12,
    "clubName": "CS Club",
    "description": "Club description",
    "meetingTime": "Thursdays at 6:30 PM",
    "imagePath": "/images/cs-club.png",
    "externalLink": "https://example.com",
    "contactInfo": "club@example.edu",
    "includeOfficerEmails": false,
    "updatedAt": "2026-04-06T18:30:00Z",
    "officers": []
  }
]
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON:

```json
{
  "error": "Forbidden"
}
```

- `500 Internal Server Error` JSON:

```json
{
  "error": "Failed to fetch orgs"
}
```

### `POST /api/admin/create`

Creates a new organization.

Authentication:

- Required.
- User must have role `admin`.

Headers:

- `Content-Type: application/json`

Request body:

```json
{
  "clubName": "New Club",
  "description": "Club description",
  "meetingTime": "Fridays at 5:00 PM",
  "imagePath": "/images/new-club.png",
  "externalLink": "https://example.com/new-club",
  "contactInfo": "new-club@example.edu",
  "includeOfficerEmails": false,
  "officers": []
}
```

Important behavior:

- `clubName` is required and must be non-empty.
- `id`, `updatedAt`, and `officers` from the request body are ignored by the server.

Success response:

- Status `201 Created`
- Response body uses the Organization object shape and includes server-generated `id` and `updatedAt`.

Example success body:

```json
{
  "id": 34,
  "clubName": "New Club",
  "description": "Club description",
  "meetingTime": "Fridays at 5:00 PM",
  "imagePath": "/images/new-club.png",
  "externalLink": "https://example.com/new-club",
  "contactInfo": "new-club@example.edu",
  "includeOfficerEmails": false,
  "updatedAt": "2026-04-06T18:30:00Z",
  "officers": []
}
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON:

```json
{
  "error": "Forbidden"
}
```

- `405 Method Not Allowed` JSON:

```json
{
  "error": "Only POST method is allowed"
}
```

- `400 Bad Request` JSON:

```json
{
  "error": "Invalid JSON"
}
```

- `400 Bad Request` JSON when `clubName` is missing/blank:

```json
{
  "error": "clubName is required"
}
```

## Frontend Integration Checklist

- Use `credentials: 'include'` if you want the browser to send the `session_token` cookie.
- Prefer also storing the returned token from `/auth/callback` so you can send `Authorization: Bearer <token>` explicitly.
- Expect both JSON and plain-text error responses depending on the endpoint.
- Treat `/response`, `/submit`, and `/api/results` as partially implemented backend endpoints for now.
- After calling `/api/officer/update`, refetch data because success returns no body.

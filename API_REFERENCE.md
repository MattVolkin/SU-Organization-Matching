# Frontend API Reference

This document reflects the current HTTP API exposed by the Go server in `server/`.

Base URL:

```text
http://localhost:8080
```

## Authentication

Most endpoints require an authenticated session.

The server accepts auth using either:

1. `Authorization: Bearer <token>`
2. `Authorization: <token>`
3. `session_token` cookie

Token lookup order:

1. `Authorization` header
2. `session_token` cookie

Important behavior:

- `/auth/callback` sets an `HttpOnly` `session_token` cookie.
- Sessions are stored in memory and expire after 30 minutes of inactivity.
- If the frontend uses cookies, requests should include credentials.
- If the frontend stores the token returned by `/auth/callback`, it can also send `Authorization: Bearer <token>` explicitly.

Example authenticated request:

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

Role strings returned by the server:

- `admin`
- `officer`
- `member`

Role-restricted route groups:

- `/api/officer/*` requires `officer`
- `/api/admin/*` requires `admin`

## Common Notes

- Successful API responses are JSON except popup-mode `/auth/callback`, which returns HTML.
- Middleware-based auth failures return JSON.
- `/api/user` and `/api/prefill` use `http.Error(...)` for token failures and therefore return plain text on `401`.
- REST methods are enforced by the router.

## Shared Data Shapes

### User

Returned by `/api/user`.

```json
{
  "email": "student@example.edu",
  "role": "member"
}
```

### Prefill payload

Returned by `/api/prefill`.

```json
{
  "fields": {
    "email": "student@example.edu",
    "name": "Student Name"
  }
}
```

### Organization

Used by:

- `GET /api/results`
- `GET /api/officer/orgs`
- `GET /api/admin/orgs`
- `POST /api/admin/orgs` request and response
- `PATCH /api/officer/orgs` response
- `PATCH /api/admin/orgs` response

```json
{
  "id": 12,
  "clubName": "CS Club",
  "description": "Club description",
  "meetingTime": "Thursdays at 6:30 PM",
  "imagePath": "/api/images/cs-club-12-1712536000000000000.png",
  "externalLink": "https://example.com",
  "contactInfo": "club@example.edu",
  "includeOfficerEmails": false,
  "officers": ["president@example.edu", "vp@example.edu"],
  "personality": ["Welcoming", "Curious"],
  "activities": ["Community Service", "Guest Speakers"],
  "genders": ["Any"],
  "ethnicities": ["Any"],
  "religions": ["Any"],
  "strict_genders": false,
  "dedicated_majors": ["Computer Science"],
  "other": ["Open to first-year students"],
  "updatedAt": "2026-04-09T18:30:00Z"
}
```

Notes:

- `officers` is an array of officer email addresses derived from the club leaders relationship.
- `strict_genders` and `dedicated_majors` intentionally use snake_case in JSON.
- `updatedAt` is returned by read endpoints and create/patch responses.

### Organization PATCH body

Used by:

- `PATCH /api/officer/orgs`
- `PATCH /api/admin/orgs`

The club id is required in the JSON body.

Only include fields that should change:

```json
{
  "id": 12,
  "clubName": "Updated Club Name",
  "description": "Updated description",
  "meetingTime": "Fridays at 5:00 PM",
  "imagePath": "/api/images/new-image.png",
  "externalLink": "https://example.com",
  "contactInfo": "updated@example.edu",
  "includeOfficerEmails": true,
  "officers": ["president@example.edu", "vp@example.edu"],
  "personality": ["Welcoming", "Collaborative"],
  "activities": ["Board Games", "Discussion"],
  "genders": ["Any"],
  "ethnicities": ["Any"],
  "religions": ["Any"],
  "strict_genders": false,
  "dedicated_majors": ["Computer Science"],
  "other": ["No experience required"]
}
```

Rules:

- At least one field must be present.
- Any omitted field is left unchanged.
- If `officers` is included, it replaces the club's full officer list.
- Every officer email provided must already belong to an existing user record.

### Swipe question object

Returned by `/api/adjectives`.

Shape is dynamic because translations come from the database.

Example:

```json
{
  "id": 4,
  "question_type": "adjective",
  "translations": {
    "en": ["Creative", "A trait related to imagination and originality"],
    "es": []
  }
}
```

At minimum, each object contains:

- `id: number`
- `question_type: string`
- `translations: object`

## Endpoint Reference

### `GET /login`

Starts the Google OAuth flow.

Query params:

- `popup=1` optional. Enables popup login behavior.

Success behavior:

- Returns `307 Temporary Redirect` to Google OAuth.

Frontend notes:

- This is normally opened via navigation or popup, not `fetch`.

### `GET /auth/callback`

Completes the Google OAuth flow.

Expected query params from Google:

- `state: string`
- `code: string`

Success response in normal flow:

```json
{
  "message": "Successfully authenticated",
  "email": "student@example.edu",
  "token": "<session-token>"
}
```

Popup flow:

- Sets `session_token` cookie.
- Returns HTML that posts this object back to `window.opener`:

```json
{
  "type": "google-auth-success",
  "email": "student@example.edu",
  "token": "<session-token>"
}
```

Possible errors:

- `400 Bad Request` plain text: invalid or expired `state`
- `401 Unauthorized` plain text: token exchange or user identification failed
- `500 Internal Server Error` plain text: user info fetch/persist failed

### `GET /api/user`

Returns the authenticated user and role.

Authentication:

- Required.

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

### `GET /api/images/{filename}`

Fetches a club image stored in `Components/OrgPhotos`.

Authentication:

- Not required.

Path params:

- `filename: string` exact stored file name

Success behavior:

- Returns the image bytes with the detected file content type.

Possible errors:

- `400 Bad Request` JSON: invalid filename
- `404 Not Found` JSON: image not found
- `500 Internal Server Error` JSON: image could not be accessed

### `POST /api/officer/orgs/{id}/image`

Uploads or replaces the image for an officer-managed club.

Authentication:

- Required.
- Officer role required.
- Officers may only upload images for clubs they lead.

Path params:

- `id: number` club id

Request content type:

- `multipart/form-data`

Required form field:

- `image` file upload

Rules:

- Images are stored in `Components/OrgPhotos`.
- Successful uploads update the club `imagePath` field to `/api/images/{filename}`.
- Supported file types: JPEG, PNG, GIF, WebP.
- Max upload size: 10 MB.

Success response example:

```json
{
  "message": "Club image uploaded successfully",
  "filename": "cs-club-12-1712536000000000000.png",
  "imagePath": "/api/images/cs-club-12-1712536000000000000.png",
  "contentType": "image/png",
  "club": {
    "id": 12,
    "clubName": "CS Club"
  }
}
```

Possible errors:

- `400 Bad Request` JSON: invalid club id, missing file, unsupported type, malformed multipart body
- `401 Unauthorized` JSON: missing authenticated user
- `403 Forbidden` JSON: officer is not allowed to manage this club
- `404 Not Found` JSON: club not found
- `413 Request Entity Too Large` JSON: upload exceeds 10 MB
- `500 Internal Server Error` JSON: file storage failure

### `DELETE /api/officer/orgs/{id}/image`

Deletes the current image for an officer-managed club and clears the club `imagePath` field.

Authentication:

- Required.
- Officer role required.
- Officers may only delete images for clubs they lead.

Path params:

- `id: number` club id

Success response example:

```json
{
  "message": "Club image deleted successfully",
  "club": {
    "id": 12,
    "clubName": "CS Club",
    "imagePath": ""
  }
}
```

Possible errors:

- `400 Bad Request` JSON: invalid club id
- `401 Unauthorized` JSON: missing authenticated user
- `403 Forbidden` JSON: officer is not allowed to manage this club
- `404 Not Found` JSON: club not found or image not set
- `500 Internal Server Error` JSON: file storage failure

### `POST /api/admin/orgs/{id}/image`

Uploads or replaces the image for any club.

Authentication:

- Required.
- Admin role required.

Behavior and response shape are the same as `POST /api/officer/orgs/{id}/image`.

### `DELETE /api/admin/orgs/{id}/image`

Deletes the current image for any club.

Authentication:

- Required.
- Admin role required.

Behavior and response shape are the same as `DELETE /api/officer/orgs/{id}/image`.

### `POST /logout`

Logs out the current session.

Authentication:

- Optional.

Success response:

```json
{
  "message": "Logged out successfully"
}
```

Behavior:

- Removes the in-memory session if present.
- Expires the `session_token` cookie.

### `GET /api/prefill`

Returns session-backed profile fields used to prefill client forms.

Authentication:

- Required.

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

### `GET /api/adjectives`

Returns swipe question content for `adjective` and `personality_traits` question types.

Authentication:

- Required.

Success response example:

```json
[
  {
    "id": 4,
    "question_type": "adjective",
    "translations": {
      "en": ["Creative", "A trait related to imagination and originality"]
    }
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

Rules:

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
- `400 Bad Request` JSON: invalid JSON
- `400 Bad Request` JSON: `questionId must be a positive integer`

Frontend note:

- The current handler validates and echoes the payload. It does not persist it.

### `POST /submit`

Submits the demographics form.

Authentication:

- Required by default.
- Controlled by `SUBMIT_REQUIRES_AUTH`.

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

- `401 Unauthorized` JSON when auth is enabled and the session is invalid
- `400 Bad Request` JSON: invalid JSON
- `400 Bad Request` JSON: missing required demographics fields

Frontend note:

- The current handler validates and echoes the payload. It does not persist it.

### `GET /api/results`

Returns ranked organization matches for the authenticated user.

Authentication:

- Required.

Behavior:

- Reads the authenticated user’s stored answers.
- Fetches all clubs.
- Uses the matching algorithm to sort clubs by normalized match score.
- Returns the ranked organizations in order.

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
    "officers": ["president@example.edu", "vp@example.edu"],
    "personality": ["Welcoming", "Curious"],
    "activities": ["Community Service", "Guest Speakers"],
    "genders": ["Any"],
    "ethnicities": ["Any"],
    "religions": ["Any"],
    "strict_genders": false,
    "dedicated_majors": ["Computer Science"],
    "other": ["Open to first-year students"],
    "updatedAt": "2026-04-09T18:30:00Z"
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

### `POST /api/delete`

Deletes the authenticated user's account data.

Authentication:

- Required.

Headers:

- `Content-Type: application/json` (optional; no request body required)

Request body:

- None

Success response:

```json
{
  "message": "User data deleted successfully"
}
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session or missing authenticated user email
- `500 Internal Server Error` JSON:

```json
{
  "error": "<database error message>"
}
```

Frontend note:

- This endpoint deletes the authenticated user's stored profile and responses.
- The endpoint uses session identity (`X-User-Email` is set by auth middleware), so the client should not send a target email in the body.

### `GET /api/officer/orgs`

Returns the organizations managed by the authenticated officer.

Authentication:

- Required.
- User must have role `officer`.

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
    "officers": ["president@example.edu", "vp@example.edu"],
    "personality": ["Welcoming", "Curious"],
    "activities": ["Community Service", "Guest Speakers"],
    "genders": ["Any"],
    "ethnicities": ["Any"],
    "religions": ["Any"],
    "strict_genders": false,
    "dedicated_majors": ["Computer Science"],
    "other": ["Open to first-year students"],
    "updatedAt": "2026-04-09T18:30:00Z"
  }
]
```

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON: forbidden
- `500 Internal Server Error` JSON: failed to fetch officer orgs

### `PATCH /api/officer/orgs`

Partially updates an organization managed by the authenticated officer.

Authentication:

- Required.
- User must have role `officer`.

Headers:

- `Content-Type: application/json`

Request body:

```json
{
  "id": 12,
  "description": "Updated description",
  "meetingTime": "Fridays at 5:00 PM"
}
```

Success response:

- `200 OK`
- Response body is the full updated Organization object

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON: forbidden
- `400 Bad Request` JSON: invalid club id
- `400 Bad Request` JSON: invalid JSON
- `400 Bad Request` JSON: no fields provided for update

Validation:

- `id` must be a positive integer

### `GET /api/admin/orgs`

Returns all organizations for admin users.

Authentication:

- Required.
- User must have role `admin`.

Success response:

- JSON array of Organization objects

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON: forbidden
- `500 Internal Server Error` JSON: failed to fetch orgs

### `POST /api/admin/orgs`

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
  "officers": ["president@example.edu"],
  "personality": ["Welcoming", "Collaborative"],
  "activities": ["Study Groups", "Guest Speakers"],
  "genders": ["Any"],
  "ethnicities": ["Any"],
  "religions": ["Any"],
  "strict_genders": false,
  "dedicated_majors": ["Computer Science"],
  "other": ["Bring your laptop"]
}
```

Rules:

- `clubName` is required and must be non-empty.
- If `officers` is provided, every email must belong to an existing user record.

Success response:

- `201 Created`
- Response body is the created Organization object

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON: forbidden
- `400 Bad Request` JSON: invalid JSON
- `400 Bad Request` JSON: `clubName is required`

### `PATCH /api/admin/orgs`

Partially updates an organization as an admin.

Authentication:

- Required.
- User must have role `admin`.

Headers:

- `Content-Type: application/json`

Request body:

- Same shape as `PATCH /api/officer/orgs`

Success response:

- `200 OK`
- Response body is the full updated Organization object

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON: forbidden
- `400 Bad Request` JSON: invalid club id
- `400 Bad Request` JSON: invalid JSON
- `400 Bad Request` JSON: no fields provided for update

### `DELETE /api/admin/orgs`

Deletes an organization.

Authentication:

- Required.
- User must have role `admin`.

Request body:

```json
{
  "id": 12
}
```

Success response:

- `204 No Content`
- No response body

Possible errors:

- `401 Unauthorized` JSON: missing/invalid/expired session
- `403 Forbidden` JSON: forbidden
- `400 Bad Request` JSON: invalid club id
- `400 Bad Request` JSON: delete failed

## Frontend Checklist

- Use `credentials: 'include'` if relying on the `session_token` cookie.
- Prefer also storing the token returned by `/auth/callback` so requests can send `Authorization: Bearer <token>`.
- Use resource-based admin/officer update routes:
  - `PATCH /api/officer/orgs`
  - `PATCH /api/admin/orgs`
- Use `POST /api/admin/orgs` to create clubs.
- Use `DELETE /api/admin/orgs` to delete clubs.
- Use `POST /api/delete` from the delete-account page to delete the authenticated user's data.
- Expect both JSON and plain-text `401` responses depending on the endpoint.

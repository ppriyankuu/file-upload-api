# File Upload API
A small backend service for handling file uploads. 

It accepts multipart file uploads, validates the file type and size, saves files locally, and returns a public URL for accessing them.

### Features
- `POST /upload` endpoint to upload a single file
- Validates:
    - File size (server-level limit)
    - MIME type (e.g., PNG, JPEG, PDF)
- Stores files in a local `uploads/` directory
- Generates safe, unique filenames using UUIDs
- Serves uploaded files at `/uploads/<filename>`
- Clean project structure with separate handler, storage, config, and middleware layers

### Folder Structure
```
file-upload/
├── main.go
├── internal/
│   ├── handlers/
│   │   └── upload.go
│   ├── middleware/
│   │   └── limit.go
│   ├── storage/
│   │   └── local.go
│   └── config/
│       └── config.go
├── uploads/              # stored files (created at runtime)
├── go.mod
└── README.md
```

### How it works
1. Client sends a multipart form request with the field name file.
2. Middleware limits the maximum request size.
3. Handler reads the file header, checks MIME type, and reopens the file stream.
4. File is stored under `./uploads/` using a generated UUID + original extension.
5. Server returns a JSON response:
```
{
    "url": "/uploads/<generated-filename>"
}
```
6. Gin serves the `uploads/` directory as static content, making the file accessible publicly.

### Running the Server
1. Clone the repo
```
git clone https://github.com/ppriyankuu/file-upload-api
cd file-upload-api
```
2. Install dependencies
```
go mod tidy
```
3. Create the `.env` file
4. Run the server
```
go run main.go
```

### Upload example
Using `curl`:
```
curl -X POST -F "file=@example.png" http://localhost:8079/upload
```
Response:
```
{
    "url": "/uploads/2dfe1f4c-1b5c-4b53-b632.png"
}
```

#### Notes
- To add more allowed MIME types, update `ALLOWED_TYPES` in `.env`
- For production, replace local storage with S3 or other object storage
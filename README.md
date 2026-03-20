# Imagez - Image Processing API

A lightweight and powerful image processing API built with Go, Echo framework, and bimg library. Perform various image operations including conversion, resizing, compression, cropping, and watermarking.

## Features

- 🖼️ **Image Conversion** - Convert images between different formats
- 📏 **Image Resizing** - Resize images to custom dimensions
- 🗜️ **Image Compression** - Compress images to reduce file size
- ✂️ **Image Cropping** - Crop images with custom width and height
- 💧 **Watermarking** - Add text watermarks to images
- ⬇️ **File Download** - Download processed images

## Tech Stack

- **Language**: Go 1.25.4
- **Framework**: Echo v5 (HTTP framework)
- **Image Processing**: bimg v1.1.9

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Linux/Unix system (bimg requires libvips)

### Installation

```bash
# Clone the repository
git clone https://github.com/Zyprush18/imagez.git

# Navigate to project directory
cd imagez

# Install dependencies
go mod download

# Build and run
go run ./cmd/api/server.go
```

The API server will start on `http://api.localhost:8000/v1`

## API Documentation

### Base URL

```
http://api.localhost:8000/v1
```

### Response Format

All endpoints return JSON responses with the following structure:

#### Success Response
```json
{
  "message": "Operation successful",
  "data": {
    "file_name": "processed_image.zip"
  }
}
```

#### Error Response
```json
{
  "message": "Operation failed",
  "error": "Error description"
}
```

---

## Endpoints

### 1. Convert Image Format

**Endpoint:** `POST /v1/convert`

**Description:** Convert image(s) from one format to another.

**Request:**
- **Content-Type:** `multipart/form-data`
- **Required Fields:**
  - `images` (file, required) - Image file to convert
  - `format` (string, required) - Target format (e.g., jpg, png, webp, gif)

**Example Request:**
```bash
curl -X POST http://api.localhost:8000/v1/convert \
  -F "images=@image.jpg" \
  -F "format=png"
```

**Success Response (200):**
```json
{
  "message": "Image converted successfully",
  "data": {
    "file_name": "processed_images.zip"
  }
}
```

**Error Responses:**
- **400 Bad Request:**
  - Missing `images` or `format` field
  - Unsupported file type
  - Unsupported target format
- **500 Internal Server Error:** Processing failed

---

### 2. Resize Image

**Endpoint:** `POST /v1/resize`

**Description:** Resize image(s) to specified dimensions.

**Request:**
- **Content-Type:** `multipart/form-data`
- **Required Fields:**
  - `images` (file, required) - Image file to resize
  - `width` (integer, required) - Target width in pixels
  - `height` (integer, required) - Target height in pixels

**Example Request:**
```bash
curl -X POST http://api.localhost:8000/v1/resize \
  -F "images=@image.jpg" \
  -F "width=800" \
  -F "height=600"
```

**Success Response (200):**
```json
{
  "message": "Image resized successfully",
  "data": {
    "file_name": "processed_images.zip"
  }
}
```

**Error Responses:**
- **400 Bad Request:**
  - Missing `images`, `width`, or `height` field
  - Invalid width or height value (non-integer)
  - Unsupported image format
- **500 Internal Server Error:** Processing failed

---

### 3. Compress Image

**Endpoint:** `POST /v1/compress`

**Description:** Compress image(s) to reduce file size.

**Request:**
- **Content-Type:** `multipart/form-data`
- **Required Fields:**
  - `images` (file, required) - Image file to compress
  - `size` (integer, required) - Compression size/quality

**Example Request:**
```bash
curl -X POST http://api.localhost:8000/v1/compress \
  -F "images=@image.jpg" \
  -F "size=80"
```

**Success Response (200):**
```json
{
  "message": "Image compressed successfully",
  "data": {
    "file_name": "processed_images.zip"
  }
}
```

**Error Responses:**
- **400 Bad Request:**
  - Missing `images` or `size` field
  - Invalid size value (non-integer)
  - Unsupported image format
- **500 Internal Server Error:** Processing failed

---

### 4. Crop Image

**Endpoint:** `POST /v1/crop`

**Description:** Crop image(s) to specified dimensions.

**Request:**
- **Content-Type:** `multipart/form-data`
- **Required Fields:**
  - `images` (file, required) - Image file to crop
  - `width` (integer, required) - Crop width in pixels
  - `height` (integer, required) - Crop height in pixels

**Example Request:**
```bash
curl -X POST http://api.localhost:8000/v1/crop \
  -F "images=@image.jpg" \
  -F "width=500" \
  -F "height=400"
```

**Success Response (200):**
```json
{
  "message": "Image cropped successfully",
  "data": {
    "file_name": "processed_images.zip"
  }
}
```

**Error Responses:**
- **400 Bad Request:**
  - Missing `images`, `width`, or `height` field
  - Invalid width or height value
  - Invalid crop parameters
  - Unsupported image format
- **500 Internal Server Error:** Processing failed

---

### 5. Add Watermark

**Endpoint:** `POST /v1/watermark`

**Description:** Add text watermark to image(s).

**Request:**
- **Content-Type:** `multipart/form-data`
- **Required Fields:**
  - `images` (file, required) - Image file to watermark
  - `text` (string, required) - Watermark text

**Example Request:**
```bash
curl -X POST http://api.localhost:8000/v1/watermark \
  -F "images=@image.jpg" \
  -F "text=Copyright 2024"
```

**Success Response (200):**
```json
{
  "message": "Image watermarked successfully",
  "data": {
    "file_name": "processed_images.zip"
  }
}
```

**Error Responses:**
- **400 Bad Request:**
  - Missing `images` or `text` field
  - Unsupported image format
- **500 Internal Server Error:** Processing failed

---

### 6. Download Processed File

**Endpoint:** `GET /v1/downloads`

**Description:** Download a processed image file by filename.

**Query Parameters:**
- `filename` (string, required) - The filename to download (returned from processing endpoints)

**Example Request:**
```bash
curl -X GET "http://api.localhost:8000/v1/downloads?filename=processed_images.zip" \
  --output downloaded_file.zip
```

**Success Response (200):**
- File downloaded as binary attachment
- Content-Type: `application/octet-stream`

**Error Responses:**
- **400 Bad Request:** Missing `filename` query parameter
- **404 Not Found:** File not found
- **500 Internal Server Error:** Download failed

**Note:** The downloaded file is automatically deleted after successful download.

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Missing or invalid parameters |
| 404 | Not Found - File not found |
| 500 | Internal Server Error - Processing failed |

## Supported Image Formats

The API supports common image formats through the bimg library:
- JPEG
- PNG
- WebP
- GIF
- TIFF
- BMP

## Project Structure

```
imagez/
├── cmd/
│   └── api/
│       └── server.go           # Main server entry point
├── internal/
│   ├── handler/
│   │   └── images.go           # Request handlers for image operations
│   ├── service/
│   │   └── images_service.go   # Business logic for image processing
│   └── routes/
│       └── api.go              # API route definitions
├── pkg/
│   └── bimg.go                 # Image processing wrapper
├── utils/
│   ├── error.go                # Error definitions
│   ├── image.go                # Image utility functions
│   └── response.go             # Response formatting utilities
├── go.mod                      # Go module definition
└── README.md                   # This file
```

### Building the Project

```bash
go build -o imagez ./cmd/api/server.go
./imagez
```

## Configuration

The API server runs on:
- **Host:** `api.localhost:8000`
- **Port:** `8000`
- **API Version:** `v1`
- **Middleware:** Request logging and panic recovery

## Middleware

The API includes the following middleware:
- **Request Logger** - Logs all incoming requests
- **Recover** - Recovers from panics and returns 500 errors

## Batch Processing

The API supports batch image processing. Multiple images can be processed in a single request with multipart form data. Processed files are zipped and available for download.

## File Storage

- Processed images are temporarily stored on the server
- Files are automatically deleted after download
- Downloaded output is provided as a ZIP archive for multiple images

## Performance

The API uses Go's concurrency features with worker pools (CPU count) for efficient image processing when handling multiple files.


## Author

Zyprush18

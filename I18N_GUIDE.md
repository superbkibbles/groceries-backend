# Internationalization (i18n) Guide

This guide explains how to use the internationalization features in the groceries backend API.

## Overview

The backend now supports multiple languages with English (en) and Arabic (ar) as the primary supported languages. The system automatically detects the user's preferred language and returns localized messages.

## Features

- **Automatic Language Detection**: Detects language from multiple sources with priority order
- **Localized Error Messages**: All API responses include localized error and success messages
- **Language Switching**: Users can switch their language preference via API
- **Fallback Support**: Falls back to English if the requested language is not supported

## Language Detection Priority

The system detects the user's preferred language in the following order:

1. **Query Parameter**: `?lang=ar` (highest priority)
2. **Accept-Language Header**: `Accept-Language: ar,en;q=0.9`
3. **Cookie**: `language=ar`
4. **Default Language**: English (fallback)

## API Endpoints

### Get Supported Languages

```http
GET /api/v1/languages
```

**Response:**
```json
{
  "message": "Operation completed successfully",
  "data": {
    "languages": ["en", "ar"],
    "default_language": "en"
  }
}
```

### Switch Language

```http
POST /api/v1/languages/switch
Content-Type: application/json

{
  "language": "ar"
}
```

**Response:**
```json
{
  "message": "Operation completed successfully",
  "data": {
    "language": "ar"
  }
}
```

## Usage Examples

### 1. Using Query Parameter

```bash
# English (default)
curl "http://localhost:8080/api/v1/products"

# Arabic
curl "http://localhost:8080/api/v1/products?lang=ar"
```

### 2. Using Accept-Language Header

```bash
# Arabic
curl -H "Accept-Language: ar" "http://localhost:8080/api/v1/products"

# English
curl -H "Accept-Language: en" "http://localhost:8080/api/v1/products"
```

### 3. Using Cookie (after language switch)

```bash
# First switch language
curl -X POST -H "Content-Type: application/json" \
  -d '{"language":"ar"}' \
  "http://localhost:8080/api/v1/languages/switch"

# Then use the API (cookie will be automatically sent)
curl -b cookies.txt "http://localhost:8080/api/v1/products"
```

## Response Format

All API responses now include localized messages:

### Success Response
```json
{
  "message": "Product created successfully", // English
  "data": {
    // ... product data
  }
}
```

### Arabic Response
```json
{
  "message": "تم إنشاء المنتج بنجاح", // Arabic
  "data": {
    // ... product data
  }
}
```

### Error Response
```json
{
  "error": "Product not found",
  "message": "Product not found",
  "code": 404
}
```

## Available Translation Keys

The system includes translations for all common API messages:

### Product Messages
- `product_created`, `product_updated`, `product_deleted`
- `product_not_found`, `products_retrieved`
- `product_name_required`, `product_price_required`
- And many more...

### User Messages
- `user_created`, `user_updated`, `user_deleted`
- `user_not_found`, `user_login_success`
- `user_unauthorized`, `user_forbidden`
- And many more...

### Order Messages
- `order_created`, `order_updated`, `order_cancelled`
- `order_not_found`, `orders_retrieved`
- And many more...

### General Messages
- `success`, `error`, `validation_error`
- `internal_server_error`, `bad_request`
- `unauthorized`, `forbidden`, `not_found`
- And many more...

## Configuration

### Environment Variables

You can configure i18n settings using environment variables:

```bash
# Default language (default: en)
DEFAULT_LANGUAGE=en

# Supported languages (default: en,ar)
SUPPORTED_LANGUAGES=en,ar

# Translation files path (default: internal/locales)
I18N_BUNDLE_PATH=internal/locales
```

### Adding New Languages

1. Create a new translation file in `internal/locales/`:
   ```
   internal/locales/fr.json  # French
   internal/locales/es.json  # Spanish
   ```

2. Add the language to supported languages in configuration:
   ```bash
   SUPPORTED_LANGUAGES=en,ar,fr,es
   ```

3. Update the i18n configuration if needed.

## Frontend Integration

### JavaScript/TypeScript

```javascript
// Set language via query parameter
const apiCall = (lang = 'en') => {
  return fetch(`/api/v1/products?lang=${lang}`)
    .then(response => response.json());
};

// Set language via header
const apiCallWithHeader = (lang = 'en') => {
  return fetch('/api/v1/products', {
    headers: {
      'Accept-Language': lang
    }
  }).then(response => response.json());
};

// Switch language permanently
const switchLanguage = async (lang) => {
  await fetch('/api/v1/languages/switch', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ language: lang })
  });
};
```

### React Example

```jsx
import { useState, useEffect } from 'react';

const useApi = () => {
  const [language, setLanguage] = useState('en');

  const switchLanguage = async (newLang) => {
    try {
      await fetch('/api/v1/languages/switch', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ language: newLang })
      });
      setLanguage(newLang);
    } catch (error) {
      console.error('Failed to switch language:', error);
    }
  };

  const apiCall = async (endpoint) => {
    const response = await fetch(`/api/v1${endpoint}?lang=${language}`);
    return response.json();
  };

  return { apiCall, switchLanguage, language };
};
```

## Development

### Adding New Translation Keys

1. Add the key to both `internal/locales/en.json` and `internal/locales/ar.json`
2. Use the key in your handlers:

```go
// In your handler
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // ... business logic
    
    Created(c, "product_created", product)
}
```

### Using Translations in Handlers

```go
import "github.com/superbkibbles/ecommerce/internal/utils"

// Simple translation
message := utils.TSimple(c, "product_created")

// Translation with template data
templateData := map[string]interface{}{
    "productName": "iPhone 15",
    "price": 999.99,
}
message := utils.TWithData(c, "product_created_with_details", templateData)
```

## Testing

### Manual Testing

```bash
# Test English
curl "http://localhost:8080/api/v1/products?lang=en"

# Test Arabic
curl "http://localhost:8080/api/v1/products?lang=ar"

# Test language switching
curl -X POST -H "Content-Type: application/json" \
  -d '{"language":"ar"}' \
  "http://localhost:8080/api/v1/languages/switch"
```

### Automated Testing

```go
func TestProductHandlerWithArabic(t *testing.T) {
    router := setupTestRouter()
    
    req := httptest.NewRequest("GET", "/api/v1/products?lang=ar", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    // Check that message is in Arabic
    assert.Contains(t, response["message"], "تم")
}
```

## Best Practices

1. **Always use translation keys** in your handlers instead of hardcoded strings
2. **Test with both languages** to ensure all messages are properly translated
3. **Use descriptive translation keys** that clearly indicate their purpose
4. **Keep translations consistent** across all API endpoints
5. **Handle missing translations gracefully** by falling back to the translation key

## Troubleshooting

### Common Issues

1. **Translation not found**: Check that the key exists in both language files
2. **Wrong language detected**: Verify the Accept-Language header format
3. **Cookie not working**: Ensure the cookie domain and path are correct
4. **Fallback to English**: Check that the requested language is in the supported languages list

### Debug Mode

To debug language detection, check the request headers and cookies in your browser's developer tools or use curl with verbose output:

```bash
curl -v -H "Accept-Language: ar" "http://localhost:8080/api/v1/products"
```

## Support

For issues or questions regarding internationalization, please check:
1. This guide
2. The translation files in `internal/locales/`
3. The middleware implementation in `internal/adapters/http/middleware/i18n.go`
4. The utility functions in `internal/utils/i18n.go`

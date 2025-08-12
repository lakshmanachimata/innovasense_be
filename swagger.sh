#!/bin/bash

# Swagger documentation generation and serving script

echo "🔄 Generating Swagger documentation..."

# Generate Swagger docs
swag init

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo ""
    echo "📚 Available documentation:"
    echo "   - Swagger UI: http://localhost:8500/swagger/index.html"
    echo "   - JSON: http://localhost:8500/swagger/doc.json"
    echo "   - YAML: http://localhost:8500/swagger/doc.yaml"
    echo ""
    echo "🚀 To start the server with Swagger:"
    echo "   go run main.go"
    echo ""
    echo "📁 Generated files:"
    echo "   - docs/docs.go"
    echo "   - docs/swagger.json"
    echo "   - docs/swagger.yaml"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi

#!/bin/bash

# Comprehensive API Test Script for InnovoSens Backend
# Database: MySQL (innosense)
# Date: $(date)

BASE_URL="http://localhost:8500"
RESULTS_FILE="final_api_test_results_$(date +%Y%m%d_%H%M%S).txt"

echo "=== INNOVOSENS COMPREHENSIVE API TEST REPORT ===" > $RESULTS_FILE
echo "Test Date: $(date)" >> $RESULTS_FILE
echo "Database: MySQL (innosense)" >> $RESULTS_FILE
echo "Base URL: $BASE_URL" >> $RESULTS_FILE
echo "" >> $RESULTS_FILE

# Function to test API endpoint
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    local test_name=$5
    
    echo "Testing: $test_name" | tee -a $RESULTS_FILE
    echo "Endpoint: $method $endpoint" | tee -a $RESULTS_FILE
    echo "Request Body: $data" | tee -a $RESULTS_FILE
    echo "Headers: $headers" | tee -a $RESULTS_FILE
    
    if [ -n "$headers" ]; then
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}\nTIME:%{time_total}" -X $method "$BASE_URL$endpoint" -H "Content-Type: application/json" -H "$headers" -d "$data")
    else
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}\nTIME:%{time_total}" -X $method "$BASE_URL$endpoint" -H "Content-Type: application/json" -d "$data")
    fi
    
    echo "Response: $response" | tee -a $RESULTS_FILE
    echo "---" | tee -a $RESULTS_FILE
    echo "" | tee -a $RESULTS_FILE
}

# Wait for server to start
echo "Waiting for server to start..."
sleep 3

# Test 1: Health Check
test_api "GET" "/health" "" "" "Health Check API"

# Test 2: Root Endpoint
test_api "GET" "/" "" "" "Root Endpoint"

# Test 3: User Registration (New User) - Encrypted Data
test_api "POST" "/Services/innovoregister" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "userpin": "8UkrqMI8cqtn+jndvnjAoA==",
    "username": "Final Test User 1",
    "gender": "Male",
    "age": 25,
    "height": 170.5,
    "weight": 70.0
}' "" "User Registration - New User (Encrypted)"

# Test 4: User Registration (With Contact Number) - Encrypted Data
test_api "POST" "/Services/innovoregister" '{
    "email": "5HfZnB40gzld/98oJ3U4p1HQlwLfY22jxhCzBKoCGsk=",
    "cnumber": "tkB3SU1u0Q2CaV1wpIJzhA==",
    "userpin": "aF2GYlwtO9el+uEUho7sOg==",
    "username": "Final Test User 2",
    "gender": "Female",
    "age": 28,
    "height": 165.0,
    "weight": 60.0
}' "" "User Registration - With Contact Number (Encrypted)"

# Test 5: User Registration (Duplicate Email) - Encrypted Data
test_api "POST" "/Services/innovoregister" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "userpin": "sh3MNi3dYHfOV6jyEMthog==",
    "username": "Final Test User 3",
    "gender": "Male",
    "age": 30,
    "height": 175.0,
    "weight": 75.0
}' "" "User Registration - Duplicate Email (Should Fail) (Encrypted)"

# Test 6: User Login (Valid Credentials) - Encrypted Data
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8500/Services/innovologin" -H "Content-Type: application/json" -d '{"email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=", "userpin": "8UkrqMI8cqtn+jndvnjAoA=="}')
JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.jwt_token')

test_api "POST" "/Services/innovologin" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "userpin": "8UkrqMI8cqtn+jndvnjAoA=="
}' "" "User Login - Valid Credentials (Encrypted)"

# Test 7: User Login (Invalid Credentials) - Encrypted Data
test_api "POST" "/Services/innovologin" '{
    "email": "GIQSlzqh445148nIVbKyXnY6u5A0cGWnJUhPC5y7+sk=",
    "userpin": "wrong123"
}' "" "User Login - Invalid Credentials (Should Fail) (Encrypted)"

# Test 8: Banner Images
test_api "POST" "/Services/getBannerImages" '{}' "" "Get Banner Images"

# Test 9: Home Images
test_api "POST" "/Services/getHomeImages" '{}' "" "Get Home Images"

# Test 10: Devices
test_api "POST" "/Services/getDevices" '{}' "" "Get Devices"

# Test 11: Hydration APIs with JWT Token
echo "Testing Hydration APIs with JWT Token..." | tee -a $RESULTS_FILE
echo "JWT Token: $JWT_TOKEN" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Test 12: Basic Hydration Data Submission
test_api "POST" "/Services/protected/innovoHyderation" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "weight": 70.0,
    "height": 170.5,
    "sweat_position": 0.6,
    "time_taken": 45.0,
    "device_type": 1,
    "image_path": "/test/image1.jpg",
    "image_id": 1
}' "Authorization: Bearer $JWT_TOKEN" "Basic Hydration Data Submission"

# Test 13: Enhanced Hydration Data Submission
test_api "POST" "/Services/protected/newinnovoHyderation" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "weight": 70.0,
    "height": 170.5,
    "sweat_position": 0.7,
    "time_taken": 30.0,
    "device_type": 2,
    "image_path": "/test/image2.jpg",
    "image_id": 2
}' "Authorization: Bearer $JWT_TOKEN" "Enhanced Hydration Data Submission"

# Test 14: Update Hydration Value
test_api "POST" "/Services/protected/updateHyderationValue" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "id": 1,
    "weight": 71.0,
    "height": 170.5,
    "sweat_position": 0.8,
    "time_taken": 50.0,
    "bmi": 24.4,
    "tbsa": 1.8,
    "sweat_rate": 35.0,
    "sweat_loss": 25.0,
    "device_type": 1
}' "Authorization: Bearer $JWT_TOKEN" "Update Hydration Value"

# Test 15: Update Sweat Data
test_api "POST" "/Services/protected/updateSweatData" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "image_id": 1,
    "sweat_rate": 40.0,
    "sweat_loss": 30.0
}' "Authorization: Bearer $JWT_TOKEN" "Update Sweat Data"

# Test 16: Get Summary
test_api "POST" "/Services/protected/getSummary" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "sweat_position": 0.6
}' "Authorization: Bearer $JWT_TOKEN" "Get Summary"

# Test 17: Get User Detailed Summary
test_api "POST" "/Services/protected/getUserDetailedSummary" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "id": 1
}' "Authorization: Bearer $JWT_TOKEN" "Get User Detailed Summary"

# Test 18: Get Hydration Summary Screen
test_api "POST" "/Services/protected/getHydrationSummaryScreen" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "id": 1
}' "Authorization: Bearer $JWT_TOKEN" "Get Hydration Summary Screen"

# Test 19: Get Client History
test_api "POST" "/Services/protected/getClientHistory" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "from_date": "2024-01-01",
    "to_date": "2024-12-31"
}' "Authorization: Bearer $JWT_TOKEN" "Get Client History"

# Test 20: Get Hydration History
test_api "POST" "/Services/protected/getHyderartionHistory" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "from_date": "2024-01-01",
    "to_date": "2024-12-31"
}' "Authorization: Bearer $JWT_TOKEN" "Get Hydration History"

# Test 21: Get Electrolyte History
test_api "POST" "/Services/protected/getElectrolyteHistory" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9,
    "from_date": "2024-01-01",
    "to_date": "2024-12-31"
}' "Authorization: Bearer $JWT_TOKEN" "Get Electrolyte History"

# Test 22: Get Sweat Images
test_api "POST" "/Services/protected/getSweatImages" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "userid": 9
}' "Authorization: Bearer $JWT_TOKEN" "Get Sweat Images"

# Test 23: Protected Route (Without JWT Token)
test_api "POST" "/Services/protected/getSummary" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "sweat_position": 0.6
}' "" "Protected Route - Without JWT Token (Should Fail)"

# Test 24: Protected Route (With Invalid JWT Token)
test_api "POST" "/Services/protected/getSummary" '{
    "email": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "username": "Final Test User 1",
    "sweat_position": 0.6
}' "Authorization: Bearer invalid.jwt.token" "Protected Route - With Invalid JWT Token (Should Fail)"

# Test 25: Organization APIs (These should fail due to PostgreSQL placeholder issue)
test_api "POST" "/Services/getHydrationRecommendation" '{
    "name": "Test User",
    "contact": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "gender": "Male",
    "age": 25,
    "sweat_position": 0.5,
    "workout_time": 30.0,
    "height": 170.5,
    "weight": 70.0
}' "apikey: innosense-api-key-2024" "Hydration Recommendation - Organization API (Expected to Fail)"

# Test 26: Historical Data - Organization API
test_api "POST" "/Services/getHistoricalData" '{
    "contact": "sgFU5M2NNmvJpwwY7aTqrrv+HUfcaqfkRSjffKP/pE8=",
    "start_date": "2024-01-01",
    "end_date": "2024-12-31"
}' "apikey: innosense-api-key-2024" "Historical Data - Organization API (Expected to Fail)"

# Test 27: Swagger Documentation
test_api "GET" "/swagger/index.html" "" "" "Swagger Documentation"

echo "=== COMPREHENSIVE API TEST COMPLETED ===" | tee -a $RESULTS_FILE
echo "Results saved to: $RESULTS_FILE" | tee -a $RESULTS_FILE

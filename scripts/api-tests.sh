#!/bin/bash
set -euo pipefail

API_URL="http://localhost:3000"

function check_status() {
  local expected=$1
  local actual=$2
  local msg=$3
  if [[ "$actual" == "$expected" ]]; then
    echo "[PASS] $msg (HTTP $actual)"
  else
    echo "[FAIL] $msg - Expected HTTP $expected but got $actual"
    exit 1
  fi
}

echo "Test 1: Add a computer and get it by MAC"

resp_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/computers" \
  -H "Content-Type: application/json" \
  -d '{
    "mac_address": "00:11:22:33:44:66",
    "computer_name": "testcomp",
    "ip_address": "192.168.1.99",
    "employee_abbreviation": "abc",
    "description": "test machine"
  }')

check_status 201 "$resp_code" "POST /computers create"

get_resp=$(curl -s -w "\n%{http_code}" "$API_URL/computers/00:11:22:33:44:66")
body=$(echo "$get_resp" | sed '$d')
code=$(echo "$get_resp" | tail -n1)

check_status 200 "$code" "GET /computers/00:11:22:33:44:66 fetch"

name=$(echo "$body" | jq -r '.computer_name')
if [[ "$name" == "testcomp" ]]; then
  echo "[PASS] GET returned correct computer_name"
else
  echo "[FAIL] GET returned wrong computer_name: $name"
  exit 1
fi

echo -e "\nTest 2: Add multiple computers and get all"

for i in alpha beta; do
  mac="00:11:22:33:44:7$( [[ $i == alpha ]] && echo 7 || echo 8 )"
  name="$i"
  ip="192.168.0.$([[ $i == alpha ]] && echo 1 || echo 2)"
  emp="$([[ $i == alpha ]] && echo abc || echo xyz)"
  
  resp_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/computers" \
    -H "Content-Type: application/json" \
    -d "{
      \"mac_address\": \"$mac\",
      \"computer_name\": \"$name\",
      \"ip_address\": \"$ip\",
      \"employee_abbreviation\": \"$emp\"
    }")

  check_status 201 "$resp_code" "POST /computers create $name"
done

all_resp=$(curl -s -w "\n%{http_code}" "$API_URL/computers")
all_code=$(echo "$all_resp" | tail -n 1)
all_body=$(echo "$all_resp" | sed '$d')

check_status 200 "$all_code" "GET /computers list all"

count=$(echo "$all_body" | jq 'length')
if [[ "$count" -ge 2 ]]; then
  echo "[PASS] GET /computers returned $count entries"
else
  echo "[FAIL] GET /computers returned less than 2 entries"
  exit 1
fi

echo -e "\nTest 3: Get computers by employee"

resp_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/computers" \
  -H "Content-Type: application/json" \
  -d '{
    "mac_address": "00:11:22:33:44:99",
    "computer_name": "empMachine",
    "ip_address": "192.168.1.20",
    "employee_abbreviation": "jdo",
    "description": "assigned to employee jdoe"
  }')

check_status 201 "$resp_code" "POST /computers create empMachine"

by_emp_resp=$(curl -s -w "\n%{http_code}" "$API_URL/employee/jdo/computers")
by_emp_code=$(echo "$by_emp_resp" | tail -n 1)
by_emp_body=$(echo "$by_emp_resp" | sed '$d')

check_status 200 "$by_emp_code" "GET /employee/jdo/computers"

emp_count=$(echo "$by_emp_body" | jq 'length')
if [[ "$emp_count" -ge 1 ]]; then
  echo "[PASS] GET by employee returned $emp_count entries"
else
  echo "[FAIL] GET by employee returned no entries"
  exit 1
fi

emp_name=$(echo "$by_emp_body" | jq -r '.[0].computer_name')
if [[ "$emp_name" == "empMachine" ]]; then
  echo "[PASS] GET by employee returned correct computer_name"
else
  echo "[FAIL] GET by employee returned wrong computer_name: $emp_name"
  exit 1
fi


echo -e "\nTest 4: Add 3 computers to the same employee (admin-notification should be called in the logs)"

emp="emp3"
for i in 1 2 3; do
  mac="00:11:22:33:55:0$i"
  name="emp3machine$i"
  ip="10.0.0.$i"

  resp_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/computers" \
    -H "Content-Type: application/json" \
    -d "{
      \"mac_address\": \"$mac\",
      \"computer_name\": \"$name\",
      \"ip_address\": \"$ip\",
      \"employee_abbreviation\": \"$emp\",
      \"description\": \"device $i for $emp\"
    }")

  check_status 201 "$resp_code" "POST create $name for $emp"
done

emp3_resp=$(curl -s -w "\n%{http_code}" "$API_URL/employee/$emp/computers")
emp3_code=$(echo "$emp3_resp" | tail -n1)
emp3_body=$(echo "$emp3_resp" | sed '$d')

check_status 200 "$emp3_code" "GET /employee/$emp/computers"

emp3_count=$(echo "$emp3_body" | jq 'length')
if [[ "$emp3_count" -eq 3 ]]; then
  echo "[PASS] $emp has 3 computers"
else
  echo "[FAIL] $emp has $emp3_count computers"
  exit 1
fi

echo -e "\n✅ All API tests passed!"

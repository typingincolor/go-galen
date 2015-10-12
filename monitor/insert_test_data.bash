#!/usr/bin/env bash
mongo monitors --eval 'db.apis.insert({"url":"http://localhost:8080/1", "method": "GET"});'
mongo monitors --eval 'db.apis.insert({"url":"http://localhost:8080/2", "method": "GET"});'
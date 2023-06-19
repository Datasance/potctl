#!/usr/bin/env bash

# These functions are designed to be used with the bats `run` command

function curlMsvc(){
    PUBLIC_ENDPOINT="$1"
    curl -s --max-time 120 ${PUBLIC_ENDPOINT}/api/raw
}

function jqMsvcArray(){
    ARR="$1"
    echo "$ARR" | jq '. | length'
}

function findMsvcState(){
    NS="$1"
    MS="$2"
    STATE="$3"
    potctl -n $NS get microservices | grep $MS | grep $STATE
}

function runNoExecutors(){
  echo '' > test/conf/nothing.yaml
  potctl deploy -f test/conf/nothing.yaml -n "$NS"
}

function runWrongNamespace(){
  echo "---
apiVersion: datasance.com/v1
kind: LocalControlPlane
metadata:
  namespace: wrong
spec:
  iofogUser:
    name: Testing
    surname: Functional
    email: user@domain.com
    password: S5gYVgLEZV
  controller:
    name: func-test" > test/conf/wrong-ns.yaml
  potctl deploy -f test/conf/wrong-ns.yaml -n "$NS"
}
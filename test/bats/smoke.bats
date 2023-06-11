#!/usr/bin/env bash

. test/func/include.bash

@test "Help" {
  potctl --help
}

@test "Help w/o Flag" {
  potctl help
}

@test "Create Help" {
  potctl create --help
}

@test "Delete Help" {
  potctl delete --help
}

@test "Deploy Help" {
  potctl deploy --help
}

@test "Describe Help" {
  potctl describe --help
}

@test "Connect Help" {
  potctl connect --help
}

@test "Disconnect Help" {
  potctl disconnect --help
}

@test "Legacy Help" {
  potctl legacy --help
}

@test "Logs Help" {
  potctl logs --help
}

@test "Get Help" {
  potctl get --help
}

@test "Version" {
  potctl version
}

@test "Get All" {
  potctl get all
}

@test "Get Namespaces" {
  potctl get namespaces
}

@test "Get Controllers" {
  potctl get controllers
}

@test "Get Agents" {
  potctl get agents
}

@test "Get Microservices" {
  potctl get microservices
}

@test "Get Applications" {
  potctl get applications
}

@test "Create Namespace" {
  potctl create namespace smoketestsnamespace1234
}

@test "Set Default Namespace" {
  potctl configure current-namespace smoketestsnamespace1234
  potctl get all
}

@test "Delete Namespace" {
  potctl delete namespace smoketestsnamespace1234
  potctl get all
  potctl get namespaces
}
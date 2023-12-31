#!/usr/bin/env bash

NAME="func-test"
APPLICATION_NAME="func-app"
MICROSERVICE_NAME="func-msvc"
MSVC1_NAME="func-app-server"
MSVC2_NAME="func-app-ui"
VOL_SRC="/tmp/potctl_tests"
WIN_VOL_SRC="C:/tests"
VOL_DEST="/tmp/iofog"
VOL_CONT_DEST="/data"
VOL_NAME="volume-data"
USER_PW="S5gYVgLEZV"
USER_PW_B64="UzVnWVZnTEVaVg=="
USER_EMAIL="user@domain.com"
ROUTE_NAME="route-1"
EDGE_RESOURCE_NAME="smart-door"
EDGE_RESOURCE_VERSION="v1.0.0"
EDGE_RESOURCE_DESC="Very smart door"
EDGE_RESOURCE_PROTOCOL="https"
APP_TEMPLATE_NAME="app-tpl-1"
APP_TEMPLATE_DESC="This is an application template to test with"
APP_TEMPLATE_KEY="agent-name"
APP_TEMPLATE_KEY_DESC="Name of Agent to deploy Microservices to"
APP_TEMPLATE_DEF_VAL="func-test-0"
VOL_INVALID_DEST="/img/iofog"
VOL_CONT_INVALID_DEST="/img/data"
MSVC3_NAME="func-app-server-1"
MSVC4_NAME="func-app-ui-1"
MSVC5_NAME="func-heart-rate-ui"
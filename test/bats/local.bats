#!/usr/bin/env bash

. test/func/include.bash

NS="$NAMESPACE"

@test "Initialize tests" {
  stopTest
}

@test "Create namespace" {
  startTest
  potctl create namespace "$NS"
  stopTest
}

@test "Test no executors" {
  startTest
  testNoExecutors
  stopTest
}

@test "Test wrong namespace metadata" {
  startTest
  testWrongNamespace
  stopTest
}

@test "Deploy local Controller" {
  startTest
  initLocalControllerFile
  potctl -v -n "$NS" deploy -f test/conf/local.yaml
  checkControllerLocal
  stopTest
}

@test "Controller legacy commands after deploy" {
  startTest
  potctl -v -n "$NS" legacy controller "$NAME" iofog list
  checkLegacyController
  stopTest
}

@test "Deploy Agents against local Controller" {
  startTest
  initLocalAgentFile
  potctl -v -n "$NS" deploy -f test/conf/local-agent.yaml
  checkAgent "${NAME}-0"
  stopTest
}

@test "Deploy Application for docker pull stats" {
  startTest
  initDockerPullStatsApplicationFiles
  potctl -v -n "$NS" deploy -f test/conf/application_pull_stat.yaml
  waitForPullingMsvc "$MSVC5_NAME" "$NS"
  checkPullPercentageOfMicroservice "$MSVC5_NAME" "$NS"
  stopTest
}

@test "Edge Resources" {
  startTest
  testEdgeResources
  stopTest
}

@test "Agent legacy commands" {
  startTest
  potctl -v -n "$NS" legacy agent "${NAME}-0" status
  checkLegacyAgent "${NAME}-0"
  stopTest
}

@test "Agent config dev mode" {
  startTest  
  [[ ! -z $(potctl -v -n "$NS" legacy agent "${NAME}-0" 'config -dev on') ]]
  stopTest
}

@test "Deploy local Controller again for indempotence" {
  startTest
  initLocalControllerFile
  potctl -v -n "$NS" deploy -f test/conf/local.yaml
  checkControllerLocal
  stopTest
}

@test "Deploy Agents against local Controller again for indempotence" {
  startTest
  initLocalAgentFile
  potctl -v -n "$NS" deploy -f test/conf/local-agent.yaml
  checkAgent "${NAME}-0"
  stopTest
}

@test "Deploy Volumes" {
  startTest
  testDeployLocalVolume
  testGetDescribeLocalVolume
  stopTest
}

@test "Deploy Volumes Idempotent" {
  startTest
  testDeployLocalVolume
  testGetDescribeLocalVolume
  stopTest
}

@test "Delete Volumes and Redeploy" {
  startTest
  testDeleteLocalVolume
  testDeployLocalVolume
  testGetDescribeLocalVolume
  stopTest
}

@test "Deploy Application Template and Templated Application" {
  startTest
  testApplicationTemplates
  stopTest
}

@test "Deploy Application" {
  startTest
  initApplicationFiles
  potctl -v -n "$NS" deploy -f test/conf/application.yaml
  checkApplication
  waitForMsvc "$MSVC1_NAME" "$NS"
  waitForMsvc "$MSVC2_NAME" "$NS"
  stopTest
}

@test "Deploy Microservice" {
  startTest
  initMicroserviceFile
  potctl -v -n "$NS" deploy -f test/conf/microservice.yaml
  checkMicroservice
  stopTest
}

@test "Deploy Route" {
  startTest
  initRouteFile
  potctl -v -n "$NS" deploy -f test/conf/route.yaml
  checkRoute "$ROUTE_NAME" "$MSVC1_NAME" "$MSVC2_NAME"
  stopTest
}

@test "Update Microservice" {
  startTest
  initMicroserviceUpdateFile
  potctl --debug -n "$NS" deploy -f test/conf/updatedMicroservice.yaml
  checkUpdatedMicroservice
  checkRoute "$ROUTE_NAME" "$MSVC1_NAME" "$MSVC2_NAME"
  stopTest
}

@test "Rename and Delete Route" {
  startTest
  local NEW_ROUTE_NAME="route-2"
  potctl -v -n "$NS" rename route $APPLICATION_NAME/"$ROUTE_NAME" "$NEW_ROUTE_NAME"
  potctl -v -n "$NS" delete route $APPLICATION_NAME/"$NEW_ROUTE_NAME"
  checkRouteNegative "$NEW_ROUTE_NAME" "$MSVC1_NAME" "$MSVC2_NAME"
  checkRouteNegative "$ROUTE_NAME" "$MSVC1_NAME" "$MSVC2_NAME"
  stopTest
}

@test "Delete Microservice using file option" {
  startTest
  potctl -v -n "$NS" delete -f test/conf/updatedMicroservice.yaml
  checkMicroserviceNegative
  stopTest
}

@test "Deploy Microservice in Application" {
  startTest
  initMicroserviceFile
  potctl -v -n "$NS" deploy -f test/conf/microservice.yaml
  checkMicroservice
  stopTest
}

@test "Get local json logs" {
  startTest
  [[ ! -z $(potctl -v -n "$NS" logs controller $NAME) ]]
  [[ ! -z $(potctl -v -n "$NS" logs agent ${NAME}-0) ]]
  [potctl logs ${NAME}-0 -n "$NS" | jq -e . >/dev/null 2>&1  | echo ${PIPESTATUS[1]}  -eq 0 ]
  stopTest
}

@test "Deploy Registry" {
  startTest
  initGCRRegistryFile
  potctl -v -n "$NS" deploy -f test/conf/gcr.yaml
  checkGCRRegistry
  initUpdatedGCRRegistryFile
  potctl -v -n "$NS" deploy -f test/conf/gcr.yaml
  checkUpdatedGCRRegistry
  potctl -v -n "$NS" delete registry 3
  checkGCRRegistryNegative
  stopTest
}

@test "Detach should fail because of running msvc" {
  startTest
  run potctl -v -n "$NS" detach agent ${NAME}-0
  [ "$status" -eq 1 ]
  echo "$output" | grep "because it still has microservices running. Remove the microservices first, or use the --force option."
  stopTest
}

@test "Detach/attach Agent" {
  startTest
  potctl -v -n "$NS" detach agent ${NAME}-0 --force
  potctl -v describe agent ${NAME}-0 --detached
  potctl -v -n "$NS" attach agent ${NAME}-0
  stopTest
}

@test "Deploy Application from file and test Application update" {
  startTest
  potctl -v -n "$NS" deploy -f test/conf/application.yaml
  checkApplication
  stopTest
}

@test "Deploy Application with volume missing " {
  startTest
  initInvalidApplicationFiles
  potctl -v -n "$NS" deploy -f test/conf/application_volume_missing.yaml
  waitForFailedMsvc "$MSVC4_NAME" "$NS"
  [[ ! -z $(potctl get microservices -n "$NS"  | grep "Volume missing") ]]
  potctl get microservices -n "$NS"
  stopTest
}

@test "Delete Agent should fail because of running msvc" {
  startTest
  run potctl -v -n "$NS" delete agent ${NAME}-0
  [ "$status" -eq 1 ]
  echo "$output" | grep "because it still has microservices running. Remove the microservices first, or use the --force option."
  stopTest
}

@test "Delete agent should work with --force option" {
  startTest
  potctl -v -n "$NS" delete agent ${NAME}-0 --force
  stopTest
}

@test "Delete all using file" {
  startTest
  initAllLocalDeleteFile
  potctl -v -n "$NS" delete -f test/conf/all-local.yaml
  checkApplicationNegative
  checkControllerNegative
  checkAgentNegative "${NAME}-0"
  stopTest
}

@test "Delete Namespace" {
  startTest
  potctl -v delete namespace "$NS"
  [[ -z $(potctl get namespaces | grep "$NS") ]]
  stopTest
}

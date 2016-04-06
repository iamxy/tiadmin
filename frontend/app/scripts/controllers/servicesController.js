'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller: ServicesController
 * @description
 * # MainCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .controller('ServicesController', ['$scope', '$http', '$timeout', function($scope,$http,$timeout) {
        var refreshServices = function() {
            $http.get("http://localhost:8080/api/v1/services").then(function(resp) {
                $scope.services = resp.data;
            });
        };
        refreshServices();

        var refreshProcesses = function() {
            $http.get("http://localhost:8080/api/v1/processes").then(function(resp) {
                $scope.processes = resp.data;
            });
        };
        refreshProcesses();

        var refreshHosts = function() {
            $http.get("http://localhost:8080/api/v1/hosts").then(function(resp) {
                $scope.hosts = resp.data;
            });
        };
        refreshHosts();

        $scope.spawnProcess = function() {
          alert("here!");
        };

    }]);

'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller: ServicesController
 * @description
 * # MainCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .controller('ServicesController', ['$scope', '$http', '$timeout', '$modal', function($scope, $http, $timeout, $modal) {

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

        $scope.openNewProcessDialog = function() {
            var modalInstance = $modal.open({
                animation: $scope.animationsEnabled,
                templateUrl: 'NewProcessModal.html',
                resolve: {
                    services: function() {
                        return $scope.services;
                    },
                    hosts: function() {
                        return $scope.hosts;
                    }
                },
                controller: function($scope, $modalInstance, services, hosts) {
                    $scope.services = services;
                    $scope.hosts = hosts;
                    $scope.newProcData = {};

                    $scope.ok = function() {
                        if ($scope.newProcData.serviceName && $scope.newProcData.machineID) {
                            $http.post("http://localhost:8080/api/v1/processes", {
                                svcName: $scope.newProcData.serviceName,
                                machID: $scope.newProcData.machineID
                            }).then(function(resp){
                                console.log(resp);
                                $modalInstance.close();
                            });
                        } else {
                            alert("invalid selection");
                        }
                    };

                    $scope.cancel = function() {
                        $modalInstance.dismiss('cancel');
                    };
                }
            });
        };
    }]);

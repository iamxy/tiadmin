'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller: ServicesController
 * @description
 * # MainCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .filter('processFilter', function() {
        return function(items, svcName) {
            var filtered = [];
            if (svcName === undefined || svcName === '') {
                return items;
            }
            angular.forEach(items, function(item) {
                if (svcName === item.svcName) {
                    filtered.push(item);
                }
            });
            return filtered;
        };
    })
    .controller('ServicesController', ['$scope', '$http', '$timeout', '$modal', function($scope, $http, $timeout, $modal) {
        var refresh = function() {
            $http.get("api/v1/services").then(function(resp) {
                $scope.services = resp.data;
            });
            $http.get("api/v1/processes").then(function(resp) {
                $scope.processes = resp.data;
            });
            $http.get("api/v1/hosts").then(function(resp) {
                $scope.hosts = resp.data;
            });
        };
        refresh();
        setInterval(refresh, 3000);

        // start process
        $scope.start = function(p) {
            $http.get("api/v1/processes/" + p.procID + "/start").then(function(resp) {
                refresh();
            });
        };

        $scope.stop = function(p) {
            $http.get("api/v1/processes/" + p.procID + "/stop").then(function(resp) {
                refresh();
            });
        };

        $scope.delete = function(p) {
            $http.delete("api/v1/processes/" + p.procID).then(function(resp) {
                refresh();
            });
        };

        // new process dialog
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
                            // create process
                            $http.post("api/v1/processes", {
                                svcName: $scope.newProcData.serviceName,
                                machID: $scope.newProcData.machineID,
                                desiredState: "started"
                            }).then(function(resp) {
                                refresh();
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

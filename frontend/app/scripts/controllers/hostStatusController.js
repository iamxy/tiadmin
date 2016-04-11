'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller:HostStatusCtrl
 * @description
 * # HostStatusCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .controller('HostStatusCtrl', ['$stateParams', '$scope', '$http', '$timeout', '$modal', function($stateParams, $scope, $http, $timeout, $modal) {
        $scope.machID = $stateParams.machID;
        var refresh = function() {
            $http.get("api/v1/services").then(function(resp) {
                $scope.services = resp.data;
            });
            $http.get("api/v1/hosts/" + $scope.machID).then(function(resp) {
                $scope.host = resp.data;
            });
            $http.get("api/v1/processes/findByHost?machID=" + $scope.machID).then(function(resp) {
                $scope.processes = resp.data;
            });
        }
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
                    },
                    machID: function() {
                        return $scope.machID;
                    }
                },
                controller: function($scope, $modalInstance, services, hosts, machID) {
                    $scope.services = services;
                    $scope.hosts = hosts;
                    $scope.machID = machID;
                    $scope.newProcData = {};

                    $scope.ok = function() {
                        if ($scope.newProcData.serviceName) {
                            // create process
                            $http.post("api/v1/processes", {
                                svcName: $scope.newProcData.serviceName,
                                machID: $scope.machID,
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

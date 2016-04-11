'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .controller('HomeCtrl', function($scope, $position, $http) {
        $scope.options = {
            chart: {
                type: 'lineChart',
                height: 180,
                margin: {
                    top: 20,
                    right: 20,
                    bottom: 40,
                    left: 55
                },
                x: function(d) { return d.x; },
                y: function(d) { return d.y; },
                useInteractiveGuideline: true,
                duration: 0,
                //yDomain: [-10, 10],
                yAxis: {
                    tickFormat: function(d) {
                        return d3.format('.01f')(d);
                    }
                }
            }
        };

        $scope.data = [{
            values: [],
            key: 'TPS',
        }];


        var x = 0;
        setInterval(function() {
          $http.get("api/v1/monitor/real/tidb_perf").then(function(resp) {
            $scope.data[0].values.push({ x: x, y: resp.data.tps });
            if ($scope.data[0].values.length > 20)
                $scope.data[0].values.shift();
            x++;
          });
        }, 1000);

        var refreshNodes = function() {
            $http.get("api/v1/hosts").then(function(resp) {
                $scope.hosts = resp.data;
                $scope.numOfNodes = resp.data.filter(function(x) {
                    return x.isAlive }).length;
            });
            $http.get("api/v1/monitor/real/tikv_storage").then(function(resp){
                $scope.storageInfo = resp.data;
            });
        };
        refreshNodes();
        setInterval(refreshNodes, 3000);

        var refreshServices = function() {
            $http.get("api/v1/services").then(function(resp) {
                $scope.services = resp.data;
            });
        };
        refreshServices();

    });

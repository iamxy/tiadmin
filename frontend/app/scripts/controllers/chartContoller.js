'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller:ChartCtrl
 * @description
 * # ChartCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .controller('ChartCtrl', ['$scope', '$http', '$timeout', function($scope, $http, $timeout) {
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
                x: function(d) {
                    return d.x; },
                y: function(d) {
                    return d.y; },
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

        $scope.tps = [{
            values: [],
            key: 'TPS',
        }];

        $scope.qps = [{
            values: [],
            key: 'QPS',
        }];

        $scope.iops = [{
            values: [],
            key: 'IOPS',
        }];

        $scope.conns = [{
            values: [],
            key: 'Conn Number',
        }];

        var x = 0;
        var refresh = function() {
            $http.get("http://localhost:8080/api/v1/monitor/real/tidb_perf").then(function(resp) {
                $scope.tps[0].values.push({ x: x, y: resp.data.tps });
                if ($scope.tps[0].values.length > 20)
                    $scope.tps[0].values.shift();

                $scope.qps[0].values.push({ x: x, y: resp.data.qps });
                if ($scope.qps[0].values.length > 20)
                    $scope.qps[0].values.shift();

                $scope.iops[0].values.push({ x: x, y: resp.data.iops });
                if ($scope.iops[0].values.length > 20)
                    $scope.iops[0].values.shift();

                $scope.conns[0].values.push({ x: x, y: resp.data.conns });
                if ($scope.conns[0].values.length > 20)
                    $scope.conns[0].values.shift();

                x++;
            });
        };
        refresh();

        setInterval(function() {
            refresh();
        }, 1000);


        
    }]);

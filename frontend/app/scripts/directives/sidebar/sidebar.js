'use strict';

/**
 * @ngdoc directive
 * @name izzyposWebApp.directive:adminPosHeader
 * @description
 * # adminPosHeader
 */

angular.module('tiAdminApp')
    .directive('sidebar', ['$location', function($interval, $http) {
        return {
            templateUrl: 'scripts/directives/sidebar/sidebar.html',
            restrict: 'E',
            replace: true,
            scope: {},
            controller: function($scope,$interval, $http) {
                $scope.selectedMenu = 'dashboard';
                $scope.collapseVar = 0;
                $scope.multiCollapseVar = 0;

                var refreshNodes = function() {
                    $http.get("http://localhost:8080/api/v1/hosts").then(function(resp) {
                        $scope.hosts = resp.data
                    });
                };
                refreshNodes();
                setInterval(refreshNodes, 5000);

                $scope.check = function(x) {
                    if (x == $scope.collapseVar)
                        $scope.collapseVar = 0;
                    else
                        $scope.collapseVar = x;
                };

                $scope.multiCheck = function(y) {
                    if (y == $scope.multiCollapseVar)
                        $scope.multiCollapseVar = 0;
                    else
                        $scope.multiCollapseVar = y;
                };
            }
        }
    }]);

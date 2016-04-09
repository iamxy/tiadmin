'use strict';

/**
 * @ngdoc directive
 * @description
 * # process list
 */

angular.module('tiAdminApp')
    .directive('processList', ['$location','DataService', 'ProcessService', function($interval, $http, DataService, ProcessService) {
        return {
            templateUrl: 'scripts/directives/processlist/processlist.html',
            restrict: 'E',
            replace: true,
            scope: {},
            controller: function($scope,$interval, $http, DataService, ProcessService) {
            }
        }
    }]);

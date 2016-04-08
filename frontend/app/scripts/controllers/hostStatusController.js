'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller:HostStatusCtrl
 * @description
 * # HostStatusCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
    .controller('HostStatusCtrl', ['$stateParams', '$scope', '$http', '$timeout', function($stateParams, $scope, $http, $timeout) {
        $scope.machID = $stateParams.machID;

        var refresh = function() {
            $http.get("http://localhost:8080/api/v1/hosts/" + $scope.machID).then(function(resp) {
                $scope.data = resp.data;
                console.log(resp.data);
            });

            $http.get("http://localhost:8080/api/v1/processes/findByHost?machID=" + $scope.machID).then(function(resp) {
                $scope.processes = resp.data;
                console.log(resp.data);
            });
        }
        refresh();

    }]);

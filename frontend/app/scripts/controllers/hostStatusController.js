'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
  .controller('HostStatusCtrl', ['$stateParams', '$scope', '$http', '$timeout', function ($stateParams, $scope,$http,$timeout) {
    $scope.machID = $stateParams.machID;
    $http.get("http://localhost:8080/api/v1/hosts/" + $scope.machID).then(function(resp) {
        $scope.data = resp.data;
    });
  }]);

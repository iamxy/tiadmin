'use strict';
/**
 * @ngdoc function
 * @name tiAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the tiAdminApp
 */
angular.module('tiAdminApp')
  .controller('HostStatusCtrl', ['$stateParams', '$scope', '$timeout', function ($stateParams, $scope, $timeout) {
    $scope.hostName = $stateParams.hostName;
  }]);

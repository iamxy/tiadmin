'use strict';
/**
 * @ngdoc function
 * @name sbAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the sbAdminApp
 */
angular.module('sbAdminApp')
  .controller('HostStatusCtrl', ['$stateParams', '$scope', '$timeout', function ($stateParams, $scope, $timeout) {
    $scope.hostName = $stateParams.hostName;
  }]);
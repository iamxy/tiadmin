'use strict';
/**
 * @ngdoc function
 * @name sbAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the sbAdminApp
 */
angular.module('sbAdminApp')
  .controller('ChartCtrl', ['$scope', '$timeout', function ($scope, $timeout) {
    $scope.options = {
        chart: {
            type: 'lineChart',
            height: 180,
            margin : {
                top: 20,
                right: 20,
                bottom: 40,
                left: 55
            },
            x: function(d){ return d.x; },
            y: function(d){ return d.y; },
            useInteractiveGuideline: true,
            duration: 0,
            yDomain: [-10,10],    
            yAxis: {
                tickFormat: function(d){
                   return d3.format('.01f')(d);
                }
            }
        }
    };

    $scope.data = [{
        values: [],
        key: 'Random Walk',
    }];

    $scope.run = true;
    
    var x = 0;
    setInterval(function(){
        if (!$scope.run) return;
        $scope.data[0].values.push({ x: x,  y: Math.random() - 0.5});
      if ($scope.data[0].values.length > 20) $scope.data[0].values.shift();
        x++;
        
      $scope.$apply(); // update both chart
    }, 1000);
  }]);
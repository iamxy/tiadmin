'use strict';
/**
 * @ngdoc overview
 * @name tiAdminApp
 * @description
 * # tiAdminApp
 *
 * Main module of the application.
 */

angular.module('tiAdminApp')
    .factory('DataService', ['$http', function($http) {
        var i = 0;
        var service = {
            add: function() { i += 1; },
            get: function() {
                return i; }
        };
        return service;
    }])
    .factory('HostService', ['$http', function($http) {
        var hosts = []; 
        var refreshNodes = function() {
            $http.get("api/v1/hosts").then(function(resp) {
                console.log("load hosts");
                hosts = resp.data
            });
        };
        refreshNodes();
        setInterval(refreshNodes, 3000);

        var service = {
            refresh: function() {
                refreshNodes();
            },
            getHosts: function() {
                return hosts;
            }
        };
        return service;
    }])

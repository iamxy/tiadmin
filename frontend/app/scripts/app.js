'use strict';
/**
 * @ngdoc overview
 * @name tiAdminApp
 * @description
 * # tiAdminApp
 *
 * Main module of the application.
 */
angular
    .module('tiAdminApp', [
        'oc.lazyLoad',
        'ui.router',
        'ui.bootstrap',
        'nvd3',
        'angular-loading-bar',
    ])
    .config(['$stateProvider', '$urlRouterProvider', '$ocLazyLoadProvider', function($stateProvider, $urlRouterProvider, $ocLazyLoadProvider) {

        $ocLazyLoadProvider.config({
            debug: false,
            events: true,
        });

        $urlRouterProvider.otherwise('/dashboard/home');

        $stateProvider
            .state('dashboard', {
                url: '/dashboard',
                templateUrl: 'views/dashboard/main.html',
                resolve: {
                    loadMyDirectives: function($ocLazyLoad) {
                        return $ocLazyLoad.load({
                                name: 'tiAdminApp',
                                files: [
                                    'scripts/services.js',
                                    'scripts/directives/header/header.js',
                                    'scripts/directives/sidebar/sidebar.js',
                                ]
                            }),
                            $ocLazyLoad.load({
                                name: 'toggle-switch',
                                files: ["bower_components/angular-toggle-switch/angular-toggle-switch.min.js",
                                    "bower_components/angular-toggle-switch/angular-toggle-switch.css"
                                ]
                            }),
                            $ocLazyLoad.load({
                                name: 'ngAnimate',
                                files: ['bower_components/angular-animate/angular-animate.js']
                            }),
                            $ocLazyLoad.load({
                                name: 'ngCookies',
                                files: ['bower_components/angular-cookies/angular-cookies.js']
                            }),
                            $ocLazyLoad.load({
                                name: 'ngResource',
                                files: ['bower_components/angular-resource/angular-resource.js']
                            }),
                            $ocLazyLoad.load({
                                name: 'ngSanitize',
                                files: ['bower_components/angular-sanitize/angular-sanitize.js']
                            })
                    }
                }
            })
            .state('dashboard.home', {
                url: '/home',
                controller: 'MainCtrl',
                templateUrl: 'views/dashboard/home.html',
                resolve: {
                    loadMyFiles: function($ocLazyLoad) {
                        return $ocLazyLoad.load({
                            name: 'tiAdminApp',
                            files: [
                                'scripts/controllers/main.js',
                            ]
                        })
                    }
                }
            })
            .state('dashboard.host-status', {
                controller: 'HostStatusCtrl',
                templateUrl: 'views/dashboard/host_status.html',
                url: '/host-status/:machID',
                resolve: {
                    loadMyFile: function($ocLazyLoad) {
                        return $ocLazyLoad.load({
                            name: 'tiAdminApp',
                            files: [
                                'scripts/directives/processlist/processlist.js',
                                'scripts/controllers/hostStatusController.js'
                            ]
                        })
                    }
                }
            })
            .state('dashboard.chart', {
                templateUrl: 'views/dashboard/chart.html',
                url: '/chart',
                controller: 'ChartCtrl',
                resolve: {
                    loadMyFile: function($ocLazyLoad) {
                        return $ocLazyLoad.load({
                                name: 'chart.js',
                                files: [
                                    'bower_components/angular-chart.js/dist/angular-chart.min.js',
                                    'bower_components/angular-chart.js/dist/angular-chart.css',
                                ]
                            }),
                            $ocLazyLoad.load({
                                name: 'tiAdminApp',
                                files: ['scripts/controllers/chartContoller.js']
                            })
                    }
                }
            })
            .state('dashboard.services', {
                controller: 'ServicesController',
                templateUrl: 'views/dashboard/services.html',
                url: '/services',
                resolve: {
                    loadMyFile: function($ocLazyLoad) {
                        return $ocLazyLoad.load({
                            name: 'tiAdminApp',
                            files: [
                                'scripts/controllers/servicesController.js',
                            ]
                        })
                    }
                }
            })
    }]);

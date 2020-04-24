(function (routeConfig) {

  'use strict';

  routeConfig.init = function (app) {

    // *** routes *** //
    const authRoutes = require('../routes/auth');
    const userRoutes = require('../routes/user');

    // *** register routes *** //
    app.use('/auth', authRoutes);
    app.use('/users', userRoutes);

  };

})(module.exports);

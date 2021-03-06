import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    const activationToken = params.token;
    const response = this.get('store').queryRecord('activate', { token: activationToken }).then(function() {}, function(adapterError) {
      console.log(adapterError);
      // response.errors = adapterError.errors.toArray();
    });
    return response;
  }
});

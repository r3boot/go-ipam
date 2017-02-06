import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    const activationToken = params.token;
    const response = this.get('store').query('activate', { token: activationToken });
    return response;
  }
});

import Ember from 'ember';

export default Ember.Controller.extend({
  username: '',
  password: '',

  isValid: Ember.computed.match('username', /^[a-z0-9]{3,32}$/) &&
           Ember.computed.match('password', /^.{8,128}$/),
  isDisabled: Ember.computed.not('isValid'),

  actions: {
    getSessionToken() {
      const username = this.get('username');
      const password = this.get('password');

      const response = this.store.queryRecord('auth', {
        username: username,
        password: password,
      });

      console.log(response);
    }
  }
});

import Ember from 'ember';

export default Ember.Controller.extend({

  username: '',
  fullname: '',
  email: '',
  password: '',
  verify: '',

  isValid: Ember.computed.match('username', /^[a-z0-9]{3,32}$/) &&
           Ember.computed.match('fullname', /^[a-zA-Z0-9\ -]{3,64}$/) &&
           Ember.computed.match('email', '/^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$/') &&
           Ember.computed.match('password', /^.{8,128}$/) &&
           Ember.computed.match('verify', /^.{8,128}$/),

  isDisabled: Ember.computed.not('isValid'),

  actions: {
    newAccountRegistration() {
      if (this.get('password') !== this.get('verify')) {
        this.set('passwordMismatch', `Passwords do not match`);
        return;
      }
      this.set('passwordMismatch', ``);

      const username = this.get('username');
      const fullname = this.get('fullname');
      const email = this.get('email');
      const password = this.get('password');

      const newAccount = this.store.createRecord('signup', {
        username: username,
        fullname: fullname,
        email: email,
        password: password,
      });
      newAccount.save();

      this.set('responseMessage', `An email will be sent to ${this.get('email')} containing a link to activate your account`);

    }
  }
});

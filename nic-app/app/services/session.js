import Ember from 'ember';

export default Ember.Service.extend({
  sessionKey: null,

  init() {
    this._super(...arguments);
    this.set('sessionKey', '');
  },

  setSessionKey(key) {
    this.set('sessionKey', key);
  },

  unsetSessionKey() {
    this.set('sessionKey', '');
  },
});

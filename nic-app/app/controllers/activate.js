import Ember from 'ember';

export default Ember.Controller.extend({

  validToken: Ember.computed.equal('model.token', true)
});

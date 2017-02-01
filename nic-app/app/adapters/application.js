import DS from 'ember-data';
import Ember from 'ember';

export default DS.JSONAPIAdapter.extend({
  namespace: '/v1',
  host: 'http://127.0.0.1:8080',
  pathForType: function(type) {
    return Ember.String.htmlSafe(type);
  }
});

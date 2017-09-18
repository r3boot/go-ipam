import DS from 'ember-data';
import Ember from 'ember';

export default DS.RESTAdapter.extend({
  namespace: '/v1',
  pathForType: function(type) {
    return Ember.String.htmlSafe(type);
  }
});

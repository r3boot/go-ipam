import { test } from 'qunit';
import moduleForAcceptance from 'nic-app/tests/helpers/module-for-acceptance';

moduleForAcceptance('Acceptance | signup');

test('visiting /signup', function(assert) {
  visit('/signup');

  andThen(function() {
    assert.equal(currentURL(), '/signup');
  });
});

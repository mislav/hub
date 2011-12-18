require 'test_helper'

class CheckoutTest < Test::Unit::TestCase
  def test_checkout_no_changes
    assert_forwarded "checkout master"
  end

  def test_checkout_pullrequest
    stub_request(:get, "https://#{auth}github.com/api/v2/json/pulls/defunkt/hub/73").
    to_return(:body => mock_pull_response('blueyed:feature'))

    assert_commands 'git remote add -f -t feature blueyed git://github.com/blueyed/hub.git',
      'git checkout -b blueyed-feature blueyed/feature',
      "checkout https://github.com/defunkt/hub/pull/73/files"
  end

  def test_checkout_pullrequest_custom_branch
    stub_request(:get, "https://#{auth}github.com/api/v2/json/pulls/defunkt/hub/73").
    to_return(:body => mock_pull_response('blueyed:feature'))

    assert_commands 'git remote add -f -t feature blueyed git://github.com/blueyed/hub.git',
      'git checkout -b review blueyed/feature',
      "checkout https://github.com/defunkt/hub/pull/73/files review"
  end

  def test_checkout_pullrequest_existing_remote
    stub_command_output 'remote', "origin\nblueyed"

    stub_request(:get, "https://#{auth}github.com/api/v2/json/pulls/defunkt/hub/73").
    to_return(:body => mock_pull_response('blueyed:feature'))

    assert_commands 'git remote set-branches --add blueyed feature',
      'git fetch blueyed +refs/heads/feature:refs/remotes/blueyed/feature',
      'git checkout -b blueyed-feature blueyed/feature',
      "checkout https://github.com/defunkt/hub/pull/73/files"
  end
end

require 'test_helper'

class PullRequestTest < Test::Unit::TestCase
  def test_pullrequest
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "tpw:master", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -f"
  end

  def test_pullrequest_with_checks
    @git["rev-list --cherry origin/master..."] = "+abcd1234\n+bcde2345"

    expected = "Aborted: 2 commits are not yet pushed to origin/master\n" <<
      "(use `-f` to force submit a pull request anyway)\n"
    assert_output expected, "pull-request hereyougo"
  end

  def test_pullrequest_from_branch
    stub_branch('refs/heads/feature')
    stub_tracking_nothing('feature')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "tpw:feature", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -f"
  end

  def test_pullrequest_from_tracking_branch
    stub_branch('refs/heads/feature')
    stub_tracking('feature', 'tpw', 'yay-feature')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "tpw:yay-feature", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -f"
  end

  def test_pullrequest_explicit_head
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "tpw:yay-feature", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -h yay-feature -f"
  end

  def test_pullrequest_explicit_head_with_owner
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "mojombo:feature", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -h mojombo:feature -f"
  end

  def test_pullrequest_explicit_base
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "feature", 'head' => "tpw:master", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -b feature -f"
  end

  def test_pullrequest_explicit_base_with_owner
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/mojombo/hub").
      with(:body => { 'pull' => {'base' => "feature", 'head' => "tpw:master", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -b mojombo:feature -f"
  end

  def test_pullrequest_explicit_base_with_repo
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/mojombo/hubbub").
      with(:body => { 'pull' => {'base' => "feature", 'head' => "tpw:master", 'title' => "hereyougo"} }).
      to_return(:body => mock_pullreq_response(1))

    expected = "https://github.com/defunkt/hub/pull/1\n"
    assert_output expected, "pull-request hereyougo -b mojombo/hubbub:feature -f"
  end

  def test_pullrequest_existing_issue
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/defunkt/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "tpw:master", 'issue' => '92'} }).
      to_return(:body => mock_pullreq_response(92))

    expected = "https://github.com/defunkt/hub/pull/92\n"
    assert_output expected, "pull-request #92 -f"
  end

  def test_pullrequest_existing_issue_url
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/pulls/mojombo/hub").
      with(:body => { 'pull' => {'base' => "master", 'head' => "tpw:master", 'issue' => '92'} }).
      to_return(:body => mock_pullreq_response(92, 'mojombo/hub'))

    expected = "https://github.com/mojombo/hub/pull/92\n"
    assert_output expected, "pull-request https://github.com/mojombo/hub/issues/92#comment_4 -f"
  end
end

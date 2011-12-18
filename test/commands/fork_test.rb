require 'test_helper'

class ForkTest < Test::Unit::TestCase
  def test_fork
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/fork/defunkt/hub")

    expected = "remote add -f tpw git@github.com:tpw/hub.git\n"
    expected << "new remote: tpw\n"
    assert_equal expected, hub("fork") { ENV['GIT'] = 'echo' }
  end

  def test_fork_failed
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/fork/defunkt/hub").
      to_return(:status => [500, "Your fork is fail"])

    expected = "Error creating fork: Your fork is fail (HTTP 500)\n"
    assert_equal expected, hub("fork") { ENV['GIT'] = 'echo' }
  end

  def test_fork_no_remote
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/fork/defunkt/hub")

    assert_equal "", hub("fork --no-remote") { ENV['GIT'] = 'echo' }
  end

  def test_fork_already_exists
    stub_existing_fork('tpw')

    expected = "tpw/hub already exists on GitHub\n"
    expected << "remote add -f tpw git@github.com:tpw/hub.git\n"
    expected << "new remote: tpw\n"
    assert_equal expected, hub("fork") { ENV['GIT'] = 'echo' }
  end

  def test_fork_https_protocol
    stub_existing_fork('tpw')
    stub_https_is_preferred

    expected = "tpw/hub already exists on GitHub\n"
    expected << "remote add -f tpw https://github.com/tpw/hub.git\n"
    expected << "new remote: tpw\n"
    assert_equal expected, hub("fork") { ENV['GIT'] = 'echo' }
  end
end

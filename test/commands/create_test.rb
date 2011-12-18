require 'test_helper'

class CreateTest < Test::Unit::TestCase
  def test_create
    stub_no_remotes
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'hub' })

    expected = "remote add -f origin git@github.com:tpw/hub.git\n"
    expected << "created repository: tpw/hub\n"
    assert_equal expected, hub("create") { ENV['GIT'] = 'echo' }
  end

  def test_create_custom_name
    stub_no_remotes
    stub_nonexisting_fork('tpw', 'hubbub')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'hubbub' })

    expected = "remote add -f origin git@github.com:tpw/hubbub.git\n"
    expected << "created repository: tpw/hubbub\n"
    assert_equal expected, hub("create hubbub") { ENV['GIT'] = 'echo' }
  end

  def test_create_in_organization
    stub_no_remotes
    stub_nonexisting_fork('acme', 'hubbub')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'acme/hubbub' })

    expected = "remote add -f origin git@github.com:acme/hubbub.git\n"
    expected << "created repository: acme/hubbub\n"
    assert_equal expected, hub("create acme/hubbub") { ENV['GIT'] = 'echo' }
  end

  def test_create_no_openssl
    stub_no_remotes
    # stub_nonexisting_fork('tpw')
    stub_request(:get, "http://#{auth}github.com/api/v2/yaml/repos/show/tpw/hub").
      to_return(:status => 404)

    stub_request(:post, "http://#{auth}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'hub' })

    expected = "remote add -f origin git@github.com:tpw/hub.git\n"
    expected << "created repository: tpw/hub\n"

    assert_equal expected, hub("create") {
      ENV['GIT'] = 'echo'
      require 'net/https'
      Object.send :remove_const, :OpenSSL
    }
  end

  def test_create_failed
    stub_no_remotes
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").
      to_return(:status => [401, "Your token is fail"])

    expected = "Error creating repository: Your token is fail (HTTP 401)\n"
    expected << "Check your token configuration (`git config github.token`)\n"
    assert_equal expected, hub("create") { ENV['GIT'] = 'echo' }
  end

  def test_create_with_env_authentication
    stub_no_remotes

    old_user  = ENV['GITHUB_USER']
    old_token = ENV['GITHUB_TOKEN']
    ENV['GITHUB_USER']  = 'mojombo'
    ENV['GITHUB_TOKEN'] = '123abc'

    # stub_nonexisting_fork('mojombo')
    stub_request(:get, "https://#{auth('mojombo', '123abc')}github.com/api/v2/yaml/repos/show/mojombo/hub").
      to_return(:status => 404)

    stub_request(:post, "https://#{auth('mojombo', '123abc')}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'hub' })

    expected = "remote add -f origin git@github.com:mojombo/hub.git\n"
    expected << "created repository: mojombo/hub\n"
    assert_equal expected, hub("create") { ENV['GIT'] = 'echo' }

  ensure
    ENV['GITHUB_USER']  = old_user
    ENV['GITHUB_TOKEN'] = old_token
  end

  def test_create_private_repository
    stub_no_remotes
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'hub', 'public' => '0' })

    expected = "remote add -f origin git@github.com:tpw/hub.git\n"
    expected << "created repository: tpw/hub\n"
    assert_equal expected, hub("create -p") { ENV['GIT'] = 'echo' }
  end

  def test_create_with_description_and_homepage
    stub_no_remotes
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").with(:body => {
      'name' => 'hub', 'description' => 'toyproject', 'homepage' => 'http://example.com'
    })

    expected = "remote add -f origin git@github.com:tpw/hub.git\n"
    expected << "created repository: tpw/hub\n"
    assert_equal expected, hub("create -d toyproject -h http://example.com") { ENV['GIT'] = 'echo' }
  end

  def test_create_with_invalid_arguments
    assert_equal "invalid argument: -a\n",   hub("create -a blah")   { ENV['GIT'] = 'echo' }
    assert_equal "invalid argument: bleh\n", hub("create blah bleh") { ENV['GIT'] = 'echo' }
  end

  def test_create_with_existing_repository
    stub_no_remotes
    stub_existing_fork('tpw')

    expected = "tpw/hub already exists on GitHub\n"
    expected << "remote add -f origin git@github.com:tpw/hub.git\n"
    expected << "set remote origin: tpw/hub\n"
    assert_equal expected, hub("create") { ENV['GIT'] = 'echo' }
  end

  def test_create_https_protocol
    stub_no_remotes
    stub_existing_fork('tpw')
    stub_https_is_preferred

    expected = "tpw/hub already exists on GitHub\n"
    expected << "remote add -f origin https://github.com/tpw/hub.git\n"
    expected << "set remote origin: tpw/hub\n"
    assert_equal expected, hub("create") { ENV['GIT'] = 'echo' }
  end

  def test_create_no_user
    stub_no_remotes
    out = hub("create") do
      stub_github_token(nil)
    end
    assert_equal "** No GitHub token set. See http://help.github.com/set-your-user-name-email-and-github-token/\n", out
  end

  def test_create_outside_git_repo
    stub_no_git_repo
    assert_equal "'create' must be run from inside a git repository\n", hub("create")
  end

  def test_create_origin_already_exists
    stub_nonexisting_fork('tpw')
    stub_request(:post, "https://#{auth}github.com/api/v2/yaml/repos/create").
      with(:body => { 'name' => 'hub' })

    expected = "remote -v\ncreated repository: tpw/hub\n"
    assert_equal expected, hub("create") { ENV['GIT'] = 'echo' }
  end
end

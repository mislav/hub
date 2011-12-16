require 'test_helper'

class AmTest < Test::Unit::TestCase
  def test_am_untouched
    assert_forwarded "am some.patch"
  end

  def test_am_pull_request
    with_tmpdir('/tmp/') do
      assert_commands "curl -#LA 'hub #{Hub::Version}' https://github.com/defunkt/hub/pull/55.patch -o /tmp/55.patch",
                      "git am --signoff /tmp/55.patch -p2",
                      "am --signoff https://github.com/defunkt/hub/pull/55 -p2"

      cmd = Hub("am https://github.com/defunkt/hub/pull/55/files").command
      assert_includes '/pull/55.patch', cmd
    end
  end

  def test_am_commit_url
    with_tmpdir('/tmp/') do
      url = 'https://github.com/davidbalbert/hub/commit/fdb9921'

      assert_commands "curl -#LA 'hub #{Hub::Version}' #{url}.patch -o /tmp/fdb9921.patch",
                      "git am --signoff /tmp/fdb9921.patch -p2",
                      "am --signoff #{url} -p2"
    end
  end

  def test_am_gist
    with_tmpdir('/tmp/') do
      url = 'https://gist.github.com/8da7fb575debd88c54cf'

      assert_commands "curl -#LA 'hub #{Hub::Version}' #{url}.txt -o /tmp/gist-8da7fb575debd88c54cf.txt",
                      "git am --signoff /tmp/gist-8da7fb575debd88c54cf.txt -p2",
                      "am --signoff #{url} -p2"
    end
  end
end

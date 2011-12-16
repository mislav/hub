require 'test_helper'

class ApplyTest < Test::Unit::TestCase
  def test_apply_untouched
    assert_forwarded "apply some.patch"
  end

  def test_apply_pull_request
    with_tmpdir('/tmp/') do
      assert_commands "curl -#LA 'hub #{Hub::Version}' https://github.com/defunkt/hub/pull/55.patch -o /tmp/55.patch",
                      "git apply /tmp/55.patch -p2",
                      "apply https://github.com/defunkt/hub/pull/55 -p2"

      cmd = Hub("apply https://github.com/defunkt/hub/pull/55/files").command
      assert_includes '/pull/55.patch', cmd
    end
  end

  def test_apply_commit_url
    with_tmpdir('/tmp/') do
      url = 'https://github.com/davidbalbert/hub/commit/fdb9921'

      assert_commands "curl -#LA 'hub #{Hub::Version}' #{url}.patch -o /tmp/fdb9921.patch",
                      "git apply /tmp/fdb9921.patch -p2",
                      "apply #{url} -p2"
    end
  end

  def test_apply_gist
    with_tmpdir('/tmp/') do
      url = 'https://gist.github.com/8da7fb575debd88c54cf'

      assert_commands "curl -#LA 'hub #{Hub::Version}' #{url}.txt -o /tmp/gist-8da7fb575debd88c54cf.txt",
                      "git apply /tmp/gist-8da7fb575debd88c54cf.txt -p2",
                      "apply #{url} -p2"
    end
  end
end

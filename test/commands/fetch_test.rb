require 'test_helper'

class FetchTest < Test::Unit::TestCase
  def test_fetch_existing_remote
    assert_forwarded "fetch mislav"
  end

  def test_fetch_new_remote
    stub_remotes_group('xoebus', nil)
    stub_existing_fork('xoebus')

    assert_commands "git remote add xoebus git://github.com/xoebus/hub.git",
                    "git fetch xoebus",
                    "fetch xoebus"
  end

  def test_fetch_new_remote_with_options
    stub_remotes_group('xoebus', nil)
    stub_existing_fork('xoebus')

    assert_commands "git remote add xoebus git://github.com/xoebus/hub.git",
                    "git fetch --depth=1 --prune xoebus",
                    "fetch --depth=1 --prune xoebus"
  end

  def test_fetch_multiple_new_remotes
    stub_remotes_group('xoebus', nil)
    stub_remotes_group('rtomayko', nil)
    stub_existing_fork('xoebus')
    stub_existing_fork('rtomayko')

    assert_commands "git remote add xoebus git://github.com/xoebus/hub.git",
                    "git remote add rtomayko git://github.com/rtomayko/hub.git",
                    "git fetch --multiple xoebus rtomayko",
                    "fetch --multiple xoebus rtomayko"
  end

  def test_fetch_multiple_comma_separated_remotes
    stub_remotes_group('xoebus', nil)
    stub_remotes_group('rtomayko', nil)
    stub_existing_fork('xoebus')
    stub_existing_fork('rtomayko')

    assert_commands "git remote add xoebus git://github.com/xoebus/hub.git",
                    "git remote add rtomayko git://github.com/rtomayko/hub.git",
                    "git fetch --multiple xoebus rtomayko",
                    "fetch xoebus,rtomayko"
  end

  def test_fetch_multiple_new_remotes_with_filtering
    stub_remotes_group('xoebus', nil)
    stub_remotes_group('mygrp', 'one two')
    stub_remotes_group('typo', nil)
    stub_existing_fork('xoebus')
    stub_nonexisting_fork('typo')

    # mislav: existing remote; skipped
    # xoebus: new remote, fork exists; added
    # mygrp:  a remotes group; skipped
    # URL:    can't be a username; skipped
    # typo:   fork doesn't exist; skipped
    assert_commands "git remote add xoebus git://github.com/xoebus/hub.git",
                    "git fetch --multiple mislav xoebus mygrp git://example.com typo",
                    "fetch --multiple mislav xoebus mygrp git://example.com typo"
  end
end

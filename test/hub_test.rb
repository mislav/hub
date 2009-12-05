require 'test/unit'
load File.dirname(__FILE__) + '/../bin/hub'

class HubTest < Test::Unit::TestCase
  def hub(args)
    args = args.split(' ')
    hub = Hub.new(*args)
    hub.send(args[0])
    hub
  end

  def test_private_clone
    h = hub("clone -p rtomayko/ron")
    assert_equal 'git@github.com:rtomayko/ron.git', h.args.last
  end

  def test_public_clone
    h = hub("clone rtomayko/ron")
    assert_equal 'git://github.com/rtomayko/ron.git', h.args.last
  end

  def test_private_remote
    h = hub("remote add -g -p rtomayko")
    assert_equal 'git@github.com:rtomayko/hub.git', h.args.last
  end

  def test_public_remote
    h = hub("remote add -g rtomayko")
    assert_equal 'git://github.com/rtomayko/hub.git', h.args.last
  end

  def test_init
    h = hub("init -g")
    assert_equal 'git remote add origin git@github.com:defunkt/hub.git', h.after
  end
end

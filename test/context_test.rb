require 'helper'

class ContextTest < Minitest::Test
  class Context
    include Hub::Context

    def initialize(&block)
      @git_reader = Hub::Context::GitReader.new('git', &block)
    end

    public :git_editor
  end

  attr_reader :context

  def setup
    super
    @stubs = {}
    @context = Context.new do |_, cmd|
      @stubs.fetch(cmd)
    end
  end

  def test_editor
    stub_command_output 'var GIT_EDITOR', 'vim'
    assert_equal %w'vim', context.git_editor
  end

  def test_editor_with_argument
    stub_command_output 'var GIT_EDITOR', 'subl -w'
    assert_equal %w'subl -w', context.git_editor
  end

  def test_editor_with_spaces
    stub_command_output 'var GIT_EDITOR', '"my editor" -w arg2'
    assert_equal %w'my\ editor -w arg2', context.git_editor
  end

  def test_editor_with_tilde
    stub_command_output 'var GIT_EDITOR', '~/bin/vi'
    with_env('HOME', '/home/mislav') do
      assert_equal %w'/home/mislav/bin/vi', context.git_editor
    end
  end

  def test_editor_with_env_variable
    stub_command_output 'var GIT_EDITOR', '$EDITOR'
    with_env('EDITOR', 'subl -w') do
      assert_equal %w'subl -w', context.git_editor
    end
  end

  def test_editor_with_embedded_env_variable
    stub_command_output 'var GIT_EDITOR', '$EDITOR -w'
    with_env('EDITOR', 'subl') do
      assert_equal %w'subl -w', context.git_editor
    end
  end

  def test_editor_with_curly_brackets_embedded_env_variable
    stub_command_output 'var GIT_EDITOR', 'my${EDITOR}2 -w'
    with_env('EDITOR', 'subl') do
      assert_equal %w'mysubl2 -w', context.git_editor
    end
  end

  private

  def stub_command_output(cmd, value)
    @stubs[cmd] = value.nil? ? nil : value.to_s
  end

  def with_env(name, value)
    dir, ENV[name] = ENV[name], value
    yield
  ensure
    ENV[name] = dir
  end
end

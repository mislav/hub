# Avoids over-zealous sanitize_text
# https://github.com/cucumber/aruba/blob/v1.0.4/lib/aruba/matchers/string/output_string_eq.rb

sanitize_text = ->(expected) {
  expected.to_s.
      # convert "\n" in expectations to literal newline, unless it is preceded by another backslash
      gsub(/(?<!\\)\\n/, "\n").
      # convert "\e" in expectations to a literal ESC, unless it is preceded by another backslash
      gsub(/(?<!\\)\\e/, "\e").
      gsub('\\\\', '\\')
}

RSpec::Matchers.define :output_string_eq do |expected|
  match do |actual|
    @expected = sanitize_text.(expected)
    @actual = actual.to_s
    @actual = extract_text(@actual) if aruba.config.remove_ansi_escape_sequences

    @expected == @actual
  end

  diffable

  description { "output string is eq: #{description_of(self.expected)}" }
end

RSpec::Matchers.define :have_output do |expected|
  match do |actual|
    @old_actual = actual

    unless @old_actual.respond_to? :output
      raise "Expected #{@old_actual} to respond to #output"
    end

    @old_actual.stop

    @actual = actual.output
    @actual = extract_text(@actual) if aruba.config.remove_ansi_escape_sequences

    expected === @actual
  end

  diffable

  description { "have output: #{description_of(expected)}" }

  failure_message do |_actual|
    "expected `#{@old_actual.commandline}` to #{description_of(expected)}\n" \
      "but was: #{description_of(@actual)}"
  end
end

RSpec::Matchers.alias_matcher :an_output_string_being_eq, :output_string_eq

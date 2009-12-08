# take control of the growl notifications
module GrowlHacks
  def growl(type, subject, body, *args, &block)
    case type
    when Kicker::GROWL_NOTIFICATIONS[:succeeded]
      puts subject = "Success"
      body = body.split("\n").last
    when Kicker::GROWL_NOTIFICATIONS[:failed]
      subject = "Failure"
      puts body
      body = body.split("\n").last
    else
      return nil
    end
    super(type, subject, body, *args, &block)
  end
end

Kicker.send :extend, GrowlHacks

# no logging
Kicker::Utils.module_eval do
  def log(message)
    nil
  end
end
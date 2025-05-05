function priToLevel(tag, timestamp, record)
    -- Convert the syslog priority to severity and facility
    priority = tonumber(record["pri"]);
    record["severity"]= priority % 8
    record["facility"]= math.floor(priority / 8)
    -- the 2 denotes that the message was modified
    return 2, timestamp, record
end
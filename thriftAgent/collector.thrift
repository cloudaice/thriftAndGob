namespace go translate


enum ResultCode 
{
    OK,
    TRY_LATER
}

struct LogEntry
{
    1: string hostname,
    2: string message
}

service proxyTrans {
    ResultCode Log(1: list<LogEntry> messages);
}

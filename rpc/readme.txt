header size = 1 ubyte
header = version(1 ubyte) + st(short) + status info(1 ubyte) + reserved(2 ubyte) + function hash(4/8字节)这里默认8字节 + sequence Id(uint16) + body size (ushort)
magic(2 ubyte) -- 修改 magic(3 ubyte)MIG
version = 0x01  

st = compress type(1 ubyte) + protocol type(1 ubyte)
status info = is nonblock(<<4) + status code(0xF)


data package = magic + header size + header + body    


compress {
    MRPC_CT_NO,
    MRPC_CT_DYNAIC,
    MRPC_CT_COMPRESS
}

protocol {
    MRPC_PT_JSON,
    MRPC_PT_XML,
    MRPC_PT_PROTOBUFF,
    MRPC_PT_FLATBUFF,
}
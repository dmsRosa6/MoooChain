#!lua name=utilslib

-- CONSTANTS

local CHAIN_NAME_KEYWORD = "BlockChainName"
local GENESIS_KEYWORD = "GenesisBlockHash"
local BLOCK_KEYWORD = "Block:"
local PREV_BLOCK_HASH_KEYWORD = "Block:prev:"
local LAST_HASH_KEYWORD = "LastHash"
local DEBUG_CHAIN_KEYWORD = "DebugChain"

-- HELPER FUNCTIONS

local function tohex(data)
    if not data or data == "" then return "" end
    local t = {}
    for i = 1, #data do
        t[i] = string.format("%02x", string.byte(data, i))
    end
    return table.concat(t)
end


local function build_chain()

	local lastHash = redis.call("GET", LAST_HASH_KEYWORD)
	
	if not lastHash then 
		return "Last Hash field does not exist"
	end
	
	local curr = redis.call("GET", lastHash)
	
	if not curr then
		return "Block does not exist"
	end
	
	while curr do
		redis.call("LPUSH", DEBUG_CHAIN_KEYWORD, curr)
		local prev = redis.call("GET", PREV_BLOCK_HASH_KEYWORD .. curr)
		if not prev or prev == "" then break end
		curr = prev
	end
	
	return "OK"
end

-- FUNCTIONS TO EXPORT

-- key: {}
-- args: {}
local function init_debug_chain(keys, args)

	local chainExists = redis.call('EXISTS', CHAIN_NAME_KEYWORD)
	
	-- If the chain doesn’t exist, nothing to build
	if chainExists == 1 then 
		local deleted = redis.call('DEL', DEBUG_CHAIN_KEYWORD)
		
		local res = build_chain()
		
		return res		
	else
		return "BlockChain does not exist"
	end
end

-- key: {}
-- ARGS:
-- lastHash: starting hash
-- count: number of blocks to fetch
local function iterate_chain(keys, args)
    
    if #args ~= 2 then 
        return redis.error_reply("Invalid number of args")
    end
    
    local lastHash = args[1]
    local count = tonumber(args[2])

    if not count or count <= 0 then
        return redis.error_reply("Invalid count")
    end

    local curr = lastHash

    if lastHash == "" then
        curr = redis.call('GET', LAST_HASH_KEYWORD)
        if not curr or curr == "" then
            return { "", {}, 0 }
        end
    end

    local result = {}
    local fetched = 0
    curr = tohex(curr)
    while curr and fetched < count do
        table.insert(result, curr)
        fetched = fetched + 1
        curr = tohex(redis.call('GET', PREV_BLOCK_HASH_KEYWORD .. curr))
        if not curr or curr == "" then break end
    end

    local more = 0
    if curr and curr ~= "" then
        more = 1
    end

    local nextHash = curr or ""
    return { nextHash, result, more }
end

-- REGISTER FUNCTIONS
redis.register_function('init_debug_chain', init_debug_chain)
redis.register_function('iterate_chain', iterate_chain)
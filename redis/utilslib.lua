#!lua name=utilslib

-- CONSTANTS

local CHAIN_NAME_KEYWORD = "BlockChainName"
local GENESIS_KEYWORD = "GenesisBlockHash"
local BLOCK_KEYWORD = "Block:"
local PREV_BLOCK_HASH_KEYWORD = "Block:prev:"
local LAST_HASH_KEYWORD = "LastHash"
local DEBUG_CHAIN_KEYWORD = "DebugChain"
local ITERATE_PAGE_SIZE = 500


-- HELPER FUNCTIONS

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
-- args: {cursorHash, count}
-- cursorHash = "" (empty string) to start from the chain tip (LastHash)
-- count      = number of blocks to fetch in this call
local function iterate_chain(keys, args)
    local cursor = args[1]
    local count  = tonumber(args[2])

    if not count or count <= 0 then
        return {err = "Invalid count"}
    end
	
    local curr
    if cursor == "" then
        curr = redis.call('GET', LAST_HASH_KEYWORD)
        if not curr then
            return { "", {} } -- chain is empty
        end
    else
        -- move one step back so we don't return the cursor itself again
        curr = redis.call('GET', PREV_BLOCK_HASH_KEYWORD .. cursor)
        if not curr or curr == "" then
            return { "", {} } -- reached genesis or invalid cursor
        end
    end

    local result = {}
    local fetched = 0
    while curr and fetched < count do
        table.insert(result, curr)
        fetched = fetched + 1
        curr = redis.call('GET', PREV_BLOCK_HASH_KEYWORD .. curr)
        if not curr or curr == "" then break end
    end

    local nextCursor = ""
    if #result > 0 then
        nextCursor = result[#result]
    end
    return { nextCursor, result }
end

-- REGISTER FUNCTIONS
redis.register_function('init_debug_chain', init_debug_chain)
redis.register_function('iterate_chain', iterate_chain)
-- Базовая конфигурация
box.cfg{
    listen = '0.0.0.0:3301',
    wal_mode = 'write',
    memtx_memory = 256 * 1024 * 1024,
    log_level = 5
}

-- Ждем готовности Tarantool (альтернатива для версии 2.10)
local function wait_ready()
    local attempts = 0
    while attempts < 20 do
        if box.info.status == 'running' then
            return true
        end
        require('fiber').sleep(0.1)
        attempts = attempts + 1
    end
    return false
end

if not wait_ready() then
    print("ERROR: Tarantool not ready after 2 seconds")
    os.exit(1)
end

-- Конфигурация
local USERNAME = os.getenv('TARANTOOL_USER_NAME')
local PASSWORD = os.getenv('TARANTOOL_USER_PASSWORD')

-- 1. Создаем пространство polls
print("Creating space 'polls'...")
local ok, err = pcall(function()
    if box.space.polls == nil then
        -- Создаем sequence для автоинкремента ID
        box.schema.sequence.create('polls_id', {
            if_not_exists = true
        })

        -- Создаем пространство polls с правильным форматом
        box.schema.space.create('polls', {
            if_not_exists = true,
            format = {
                {name = 'id', type = 'unsigned', is_nullable = false},
                {name = 'question', type = 'string'},
                {name = 'options', type = 'array'},
                {name = 'votes', type = 'map'},
                {name = 'end_date', type = 'unsigned'},
                {name = 'creator_id', type = 'unsigned'},
                {name = 'creator_name', type = 'string'},
                {name = 'voters', type = 'map'},
                {name = 'created_at', type = 'unsigned'}
            }
        })

        -- Создаем первичный индекс с sequence
        box.space.polls:create_index('primary', {
            type = 'tree',
            parts = {'id'},
            sequence = 'polls_id',
            if_not_exists = true
        })
        print("Space 'polls' created successfully")
    else
        print("Space 'polls' already exists")
    end
end)

if not ok then
    print("ERROR creating space: "..tostring(err))
    os.exit(1)
end

-- 2. Настройка пользователя
print("Configuring user...")
ok, err = pcall(function()
    if not box.schema.user.exists(USERNAME) then
        box.schema.user.create(USERNAME, {password = PASSWORD})
    end
    box.schema.user.grant(USERNAME, 'read,write,execute', 'space', 'polls')
    print("User configured successfully")
end)

if not ok then
    print("ERROR configuring user: "..tostring(err))
    os.exit(1)
end

-- 3. Healthcheck функция
rawset(_G, 'healthcheck', function()
    return {
        status = box.space.polls ~= nil and 'OK' or 'FAIL',
        timestamp = os.time()
    }
end)
local function setup_functions()
    -- Удаление опроса
    box.schema.func.create('delete_poll', {
        if_not_exists = true,
        body = [[function(poll_id, user_id)
            local poll = box.space.polls:get(poll_id)
            if not poll then return {error = 'PollNotFound'} end
            if poll.creator_id ~= user_id then return {error = 'UserNotCreator'} end
            box.space.polls:delete(poll_id)
            return {ok = true}
        end]]
    })

    -- Получение результатов
    box.schema.func.create('get_poll_results', {
        if_not_exists = true,
        body = [[function(poll_id)
            local poll = box.space.polls:get(poll_id)
            if not poll then return nil end

            return {
                id = poll[1],
                question = poll[2],
                options = poll[3],
                votes = poll[4] or {}, -- Важно: возвращаем пустой map если голосов нет
                end_date = poll[5],
                creator_id = poll[6],
                creator_name = poll[7],
                voters = poll[8] or {}
            }
        end]]
    })

    -- Досрочное завершение
    box.schema.func.create('complete_poll_early', {
        if_not_exists = true,
        body = [[function(poll_id, user_id)
            local poll = box.space.polls:get(poll_id)
            if not poll then
                return {error = 'PollNotFound'}
            end

            -- Проверяем creator_id (поле 6)
            if poll[6] ~= user_id then
                return {error = 'UserNotCreator'}
            end

            -- Устанавливаем текущее время как end_date (поле 5)
            local new_end_date = os.time()
            box.space.polls:update(poll_id, {{'=', 5, new_end_date}})

            -- Возвращаем только необходимые данные
            return {
                ok = true,
                end_date = new_end_date  -- Возвращаем ТОЛЬКО новую дату окончания
            }
        end]]
    })

    -- Создание голоса
    box.schema.func.create('create_vote', {
        if_not_exists = true,
        body = [[function(poll_id, option, user_id_str)
            local poll = box.space.polls:get(poll_id)
            if not poll then return {error = 'PollNotFound'} end

            -- Проверка, что опрос ещё активен
            if os.time() > poll[5] then return {error = 'PollAlreadyClosed'} end

            -- Проверка, что пользователь ещё не голосовал
            local voters = poll[8] or {}
            if voters[user_id_str] then return {error = 'UserAlreadyVoted'} end

            -- Проверка, что вариант ответа существует
            local options = poll[3] or {}
            local valid_option = false
            for _, opt in ipairs(options) do
                if opt == option then
                    valid_option = true
                    break
                end
            end
            if not valid_option then return {error = 'InvalidVoteOption'} end

            -- Обновляем голоса
            local votes = poll[4] or {}
            votes[option] = (votes[option] or 0) + 1
            voters[user_id_str] = true

            -- Сохраняем изменения
            box.space.polls:update(poll_id, {
                {'=', 4, votes},
                {'=', 8, voters}
            })

            return {ok = true, votes = votes}
        end]]
    })

    print("All functions created successfully")
    return true
end
print("Setting up functions...")
ok, err = pcall(setup_functions)  -- Вот этот вызов!
if not ok then
    print("ERROR setting up functions: "..tostring(err))
    os.exit(1)
end

-- Даем права на выполнение функций
box.schema.user.grant(USERNAME, 'execute', 'function', 'delete_poll')
box.schema.user.grant(USERNAME, 'execute', 'function', 'get_poll_results')
box.schema.user.grant(USERNAME, 'execute', 'function', 'complete_poll_early')
box.schema.user.grant(USERNAME, 'execute', 'function', 'create_vote')

print("Tarantool initialization COMPLETED successfully!")
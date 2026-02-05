#!/bin/sh
set -e

# Функция для логирования
log() {
    echo "[xTeVe Entrypoint] $1"
}

# Обработка шаблонных файлов с переменными в формате {XTEVE_*}
process_templates() {
    log "Processing template files..."
    
    # Находим все файлы с расширением .template
    find ${XTEVE_HOME} -name "*.template" -type f | while read file; do
        target="${file%.template}"
        log "Processing template: $file -> $target"
        
        # Заменяем {VAR} на $VAR для envsubst
        sed 's/\${\([A-Z_][A-Z0-9_]*\)}/\$\1/g; s/{\([A-Z_][A-Z0-9_]*\)}/\$\1/g' "$file" | envsubst > "$target"
        
        # Проверяем результат
        if [ -f "$target" ]; then
            log "Successfully processed: $target"
            # Удаляем шаблонный файл после успешной обработки
            rm "$file"
        else
            log "ERROR: Failed to process template: $file"
            exit 1
        fi
    done
    
    # Проверяем наличие settings.json и создаем его из шаблона, если необходимо
    if [ ! -f "${XTEVE_HOME}/settings.json" ] && [ -f "${XTEVE_HOME}/settings.json.template" ]; then
        log "Creating settings.json from template..."
        cp "${XTEVE_HOME}/settings.json.template" "${XTEVE_HOME}/settings.json"
    fi
    
    # Дополнительная обработка для файлов конфигурации, которые могут содержать {XTEVE_*}
    # но не имеют расширения .template
    if [ -f "${XTEVE_HOME}/settings.json" ]; then
        log "Checking settings.json for template variables..."
        # Создаем временную копию
        cp "${XTEVE_HOME}/settings.json" "${XTEVE_HOME}/settings.json.tmp"
        
        # Обрабатываем переменные в формате {XTEVE_*}
        sed 's/\${\([A-Z_][A-Z0-9_]*\)}/\$\1/g; s/{\([A-Z_][A-Z0-9_]*\)}/\$\1/g' "${XTEVE_HOME}/settings.json.tmp" | envsubst > "${XTEVE_HOME}/settings.json"
        
        # Удаляем временный файл
        rm "${XTEVE_HOME}/settings.json.tmp"
        log "Processed settings.json"
    fi
}

# Создание необходимых директорий, если они не существуют
create_directories() {
    log "Creating necessary directories..."
    
    # Создаем директории на основе переменных окружения
    mkdir -p "${XTEVE_HOME}"
    mkdir -p "${XTEVE_TEMP}"
    mkdir -p "${XTEVE_HOME}/backup"
    mkdir -p "${XTEVE_HOME}/data"
    mkdir -p "${XTEVE_HOME}/cache"
    mkdir -p "${XTEVE_HOME}/cache/images"
    
    # Устанавливаем правильные права доступа
    chown -R xteve:xteve "${XTEVE_HOME}" "${XTEVE_TEMP}" 2>/dev/null || true
    chmod -R 755 "${XTEVE_HOME}" "${XTEVE_TEMP}" 2>/dev/null || true
    
    log "Directories created and permissions set"
}

# Проверка переменных окружения
check_environment() {
    log "Checking environment variables..."
    
    # Проверяем обязательные переменные
    if [ -z "${XTEVE_HOME}" ]; then
        log "WARNING: XTEVE_HOME is not set, using default: /home/xteve"
        export XTEVE_HOME="/home/xteve"
    fi
    
    if [ -z "${XTEVE_PORT}" ]; then
        log "WARNING: XTEVE_PORT is not set, using default: 34400"
        export XTEVE_PORT="34400"
    fi
    
    if [ -z "${XTEVE_TEMP}" ]; then
        log "WARNING: XTEVE_TEMP is not set, using default: /tmp/xteve"
        export XTEVE_TEMP="/tmp/xteve"
    fi
    
    # Выводим информацию о конфигурации
    log "Configuration:"
    log "  XTEVE_HOME: ${XTEVE_HOME}"
    log "  XTEVE_PORT: ${XTEVE_PORT}"
    log "  XTEVE_TEMP: ${XTEVE_TEMP}"
    log "  TZ: ${TZ:-UTC}"
}

# Основная логика
main() {
    log "Starting xTeVe container initialization..."
    
    # Проверяем переменные окружения
    check_environment
    
    # Создаем необходимые директории
    create_directories
    
    # Обрабатываем шаблонные файлы
    process_templates
    
    log "Initialization completed successfully"
    
    # Запускаем основное приложение с переданными аргументами
    # Заменяем переменные в аргументах командной строки
    ARGS=""
    for arg in "$@"; do
        # Заменяем ${VAR} и {VAR} на реальные значения переменных окружения
        processed_arg=$(echo "$arg" | sed 's/\${\([A-Z_][A-Z0-9_]*\)}/\$\1/g; s/{\([A-Z_][A-Z0-9_]*\)}/\$\1/g' | envsubst)
        ARGS="$ARGS $processed_arg"
    done
    
    log "Starting xTeVe with arguments: $ARGS"
    exec /usr/local/bin/xteve $ARGS
}

# Запускаем основную функцию со всеми аргументами
main "$@"
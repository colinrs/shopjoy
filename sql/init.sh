#!/bin/bash
# ============================================
# ShopJoy 数据库初始化脚本
#
# 功能：
#   1. 自动连接到 MySQL 数据库
#   2. 执行所有 schema.sql 文件创建表
#   3. 插入测试数据
#
# 使用方法:
#   # 使用环境变量 (推荐)
#   export MYSQL_HOST=localhost
#   export MYSQL_PORT=3306
#   export MYSQL_USER=root
#   export MYSQL_PASSWORD=your_password
#   export MYSQL_DATABASE=shopjoy
#   ./init.sh
#
#   # 或使用命令行参数
#   ./init.sh -h localhost -P 3306 -u root -p your_password -d shopjoy
#
#   # 仅执行特定领域
#   ./init.sh --only=user,product
#
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-shopjoy}"
MYSQL_CHARSET="${MYSQL_CHARSET:-utf8mb4}"

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 执行顺序 (按依赖关系)
SCHEMA_ORDER=(
    "user"
    "product"
    "promotion"
    "order"
    "payment"
    "fulfillment"
    "storefront"
    "points"
    "shop"
    "review"
)

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助
show_help() {
    cat << EOF
ShopJoy 数据库初始化脚本

用法: $0 [选项]

选项:
    -h, --host HOST      MySQL 主机地址 (默认: localhost)
    -P, --port PORT      MySQL 端口 (默认: 3306)
    -u, --user USER      MySQL 用户名 (默认: root)
    -p, --password PASS  MySQL 密码 (默认: 环境变量 MYSQL_PASSWORD)
    -d, --database NAME  数据库名称 (默认: shopjoy)
    -c, --charset CHARSET 字符集 (默认: utf8mb4)
    --only MODULES       仅执行指定模块，用逗号分隔
    --skip MODULES       跳过指定模块，用逗号分隔
    --tables-only        仅创建表，不插入数据
    --data-only          仅插入数据，不创建表 (需表已存在)
    --all                创建表并插入数据 (默认)
    --dry-run            模拟运行，不执行任何操作
    --help               显示帮助信息

环境变量:
    MYSQL_HOST           MySQL 主机地址
    MYSQL_PORT           MySQL 端口
    MYSQL_USER           MySQL 用户名
    MYSQL_PASSWORD       MySQL 密码
    MYSQL_DATABASE       数据库名称

示例:
    $0                                    # 使用环境变量执行
    $0 -h localhost -u root -p pass -d shopjoy  # 使用命令行参数
    $0 --only=user,product                # 仅初始化 user 和 product
    $0 --skip=points,review               # 跳过 points 和 review
    $0 --tables-only                      # 仅创建表

EOF
}

# 解析命令行参数
PARSE_ARGS=()
SKIP_MODULES=""
ONLY_MODULES=""
TABLES_ONLY=false
DATA_ONLY=false
ALL_MODE=false
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--host)
            MYSQL_HOST="$2"
            shift 2
            ;;
        -P|--port)
            MYSQL_PORT="$2"
            shift 2
            ;;
        -u|--user)
            MYSQL_USER="$2"
            shift 2
            ;;
        -p|--password)
            MYSQL_PASSWORD="$2"
            shift 2
            ;;
        -d|--database)
            MYSQL_DATABASE="$2"
            shift 2
            ;;
        -c|--charset)
            MYSQL_CHARSET="$2"
            shift 2
            ;;
        --only)
            ONLY_MODULES="$2"
            shift 2
            ;;
        --skip)
            SKIP_MODULES="$2"
            shift 2
            ;;
        --tables-only)
            TABLES_ONLY=true
            ALL_MODE=false
            shift
            ;;
        --data-only)
            DATA_ONLY=true
            ALL_MODE=false
            shift
            ;;
        --all)
            ALL_MODE=true
            TABLES_ONLY=false
            DATA_ONLY=false
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            log_error "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

# 构建 MySQL 连接命令
build_mysql_cmd() {
    local cmd="mysql -h '${MYSQL_HOST}' -P '${MYSQL_PORT}' -u '${MYSQL_USER}'"
    if [[ -n "${MYSQL_PASSWORD}" ]]; then
        cmd="${cmd} -p'${MYSQL_PASSWORD}'"
    fi
    cmd="${cmd} '${MYSQL_DATABASE}'"
    echo "${cmd}"
}

# 执行 SQL 文件
execute_sql_file() {
    local file="$1"
    local description="$2"
    local mysql_cmd

    if [[ ! -f "${file}" ]]; then
        log_warning "文件不存在: ${file}"
        return 1
    fi

    log_info "执行: ${description}"
    log_info "  文件: ${file}"

    if [[ "${DRY_RUN}" == true ]]; then
        log_info "  [DRY-RUN] 跳过执行"
        return 0
    fi

    mysql_cmd=$(build_mysql_cmd)

    # 执行 SQL 文件，忽略警告
    if eval "${mysql_cmd}" < "${file}" 2>&1 | grep -v "Using a password on the command line interface" | grep -v "Warning"; then
        # 如果有错误输出
        local output
        output=$(eval "${mysql_cmd}" < "${file}" 2>&1)
        if echo "${output}" | grep -qi "error"; then
            log_error "执行失败: ${description}"
            echo "${output}" | head -20
            return 1
        fi
    fi

    log_success "完成: ${description}"
    return 0
}

# 执行 SQL 语句
execute_sql() {
    local sql="$1"
    local description="$2"
    local mysql_cmd

    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY-RUN] ${description}"
        log_info "  SQL: ${sql:0:100}..."
        return 0
    fi

    mysql_cmd=$(build_mysql_cmd)
    if eval "echo \"${sql}\" | ${mysql_cmd}" 2>&1 | grep -qi "error"; then
        log_error "执行失败: ${description}"
        return 1
    fi

    log_success "完成: ${description}"
    return 0
}

# 检查数据库连接
check_connection() {
    log_info "检查数据库连接..."

    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY-RUN] 跳过连接检查"
        return 0
    fi

    local mysql_cmd=$(build_mysql_cmd)
    if ! eval "${mysql_cmd}" -e "SELECT 1" >/dev/null 2>&1; then
        log_error "无法连接到 MySQL 数据库"
        log_error "请检查以下配置:"
        log_error "  主机: ${MYSQL_HOST}:${MYSQL_PORT}"
        log_error "  用户: ${MYSQL_USER}"
        log_error "  数据库: ${MYSQL_DATABASE}"
        return 1
    fi

    log_success "数据库连接成功"
    return 0
}

# 创建数据库
create_database() {
    log_info "创建数据库 (如果不存在): ${MYSQL_DATABASE}"

    if [[ "${DRY_RUN}" == true ]]; then
        log_info "[DRY-RUN] 跳过数据库创建"
        return 0
    fi

    local create_cmd="mysql -h '${MYSQL_HOST}' -P '${MYSQL_PORT}' -u '${MYSQL_USER}'"
    if [[ -n "${MYSQL_PASSWORD}" ]]; then
        create_cmd="${create_cmd} -p'${MYSQL_PASSWORD}'"
    fi

    ${create_cmd} -e "CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\` DEFAULT CHARACTER SET ${MYSQL_CHARSET} COLLATE ${MYSQL_CHARSET}_unicode_ci;" 2>/dev/null

    if [[ $? -eq 0 ]]; then
        log_success "数据库 ${MYSQL_DATABASE} 就绪"
    else
        log_error "创建数据库失败"
        return 1
    fi
}

# 获取要执行的模块列表
get_modules_to_execute() {
    local modules=()

    if [[ -n "${ONLY_MODULES}" ]]; then
        # 使用 --only 指定的模块
        IFS=',' read -ra ADDR <<< "${ONLY_MODULES}"
        for module in "${ADDR[@]}"; do
            # 去除空白
            module=$(echo "${module}" | xargs)
            modules+=("${module}")
        done
    else
        # 使用默认顺序，排除 --skip 的模块
        for module in "${SCHEMA_ORDER[@]}"; do
            local skip=false
            if [[ -n "${SKIP_MODULES}" ]]; then
                IFS=',' read -ra ADDR <<< "${SKIP_MODULES}"
                for skip_module in "${ADDR[@]}"; do
                    skip_module=$(echo "${skip_module}" | xargs)
                    if [[ "${module}" == "${skip_module}" ]]; then
                        skip=true
                        break
                    fi
                done
            fi
            if [[ "${skip}" == false ]]; then
                modules+=("${module}")
            fi
        done
    fi

    echo "${modules[@]}"
}

# 初始化指定模块
init_module() {
    local module="$1"
    local schema_file="${SCRIPT_DIR}/${module}/schema.sql"

    if [[ ! -f "${schema_file}" ]]; then
        log_warning "模块 ${module} 的 schema 文件不存在: ${schema_file}"
        return 1
    fi

    log_info ""
    log_info "=========================================="
    log_info "初始化模块: ${module}"
    log_info "=========================================="

    # --all 模式或默认模式: 创建表并插入数据
    if [[ "${ALL_MODE}" == true ]] || [[ "${TABLES_ONLY}" == false && "${DATA_ONLY}" == false ]]; then
        execute_sql_file "${schema_file}" "${module} - 创建表和插入测试数据"
    # --tables-only 模式: 仅创建表
    elif [[ "${TABLES_ONLY}" == true ]]; then
        log_info "[TABLES-ONLY 模式] 仅创建表，跳过数据插入"
        execute_sql_file "${schema_file}" "${module} - 仅创建表"
    # --data-only 模式: 仅插入数据 (需要表已存在)
    elif [[ "${DATA_ONLY}" == true ]]; then
        log_warning "[DATA-ONLY 模式] 需要表已存在，跳过表创建"
        log_info "[DATA-ONLY 模式] 注意: 当前 schema 文件中表结构和数据在同一文件中"
    fi

    return 0
}

# 主函数
main() {
    echo ""
    echo "=========================================="
    echo "ShopJoy 数据库初始化脚本"
    echo "=========================================="
    echo ""
    echo "配置信息:"
    echo "  主机: ${MYSQL_HOST}:${MYSQL_PORT}"
    echo "  用户: ${MYSQL_USER}"
    echo "  数据库: ${MYSQL_DATABASE}"
    echo "  字符集: ${MYSQL_CHARSET}"
    echo ""

    if [[ "${DRY_RUN}" == true ]]; then
        log_warning "[DRY-RUN 模式] 不会执行任何实际操作"
        echo ""
    fi

    # 检查连接
    check_connection || exit 1

    # 创建数据库
    create_database || exit 1

    # 获取要执行的模块
    modules=$(get_modules_to_execute)

    echo ""
    log_info "将执行以下模块:"
    for module in ${modules}; do
        echo "  - ${module}"
    done
    echo ""

    # 执行各个模块
    local failed=0
    for module in ${modules}; do
        if ! init_module "${module}"; then
            log_error "模块 ${module} 初始化失败"
            failed=1
        fi
    done

    echo ""
    if [[ ${failed} -eq 0 ]]; then
        log_success "=========================================="
        log_success "所有模块初始化完成!"
        log_success "=========================================="
        echo ""
        log_info "下一步:"
        log_info "  - 检查表是否创建成功: mysql -h ${MYSQL_HOST} -P ${MYSQL_PORT} -u ${MYSQL_USER} -p -e 'USE ${MYSQL_DATABASE}; SHOW TABLES;'"
        log_info "  - 查看数据: mysql -h ${MYSQL_HOST} -P ${MYSQL_PORT} -u ${MYSQL_USER} -p -e 'USE ${MYSQL_DATABASE}; SELECT COUNT(*) FROM users;'"
        echo ""
    else
        log_error "=========================================="
        log_error "部分模块初始化失败，请检查错误信息"
        log_error "=========================================="
        exit 1
    fi
}

# 脚本入口
main

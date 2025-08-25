#!/bin/bash
# 更新数据库表和列注释的脚本
# 使用方法: ./scripts/update-db-comments.sh

set -e

# 数据库连接参数
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_USER=${DB_USER:-"postgres"}
DB_NAME=${DB_NAME:-"proxy_enhancer_ultra"}

echo "🚀 开始更新数据库注释..."
echo "数据库: $DB_NAME@$DB_HOST:$DB_PORT"

# 检查是否存在SQL文件
SQL_FILE="internal/database/add_column_comments.sql"
if [ ! -f "$SQL_FILE" ]; then
    echo "❌ 错误: 找不到SQL文件 $SQL_FILE"
    exit 1
fi

# 执行SQL脚本
echo "📝 应用数据库注释..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$SQL_FILE"

if [ $? -eq 0 ]; then
    echo "✅ 数据库注释更新完成！"
    
    # 显示更新结果统计
    echo ""
    echo "📊 更新统计:"
    
    # 统计有注释的表数量
    TABLE_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM pg_tables JOIN pg_class ON pg_tables.tablename = pg_class.relname 
        WHERE schemaname = 'public' AND obj_description(oid, 'pg_class') IS NOT NULL;
    " | tr -d ' ')
    
    echo "   - 已添加注释的表: $TABLE_COUNT 个"
    
    # 统计有注释的列数量  
    COLUMN_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND col_description((SELECT oid FROM pg_class WHERE relname = table_name), ordinal_position) IS NOT NULL;
    " | tr -d ' ')
    
    echo "   - 已添加注释的列: $COLUMN_COUNT 个"
    
    echo ""
    echo "🎉 所有数据库字段和表都已添加中文注释！"
else
    echo "❌ 数据库注释更新失败！"
    exit 1
fi

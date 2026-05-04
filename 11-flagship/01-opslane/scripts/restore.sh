#!/bin/bash
# Database restore script for Opslane
# Usage: ./scripts/restore.sh <backup_file>

set -e

# Configuration
DB_NAME="${OPSLANE_DB_NAME:-opslane}"
DB_USER="${OPSLANE_DB_USER:-opslane}"
DB_HOST="${OPSLANE_DB_HOST:-localhost}"
DB_PORT="${OPSLANE_DB_PORT:-5432}"
BACKUP_DIR="${OPSLANE_BACKUP_DIR:-./backups}"

if [ -z "$1" ]; then
    echo "Error: Backup file required"
    echo "Usage: $0 <backup_file>"
    echo ""
    echo "Available backups:"
    ls -la "${BACKUP_DIR}"/*.sql 2>/dev/null || echo "  No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Expand path if relative
if [[ ! "$BACKUP_FILE" = /* ]]; then
    BACKUP_FILE="${BACKUP_DIR}/${BACKUP_FILE}"
fi

if [ ! -f "${BACKUP_FILE}" ]; then
    echo "Error: Backup file not found: ${BACKUP_FILE}"
    exit 1
fi

echo "WARNING: This will DROP all existing data and restore from backup!"
echo "Backup file: ${BACKUP_FILE}"
echo ""

read -p "Type 'yes' to continue: " CONFIRM
if [ "${CONFIRM}" != "yes" ]; then
    echo "Aborted."
    exit 0
fi

echo ""
echo "Restoring database: ${DB_NAME}"

# Check if we're in Docker
if [ -f /.dockerenv ] || grep -q docker /proc/1/cgroup 2>/dev/null; then
    # Drop and recreate database
    psql -h db -U "${DB_USER}" -c "DROP DATABASE IF EXISTS ${DB_NAME};"
    psql -h db -U "${DB_USER}" -c "CREATE DATABASE ${DB_NAME};"
    psql -h db -U "${DB_USER}" -d "${DB_NAME}" -f "${BACKUP_FILE}"
else
    # Drop and recreate database
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -c "DROP DATABASE IF EXISTS ${DB_NAME};"
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -c "CREATE DATABASE ${DB_NAME};"
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "${BACKUP_FILE}"
fi

echo ""
echo "✓ Database restored successfully"
# Export API - Go Version

A high-performance Go API for exporting large datasets from MySQL databases to Excel files. This application is designed to handle unlimited data exports without the limitations of the original TypeScript/Next.js version.

## Features

- **Unlimited Data Export**: No artificial limits on record counts
- **Chunked Processing**: Efficiently processes large datasets in 10k record chunks
- **Memory Optimized**: Uses streaming and chunked processing to handle massive datasets
- **Excel Generation**: Creates properly formatted Excel files with headers and data
- **Date Filtering**: Support for date range filtering
- **Pretty Formatting**: Option to export data in a user-friendly format with proper column ordering
- **Fault Detection**: Automatic detection and formatting of fault-related data
- **CORS Support**: Built-in CORS middleware for web applications
- **Health Checks**: Built-in health check endpoint for monitoring

## Architecture

- **Framework**: Gin (high-performance HTTP web framework)
- **Database**: MySQL with connection pooling
- **Excel**: Excelize library for Excel file generation
- **Containerization**: Docker with multi-stage builds
- **Deployment**: Ready for Render deployment

## Prerequisites

- Go 1.21 or higher
- MySQL database
- Docker (for containerization)

## Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd export-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USER=your_username
   export DB_PASSWORD=your_password
   export DB_NAME=your_database
   export PORT=8080
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

### Docker

1. **Build the image**
   ```bash
   docker build -t export-api .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 \
     -e DB_HOST=your_db_host \
     -e DB_PORT=3306 \
     -e DB_USER=your_username \
     -e DB_PASSWORD=your_password \
     -e DB_NAME=your_database \
     export-api
   ```

## API Endpoints

### Export Data
```
GET /export?table=<table_name>&fromDate=<YYYY-MM-DD>&toDate=<YYYY-MM-DD>&all=<true|false>&order=<asc|desc>
```

**Parameters:**
- `table` (required): Name of the table to export
- `fromDate` (optional): Start date for filtering (YYYY-MM-DD format)
- `toDate` (optional): End date for filtering (YYYY-MM-DD format)
- `all` (optional): Whether to use pretty formatting (default: true)
- `order` (optional): Sort order - "asc" or "desc" (default: "desc")

**Response:** Excel file download

### Health Check
```
GET /health
```

**Response:** JSON status response

### CORS Preflight
```
OPTIONS /export
```

## Supported Tables

The API supports the following tables:
- GTPL_108_gT_40E_P_S7_200_Germany
- GTPL_109_gT_40E_P_S7_200_Germany
- GTPL_110_gT_40E_P_S7_200_Germany
- GTPL_111_gT_80E_P_S7_200_Germany
- GTPL_112_gT_80E_P_S7_200_Germany
- GTPL_113_gT_80E_P_S7_200_Germany
- kabomachinedatasmart200
- GTPL_114_GT_140E_S7_1200
- GTPL_115_GT_180E_S7_1200
- GTPL_119_GT_180E_S7_1200
- GTPL_120_GT_180E_S7_1200
- GTPL_116_GT_240E_S7_1200
- GTPL_117_GT_320E_S7_1200
- GTPL_121_GT1000T
- gtpl_122_s7_1200_01
- GTPL_124_GT_450T_S7_1200
- GTPL_131_GT_650T_S7_1200
- GTPL_132_GT_650T_S7_1200

## Data Format

### Pretty Format (all=true)
- Separates date and time into separate columns
- Orders numeric columns based on predefined preferences
- Formats fault information into a readable format
- Applies proper Excel formatting

### Raw Format (all=false)
- Exports all data as-is from the database
- Maintains original column structure
- Minimal data transformation

## Performance Features

- **Chunked Processing**: Processes data in 10,000 record chunks to manage memory
- **Connection Pooling**: Efficient database connection management
- **Streaming**: Generates Excel files without loading entire dataset into memory
- **Progress Logging**: Logs processing progress for large exports

## Deployment on Render

1. **Connect your repository** to Render
2. **Create a new Web Service**
3. **Set the environment variables**:
   - `DB_HOST`: Your MySQL host
   - `DB_PORT`: MySQL port (usually 3306)
   - `DB_USER`: MySQL username
   - `DB_PASSWORD`: MySQL password
   - `DB_NAME`: Database name
4. **Deploy** - Render will automatically build and deploy your application

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | MySQL host address | localhost |
| `DB_PORT` | MySQL port | 3306 |
| `DB_USER` | MySQL username | root |
| `DB_PASSWORD` | MySQL password | (none) |
| `DB_NAME` | Database name | test |
| `PORT` | Application port | 8080 |

## Security Features

- **Input Validation**: All table names and parameters are validated
- **SQL Injection Protection**: Uses parameterized queries
- **CORS Configuration**: Configurable CORS settings
- **Non-root Container**: Runs as non-root user in Docker

## Monitoring and Health Checks

- Built-in health check endpoint at `/health`
- Progress logging for large exports
- Error logging with detailed error messages
- Docker health checks for container orchestration

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify database credentials and network connectivity
   - Check if MySQL server is running
   - Ensure database user has proper permissions

2. **Memory Issues with Large Exports**
   - The application is designed to handle large datasets
   - Monitor system memory usage
   - Consider increasing container memory limits if needed

3. **Export Timeout**
   - Large exports may take time depending on data size
   - Monitor application logs for progress
   - Consider implementing client-side progress tracking

### Logs

The application provides detailed logging:
- Database connection status
- Export progress for large datasets
- Error details and stack traces
- Request processing information

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue in the repository or contact the development team.
# golang_excel_download

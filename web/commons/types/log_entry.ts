export interface LogEntry {
  level: string;
  message: string;
  resourceId: string;
  timestamp: string;
  traceId: string;
  spanId: string;
  commit: string;
  metadata: {
    parentResourceId: string;
  };
}
export const sampleLogEntries: LogEntry[] = [
  {
    level: "error",
    message: "Failed to connect to DB",
    resourceId: "server-1234",
    timestamp: "2023-09-15T08:00:00Z",
    traceId: "abc-xyz-123",
    spanId: "span-456",
    commit: "5e5342f",
    metadata: {
      parentResourceId: "server-0987",
    },
  },
  {
    level: "info",
    message: "Application started",
    resourceId: "app-5678",
    timestamp: "2023-09-15T08:05:00Z",
    traceId: "def-uvw-456",
    spanId: "span-789",
    commit: "a1b2c3d",
    metadata: {
      parentResourceId: "app-9876",
    },
  },
  {
    level: "warning",
    message: "High CPU usage detected",
    resourceId: "server-9876",
    timestamp: "2023-09-15T08:10:00Z",
    traceId: "ghi-jkl-789",
    spanId: "span-012",
    commit: "f4e5d6c",
    metadata: {
      parentResourceId: "server-5432",
    },
  },
  {
    level: "debug",
    message: "Request received",
    resourceId: "api-1234",
    timestamp: "2023-09-15T08:15:00Z",
    traceId: "mno-pqr-012",
    spanId: "span-345",
    commit: "9876543",
    metadata: {
      parentResourceId: "api-5678",
    },
  },
  {
    level: "error",
    message: "Database connection lost",
    resourceId: "db-5678",
    timestamp: "2023-09-15T08:20:00Z",
    traceId: "stu-vwx-345",
    spanId: "span-678",
    commit: "1a2b3c4",
    metadata: {
      parentResourceId: "db-9012",
    },
  },
];

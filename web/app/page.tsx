"use client";
import { Navbar } from "@/components/navbar/navbar";
import { useRouter, useSearchParams } from "next/navigation";
import React, { useEffect, useRef } from "react";
import { Input } from "@/components/form/input";
import { Button } from "@/components/form/button";
import { LogEntry } from "@/commons/types/log_entry";
import Table from "@/components/table/logtable";
import axios from "axios";
import moment from "moment";
import { LogFilters } from "@/components/filters/logFilters";
import { FilterOptions } from "@/commons/types/searchFilter";
import { isEmptyObject } from "@/utils/object";

interface SearchFilters {
  q?: string;
  level?: string;
  message?: string;
  resourceId?: string;
  traceId?: string;
  spanId?: string;
  commit?: string;
  parentResourceId?: string;
  from?: string;
  to?: string;
  limit?: number;
  page?: number;
}

interface SearchResults {
  logEntries: LogEntry[];
  total: number;
}

const SearchPage: React.FC = () => {
  const searchParams = useSearchParams();
  const router = useRouter();

  const [searchFilters, setSearchFilters] = React.useState<SearchFilters>({
    q: "",
    limit: 20,
    page: 1,
  });

  // const [logEntries, setLogEntries] = React.useState<LogEntry[]>([]);
  const [searchResults, setSearchResults] = React.useState<SearchResults>({
    logEntries: [],
    total: 0,
  });

  const [showFilters, setShowFilters] = React.useState<boolean>(false);

  useEffect(() => {
    const sQuery = searchParams.get("q");
    if (sQuery) {
      setSearchFilters({
        q: sQuery,
        ...searchFilters,
      });
    }
  }, [searchParams]);
  
  let timeout = useRef<NodeJS.Timeout>()
  
  useEffect(() => {
    if (timeout.current) {
      clearTimeout(timeout.current);
    }
    timeout.current = setTimeout(() => {
      handleSubmit();
    }, 500);
  }, [searchFilters]);

  const handleFilterSubmit = (filters: FilterOptions) => {
    if (!isEmptyObject(filters)) {
      setSearchFilters({
        ...searchFilters,
        ...filters,
      });
    } else {
      setSearchFilters({
        q: "",
        limit: 20,
        page: 1,
        from: '',
        to: '',
        level: "",
        message: "",
        resourceId: "",
        traceId: "",
        spanId: "",
        commit: "",
        parentResourceId: "",
      });
    }
  };

  // console.log(searchFilters)
  const handleSubmit = () => {
    axios
      .get(`${process.env.NEXT_PUBLIC_LOGSWIFT_API_URL}/query`, {
        params: {
          ...searchFilters,
          from:
            moment(searchFilters?.from, "YYYY/MM/DD").format(
              "YYYY-MM-DDTHH:mm:ss[Z]"
            ) === "Invalid date"
              ? undefined
              : moment(searchFilters?.from, "YYYY/MM/DD").format(
                  "YYYY-MM-DDTHH:mm:ss[Z]"
                ),
          to:
            moment(searchFilters?.to, "YYYY/MM/DD").format(
              "YYYY-MM-DDTHH:mm:ss[Z]"
            ) === "Invalid date"
              ? undefined
              : moment(searchFilters?.to, "YYYY/MM/DD").format(
                  "YYYY-MM-DDTHH:mm:ss[Z]"
                ),
        },
      })
      .then((res) => {
        if (res.data?.data?.logs) {
          setSearchResults({
            logEntries: res.data.data.logs,
            total: res.data.data.total,
          });
        } else {
          setSearchResults({
            logEntries: [],
            total: 0,
          });
        }
      })
      .catch((err) => {
        alert("error in getting values");
      });
  };

  return (
    <main className="flex h-screen flex-col w-full">
      <Navbar />
      <div className="w-full flex flex-col items-center gap-2">
        <div className="w-4/5 mt-4 flex flex-col gap-2">
          <div className="flex items-center gap-2">
            <Input
              value={searchFilters.q}
              onChange={(e) => {

                setSearchFilters({
                  q: e.target.value,
                });
              }}
              handleEnter={(query) => {
                router.push(`/search?q=${query}`);
              }}
            />
            <Button
              text="Search"
              onClick={() => {
                handleSubmit();
              }}
            />
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-500">
              Showing {searchResults?.logEntries?.length} of{" "}
              {searchResults.total} results found
            </span>
            <span
              className="text-sm text-blue-700 cursor-pointer hover:underline"
              onClick={() => {
                setShowFilters(!showFilters);
              }}
            >
              {showFilters ? "Hide" : "Show"} Filters
            </span>
          </div>
          {showFilters && (
            <LogFilters
              initialFilters={searchFilters}
              onFilterSubmit={(filters) => {
                handleFilterSubmit(filters);
              }}
            />
          )}
          <div className="overflow-y-scroll h-[76vh] no-scrollbar border-t border-b">
            <Table
              data={searchResults?.logEntries.map((logEntry) => ({
                ...logEntry,
                timestamp: moment(logEntry.timestamp).format("ll"),
              }))}
              columns={[
                {
                  header: "Level",
                  accessor: "level",
                  width: "w-1/12",
                  usetag: true,
                },
                { header: "Message", accessor: "message", width: "w-3/12" },
                {
                  header: "Resource ID",
                  accessor: "resourceId",
                  width: "w-2/12",
                },
                { header: "Timestamp", accessor: "timestamp", width: "w-2/12" },
                { header: "Trace ID", accessor: "traceId", width: "w-2/12" },
                { header: "Span ID", accessor: "spanId", width: "w-1/12" },
                { header: "Commit", accessor: "commit", width: "w-1/12" },
                {
                  header: "Parent Resource ID",
                  accessor: "metadata.parentResourceId",
                  width: "w-2/12",
                },
              ]}
            />
          </div>
          {/* basic pagination */}
          <div className="flex items-center justify-end gap-2 my-2">
            {(searchFilters.page || 0) > 1 && (
              <span
                className="text-blue-500 hover:underline cursor-pointer text-sm"
                onClick={() => {
                  setSearchFilters({
                    ...searchFilters,
                    page: (searchFilters.page || 1) - 1,
                  });
                }}
              >
                Previous
              </span>
            )}
            {searchResults.total > (searchFilters.limit || 20) && (
              <span
                className="text-blue-500 hover:underline cursor-pointer text-sm"
                onClick={() => {
                  setSearchFilters({
                    ...searchFilters,
                    page: (searchFilters.page || 1) + 1,
                  });
                }}
              >
                Next
              </span>
            )}
          </div>
        </div>
      </div>
    </main>
  );
};

export default SearchPage;

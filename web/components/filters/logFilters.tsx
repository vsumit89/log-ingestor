import React, { useEffect, useState } from "react";
import { Input } from "../form/input";
import { FilterOptions } from "../../commons/types/searchFilter";

interface LogFiltersProps {
  initialFilters?: FilterOptions;
  onFilterSubmit: (filters: FilterOptions) => void;
}

export const LogFilters: React.FC<LogFiltersProps> = ({
  onFilterSubmit,
  initialFilters,
}) => {
  const [filters, setFilters] = useState<FilterOptions>({});

  useEffect(() => {
    setFilters(initialFilters || {});
  }, [initialFilters]);

  const handleFilterSubmit = () => {
    onFilterSubmit(filters);
  };

  const handleFilterClear = () => {
    setFilters(() => {
      onFilterSubmit({});
      return {};
    });
  };

  console.log(filters);
  return (
    <div className="grid grid-cols-5 gap-2">
      <Input
        label="level"
        placeholder="Enter level to filter"
        value={filters.level}
        onChange={(e) => {
          setFilters({
            ...filters,
            level: e.target.value,
          });
        }}
      />
      <Input
        label="message"
        placeholder="Enter message to filter"
        value={filters.message}
        onChange={(e) => {
          setFilters({
            ...filters,
            message: e.target.value,
          });
        }}
      />
      <Input
        label="resourceId"
        placeholder="Enter resourceId to filter"
        value={filters.resourceId}
        onChange={(e) => {
          setFilters({
            ...filters,
            resourceId: e.target.value,
          });
        }}
      />
      <Input
        label="traceId"
        placeholder="Enter traceId to filter"
        value={filters.traceId}
        onChange={(e) => {
          setFilters({
            ...filters,
            traceId: e.target.value,
          });
        }}
      />
      <Input
        label="spanId"
        placeholder="Enter spanId to filter"
        value={filters.spanId}
        onChange={(e) => {
          setFilters({
            ...filters,
            spanId: e.target.value,
          });
        }}
      />
      <Input
        label="commit"
        placeholder="Enter commit to filter"
        value={filters.commit}
        onChange={(e) => {
          setFilters({
            ...filters,
            commit: e.target.value,
          });
        }}
      />
      <Input
        label="parentResourceId"
        placeholder="Enter parentResourceId to filter"
        value={filters.parentResourceId}
        onChange={(e) => {
          setFilters({
            ...filters,
            parentResourceId: e.target.value,
          });
        }}
      />
      <Input
        label="from date"
        placeholder="Enter from date to filter"
        type="date"
        value={filters.from}
        onChange={(e) => {
          setFilters({
            ...filters,
            from: e.target.value,
          });
        }}
      />
      <Input
        label="to date"
        placeholder="Enter to date to filter"
        type="date"
        value={filters.to}
        onChange={(e) => {
          setFilters({
            ...filters,
            to: e.target.value,
          });
        }}
      />
      <Input
        label="number of rows"
        placeholder="Enter level to filter"
        type="select"
        onSelected={(value) => {
          let limitInt = parseInt(value);
          setFilters({
            ...filters,
            limit: limitInt,
          });
        }}
        value={filters?.limit?.toString()}
        options={[
          {
            label: "20",
            value: "20",
          },
          {
            label: "50",
            value: "50",
          },
          {
            label: "100",
            value: "100",
          },
        ]}
      />
      <div className="flex items-center gap-2">
        <button
          className="bg-blue-500 text-white p-1 rounded-md hover:bg-blue-600 focus:outline-none focus:shadow-outline-blue active:bg-blue-800 w-1/2 text-sm"
          onClick={() => handleFilterSubmit()}
        >
          Apply Filters
        </button>
        <button
          className="bg-gray-300 text-gray-700 p-1 rounded-md hover:bg-gray-400 focus:outline-none focus:shadow-outline-gray active:bg-gray-500 w-1/2 text-sm"
          onClick={() => handleFilterClear()}
        >
          Clear Filters
        </button>
      </div>
      {/* small button with apply filter */}
    </div>
  );
};

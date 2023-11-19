import React from "react";
import Tag from "../tag/tag";

interface TableProps {
  columns: {
    header: string;
    accessor: string;
    width: string;
    usetag?: boolean;
  }[];
  data: any[];
}

function Table({ columns, data }: TableProps) {
  return (
    <div>
      <table className="min-w-full bg-white border border-gray-300">
        <thead>
          <tr>
            {columns.map((column) => (
              <th
                key={column.header}
                className={`py-2 px-4 border-b border-gray-300 bg-gray-100 sticky top-0 text-sm text-gray-500 ${
                  column.width || ""
                }`}
              >
                {column.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => (
            <tr key={rowIndex}>
              {columns.map((column, colIndex) => (
                <td
                  key={colIndex}
                  className={`py-2 px-4 border-b border-gray-300 ${
                    colIndex === columns?.length - 1 ? "pr-4" : ""
                  } text-sm text-gray-700 ${column.width || ""}} text-center`}
                >
                  {column.usetag ? (
                    <Tag text={row[column?.accessor]} />
                  ) : column?.accessor?.includes(".") ? (
                    getNestedPropertyValue(row, column?.accessor)
                  ) : (
                    row[column?.accessor]
                  )}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function getNestedPropertyValue(obj: any, accessor: string): any {
  const keys = accessor.split(".");
  let value = obj;

  for (const key of keys) {
    value = value?.[key];
    if (value === undefined) {
      break;
    }
  }

  return value;
}

export default Table;

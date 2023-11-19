import moment from "moment";
import React from "react";
import Select from "react-select";

type InputType = "text" | "date" | "select";
interface InputProps {
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string;
  value?: string;
  fullWidth?: boolean;
  width?: string;
  handleEnter?: (query: string) => void;
  label?: string;
  type?: InputType;
  options?: any[];
  onSelected?: (value: string) => void;
}

export const Input: React.FC<InputProps> = ({
  onChange,
  placeholder = "Enter your search query",
  value,
  fullWidth = true,
  width,
  handleEnter,
  label,
  type,
  options,
  onSelected,
}) => {
  const renderInput = () => {
    switch (type) {
      case "select":
        return (
          <Select
            className="w-full"
            options={options}
            onChange={(e) => {
              if (onSelected) {
                onSelected(e.value);
              }
            }}
            value={options?.find((option) => option.value === value)}
          />
        );
      case "date":
        return (
          <input
            type="date"
            className={`py-1 px-2 border outline-none w-full focus:border-[#8478E2]  text-base rounded-sm placeholder:text-gray-200`}
            onChange={onChange}
            placeholder={placeholder}
            value={value}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                // @ts-ignore
                handleEnter(value);
              }
            }}
          ></input>
        );
      default:
        return (
          <input
            className={`py-1 px-2 border outline-none w-full focus:border-[#8478E2]  text-base rounded-sm placeholder:text-gray-200`}
            onChange={onChange}
            placeholder={placeholder}
            value={value}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                // @ts-ignore
                handleEnter(value);
              }
            }}
          ></input>
        );
    }
  };

  return (
    <div className={`${fullWidth ? "w-full" : width} flex flex-col gap-1`}>
      {label && <label className="text-sm text-gray-500">{label}</label>}
      {renderInput()}
      {/* error text here */}
    </div>
  );
};

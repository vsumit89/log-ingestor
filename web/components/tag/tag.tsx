import React from "react";

interface TagProps {
  text: string;
}

const Tag: React.FC<TagProps> = ({ text }) => {
  // Define styles based on the log text
  const gettextStyles = () => {
    switch (text.toLowerCase()) {
      case "info":
        return "bg-blue-200 text-blue-800";
      case "warn":
        return "bg-yellow-200 text-yellow-800";
      case "error":
        return "bg-red-200 text-red-800";
      case "debug":
        return "bg-green-200 text-green-800";
      default:
        return "bg-gray-200 text-gray-800";
    }
  };

  return (
    <span className={`inline-block px-2 py-1 rounded ${gettextStyles()}`}>
      {text}
    </span>
  );
};

export default Tag;
